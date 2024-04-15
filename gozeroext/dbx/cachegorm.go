package dbx

import (
	"context"
	redis "github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"gorm.io/gorm"
	"time"
)

type CacheGorm struct {
	db    *gorm.DB
	cache cache.Cache
}
type RedisCache struct {
	rdb *redis.Client
}

func (r RedisCache) Del(keys ...string) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisCache) DelCtx(ctx context.Context, keys ...string) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisCache) Get(key string, val any) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisCache) GetCtx(ctx context.Context, key string, val any) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisCache) IsNotFound(err error) bool {
	//TODO implement me
	panic("implement me")
}

func (r RedisCache) Set(key string, val any) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisCache) SetCtx(ctx context.Context, key string, val any) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisCache) SetWithExpire(key string, val any, expire time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisCache) SetWithExpireCtx(ctx context.Context, key string, val any, expire time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisCache) Take(val any, key string, query func(val any) error) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisCache) TakeCtx(ctx context.Context, val any, key string, query func(val any) error) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisCache) TakeWithExpire(val any, key string, query func(val any, expire time.Duration) error) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisCache) TakeWithExpireCtx(ctx context.Context, val any, key string, query func(val any, expire time.Duration) error) error {
	//TODO implement me
	panic("implement me")
}

func NewConn(db *gorm.DB, rdb *redis.Client, opts ...cache.Option) *CacheGorm {
	return &CacheGorm{
		db:    db,
		cache: &RedisCache{rdb: rdb},
	}
}
