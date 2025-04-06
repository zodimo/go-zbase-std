package optional

import (
	"errors"
	"reflect"
	"testing"

	"github.com/zodimo/go-zstd/complete"
)

// Mock type for testing Complete interface
type MockComplete struct {
	isComplete bool
}

// Complete implementation for MockComplete
func (m MockComplete) Complete() bool {
	return m.isComplete
}

func TestNone(t *testing.T) {
	// Act
	opt := None[int]()

	// Assert
	value, some := opt.Value()
	if some {
		t.Errorf("expected opt to be None, got Some with value %v", value)
	}

	if !reflect.DeepEqual(value, 0) {
		t.Errorf("expected None to return a zero-value, got %v", value)
	}
}

func TestSome(t *testing.T) {
	// Arrange
	value := "hello"

	// Act
	opt := Some(value)

	// Assert
	optValue, some := opt.Value()
	if !some {
		t.Errorf("expected opt to be Some, got None")
	}
	if optValue != value {
		t.Errorf("expected Some value to be %q, got %q", value, optValue)
	}
}

func TestSome_ValidValue(t *testing.T) {
	// Arrange
	value := 123

	// Act
	opt, err := SomeComplete(value)

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	optValue, some := opt.Value()
	if !some {
		t.Errorf("expected opt to be Some, got None")
	}
	if optValue != value {
		t.Errorf("expected Some value to be %v, got %v", value, optValue)
	}
}

func TestSome_IncompleteValue(t *testing.T) {
	// Arrange
	incomplete := MockComplete{isComplete: false}

	// Act
	_, err := SomeComplete(incomplete)

	// Assert
	if err == nil {
		t.Error("expected an error when passing an incomplete value, got nil")
	}

	var incompleteError *complete.IncompleteTypeError
	if !errors.As(err, &incompleteError) {
		t.Errorf("expected error of type *IncompleteTypeError, got %T", err)
	}
}

func TestSome_CompleteValue(t *testing.T) {
	// Arrange
	completeValue := MockComplete{isComplete: true}

	// Act
	opt, err := SomeComplete(completeValue)

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	optValue, some := opt.Value()
	if !some {
		t.Errorf("expected opt to be Some, got None")
	}
	if !reflect.DeepEqual(optValue, completeValue) {
		t.Errorf("expected Some value to be %+v, got %+v", completeValue, optValue)
	}
}

func TestOption_CompleteValue(t *testing.T) {
	// Arrange
	value := "test value"
	opt, _ := SomeComplete(value)

	// Act
	optValue, some := opt.Value()

	// Assert
	if !some {
		t.Errorf("expected opt to be Some, got None")
	}
	if optValue != value {
		t.Errorf("expected Some value to be %q, got %q", value, optValue)
	}
}

func TestPartiallyComplete_Partial(t *testing.T) {
	// Arrange
	partial := MockComplete{isComplete: false}

	// Act
	result := partiallyComplete(partial)

	// Assert
	if !result {
		t.Error("expected partiallyComplete to return true for an incomplete value, got false")
	}
}

func TestPartiallyComplete_FullyComplete(t *testing.T) {
	// Arrange
	completeValue := MockComplete{isComplete: true}

	// Act
	result := partiallyComplete(completeValue)

	// Assert
	if result {
		t.Error("expected partiallyComplete to return false for a complete value, got true")
	}
}

func TestPartiallyComplete_Nil(t *testing.T) {
	// Act
	result := partiallyComplete(nil)

	// Assert
	if !result {
		t.Error("expected partiallyComplete to return true for a nil value, got false")
	}
}

func TestPartiallyComplete_ZeroValue(t *testing.T) {
	// Arrange
	var zeroValue MockComplete

	// Act
	result := partiallyComplete(zeroValue)

	// Assert
	if !result {
		t.Error("expected partiallyComplete to return true for a zero-value, got false")
	}
}
