package producer

import (
	"context"
	"github.com/heshiyingx/gotool/mqext/rabbitmqext/producer/simple"
	"strconv"
	"sync"
	"testing"
)

func TestSimpleMustConvertAndSend(t *testing.T) {
	MustInitSimple("amqp://john:john123@myshuju.top:5672/", "/", "testxx", func(config *simple.ProducerConfig) {
		config.WaitToConfirm = false
	})
	defer Close()
	wg := sync.WaitGroup{}
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {
			SimpleMustConvertAndSend(context.Background(), []byte("xhello world:"+strconv.Itoa(i)))
			wg.Done()
		}()

	}
	wg.Wait()

}
