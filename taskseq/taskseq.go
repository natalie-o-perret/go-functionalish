// Package taskseq provides lazy, sequential, context-aware sequences for work
// that can wait on I/O or fail.
//
// Task sequences do not start goroutines or materialise their input. Each
// element is produced only when a terminal operation requests it.
// If an operation and source cleanup both fail, the returned error joins both
// errors and supports inspection with [errors.Is] and [errors.As].
package taskseq

import (
	"context"
	"errors"
)

// Seq is a lazy sequence that returns a terminal error.
//
// Implementations must call yield synchronously and stop when it returns false.
// The context is supplied when the sequence is consumed, allowing the same
// pipeline to be reused with different deadlines.
type Seq[T any] func(context.Context, func(T) bool) error

func joinErrors(ctx context.Context, errs ...error) error {
	err := errors.Join(errs...)
	if cause := context.Cause(ctx); cause != nil && !errors.Is(err, cause) {
		return errors.Join(err, cause)
	}
	return err
}

// FromSeq lifts a synchronous range-over-function sequence into a task sequence.
// It accepts seq.Seq and iter.Seq values without conversion.
// Cancellation is checked between elements; it cannot interrupt a synchronous
// source that is blocked while producing its next element.
func FromSeq[T any](source func(func(T) bool)) Seq[T] {
	return func(ctx context.Context, yield func(T) bool) error {
		if err := context.Cause(ctx); err != nil {
			return err
		}
		for value := range source {
			if err := context.Cause(ctx); err != nil {
				return err
			}
			if !yield(value) {
				return nil
			}
			if err := context.Cause(ctx); err != nil {
				return err
			}
		}
		return context.Cause(ctx)
	}
}

// MapAsync lazily maps each element with a context-aware operation.
// Elements are processed sequentially and in source order.
func (s Seq[T]) MapAsync[R any](fn func(context.Context, T) (R, error)) Seq[R] {
	return func(ctx context.Context, yield func(R) bool) error {
		if err := context.Cause(ctx); err != nil {
			return err
		}

		var operationErr error
		upstreamErr := s(ctx, func(value T) bool {
			if context.Cause(ctx) != nil {
				return false
			}

			mapped, err := fn(ctx, value)
			if err != nil {
				operationErr = err
				return false
			}
			if context.Cause(ctx) != nil {
				return false
			}
			if !yield(mapped) {
				return false
			}
			if context.Cause(ctx) != nil {
				return false
			}
			return true
		})
		return joinErrors(ctx, operationErr, upstreamErr)
	}
}

// FilterAsync lazily keeps elements accepted by a context-aware predicate.
// Elements are tested sequentially and in source order.
func (s Seq[T]) FilterAsync(fn func(context.Context, T) (bool, error)) Seq[T] {
	return func(ctx context.Context, yield func(T) bool) error {
		if err := context.Cause(ctx); err != nil {
			return err
		}

		var operationErr error
		upstreamErr := s(ctx, func(value T) bool {
			if context.Cause(ctx) != nil {
				return false
			}

			keep, err := fn(ctx, value)
			if err != nil {
				operationErr = err
				return false
			}
			if context.Cause(ctx) != nil {
				return false
			}
			if !keep {
				return true
			}
			if !yield(value) {
				return false
			}
			if context.Cause(ctx) != nil {
				return false
			}
			return true
		})
		return joinErrors(ctx, operationErr, upstreamErr)
	}
}

// Take lazily yields at most the first n elements.
// Non-positive n yields nothing without consuming the source.
func (s Seq[T]) Take(n int) Seq[T] {
	if n <= 0 {
		return func(ctx context.Context, _ func(T) bool) error {
			return context.Cause(ctx)
		}
	}

	return func(ctx context.Context, yield func(T) bool) error {
		if err := context.Cause(ctx); err != nil {
			return err
		}

		remaining := n
		upstreamErr := s(ctx, func(value T) bool {
			if context.Cause(ctx) != nil {
				return false
			}

			remaining--
			if !yield(value) {
				return false
			}
			if context.Cause(ctx) != nil {
				return false
			}
			return remaining > 0
		})
		return joinErrors(ctx, upstreamErr)
	}
}

// ForEach consumes the sequence, calling fn once for each element.
// It stops at the first source, cancellation, or callback error.
func (s Seq[T]) ForEach(ctx context.Context, fn func(context.Context, T) error) error {
	if err := context.Cause(ctx); err != nil {
		return err
	}

	var operationErr error
	upstreamErr := s(ctx, func(value T) bool {
		if context.Cause(ctx) != nil {
			return false
		}

		if err := fn(ctx, value); err != nil {
			operationErr = err
			return false
		}
		if context.Cause(ctx) != nil {
			return false
		}
		return true
	})
	return joinErrors(ctx, operationErr, upstreamErr)
}

// ToSlice consumes the sequence and returns the values produced before any
// error together with that error.
func (s Seq[T]) ToSlice(ctx context.Context) ([]T, error) {
	var values []T
	err := s.ForEach(ctx, func(_ context.Context, value T) error {
		values = append(values, value)
		return nil
	})
	return values, err
}
