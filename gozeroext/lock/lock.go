package lock

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

func GetLock(ctx context.Context, rds *redis.Redis, key string, ttlSec int) *redis.RedisLock {
	lock := redis.NewRedisLock(rds, key)
	lock.SetExpire(ttlSec)
	return lock
}
