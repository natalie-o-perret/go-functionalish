package seq

import (
	"cmp"

	"github.com/natalie-o-perret/go-functionalish/option"
)

// Then threads a Seq[T] through a type-changing function.
//
// Deprecated: with Go 1.27 generic methods, type-changing operations like
// Map, Collect, Choose, and Fold are now directly chainable as methods:
//
//	// Before:
//	seq.Then(seq.OfSlice(cars).Filter(pred), seq.MapFn(toOwner)).ToSlice()
//
//	// After:
//	seq.OfSlice(cars).Filter(pred).Map(toOwner).ToSlice()
//
// Then is retained for cases where a function value of type func(Seq[T]) R
// must be passed to pipe.Pipe* or similar combinators.
func Then[T, R any](s Seq[T], fn func(Seq[T]) R) R {
	return fn(s)
}

// -- curried same-type helpers (for use with pipe) -----------------------------

// FilterFn returns a transform: Seq[T] => Seq[T], keeping elements where fn returns true.
func FilterFn[T any](fn func(T) bool) func(Seq[T]) Seq[T] {
	return func(s Seq[T]) Seq[T] { return s.Filter(fn) }
}

// ExcludeFn returns a transform: Seq[T] => Seq[T], removing elements where fn returns true.
func ExcludeFn[T any](fn func(T) bool) func(Seq[T]) Seq[T] {
	return func(s Seq[T]) Seq[T] { return s.Exclude(fn) }
}

// TruncateFn returns a transform: Seq[T] => Seq[T], keeping at most n elements.
func TruncateFn[T any](n uint) func(Seq[T]) Seq[T] {
	return func(s Seq[T]) Seq[T] { return s.Truncate(n) }
}

// TakeWhileFn returns a transform: Seq[T] => Seq[T], yielding while fn returns true.
func TakeWhileFn[T any](fn func(T) bool) func(Seq[T]) Seq[T] {
	return func(s Seq[T]) Seq[T] { return s.TakeWhile(fn) }
}

// SkipFn returns a transform: Seq[T] => Seq[T], skipping the first n elements.
func SkipFn[T any](n uint) func(Seq[T]) Seq[T] {
	return func(s Seq[T]) Seq[T] { return s.Skip(n) }
}

// SkipWhileFn returns a transform: Seq[T] => Seq[T], skipping while fn returns true.
func SkipWhileFn[T any](fn func(T) bool) func(Seq[T]) Seq[T] {
	return func(s Seq[T]) Seq[T] { return s.SkipWhile(fn) }
}

// AppendFn returns a transform: Seq[T] => Seq[T], appending other sequences.
func AppendFn[T any](others ...Seq[T]) func(Seq[T]) Seq[T] {
	return func(s Seq[T]) Seq[T] { return s.Append(others...) }
}

// SortWithFn returns a transform: Seq[T] => Seq[T], sorted using a comparison function.
func SortWithFn[T any](cmpFn func(T, T) int) func(Seq[T]) Seq[T] {
	return func(s Seq[T]) Seq[T] { return s.SortWith(cmpFn) }
}

// RevFn returns a transform: Seq[T] => Seq[T], reversed.
func RevFn[T any]() func(Seq[T]) Seq[T] {
	return func(s Seq[T]) Seq[T] { return s.Rev() }
}

// CycleFn returns a transform: Seq[T] => Seq[T], repeating infinitely.
func CycleFn[T any]() func(Seq[T]) Seq[T] {
	return func(s Seq[T]) Seq[T] { return s.Cycle() }
}

// TailFn returns a transform: Seq[T] => Seq[T], dropping the first element.
func TailFn[T any]() func(Seq[T]) Seq[T] {
	return func(s Seq[T]) Seq[T] { return s.Tail() }
}

// -- curried type-changing helpers (for use with Then / pipe) ------------------
// These produce function values for use with pipe.Pipe* or as callbacks.
// For simple pipelines, prefer the equivalent methods: s.Map(fn), s.Choose(fn), etc.

// MapFn returns a transform: Seq[T] => Seq[R].
// Prefer s.Map(fn) for direct chaining; use MapFn when a function value is required.
func MapFn[T, R any](fn func(T) R) func(Seq[T]) Seq[R] {
	return func(s Seq[T]) Seq[R] { return s.Map(fn) }
}

// MapiFn returns a transform: Seq[T] => Seq[R] with index.
// Prefer s.Mapi(fn) for direct chaining; use MapiFn when a function value is required.
func MapiFn[T, R any](fn func(int, T) R) func(Seq[T]) Seq[R] {
	return func(s Seq[T]) Seq[R] { return s.Mapi(fn) }
}

// CollectFn returns a transform: Seq[T] => Seq[R] via one-to-many mapping.
// Prefer s.Collect(fn) for direct chaining; use CollectFn when a function value is required.
func CollectFn[T, R any](fn func(T) []R) func(Seq[T]) Seq[R] {
	return func(s Seq[T]) Seq[R] { return s.Collect(fn) }
}

// ChooseFn returns a transform: Seq[T] => Seq[R], keeping Some values.
// Prefer s.Choose(fn) for direct chaining; use ChooseFn when a function value is required.
func ChooseFn[T, R any](fn func(T) option.Option[R]) func(Seq[T]) Seq[R] {
	return func(s Seq[T]) Seq[R] { return s.Choose(fn) }
}

// SortByFn returns a transform: Seq[T] => Seq[T], sorted ascending by key.
// Prefer s.SortBy(key) for direct chaining; use SortByFn when a function value is required.
func SortByFn[T any, K cmp.Ordered](key func(T) K) func(Seq[T]) Seq[T] {
	return func(s Seq[T]) Seq[T] { return s.SortBy(key) }
}

// SortByDescendingFn returns a transform: Seq[T] => Seq[T], sorted descending by key.
// Prefer s.SortByDescending(key) for direct chaining; use SortByDescendingFn when a function value is required.
func SortByDescendingFn[T any, K cmp.Ordered](key func(T) K) func(Seq[T]) Seq[T] {
	return func(s Seq[T]) Seq[T] { return s.SortByDescending(key) }
}

// DistinctFn returns a transform: Seq[T] => Seq[T], removing duplicates.
func DistinctFn[T comparable]() func(Seq[T]) Seq[T] {
	return Distinct[T]
}

// DistinctByFn returns a transform: Seq[T] => Seq[T], removing duplicates by key.
// Prefer s.DistinctBy(key) for direct chaining; use DistinctByFn when a function value is required.
func DistinctByFn[T any, K comparable](key func(T) K) func(Seq[T]) Seq[T] {
	return func(s Seq[T]) Seq[T] { return s.DistinctBy(key) }
}

// IndexedFn returns a transform: Seq[T] => Seq[Pair[int, T]].
func IndexedFn[T any]() func(Seq[T]) Seq[Pair[int, T]] {
	return Indexed[T]
}

// PairwiseFn returns a transform: Seq[T] => Seq[Pair[T, T]].
func PairwiseFn[T any]() func(Seq[T]) Seq[Pair[T, T]] {
	return Pairwise[T]
}

// WindowedFn returns a transform: Seq[T] => Seq[[]T].
func WindowedFn[T any](size int) func(Seq[T]) Seq[[]T] {
	return func(s Seq[T]) Seq[[]T] { return Windowed(s, size) }
}

// ChunkBySizeFn returns a transform: Seq[T] => Seq[[]T].
func ChunkBySizeFn[T any](size int) func(Seq[T]) Seq[[]T] {
	return func(s Seq[T]) Seq[[]T] { return ChunkBySize(s, size) }
}

// ScanFn returns a transform: Seq[T] => Seq[S].
// Prefer s.Scan(initial, fn) for direct chaining; use ScanFn when a function value is required.
func ScanFn[T, S any](initial S, fn func(S, T) S) func(Seq[T]) Seq[S] {
	return func(s Seq[T]) Seq[S] { return s.Scan(initial, fn) }
}

// FoldFn returns a terminal: Seq[T] => A.
// Prefer s.Fold(initial, fn) for direct chaining; use FoldFn when a function value is required.
func FoldFn[T, A any](initial A, fn func(A, T) A) func(Seq[T]) A {
	return func(s Seq[T]) A { return s.Fold(initial, fn) }
}

// GroupByFn returns a terminal: Seq[T] => map[K][]T.
// Prefer s.GroupBy(key) for direct chaining; use GroupByFn when a function value is required.
func GroupByFn[T any, K comparable](key func(T) K) func(Seq[T]) map[K][]T {
	return func(s Seq[T]) map[K][]T { return s.GroupBy(key) }
}

// ExceptFn returns a transform: Seq[T] => Seq[T], excluding elements in the given set.
func ExceptFn[T comparable](exclusion Seq[T]) func(Seq[T]) Seq[T] {
	return func(s Seq[T]) Seq[T] { return Except(s, exclusion) }
}

// -- curried terminal helpers --------------------------------------------------

// ToSliceFn returns a terminal: Seq[T] => []T.
func ToSliceFn[T any]() func(Seq[T]) []T {
	return func(s Seq[T]) []T { return s.ToSlice() }
}

// LengthFn returns a terminal: Seq[T] => int.
func LengthFn[T any]() func(Seq[T]) int {
	return func(s Seq[T]) int { return s.Length() }
}

// ZipWithFn returns a transform: Seq[T] => Seq[R], combining with other using fn.
// Prefer s.ZipWith(other, fn) for direct chaining; use ZipWithFn when a function value is required.
func ZipWithFn[T, U, R any](other Seq[U], fn func(T, U) R) func(Seq[T]) Seq[R] {
	return func(s Seq[T]) Seq[R] { return s.ZipWith(other, fn) }
}

// ZipWith3Fn returns a transform: Seq[T] => Seq[R], combining with b and c using fn.
// Prefer s.ZipWith3(b, c, fn) for direct chaining; use ZipWith3Fn when a function value is required.
func ZipWith3Fn[T, U, V, R any](b Seq[U], c Seq[V], fn func(T, U, V) R) func(Seq[T]) Seq[R] {
	return func(s Seq[T]) Seq[R] { return s.ZipWith3(b, c, fn) }
}
