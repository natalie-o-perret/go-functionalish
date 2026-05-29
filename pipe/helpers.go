package pipe

// -- arithmetic helpers -------------------------------------------------------

// Inc returns v + 1. Wraps on integer overflow per the Go spec.
func Inc[T Number](v T) T { return v + 1 }

// Dec returns v - 1. Wraps on unsigned underflow per the Go spec.
func Dec[T Number](v T) T { return v - 1 }

// Add returns a function that adds n to its argument.
//
//	seq.Range(1, 6).Map(pipe.Add(10)).ToSlice() // => [11 12 13 14 15]
func Add[T Number](n T) func(T) T { return func(v T) T { return v + n } }

// Sub returns a function that subtracts n from its argument.
func Sub[T Number](n T) func(T) T { return func(v T) T { return v - n } }

// Mul returns a function that multiplies its argument by n.
func Mul[T Number](n T) func(T) T { return func(v T) T { return v * n } }

// Div returns a function that divides its argument by n.
// Integer division by zero panics; float division by zero returns ±Inf or NaN.
func Div[T Number](n T) func(T) T { return func(v T) T { return v / n } }

// Negate returns -v. Defined for signed integers and floats only.
func Negate[T Signed](v T) T { return -v }

// Abs returns the absolute value of v. Defined for signed integers and floats.
// Note: Abs[int8](-128) overflows to -128 per Go's two's complement arithmetic.
func Abs[T Signed](v T) T {
	if v < 0 {
		return -v
	}
	return v
}

// Clamp returns a function that clamps its argument to [lo, hi].
//
//	pipe.Clamp(0, 100)(150) // => 100
func Clamp[T Number](lo, hi T) func(T) T {
	return func(v T) T {
		if v < lo {
			return lo
		}
		if v > hi {
			return hi
		}
		return v
	}
}

// -- predicate combinators ----------------------------------------------------

// And returns a predicate that is true only when every given predicate is true.
// Returns true when called with zero predicates (vacuous truth).
//
//	isShortEven := pipe.And(pipe.IsEven[int], func(n int) bool { return n < 10 })
func And[T any](preds ...func(T) bool) func(T) bool {
	return func(v T) bool {
		for _, p := range preds {
			if !p(v) {
				return false
			}
		}
		return true
	}
}

// Or returns a predicate that is true when at least one given predicate is true.
// Returns false when called with zero predicates.
//
//	isEdge := pipe.Or(pipe.IsZero[int], func(n int) bool { return n == 100 })
func Or[T any](preds ...func(T) bool) func(T) bool {
	return func(v T) bool {
		for _, p := range preds {
			if p(v) {
				return true
			}
		}
		return false
	}
}

// -- general combinators ------------------------------------------------------

// Const returns a function that ignores its argument and always returns v.
// Useful as a stub or default in higher-order APIs.
//
//	seq.Range(1, 4).Map(pipe.Const[int, int](42)).ToSlice() // => [42 42 42]
func Const[T, U any](v T) func(U) T {
	return func(_ U) T { return v }
}

// Flip returns a new function with the two argument positions swapped.
//
//	divide := func(a, b float64) float64 { return a / b }
//	divideInto := pipe.Flip(divide) // divideInto(b, a) == a/b
func Flip[A, B, C any](f func(A, B) C) func(B, A) C {
	return func(b B, a A) C { return f(a, b) }
}

// On combines two values by first projecting each through project, then
// passing both results to combine. Useful for comparison and aggregation.
//
//	// sort words by length
//	slices.SortFunc(words, pipe.On(cmp.Compare, len))
func On[A, B, C any](combine func(B, B) C, project func(A) B) func(A, A) C {
	return func(a1, a2 A) C { return combine(project(a1), project(a2)) }
}

// Curry2 converts a two-argument function into a chain of one-argument functions.
// Curry2(f)(a)(b) == f(a, b).
//
//	add := pipe.Curry2(func(a, b int) int { return a + b })
//	add5 := add(5) // func(int) int
//	add5(3)        // => 8
func Curry2[A, B, C any](f func(A, B) C) func(A) func(B) C {
	return func(a A) func(B) C {
		return func(b B) C { return f(a, b) }
	}
}

// Uncurry2 is the inverse of Curry2: converts a curried function back into a
// two-argument function.
//
//	pipe.Uncurry2(pipe.Curry2(f))(a, b) == f(a, b)
func Uncurry2[A, B, C any](f func(A) func(B) C) func(A, B) C {
	return func(a A, b B) C { return f(a)(b) }
}

// Partial1 partially applies the first argument of a two-argument function,
// returning a one-argument function.
//
//	addFive := pipe.Partial1(func(a, b int) int { return a + b }, 5)
//	addFive(3) // => 8
func Partial1[A, B, C any](f func(A, B) C, a A) func(B) C {
	return func(b B) C { return f(a, b) }
}

// Partial2 partially applies the second argument of a two-argument function,
// returning a one-argument function.
//
//	double := pipe.Partial2(func(a, b int) int { return a * b }, 2)
//	double(7) // => 14
func Partial2[A, B, C any](f func(A, B) C, b B) func(A) C {
	return func(a A) C { return f(a, b) }
}
