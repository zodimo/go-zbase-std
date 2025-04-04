package complete

import (
	"errors"
	"testing"
)

// Mock implementation of the Complete interface
type MockComplete struct {
	isComplete bool
}

// Complete implementation for MockComplete
func (m MockComplete) Complete() bool {
	return m.isComplete
}

func TestValidateCompleteness_AllComplete(t *testing.T) {
	// Arrange
	c1 := MockComplete{isComplete: true}
	c2 := MockComplete{isComplete: true}

	// Act
	err := ValidateCompleteness(c1, c2)

	// Assert
	if err != nil {
		t.Errorf("expected no error, but got: %v", err)
	}
}

func TestValidateCompleteness_Incomplete(t *testing.T) {
	// Arrange
	c1 := MockComplete{isComplete: true}
	c2 := MockComplete{isComplete: false}

	// Act
	err := ValidateCompleteness(c1, c2)

	// Assert
	if err == nil {
		t.Error("expected an error, but got nil")
	}

	var incompleteError *IncompleteTypeError
	if !errors.As(err, &incompleteError) {
		t.Errorf("expected error of type *IncompleteTypeError, but got: %T", err)
	}
}

func TestIncompleteTypeError_ErrorMethod(t *testing.T) {
	// Arrange
	incomplete := MockComplete{isComplete: false}
	err := &IncompleteTypeError{Incomplete: incomplete}

	// Act
	got := err.Error()
	expected := "value of type complete.MockComplete implements Complete but was incomplete: complete.MockComplete{isComplete:false}"

	// Assert
	if got != expected {
		t.Errorf("Error() = %q; want %q", got, expected)
	}
}
