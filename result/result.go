// Package result provides a type-safe Result type inspired by F# and Rust.
// It represents either a successful value (Ok) or an error (Err),
// enabling railway-oriented error handling without explicit nil checks.
package result

import "github.com/natalie-o-perret/go-functionalish/option"

// Result[T, E] holds either a successful value of type T or an error of type E.
type Result[T, E any] struct {
	value T
	err   E
	ok    bool
}

// Ok wraps a successful value.
func Ok[T, E any](v T) Result[T, E] { return Result[T, E]{value: v, ok: true} }

// Err wraps an error value.
func Err[T, E any](e E) Result[T, E] { return Result[T, E]{err: e, ok: false} }

// Try calls fn and wraps the returned (T, error) pair into a Result[T, error].
func Try[T any](fn func() (T, error)) Result[T, error] {
	v, err := fn()
	if err != nil {
		return Err[T, error](err)
	}
	return Ok[T, error](v)
}

// IsOk reports whether the Result is Ok.
func (r Result[T, E]) IsOk() bool { return r.ok }

// IsErr reports whether the Result is Err.
func (r Result[T, E]) IsErr() bool { return !r.ok }

// Unwrap returns the value. Panics if Err.
func (r Result[T, E]) Unwrap() T {
	if !r.ok {
		panic("result: Unwrap called on Err")
	}
	return r.value
}

// UnwrapOr returns the value, or def if Err.
func (r Result[T, E]) UnwrapOr(def T) T {
	if r.ok {
		return r.value
	}
	return def
}

// UnwrapOrElse returns the value, or calls fn with the error if Err.
func (r Result[T, E]) UnwrapOrElse(fn func(E) T) T {
	if r.ok {
		return r.value
	}
	return fn(r.err)
}

// UnwrapErr returns the error. Panics if Ok.
func (r Result[T, E]) UnwrapErr() E {
	if r.ok {
		panic("result: UnwrapErr called on Ok")
	}
	return r.err
}

// ToOption converts Ok to Some, Err to None.
func (r Result[T, E]) ToOption() option.Option[T] {
	if r.ok {
		return option.Some(r.value)
	}
	return option.None[T]()
}

// -- type-transforming package-level functions ---------------------------------

// Map applies fn to the Ok value, passing Err unchanged.
func Map[T, U, E any](r Result[T, E], fn func(T) U) Result[U, E] {
	if r.ok {
		return Ok[U, E](fn(r.value))
	}
	return Err[U, E](r.err)
}

// MapErr applies fn to the Err value, passing Ok unchanged.
func MapErr[T, E, F any](r Result[T, E], fn func(E) F) Result[T, F] {
	if !r.ok {
		return Err[T, F](fn(r.err))
	}
	return Ok[T, F](r.value)
}

// Bind applies fn to the Ok value, flattening the resulting Result.
func Bind[T, U, E any](r Result[T, E], fn func(T) Result[U, E]) Result[U, E] {
	if r.ok {
		return fn(r.value)
	}
	return Err[U, E](r.err)
}

// FromOption converts Some to Ok and None to Err using the provided error.
func FromOption[T, E any](o option.Option[T], errIfNone E) Result[T, E] {
	if o.IsSome() {
		return Ok[T, E](o.Unwrap())
	}
	return Err[T, E](errIfNone)
}

// Zip combines two Results into a Result of a pair.
// Returns the first Err encountered if either is Err.
func Zip[T, U, E any](a Result[T, E], b Result[U, E]) Result[option.Pair[T, U], E] {
	if a.ok && b.ok {
		return Ok[option.Pair[T, U], E](option.Pair[T, U]{First: a.value, Second: b.value})
	}
	if !a.ok {
		return Err[option.Pair[T, U], E](a.err)
	}
	return Err[option.Pair[T, U], E](b.err)
}

// Flatten unwraps a nested Result[Result[T,E],E] into Result[T,E].
func Flatten[T, E any](r Result[Result[T, E], E]) Result[T, E] {
	if r.ok {
		return r.value
	}
	return Err[T, E](r.err)
}

// OrElse returns r if Ok, otherwise calls fn with the error to produce a fallback Result.
func OrElse[T, E any](r Result[T, E], fn func(E) Result[T, E]) Result[T, E] {
	if r.ok {
		return r
	}
	return fn(r.err)
}

// Tee calls fn with the Ok value as a side effect and returns r unchanged.
// Useful for logging or metrics in a pipeline without breaking the chain.
func Tee[T, E any](r Result[T, E], fn func(T)) Result[T, E] {
	if r.ok {
		fn(r.value)
	}
	return r
}

// TeeErr calls fn with the Err value as a side effect and returns r unchanged.
func TeeErr[T, E any](r Result[T, E], fn func(E)) Result[T, E] {
	if !r.ok {
		fn(r.err)
	}
	return r
}

// Map2 combines two Results with fn. Short-circuits on the first Err.
func Map2[A, B, R, E any](ra Result[A, E], rb Result[B, E], fn func(A, B) R) Result[R, E] {
	if ra.ok && rb.ok {
		return Ok[R, E](fn(ra.value, rb.value))
	}
	if !ra.ok {
		return Err[R, E](ra.err)
	}
	return Err[R, E](rb.err)
}

// Contains reports whether r is Ok and its value equals v.
func Contains[T comparable, E any](r Result[T, E], v T) bool {
	return r.ok && r.value == v
}

// Sequence converts a slice of Results into a Result of a slice.
// Short-circuits on the first Err.
func Sequence[T, E any](rs []Result[T, E]) Result[[]T, E] {
	values := make([]T, 0, len(rs))
	for _, r := range rs {
		if !r.ok {
			return Err[[]T, E](r.err)
		}
		values = append(values, r.value)
	}
	return Ok[[]T, E](values)
}

// Traverse applies fn to each item and sequences the results.
// Short-circuits on the first Err.
func Traverse[T, R, E any](items []T, fn func(T) Result[R, E]) Result[[]R, E] {
	values := make([]R, 0, len(items))
	for _, item := range items {
		r := fn(item)
		if !r.ok {
			return Err[[]R, E](r.err)
		}
		values = append(values, r.value)
	}
	return Ok[[]R, E](values)
}
