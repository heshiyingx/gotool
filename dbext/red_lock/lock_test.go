package red_lock

import (
	"context"
	"github.com/redis/go-redis/v9"
	"testing"
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
		panic(err)
		//return
	}
	locked, err = locker.Lock()
	if err != nil {
		panic(err)
		//return
	}
	locked, err = locker.Lock()
	if err != nil {
		panic(err)
		//return
	}
	//time.Sleep(time.Hour)
	unlocked, err := locker.Unlock()
	if err != nil {
		panic(err)
	}
	t.Log(locked, unlocked)
}
