package simple

import (
	"context"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"strconv"
	"time"
)

type Producer struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

var (
	PublishingTimeoutERR = errors.New("publishing timeout")
)

type ProducerResultFunc func(msgID string)

//	func NewProducer(conn *amqp.Connection, exchange string, qNames []string) (*Producer, error) {
//		ch, err := conn.Channel()
//		if err != nil {
//			return nil, err
//		}
//		if err = ch.ExchangeDeclare(
//			exchange, // name
//			"fanout", // type
//			true,     // durable
//			false,    // auto-delete
//			false,    // internal
//			false,    // noWait
//			nil,      // arguments
//		); err != nil {
//			return nil, err
//		}
//		for _, qName := range qNames {
//			_, err = ch.QueueDeclare(
//				qName, // name of the queue
//				true,  // 是否持久化
//				false, // 是否自动删除
//				false, // 是否排他,true:可以有多个消费者
//				false, // no-wait,true:不等待服务器的响应,如果noWait参数为false，那么RabbitMQ服务器在创建队列后会向客户端发送一个确认消息。客户端会等待这个确认消息，然后再继续执行后续的操作。这种方式可以确保队列已经成功创建，但是会增加一些延迟。
//				nil,   // arguments
//			)
//			if err != nil {
//				return nil, err
//			}
//			if err = ch.QueueBind(qName, "", exchange, false, nil); err != nil {
//				return nil, err
//			}
//			if err = ch.Confirm(false); err != nil {
//				return nil, err
//			}
//		}
//
//		return &Producer{
//			Conn:    conn,
//			Channel: ch,
//		}, nil
//	}
//
//	func (p *Producer) SimpleConvertAndSend(ctx context.Context, exchange, qName string, message string) error {
//		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//		defer cancel()
//		err := p.Channel.PublishWithContext(ctx,
//			exchange, // exchange
//			qName,    // routing key
//			false,    // mandatory参数决定了当消息无法路由到任何队列时，服务器应该如何处理。如果mandatory参数为true，那么服务器会将无法路由的消息返回给生产者。如果mandatory参数为false，那么服务器会直接丢弃无法路由的消息。
//			false,    // immediate参数决定了当消息路由到队列，但是队列中没有消费者时，服务器应该如何处理。如果immediate参数为true，那么服务器会将这种消息返回给生产者。如果immediate参数为false，那么服务器会将消息存储在队列中，等待消费者来消费
//			amqp.Publishing{
//				ContentType: "text/plain",
//				Body:        []byte(message),
//			})
//
//		return err
//	}
//
//	func (p *Producer) PublishWithContext(ctx context.Context, exchange, key string, mandatory, immediate bool, message string) (string, error) {
//		timeoutCtx, cancelFunc := context.WithTimeout(ctx, time.Second*200)
//		defer cancelFunc()
//		confirmation, err := p.Channel.PublishWithDeferredConfirmWithContext(
//			timeoutCtx,
//			exchange,  // exchange
//			key,       // routing key
//			mandatory, // mandatory参数决定了当消息无法路由到任何队列时，服务器应该如何处理。如果mandatory参数为true，那么服务器会将无法路由的消息返回给生产者。如果mandatory参数为false，那么服务器会直接丢弃无法路由的消息。
//			immediate, // immediate参数决定了当消息路由到队列，但是队列中没有消费者时，服务器应该如何处理。如果immediate参数为true，那么服务器会将这种消息返回给生产者。如果immediate参数为false，那么服务器会将消息存储在队列中，等待消费者来消费
//			amqp.Publishing{
//				Headers:         amqp.Table{},
//				ContentType:     "text/plain",
//				ContentEncoding: "",
//				DeliveryMode:    amqp.Persistent,
//				Priority:        0,
//				AppId:           "sequential-producer",
//				Body:            []byte(message),
//			})
//		if err != nil {
//			return "", err
//		}
//		select {
//		case <-confirmation.Done():
//		case <-timeoutCtx.Done():
//			return "", PublishingTimeoutERR
//		}
//		if confirmation.Acked() {
//			return strconv.FormatUint(confirmation.DeliveryTag, 10), nil
//		}
//		return "", nil
//	}
func (p *Producer) Close() {
	p.Channel.Close()
	p.Conn.Close()

}

func NewProducer(conn *amqp.Connection, qNames []string) (*Producer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	//if err = ch.ExchangeDeclare(
	//	exchange, // name
	//	"fanout", // type
	//	true,     // durable
	//	false,    // auto-delete
	//	false,    // internal
	//	false,    // noWait
	//	nil,      // arguments
	//); err != nil {
	//	return nil, err
	//}
	for _, qName := range qNames {
		_, err = ch.QueueDeclare(
			qName, // name of the queue
			true,  // 是否持久化
			false, // 是否自动删除
			false, // 是否排他,true:可以有多个消费者
			false, // no-wait,true:不等待服务器的响应,如果noWait参数为false，那么RabbitMQ服务器在创建队列后会向客户端发送一个确认消息。客户端会等待这个确认消息，然后再继续执行后续的操作。这种方式可以确保队列已经成功创建，但是会增加一些延迟。
			nil,   // arguments
		)
		if err != nil {
			return nil, err
		}
		//if err = ch.QueueBind(qName, "", exchange, false, nil); err != nil {
		//	return nil, err
		//}
		if err = ch.Confirm(false); err != nil {
			return nil, err
		}
	}

	return &Producer{
		Conn:    conn,
		Channel: ch,
	}, nil
}
func (p *Producer) SimpleConvertAndSend(ctx context.Context, qName string, message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := p.Channel.PublishWithContext(ctx,
		"",    // exchange
		qName, // routing key
		false, // mandatory参数决定了当消息无法路由到任何队列时，服务器应该如何处理。如果mandatory参数为true，那么服务器会将无法路由的消息返回给生产者。如果mandatory参数为false，那么服务器会直接丢弃无法路由的消息。
		false, // immediate参数决定了当消息路由到队列，但是队列中没有消费者时，服务器应该如何处理。如果immediate参数为true，那么服务器会将这种消息返回给生产者。如果immediate参数为false，那么服务器会将消息存储在队列中，等待消费者来消费
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

	return err
}
func (p *Producer) PublishWithContext(ctx context.Context, key string, mandatory, immediate bool, message string) (string, error) {
	timeoutCtx, cancelFunc := context.WithTimeout(ctx, time.Second*200)
	defer cancelFunc()
	confirmation, err := p.Channel.PublishWithDeferredConfirmWithContext(
		timeoutCtx,
		"",        // exchange
		key,       // routing key
		mandatory, // mandatory参数决定了当消息无法路由到任何队列时，服务器应该如何处理。如果mandatory参数为true，那么服务器会将无法路由的消息返回给生产者。如果mandatory参数为false，那么服务器会直接丢弃无法路由的消息。
		immediate, // immediate参数决定了当消息路由到队列，但是队列中没有消费者时，服务器应该如何处理。如果immediate参数为true，那么服务器会将这种消息返回给生产者。如果immediate参数为false，那么服务器会将消息存储在队列中，等待消费者来消费
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			DeliveryMode:    amqp.Persistent,
			Priority:        0,
			AppId:           "sequential-producer",
			Body:            []byte(message),
		})
	if err != nil {
		return "", err
	}
	select {
	case <-confirmation.Done():
	case <-timeoutCtx.Done():
		return "", PublishingTimeoutERR
	}
	if confirmation.Acked() {
		return strconv.FormatUint(confirmation.DeliveryTag, 10), nil
	}
	return "", nil
}
