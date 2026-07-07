package mq

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"requests-service/infra/models/domain"
	"testing"
)

func TestPublisher(t *testing.T) {
	l := zerolog.Nop()
	_, e := New("invalid://", &l)
	require.Error(t, e)
	c := gomock.NewController(t)
	ch := NewMockpublisherChannel(c)
	p := &Publisher{channel: ch, logger: &l}
	evt := domain.RequestEvent{Type: "request.created", RequestID: "r"}
	for _, call := range []struct {
		key string
		run func() error
	}{{RoutingRequestCreated, func() error { return p.PublishRequestCreated(context.Background(), evt) }}, {RoutingRequestStatusUpdated, func() error { return p.PublishRequestStatusUpdated(context.Background(), evt) }}, {RoutingRequestCommentAdded, func() error { return p.PublishRequestCommentAdded(context.Background(), evt) }}} {
		ch.EXPECT().PublishWithContext(gomock.Any(), ExchangeEvents, call.key, false, false, gomock.Any()).DoAndReturn(func(_ context.Context, _, _ string, _, _ bool, msg amqp.Publishing) error {
			assert.Contains(t, string(msg.Body), `"request_id":"r"`)
			return nil
		})
		require.NoError(t, call.run())
	}
	ch.EXPECT().PublishWithContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mq"))
	e = p.PublishRequestCreated(context.Background(), evt)
	assert.EqualError(t, e, "publish request.created: mq")
	ch.EXPECT().Close().Return(nil)
	require.NoError(t, p.Close())
	require.NoError(t, (&Publisher{}).Close())
}
