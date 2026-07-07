package mq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"

	"news-service/infra/models/domain"
)

const (
	ExchangeEvents = "hms.events"

	RoutingNewsCreated = "news.created"

	QueueNewsNotifications = "notifications.news"
)

type Publisher struct {
	conn    *amqp.Connection
	channel publisherChannel
	logger  *zerolog.Logger
}

type publisherChannel interface {
	PublishWithContext(context.Context, string, string, bool, bool, amqp.Publishing) error
	Close() error
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

func (p *Publisher) PublishNewsCreated(ctx context.Context, evt domain.NewsCreatedEvent) error {
	body, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	err = p.channel.PublishWithContext(
		ctx,
		ExchangeEvents,
		RoutingNewsCreated,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("publish news.created: %w", err)
	}

	p.logger.Info().
		Str("news_id", evt.NewsID).
		Str("house_id", evt.HouseID).
		Msg("event news.created published")

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
