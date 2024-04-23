package producer

import (
	"context"
	"strconv"
	"sync"
	"testing"
)

func TestSimpleMustConvertAndSend(t *testing.T) {
	MustInit("amqp://john:john123@myshuju.top:5672/", "/", "testxx", []string{"producer_exchange_fanout_test-queue"})
	defer Close()
	wg := sync.WaitGroup{}
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {
			SimpleMustConvertAndSend(context.Background(), "testxx", "producer_exchange_fanout_test-queue", "xhello world:"+strconv.Itoa(i))
			wg.Done()
		}()

	}
	wg.Wait()

}
