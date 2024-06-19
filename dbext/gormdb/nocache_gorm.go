package gormdb

import (
	"context"
	"gorm.io/gorm"
	"log"
	"time"
)

type (
	GormDB[T any, P int64 | uint64 | string] struct {
		db *gorm.DB
	}
)

func MustNewGormDB[T any, P int64 | uint64 | string](c Config) *GormDB[T, P] {
	gormDB, err := NewGormDB[T, P](c)
	if err != nil {
		log.Fatalf("NewCacheGormDB err:%v", err)
		return nil
	}
	return gormDB
}

func NewGormDB[T any, P int64 | uint64 | string](c Config) (*GormDB[T, P], error) {
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

	return &GormDB[T, P]{
		db: db,
	}, nil
}

func (cg *GormDB[T, P]) QueryCtx(ctx context.Context, result any, queryFn QueryCtxFn) error {
	return queryFn(ctx, result, cg.db)
}

func (cg *GormDB[T, P]) ExecCtx(ctx context.Context, execFn ExecCtxFn) (int64, error) {

	return execFn(ctx, cg.db)
}
