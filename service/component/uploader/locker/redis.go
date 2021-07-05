package locker

import (
	"context"
	"errors"
	"time"

	"github.com/latolukasz/orm"
	tusd "github.com/tus/tusd/pkg/handler"
)

type RedisLocker struct {
	ormService *orm.Engine
}

func (locker *RedisLocker) NewLock(id string) (tusd.Lock, error) {
	return &redisLock{id: id, redis: locker.ormService.GetRedis()}, nil
}

type redisLock struct {
	id        string
	redis     *orm.RedisCache
	redisLock *orm.Lock
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
