package oauth2_util

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

func Init(redisCli *redis.Client) error {
	if _redis != nil {
		return errors.New("already init")
	}
	_redis = redisCli
	return nil
}
