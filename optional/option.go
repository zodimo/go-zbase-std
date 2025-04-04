package optional

//@see https://github.com/AngusGMorrison/typedd-gophers-talk
import (
	"github.com/zodimo/go-zstd/complete"
)

type Option[T any] struct {
	value T
	some  bool
}

func None[T any]() Option[T] {
	return Option[T]{}
}

func Some[T any](value T) Option[T] {
	return Option[T]{
		value: value,
		some:  true,
	}
}

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

func (o *Option[T]) Value() (T, bool) {
	return o.value, o.some
}

func partiallyComplete(maybePartial complete.Complete) bool {
	// Check for incomplete values
	if maybePartial == nil {
		return true
	}
	return !maybePartial.Complete()
}
