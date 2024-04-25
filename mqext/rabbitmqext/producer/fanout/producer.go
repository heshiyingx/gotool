package fanout

import (
	"context"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"strconv"
	"time"
)

type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	cfg     *ProducerConfig
}
type ProducerConfig struct {
	Durable       bool
	AutoDelete    bool // 当最后一个消费者断开连接之后，队列是否自动删除
	Exclusive     bool
	NoWait        bool
	args          amqp.Table
	AppId         string
	DeliveryMode  uint8
	qNames        []string
	Mandatory     bool // 如果为true，当消息无法路由到队列时，会通过basic.return方法将消息返回给生产者,如果mandatory参数为false，那么服务器会直接丢弃无法路由的消息。
	Immediate     bool // 如果为true，当消息无法路由到队列时，会通过basic.return方法将消息返回给生产者,如果immediate参数为false，那么服务器会将消息存储在队列中，等待消费者来消费
	WaitToConfirm bool

	ExChangeInternal bool // 如果将 Internal 设置为 true，则表示创建的交换器是内部的。内部交换器不能直接接收生产者的消息，只能通过交换器到交换器的绑定来接收消息。这对于创建复杂的路由结构或者隐藏某些交换器的存在非常有用。
	ExChangeAutoDel  bool
	ExChangeDuring   bool
	ExChangeNowait   bool
	ExChangeArgs     amqp.Table
	exChangeName     string
}

var (
	PublishingTimeoutERR = errors.New("publishing timeout")
	ExChangeNameERR      = errors.New("exchange name can not be empty")
	QueueNameERR         = errors.New("queue name can not be empty")
)

type Option func(*ProducerConfig)

type ProducerResultFunc func(msgID string)

func MustInitProducer(url string, vHost string, exChangeName string, qNames []string, opts ...Option) *Producer {
	p, err := InitProducer(url, vHost, exChangeName, qNames, opts...)
	if err != nil {
		log.Fatal(err)
	}
	return p

}

func InitProducer(url string, vHost string, exChangeName string, qNames []string, opts ...Option) (*Producer, error) {
	if vHost == "" {
		vHost = "/"
	}
	config := amqp.Config{
		Vhost:      vHost,
		Properties: amqp.NewConnectionProperties(),
	}
	config.Properties.SetClientConnectionName("sendMsg")

	//Log.Printf("producer: dialing %s", url)
	conn, err := amqp.DialConfig(url, config)
	if err != nil {
		log.Printf("producer: error in dial: %s", err)
		return nil, err
	}
	p, err := NewProducer(conn, exChangeName, qNames, opts...)
	if err != nil {
		log.Printf("producer: error in NewProducer: %s", err)
		return nil, err
	}

	return p, nil
}
func (p *Producer) Close() {
	p.channel.Close()
	p.conn.Close()

}
func NewProducer(conn *amqp.Connection, exChangeName string, qNames []string, opts ...Option) (*Producer, error) {
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("producer: error in channel: %s", err)
		return nil, err
	}
	if exChangeName == "" {
		return nil, ExChangeNameERR
	}
	if len(qNames) == 0 {
		return nil, QueueNameERR
	}
	cfg := &ProducerConfig{
		Durable:          true,
		AutoDelete:       false,
		Exclusive:        false,
		NoWait:           false,
		WaitToConfirm:    true,
		args:             nil,
		AppId:            "",
		DeliveryMode:     amqp.Persistent,
		qNames:           qNames,
		ExChangeInternal: false,
		ExChangeAutoDel:  false,
		ExChangeDuring:   true,
		ExChangeNowait:   false,
		ExChangeArgs:     nil,
		exChangeName:     exChangeName,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	err = ch.ExchangeDeclare(
		exChangeName,         // name
		"fanout",             // type
		cfg.ExChangeDuring,   // durable
		cfg.ExChangeAutoDel,  // auto-deleted
		cfg.ExChangeInternal, // internal
		cfg.ExChangeNowait,   // no-wait
		cfg.ExChangeArgs,     // arguments
	)
	if err != nil {
		return nil, err
	}

	for _, qName := range qNames {
		_, err = ch.QueueDeclare(
			qName,          // name of the queue
			cfg.Durable,    // 是否持久化
			cfg.AutoDelete, // 是否自动删除
			cfg.Exclusive,  // 是否排他
			cfg.NoWait,     // no-wait,true:不等待服务器的响应,如果noWait参数为false，那么RabbitMQ服务器在创建队列后会向客户端发送一个确认消息。客户端会等待这个确认消息，然后再继续执行后续的操作。这种方式可以确保队列已经成功创建，但是会增加一些延迟。
			cfg.args,       // arguments
		)
		if err != nil {
			return nil, err
		}
		err = ch.QueueBind(qName, "", exChangeName, false, nil)
		if err != nil {
			return nil, err
		}
	}

	if cfg.WaitToConfirm {
		if err = ch.Confirm(false); err != nil {
			return nil, err
		}
	}

	return &Producer{
		conn:    conn,
		channel: ch,
		cfg:     cfg,
	}, nil
}
func (p *Producer) PublishWithContext(ctx context.Context, message []byte) (string, error) {
	timeoutCtx, cancelFunc := context.WithTimeout(ctx, time.Second*10)
	defer cancelFunc()
	confirmation, err := p.channel.PublishWithDeferredConfirmWithContext(
		timeoutCtx,
		p.cfg.exChangeName, // exchange
		"",                 // routing key
		p.cfg.Mandatory,    // mandatory参数决定了当消息无法路由到任何队列时，服务器应该如何处理。如果mandatory参数为true，那么服务器会将无法路由的消息返回给生产者。如果mandatory参数为false，那么服务器会直接丢弃无法路由的消息。
		p.cfg.Immediate,    // immediate参数决定了当消息路由到队列，但是队列中没有消费者时，服务器应该如何处理。如果immediate参数为true，那么服务器会将这种消息返回给生产者。如果immediate参数为false，那么服务器会将消息存储在队列中，等待消费者来消费
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			DeliveryMode:    p.cfg.DeliveryMode,
			Priority:        0,
			AppId:           p.cfg.AppId,
			Body:            message,
		})
	if err != nil {
		return "", err
	}
	if p.cfg.WaitToConfirm {
		select {
		case <-confirmation.Done():
		case <-timeoutCtx.Done():
			return "", PublishingTimeoutERR
		}
		if confirmation.Acked() {
			return strconv.FormatUint(confirmation.DeliveryTag, 10), nil
		}
	}

	return "", nil
}
