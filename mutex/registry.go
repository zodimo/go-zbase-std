package mutex

import (
	"errors"
	"github.com/zodimo/go-zstd/optional"
	"sync"
	"sync/atomic"
)

var AlreadyRegisteredError = errors.New("mutex already registered")

var registry = newAtomicRegistry()

var _ MutexRegistry = (*mutexRegistry)(nil)

type mutexRegistry struct {
	mutexMap sync.Map
}

type mutexRegistryHolder struct {
	rh MutexRegistry
}

type MutexRegistry interface {
	HasMutex(key string) bool
	GetMutex(key string) optional.Option[CancellableMutex]
	Register(mutex CancellableMutex) error
}

func resetRegistry() {
	registry.Store(mutexRegistryHolder{
		rh: &mutexRegistry{
			mutexMap: sync.Map{},
		},
	})
}

func newAtomicRegistry() *atomic.Value {
	v := &atomic.Value{}
	v.Store(mutexRegistryHolder{
		rh: &mutexRegistry{
			mutexMap: sync.Map{},
		},
	})
	return v
}

func GetMutexRegistry() MutexRegistry {
	return registry.Load().(mutexRegistryHolder).rh
}

func (mr *mutexRegistry) HasMutex(key string) bool {
	if _, ok := mr.mutexMap.Load(key); ok {
		return true
	}
	return false
}

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

func (mr *mutexRegistry) Register(mutex CancellableMutex) error {
	if mr.HasMutex(mutex.GetKey()) {
		return AlreadyRegisteredError
	}
	mr.mutexMap.Store(mutex.GetKey(), mutex)
	return nil
}
