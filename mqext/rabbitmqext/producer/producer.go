package producer

import (
	"context"
	"github.com/heshiyingx/gotool/mqext/rabbitmqext/producer/simple"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

var (
// producerConfig = config.Config{}
// conn           *amqp.Connection
)

//type Option func(*config.Config)

var simpleProducer *simple.Producer

func Close() {
	simpleProducer.Close()

}
func MustInitSimple(url string, vHost string, qName string, opts ...simple.Option) {
	if vHost == "" {
		vHost = "/"
	}
	config := amqp.Config{
		Vhost:      vHost,
		Properties: amqp.NewConnectionProperties(),
	}
	config.Properties.SetClientConnectionName("producer-with-confirms")

	//Log.Printf("producer: dialing %s", url)
	conn, err := amqp.DialConfig(url, config)
	if err != nil {
		log.Fatalf("producer: error in dial: %s", err)
	}
	if simpleProducer == nil {
		p, err := simple.NewProducer(conn, qName, opts...)
		if err != nil {
			log.Fatalf("producer: error in NewProducer: %s", err)
		}
		simpleProducer = p

	}
}

func SimpleMustConvertAndSend(ctx context.Context, message []byte) (string, error) {
	// 需要再rabbitmq中把队列创建出来
	msgMQID, err := simpleProducer.PublishWithContext(ctx, message)
	if err != nil {
		return "", err
	}
	return msgMQID, nil

}
