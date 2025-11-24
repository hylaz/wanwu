package lock

import (
	"context"
	"time"
)

type DistributedLockConfig struct {
	Timeout     time.Duration // 锁超时时间
	SyncAcquire bool          // 同步获取锁
}

type DistributedLockService interface {
	AcquireLock(ctx context.Context, lockKey string, config *DistributedLockConfig) error
	ReleaseLock(ctx context.Context, lockKey string) error
}
