package pseq

import (
	"github.com/natalie-o-perret/go-functional-ish/option"
	"github.com/natalie-o-perret/go-functional-ish/seq"
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
