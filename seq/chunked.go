package seq

import (
	"iter"
	"slices"
)

// ChunkedSeq is a lazy sequence of []T chunks, returned by [Seq.Windowed],
// [Seq.ChunkBySize] and [Seq.SplitInto]. Using a distinct named type breaks the
// instantiation cycle that would arise if those methods returned Seq[[]T] directly
// (Go's type checker would need to instantiate Seq[[]T], then Seq[[][]T], etc.).
//
// To continue chaining with regular Seq operations, use the [ChunkedToSeq]
// package-level function; a method cannot be used here as it would re-introduce
// the same instantiation cycle.
type ChunkedSeq[T any] iter.Seq[[]T]

// ChunkedToSeq converts a ChunkedSeq[T] to a Seq[[]T] for further composition.
// It is intentionally a package-level function rather than a method to avoid an
// instantiation cycle in the type checker.
func ChunkedToSeq[T any](c ChunkedSeq[T]) Seq[[]T] { return Seq[[]T](c) }

// ToSlice materialises the chunks into a [][]T.
func (c ChunkedSeq[T]) ToSlice() [][]T { return slices.Collect(iter.Seq[[]T](c)) }

// ForEach calls fn for every chunk, consuming the sequence.
func (c ChunkedSeq[T]) ForEach(fn func([]T)) {
	for chunk := range c {
		fn(chunk)
	}
}
