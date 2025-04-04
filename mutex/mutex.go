package mutex

import (
	"context"
	"github.com/zodimo/go-zstd/complete"
)

var _ CancellableMutex = (*cancellableMutex)(nil)

type CancellableMutex interface {
	complete.Complete
	Lock(context.Context) error
	Unlock()
	GetKey() string
}

type cancellableMutex struct {
	key         string
	lockChannel chan struct{} // Global, unexported lock channel
	locked      bool
}

func (cm *cancellableMutex) GetKey() string {
	return cm.key
}

func (cm *cancellableMutex) Complete() bool {
	return cm.key != ""
}

func NewCancellableMutex(key string) CancellableMutex {
	register := GetMutexRegistry()
	maybeMutex := register.GetMutex(key)

	rMutex, some := maybeMutex.Value()
	if some {
		return rMutex
	}

	mutex := cancellableMutex{
		lockChannel: make(chan struct{}, 1),
		key:         key,
	}
	mutexMap.Store(key, mutex)
	return &mutex
}

func (cm *cancellableMutex) Lock(ctx context.Context) error {
	select {
	case cm.lockChannel <- struct{}{}:
		cm.locked = true
		return nil // Lock acquired
	case <-ctx.Done():
		return ctx.Err() // Context cancelled or timeout
	}
}

func (cm *cancellableMutex) Unlock() {
	if cm.locked {
		<-cm.lockChannel // Release the lock
	}
}
