package ddos

import (
	"strconv"

	"github.com/latolukasz/beeorm"
)

type DDOS struct {
}

func (t *DDOS) ProtectManyAttempts(redis *beeorm.RedisCache, protectCriterion string, maxAttempts int, ttl int) bool {
	attempts, has := redis.Get("ddos_" + protectCriterion)
	count := 0
	if len(attempts) > 0 {
		var err error
		count, err = strconv.Atoi(attempts)
		if err != nil {
			panic(err)
		}
	}
	if has && count >= maxAttempts {
		return false
	}

	redis.Set("ddos_"+protectCriterion, count+1, ttl)
	return true
}
