package consumer

import (
	"context"
	"github.com/heshiyingx/gotool/mqext/rabbitmqext/consumer/simple"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

//type ConsumeFunc func(msg amqp.Delivery) error

var (
	//simpleConfig = SimpleConfig{}
	// conn *amqp.Connection
	simpleProducer *simple.Consumer
)

type SimpleConfig struct {
	Url   string
	VHost string
}

func MustSampleConsumeWithQName(ctx context.Context, sc SimpleConfig, qName string, f simple.ConsumeFunc) error {
	config := amqp.Config{Properties: amqp.NewConnectionProperties(), Vhost: sc.VHost}
	config.Properties.SetClientConnectionName("sample-consumer")

	conn, err := amqp.DialConfig(sc.Url, config)
	if err != nil {
		log.Fatalf("consumer: error in dial: %s", err)
	}
	consumer, err := simple.NewConsume(conn)
	if err != nil {
		return err
	}
	err = consumer.SampleConsumeWithQName(ctx, qName, f)
	if err != nil {
		return err
	}
	return nil
}
