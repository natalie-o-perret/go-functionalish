// Package seq provides a lazy, generic Seq[T] sequence type inspired
// by F# sequences. Pipelines are built by chaining methods; nothing
// runs until a terminal method (ToSlice, Head, Length, ...) is called.
//
// Design note : methods vs package-level functions:
// Go does not allow methods to introduce new type parameters. Operations that
// transform the element type (Map, Collect, GroupBy, ...) are therefore
// package-level functions. Same-type operations are methods.
package seq

import (
	"cmp"
	"iter"
	"slices"

	"github.com/natalie-o-perret/go-functionalish/option"
)

// Seq[T] is a lazy sequence. It is a named type over iter.Seq[T] so that
// methods can be attached directly without a wrapper struct.
type Seq[T any] iter.Seq[T]

// -- constructors --------------------------------------------------------------

// OfSlice wraps a slice into a Seq.
func OfSlice[T any](items []T) Seq[T] {
	return Seq[T](slices.Values(items))
}

// FromIter wraps an iter.Seq[T] into a Seq.
func FromIter[T any](seq iter.Seq[T]) Seq[T] {
	return Seq[T](seq)
}

// Range produces a Seq of integers from start (inclusive) to end (exclusive).
func Range(start, end int) Seq[int] {
	return func(yield func(int) bool) {
		for i := start; i < end; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

// Replicate produces a Seq that yields v exactly n times.
func Replicate[T any](v T, n int) Seq[T] {
	return func(yield func(T) bool) {
		for range n {
			if !yield(v) {
				return
			}
		}
	}
}

// Empty returns a Seq with no elements.
func Empty[T any]() Seq[T] {
	return func(func(T) bool) {}
}

// -- lazy pipeline methods -----------------------------------------------------

// Filter yields elements for which fn returns true.
func (s Seq[T]) Filter(fn func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if fn(v) && !yield(v) {
				return
			}
		}
	}
}

// Exclude yields elements for which fn returns false (inverse of Filter).
func (s Seq[T]) Exclude(fn func(T) bool) Seq[T] {
	return s.Filter(func(v T) bool { return !fn(v) })
}

// Truncate yields at most the first n elements.
func (s Seq[T]) Truncate(n uint) Seq[T] {
	return func(yield func(T) bool) {
		remaining := n
		for v := range s {
			if remaining == 0 {
				return
			}
			remaining--
			if !yield(v) {
				return
			}
		}
	}
}

// TakeWhile yields elements as long as fn returns true, then stops.
func (s Seq[T]) TakeWhile(fn func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if !fn(v) || !yield(v) {
				return
			}
		}
	}
}

// Skip skips the first n elements.
func (s Seq[T]) Skip(n uint) Seq[T] {
	return func(yield func(T) bool) {
		remaining := n
		for v := range s {
			if remaining > 0 {
				remaining--
				continue
			}
			if !yield(v) {
				return
			}
		}
	}
}

// SkipWhile skips elements as long as fn returns true, then yields the rest.
func (s Seq[T]) SkipWhile(fn func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		dropping := true
		for v := range s {
			if dropping && fn(v) {
				continue
			}
			dropping = false
			if !yield(v) {
				return
			}
		}
	}
}

// Append appends other Seqs to the end of this one.
func (s Seq[T]) Append(others ...Seq[T]) Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if !yield(v) {
				return
			}
		}
		for _, o := range others {
			for v := range o {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Cycle repeats the sequence indefinitely. Always pair with Truncate.
func (s Seq[T]) Cycle() Seq[T] {
	return func(yield func(T) bool) {
		for {
			for v := range s {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// SortWith materialises, sorts using a comparison function, then re-wraps lazily.
func (s Seq[T]) SortWith(cmpFn func(a, b T) int) Seq[T] {
	cp := slices.Collect(iter.Seq[T](s))
	slices.SortFunc(cp, cmpFn)
	return OfSlice(cp)
}

// SortBy materialises and sorts ascending by an ordered key.
func SortBy[T any, K cmp.Ordered](s Seq[T], key func(T) K) Seq[T] {
	return s.SortWith(func(a, b T) int { return cmp.Compare(key(a), key(b)) })
}

// SortByDescending materialises and sorts descending by an ordered key.
func SortByDescending[T any, K cmp.Ordered](s Seq[T], key func(T) K) Seq[T] {
	return s.SortWith(func(a, b T) int { return cmp.Compare(key(b), key(a)) })
}

// Rev materialises, reverses, then re-wraps lazily.
func (s Seq[T]) Rev() Seq[T] {
	cp := slices.Collect(iter.Seq[T](s))
	slices.Reverse(cp)
	return OfSlice(cp)
}

// Iter calls fn for every element, consuming the sequence.
func (s Seq[T]) Iter(fn func(T)) {
	for v := range s {
		fn(v)
	}
}

// -- terminal methods ----------------------------------------------------------

// ToSlice materialises the sequence into a slice.
//
// Deprecated: ToArray was renamed to ToSlice for Go idiom accuracy.
func (s Seq[T]) ToArray() []T { return s.ToSlice() }

// ToSlice materialises the sequence into a slice.
func (s Seq[T]) ToSlice() []T { return slices.Collect(iter.Seq[T](s)) }

// Length returns the total number of elements.
func (s Seq[T]) Length() int {
	n := 0
	for range s {
		n++
	}
	return n
}

// CountBy returns the number of elements satisfying fn.
func (s Seq[T]) CountBy(fn func(T) bool) int {
	return s.Filter(fn).Length()
}

// Exists returns true if at least one element satisfies fn. Short-circuits.
func (s Seq[T]) Exists(fn func(T) bool) bool {
	for v := range s {
		if fn(v) {
			return true
		}
	}
	return false
}

// ForAll returns true if all elements satisfy fn. Short-circuits on first failure.
func (s Seq[T]) ForAll(fn func(T) bool) bool {
	for v := range s {
		if !fn(v) {
			return false
		}
	}
	return true
}

// Head returns the first element and true, or the zero value and false.
func (s Seq[T]) Head() (T, bool) {
	for v := range s {
		return v, true
	}
	var zero T
	return zero, false
}

// TryHead returns the first element as an Option.
func (s Seq[T]) TryHead() option.Option[T] {
	v, ok := s.Head()
	if ok {
		return option.Some(v)
	}
	return option.None[T]()
}

// Last returns the last element and true, or the zero value and false.
func (s Seq[T]) Last() (T, bool) {
	var last T
	found := false
	for v := range s {
		last, found = v, true
	}
	return last, found
}

// TryLast returns the last element as an Option.
func (s Seq[T]) TryLast() option.Option[T] {
	v, ok := s.Last()
	if ok {
		return option.Some(v)
	}
	return option.None[T]()
}

// Fold folds the sequence using fn, starting with initial.
func Fold[T, A any](s Seq[T], initial A, fn func(A, T) A) A {
	acc := initial
	for v := range s {
		acc = fn(acc, v)
	}
	return acc
}

// -- type-transforming package-level functions ---------------------------------

// Map transforms Seq[T] => Seq[R] lazily.
func Map[T, R any](s Seq[T], fn func(T) R) Seq[R] {
	return func(yield func(R) bool) {
		for v := range s {
			if !yield(fn(v)) {
				return
			}
		}
	}
}

// Collect transforms Seq[T] => Seq[R] via a one-to-many mapping.
func Collect[T, R any](s Seq[T], fn func(T) []R) Seq[R] {
	return func(yield func(R) bool) {
		for v := range s {
			for _, r := range fn(v) {
				if !yield(r) {
					return
				}
			}
		}
	}
}

// GroupBy materialises and groups elements by key.
func GroupBy[T any, K comparable](s Seq[T], key func(T) K) map[K][]T {
	out := make(map[K][]T)
	for v := range s {
		k := key(v)
		out[k] = append(out[k], v)
	}
	return out
}

// Distinct lazily removes duplicate comparable elements, preserving first-seen order.
func Distinct[T comparable](s Seq[T]) Seq[T] {
	return func(yield func(T) bool) {
		seen := make(map[T]struct{})
		for v := range s {
			if _, ok := seen[v]; !ok {
				seen[v] = struct{}{}
				if !yield(v) {
					return
				}
			}
		}
	}
}

// DistinctBy lazily removes duplicates by key, preserving first-seen order.
func DistinctBy[T any, K comparable](s Seq[T], key func(T) K) Seq[T] {
	return func(yield func(T) bool) {
		seen := make(map[K]struct{})
		for v := range s {
			k := key(v)
			if _, ok := seen[k]; !ok {
				seen[k] = struct{}{}
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Pair holds two values produced by Zip.
type Pair[T, U any] struct {
	First  T
	Second U
}

// Zip lazily pairs elements from two Seqs. Stops at the shorter one.
func Zip[T, U any](a Seq[T], b Seq[U]) Seq[Pair[T, U]] {
	return func(yield func(Pair[T, U]) bool) {
		nextB, stopB := iter.Pull(iter.Seq[U](b))
		defer stopB()
		for v := range a {
			u, ok := nextB()
			if !ok {
				return
			}
			if !yield(Pair[T, U]{First: v, Second: u}) {
				return
			}
		}
	}
}

// Choose applies fn to each element, keeping Some values and discarding None.
func Choose[T, R any](s Seq[T], fn func(T) option.Option[R]) Seq[R] {
	return func(yield func(R) bool) {
		for v := range s {
			if r := fn(v); r.IsSome() {
				if !yield(r.Unwrap()) {
					return
				}
			}
		}
	}
}
