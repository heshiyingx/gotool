package simple

import (
	"context"
	_const "github.com/heshiyingx/gotool/mqext/rabbitmqext/const"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	cfg     *ConsumerConfig
}
type ConsumerConfig struct {
	Tag           string
	PrefetchCount int  // 限制每次从队列中获取的消息数量
	PrefetchSize  int  // 这是预取大小，它限制了服务器将未确认的消息发送给消费者的总体积。然而，需要注意的是，这个参数在RabbitMQ中并未实现，所以通常设置为
	AutoAsk       bool // 是否自动确认
	Exclusive     bool // 是否排他
	NoLocal       bool // 换句话说，如果你在一个连接上既发布消息又消费消息，设置 NoLocal 为 true 可以防止消费者接收到自己发布的消息。
	NoWait        bool // 如果noWait参数为false，那么RabbitMQ服务器在创建队列后会向客户端发送一个确认消息。客户端会等待这个确认消息，然后再继续执行后续的操作。这种方式可以确保队列已经成功创建，但是会增加一些延迟。
	Args          amqp.Table
}
type Option func(*ConsumerConfig)
type ConsumeFunc func(msg amqp.Delivery) error

func NewConsume(conn *amqp.Connection, opts ...Option) (*Consumer, error) {

	config := &ConsumerConfig{
		//PrefetchCount: prefetchCount,
	}
	for _, opt := range opts {
		opt(config)
	}
	if config.PrefetchCount == 0 {
		config.PrefetchCount = 10
	}
	if config.Tag == "" {
		return nil, _const.RabbitConsumerConfigErr.Wrap("tag is empty")
	}
	consumer := &Consumer{
		conn: conn,
		cfg:  config,
	}
	channel, err := consumer.conn.Channel()
	if err != nil {
		return nil, err
	}
	consumer.channel = channel
	err = consumer.channel.Qos(config.PrefetchCount, config.PrefetchSize, false)
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
		qName,           // name
		c.cfg.Tag,       // consumerTag,
		c.cfg.AutoAsk,   // autoAck
		c.cfg.Exclusive, // exclusive
		c.cfg.NoLocal,   // noLocal
		c.cfg.NoWait,    // noWait
		c.cfg.Args,      // arguments
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
