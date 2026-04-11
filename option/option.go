// Package option provides a type-safe Optional value inspired by F#'s Option type.
// It eliminates nil pointer errors by making the presence or absence of a value explicit.
package option

// Option[T] represents a value that may or may not be present.
// Use Some to wrap a value, None to represent absence.
type Option[T any] struct {
	value T
	valid bool
}

// Some wraps a value in an Option.
func Some[T any](v T) Option[T] { return Option[T]{value: v, valid: true} }

// None returns an empty Option.
func None[T any]() Option[T] { return Option[T]{} }

// IsSome reports whether the Option contains a value.
func (o Option[T]) IsSome() bool { return o.valid }

// IsNone reports whether the Option is empty.
func (o Option[T]) IsNone() bool { return !o.valid }

// Unwrap returns the value. Panics if the Option is None.
func (o Option[T]) Unwrap() T {
	if !o.valid {
		panic("option: Unwrap called on None")
	}
	return o.value
}

// UnwrapOr returns the value, or def if None.
func (o Option[T]) UnwrapOr(def T) T {
	if o.valid {
		return o.value
	}
	return def
}

// UnwrapOrElse returns the value, or calls fn if None.
func (o Option[T]) UnwrapOrElse(fn func() T) T {
	if o.valid {
		return o.value
	}
	return fn()
}

// Filter returns Some if the predicate holds, otherwise None.
func (o Option[T]) Filter(fn func(T) bool) Option[T] {
	if o.valid && fn(o.value) {
		return o
	}
	return None[T]()
}

// ToSlice returns a slice of one element if Some, or an empty slice if None.
func (o Option[T]) ToSlice() []T {
	if o.valid {
		return []T{o.value}
	}
	return []T{}
}

// ── type-transforming package-level functions ─────────────────────────────────
// These must be package-level because Go methods cannot introduce new type parameters.

// Map applies fn to the value inside Some, returning None unchanged.
func Map[T, R any](o Option[T], fn func(T) R) Option[R] {
	if o.valid {
		return Some(fn(o.value))
	}
	return None[R]()
}

// FlatMap applies fn to the value inside Some, flattening the resulting Option.
func FlatMap[T, R any](o Option[T], fn func(T) Option[R]) Option[R] {
	if o.valid {
		return fn(o.value)
	}
	return None[R]()
}

// Zip combines two Options into an Option of a pair. None if either is None.
func Zip[T, U any](a Option[T], b Option[U]) Option[Pair[T, U]] {
	if a.valid && b.valid {
		return Some(Pair[T, U]{First: a.value, Second: b.value})
	}
	return None[Pair[T, U]]()
}

// Pair holds two values of potentially different types.
type Pair[T, U any] struct {
	First  T
	Second U
}

