package fanout

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"testing"
)

func TestProducer_PublishWithContext(t *testing.T) {
	url := "amqp://john:john123@myshuju.top:5672/"
	vHost := "/"
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
	producer, err := NewProducer(conn, "oneToTwo", []string{"one", "two"})
	if err != nil {
		return
	}
	mqID, err := producer.PublishWithContext(context.Background(), []byte("hello world"))
	if err != nil {
		return
	}
	t.Log(mqID)
}
