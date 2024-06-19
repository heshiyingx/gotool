package red_lock

import (
	"context"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

func TestMustInitRedLock(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     ":6379",
		Password: "root",
	})
	MustInitRedLock(client)

	locker, err := NewLock(context.Background(), "my-global-mutex")
	if err != nil {
		return
	}
	locked, err := locker.Lock()
	if err != nil {
		//return
	}
	locked, err = locker.Lock()
	if err != nil {
		//return
	}
	time.Sleep(time.Hour)
	unlocked, err := locker.Unlock()
	if err != nil {
		return
	}
	t.Log(locked, unlocked)
}
