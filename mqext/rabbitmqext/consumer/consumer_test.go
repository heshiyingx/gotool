package consumer

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"testing"
)

func TestMustSampleConsumeWithQName(t *testing.T) {
	MustSampleConsumeWithQName(context.Background(), SimpleConfig{
		Url:   "amqp://john:john123@myshuju.top:5672/",
		VHost: "/",
	}, "producer_exchange_fanout_test-queue", func(msg amqp.Delivery) error {
		t.Log(string(msg.Body))
		return nil
	})
}
func TestMustSampleConsumeAckWithQName(t *testing.T) {

	conn, err := amqp.Dial("amqp://john:john123@myshuju.top:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"testxx", // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	//err = ch.Qos(
	//	10,    // prefetch count
	//	0,     // prefetch size
	//	false, // global
	//)
	//failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	batchSize := 10
	batch := make([]amqp.Delivery, 0, batchSize)

	for d := range msgs {
		batch = append(batch, d)
		//time.Sleep(time.Millisecond * 100)
		if len(batch) >= batchSize {
			dTag := uint64(0)
			for _, delivery := range batch {
				log.Println("###", string(delivery.Body), "::::::", delivery.DeliveryTag)
				if delivery.DeliveryTag > dTag {
					dTag = delivery.DeliveryTag
				}
			}
			log.Println("最大Tag###:", dTag)

			// Process batch of messages...
			// After processing, acknowledge all messages in the batch
			for _, d := range batch {
				d.Ack(false)
			}
			// Clear the batch
			batch = batch[:0]
		}
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
