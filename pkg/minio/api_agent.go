package minio

import (
	"context"
)

var (
	_minioAgent *client
)

func InitAgent(ctx context.Context, cfg Config) error {
	if _minioAgent == nil {
		c, err := newClient(cfg)
		if err != nil {
			return err
		}
		_minioAgent = c
	}
	return nil
}

func Agent() *client {
	return _minioAgent
}
