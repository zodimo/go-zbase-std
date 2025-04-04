// Package optional provides a generic implementation of an optional type
// in Go, similar to optional types in other programming languages. It
// offers functionality to represent a value that may or may not exist.
package optional

import (
	"github.com/zodimo/go-zstd/complete"
)

// Option represents an optional value of a generic type T.
// It can either hold a value or indicate the absence of a value.
type Option[T any] struct {
	value T    // The value of type T.
	some  bool // Indicates whether the value is present.
}

// None initializes an Option without a value, representing the absence
// of a value.
//
// Example:
//
//	var noneOption = None[string]()
func None[T any]() Option[T] {
	return Option[T]{}
}

// Some initializes an Option with a given value.
//
// Example:
//
//	someOption := Some(42)
func Some[T any](value T) Option[T] {
	return Option[T]{
		value: value,
		some:  true,
	}
}

// SomeComplete initializes an Option with a given value, performing a check
// to ensure the value is "complete." If the value implements the
// complete.Complete interface and is found to be incomplete, an
// IncompleteTypeError is returned.
//
// This function is particularly useful when working with types that require
// additional validation.
//
// Parameters:
//   - value: The value of type T to be wrapped by the Option.
//
// Returns:
//   - Option[T]: If the value is valid and complete.
//   - error: If the value is incomplete and fails validation.
//
// Example:
//
//	validOption, err := SomeComplete(myCompleteTypeInstance)
func SomeComplete[T any](value T) (Option[T], error) {
	if c, ok := any(value).(complete.Complete); ok {
		if partiallyComplete(c) {
			return Option[T]{}, &complete.IncompleteTypeError{Incomplete: c}
		}
	}

	return Option[T]{
		value: value,
		some:  true,
	}, nil
}

// Value retrieves the wrapped value from the Option and a boolean
// to indicate whether the value is present.
//
// Returns:
//   - T: The contained value of type T.
//   - bool: True if the value is present, false if not.
//
// Example:
//
//	value, ok := option.Value()
func (o *Option[T]) Value() (T, bool) {
	return o.value, o.some
}

// partiallyComplete checks whether a value of type complete.Complete is
// incomplete. A value is considered incomplete if it is nil or its Complete()
// method returns false.
//
// Parameters:
//   - maybePartial: A value that implements the complete.Complete interface.
//
// Returns:
//   - bool: True if the value is incomplete, false otherwise.
func partiallyComplete(maybePartial complete.Complete) bool {
	// Check for incomplete values
	if maybePartial == nil {
		return true
	}
	return !maybePartial.Complete()
}
