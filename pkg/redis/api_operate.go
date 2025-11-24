package redis

import (
	"context"
	"fmt"
)

const (
	_dbOP = 7
)

var (
	_redisOP *client
)

func InitOP(ctx context.Context, cfg Config) error {
	if _redisOP != nil {
		return fmt.Errorf("redis operate client already init")
	}
	c, err := newClient(ctx, cfg, _dbOP)
	if err != nil {
		return err
	}
	_redisOP = c
	return nil
}

func StopOP() {
	if _redisOP != nil {
		_redisOP.Stop()
		_redisOP = nil
	}
}

func OP() *client {
	return _redisOP
}
