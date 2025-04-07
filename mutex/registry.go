// Package mutex provides a mechanism for managing and registering cancellable
// mutexes within a thread-safe registry. It ensures efficient retrieval,
// registration, and validation of mutexes across various operations.
package mutex

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/zodimo/go-zbase-std/optional"
)

// AlreadyRegisteredError is returned when attempting to register a mutex
// that is already present in the MutexRegistry.
var AlreadyRegisteredError = errors.New("mutex already registered")

// registry holds the atomic reference to the global mutex registry.
var registry = newAtomicRegistry()

// mutexRegistry implements the MutexRegistry interface and provides
// thread-safe operations on a map of cancellable mutexes.
type mutexRegistry struct {
	mutexMap sync.Map // Synchronizes access to the registered mutexes.
}

// mutexRegistryHolder wraps a MutexRegistry for atomic operations,
// allowing concurrent access and replacement of the underlying registry.
type mutexRegistryHolder struct {
	rh MutexRegistry // The wrapped MutexRegistry implementation.
}

// MutexRegistry defines the interface for managing cancellable mutexes.
// It allows checking for the existence of a mutex, retrieving it,
// and registering new mutexes.
type MutexRegistry interface {

	// HasMutex checks whether a mutex with the given key is already present
	// in the registry.
	//
	// Parameters:
	//   - key: The unique key identifying the mutex.
	//
	// Returns:
	//   - bool: True if the mutex exists; false otherwise.
	HasMutex(key string) bool

	// GetMutex retrieves the mutex associated with the given key from
	// the registry. If the mutex does not exist or is incomplete, it
	// returns an empty optional.
	//
	// Parameters:
	//   - key: The unique key identifying the mutex.
	//
	// Returns:
	//   - optional.Option[CancellableMutex]: The optional containing
	//     the mutex if it exists and is complete; otherwise, an empty optional.
	GetMutex(key string) optional.Option[CancellableMutex]

	// Register adds a new mutex to the registry. If a mutex with the
	// same key already exists, it returns an error.
	//
	// Parameters:
	//   - mutex: The CancellableMutex to be registered.
	//
	// Returns:
	//   - error: AlreadyRegisteredError if a mutex with the same key exists;
	//     nil otherwise.
	Register(mutex CancellableMutex) error
}

// resetRegistry resets the global mutex registry to its initial state.
// This is useful for testing or reinitialization purposes.
func resetRegistry() {
	registry.Store(mutexRegistryHolder{
		rh: &mutexRegistry{
			mutexMap: sync.Map{},
		},
	})
}

// newAtomicRegistry creates and initializes a new atomic registry holder.
// It ensures thread-safe access to the registry across multiple goroutines.
//
// Returns:
//   - *atomic.Value: A pointer to the registry's atomic storage.
func newAtomicRegistry() *atomic.Value {
	v := &atomic.Value{}
	v.Store(mutexRegistryHolder{
		rh: &mutexRegistry{
			mutexMap: sync.Map{},
		},
	})
	return v
}

// GetMutexRegistry retrieves the current global mutex registry.
// It enables access to the centralized registry for all operations.
//
// Returns:
//   - MutexRegistry: The current MutexRegistry instance.
func GetMutexRegistry() MutexRegistry {
	return registry.Load().(mutexRegistryHolder).rh
}

// HasMutex checks if a mutex with the given key exists in the registry.
//
// Parameters:
//   - key: The unique key identifying the mutex.
//
// Returns:
//   - bool: True if a mutex with the key is found; false otherwise.
func (mr *mutexRegistry) HasMutex(key string) bool {
	if _, ok := mr.mutexMap.Load(key); ok {
		return true
	}
	return false
}

// GetMutex retrieves the mutex associated with the given key from the
// mutex registry. If the mutex exists and is complete, it is returned
// as an optional; otherwise, an empty optional is returned.
//
// Parameters:
//   - key: The unique key identifying the mutex.
//
// Returns:
//   - optional.Option[CancellableMutex]: The mutex wrapped in an optional
//     if it exists and is complete; otherwise, an empty optional.
func (mr *mutexRegistry) GetMutex(key string) optional.Option[CancellableMutex] {
	if mutex, ok := mr.mutexMap.Load(key); ok {
		cm, ok := mutex.(*cancellableMutex)
		if ok {
			option, err := optional.SomeComplete[CancellableMutex](cm)
			if err == nil {
				return option
			}
		}
		mr.mutexMap.Delete(key)
	}
	return optional.None[CancellableMutex]()
}

// Register adds a new cancellable mutex to the registry. If a mutex
// with the same key is already registered, the method returns an error.
//
// Parameters:
//   - mutex: The CancellableMutex to be registered.
//
// Returns:
//   - error: AlreadyRegisteredError if the mutex is already registered;
//     nil otherwise.
func (mr *mutexRegistry) Register(mutex CancellableMutex) error {
	if mr.HasMutex(mutex.GetKey()) {
		return AlreadyRegisteredError
	}
	mr.mutexMap.Store(mutex.GetKey(), mutex)
	return nil
}
