package mutex

import (
	"errors"
	"github.com/zodimo/go-zstd/optional"
	"sync"
)

var mutexMap = sync.Map{}

type MutexRegistry struct{}

func GetMutexRegistry() MutexRegistry {
	return MutexRegistry{}
}

func (mr *MutexRegistry) HasMutex(key string) bool {
	if _, ok := mutexMap.Load(key); ok {

		return true
	}
	return false
}

func (mr *MutexRegistry) GetMutex(key string) optional.Option[CancellableMutex] {
	if mutex, ok := mutexMap.Load(key); ok {
		cm, ok := mutex.(cancellableMutex)
		if ok {
			option, err := optional.Some[CancellableMutex](&cm)
			if err == nil {
				return option
			}
		}
		mutexMap.Delete(key)
	}
	return optional.None[CancellableMutex]()
}

func (mr *MutexRegistry) Register(mutex CancellableMutex) error {
	if mr.HasMutex(mutex.GetKey()) {
		return errors.New("mutex already registered")
	}
	mutexMap.Store(mutex.GetKey(), mutex)
	return nil
}
