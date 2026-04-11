// Package validation provides a type-safe Validation type for applicative error accumulation.
//
// Unlike [github.com/natalie-o-perret/gof/result] which short-circuits on the first error,
// Validation runs all checks and collects every failure. This makes it ideal for
// form validation, config parsing, and anywhere you want to report all problems at once.
//
//	v := validation.Map3(
//	    validateName(input.Name),
//	    validateAge(input.Age),
//	    validateEmail(input.Email),
//	    NewUser,
//	)
//	// If all pass: Success(user)
//	// If any fail: Failure with ALL errors combined
package validation

import (
	"github.com/natalie-o-perret/gof/option"
	"github.com/natalie-o-perret/gof/result"
)

// Validation[T, E] holds either a successful value of type T
// or a slice of errors of type E. Unlike Result, multiple errors
// are accumulated when combining validations via Map2-Map5, Apply,
// Sequence, or Traverse.
type Validation[T, E any] struct {
	value  T
	errors []E
	ok     bool
}

// -- constructors -------------------------------------------------------------

// Success wraps a value in a successful Validation.
func Success[T, E any](v T) Validation[T, E] {
	return Validation[T, E]{value: v, ok: true}
}

// Failure creates a failed Validation with one or more errors.
// Panics if no errors are provided.
func Failure[T, E any](errs ...E) Validation[T, E] {
	if len(errs) == 0 {
		panic("validation: Failure requires at least one error")
	}
	return Validation[T, E]{errors: errs, ok: false}
}

// FromResult converts a Result into a Validation.
// Ok becomes Success; Err becomes Failure with a single error.
func FromResult[T, E any](r result.Result[T, E]) Validation[T, E] {
	if r.IsOk() {
		return Success[T, E](r.Unwrap())
	}
	return Failure[T, E](r.UnwrapErr())
}

// FromOption converts an Option into a Validation.
// Some becomes Success; None becomes Failure with the provided error.
func FromOption[T, E any](o option.Option[T], errIfNone E) Validation[T, E] {
	if o.IsSome() {
		return Success[T, E](o.Unwrap())
	}
	return Failure[T, E](errIfNone)
}

// -- methods ------------------------------------------------------------------

// IsSuccess reports whether the Validation is successful.
func (v Validation[T, E]) IsSuccess() bool { return v.ok }

// IsFailure reports whether the Validation has errors.
func (v Validation[T, E]) IsFailure() bool { return !v.ok }

// Unwrap returns the success value. Panics if Failure.
func (v Validation[T, E]) Unwrap() T {
	if !v.ok {
		panic("validation: Unwrap called on Failure")
	}
	return v.value
}

// UnwrapOr returns the success value, or def if Failure.
func (v Validation[T, E]) UnwrapOr(def T) T {
	if v.ok {
		return v.value
	}
	return def
}

// UnwrapOrElse returns the success value, or calls fn with the errors if Failure.
func (v Validation[T, E]) UnwrapOrElse(fn func([]E) T) T {
	if v.ok {
		return v.value
	}
	return fn(v.errors)
}

// UnwrapErrors returns the error slice. Panics if Success.
func (v Validation[T, E]) UnwrapErrors() []E {
	if v.ok {
		panic("validation: UnwrapErrors called on Success")
	}
	return v.errors
}

// ToResult converts Success to Ok and Failure to Err (with the error slice).
func (v Validation[T, E]) ToResult() result.Result[T, []E] {
	if v.ok {
		return result.Ok[T, []E](v.value)
	}
	return result.Err[T, []E](v.errors)
}

// ToOption converts Success to Some and Failure to None.
func (v Validation[T, E]) ToOption() option.Option[T] {
	if v.ok {
		return option.Some(v.value)
	}
	return option.None[T]()
}

// -- type-transforming functions (package-level) ------------------------------

// Map applies fn to the success value, passing Failure unchanged.
func Map[T, R, E any](v Validation[T, E], fn func(T) R) Validation[R, E] {
	if v.ok {
		return Success[R, E](fn(v.value))
	}
	return Validation[R, E]{errors: v.errors}
}

// MapError applies fn to each error, passing Success unchanged.
func MapError[T, E, F any](v Validation[T, E], fn func(E) F) Validation[T, F] {
	if v.ok {
		return Success[T, F](v.value)
	}
	mapped := make([]F, len(v.errors))
	for i, e := range v.errors {
		mapped[i] = fn(e)
	}
	return Validation[T, F]{errors: mapped}
}

// FlatMap applies fn to the success value, flattening the result.
// Note: this short-circuits like Result. For error accumulation, use Map2-Map5 or Apply.
func FlatMap[T, R, E any](v Validation[T, E], fn func(T) Validation[R, E]) Validation[R, E] {
	if v.ok {
		return fn(v.value)
	}
	return Validation[R, E]{errors: v.errors}
}

// -- applicative combinators (error-accumulating) -----------------------------

// Apply applies a wrapped function to a wrapped value, accumulating errors from both.
// If both are Failure, all errors are merged. If both are Success, fn(val) is returned.
func Apply[T, R, E any](vFn Validation[func(T) R, E], vVal Validation[T, E]) Validation[R, E] {
	switch {
	case vFn.ok && vVal.ok:
		return Success[R, E](vFn.value(vVal.value))
	case !vFn.ok && !vVal.ok:
		return Validation[R, E]{errors: append(vFn.errors, vVal.errors...)}
	case !vFn.ok:
		return Validation[R, E]{errors: vFn.errors}
	default:
		return Validation[R, E]{errors: vVal.errors}
	}
}

// Map2 combines two Validations with fn, accumulating errors from both.
func Map2[A, B, R, E any](
	va Validation[A, E],
	vb Validation[B, E],
	fn func(A, B) R,
) Validation[R, E] {
	if va.ok && vb.ok {
		return Success[R, E](fn(va.value, vb.value))
	}
	return Validation[R, E]{errors: mergeErrors(va.errors, vb.errors)}
}

// Map3 combines three Validations with fn, accumulating all errors.
func Map3[A, B, C, R, E any](
	va Validation[A, E],
	vb Validation[B, E],
	vc Validation[C, E],
	fn func(A, B, C) R,
) Validation[R, E] {
	if va.ok && vb.ok && vc.ok {
		return Success[R, E](fn(va.value, vb.value, vc.value))
	}
	return Validation[R, E]{errors: mergeErrors(va.errors, vb.errors, vc.errors)}
}

// Map4 combines four Validations with fn, accumulating all errors.
func Map4[A, B, C, D, R, E any](
	va Validation[A, E],
	vb Validation[B, E],
	vc Validation[C, E],
	vd Validation[D, E],
	fn func(A, B, C, D) R,
) Validation[R, E] {
	if va.ok && vb.ok && vc.ok && vd.ok {
		return Success[R, E](fn(va.value, vb.value, vc.value, vd.value))
	}
	return Validation[R, E]{errors: mergeErrors(va.errors, vb.errors, vc.errors, vd.errors)}
}

// Map5 combines five Validations with fn, accumulating all errors.
func Map5[A, B, C, D, F, R, E any](
	va Validation[A, E],
	vb Validation[B, E],
	vc Validation[C, E],
	vd Validation[D, E],
	vf Validation[F, E],
	fn func(A, B, C, D, F) R,
) Validation[R, E] {
	if va.ok && vb.ok && vc.ok && vd.ok && vf.ok {
		return Success[R, E](fn(va.value, vb.value, vc.value, vd.value, vf.value))
	}
	return Validation[R, E]{errors: mergeErrors(va.errors, vb.errors, vc.errors, vd.errors, vf.errors)}
}

// -- sequence / traverse ------------------------------------------------------

// Sequence converts a slice of Validations into a Validation of a slice.
// If all succeed, returns Success with all values. If any fail, returns
// Failure with ALL errors from ALL failures accumulated.
func Sequence[T, E any](vs []Validation[T, E]) Validation[[]T, E] {
	values := make([]T, 0, len(vs))
	var errs []E
	for _, v := range vs {
		if v.ok {
			values = append(values, v.value)
		} else {
			errs = append(errs, v.errors...)
		}
	}
	if errs != nil {
		return Validation[[]T, E]{errors: errs}
	}
	return Success[[]T, E](values)
}

// Traverse applies fn to each item and sequences the results.
// If all succeed, returns Success with all mapped values. If any fail,
// returns Failure with ALL errors accumulated.
func Traverse[T, R, E any](items []T, fn func(T) Validation[R, E]) Validation[[]R, E] {
	values := make([]R, 0, len(items))
	var errs []E
	for _, item := range items {
		v := fn(item)
		if v.ok {
			values = append(values, v.value)
		} else {
			errs = append(errs, v.errors...)
		}
	}
	if errs != nil {
		return Validation[[]R, E]{errors: errs}
	}
	return Success[[]R, E](values)
}

// -- internal helpers ---------------------------------------------------------

func mergeErrors[E any](slices ...[]E) []E {
	var out []E
	for _, s := range slices {
		out = append(out, s...)
	}
	return out
}
