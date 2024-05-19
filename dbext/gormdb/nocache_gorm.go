package gormdb

import (
	"context"
	"errors"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"log"
	"time"
)

type GormDB struct {
	singleFlight *singleflight.Group
	db           *gorm.DB
}

func MustNewGormDB(c Config) *GormDB {
	gormDB, err := NewGormDB(c)
	if err != nil {
		log.Fatalf("NewCacheGormDB err:%v", err)
		return nil
	}
	return gormDB
}

func NewGormDB(c Config) (*GormDB, error) {
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

	return &GormDB{
		singleFlight: &singleflight.Group{},
		db:           db,
	}, nil
}
func (cg *GormDB) QueryCtx(ctx context.Context, result any, fn QueryCtxFn) error {
	return fn(ctx, result, cg.db)
}
func (cg *GormDB) QueryNoCacheCtx(ctx context.Context, result any, fn QueryCtxFn) error {
	return fn(ctx, result, cg.db)
}
func (cg *GormDB) takeCtx(ctx context.Context, key string, result any, query QueryCtxFn) error {

	_, err, _ := cg.singleFlight.Do(key, func() (interface{}, error) {

		if err := query(ctx, result, cg.db); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if cg.db.Logger != nil {
					cg.db.Logger.Error(ctx, "setCacheWithNotFound err: %v key:%v", err, key)
				}
				return nil, gorm.ErrRecordNotFound
			} else {
				return nil, err
			}
		}
		return result, nil
	})
	return err
}
func (cg *GormDB) ExecCtx(ctx context.Context, execFn ExecCtxFn, keys ...string) (int64, error) {

	result, err := execFn(ctx, cg.db)
	if err != nil {
		return 0, err
	}
	return result, nil
}
