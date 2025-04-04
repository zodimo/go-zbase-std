package mutex

import (
	"context"
	"github.com/zodimo/go-zstd/complete"
)

var _ CancellableMutex = (*cancellableMutex)(nil)

var _ complete.Complete = (*cancellableMutex)(nil)

type CancellableMutex interface {
	Lock(context.Context) error
	Unlock()
	GetKey() string
	IsLocked() bool
}

type cancellableMutex struct {
	key         string
	lockChannel chan struct{} // Global, unexported lock channel
	locked      bool
}

func (cm *cancellableMutex) IsLocked() bool {
	return cm.locked
}

func (cm *cancellableMutex) GetKey() string {
	return cm.key
}

func GetOrNewCancellableMutex(key string) CancellableMutex {
	registry = GetMutexRegistry()
	optionalRegistry := registry.GetMutex(key)
	maybeMutex, some := optionalRegistry.Value()
	if some {
		return maybeMutex.(CancellableMutex)
	}
	mutex := NewCancellableMutex(key)
	_ = registry.Register(mutex)
	return mutex
}
func NewCancellableMutex(key string) CancellableMutex {
	return &cancellableMutex{
		lockChannel: make(chan struct{}, 1),
		key:         key,
	}
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
		cm.locked = false
	}
}

func (cm *cancellableMutex) Complete() bool {
	return cm.key != ""
}
