package gormdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/heshiyingx/gotool/dbext/red_lock"
	"github.com/heshiyingx/gotool/dbext/redis_script"
	"github.com/heshiyingx/gotool/strext"
	"github.com/panjf2000/ants/v2"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"runtime"
	"time"
)

const (
	notFoundPlaceholder = "*"
	keyUpdatePrefix     = "updating:"
	// make the expiry unstable to avoid lots of cached items expire at the same time
	// make the unstable expiry to be [0.95, tempd1.05] * seconds
	expiryDeviation = 0.05
)

var TypeErr = errors.New("type is err")

type (
	CacheGormDB[T any, P int64 | uint64 | string] struct {
		rdb               redis.UniversalClient
		singleFlight      *singleflight.Group
		notFoundExpireSec int
		cacheExpireSec    int
		randSec           int
		db                *gorm.DB
		antPool           *ants.Pool
		//antFailChan       chan []string
	}
	pkInfoDefine[P int64 | uint64 | string] struct {
		pkCacheKey string
		p          P
	}
)

func MustNewCacheGormDB[T any, P int64 | uint64 | string](c Config) *CacheGormDB[T, P] {
	gormDB, err := NewCacheGormDB[T, P](c)
	if err != nil {
		log.Fatalf("NewCacheGormDB err:%v", err)
		return nil
	}
	return gormDB
}

func NewCacheGormDB[T any, P int64 | uint64 | string](c Config) (*CacheGormDB[T, P], error) {
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
	pool, err := ants.NewPool(runtime.NumCPU(), ants.WithExpiryDuration(time.Minute*5))
	if err != nil {
		return nil, err
	}
	cacheGromDB := &CacheGormDB[T, P]{
		rdb:               c.Rdb,
		singleFlight:      &singleflight.Group{},
		notFoundExpireSec: c.NotFoundExpireSec,
		cacheExpireSec:    c.CacheExpireSec,
		randSec:           c.RandSec,
		db:                db,
		antPool:           pool,
		//antFailChan:       make(chan []string, 20000),
	}
	if c.PreFunc != nil {
		c.PreFunc(cacheGromDB.db)
	}
	//go func() {
	//	for {
	//		keys, ok := <-cacheGromDB.antFailChan
	//		if !ok {
	//			break
	//		}
	//		if len(keys) > 0 {
	//			err = cacheGromDB.antPool.Submit(func() {
	//				err = cacheGromDB.rdb.Del(context.Background(), keys...).Err()
	//				if err != nil {
	//					log.Printf("ant pool task doing err:%v", err)
	//					cacheGromDB.antFailChan <- keys
	//				}
	//			})
	//			if err != nil {
	//				log.Printf("ant pool task antFailChan Submit err:%v", err)
	//			}
	//		}
	//	}
	//}()
	return cacheGromDB, nil
}

/*---------------*/

func (cg *CacheGormDB[T, P]) QueryOneCtx(ctx context.Context, result any, key string, queryPrimaryFn QueryPrimaryKeyFn[P], primaryCachePrefix string, queryModelFn QueryModelByPKFn[T, P]) error {
	defer func() {
		logx.WithContext(ctx).Debugf("QueryOneCtx key:%v,result:%v", key, strext.ToJsonStr(result))
	}()
	var primaryValue P
	err := cg.takeCtx(ctx, key, &primaryValue, func(ctx context.Context, ret any, db *gorm.DB) error {
		//var p P
		typedRet, ok := ret.(*P)
		if !ok {
			return fmt.Errorf("unexpected type for ret, expected *P, got: %T", ret)
		}
		err := queryPrimaryFn(ctx, typedRet, cg.db)
		if err != nil {
			return err
		}
		return err
	}, func(result string, waitUpdate bool) error {
		if waitUpdate {
			_, err := cg.rdb.Set(ctx, key, result, time.Second*2).Result()
			return err
		} else {
			isSet, err := cg.rdb.SetNX(ctx, key, result, genDuring(cg.cacheExpireSec, cg.randSec)).Result()
			if err != nil {
				return err
			}
			if !isSet {
				_, err = cg.rdb.Set(ctx, key, result, time.Second*2).Result()
				return err
			}
			return nil
		}

	})
	if err != nil {
		return err
	}

	primaryCacheKey := fmt.Sprintf("%v%v", primaryCachePrefix, primaryValue)

	err = cg.takeCtx(ctx, primaryCacheKey, result, func(ctx context.Context, r any, db *gorm.DB) error {
		rModel, ok := r.(*T)
		if !ok {
			return fmt.Errorf("unexpected type for ret, expected *P, got: %T", r)
		}
		err = queryModelFn(ctx, rModel, primaryValue, db)
		if err != nil {
			return err
		}
		return nil
	}, func(result string, waitUpdate bool) error {
		if waitUpdate {
			_, err = cg.rdb.Set(ctx, primaryCacheKey, result, time.Second*2).Result()
			return err
		} else {
			isSet, err := cg.rdb.SetNX(ctx, primaryCacheKey, result, genDuring(cg.cacheExpireSec, cg.randSec)).Result()
			if err != nil {
				return err
			}
			if !isSet {
				_, err = cg.rdb.Set(ctx, primaryCacheKey, result, time.Second*2).Result()
				return err
			}
			return nil

		}

	})
	return err
}
func (cg *CacheGormDB[T, P]) QueryOneByPKCtx(ctx context.Context, r *T, key string, queryFn QueryCtxFn) error {

	return cg.QueryCtx(ctx, r, key, queryFn)
}
func (cg *CacheGormDB[T, P]) QueryToGetPKCtx(ctx context.Context, key string, p *P, queryFn QueryPrimaryKeyFn[P]) error {
	return cg.QueryCtx(ctx, p, key, func(ctx context.Context, r any, db *gorm.DB) error {
		pV, ok := r.(*P)
		if !ok {
			return fmt.Errorf("P type is wrong")
		}
		return queryFn(ctx, pV, db)
	})
}
func (cg *CacheGormDB[T, P]) QueryManyCtx(ctx context.Context, result *[]T, queryPKsFn QueryPrimaryKeysFn[P], primaryCachePrefix string, queryModelFn QueryModelByPKFn[T, P]) error {
	defer func() {
		logx.WithContext(ctx).Debugf("QueryManyCtx  result:%v", strext.ToJsonStr(result))
	}()
	var pks []P
	err := queryPKsFn(ctx, &pks, cg.db)
	if err != nil {
		return err
	}
	if len(pks) == 0 {
		return nil
	}
	pkInfos := make([]pkInfoDefine[P], 0, len(pks))
	for _, pk := range pks {
		pkCacheKey := fmt.Sprintf("%v%v", primaryCachePrefix, pk)
		pkInfos = append(pkInfos, pkInfoDefine[P]{pkCacheKey: pkCacheKey, p: pk})
	}
	//res := make([]T, 0, len(pkInfos))
	return cg.QueryManyByPKsCtx(ctx, result, pks, primaryCachePrefix, queryModelFn)

	//for _, pkInfo := range pkInfos {
	//	var t T
	//	err = cg.takeCtx(ctx, pkInfo.pkCacheKey, &t, func(ctx context.Context, r any, db *gorm.DB) error {
	//		rm, ok := r.(*T)
	//		if !ok {
	//			return fmt.Errorf("unexpected type:%T", r)
	//		}
	//		err = queryModelFn(ctx, rm, pkInfo.p, cg.db)
	//		if err != nil {
	//			return err
	//		}
	//		return nil
	//	}, func(result string) error {
	//		_, err = cg.rdb.Set(ctx, pkInfo.pkCacheKey, result, genDuring(cg.cacheExpireSec, cg.randSec)).Result()
	//		return err
	//	})
	//	if err != nil {
	//		return err
	//	}
	//	*result = append(*result, t)
	//}
	//return nil

}
func (cg *CacheGormDB[T, P]) QueryManyByPKsCtx(ctx context.Context, result *[]T, pks []P, primaryCachePrefix string, queryDBFn QueryModelByPKFn[T, P]) error {
	defer func() {
		logx.WithContext(ctx).Debugf("QueryManyByPKsCtx  pks:%v,result:%v", pks, strext.ToJsonStr(result))
	}()
	pkInfos := make([]pkInfoDefine[P], 0, len(pks))
	for _, pk := range pks {
		pkCacheKey := fmt.Sprintf("%v%v", primaryCachePrefix, pk)
		pkInfos = append(pkInfos, pkInfoDefine[P]{pkCacheKey: pkCacheKey, p: pk})
	}
	for _, pkInfo := range pkInfos {
		var t T
		err := cg.takeCtx(ctx, pkInfo.pkCacheKey, &t, func(ctx context.Context, r any, db *gorm.DB) error {
			rm, ok := r.(*T)
			if !ok {
				return fmt.Errorf("unexpected type:%T", r)
			}
			err := queryDBFn(ctx, rm, pkInfo.p, cg.db)
			if err != nil {
				return err
			}
			return nil
		}, func(result string, waitUpdate bool) error {
			if waitUpdate {
				_, err := cg.rdb.Set(ctx, pkInfo.pkCacheKey, result, time.Second).Result()
				return err
			} else {
				isSet, err := cg.rdb.SetNX(ctx, pkInfo.pkCacheKey, result, genDuring(cg.cacheExpireSec, cg.randSec)).Result()
				if err != nil {
					return err
				}
				if !isSet {
					_, err = cg.rdb.Set(ctx, pkInfo.pkCacheKey, result, time.Second*2).Result()
					return err
				}
				return nil
			}

		})
		if err != nil {
			return err
		}
		*result = append(*result, t)
	}
	return nil
}
func (cg *CacheGormDB[T, P]) QueryNoCacheCtx(ctx context.Context, result any, fn QueryCtxFn) error {
	return fn(ctx, result, cg.db)
}
func (cg *CacheGormDB[T, P]) QuerySafeSingleFromDB(ctx context.Context, key string, result any, queryFn QueryCtxFn, expire int) error {
	defer func() {
		logx.WithContext(ctx).Debugf("QuerySafeSingleFromDB  key:%v,result:%v", key, strext.ToJsonStr(result))
	}()
	val, err := cg.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		err = nil
	}
	if val == notFoundPlaceholder {
		return gorm.ErrRecordNotFound
	}
	if val != "" {
		logx.WithContext(ctx).Debugf("QuerySafeSingleFromDB->cache  key:%v,val:%v", key, val)
		err = json.Unmarshal([]byte(val), result)
		return err

	}
	redSync, err := red_lock.NewRedSync(cg.rdb)
	if err != nil {
		logx.WithContext(ctx).Debugf("QuerySafeSingleFromDB->NewRedSyncErr  key:%v,err:%v", key, err)
		return err
	}
	locker, err := red_lock.NewLockWithRS(ctx, redSync, key)
	if err != nil {
		logx.WithContext(ctx).Debugf("QuerySafeSingleFromDB->NewLockWithRSRR  key:%v,err:%v", key, err)
		return err
	}
	defer locker.Unlock()
	for {
		lock, err := locker.Lock()
		if err != nil {
			logx.WithContext(ctx).Debugf("QuerySafeSingleFromDB->locker.LockErr  key:%v,err:%v", key, err)
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
				logx.WithContext(ctx).Debugf("QuerySafeSingleFromDB->afterLocker-cache  key:%v,val:%v", key, val)
				err = json.Unmarshal([]byte(val), result)
				return err
			}
			err = queryFn(ctx, result, cg.db)
			if err != nil {
				logx.WithContext(ctx).Debugf("QuerySafeSingleFromDB->queryFnErr  key:%v,err:%v", key, err)
				return err
			}
			resultBytes, err := json.Marshal(result)
			if err != nil {
				logx.WithContext(ctx).Debugf("QuerySafeSingleFromDB->json.Marsha  key:%v,jsonStr:%v,err:%v", key, string(resultBytes), err)
				return err
			}
			isSet, err := cg.rdb.SetNX(ctx, key, string(resultBytes), genDuring(expire, cg.randSec)).Result()
			if err != nil {
				return err
			}
			if !isSet {
				_, err = cg.rdb.Set(ctx, key, string(resultBytes), time.Second*2).Result()
				return err
			}
			return nil
		}
	}
}
func (cg *CacheGormDB[T, P]) QueryCtx(ctx context.Context, result any, key string, fn QueryCtxFn) error {
	err := cg.takeCtx(ctx, key, result, fn, func(result string, waitUpdate bool) error {
		if waitUpdate {
			_, err := cg.rdb.Set(ctx, key, result, time.Second*2).Result()
			return err
		} else {
			isSet, err := cg.rdb.SetNX(ctx, key, result, genDuring(cg.cacheExpireSec, cg.randSec)).Result()
			if err != nil {
				return err
			}
			if !isSet {
				_, err = cg.rdb.Set(ctx, key, result, time.Second*2).Result()
				return err
			}
			return nil

		}

	})
	logx.WithContext(ctx).Debugf("QueryCtx   key:%v,result:%v,err:%v", key, strext.ToJsonStr(result), err)

	return err
}
func (cg *CacheGormDB[T, P]) QuerySlicesCtxCustom(ctx context.Context, result *[]T, key string, queryDBFn QuerySlicesFn[T], queryCacheFn QueryCacheSlicesCtxFn) error {

	_, err, _ := cg.singleFlight.Do(key, func() (interface{}, error) {
		tString, err := cg.rdb.Type(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		if tString == "string" {
			return nil, gorm.ErrRecordNotFound
		} else if tString == "set" {
			res, isSuccess, err := queryCacheFn(ctx, cg.rdb)
			if err != nil {
				if errors.Is(err, redis.Nil) {
					err = nil
				} else {
					return nil, err
				}
			}
			if isSuccess {
				for _, ele := range res {
					var eObj T
					err = json.Unmarshal([]byte(ele), &eObj)
					if err != nil {
						return nil, err
					}
					*result = append(*result, eObj)
				}

				return result, nil
			}
		} else if tString != "none" {
			return nil, TypeErr
		}

		err = queryDBFn(ctx, result, cg.db)
		if err != nil {
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
		if result == nil || len(*result) == 0 {
			err = cg.setCacheWithNotFound(ctx, key)
			if cg.db.Logger != nil && err != nil {
				cg.db.Logger.Error(ctx, "setCacheWithNotFound err: %v key:%v", err, key)
			}
			return nil, gorm.ErrRecordNotFound
		}
		results := make([]interface{}, 0, len(*result))
		for _, element := range *result {
			retStr, err := json.Marshal(element)
			if err != nil {
				return nil, err
			}
			results = append(results, retStr)
		}
		_, err = cg.rdb.SAdd(ctx, key, results...).Result()
		if err != nil {
			return nil, err
		}
		cg.rdb.Expire(ctx, key, genDuring(cg.cacheExpireSec, cg.randSec))
		*result = nil
		res, isSuccess, err := queryCacheFn(ctx, cg.rdb)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				err = nil
			} else {
				return nil, err
			}
		}
		if isSuccess {
			resBytes, err := json.Marshal(res)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(resBytes, result)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return result, err

	})
	logx.WithContext(ctx).Debugf("QuerySlicesCtxCustom   key:%v,result:%v,err:%v", key, strext.ToJsonStr(result), err)
	return err

}
func (cg *CacheGormDB[T, P]) DelCacheKeys(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return cg.rdb.Del(ctx, keys...).Err()
}

/*---------------*/

func (cg *CacheGormDB[T, P]) takeCtx(ctx context.Context, key string, result any, query QueryCtxFn, cacheFn CacheFn) error {

	_, err, _ := cg.singleFlight.Do(key, func() (interface{}, error) {
		//fmt.Println("进入redis缓存")
		val, err := cg.rdb.Get(ctx, key).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				err = nil
			} else {
				return nil, err
			}
		}

		if val == notFoundPlaceholder {
			return nil, gorm.ErrRecordNotFound
		}
		if val != "" {
			err = json.Unmarshal([]byte(val), result)
			logx.WithContext(ctx).Debugf("takeCtx->Cache   key:%v,result:%v,jsonStr:%v,err:%v", key, strext.ToJsonStr(result), val, err)
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
			logx.WithContext(ctx).Debugf("takeCtx->jsonMarshalErr   key:%v,result:%v,jsonStr_val:%v,err:%v", key, strext.ToJsonStr(result), val, err)
			return nil, err
		}

		isUpdating := true
		_, err = cg.rdb.Get(ctx, keyUpdatePrefix+key).Result()
		if errors.Is(err, redis.Nil) {
			isUpdating = false
			err = nil
		}
		err = cacheFn(string(resultBytes), isUpdating)
		if err != nil {
			logx.WithContext(ctx).Debugf("takeCtx->jsonMarshalErr   key:%v,result:%v,jsonStr_val:%v,err:%v", key, strext.ToJsonStr(result), val, err)
			return nil, err
		}
		return result, nil
	})
	return err
}
func (cg *CacheGormDB[T, P]) ExecCtx(ctx context.Context, execFn ExecCtxFn, keys ...string) (int64, error) {
	if len(keys) > 0 {
		err := cg.rdb.Del(ctx, keys...).Err()
		if err != nil {
			return 0, err
		}
	}
	for _, key := range keys {
		_, err := redis_script.IncrExpireScript.Run(ctx, cg.rdb, []string{keyUpdatePrefix + key}, 20).Result()
		if err != nil {
			return 0, err
		}
	}
	defer func() {
		for _, key := range keys {
			_, _ = redis_script.DecrZeroDelScript.Run(ctx, cg.rdb, []string{keyUpdatePrefix + key}).Result()

		}
	}()
	result, err := execFn(ctx, cg.db)
	if err != nil {
		return 0, err
	}

	if len(keys) > 0 {
		err = cg.rdb.Del(ctx, keys...).Err()
		if err != nil {
			return 0, err
		}
	}
	if len(keys) > 0 {
		err = cg.antPool.Submit(func() {
			deadline, cancelFunc := context.WithDeadline(ctx, time.Now().Add(time.Second))
			defer cancelFunc()
			select {
			case <-deadline.Done():
			}
			err = cg.rdb.Del(context.Background(), keys...).Err()
			if err != nil {
				log.Printf("ant pool task doing err:%v", err)
				//cg.antFailChan <- keys
			}
		})
		if err != nil {
			log.Printf("ant pool task Submit err:%v", err)
		}
	}

	return result, nil
}

func (cg *CacheGormDB[T, P]) setCacheWithNotFound(ctx context.Context, key string) error {
	expire := time.Second*time.Duration(cg.notFoundExpireSec) + genDuring(cg.randSec, cg.notFoundExpireSec)
	_, err := cg.rdb.SetNX(ctx, key, notFoundPlaceholder, expire).Result()
	return err
}
func (cg *CacheGormDB[T, P]) GetRdb() redis.UniversalClient {
	return cg.rdb
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
