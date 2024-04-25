package main

import (
	"context"
	"fmt"
	"github.com/heshiyingx/gotool/mqext/rabbitmqext/consumer/simple"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func main() {
	consumer1()
	//go consumer1()
	//go consumer2()
	time.Sleep(time.Hour)
}
func consumer1() {
	//msgs := make([]amqp.Delivery, 0, 10)
	config := amqp.Config{Properties: amqp.NewConnectionProperties(), Vhost: "/"}
	config.Properties.SetClientConnectionName("sample-consumer")

	conn, err := amqp.DialConfig("amqp://john:john123@myshuju.top:5672/", config)
	if err != nil {
		log.Fatalf("consumer: error in dial: %s", err)
	}
	consumerXX, err := simple.NewConsume(conn, func(config *simple.ConsumerConfig) {
		config.Tag = "consumer1"
		config.PrefetchCount = 10
		config.AutoAsk = false
	})
	if err != nil {
		log.Fatalf("consumer: error in dial: %s", err)
		//return err
	}
	consumerXX.SampleBatchConsumeWithQName(context.Background(), "testxx", time.Second, 20, func(msgs []amqp.Delivery) ([]uint64, []uint64) {
		for _, msg := range msgs {
			fmt.Println("consumer1", string(msg.Body), ":", msg.DeliveryTag)
			successTags = append(successTags, msg.DeliveryTag)
		}
		return successTags, nil
	})
}
