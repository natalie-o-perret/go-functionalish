package seq

import (
	"cmp"

	"github.com/natalie-o-perret/go-functionalish/option"
)

// Then threads a Seq[T] through a type-changing function, enabling a more
// linear style when combined with curried helpers like MapFn, SortByFn, etc.
//
//	seq.Then(
//	    seq.OfSlice(cars).Filter(pred),
//	    seq.MapFn(func(c Car) string { return c.Owner }),
//	).ToSlice()
//
// Also composable with the pipe package:
//
//	pipe.Pipe3(seq.OfSlice(cars).Filter(pred), seq.MapFn(toOwner), seq.MapFn(toUpper))
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

// MapFn returns a transform: Seq[T] => Seq[R].
func MapFn[T, R any](fn func(T) R) func(Seq[T]) Seq[R] {
	return func(s Seq[T]) Seq[R] { return Map(s, fn) }
}

// MapiFn returns a transform: Seq[T] => Seq[R] with index.
func MapiFn[T, R any](fn func(int, T) R) func(Seq[T]) Seq[R] {
	return func(s Seq[T]) Seq[R] { return Mapi(s, fn) }
}

// CollectFn returns a transform: Seq[T] => Seq[R] via one-to-many mapping.
func CollectFn[T, R any](fn func(T) []R) func(Seq[T]) Seq[R] {
	return func(s Seq[T]) Seq[R] { return Collect(s, fn) }
}

// ChooseFn returns a transform: Seq[T] => Seq[R], keeping Some values.
func ChooseFn[T, R any](fn func(T) option.Option[R]) func(Seq[T]) Seq[R] {
	return func(s Seq[T]) Seq[R] { return Choose(s, fn) }
}

// SortByFn returns a transform: Seq[T] => Seq[T], sorted ascending by key.
func SortByFn[T any, K cmp.Ordered](key func(T) K) func(Seq[T]) Seq[T] {
	return func(s Seq[T]) Seq[T] { return SortBy(s, key) }
}

// SortByDescendingFn returns a transform: Seq[T] => Seq[T], sorted descending by key.
func SortByDescendingFn[T any, K cmp.Ordered](key func(T) K) func(Seq[T]) Seq[T] {
	return func(s Seq[T]) Seq[T] { return SortByDescending(s, key) }
}

// DistinctFn returns a transform: Seq[T] => Seq[T], removing duplicates.
func DistinctFn[T comparable]() func(Seq[T]) Seq[T] {
	return func(s Seq[T]) Seq[T] { return Distinct(s) }
}

// DistinctByFn returns a transform: Seq[T] => Seq[T], removing duplicates by key.
func DistinctByFn[T any, K comparable](key func(T) K) func(Seq[T]) Seq[T] {
	return func(s Seq[T]) Seq[T] { return DistinctBy(s, key) }
}

// ExceptFn returns a transform: Seq[T] => Seq[T], excluding elements in the given set.
func ExceptFn[T comparable](exclusion Seq[T]) func(Seq[T]) Seq[T] {
	return func(s Seq[T]) Seq[T] { return Except(s, exclusion) }
}

// IndexedFn returns a transform: Seq[T] => Seq[Pair[int, T]].
func IndexedFn[T any]() func(Seq[T]) Seq[Pair[int, T]] {
	return func(s Seq[T]) Seq[Pair[int, T]] { return Indexed(s) }
}

// PairwiseFn returns a transform: Seq[T] => Seq[Pair[T, T]].
func PairwiseFn[T any]() func(Seq[T]) Seq[Pair[T, T]] {
	return func(s Seq[T]) Seq[Pair[T, T]] { return Pairwise(s) }
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
func ScanFn[T, S any](initial S, fn func(S, T) S) func(Seq[T]) Seq[S] {
	return func(s Seq[T]) Seq[S] { return Scan(s, initial, fn) }
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

// FoldFn returns a terminal: Seq[T] => A.
func FoldFn[T, A any](initial A, fn func(A, T) A) func(Seq[T]) A {
	return func(s Seq[T]) A { return Fold(s, initial, fn) }
}

// GroupByFn returns a terminal: Seq[T] => map[K][]T.
func GroupByFn[T any, K comparable](key func(T) K) func(Seq[T]) map[K][]T {
	return func(s Seq[T]) map[K][]T { return GroupBy(s, key) }
}
