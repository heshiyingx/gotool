package gormdb

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/heshiyingx/gotool/dbext/red_lock"
	"github.com/panjf2000/ants/v2"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"time"
)

const (
	notFoundPlaceholder = "*"
	// make the expiry unstable to avoid lots of cached items expire at the same time
	// make the unstable expiry to be [0.95, tempd1.05] * seconds
	expiryDeviation = 0.05
)

type (
	CacheMongoDB struct {
		rdb               redis.UniversalClient
		singleFlight      *singleflight.Group
		notFoundExpireSec int
		cacheExpireSec    int
		randSec           int
		db                *mongo.Client
		collection        *mongo.Collection
		antPool           *ants.Pool
	}
)

func MustCacheMongoDB(c Config, opts ...Option) *CacheMongoDB {
	gormDB, err := NewCacheMongoDB(c, opts...)
	if err != nil {
		log.Fatalf("NewCacheGormDB err:%v", err)
		return nil
	}
	return gormDB
}

func NewCacheMongoDB(c Config, opts ...Option) (*CacheMongoDB, error) {

	pool, err := ants.NewPool(20, ants.WithExpiryDuration(time.Minute*5))
	if err != nil {
		return nil, err
	}
	// 设置客户端选项
	clientOptions := options.Client().
		ApplyURI(c.URI).
		//ApplyURI("mongodb://localhost:27017").
		SetMaxPoolSize(150).
		SetMinPoolSize(10).
		SetConnectTimeout(time.Second * 10).
		SetRetryReads(true).
		SetRetryWrites(true)
	if c.User != "" && c.Pwd != "" {
		clientOptions.SetAuth(options.Credential{
			AuthSource:  "admin",
			Username:    c.User,
			Password:    c.Pwd,
			PasswordSet: true,
		})
	}
	for _, opt := range opts {
		opt(clientOptions)
	}

	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Cannot connect to MongoDB!")
	}
	collection := client.Database(c.DataBase).Collection(c.Collection)
	return &CacheMongoDB{
		rdb:               c.Rdb,
		singleFlight:      &singleflight.Group{},
		notFoundExpireSec: c.NotFoundExpireSec,
		cacheExpireSec:    c.CacheExpireSec,
		randSec:           c.RandSec,
		db:                client,
		antPool:           pool,
		collection:        collection,
	}, nil
}

func (cg *CacheMongoDB) QueryCtx(ctx context.Context, result any, key string, fn QueryCtxFn) error {
	return cg.takeCtx(ctx, key, result, fn, func(result string) error {
		_, err := cg.rdb.Set(ctx, key, result, genDuring(cg.cacheExpireSec, cg.randSec)).Result()
		return err
	})
}
func (cg *CacheMongoDB) QueryNoCacheCtx(ctx context.Context, result any, fn QueryCtxFn) error {
	return fn(ctx, result, cg.collection)
}
func (cg *CacheMongoDB) takeCtx(ctx context.Context, key string, result any, query QueryCtxFn, cacheFn CacheFn) error {

	_, err, _ := cg.singleFlight.Do(key, func() (interface{}, error) {
		val, err := cg.rdb.Get(ctx, key).Result()
		if errors.Is(err, redis.Nil) {
			err = nil
		}
		if val == notFoundPlaceholder {
			return nil, gorm.ErrRecordNotFound
		}
		if val != "" {
			err = json.Unmarshal([]byte(val), result)
			return result, err

		}

		if err = query(ctx, result, cg.collection); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = cg.setCacheWithNotFound(ctx, key)
				if err != nil {
					log.Printf("setCacheWithNotFound err: %v key:%v\n", err, key)
				}
				return nil, gorm.ErrRecordNotFound
			} else {
				return nil, err
			}
		}
		resultBytes, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}
		err = cacheFn(string(resultBytes))
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	return err
}
func (cg *CacheMongoDB) ExecCtx(ctx context.Context, execFn ExecCtxFn, keys ...string) (int64, error) {
	err := cg.rdb.Del(ctx, keys...).Err()
	if err != nil {
		return 0, err
	}
	result, err := execFn(ctx, cg.collection)
	if err != nil {
		return 0, err
	}
	err = cg.rdb.Del(ctx, keys...).Err()
	if err != nil {
	}
	err = cg.antPool.Submit(func() {
		deadline, cancelFunc := context.WithDeadline(ctx, time.Now().Add(time.Millisecond*200))
		defer cancelFunc()
		select {
		case <-deadline.Done():
		}
		err = cg.rdb.Del(ctx, keys...).Err()
		if err != nil {
		}
	})
	if err != nil {
		//return 0, err
	}
	return result, nil
}
func (cg *CacheMongoDB) QuerySafeSingleFromDB(ctx context.Context, key string, result any, queryFn QueryCtxFn, expire int) error {
	val, err := cg.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		err = nil
	}
	if val == notFoundPlaceholder {
		return gorm.ErrRecordNotFound
	}
	if val != "" {
		err = json.Unmarshal([]byte(val), result)
		return err

	}
	redSync, err := red_lock.NewRedSync(cg.rdb)
	if err != nil {
		return err
	}
	locker, err := red_lock.NewLockWithRS(ctx, redSync, key)
	if err != nil {
		return err
	}
	defer locker.Unlock()
	for {
		lock, err := locker.Lock()
		if err != nil {
			return err
		}
		if lock {
			val, err = cg.rdb.Get(ctx, key).Result()
			if errors.Is(err, redis.Nil) {
				err = nil
			}
			if val == notFoundPlaceholder {
				return gorm.ErrRecordNotFound
			}
			if val != "" {
				err = json.Unmarshal([]byte(val), result)
				return err
			}
			err = queryFn(ctx, result, cg.collection)
			if err != nil {
				return err
			}
			resultBytes, err := json.Marshal(result)
			if err != nil {
				return err
			}
			_, err = cg.rdb.Set(ctx, key, resultBytes, genDuring(expire, cg.randSec)).Result()
			return nil
		}
	}
}

func (cg *CacheMongoDB) setCacheWithNotFound(ctx context.Context, key string) error {
	expire := time.Second*time.Duration(cg.notFoundExpireSec) + genDuring(cg.randSec, cg.notFoundExpireSec)
	_, err := cg.rdb.SetNX(ctx, key, notFoundPlaceholder, expire).Result()
	return err
}
func genDuring(oriSec int, randSec int) time.Duration {
	if oriSec == 0 {
		return 0
	}
	if randSec == 0 {
		randSec = 5
	}
	n := rand.Int31n(int32(time.Duration(randSec) * time.Second / time.Millisecond))
	return time.Duration(n)*time.Millisecond + time.Duration(oriSec)*time.Second
}
