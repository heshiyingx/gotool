package simple

import (
	"context"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
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
	qName         string
}
type Option func(*ConsumerConfig)
type ConsumeFunc func(msg amqp.Delivery) error

// ConsumeBatchFunc 批量消费,返回成功的消息tags和失败的消息tags
type ConsumeBatchFunc func(msgs []amqp.Delivery) ([]uint64, []uint64)
type AckFunc func(msg amqp.Delivery, err error) error

func MustNewConsume(url string, vHost string, qName string, opts ...Option) *Consumer {
	config := amqp.Config{Properties: amqp.NewConnectionProperties(), Vhost: vHost}
	config.Properties.SetClientConnectionName("sample-consumer")

	conn, err := amqp.DialConfig(url, config)
	if err != nil {
		log.Fatalf("consumer: error in dial: %s", err)
	}
	consumer, err := NewConsume(conn, qName, opts...)
	if err != nil {
		log.Fatalf("consumer: error in dial: %s", err)
	}
	return consumer
}
func NewConsume(conn *amqp.Connection, qName string, opts ...Option) (*Consumer, error) {

	config := &ConsumerConfig{
		Tag:           "",
		PrefetchCount: 0,
		PrefetchSize:  0,
		AutoAsk:       false,
		Exclusive:     false,
		NoLocal:       false,
		NoWait:        false,
		Args:          nil,
		qName:         qName,
	}
	for _, opt := range opts {
		opt(config)
	}
	//if config.Tag == "" {
	//	return nil, _const.RabbitConsumerConfigErr.Wrap("tag is empty")
	//}
	consumer := &Consumer{
		conn: conn,
		cfg:  config,
	}
	channel, err := consumer.conn.Channel()
	if err != nil {
		return nil, err
	}
	queue, err := channel.QueueDeclare(
		qName, // name of the queue
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}
	err = channel.QueueBind(queue.Name, qName, "", false, nil)
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
func (c *Consumer) SampleConsumeWithFinishAckQName(ctx context.Context, f ConsumeFunc) error {
	defer c.close()

	deliveries, err := c.channel.Consume(
		c.cfg.qName,     // name
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
			if !c.cfg.AutoAsk {
				c.channel.Nack(delivery.DeliveryTag, false, true)
			}
			continue
		}
		if !c.cfg.AutoAsk {
			c.channel.Ack(delivery.DeliveryTag, false)
		}

	}

	return nil
}
func (c *Consumer) SampleConsumeWithQName(ctx context.Context, f ConsumeFunc) error {
	defer c.close()

	deliveries, err := c.channel.Consume(
		c.cfg.qName,     // name
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
		_ = f(delivery)
	}

	return nil
}
func (c *Consumer) SampleBatchConsumeWithQName(ctx context.Context, during time.Duration, batchNum int, f ConsumeBatchFunc) error {
	defer c.close()
	err := c.channel.Qos(batchNum, 0, false)
	if err != nil {
		return err
	}
	deliveries, err := c.channel.Consume(
		c.cfg.qName,     // name
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
	ticker := time.NewTicker(during)

	msgs := make([]amqp.Delivery, 0, batchNum)
	for {
		select {
		case msg, ok := <-deliveries:
			if !ok {
				return errors.New("channel closed")
			}
			msgs = append(msgs, msg)
			if len(msgs) >= batchNum {
				successTags, failTags := f(msgs)
				if len(successTags)+len(failTags) != len(msgs) {
					log.Fatal("successTags and failTags length not equal msgs length")
				}
				for _, tag := range successTags {
					c.channel.Ack(tag, false)
				}
				for _, tag := range failTags {
					c.channel.Nack(tag, false, true)
				}
				msgs = msgs[:0]
			}
		case <-ticker.C:
			if len(msgs) > 0 {
				successTags, failTags := f(msgs)
				if len(successTags)+len(failTags) != len(msgs) {
					log.Fatal("successTags and failTags length not equal msgs length")
				}
				for _, tag := range successTags {
					c.channel.Ack(tag, false)
				}
				for _, tag := range failTags {
					c.channel.Nack(tag, false, true)
				}
				msgs = msgs[:0]
			}

		}
	}
}
func (c *Consumer) Ack(deliveryTag uint64, success bool, multipe bool, requeue bool) {
	if success {
		log.Println(deliveryTag)
		err := c.channel.Ack(deliveryTag, multipe)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		c.channel.Nack(deliveryTag, multipe, requeue)
	}
}
