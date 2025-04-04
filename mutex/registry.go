package mutex

import (
	"errors"
	"github.com/zodimo/go-zstd/optional"
	"sync"
)

var AlreadyRegisteredError = errors.New("mutex already registered")
var registry MutexRegistry

var _ MutexRegistry = (*mutexRegistry)(nil)

type mutexRegistry struct {
	mutexMap sync.Map
}

type MutexRegistry interface {
	HasMutex(key string) bool
	GetMutex(key string) optional.Option[CancellableMutex]
	Register(mutex CancellableMutex) error
}

func initRegistry() {
	registry = &mutexRegistry{
		mutexMap: sync.Map{},
	}
}

func init() {
	initRegistry()
}

func InitAndGetMutexRegistry() MutexRegistry {
	initRegistry()
	return registry
}

func GetMutexRegistry() MutexRegistry {
	if registry == nil {
		return InitAndGetMutexRegistry()
	}
	return registry
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
			option, err := optional.Some[CancellableMutex](cm)
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
