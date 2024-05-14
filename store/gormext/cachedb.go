package gormext

import (
	"context"
	"database/sql"
	"github.com/heshiyingx/gotool/store/cacheext"
	"github.com/panjf2000/ants/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"log"
	"time"
)

type (
	GormCacheDB struct {
		db          *gorm.DB
		cache       cacheext.Cache
		pool        *ants.Pool
		delayDuring time.Duration
	}
	ExecFn func(conn *gorm.DB) (sql.Result, error)
	// ExecCtxFn defines the sql exec method.
	ExecCtxFn  func(ctx context.Context, db *gorm.DB) (sql.Result, error)
	QueryCtxFn func(ctx context.Context, db *gorm.DB, v any) error
)

func NewGormCacheDB(db *gorm.DB, c cacheext.Cache) *GormCacheDB {
	pool, err := ants.NewPool(100)
	if err != nil {
		log.Fatal(err)
	}
	return &GormCacheDB{
		db:    db,
		cache: c,
		pool:  pool,
	}
}
func (d *GormCacheDB) UpdateExec(ctx context.Context, exec ExecCtxFn, delKeys ...string) (
	sql.Result, error) {

	result, err := exec(ctx, d.db)
	if err != nil {
		return nil, err
	}
	delTime := time.Now().Add(d.delayDuring)
	err = d.pool.Submit(func() {
		remainTime := delTime.Sub(time.Now())
		if remainTime > 0 {
			time.Sleep(remainTime)
		}
		err = d.cache.DelCtx(ctx, delKeys...)
		if err != nil {
			logx.Error(err)
			return
		}
	})
	if err != nil {
		logx.Error(err)
	}
	return result, d.cache.DelCtx(ctx, delKeys...)
}
func (d *GormCacheDB) QueryExec(ctx context.Context) {

}
