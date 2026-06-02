// Package seq provides a lazy, generic Seq[T] sequence type inspired
// by F# sequences. Pipelines are built by chaining methods; nothing
// runs until a terminal method (ToSlice, Head, Length, ...) is called.
//
// Generic methods (Go 1.27+) allow type-changing operations such as Map,
// Collect, Choose, Fold, GroupBy, and SortBy to be called directly on
// a Seq value, enabling fully fluent left-to-right pipelines:
// (Zip remains a package-level function due to a type-checker instantiation cycle.)
//
//	seq.OfSlice(items).
//	    Filter(isActive).
//	    Map(toName).
//	    SortBy(strings.ToLower).
//	    ToSlice()
package seq

import (
	"cmp"
	"iter"
	"slices"

	"github.com/natalie-o-perret/go-functionalish/tuple"
)

// Seq is a lazy sequence (named type over iter.Seq[T]) so that
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
func (s Seq[T]) Truncate(n int) Seq[T] {
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
func (s Seq[T]) Skip(n int) Seq[T] {
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

// -- type-transforming methods (enabled by Go 1.27 generic methods) ------------

// Map transforms each element T => R lazily.
func (s Seq[T]) Map[R any](fn func(T) R) Seq[R] {
	return func(yield func(R) bool) {
		for v := range s {
			if !yield(fn(v)) {
				return
			}
		}
	}
}

// Collect transforms T => []R lazily, flattening the results into a single Seq[R].
func (s Seq[T]) Collect[R any](fn func(T) []R) Seq[R] {
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

// Choose applies fn to each element, keeping Some values and discarding None.
func (s Seq[T]) Choose[R any](fn func(T) option.Option[R]) Seq[R] {
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

// Fold folds the sequence left using fn, starting with initial.
func (s Seq[T]) Fold[A any](initial A, fn func(A, T) A) A {
	acc := initial
	for v := range s {
		acc = fn(acc, v)
	}
	return acc
}

// GroupBy materialises the sequence and groups elements by key.
func (s Seq[T]) GroupBy[K comparable](key func(T) K) map[K][]T {
	out := make(map[K][]T)
	for v := range s {
		k := key(v)
		out[k] = append(out[k], v)
	}
	return out
}

// SortBy materialises and sorts ascending by an ordered key.
func (s Seq[T]) SortBy[K cmp.Ordered](key func(T) K) Seq[T] {
	return s.SortWith(func(a, b T) int { return cmp.Compare(key(a), key(b)) })
}

// SortByDescending materialises and sorts descending by an ordered key.
func (s Seq[T]) SortByDescending[K cmp.Ordered](key func(T) K) Seq[T] {
	return s.SortWith(func(a, b T) int { return cmp.Compare(key(b), key(a)) })
}

// DistinctBy lazily removes duplicates by key, preserving first-seen order.
func (s Seq[T]) DistinctBy[K comparable](key func(T) K) Seq[T] {
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

// -- lazy pipeline methods (continued) ----------------------------------------

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

// Pair is a type alias for tuple.Pair, kept for backward compatibility.
// Prefer tuple.Pair in new code.
type Pair[T, U any] = tuple.Pair[T, U]

// PairOf is a type alias for tuple.PairOf, kept for backward compatibility.
// Prefer tuple.PairOf in new code.
func PairOf[T, U any](first T, second U) Pair[T, U] { return tuple.PairOf(first, second) }

// ZipWith lazily combines s and other element-wise using fn. Stops at the shorter Seq.
// See also the package-level [Zip] function which lazily pairs elements from two Seqs.
func (s Seq[T]) ZipWith[U, R any](other Seq[U], fn func(T, U) R) Seq[R] {
	return func(yield func(R) bool) {
		nextU, stopU := iter.Pull(iter.Seq[U](other))
		defer stopU()
		for v := range s {
			u, ok := nextU()
			if !ok {
				return
			}
			if !yield(fn(v, u)) {
				return
			}
		}
	}
}

// ZipWith3 lazily combines s, b, and c element-wise using fn. Stops at the shortest Seq.
func (s Seq[T]) ZipWith3[U, V, R any](b Seq[U], c Seq[V], fn func(T, U, V) R) Seq[R] {
	return func(yield func(R) bool) {
		nextU, stopU := iter.Pull(iter.Seq[U](b))
		defer stopU()
		nextV, stopV := iter.Pull(iter.Seq[V](c))
		defer stopV()
		for v := range s {
			u, okU := nextU()
			if !okU {
				return
			}
			w, okV := nextV()
			if !okV {
				return
			}
			if !yield(fn(v, u, w)) {
				return
			}
		}
	}
}

// Zip cannot be a method because returning Seq[Pair[T,U]] would create
// an instantiation cycle in the type checker (Go spec §Generic methods).
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
