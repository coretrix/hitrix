package locker

import (
	"context"
	"errors"
	"time"

	"github.com/latolukasz/beeorm/v2"
	tusd "github.com/tus/tusd/pkg/handler"

	"github.com/coretrix/hitrix/datalayer"
)

type RedisLocker struct {
	ormService *datalayer.ORM
}

func (locker *RedisLocker) NewLock(id string) (tusd.Lock, error) {
	return &redisLock{id: id, redis: locker.ormService.GetRedis()}, nil
}

type redisLock struct {
	id        string
	redis     beeorm.RedisCache
	redisLock *beeorm.Lock
}

func (lock *redisLock) Lock() error {
	redisLock, obtained := lock.redis.GetLocker().Obtain(context.Background(), "tusd:upload:lock:"+lock.id, time.Hour*24, time.Second*2)
	if !obtained {
		return errors.New("cannot obtain lock")
	}

	lock.redisLock = redisLock

	return nil
}

func (lock *redisLock) Unlock() error {
	lock.redisLock.Release()

	return nil
}
