package mutex

import (
	"errors"
	"testing"
)

func TestGetMutexRegistry(t *testing.T) {
	// Arrange: Ensure global initialization
	resetRegistry()

	// Act: Fetch the global registry
	reg := GetMutexRegistry()

	// Assert
	if reg == nil {
		t.Error("expected a valid MutexRegistry from GetMutexRegistry, got nil")
	}
}

func TestMutexRegistry_RegisterAndHasMutex(t *testing.T) {
	// Arrange
	resetRegistry()
	reg := GetMutexRegistry()
	key := "test-mutex"
	mutex := NewCancellableMutex(key)

	// Act: Register the mutex
	err := reg.Register(mutex)

	// Assert
	if err != nil {
		t.Errorf("expected no error when registering new mutex, got %v", err)
	}

	// Act: Verify the mutex existence
	if !reg.HasMutex(key) {
		t.Errorf("expected registry to have mutex with key %q, but it did not", key)
	}

	// Act: Attempt duplicate registration
	err = reg.Register(mutex)

	// Assert Duplicate Registration
	if !errors.Is(err, AlreadyRegisteredError) {
		t.Errorf("expected error when re-registering a mutex, got %v", err)
	}
}

func TestMutexRegistry_GetMutex(t *testing.T) {
	// Arrange
	resetRegistry()
	reg := GetMutexRegistry()
	key := "test-mutex"
	mutex := NewCancellableMutex(key)

	// Arrange: Register the mutex
	err := reg.Register(mutex)
	if err != nil {
		t.Fatalf("unexpected error during registration: %v", err)
	}

	// Act: Retrieve the mutex by key
	optMutex := reg.GetMutex(key)

	// Assert
	value, some := optMutex.Value()
	if !some {
		t.Errorf("expected Some[CancellableMutex] from GetMutex, got None")
	}

	if value.GetKey() != key {
		t.Errorf("expected to get mutex with key %q, got %q", key, value.GetKey())
	}
}

func TestMutexRegistry_GetMutex_NotFound(t *testing.T) {
	// Arrange
	resetRegistry()
	reg := GetMutexRegistry()
	nonExistentKey := "non-existent-key"

	// Act: Try to get a mutex with a key that doesn’t exist
	optMutex := reg.GetMutex(nonExistentKey)

	// Assert
	_, some := optMutex.Value()
	if some {
		t.Error("expected GetMutex to return None for a non-existent key, but got Some")
	}
}

func TestMutexRegistry_RegisterAndRetrieveMultipleKeys(t *testing.T) {
	// Arrange
	resetRegistry()
	reg := GetMutexRegistry()
	keys := []string{"mutex-1", "mutex-2", "mutex-3"}

	// Act: Register multiple mutexes
	for _, key := range keys {
		mutex := NewCancellableMutex(key)
		err := reg.Register(mutex)
		if err != nil {
			t.Fatalf("unexpected error during registration of key %q: %v", key, err)
		}
	}

	// Assert: Retrieve and verify each mutex
	for _, key := range keys {
		// Act: Check if mutex exists
		if !reg.HasMutex(key) {
			t.Errorf("expected registry to have mutex with key %q, but it did not", key)
		}

		// Act: Fetch the mutex by key
		maybeMutex := reg.GetMutex(key)
		mutex, some := maybeMutex.Value()

		// Assert: Verify mutex details
		if !some {
			t.Errorf("expected Some[CancellableMutex] for key %q, got None", key)
		}
		if mutex.GetKey() != key {
			t.Errorf("expected mutex with key %q, got key %q", key, mutex.GetKey())
		}
	}
}
