package mutex

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewCancellableMutex_NewInstance(t *testing.T) {
	// Arrange
	key := "test-mutex"

	// Act
	mutex := NewCancellableMutex(key)

	// Assert
	if mutex.GetKey() != key {
		t.Errorf("expected key to be %q, got %q", key, mutex.GetKey())
	}

}

func TestNewCancellableMutex_ExistingInstance(t *testing.T) {
	//reset
	resetRegistry()

	// Arrange
	key := "test-mutex"
	mutex1 := GetOrNewCancellableMutex(key)

	// Act
	mutex2 := GetOrNewCancellableMutex(key)

	// Assert
	if mutex1 != mutex2 {
		t.Error("expected to retrieve the same instance of mutex for the same key")
	}
}

func TestCancellableMutex_LockAndUnlock(t *testing.T) {
	// Arrange
	key := "test-mutex"
	mutex := NewCancellableMutex(key)
	ctx := context.Background()

	// Act
	err := mutex.Lock(ctx) // Acquire lock
	if err != nil {
		t.Fatalf("expected no error when locking mutex, got %v", err)
	}

	// Assert after locking
	contextWithTimeoutLocked, cancelLocked := context.WithTimeout(ctx, time.Millisecond)
	defer cancelLocked()
	err = mutex.Lock(contextWithTimeoutLocked)
	if err == nil {
		t.Error("Lock should fail when mutex is already locked")
	}
	mutex.Unlock() // Release lock

	// Assert after unlocking
	contextWithTimeoutUnLocked, cancelUnLocked := context.WithTimeout(ctx, time.Millisecond)
	defer cancelUnLocked()
	err = mutex.Lock(contextWithTimeoutUnLocked)
	if err != nil {
		t.Errorf("expected no error when locking mutex, got %v", err)
	}
}

func TestCancellableMutex_LockWithContextCancel(t *testing.T) {
	//reset
	resetRegistry()

	// Arrange
	key := "test-mutex"
	mutex := NewCancellableMutex(key)
	otherCtx := context.Background()

	// Lock mutex to simulate it already being locked
	err := mutex.Lock(otherCtx)
	if err != nil {
		t.Fatalf("expected no error when pre-locking mutex, got %v", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Act
	err = mutex.Lock(ctx)

	// Assert
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context deadline exceeded error, got %v", err)
	}

	// Cleanup
	mutex.Unlock()
}

func TestCancellableMutex_MultipleLocks(t *testing.T) {
	//reset
	resetRegistry()

	// Arrange
	key := "test-mutex"
	mutex := NewCancellableMutex(key)

	// Act
	ctx := context.Background()

	// First lock
	err := mutex.Lock(ctx)
	if err != nil {
		t.Errorf("expected no error on first lock, got %v", err)
	}

	// Simulate unlocking and relocking
	mutex.Unlock()

	err = mutex.Lock(ctx)
	if err != nil {
		t.Errorf("expected no error on relock, got %v", err)
	}

	// Assert
	contextWithTimeoutLocked, cancelLocked := context.WithTimeout(ctx, time.Millisecond)
	defer cancelLocked()
	err = mutex.Lock(contextWithTimeoutLocked)
	if err == nil {
		t.Error("Lock should fail when mutex is already locked")
	}

	mutex.Unlock() // Cleanup the lock

}
