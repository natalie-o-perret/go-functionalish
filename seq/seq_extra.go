package seq

import (
	"cmp"
	"iter"
	"slices"

	"github.com/natalie-o-perret/go-functionalish/option"
)

// Numeric is a constraint for types supporting arithmetic operations.
type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Triple holds three values produced by Zip3.
type Triple[T, U, V any] struct {
	First  T
	Second U
	Third  V
}

// -- additional constructors ---------------------------------------------------

// Init generates a Seq of n elements using fn(index).
func Init[T any](n int, fn func(int) T) Seq[T] {
	return func(yield func(T) bool) {
		for i := range n {
			if !yield(fn(i)) {
				return
			}
		}
	}
}

// InitInfinite generates an infinite Seq using fn(index).
// Always pair with Truncate.
func InitInfinite[T any](fn func(int) T) Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; ; i++ {
			if !yield(fn(i)) {
				return
			}
		}
	}
}

// Singleton returns a Seq containing exactly one element.
func Singleton[T any](v T) Seq[T] {
	return func(yield func(T) bool) { yield(v) }
}

// -- additional lazy pipeline methods ------------------------------------------

// Tail returns all elements except the first.
func (s Seq[T]) Tail() Seq[T] { return s.Skip(1) }

// InsertAt inserts v at the given index lazily.
func (s Seq[T]) InsertAt(index int, v T) Seq[T] {
	return func(yield func(T) bool) {
		i := 0
		for x := range s {
			if i == index {
				if !yield(v) {
					return
				}
			}
			if !yield(x) {
				return
			}
			i++
		}
		if i == index {
			yield(v)
		}
	}
}

// RemoveAt removes the element at the given index lazily.
func (s Seq[T]) RemoveAt(index int) Seq[T] {
	return func(yield func(T) bool) {
		i := 0
		for v := range s {
			if i != index {
				if !yield(v) {
					return
				}
			}
			i++
		}
	}
}

// UpdateAt replaces the element at the given index with v lazily.
func (s Seq[T]) UpdateAt(index int, v T) Seq[T] {
	return func(yield func(T) bool) {
		i := 0
		for x := range s {
			if i == index {
				if !yield(v) {
					return
				}
			} else {
				if !yield(x) {
					return
				}
			}
			i++
		}
	}
}

// Permute reorders elements by index permutation. Materialises.
func (s Seq[T]) Permute(fn func(int) int) Seq[T] {
	items := slices.Collect(iter.Seq[T](s))
	result := make([]T, len(items))
	for i, v := range items {
		result[fn(i)] = v
	}
	return OfSlice(result)
}

// -- additional terminal methods -----------------------------------------------

// Iteri calls fn(index, element) for every element, consuming the sequence.
func (s Seq[T]) Iteri(fn func(int, T)) {
	i := 0
	for v := range s {
		fn(i, v)
		i++
	}
}

// IsEmpty returns true if the sequence has no elements.
func (s Seq[T]) IsEmpty() bool {
	for range s {
		return false
	}
	return true
}

// Item returns the element at the given index and true, or the zero value and false.
func (s Seq[T]) Item(index int) (T, bool) {
	i := 0
	for v := range s {
		if i == index {
			return v, true
		}
		i++
	}
	var zero T
	return zero, false
}

// TryFind returns the first element matching fn as an Option. Short-circuits.
func (s Seq[T]) TryFind(fn func(T) bool) option.Option[T] {
	for v := range s {
		if fn(v) {
			return option.Some(v)
		}
	}
	return option.None[T]()
}

// TryFindIndex returns the index of the first element matching fn as an Option. Short-circuits.
func (s Seq[T]) TryFindIndex(fn func(T) bool) option.Option[int] {
	i := 0
	for v := range s {
		if fn(v) {
			return option.Some(i)
		}
		i++
	}
	return option.None[int]()
}

// TryFindBack returns the last element matching fn as an Option. Materialises.
func (s Seq[T]) TryFindBack(fn func(T) bool) option.Option[T] {
	var found T
	ok := false
	for v := range s {
		if fn(v) {
			found, ok = v, true
		}
	}
	if ok {
		return option.Some(found)
	}
	return option.None[T]()
}

// TryFindIndexBack returns the index of the last element matching fn as an Option. Materialises.
func (s Seq[T]) TryFindIndexBack(fn func(T) bool) option.Option[int] {
	i, result := 0, -1
	for v := range s {
		if fn(v) {
			result = i
		}
		i++
	}
	if result >= 0 {
		return option.Some(result)
	}
	return option.None[int]()
}

// Reduce applies fn using the first element as the initial accumulator.
// Returns (result, true) or (zero, false) if empty.
func (s Seq[T]) Reduce(fn func(T, T) T) (T, bool) {
	first := true
	var acc T
	for v := range s {
		if first {
			acc, first = v, false
		} else {
			acc = fn(acc, v)
		}
	}
	if first {
		var zero T
		return zero, false
	}
	return acc, true
}

// ReduceBack applies fn from the right using the last element as the initial accumulator.
// Materialises.
func (s Seq[T]) ReduceBack(fn func(T, T) T) (T, bool) {
	items := slices.Collect(iter.Seq[T](s))
	if len(items) == 0 {
		var zero T
		return zero, false
	}
	acc := items[len(items)-1]
	for i := len(items) - 2; i >= 0; i-- {
		acc = fn(items[i], acc)
	}
	return acc, true
}

// ExactlyOne returns the single element and true if the sequence has exactly one element.
// Returns (zero, false) if empty or more than one.
func (s Seq[T]) ExactlyOne() (T, bool) {
	count := 0
	var val T
	for v := range s {
		val = v
		count++
		if count > 1 {
			var zero T
			return zero, false
		}
	}
	if count == 1 {
		return val, true
	}
	var zero T
	return zero, false
}

// TryExactlyOne returns the single element as an Option.
func (s Seq[T]) TryExactlyOne() option.Option[T] {
	v, ok := s.ExactlyOne()
	if ok {
		return option.Some(v)
	}
	return option.None[T]()
}

// -- additional type-transforming package-level functions -----------------------

// Mapi transforms Seq[T] => Seq[R] lazily, providing the index to fn.
func Mapi[T, R any](s Seq[T], fn func(int, T) R) Seq[R] {
	return func(yield func(R) bool) {
		i := 0
		for v := range s {
			if !yield(fn(i, v)) {
				return
			}
			i++
		}
	}
}

// Indexed pairs each element with its zero-based index.
func Indexed[T any](s Seq[T]) Seq[Pair[int, T]] {
	return Mapi(s, func(i int, v T) Pair[int, T] {
		return Pair[int, T]{First: i, Second: v}
	})
}

// Pairwise yields sliding pairs of consecutive elements.
func Pairwise[T any](s Seq[T]) Seq[Pair[T, T]] {
	return func(yield func(Pair[T, T]) bool) {
		first := true
		var prev T
		for v := range s {
			if !first {
				if !yield(Pair[T, T]{First: prev, Second: v}) {
					return
				}
			}
			prev, first = v, false
		}
	}
}

// Windowed yields sliding windows of the given size as slices.
func Windowed[T any](s Seq[T], size int) Seq[[]T] {
	return func(yield func([]T) bool) {
		buf := make([]T, 0, size)
		for v := range s {
			if len(buf) < size {
				buf = append(buf, v)
			} else {
				copy(buf, buf[1:])
				buf[size-1] = v
			}
			if len(buf) == size {
				w := make([]T, size)
				copy(w, buf)
				if !yield(w) {
					return
				}
			}
		}
	}
}

// ChunkBySize yields non-overlapping chunks of the given size.
func ChunkBySize[T any](s Seq[T], size int) Seq[[]T] {
	return func(yield func([]T) bool) {
		chunk := make([]T, 0, size)
		for v := range s {
			chunk = append(chunk, v)
			if len(chunk) == size {
				if !yield(chunk) {
					return
				}
				chunk = make([]T, 0, size)
			}
		}
		if len(chunk) > 0 {
			yield(chunk)
		}
	}
}

// SplitInto splits the sequence into at most count roughly-equal chunks. Materialises.
func SplitInto[T any](s Seq[T], count int) Seq[[]T] {
	items := slices.Collect(iter.Seq[T](s))
	if count <= 0 || len(items) == 0 {
		return Empty[[]T]()
	}
	return func(yield func([]T) bool) {
		n := len(items)
		base, rem := n/count, n%count
		offset := 0
		for i := range count {
			sz := base
			if i < rem {
				sz++
			}
			if sz == 0 {
				break
			}
			chunk := make([]T, sz)
			copy(chunk, items[offset:offset+sz])
			if !yield(chunk) {
				return
			}
			offset += sz
		}
	}
}

// Scan is like Fold but yields each intermediate accumulator value, starting with initial.
func Scan[T, S any](s Seq[T], initial S, fn func(S, T) S) Seq[S] {
	return func(yield func(S) bool) {
		acc := initial
		if !yield(acc) {
			return
		}
		for v := range s {
			acc = fn(acc, v)
			if !yield(acc) {
				return
			}
		}
	}
}

// ScanBack is like FoldBack but yields each intermediate accumulator value. Materialises.
func ScanBack[T, S any](s Seq[T], initial S, fn func(T, S) S) Seq[S] {
	items := slices.Collect(iter.Seq[T](s))
	return func(yield func(S) bool) {
		results := make([]S, len(items)+1)
		results[len(items)] = initial
		for i := len(items) - 1; i >= 0; i-- {
			results[i] = fn(items[i], results[i+1])
		}
		for _, r := range results {
			if !yield(r) {
				return
			}
		}
	}
}

// TryPick applies fn to each element, returning the first Some result. Short-circuits.
func TryPick[T, R any](s Seq[T], fn func(T) option.Option[R]) option.Option[R] {
	for v := range s {
		if r := fn(v); r.IsSome() {
			return r
		}
	}
	return option.None[R]()
}

// Contains returns true if the sequence contains the given value. Short-circuits.
func Contains[T comparable](s Seq[T], value T) bool {
	for v := range s {
		if v == value {
			return true
		}
	}
	return false
}

// Except yields elements from s that are not in the exclusion sequence.
// The exclusion sequence is materialised once.
func Except[T comparable](s Seq[T], exclusion Seq[T]) Seq[T] {
	return func(yield func(T) bool) {
		set := make(map[T]struct{})
		for v := range exclusion {
			set[v] = struct{}{}
		}
		for v := range s {
			if _, ok := set[v]; !ok {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Sum returns the sum of all elements.
func Sum[T Numeric](s Seq[T]) T {
	var sum T
	for v := range s {
		sum += v
	}
	return sum
}

// SumBy returns the sum of fn(element) for all elements.
func SumBy[T any, N Numeric](s Seq[T], fn func(T) N) N {
	var sum N
	for v := range s {
		sum += fn(v)
	}
	return sum
}

// Min returns the minimum element and true, or (zero, false) if empty.
func Min[T cmp.Ordered](s Seq[T]) (T, bool) {
	first := true
	var m T
	for v := range s {
		if first || v < m {
			m, first = v, false
		}
	}
	return m, !first
}

// Max returns the maximum element and true, or (zero, false) if empty.
func Max[T cmp.Ordered](s Seq[T]) (T, bool) {
	first := true
	var m T
	for v := range s {
		if first || v > m {
			m, first = v, false
		}
	}
	return m, !first
}

// MinBy returns the element with the minimum key and true, or (zero, false) if empty.
func MinBy[T any, K cmp.Ordered](s Seq[T], fn func(T) K) (T, bool) {
	first := true
	var minVal T
	var minKey K
	for v := range s {
		k := fn(v)
		if first || k < minKey {
			minVal, minKey, first = v, k, false
		}
	}
	return minVal, !first
}

// MaxBy returns the element with the maximum key and true, or (zero, false) if empty.
func MaxBy[T any, K cmp.Ordered](s Seq[T], fn func(T) K) (T, bool) {
	first := true
	var maxVal T
	var maxKey K
	for v := range s {
		k := fn(v)
		if first || k > maxKey {
			maxVal, maxKey, first = v, k, false
		}
	}
	return maxVal, !first
}

// Average returns the arithmetic mean as float64, and true.
// Returns (0, false) if empty.
func Average[T Numeric](s Seq[T]) (float64, bool) {
	var sum T
	count := 0
	for v := range s {
		sum += v
		count++
	}
	if count == 0 {
		return 0, false
	}
	return float64(sum) / float64(count), true
}

// AverageBy returns the arithmetic mean of fn(element) as float64, and true.
// Returns (0, false) if empty.
func AverageBy[T any, N Numeric](s Seq[T], fn func(T) N) (float64, bool) {
	var sum N
	count := 0
	for v := range s {
		sum += fn(v)
		count++
	}
	if count == 0 {
		return 0, false
	}
	return float64(sum) / float64(count), true
}

// Zip3 lazily combines three Seqs element-wise. Stops at the shortest.
func Zip3[T, U, V any](a Seq[T], b Seq[U], c Seq[V]) Seq[Triple[T, U, V]] {
	return func(yield func(Triple[T, U, V]) bool) {
		nextB, stopB := iter.Pull(iter.Seq[U](b))
		defer stopB()
		nextC, stopC := iter.Pull(iter.Seq[V](c))
		defer stopC()
		for v := range a {
			u, okU := nextB()
			if !okU {
				return
			}
			w, okW := nextC()
			if !okW {
				return
			}
			if !yield(Triple[T, U, V]{First: v, Second: u, Third: w}) {
				return
			}
		}
	}
}

// Unzip splits a Seq of Pairs into two slices. Materialises.
func Unzip[T, U any](s Seq[Pair[T, U]]) ([]T, []U) {
	var ts []T
	var us []U
	for p := range s {
		ts = append(ts, p.First)
		us = append(us, p.Second)
	}
	return ts, us
}

// Concat flattens a Seq of Seqs into a single Seq.
func Concat[T any](seqs Seq[Seq[T]]) Seq[T] {
	return func(yield func(T) bool) {
		for inner := range seqs {
			for v := range inner {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// AllPairs yields the cartesian product of two Seqs. Materialises b.
func AllPairs[T, U any](a Seq[T], b Seq[U]) Seq[Pair[T, U]] {
	return func(yield func(Pair[T, U]) bool) {
		bItems := slices.Collect(iter.Seq[U](b))
		for va := range a {
			for _, vb := range bItems {
				if !yield(Pair[T, U]{First: va, Second: vb}) {
					return
				}
			}
		}
	}
}

// Map2 applies fn to pairs of elements from two Seqs lazily. Stops at the shorter.
func Map2[T, U, V any](a Seq[T], b Seq[U], fn func(T, U) V) Seq[V] {
	return func(yield func(V) bool) {
		nextB, stopB := iter.Pull(iter.Seq[U](b))
		defer stopB()
		for v := range a {
			u, ok := nextB()
			if !ok {
				return
			}
			if !yield(fn(v, u)) {
				return
			}
		}
	}
}

// Iter2 calls fn for pairs of elements from two Seqs. Stops at the shorter.
func Iter2[T, U any](a Seq[T], b Seq[U], fn func(T, U)) {
	nextB, stopB := iter.Pull(iter.Seq[U](b))
	defer stopB()
	for v := range a {
		u, ok := nextB()
		if !ok {
			return
		}
		fn(v, u)
	}
}

// Exists2 returns true if any parallel pair satisfies fn. Short-circuits.
func Exists2[T, U any](a Seq[T], b Seq[U], fn func(T, U) bool) bool {
	nextB, stopB := iter.Pull(iter.Seq[U](b))
	defer stopB()
	for v := range a {
		u, ok := nextB()
		if !ok {
			return false
		}
		if fn(v, u) {
			return true
		}
	}
	return false
}

// ForAll2 returns true if all parallel pairs satisfy fn. Short-circuits.
func ForAll2[T, U any](a Seq[T], b Seq[U], fn func(T, U) bool) bool {
	nextB, stopB := iter.Pull(iter.Seq[U](b))
	defer stopB()
	for v := range a {
		u, ok := nextB()
		if !ok {
			return true
		}
		if !fn(v, u) {
			return false
		}
	}
	return true
}

// CompareWith compares two Seqs element-wise using cmpFn.
// Returns negative if a < b, zero if equal, positive if a > b.
func CompareWith[T any](a, b Seq[T], cmpFn func(T, T) int) int {
	nextB, stopB := iter.Pull(iter.Seq[T](b))
	defer stopB()
	for va := range a {
		vb, ok := nextB()
		if !ok {
			return 1 // a is longer
		}
		if c := cmpFn(va, vb); c != 0 {
			return c
		}
	}
	if _, ok := nextB(); ok {
		return -1 // b is longer
	}
	return 0
}

// FoldBack folds from the right with an accumulator. Materialises.
func FoldBack[T, S any](s Seq[T], initial S, fn func(T, S) S) S {
	items := slices.Collect(iter.Seq[T](s))
	acc := initial
	for i := len(items) - 1; i >= 0; i-- {
		acc = fn(items[i], acc)
	}
	return acc
}

// Transpose transposes a Seq of Seqs (rows to columns). Materialises.
func Transpose[T any](s Seq[Seq[T]]) Seq[Seq[T]] {
	var rows [][]T
	for inner := range s {
		rows = append(rows, slices.Collect(iter.Seq[T](inner)))
	}
	if len(rows) == 0 {
		return Empty[Seq[T]]()
	}
	cols := 0
	for _, row := range rows {
		if len(row) > cols {
			cols = len(row)
		}
	}
	return func(yield func(Seq[T]) bool) {
		for c := range cols {
			col := func(y func(T) bool) {
				for _, row := range rows {
					if c < len(row) {
						if !y(row[c]) {
							return
						}
					}
				}
			}
			if !yield(Seq[T](col)) {
				return
			}
		}
	}
}

// MapFold combines map and fold in one pass. Materialises.
// Returns the mapped results as a slice and the final state.
func MapFold[S, T, R any](s Seq[T], initial S, fn func(S, T) (R, S)) ([]R, S) {
	var results []R
	acc := initial
	for v := range s {
		var r R
		r, acc = fn(acc, v)
		results = append(results, r)
	}
	return results, acc
}

// -- additional constructors ---------------------------------------------------

// Unfold generates a Seq by repeatedly applying fn to a seed state.
// fn returns None to stop, or Some(Pair{value, nextState}) to continue.
//
// Example: Fibonacci
//
//	fibs := seq.Unfold([2]int{0, 1}, func(s [2]int) option.Option[seq.Pair[int, [2]int]] {
//	    return option.Some(seq.Pair[int, [2]int]{First: s[0], Second: [2]int{s[1], s[0] + s[1]}})
//	})
func Unfold[S, T any](seed S, fn func(S) option.Option[Pair[T, S]]) Seq[T] {
	return func(yield func(T) bool) {
		state := seed
		for {
			result := fn(state)
			if result.IsNone() {
				return
			}
			p := result.Unwrap()
			if !yield(p.First) {
				return
			}
			state = p.Second
		}
	}
}

// RangeStep produces a Seq of integers from start (inclusive) to end (exclusive) with a given step.
// Panics if step is zero.
func RangeStep(start, end, step int) Seq[int] {
	if step == 0 {
		panic("seq: RangeStep step must not be zero")
	}
	return func(yield func(int) bool) {
		for i := start; (step > 0 && i < end) || (step < 0 && i > end); i += step {
			if !yield(i) {
				return
			}
		}
	}
}

// OfMap returns a Seq of Pair[K,V] from a map. Iteration order is not guaranteed.
func OfMap[K comparable, V any](m map[K]V) Seq[Pair[K, V]] {
	return func(yield func(Pair[K, V]) bool) {
		for k, v := range m {
			if !yield(Pair[K, V]{First: k, Second: v}) {
				return
			}
		}
	}
}

// OfOption returns a Seq with one element if o is Some, or an empty Seq if None.
func OfOption[T any](o option.Option[T]) Seq[T] {
	if o.IsSome() {
		return Singleton(o.Unwrap())
	}
	return Empty[T]()
}

// -- additional pipeline functions ---------------------------------------------

// Interleave alternates elements from a and b: [a1,b1,a2,b2,...].
// Stops when either sequence is exhausted.
func Interleave[T any](a, b Seq[T]) Seq[T] {
	return func(yield func(T) bool) {
		nextB, stopB := iter.Pull(iter.Seq[T](b))
		defer stopB()
		for va := range a {
			if !yield(va) {
				return
			}
			vb, ok := nextB()
			if !ok {
				return
			}
			if !yield(vb) {
				return
			}
		}
	}
}

// -- additional terminal functions ---------------------------------------------

// Partition splits s into two slices in one pass:
// the first contains elements where fn returns true, the second where fn returns false.
func Partition[T any](s Seq[T], fn func(T) bool) ([]T, []T) {
	var yes, no []T
	for v := range s {
		if fn(v) {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}
	return yes, no
}

// CountByKey counts occurrences of each key produced by fn.
func CountByKey[T any, K comparable](s Seq[T], fn func(T) K) map[K]int {
	counts := make(map[K]int)
	for v := range s {
		counts[fn(v)]++
	}
	return counts
}
