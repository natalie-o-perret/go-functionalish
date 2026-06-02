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

// Apply calls fn with the two fields of t and returns the result.
func (t T2[A, B]) Apply[R any](fn func(A, B) R) R { return fn(t.First, t.Second) }

// MapFirst transforms the First element, leaving Second unchanged.
func (t T2[A, B]) MapFirst[R any](fn func(A) R) T2[R, B] {
	return T2[R, B]{First: fn(t.First), Second: t.Second}
}

// MapSecond transforms the Second element, leaving First unchanged.
func (t T2[A, B]) MapSecond[R any](fn func(B) R) T2[A, R] {
	return T2[A, R]{First: t.First, Second: fn(t.Second)}
}

// Map transforms both elements of a T2 independently.
func (t T2[A, B]) Map[RA, RB any](fnA func(A) RA, fnB func(B) RB) T2[RA, RB] {
	return T2[RA, RB]{First: fnA(t.First), Second: fnB(t.Second)}
}

// Extend returns a T3 with c appended as the third element.
func (t T2[A, B]) Extend[C any](c C) T3[A, B, C] {
	return T3[A, B, C]{First: t.First, Second: t.Second, Third: c}
}

// Unpack returns the three fields as individual return values.
func (t T3[A, B, C]) Unpack() (A, B, C) { return t.First, t.Second, t.Third }

// Apply calls fn with the three fields of t and returns the result.
func (t T3[A, B, C]) Apply[R any](fn func(A, B, C) R) R {
	return fn(t.First, t.Second, t.Third)
}

// MapFirst transforms the First element, leaving Second and Third unchanged.
func (t T3[A, B, C]) MapFirst[R any](fn func(A) R) T3[R, B, C] {
	return T3[R, B, C]{First: fn(t.First), Second: t.Second, Third: t.Third}
}

// MapSecond transforms the Second element, leaving First and Third unchanged.
func (t T3[A, B, C]) MapSecond[R any](fn func(B) R) T3[A, R, C] {
	return T3[A, R, C]{First: t.First, Second: fn(t.Second), Third: t.Third}
}

// MapThird transforms the Third element, leaving First and Second unchanged.
func (t T3[A, B, C]) MapThird[R any](fn func(C) R) T3[A, B, R] {
	return T3[A, B, R]{First: t.First, Second: t.Second, Third: fn(t.Third)}
}

// Map transforms all three elements independently.
func (t T3[A, B, C]) Map[RA, RB, RC any](fnA func(A) RA, fnB func(B) RB, fnC func(C) RC) T3[RA, RB, RC] {
	return T3[RA, RB, RC]{First: fnA(t.First), Second: fnB(t.Second), Third: fnC(t.Third)}
}

// DropFirst returns a T2 containing only Second and Third.
func (t T3[A, B, C]) DropFirst() T2[B, C] { return T2[B, C]{First: t.Second, Second: t.Third} }

// DropSecond returns a T2 containing only First and Third.
func (t T3[A, B, C]) DropSecond() T2[A, C] { return T2[A, C]{First: t.First, Second: t.Third} }

// DropThird returns a T2 containing only First and Second.
func (t T3[A, B, C]) DropThird() T2[A, B] { return T2[A, B]{First: t.First, Second: t.Second} }

// Extend returns a T4 with d appended as the fourth element.
func (t T3[A, B, C]) Extend[D any](d D) T4[A, B, C, D] {
	return T4[A, B, C, D]{First: t.First, Second: t.Second, Third: t.Third, Fourth: d}
}

// Unpack returns the four fields as individual return values.
func (t T4[A, B, C, D]) Unpack() (A, B, C, D) { return t.First, t.Second, t.Third, t.Fourth }

// Apply calls fn with the four fields of t and returns the result.
func (t T4[A, B, C, D]) Apply[R any](fn func(A, B, C, D) R) R {
	return fn(t.First, t.Second, t.Third, t.Fourth)
}

// MapFirst transforms the First element, leaving the others unchanged.
func (t T4[A, B, C, D]) MapFirst[R any](fn func(A) R) T4[R, B, C, D] {
	return T4[R, B, C, D]{First: fn(t.First), Second: t.Second, Third: t.Third, Fourth: t.Fourth}
}

// MapSecond transforms the Second element, leaving the others unchanged.
func (t T4[A, B, C, D]) MapSecond[R any](fn func(B) R) T4[A, R, C, D] {
	return T4[A, R, C, D]{First: t.First, Second: fn(t.Second), Third: t.Third, Fourth: t.Fourth}
}

// MapThird transforms the Third element, leaving the others unchanged.
func (t T4[A, B, C, D]) MapThird[R any](fn func(C) R) T4[A, B, R, D] {
	return T4[A, B, R, D]{First: t.First, Second: t.Second, Third: fn(t.Third), Fourth: t.Fourth}
}

// MapFourth transforms the Fourth element, leaving the others unchanged.
func (t T4[A, B, C, D]) MapFourth[R any](fn func(D) R) T4[A, B, C, R] {
	return T4[A, B, C, R]{First: t.First, Second: t.Second, Third: t.Third, Fourth: fn(t.Fourth)}
}

// Map transforms all four elements independently.
func (t T4[A, B, C, D]) Map[RA, RB, RC, RD any](fnA func(A) RA, fnB func(B) RB, fnC func(C) RC, fnD func(D) RD) T4[RA, RB, RC, RD] {
	return T4[RA, RB, RC, RD]{First: fnA(t.First), Second: fnB(t.Second), Third: fnC(t.Third), Fourth: fnD(t.Fourth)}
}

// DropFirst returns a T3 containing only Second, Third, and Fourth.
func (t T4[A, B, C, D]) DropFirst() T3[B, C, D] {
	return T3[B, C, D]{First: t.Second, Second: t.Third, Third: t.Fourth}
}

// DropSecond returns a T3 containing only First, Third, and Fourth.
func (t T4[A, B, C, D]) DropSecond() T3[A, C, D] {
	return T3[A, C, D]{First: t.First, Second: t.Third, Third: t.Fourth}
}

// DropThird returns a T3 containing only First, Second, and Fourth.
func (t T4[A, B, C, D]) DropThird() T3[A, B, D] {
	return T3[A, B, D]{First: t.First, Second: t.Second, Third: t.Fourth}
}

// DropFourth returns a T3 containing only First, Second, and Third.
func (t T4[A, B, C, D]) DropFourth() T3[A, B, C] {
	return T3[A, B, C]{First: t.First, Second: t.Second, Third: t.Third}
}

// -- package-level functions (kept for backward compatibility) -----------------
// Prefer the method forms: t.Apply(fn), t.MapFirst(fn), t.MapSecond(fn), t.Map(fnA, fnB).

// Fst returns the first element of a T2. Useful as a first-class function value.
func Fst[A, B any](t T2[A, B]) A { return t.First }

// Snd returns the second element of a T2. Useful as a first-class function value.
func Snd[A, B any](t T2[A, B]) B { return t.Second }

// Thd returns the third element of a T3. Useful as a first-class function value.
func Thd[A, B, C any](t T3[A, B, C]) C { return t.Third }

// Fth returns the fourth element of a T4. Useful as a first-class function value.
func Fth[A, B, C, D any](t T4[A, B, C, D]) D { return t.Fourth }

// Map2 maps both elements of a homogeneous T2[A, A] with a single function.
// For heterogeneous tuples use the Map method, which takes one function per field.
func Map2[A, R any](t T2[A, A], fn func(A) R) T2[R, R] {
	return T2[R, R]{First: fn(t.First), Second: fn(t.Second)}
}

// Map3 maps all three elements of a homogeneous T3[A, A, A] with a single function.
func Map3[A, R any](t T3[A, A, A], fn func(A) R) T3[R, R, R] {
	return T3[R, R, R]{First: fn(t.First), Second: fn(t.Second), Third: fn(t.Third)}
}

// Map4 maps all four elements of a homogeneous T4[A, A, A, A] with a single function.
func Map4[A, R any](t T4[A, A, A, A], fn func(A) R) T4[R, R, R, R] {
	return T4[R, R, R, R]{First: fn(t.First), Second: fn(t.Second), Third: fn(t.Third), Fourth: fn(t.Fourth)}
}

// Curry converts a binary function (A, B) -> C into a curried form A -> B -> C.
func Curry[A, B, C any](fn func(A, B) C) func(A) func(B) C {
	return func(a A) func(B) C {
		return func(b B) C { return fn(a, b) }
	}
}

// Uncurry converts a curried function A -> B -> C into a binary function (A, B) -> C.
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

// Pair holds two values of potentially different types.
// Prefer this over option.Pair for non-Option use cases.
type Pair[T, U any] struct {
	First  T
	Second U
}

// PairOf constructs a Pair with positional arguments.
func PairOf[T, U any](first T, second U) Pair[T, U] { return Pair[T, U]{First: first, Second: second} }
