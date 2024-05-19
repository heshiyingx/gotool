package gormdb

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/heshiyingx/gotool/dbext/red_lock"
	"github.com/panjf2000/ants/v2"
	"github.com/redis/go-redis/v9"
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
	CacheGormDB struct {
		rdb               redis.UniversalClient
		singleFlight      *singleflight.Group
		notFoundExpireSec int
		cacheExpireSec    int
		randSec           int
		db                *gorm.DB
		antPool           *ants.Pool
	}
)

func MustNewCacheGormDB(c Config) *CacheGormDB {
	gormDB, err := NewCacheGormDB(c)
	if err != nil {
		log.Fatalf("NewCacheGormDB err:%v", err)
		return nil
	}
	return gormDB
}

func NewCacheGormDB(c Config) (*CacheGormDB, error) {
	db, err := gorm.Open(getDialector(c), &c.GormConfig)
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置最大打开的连接数
	sqlDB.SetMaxOpenConns(150)
	// 设置连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(10 * time.Minute)
	_, err = c.Rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	pool, err := ants.NewPool(20, ants.WithExpiryDuration(time.Minute*5))
	if err != nil {
		return nil, err
	}
	return &CacheGormDB{
		rdb:               c.Rdb,
		singleFlight:      &singleflight.Group{},
		notFoundExpireSec: c.NotFoundExpireSec,
		cacheExpireSec:    c.CacheExpireSec,
		randSec:           c.RandSec,
		db:                db,
		antPool:           pool,
	}, nil
}

func (cg *CacheGormDB) QueryCtx(ctx context.Context, result any, key string, fn QueryCtxFn) error {
	return cg.takeCtx(ctx, key, result, fn, func(result string) error {
		_, err := cg.rdb.Set(ctx, key, result, genDuring(cg.cacheExpireSec, cg.randSec)).Result()
		return err
	})
}
func (cg *CacheGormDB) QueryNoCacheCtx(ctx context.Context, result any, fn QueryCtxFn) error {
	return fn(ctx, result, cg.db)
}
func (cg *CacheGormDB) takeCtx(ctx context.Context, key string, result any, query QueryCtxFn, cacheFn CacheFn) error {

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

		if err = query(ctx, result, cg.db); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = cg.setCacheWithNotFound(ctx, key)
				if cg.db.Logger != nil && err != nil {
					cg.db.Logger.Error(ctx, "setCacheWithNotFound err: %v key:%v", err, key)
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
func (cg *CacheGormDB) ExecCtx(ctx context.Context, execFn ExecCtxFn, keys ...string) (int64, error) {
	err := cg.rdb.Del(ctx, keys...).Err()
	if err != nil {
		return 0, err
	}
	result, err := execFn(ctx, cg.db)
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
func (cg *CacheGormDB) QuerySafeSingleFromDB(ctx context.Context, key string, result any, queryFn QueryCtxFn, expire int) error {
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
			err = queryFn(ctx, result, cg.db)
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

func (cg *CacheGormDB) setCacheWithNotFound(ctx context.Context, key string) error {
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
