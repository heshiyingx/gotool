package red_lock

import (
	"context"
	"errors"
	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"log"
	"sync/atomic"
	"time"
)

var (
	pool          redsyncredis.Pool
	defaultRS     *redsync.Redsync
	timeOutDuring = time.Duration(30) * time.Second
)

var (
// TimeOutLessErr = errors.New("time out sec is less than tempd1")
)

type RedLock struct {
	ticker  *time.Ticker
	mutex   *redsync.Mutex
	ctx     context.Context
	endMark int32 //0 待lock,1请求locke 2上锁成功 3已释放锁 4:从已释放锁恢复到0的过程中 5:释放做过程中
	endChan chan struct{}
}

func MustInitRedLock(rdb redis.UniversalClient) *redsync.Redsync {
	rs, err := InitRedLock(rdb)
	if err != nil {
		log.Fatalf("InitRedLock err:%v", err)
	}
	return rs

}
func InitRedLock(rdb redis.UniversalClient) (*redsync.Redsync, error) {
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	pool = goredis.NewPool(rdb) // or, pool := redigo.NewPool(...)
	defaultRS = redsync.New(pool)
	return defaultRS, nil
}
func NewRedSync(rdb redis.UniversalClient) (*redsync.Redsync, error) {
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	pool = goredis.NewPool(rdb) // or, pool := redigo.NewPool(...)
	rs := redsync.New(pool)
	return rs, nil
}
func NewLock(ctx context.Context, name string) (*RedLock, error) {
	if defaultRS == nil {
		log.Fatalf("NewLock err:%v", "rs  is nil")
	}

	mutex := defaultRS.NewMutex(name, redsync.WithExpiry(timeOutDuring), redsync.WithRetryDelay(time.Millisecond*500), redsync.WithTries(10))
	return &RedLock{
		ticker:  time.NewTicker((timeOutDuring - time.Millisecond*100) / 3),
		mutex:   mutex,
		endChan: make(chan struct{}, 2),
		ctx:     ctx,
	}, nil
}
func NewLockWithRS(ctx context.Context, rs *redsync.Redsync, name string) (*RedLock, error) {
	if rs == nil {
		return nil, errors.New("rs is nil")
	}
	if name == "" {
		return nil, errors.New("name is empty")
	}
	name = "locker:" + name

	mutex := rs.NewMutex(name, redsync.WithExpiry(timeOutDuring), redsync.WithRetryDelay(time.Millisecond*500), redsync.WithTries(10))
	return &RedLock{
		ticker:  time.NewTicker((timeOutDuring - time.Millisecond*100) / 3),
		mutex:   mutex,
		endChan: make(chan struct{}, 2),
		ctx:     ctx,
	}, nil
}
func (l *RedLock) Lock() (bool, error) {

	if atomic.CompareAndSwapInt32(&l.endMark, 3, 4) {
		l.endChan = make(chan struct{}, 2)
		l.ticker = time.NewTicker((timeOutDuring - time.Millisecond*100) / 3)
		atomic.CompareAndSwapInt32(&l.endMark, 4, 0)
	}
	ticker := l.ticker
	if atomic.CompareAndSwapInt32(&l.endMark, 0, 1) {
		err := l.mutex.LockContext(l.ctx)
		if err != nil {
			atomic.CompareAndSwapInt32(&l.endMark, 1, 0)
			return false, err
		} else {
			atomic.CompareAndSwapInt32(&l.endMark, 1, 2)
		}
		go func() {

			defer ticker.Stop()
			for {
				select {
				case <-l.ctx.Done():
					return
				case <-ticker.C:
					_, err = l.mutex.Extend()
					if err != nil {
						log.Printf("Extend err:%v", err)
						return
					}
				case <-l.endChan:
					return
				}
			}
		}()
	}

	return false, nil

}
func (l *RedLock) Unlock() (bool, error) {
	if atomic.CompareAndSwapInt32(&l.endMark, 2, 5) || atomic.CompareAndSwapInt32(&l.endMark, 1, 5) {
		close(l.endChan)

		unlock, err := l.mutex.Unlock()
		if err != nil {
			var errNodeTaken *redsync.ErrNodeTaken
			var errTaken *redsync.ErrTaken
			if errors.Is(err, redsync.ErrLockAlreadyExpired) || errors.As(err, &errNodeTaken) || errors.As(err, &errTaken) {
				atomic.CompareAndSwapInt32(&l.endMark, 5, 3)
				return true, nil
			}
			return unlock, err
		}
		atomic.CompareAndSwapInt32(&l.endMark, 5, 3)
		return unlock, nil
	}
	return false, nil

}
