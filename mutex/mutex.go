package mutex

import (
	"context"
)

// CancellableMutex defines an interface for a mutex that supports cancellation through context.
type CancellableMutex interface {
	// Lock attempts to acquire the lock and blocks until the lock is acquired
	// or the provided context is canceled. Returns an error if the context is canceled.
	Lock(context.Context) error

	// Unlock releases the lock, allowing it to be acquired by another operation.
	Unlock()

	// GetKey returns the unique key associated with this mutex.
	GetKey() string

	// IsLocked returns whether the mutex is currently locked.
	IsLocked() bool
}

// cancellableMutex is an implementation of the CancellableMutex interface.
// It uses a channel to manage lock state and supports context-based cancellation.
type cancellableMutex struct {
	// key is the unique identifier for this mutex.
	key string

	// lockChannel is a channel used to manage the lock state of the mutex.
	lockChannel chan struct{}

	// locked indicates whether the mutex is currently locked.
	locked bool
}

// IsLocked returns whether the mutex is currently in a locked state.
func (cm *cancellableMutex) IsLocked() bool {
	return cm.locked
}

// GetKey returns the unique key associated with this mutex.
func (cm *cancellableMutex) GetKey() string {
	return cm.key
}

// GetOrNewCancellableMutex retrieves an existing CancellableMutex with the given key
// from the mutex registry, or creates a new one if it doesn't exist.
func GetOrNewCancellableMutex(key string) CancellableMutex {
	optionalRegistry := GetMutexRegistry().GetMutex(key)
	maybeMutex, some := optionalRegistry.Value()
	if some {
		return maybeMutex.(CancellableMutex)
	}
	mutex := NewCancellableMutex(key)
	_ = GetMutexRegistry().Register(mutex)
	return mutex
}

// NewCancellableMutex creates and returns a new CancellableMutex with the given key.
// The mutex uses a buffered channel to manage its lock state.
func NewCancellableMutex(key string) CancellableMutex {
	return &cancellableMutex{
		lockChannel: make(chan struct{}, 1),
		key:         key,
	}
}

// Lock attempts to acquire the lock. If the lock is acquired successfully, the method
// returns nil. If the provided context is canceled or times out before the lock
// is acquired, the method returns an error.
func (cm *cancellableMutex) Lock(ctx context.Context) error {
	select {
	case cm.lockChannel <- struct{}{}:
		cm.locked = true
		return nil // Lock acquired
	case <-ctx.Done():
		return ctx.Err() // Context cancelled or timeout
	}
}

// Unlock releases the lock, allowing it to be acquired by another operation.
// It is safe to call Unlock only if the lock is currently held.
func (cm *cancellableMutex) Unlock() {
	if cm.locked {
		<-cm.lockChannel // Release the lock
		cm.locked = false
	}
}

// Complete implements the complete.Complete interface by returning true
// if the mutex has a non-empty key.
func (cm *cancellableMutex) Complete() bool {
	return cm.key != ""
}
