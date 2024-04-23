package simple

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}
type ConsumerConfig struct {
	PrefetchCount int
	PrefetchSize  int
}
type Option func(*ConsumerConfig)
type ConsumeFunc func(msg amqp.Delivery) error

func NewConsume(conn *amqp.Connection, prefetchCount int, opts ...Option) (*Consumer, error) {
	consumer := &Consumer{
		conn: conn,
	}
	if prefetchCount == 0 {
		prefetchCount = 10
	}
	config := &ConsumerConfig{
		PrefetchCount: prefetchCount,
	}
	for _, opt := range opts {
		opt(config)
	}
	channel, err := consumer.conn.Channel()
	if err != nil {
		return nil, err
	}
	consumer.channel = channel
	err = consumer.channel.Qos(config.PrefetchCount, 0, false)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}
func (c *Consumer) close() {
	c.channel.Close()
	c.conn.Close()
}
func (c *Consumer) SampleConsumeWithQName(ctx context.Context, qName string, f ConsumeFunc) error {
	defer c.close()
	deliveries, err := c.channel.Consume(
		qName,    // name
		"simple", // consumerTag,
		false,    // autoAck
		true,     // exclusive
		false,    // noLocal
		false,    // noWait
		nil,      // arguments
	)
	if err != nil {
		return err
	}
	for delivery := range deliveries {
		handleErr := f(delivery)
		if handleErr != nil {
			c.channel.Nack(delivery.DeliveryTag, false, true)
			continue
		}

		c.channel.Ack(delivery.DeliveryTag, false)
	}

	return nil

}
