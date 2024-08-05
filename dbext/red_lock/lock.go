package red_lock

import (
	"context"
	"errors"
	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
	"time"
)

var (
	pool          redsyncredis.Pool
	defaultRS     *RedSync
	timeOutDuring = time.Duration(30) * time.Second
)

var (
// TimeOutLessErr = errors.New("time out sec is less than tempd1")
)

type RedLock struct {
	ticker     *time.Ticker
	mutex      *redsync.Mutex
	innerMutex sync.Mutex
	ctx        context.Context
	lockerMark int32 //0 未上锁，1已上锁
	endChan    chan struct{}
}
type RedSync struct {
	redSync *redsync.Redsync
}

func (r *RedSync) NewMutex(ctx context.Context, name string) (*RedLock, error) {
	mutex := r.redSync.NewMutex(name, redsync.WithExpiry(timeOutDuring), redsync.WithRetryDelay(time.Millisecond*500), redsync.WithTries(10))
	return &RedLock{
		ticker: time.NewTicker((timeOutDuring - time.Millisecond*100) / 3),
		mutex:  mutex,
		//endChan: make(chan struct{}, 2),
		ctx: ctx,
	}, nil
}
func MustInitRedLock(rdb redis.UniversalClient) {
	err := InitRedLock(rdb)
	if err != nil {
		log.Fatalf("InitRedLock err:%v", err)
	}

}
func InitRedLock(rdb redis.UniversalClient) error {
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	pool = goredis.NewPool(rdb) // or, pool := redigo.NewPool(...)
	rs := redsync.New(pool)
	defaultRS = &RedSync{redSync: rs}
	return nil
}
func NewRedSync(rdb redis.UniversalClient) (*RedSync, error) {
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	pool = goredis.NewPool(rdb) // or, pool := redigo.NewPool(...)
	rs := redsync.New(pool)
	return &RedSync{redSync: rs}, nil
}

func NewLock(ctx context.Context, name string) (*RedLock, error) {
	if defaultRS == nil {
		return nil, errors.New("defaultRS not init")
	}

	return defaultRS.NewMutex(ctx, name)
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

	l.innerMutex.Lock()
	defer l.innerMutex.Unlock()
	if l.lockerMark == 1 {
		return true, nil
	}
	err := l.mutex.LockContext(l.ctx)
	if err != nil {
		var errTaken *redsync.ErrTaken
		if errors.As(err, &errTaken) {
			return false, nil
		}
		return false, err
	}
	l.lockerMark = 1
	if l.endChan != nil {
		return false, errors.New("endChan has exists")
	}
	if l.ticker != nil {
		l.ticker.Stop()
	}

	l.ticker = time.NewTicker((timeOutDuring - time.Millisecond*100) / 3)
	l.endChan = make(chan struct{}, 0)
	go func() {
		defer l.ticker.Stop()
		for {
			select {
			case <-l.ctx.Done():
				return
			case <-l.ticker.C:
				_, err = l.mutex.Extend()
				if err != nil {
					l.innerMutex.Lock()
					defer l.innerMutex.Unlock()
					l.lockerMark = 0
					close(l.endChan)
					log.Printf("Extend err:%v", err)
					return
				}
			case <-l.endChan:
				return
			}
		}
	}()

	return true, nil

}
func (l *RedLock) Unlock() (bool, error) {
	l.innerMutex.Lock()
	defer l.innerMutex.Unlock()
	if l.lockerMark == 0 {
		return true, nil
	}
	close(l.endChan)
	l.endChan = nil
	l.lockerMark = 0
	unlock, err := l.mutex.Unlock()
	if err != nil {
		var errNodeTaken *redsync.ErrNodeTaken
		var errTaken *redsync.ErrTaken
		if errors.Is(err, redsync.ErrLockAlreadyExpired) || errors.As(err, &errNodeTaken) || errors.As(err, &errTaken) {
			return true, nil
		}
		return unlock, err
	}

	return unlock, nil

}
