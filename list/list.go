// Package list provides an immutable List[T] type backed by a private slice.
// Because the underlying slice is unexported, direct index-write (s[i] = v) is
// impossible from outside the package; all operations return a new List, giving
// F#-style list<T> semantics without the cache penalty of a linked list.
//
// Generic methods (Go 1.27+) allow type-changing operations (Map, Collect,
// Choose, Fold) to chain directly:
//
//	list.Of(1, 2, 3, 4, 5).
//	    Filter(func(n int) bool { return n%2 == 0 }).
//	    Map(func(n int) string { return strconv.Itoa(n) }).
//	    ToSlice() // ["2", "4"]
//
// To iterate, use [List.All]:
//
//	for v := range myList.All() { ... }
//
// For cases where direct indexing is convenient and immutability is not a
// concern, see the [slice] package's Slice[T] type instead.
package list

import (
	"cmp"
	"iter"
	"slices"

	"github.com/natalie-o-perret/go-functionalish/option"
)

// List[T] is an immutable list backed by a private slice.
// The zero value is an empty, ready-to-use List.
type List[T any] struct {
	items []T
}

// -- constructors --------------------------------------------------------------

// Of creates a List from variadic elements.
func Of[T any](items ...T) List[T] {
	cp := make([]T, len(items))
	copy(cp, items)
	return List[T]{items: cp}
}

// FromSlice creates a List from a slice (the slice is copied).
func FromSlice[T any](items []T) List[T] {
	cp := make([]T, len(items))
	copy(cp, items)
	return List[T]{items: cp}
}

// Replicate returns a List containing v repeated n times.
// Non-positive n returns an empty List.
func Replicate[T any](v T, n int) List[T] {
	if n <= 0 {
		return List[T]{}
	}
	items := make([]T, n)
	for i := range n {
		items[i] = v
	}
	return List[T]{items: items}
}

// Init returns a List of length n where element i is fn(i).
// Non-positive n returns an empty List.
func Init[T any](n int, fn func(int) T) List[T] {
	if n <= 0 {
		return List[T]{}
	}
	items := make([]T, n)
	for i := range n {
		items[i] = fn(i)
	}
	return List[T]{items: items}
}

// -- same-type pipeline methods ------------------------------------------------

// Filter returns a new List containing only elements for which fn returns true.
func (l List[T]) Filter(fn func(T) bool) List[T] {
	var out []T
	for _, v := range l.items {
		if fn(v) {
			out = append(out, v)
		}
	}
	return List[T]{items: out}
}

// Exclude returns elements for which fn returns false (inverse of Filter).
func (l List[T]) Exclude(fn func(T) bool) List[T] {
	return l.Filter(func(v T) bool { return !fn(v) })
}

// Append returns a new List with the given items appended.
func (l List[T]) Append(items ...T) List[T] {
	out := make([]T, len(l.items)+len(items))
	copy(out, l.items)
	copy(out[len(l.items):], items)
	return List[T]{items: out}
}

// Concat returns a new List with other Lists appended.
func (l List[T]) Concat(others ...List[T]) List[T] {
	total := len(l.items)
	for _, o := range others {
		total += len(o.items)
	}
	out := make([]T, 0, total)
	out = append(out, l.items...)
	for _, o := range others {
		out = append(out, o.items...)
	}
	return List[T]{items: out}
}

// Truncate returns the first n elements (or the whole list if n >= Length).
// Non-positive n returns an empty List.
func (l List[T]) Truncate(n int) List[T] {
	if n <= 0 {
		return List[T]{}
	}
	if n >= len(l.items) {
		return FromSlice(l.items)
	}
	cp := make([]T, n)
	copy(cp, l.items[:n])
	return List[T]{items: cp}
}

// Skip returns a new List with the first n elements dropped.
// Non-positive n returns a copy of the full List.
func (l List[T]) Skip(n int) List[T] {
	if n <= 0 {
		return FromSlice(l.items)
	}
	if n >= len(l.items) {
		return List[T]{}
	}
	cp := make([]T, len(l.items)-n)
	copy(cp, l.items[n:])
	return List[T]{items: cp}
}

// TakeWhile returns elements from the front while fn returns true.
func (l List[T]) TakeWhile(fn func(T) bool) List[T] {
	for i, v := range l.items {
		if !fn(v) {
			cp := make([]T, i)
			copy(cp, l.items[:i])
			return List[T]{items: cp}
		}
	}
	return FromSlice(l.items)
}

// SkipWhile drops elements from the front while fn returns true.
func (l List[T]) SkipWhile(fn func(T) bool) List[T] {
	for i, v := range l.items {
		if !fn(v) {
			cp := make([]T, len(l.items)-i)
			copy(cp, l.items[i:])
			return List[T]{items: cp}
		}
	}
	return List[T]{}
}

// Rev returns a new List with elements in reversed order.
func (l List[T]) Rev() List[T] {
	out := make([]T, len(l.items))
	for i, v := range l.items {
		out[len(l.items)-1-i] = v
	}
	return List[T]{items: out}
}

// SortWith returns a new List sorted according to cmpFn.
func (l List[T]) SortWith(cmpFn func(a, b T) int) List[T] {
	out := slices.Clone(l.items)
	slices.SortFunc(out, cmpFn)
	return List[T]{items: out}
}

// Tail returns all elements except the first.
func (l List[T]) Tail() List[T] { return l.Skip(1) }

// InsertAt returns a new List with v inserted before index.
// Panics if index is out of range [0, Length].
func (l List[T]) InsertAt(index int, v T) List[T] {
	if index < 0 || index > len(l.items) {
		panic("list.InsertAt: index out of range")
	}
	out := make([]T, len(l.items)+1)
	copy(out, l.items[:index])
	out[index] = v
	copy(out[index+1:], l.items[index:])
	return List[T]{items: out}
}

// RemoveAt returns a new List with the element at index removed.
// Panics if index is out of range [0, Length-1].
func (l List[T]) RemoveAt(index int) List[T] {
	if index < 0 || index >= len(l.items) {
		panic("list.RemoveAt: index out of range")
	}
	out := make([]T, len(l.items)-1)
	copy(out, l.items[:index])
	copy(out[index:], l.items[index+1:])
	return List[T]{items: out}
}

// UpdateAt returns a new List with the element at index replaced by v.
// Panics if index is out of range [0, Length-1].
func (l List[T]) UpdateAt(index int, v T) List[T] {
	if index < 0 || index >= len(l.items) {
		panic("list.UpdateAt: index out of range")
	}
	out := slices.Clone(l.items)
	out[index] = v
	return List[T]{items: out}
}

// -- terminal / query methods --------------------------------------------------

// Iter calls fn for each element.
func (l List[T]) Iter(fn func(T)) {
	for _, v := range l.items {
		fn(v)
	}
}

// Iteri calls fn with the index and each element.
func (l List[T]) Iteri(fn func(int, T)) {
	for i, v := range l.items {
		fn(i, v)
	}
}

// All returns an iter.Seq[T] over all elements in order.
func (l List[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range l.items {
			if !yield(v) {
				return
			}
		}
	}
}

// IsEmpty reports whether the list has no elements.
func (l List[T]) IsEmpty() bool { return len(l.items) == 0 }

// Length returns the number of elements.
func (l List[T]) Length() int { return len(l.items) }

// At returns the element at index as an Option (None if out of bounds).
func (l List[T]) At(index int) option.Option[T] {
	if index < 0 || index >= len(l.items) {
		return option.None[T]()
	}
	return option.Some(l.items[index])
}

// Head returns the first element and true, or the zero value and false if empty.
func (l List[T]) Head() (T, bool) {
	if len(l.items) == 0 {
		var zero T
		return zero, false
	}
	return l.items[0], true
}

// TryHead returns the first element as an Option.
func (l List[T]) TryHead() option.Option[T] {
	if len(l.items) == 0 {
		return option.None[T]()
	}
	return option.Some(l.items[0])
}

// Last returns the last element and true, or the zero value and false if empty.
func (l List[T]) Last() (T, bool) {
	if len(l.items) == 0 {
		var zero T
		return zero, false
	}
	return l.items[len(l.items)-1], true
}

// TryLast returns the last element as an Option.
func (l List[T]) TryLast() option.Option[T] {
	if len(l.items) == 0 {
		return option.None[T]()
	}
	return option.Some(l.items[len(l.items)-1])
}

// CountBy returns the number of elements satisfying fn.
func (l List[T]) CountBy(fn func(T) bool) int {
	n := 0
	for _, v := range l.items {
		if fn(v) {
			n++
		}
	}
	return n
}

// Exists reports whether any element satisfies fn.
func (l List[T]) Exists(fn func(T) bool) bool {
	for _, v := range l.items {
		if fn(v) {
			return true
		}
	}
	return false
}

// ForAll reports whether every element satisfies fn.
func (l List[T]) ForAll(fn func(T) bool) bool {
	for _, v := range l.items {
		if !fn(v) {
			return false
		}
	}
	return true
}

// TryFind returns the first element satisfying fn, or None.
func (l List[T]) TryFind(fn func(T) bool) option.Option[T] {
	for _, v := range l.items {
		if fn(v) {
			return option.Some(v)
		}
	}
	return option.None[T]()
}

// TryFindBack returns the last element satisfying fn, or None.
func (l List[T]) TryFindBack(fn func(T) bool) option.Option[T] {
	for i := len(l.items) - 1; i >= 0; i-- {
		if fn(l.items[i]) {
			return option.Some(l.items[i])
		}
	}
	return option.None[T]()
}

// Reduce reduces the list left-to-right with fn; returns (zero, false) if empty.
func (l List[T]) Reduce(fn func(T, T) T) (T, bool) {
	if len(l.items) == 0 {
		var zero T
		return zero, false
	}
	acc := l.items[0]
	for _, v := range l.items[1:] {
		acc = fn(acc, v)
	}
	return acc, true
}

// ToSlice returns a copy of the underlying elements as a []T.
func (l List[T]) ToSlice() []T {
	out := make([]T, len(l.items))
	copy(out, l.items)
	return out
}

// -- generic methods (Go 1.27+) ------------------------------------------------

// Map returns a new List with each element transformed by fn.
func (l List[T]) Map[R any](fn func(T) R) List[R] {
	out := make([]R, len(l.items))
	for i, v := range l.items {
		out[i] = fn(v)
	}
	return List[R]{items: out}
}

// Mapi returns a new List with each (index, element) pair transformed by fn.
func (l List[T]) Mapi[R any](fn func(int, T) R) List[R] {
	out := make([]R, len(l.items))
	for i, v := range l.items {
		out[i] = fn(i, v)
	}
	return List[R]{items: out}
}

// Collect applies fn to each element and flattens the results.
func (l List[T]) Collect[R any](fn func(T) []R) List[R] {
	var out []R
	for _, v := range l.items {
		out = append(out, fn(v)...)
	}
	return List[R]{items: out}
}

// Choose applies fn to each element and collects the Some values.
func (l List[T]) Choose[R any](fn func(T) option.Option[R]) List[R] {
	var out []R
	for _, v := range l.items {
		if r := fn(v); r.IsSome() {
			out = append(out, r.Unwrap())
		}
	}
	return List[R]{items: out}
}

// Fold accumulates a result by applying fn left-to-right starting from initial.
func (l List[T]) Fold[A any](initial A, fn func(A, T) A) A {
	acc := initial
	for _, v := range l.items {
		acc = fn(acc, v)
	}
	return acc
}

// FoldBack accumulates a result by applying fn right-to-left starting from initial.
func (l List[T]) FoldBack[A any](initial A, fn func(T, A) A) A {
	acc := initial
	for i := len(l.items) - 1; i >= 0; i-- {
		acc = fn(l.items[i], acc)
	}
	return acc
}

// GroupBy groups elements by the key produced by key.
func (l List[T]) GroupBy[K comparable](key func(T) K) map[K][]T {
	m := make(map[K][]T)
	for _, v := range l.items {
		k := key(v)
		m[k] = append(m[k], v)
	}
	return m
}

// SortBy returns a new List sorted ascending by the key produced by key.
func (l List[T]) SortBy[K cmp.Ordered](key func(T) K) List[T] {
	out := slices.Clone(l.items)
	slices.SortFunc(out, func(a, b T) int { return cmp.Compare(key(a), key(b)) })
	return List[T]{items: out}
}

// SortByDescending returns a new List sorted descending by the key produced by key.
func (l List[T]) SortByDescending[K cmp.Ordered](key func(T) K) List[T] {
	out := slices.Clone(l.items)
	slices.SortFunc(out, func(a, b T) int { return cmp.Compare(key(b), key(a)) })
	return List[T]{items: out}
}

// DistinctBy returns a new List keeping only the first element with each key.
func (l List[T]) DistinctBy[K comparable](key func(T) K) List[T] {
	seen := make(map[K]struct{})
	var out []T
	for _, v := range l.items {
		k := key(v)
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			out = append(out, v)
		}
	}
	return List[T]{items: out}
}

// TryPick applies fn to each element, returning the first Some result or None.
func (l List[T]) TryPick[R any](fn func(T) option.Option[R]) option.Option[R] {
	for _, v := range l.items {
		if r := fn(v); r.IsSome() {
			return r
		}
	}
	return option.None[R]()
}

// Scan returns a List of accumulated states applying fn left-to-right.
// The first element is initial; length is Length+1.
func (l List[T]) Scan[S any](initial S, fn func(S, T) S) List[S] {
	out := make([]S, len(l.items)+1)
	out[0] = initial
	for i, v := range l.items {
		out[i+1] = fn(out[i], v)
	}
	return List[S]{items: out}
}

// ScanBack returns a List of accumulated states applying fn right-to-left.
// The last element is initial; length is Length+1.
func (l List[T]) ScanBack[S any](initial S, fn func(T, S) S) List[S] {
	out := make([]S, len(l.items)+1)
	out[len(l.items)] = initial
	for i := len(l.items) - 1; i >= 0; i-- {
		out[i] = fn(l.items[i], out[i+1])
	}
	return List[S]{items: out}
}

// MapFold transforms elements and threads state left-to-right, returning the
// transformed elements and the final state.
func (l List[T]) MapFold[S, R any](initial S, fn func(S, T) (R, S)) (List[R], S) {
	out := make([]R, len(l.items))
	acc := initial
	for i, v := range l.items {
		out[i], acc = fn(acc, v)
	}
	return List[R]{items: out}, acc
}

// ZipWith combines elements from l and other pairwise using fn. Stops at the shorter List.
func (l List[T]) ZipWith[U, R any](other List[U], fn func(T, U) R) List[R] {
	n := min(len(l.items), len(other.items))
	out := make([]R, n)
	for i := range n {
		out[i] = fn(l.items[i], other.items[i])
	}
	return List[R]{items: out}
}

// ZipWith3 combines elements from l, b, and c element-wise using fn. Stops at the shortest List.
func (l List[T]) ZipWith3[U, V, R any](b List[U], c List[V], fn func(T, U, V) R) List[R] {
	n := min(len(l.items), len(b.items), len(c.items))
	out := make([]R, n)
	for i := range n {
		out[i] = fn(l.items[i], b.items[i], c.items[i])
	}
	return List[R]{items: out}
}

// -- package-level functions (constraint or multi-receiver prevents methods) ---

// Zip pairs corresponding elements from a and b. Stops at the shorter List.
//
// Note: Zip cannot be a method because returning List[option.Pair[T,U]] from a
// List[T] method creates an instantiation cycle in the Go 1.27 type checker.
func Zip[T, U any](a List[T], b List[U]) List[option.Pair[T, U]] {
	n := min(len(a.items), len(b.items))
	out := make([]option.Pair[T, U], n)
	for i := range n {
		out[i] = option.Pair[T, U]{First: a.items[i], Second: b.items[i]}
	}
	return List[option.Pair[T, U]]{items: out}
}

// Distinct returns a new List with duplicate elements removed, keeping the first occurrence.
func Distinct[T comparable](l List[T]) List[T] {
	seen := make(map[T]struct{}, len(l.items))
	var out []T
	for _, v := range l.items {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return List[T]{items: out}
}

// Contains reports whether l contains value.
func Contains[T comparable](l List[T], value T) bool {
	for _, v := range l.items {
		if v == value {
			return true
		}
	}
	return false
}

// Except returns a new List with all elements present in exclusion removed.
func Except[T comparable](l, exclusion List[T]) List[T] {
	excl := make(map[T]struct{}, len(exclusion.items))
	for _, v := range exclusion.items {
		excl[v] = struct{}{}
	}
	var out []T
	for _, v := range l.items {
		if _, ok := excl[v]; !ok {
			out = append(out, v)
		}
	}
	return List[T]{items: out}
}
