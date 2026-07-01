package mq

import (
	"context"
	"encoding/json"
	"fmt"

	"requests-service/infra/models/domain"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

const (
	ExchangeEvents = "hms.events"

	RoutingRequestCreated       = "request.created"
	RoutingRequestStatusUpdated = "request.status_updated"
	RoutingRequestCommentAdded  = "request.comment_added"
)

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	logger  *zerolog.Logger
}

func New(dsn string, logger *zerolog.Logger) (*Publisher, error) {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("rabbitmq channel: %w", err)
	}

	if err := ch.ExchangeDeclare(
		ExchangeEvents,
		amqp.ExchangeTopic,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("declare exchange: %w", err)
	}

	logger.Info().Str("exchange", ExchangeEvents).Msg("rabbitmq publisher ready")

	return &Publisher{conn: conn, channel: ch, logger: logger}, nil
}

func (p *Publisher) PublishRequestCreated(ctx context.Context, evt domain.RequestEvent) error {
	return p.publish(ctx, RoutingRequestCreated, evt)
}

func (p *Publisher) PublishRequestStatusUpdated(ctx context.Context, evt domain.RequestEvent) error {
	return p.publish(ctx, RoutingRequestStatusUpdated, evt)
}

func (p *Publisher) PublishRequestCommentAdded(ctx context.Context, evt domain.RequestEvent) error {
	return p.publish(ctx, RoutingRequestCommentAdded, evt)
}

func (p *Publisher) publish(ctx context.Context, routingKey string, evt domain.RequestEvent) error {
	body, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	if err := p.channel.PublishWithContext(
		ctx,
		ExchangeEvents,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	); err != nil {
		return fmt.Errorf("publish %s: %w", routingKey, err)
	}

	p.logger.Info().Str("routing_key", routingKey).Str("request_id", evt.RequestID).Msg("event published")

	return nil
}

func (p *Publisher) Close() error {
	if p.channel != nil {
		_ = p.channel.Close()
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}
