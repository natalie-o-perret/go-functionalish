// Package enum provides a lazy, generic Enumerable[T] sequence type inspired
// by C# LINQ and F# sequences. Pipelines are built by chaining methods; nothing
// runs until a terminal method (ToSlice, First, Count, …) is called.
//
// Design note — methods vs package-level functions:
// Go does not allow methods to introduce new type parameters. Operations that
// transform the element type (Map, FlatMap, GroupBy, …) are therefore
// package-level functions. Same-type operations are methods.
package enum

import (
	"cmp"
	"iter"
	"slices"

	"github.com/BooleanCat/go-functional/v2/it"
	"github.com/BooleanCat/go-functional/v2/it/itx"
	"github.com/natalie/gof/option"
	"github.com/natalie/gof/result"
)

// Enumerable[T] is a lazy sequence backed by iter.Seq[T].
type Enumerable[T any] struct {
	inner itx.Iterator[T]
}

// ── constructors ──────────────────────────────────────────────────────────────

// From wraps a slice into an Enumerable.
func From[T any](items []T) Enumerable[T] {
	return Enumerable[T]{inner: itx.FromSlice(items)}
}

// FromSeq wraps an iter.Seq[T] into an Enumerable.
func FromSeq[T any](seq iter.Seq[T]) Enumerable[T] {
	return Enumerable[T]{inner: itx.From[T](seq)}
}

// Range produces an Enumerable of integers from start (inclusive) to end (exclusive).
func Range(start, end int) Enumerable[int] {
	return FromSeq(func(yield func(int) bool) {
		for i := start; i < end; i++ {
			if !yield(i) {
				return
			}
		}
	})
}

// Repeat produces an Enumerable that yields v exactly n times.
func Repeat[T any](v T, n int) Enumerable[T] {
	return FromSeq(func(yield func(T) bool) {
		for range n {
			if !yield(v) {
				return
			}
		}
	})
}

// Empty returns an Enumerable with no elements.
func Empty[T any]() Enumerable[T] {
	return FromSeq(func(func(T) bool) {})
}

// ── lazy pipeline methods ─────────────────────────────────────────────────────

// Filter yields elements for which fn returns true.
func (e Enumerable[T]) Filter(fn func(T) bool) Enumerable[T] {
	return Enumerable[T]{inner: e.inner.Filter(fn)}
}

// Exclude yields elements for which fn returns false (inverse of Filter).
func (e Enumerable[T]) Exclude(fn func(T) bool) Enumerable[T] {
	return Enumerable[T]{inner: e.inner.Exclude(fn)}
}

// Take yields at most the first n elements.
func (e Enumerable[T]) Take(n uint) Enumerable[T] {
	return Enumerable[T]{inner: e.inner.Take(n)}
}

// TakeWhile yields elements as long as fn returns true, then stops.
func (e Enumerable[T]) TakeWhile(fn func(T) bool) Enumerable[T] {
	return Enumerable[T]{inner: e.inner.TakeWhile(fn)}
}

// Drop skips the first n elements.
func (e Enumerable[T]) Drop(n uint) Enumerable[T] {
	return Enumerable[T]{inner: e.inner.Drop(n)}
}

// DropWhile skips elements as long as fn returns true, then yields the rest.
func (e Enumerable[T]) DropWhile(fn func(T) bool) Enumerable[T] {
	return Enumerable[T]{inner: e.inner.DropWhile(fn)}
}

// Chain appends other Enumerables to the end of this one.
func (e Enumerable[T]) Chain(others ...Enumerable[T]) Enumerable[T] {
	iters := make([]func(func(T) bool), len(others)+1)
	iters[0] = e.inner
	for i, o := range others {
		iters[i+1] = o.inner
	}
	return Enumerable[T]{inner: itx.From[T](it.Chain(iters...))}
}

// Cycle repeats the sequence indefinitely. Always pair with Take.
func (e Enumerable[T]) Cycle() Enumerable[T] {
	return Enumerable[T]{inner: e.inner.Cycle()}
}

// SortBy materialises, sorts using the comparator, then re-wraps lazily.
func (e Enumerable[T]) SortBy(less func(a, b T) int) Enumerable[T] {
	cp := e.inner.Collect()
	slices.SortFunc(cp, less)
	return From(cp)
}

// SortAscBy materialises and sorts ascending by an ordered key.
func SortAscBy[T any, K cmp.Ordered](e Enumerable[T], key func(T) K) Enumerable[T] {
	return e.SortBy(func(a, b T) int { return cmp.Compare(key(a), key(b)) })
}

// SortDescBy materialises and sorts descending by an ordered key.
func SortDescBy[T any, K cmp.Ordered](e Enumerable[T], key func(T) K) Enumerable[T] {
	return e.SortBy(func(a, b T) int { return cmp.Compare(key(b), key(a)) })
}

// Reverse materialises, reverses, then re-wraps lazily.
func (e Enumerable[T]) Reverse() Enumerable[T] {
	cp := e.inner.Collect()
	slices.Reverse(cp)
	return From(cp)
}

// ForEach calls fn for every element, consuming the sequence.
func (e Enumerable[T]) ForEach(fn func(T)) {
	e.inner.ForEach(fn)
}

// ── terminal methods ──────────────────────────────────────────────────────────

// ToSlice materialises the sequence into a slice.
func (e Enumerable[T]) ToSlice() []T { return e.inner.Collect() }

// Count returns the total number of elements.
func (e Enumerable[T]) Count() int { return e.inner.Len() }

// CountBy returns the number of elements satisfying fn.
func (e Enumerable[T]) CountBy(fn func(T) bool) int { return e.inner.Filter(fn).Len() }

// Any returns true if at least one element satisfies fn (short-circuits).
func (e Enumerable[T]) Any(fn func(T) bool) bool {
	_, ok := e.inner.Find(fn)
	return ok
}

// Every returns true if all elements satisfy fn (short-circuits on first failure).
func (e Enumerable[T]) Every(fn func(T) bool) bool {
	_, ok := e.inner.Find(func(v T) bool { return !fn(v) })
	return !ok
}

// First returns the first element and true, or the zero value and false.
func (e Enumerable[T]) First() (T, bool) {
	return e.inner.Find(func(T) bool { return true })
}

// FirstOption returns the first element as an Option.
func (e Enumerable[T]) FirstOption() option.Option[T] {
	v, ok := e.First()
	if ok {
		return option.Some(v)
	}
	return option.None[T]()
}

// Last returns the last element and true, or the zero value and false.
func (e Enumerable[T]) Last() (T, bool) {
	var last T
	found := false
	for v := range iter.Seq[T](e.inner) {
		last, found = v, true
	}
	return last, found
}

// LastOption returns the last element as an Option.
func (e Enumerable[T]) LastOption() option.Option[T] {
	v, ok := e.Last()
	if ok {
		return option.Some(v)
	}
	return option.None[T]()
}

// Reduce folds the sequence using fn, starting with initial.
func Reduce[T, A any](e Enumerable[T], initial A, fn func(A, T) A) A {
	acc := initial
	for v := range iter.Seq[T](e.inner) {
		acc = fn(acc, v)
	}
	return acc
}

// ── type-transforming package-level functions ─────────────────────────────────

// Map transforms Enumerable[T] → Enumerable[R] lazily.
func Map[T, R any](e Enumerable[T], fn func(T) R) Enumerable[R] {
	return Enumerable[R]{inner: itx.From[R](it.Map(iter.Seq[T](e.inner), fn))}
}

// FlatMap transforms Enumerable[T] → Enumerable[R] via a one-to-many mapping.
func FlatMap[T, R any](e Enumerable[T], fn func(T) []R) Enumerable[R] {
	return Enumerable[R]{inner: func(yield func(R) bool) {
		for v := range iter.Seq[T](e.inner) {
			for _, r := range fn(v) {
				if !yield(r) {
					return
				}
			}
		}
	}}
}

// GroupBy materialises and groups elements by key.
func GroupBy[T any, K comparable](e Enumerable[T], key func(T) K) map[K][]T {
	out := make(map[K][]T)
	for v := range iter.Seq[T](e.inner) {
		k := key(v)
		out[k] = append(out[k], v)
	}
	return out
}

// Uniq lazily removes duplicate comparable elements, preserving first-seen order.
func Uniq[T comparable](e Enumerable[T]) Enumerable[T] {
	return Enumerable[T]{inner: itx.From[T](it.FilterUnique(e.inner))}
}

// UniqBy lazily removes duplicates by key, preserving first-seen order.
func UniqBy[T any, K comparable](e Enumerable[T], key func(T) K) Enumerable[T] {
	return Enumerable[T]{inner: func(yield func(T) bool) {
		seen := make(map[K]struct{})
		for v := range iter.Seq[T](e.inner) {
			k := key(v)
			if _, ok := seen[k]; !ok {
				seen[k] = struct{}{}
				if !yield(v) {
					return
				}
			}
		}
	}}
}

// Pair holds two values produced by Zip.
type Pair[T, U any] struct {
	First  T
	Second U
}

// Zip lazily pairs elements from two Enumerables. Stops at the shorter one.
func Zip[T, U any](a Enumerable[T], b Enumerable[U]) Enumerable[Pair[T, U]] {
	return Enumerable[Pair[T, U]]{inner: func(yield func(Pair[T, U]) bool) {
		nextB, stopB := iter.Pull(iter.Seq[U](b.inner))
		defer stopB()
		for v := range iter.Seq[T](a.inner) {
			u, ok := nextB()
			if !ok {
				return
			}
			if !yield(Pair[T, U]{First: v, Second: u}) {
				return
			}
		}
	}}
}

// FilterMap applies fn to each element, keeping Ok values and discarding Err.
func FilterMap[T, R, E any](e Enumerable[T], fn func(T) result.Result[R, E]) Enumerable[R] {
	return Enumerable[R]{inner: func(yield func(R) bool) {
		for v := range iter.Seq[T](e.inner) {
			if r := fn(v); r.IsOk() {
				if !yield(r.Unwrap()) {
					return
				}
			}
		}
	}}
}

