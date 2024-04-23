package producer

import (
	"context"
	"github.com/heshiyingx/gotool/mqext/rabbitmqext/producer/config"
	"github.com/heshiyingx/gotool/mqext/rabbitmqext/producer/simple"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

var (
	producerConfig = config.Config{}
	//conn           *amqp.Connection
)

type Option func(*config.Config)

var simpleProducer *simple.Producer

func Close() {
	simpleProducer.Close()

}
func MustInit(url string, vHost string, exchange string, qNames []string) {
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
		p, err := simple.NewProducer(conn, qNames)
		if err != nil {
			log.Fatalf("producer: error in NewProducer: %s", err)
		}
		simpleProducer = p

	}
}

func SimpleMustConvertAndSend(ctx context.Context, exchange, qName string, message string) (string, error) {
	// 需要再rabbitmq中把队列创建出来
	msgMQID, err := simpleProducer.PublishWithContext(ctx, qName, true, false, message)
	if err != nil {
		return "", err
	}
	return msgMQID, nil

}
