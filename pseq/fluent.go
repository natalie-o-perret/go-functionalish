package pseq

import (
	"github.com/natalie-o-perret/go-functionalish/option"
	"github.com/natalie-o-perret/go-functionalish/seq"
)

// ---------------------------------------------------------------------------
// Curried transform helpers  (for use with pipe / seq.Then)
// ---------------------------------------------------------------------------

// MapFn returns a transform: seq.Seq[T] → seq.Seq[R], applied in parallel.
func MapFn[T, R any](fn func(T) R, opts ...Option) func(seq.Seq[T]) seq.Seq[R] {
	return func(s seq.Seq[T]) seq.Seq[R] { return Map(s, fn, opts...) }
}

// FilterFn returns a transform: seq.Seq[T] → seq.Seq[T], applied in parallel.
func FilterFn[T any](fn func(T) bool, opts ...Option) func(seq.Seq[T]) seq.Seq[T] {
	return func(s seq.Seq[T]) seq.Seq[T] { return Filter(s, fn, opts...) }
}

// CollectFn returns a transform: seq.Seq[T] → seq.Seq[R] (flatMap), applied in parallel.
func CollectFn[T, R any](fn func(T) []R, opts ...Option) func(seq.Seq[T]) seq.Seq[R] {
	return func(s seq.Seq[T]) seq.Seq[R] { return Collect(s, fn, opts...) }
}

// ChooseFn returns a transform: seq.Seq[T] → seq.Seq[R], applied in parallel.
func ChooseFn[T, R any](fn func(T) option.Option[R], opts ...Option) func(seq.Seq[T]) seq.Seq[R] {
	return func(s seq.Seq[T]) seq.Seq[R] { return Choose(s, fn, opts...) }
}

// ---------------------------------------------------------------------------
// Curried terminal helpers
// ---------------------------------------------------------------------------

// ForEachFn returns a terminal that calls fn on every element in parallel.
func ForEachFn[T any](fn func(T), opts ...Option) func(seq.Seq[T]) {
	return func(s seq.Seq[T]) { ForEach(s, fn, opts...) }
}

// GroupByFn returns a terminal: seq.Seq[T] → map[K][]T, applied in parallel.
func GroupByFn[T any, K comparable](key func(T) K, opts ...Option) func(seq.Seq[T]) map[K][]T {
	return func(s seq.Seq[T]) map[K][]T { return GroupBy(s, key, opts...) }
}

// CountByKeyFn returns a terminal: seq.Seq[T] → map[K]int, applied in parallel.
func CountByKeyFn[T any, K comparable](fn func(T) K, opts ...Option) func(seq.Seq[T]) map[K]int {
	return func(s seq.Seq[T]) map[K]int { return CountByKey(s, fn, opts...) }
}

// ---------------------------------------------------------------------------
// PSeq[T] — fluent parallel sequence builder
// ---------------------------------------------------------------------------

// PSeq[T] is a fluent parallel sequence builder.
//
// Construct one with [Of] or [OfSlice], then chain operations left-to-right:
//
//	result := pseq.OfSlice(bigData, pseq.WithWorkers(4)).
//	    Filter(isValid).
//	    Map(transform).
//	    ToSlice()
//
// Workers configured on the PSeq propagate automatically to every chained
// step. Individual steps may override workers via their own opts arguments.
type PSeq[T any] struct {
	s    seq.Seq[T]
	opts []Option
}

// Of wraps a [seq.Seq] for fluent parallel processing.
func Of[T any](s seq.Seq[T], opts ...Option) PSeq[T] {
	return PSeq[T]{s: s, opts: opts}
}

// OfSlice wraps a slice for fluent parallel processing.
func OfSlice[T any](items []T, opts ...Option) PSeq[T] {
	return Of(seq.OfSlice(items), opts...)
}

// mergeOpts appends per-call overrides on top of the PSeq's default opts.
// Per-call opts are applied last and therefore take precedence.
func mergeOpts(base, extra []Option) []Option {
	if len(extra) == 0 {
		return base
	}
	merged := make([]Option, len(base)+len(extra))
	copy(merged, base)
	copy(merged[len(base):], extra)
	return merged
}

// Seq returns the underlying [seq.Seq] so the result can be handed back to
// sequential seq operations.
func (p PSeq[T]) Seq() seq.Seq[T] { return p.s }

// ToSlice materialises the sequence into a slice.
func (p PSeq[T]) ToSlice() []T { return p.s.ToSlice() }

// ---------------------------------------------------------------------------
// Same-type transforms
// ---------------------------------------------------------------------------

// Filter keeps elements where fn returns true, processing chunks in parallel.
// Output order matches input order.
func (p PSeq[T]) Filter(fn func(T) bool, opts ...Option) PSeq[T] {
	return PSeq[T]{s: Filter(p.s, fn, mergeOpts(p.opts, opts)...), opts: p.opts}
}

// ---------------------------------------------------------------------------
// Type-changing transforms (Go 1.27 generic methods)
// ---------------------------------------------------------------------------

// Map transforms each element T => R in parallel, preserving order.
func (p PSeq[T]) Map[R any](fn func(T) R, opts ...Option) PSeq[R] {
	return PSeq[R]{s: Map(p.s, fn, mergeOpts(p.opts, opts)...), opts: p.opts}
}

// Collect (flatMap) maps T => []R in parallel, flattening results in order.
func (p PSeq[T]) Collect[R any](fn func(T) []R, opts ...Option) PSeq[R] {
	return PSeq[R]{s: Collect(p.s, fn, mergeOpts(p.opts, opts)...), opts: p.opts}
}

// Choose applies fn to each element in parallel, keeping [option.Some] values
// and discarding [option.None]. Output order is preserved.
func (p PSeq[T]) Choose[R any](fn func(T) option.Option[R], opts ...Option) PSeq[R] {
	return PSeq[R]{s: Choose(p.s, fn, mergeOpts(p.opts, opts)...), opts: p.opts}
}

// GroupBy groups elements by key, processing chunks in parallel.
// Order within each group is preserved.
func (p PSeq[T]) GroupBy[K comparable](key func(T) K, opts ...Option) map[K][]T {
	return GroupBy(p.s, key, mergeOpts(p.opts, opts)...)
}

// CountByKey counts occurrences of each key in parallel.
func (p PSeq[T]) CountByKey[K comparable](fn func(T) K, opts ...Option) map[K]int {
	return CountByKey(p.s, fn, mergeOpts(p.opts, opts)...)
}

// SumBy returns the sum of fn(element) for all elements, computed in parallel.
func (p PSeq[T]) SumBy[N seq.Numeric](fn func(T) N, opts ...Option) N {
	return SumBy(p.s, fn, mergeOpts(p.opts, opts)...)
}

// ---------------------------------------------------------------------------
// Terminal operations
// ---------------------------------------------------------------------------

// ForEach calls fn for every element in parallel. Invocation order is not
// guaranteed; use [PSeq.Map] if ordered results are needed.
func (p PSeq[T]) ForEach(fn func(T), opts ...Option) {
	ForEach(p.s, fn, mergeOpts(p.opts, opts)...)
}

// Reduce combines elements using an associative fn in parallel.
// Returns (zero, false) for an empty sequence.
func (p PSeq[T]) Reduce(fn func(T, T) T, opts ...Option) (T, bool) {
	return Reduce(p.s, fn, mergeOpts(p.opts, opts)...)
}

// Partition splits elements into two slices in parallel: the first where fn
// returns true, the second where fn returns false. Order is preserved in each.
func (p PSeq[T]) Partition(fn func(T) bool, opts ...Option) ([]T, []T) {
	return Partition(p.s, fn, mergeOpts(p.opts, opts)...)
}

// Exists returns true if at least one element satisfies fn.
// Goroutines short-circuit as soon as a match is found.
func (p PSeq[T]) Exists(fn func(T) bool, opts ...Option) bool {
	return Exists(p.s, fn, mergeOpts(p.opts, opts)...)
}

// ForAll returns true if all elements satisfy fn.
// Goroutines short-circuit as soon as a non-matching element is found.
func (p PSeq[T]) ForAll(fn func(T) bool, opts ...Option) bool {
	return ForAll(p.s, fn, mergeOpts(p.opts, opts)...)
}

// ---------------------------------------------------------------------------
// Package-level helpers for numeric PSeq
// (Sum cannot be a method: the receiver constraint T any cannot be narrowed
//  to seq.Numeric on a per-method basis)
// ---------------------------------------------------------------------------

// SumOf returns the sum of all elements in a numeric PSeq, computed in parallel.
func SumOf[T seq.Numeric](p PSeq[T], opts ...Option) T {
	return Sum(p.s, mergeOpts(p.opts, opts)...)
}
