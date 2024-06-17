package redis_script

import (
	"context"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"testing"
	"time"
)

func Test_GetIncrSeqScript(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     ":6379",
		Password: "root",
		DB:       0,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 200; i++ {
		go func() {
			result, err := DecrZeroDelScript.Run(context.Background(), client, []string{"a"}).Result()
			if err != nil {
				t.Error("错误:", err)
				return
			}
			t.Log("成功:", result)
		}()
	}
	time.Sleep(time.Minute)

}
func Test_SetMustGTOldScript(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     ":6379",
		Password: "root",
		DB:       0,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 50; i++ {
		go func() {
			newV := rand.Int63n(1000)
			t.Log(newV)
			result, err := SetMustGTOldScript.Run(context.Background(), client, []string{"a"}, []interface{}{newV}).Result()
			if err != nil {
				t.Error("错误:", err)
				return
			}
			t.Log("成功:", result)
		}()
	}
	time.Sleep(time.Minute)

}
func Test_SafeDECRScript(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     ":6379",
		Password: "root",
		DB:       0,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 50; i++ {
		go func() {
			newV := rand.Int63n(1000)
			t.Log(newV)
			result, err := SafeDECRScript.Run(context.Background(), client, []string{"a"}, []interface{}{newV}).Result()
			if err != nil {
				t.Error("错误:", err)
				return
			}
			t.Log("成功:", result)
		}()
	}
	time.Sleep(time.Minute)

}
func Test_IncrExpireScript(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     ":6379",
		Password: "root",
		DB:       0,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		t.Error(err)
	}
	result, err := IncrExpireScript.Run(context.Background(), client, []string{"incr_demo"}, []interface{}{10}).Result()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result)
}
