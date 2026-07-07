package mq

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"

	infra_errors "notification-service/infra/errors"
	"notification-service/infra/models/domain"
)

const exchangeEvents = "hms.events"

type Handler interface {
	Notify(context.Context, *domain.Notification) error
}
type requestEvent struct {
	UserID    string `json:"user_id"`
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
}
type newsEvent struct {
	HouseID string `json:"house_id"`
	NewsID  string `json:"news_id"`
	Title   string `json:"title"`
}

type Consumer struct {
	dsn     string
	handler Handler
	logger  *zerolog.Logger
}

func New(dsn string, handler Handler, logger *zerolog.Logger) *Consumer {
	return &Consumer{dsn: dsn, handler: handler, logger: logger}
}

func (c *Consumer) Run(ctx context.Context) {
	for {
		if err := c.consume(ctx); err != nil && ctx.Err() == nil {
			c.logger.Error().Err(err).Msg("notification consumer disconnected")
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(2 * time.Second):
		}
	}
}

func (c *Consumer) consume(ctx context.Context) error {
	conn, err := amqp.Dial(c.dsn)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare(exchangeEvents, amqp.ExchangeTopic, true, false, false, false, nil); err != nil {
		return err
	}

	queue, err := ch.QueueDeclare("notifications.events", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for _, key := range []string{"news.created", "request.status_updated"} {
		if err := ch.QueueBind(queue.Name, key, exchangeEvents, false, nil); err != nil {
			return err
		}
	}

	deliveries, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	c.logger.Info().Msg("notification RabbitMQ consumer ready")

	for {
		select {
		case <-ctx.Done():
			return nil
		case delivery, ok := <-deliveries:
			if !ok {
				return infra_errors.ErrDeliveryChannelClosed
			}

			var notification *domain.Notification

			switch delivery.RoutingKey {
			case "request.status_updated":
				var event requestEvent
				if err := json.Unmarshal(delivery.Body, &event); err == nil {
					notification = requestEventToNotification(event)
				}
			case "news.created":
				var event newsEvent
				if err := json.Unmarshal(delivery.Body, &event); err == nil {
					notification = newsEventToNotification(event)
				}
			}

			if notification == nil {
				_ = delivery.Nack(false, false)

				continue
			}

			if err := c.handler.Notify(ctx, notification); err != nil {
				c.logger.Error().Err(err).Msg("store notification failed")

				_ = delivery.Nack(false, true)

				continue
			}

			_ = delivery.Ack(false)
		}
	}
}
