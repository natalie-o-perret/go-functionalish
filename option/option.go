// Package option provides a type-safe Optional value inspired by F#'s Option type.
// It eliminates nil pointer errors by making the presence or absence of a value explicit.
package option

// Option represents a value that may or may not be present.
// Use Some to wrap a value, None to represent absence.
type Option[T any] struct {
	value T
	valid bool
}

// Some wraps a value in an Option.
func Some[T any](v T) Option[T] { return Option[T]{value: v, valid: true} }

// None returns an empty Option. Use var n Option[T] for a zero-value None without parentheses.
func None[T any]() Option[T] { return Option[T]{} }

// Empty is an alias for None. Use var n Option[T] for a zero-value None without parentheses.
func Empty[T any]() Option[T] { return None[T]() }

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

// -- type-transforming methods -------------------------------------------------

// Map applies fn to the value inside Some, returning None unchanged.
func (o Option[T]) Map[R any](fn func(T) R) Option[R] {
	if o.valid {
		return Some(fn(o.value))
	}
	return None[R]()
}

// Bind applies fn to the value inside Some, flattening the resulting Option.
func (o Option[T]) Bind[R any](fn func(T) Option[R]) Option[R] {
	if o.valid {
		return fn(o.value)
	}
	return None[R]()
}

// OrElse returns o if Some, otherwise calls fn and returns its result.
func (o Option[T]) OrElse(fn func() Option[T]) Option[T] {
	if o.valid {
		return o
	}
	return fn()
}

// DefaultWith returns the value if Some, otherwise calls fn lazily.
// Unlike UnwrapOr, the default is only computed if needed.
func (o Option[T]) DefaultWith(fn func() T) T {
	if o.valid {
		return o.value
	}
	return fn()
}

// Tee calls fn with the value if Some, returning o unchanged.
// Useful for logging or side effects in a pipeline.
func (o Option[T]) Tee(fn func(T)) Option[T] {
	if o.valid {
		fn(o.value)
	}
	return o
}

// TeeNone calls fn if None, returning o unchanged.
// Useful for logging or side effects on the absent path.
func (o Option[T]) TeeNone(fn func()) Option[T] {
	if !o.valid {
		fn()
	}
	return o
}

// ZipWith combines o with other using fn. Returns None if either is None.
func (o Option[T]) ZipWith[U, R any](other Option[U], fn func(T, U) R) Option[R] {
	if o.valid && other.valid {
		return Some(fn(o.value, other.value))
	}
	return None[R]()
}

// -- package-level functions (only where methods are impossible) ---------------

// Zip combines two Options into an Option of a pair. None if either is None.
// Note: Zip cannot be a method because returning Option[Pair[T,U]] would create
// an instantiation cycle in the type checker.
func Zip[T, U any](a Option[T], b Option[U]) Option[Pair[T, U]] {
	if a.valid && b.valid {
		return Some(Pair[T, U]{First: a.value, Second: b.value})
	}
	return None[Pair[T, U]]()
}

// Flatten unwraps a nested Option[Option[T]] into Option[T].
// Note: Flatten cannot be a method because the receiver would need to be Option[Option[T]],
// which cannot be expressed in Go's type system.
func Flatten[T any](o Option[Option[T]]) Option[T] {
	if o.valid {
		return o.value
	}
	return None[T]()
}

// Contains reports whether o is Some and its value equals v.
// Note: Contains cannot be a method because Option[T] uses the unconstrained T any;
// the comparable constraint required here cannot be added on the method alone.
func Contains[T comparable](o Option[T], v T) bool {
	return o.valid && o.value == v
}

// Pair holds two values of potentially different types.
type Pair[T, U any] struct {
	First  T
	Second U
}

// PairOf constructs a Pair with positional arguments.
func PairOf[T, U any](first T, second U) Pair[T, U] { return Pair[T, U]{First: first, Second: second} }
