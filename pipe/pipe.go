// Package pipe provides F#-style pipe operators for Go.
// Since Go has no |> operator, Pipe functions thread a value left-to-right
// through a chain of single-argument functions.
//
// Example:
//
//	result := pipe.Pipe3(
//	    rawInput,
//	    strings.TrimSpace,
//	    strings.ToLower,
//	    validate,
//	)
package pipe

// Pipe2 threads v through f1 then f2.
func Pipe2[A, B, C any](v A, f1 func(A) B, f2 func(B) C) C {
	return f2(f1(v))
}

// Pipe3 threads v through f1 => f2 => f3.
func Pipe3[A, B, C, D any](v A, f1 func(A) B, f2 func(B) C, f3 func(C) D) D {
	return f3(f2(f1(v)))
}

// Pipe4 threads v through f1 => f2 => f3 => f4.
func Pipe4[A, B, C, D, E any](v A, f1 func(A) B, f2 func(B) C, f3 func(C) D, f4 func(D) E) E {
	return f4(f3(f2(f1(v))))
}

// Pipe5 threads v through f1 => ... => f5.
func Pipe5[A, B, C, D, E, F any](v A, f1 func(A) B, f2 func(B) C, f3 func(C) D, f4 func(D) E, f5 func(E) F) F {
	return f5(f4(f3(f2(f1(v)))))
}

// Pipe6 threads v through f1 => ... => f6.
func Pipe6[A, B, C, D, E, F, G any](v A, f1 func(A) B, f2 func(B) C, f3 func(C) D, f4 func(D) E, f5 func(E) F, f6 func(F) G) G {
	return f6(f5(f4(f3(f2(f1(v))))))
}

// Pipe7 threads v through f1 => ... => f7.
func Pipe7[A, B, C, D, E, F, G, H any](v A, f1 func(A) B, f2 func(B) C, f3 func(C) D, f4 func(D) E, f5 func(E) F, f6 func(F) G, f7 func(G) H) H {
	return f7(f6(f5(f4(f3(f2(f1(v)))))))
}

// Pipe8 threads v through f1 => ... => f8.
func Pipe8[A, B, C, D, E, F, G, H, I any](v A, f1 func(A) B, f2 func(B) C, f3 func(C) D, f4 func(D) E, f5 func(E) F, f6 func(F) G, f7 func(G) H, f8 func(H) I) I {
	return f8(f7(f6(f5(f4(f3(f2(f1(v))))))))
}

// PipeEndoN threads v through an arbitrary number of same-type (endomorphic) functions.
// All functions must share the same input and output type T.
//
// Naming convention:
//   - *EndoN variants (PipeEndoN, ComposeEndoN) : variadic, same-type only (T => T)
//   - *2..8 variants (Pipe2-Pipe8, Compose2-Compose4) : fixed-arity, type-changing (A => B => ...)
//
// For pipelines where the type changes between steps, use Pipe2-Pipe8
// with ComposeEndoN/Compose2-Compose4 to collapse multiple steps into one slot.
//
//nolint:revive // PipeEndoN intentionally includes the package prefix for discoverability.
func PipeEndoN[T any](v T, fns ...func(T) T) T {
	for _, fn := range fns {
		v = fn(v)
	}
	return v
}

// ComposeEndoN merges multiple same-type (endomorphic) transform functions into a single function.
// Useful for collapsing consecutive same-type pipeline steps into one pipe stage.
//
// Example:
//
//	pipe.Pipe3(
//	    input,
//	    pipe.ComposeEndoN(step1, step2, step3),  // all T => T
//	    transformType,                            // T => R
//	    ...
func ComposeEndoN[T any](fns ...func(T) T) func(T) T {
	return func(v T) T {
		for _, fn := range fns {
			v = fn(v)
		}
		return v
	}
}

// Compose2 composes two functions: A => B => C into A => C.
// Useful for collapsing type-changing steps into one pipe stage.
func Compose2[A, B, C any](f1 func(A) B, f2 func(B) C) func(A) C {
	return func(v A) C { return f2(f1(v)) }
}

// Compose3 composes three functions: A => B => C => D into A => D.
func Compose3[A, B, C, D any](f1 func(A) B, f2 func(B) C, f3 func(C) D) func(A) D {
	return func(v A) D { return f3(f2(f1(v))) }
}

// Compose4 composes four functions: A => B => C => D => E into A => E.
func Compose4[A, B, C, D, E any](f1 func(A) B, f2 func(B) C, f3 func(C) D, f4 func(D) E) func(A) E {
	return func(v A) E { return f4(f3(f2(f1(v)))) }
}

// Tap returns a function that calls fn on the value as a side effect, then returns it unchanged.
// Useful for logging or metrics inside a pipeline without changing the type.
//
// Example:
//
//	pipe.Pipe3(
//	    input,
//	    pipe.Tap(func(s string) { log.Println("after trim:", s) }),
//	    strings.ToUpper,
//	)
func Tap[T any](fn func(T)) func(T) T {
	return func(v T) T {
		fn(v)
		return v
	}
}

// Identity returns its argument unchanged. Useful as a no-op transform in
// pipelines, e.g. seq.OfSlice(items).Map(pipe.Identity).
func Identity[T any](v T) T { return v }

// -- numeric type constraints -------------------------------------------------

// SignedInteger covers all signed integer types (~int8 through ~int64 and ~int).
type SignedInteger interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// UnsignedInteger covers all unsigned integer types.
type UnsignedInteger interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer covers all integer types (signed and unsigned).
// Use for operations valid on any integer, e.g. IsEven/IsOdd.
type Integer interface {
	SignedInteger | UnsignedInteger
}

// Float covers all floating-point types (~float32 and ~float64).
type Float interface {
	~float32 | ~float64
}

// Signed covers all types that can represent negative values
// (signed integers and floats). Use for IsNegative.
type Signed interface {
	SignedInteger | Float
}

// Number covers all numeric primitive types (signed, unsigned, float).
// Use for IsPositive.
type Number interface {
	Integer | Float
}

// -- common predicates --------------------------------------------------------

// IsTrue reports whether v is true.
func IsTrue(v bool) bool { return v }

// IsFalse reports whether v is false.
func IsFalse(v bool) bool { return !v }

// IsZero reports whether v is the zero value for its type.
func IsZero[T comparable](v T) bool {
	var zero T
	return v == zero
}

// IsNonZero reports whether v is not the zero value for its type.
func IsNonZero[T comparable](v T) bool {
	var zero T
	return v != zero
}

// IsEven reports whether n is even. Defined for all integer types.
func IsEven[T Integer](n T) bool { return n%2 == 0 }

// IsOdd reports whether n is odd. Defined for all integer types.
func IsOdd[T Integer](n T) bool { return n%2 != 0 }

// IsPositive reports whether n > 0. Defined for all numeric types (integer and float).
func IsPositive[T Number](n T) bool { return n > 0 }

// IsNegative reports whether n < 0. Defined for signed types only (signed integers
// and floats); intentionally excludes unsigned integers since they can never be negative.
func IsNegative[T Signed](n T) bool { return n < 0 }

// Not wraps a predicate, returning its logical negation.
func Not[T any](pred func(T) bool) func(T) bool {
	return func(v T) bool { return !pred(v) }
}
