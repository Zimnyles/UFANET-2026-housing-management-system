package mq

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"news-service/infra/models/domain"
)

func TestPublisher(t *testing.T) {
	logger := zerolog.Nop()
	_, err := New("invalid://", &logger)
	require.Error(t, err)
	assert.Equal(t, "hms.events", ExchangeEvents)
	assert.Equal(t, "news.created", RoutingNewsCreated)
	assert.Equal(t, "notifications.news", QueueNewsNotifications)

	ctrl := gomock.NewController(t)
	ch := NewMockpublisherChannel(ctrl)
	p := &Publisher{channel: ch, logger: &logger}
	evt := domain.NewsCreatedEvent{Type: "news.created", NewsID: "n", Title: "title", HouseID: "h"}
	ch.EXPECT().PublishWithContext(context.Background(), ExchangeEvents, RoutingNewsCreated, false, false, gomock.Any()).DoAndReturn(
		func(_ context.Context, _, _ string, _, _ bool, msg amqp.Publishing) error {
			assert.JSONEq(t, `{"type":"news.created","news_id":"n","title":"title","house_id":"h"}`, string(msg.Body))
			assert.Equal(t, "application/json", msg.ContentType)
			return nil
		})
	require.NoError(t, p.PublishNewsCreated(context.Background(), evt))
	ch.EXPECT().PublishWithContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mq"))
	err = p.PublishNewsCreated(context.Background(), evt)
	assert.EqualError(t, err, "publish news.created: mq")
	ch.EXPECT().Close().Return(nil)
	require.NoError(t, p.Close())
	require.NoError(t, (&Publisher{logger: &logger}).Close())
}
