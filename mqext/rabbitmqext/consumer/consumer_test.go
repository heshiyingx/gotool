package consumer

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"testing"
)

func TestMustSampleConsumeWithQName(t *testing.T) {
	MustSampleConsumeWithQName(context.Background(), SimpleConfig{
		PrefetchCount: 100000,
		Url:           "amqp://john:john123@myshuju.top:5672/",
		VHost:         "/",
	}, "producer_exchange_fanout_test-queue", func(msg amqp.Delivery) error {
		t.Log(string(msg.Body))
		return nil
	})
}
