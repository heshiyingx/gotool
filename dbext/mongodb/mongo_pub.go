package gormdb

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	QueryCtxFn func(ctx context.Context, result any, collection *mongo.Collection) error
	ExecCtxFn  func(ctx context.Context, collection *mongo.Collection) (int64, error)
	CacheFn    func(result string) error
	Config     struct {
		URI               string
		DataBase          string
		Collection        string
		User              string
		Pwd               string
		Rdb               redis.UniversalClient
		NotFoundExpireSec int
		CacheExpireSec    int
		RandSec           int
	}
	Option func(o *options.ClientOptions)
)
