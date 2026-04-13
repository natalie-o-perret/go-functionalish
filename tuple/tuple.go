// Package tuple provides generic immutable tuple types T2, T3 and T4,
// inspired by F# native tuple syntax.
// Tuples are value types; all operations return new values.
package tuple

// T2 holds two values of potentially different types.
type T2[A, B any] struct {
	First  A
	Second B
}

// T3 holds three values of potentially different types.
type T3[A, B, C any] struct {
	First  A
	Second B
	Third  C
}

// T4 holds four values of potentially different types.
type T4[A, B, C, D any] struct {
	First  A
	Second B
	Third  C
	Fourth D
}

// Of creates a T2 from two values.
func Of[A, B any](a A, b B) T2[A, B] { return T2[A, B]{First: a, Second: b} }

// Of3 creates a T3 from three values.
func Of3[A, B, C any](a A, b B, c C) T3[A, B, C] {
	return T3[A, B, C]{First: a, Second: b, Third: c}
}

// Of4 creates a T4 from four values.
func Of4[A, B, C, D any](a A, b B, c C, d D) T4[A, B, C, D] {
	return T4[A, B, C, D]{First: a, Second: b, Third: c, Fourth: d}
}

// Unpack returns the two fields as individual return values.
func (t T2[A, B]) Unpack() (A, B) { return t.First, t.Second }

// Swap returns a new T2 with First and Second exchanged.
func (t T2[A, B]) Swap() T2[B, A] { return T2[B, A]{First: t.Second, Second: t.First} }

// Unpack returns the three fields as individual return values.
func (t T3[A, B, C]) Unpack() (A, B, C) { return t.First, t.Second, t.Third }

// Unpack returns the four fields as individual return values.
func (t T4[A, B, C, D]) Unpack() (A, B, C, D) { return t.First, t.Second, t.Third, t.Fourth }

// Apply calls fn with the two fields of t and returns the result.
func Apply[A, B, R any](t T2[A, B], fn func(A, B) R) R { return fn(t.First, t.Second) }

// Apply3 calls fn with the three fields of t and returns the result.
func Apply3[A, B, C, R any](t T3[A, B, C], fn func(A, B, C) R) R {
	return fn(t.First, t.Second, t.Third)
}

// Apply4 calls fn with the four fields of t and returns the result.
func Apply4[A, B, C, D, R any](t T4[A, B, C, D], fn func(A, B, C, D) R) R {
	return fn(t.First, t.Second, t.Third, t.Fourth)
}

// MapFirst transforms the First element of a T2, leaving Second unchanged.
func MapFirst[A, B, R any](t T2[A, B], fn func(A) R) T2[R, B] {
	return T2[R, B]{First: fn(t.First), Second: t.Second}
}

// MapSecond transforms the Second element of a T2, leaving First unchanged.
func MapSecond[A, B, R any](t T2[A, B], fn func(B) R) T2[A, R] {
	return T2[A, R]{First: t.First, Second: fn(t.Second)}
}

// Map transforms both elements of a T2 independently.
func Map[A, B, RA, RB any](t T2[A, B], fnA func(A) RA, fnB func(B) RB) T2[RA, RB] {
	return T2[RA, RB]{First: fnA(t.First), Second: fnB(t.Second)}
}

// Curry converts a binary function (A, B) → C into a curried form A → B → C.
func Curry[A, B, C any](fn func(A, B) C) func(A) func(B) C {
	return func(a A) func(B) C {
		return func(b B) C { return fn(a, b) }
	}
}

// Uncurry converts a curried function A → B → C into a binary function (A, B) → C.
func Uncurry[A, B, C any](fn func(A) func(B) C) func(A, B) C {
	return func(a A, b B) C { return fn(a)(b) }
}

// FromFunc2 lifts a binary function into a T2-consuming function.
func FromFunc2[A, B, C any](fn func(A, B) C) func(T2[A, B]) C {
	return func(t T2[A, B]) C { return fn(t.First, t.Second) }
}

// ToFunc2 converts a T2-consuming function into a binary function.
func ToFunc2[A, B, C any](fn func(T2[A, B]) C) func(A, B) C {
	return func(a A, b B) C { return fn(T2[A, B]{First: a, Second: b}) }
}
