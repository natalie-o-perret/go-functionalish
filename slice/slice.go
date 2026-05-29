// Package slice provides a named Slice[T] type over []T with fluent generic
// methods. Type-changing operations such as Map, Choose, Fold, and GroupBy
// are methods, enabling left-to-right pipelines without helper functions:
//
//	slice.Of(users...).
//	    Filter(isActive).
//	    Map(toName).
//	    SortBy(strings.ToLower)
//
// Generic methods require Go 1.27+.
package slice

import (
	"cmp"
	"slices"

	"github.com/natalie-o-perret/go-functionalish/option"
)

// Slice is a named type over []T so that methods (including generic ones) can
// be attached directly without a wrapper struct.
type Slice[T any] []T

// -- constructors --------------------------------------------------------------

// Of creates a Slice from variadic elements.
func Of[T any](items ...T) Slice[T] { return Slice[T](items) }

// FromSlice wraps an existing []T in a Slice without copying.
func FromSlice[T any](items []T) Slice[T] { return Slice[T](items) }

// Replicate returns a Slice containing v repeated n times.
func Replicate[T any](v T, n int) Slice[T] {
	s := make(Slice[T], n)
	for i := range n {
		s[i] = v
	}
	return s
}

// Init returns a Slice of length n where element i is fn(i).
func Init[T any](n int, fn func(int) T) Slice[T] {
	s := make(Slice[T], n)
	for i := range n {
		s[i] = fn(i)
	}
	return s
}

// -- same-type pipeline methods ------------------------------------------------

// Filter returns a new Slice containing only elements for which fn returns true.
func (s Slice[T]) Filter(fn func(T) bool) Slice[T] {
	out := make(Slice[T], 0, len(s))
	for _, v := range s {
		if fn(v) {
			out = append(out, v)
		}
	}
	return out
}

// Exclude returns elements for which fn returns false (inverse of Filter).
func (s Slice[T]) Exclude(fn func(T) bool) Slice[T] {
	return s.Filter(func(v T) bool { return !fn(v) })
}

// Append returns a new Slice with the given items appended.
func (s Slice[T]) Append(items ...T) Slice[T] {
	out := make(Slice[T], len(s), len(s)+len(items))
	copy(out, s)
	return append(out, items...)
}

// Truncate returns the first n elements (or the whole slice if n >= len).
func (s Slice[T]) Truncate(n int) Slice[T] {
	if n >= len(s) {
		return slices.Clone([]T(s))
	}
	return slices.Clone([]T(s[:n]))
}

// Skip returns a new Slice with the first n elements dropped.
func (s Slice[T]) Skip(n int) Slice[T] {
	if n >= len(s) {
		return Slice[T]{}
	}
	return slices.Clone([]T(s[n:]))
}

// TakeWhile returns elements from the front while fn returns true.
func (s Slice[T]) TakeWhile(fn func(T) bool) Slice[T] {
	for i, v := range s {
		if !fn(v) {
			return slices.Clone([]T(s[:i]))
		}
	}
	return slices.Clone([]T(s))
}

// SkipWhile drops elements from the front while fn returns true.
func (s Slice[T]) SkipWhile(fn func(T) bool) Slice[T] {
	for i, v := range s {
		if !fn(v) {
			return slices.Clone([]T(s[i:]))
		}
	}
	return Slice[T]{}
}

// Rev returns a new Slice with elements in reversed order.
func (s Slice[T]) Rev() Slice[T] {
	out := make(Slice[T], len(s))
	for i, v := range s {
		out[len(s)-1-i] = v
	}
	return out
}

// SortWith returns a new Slice sorted according to cmpFn.
func (s Slice[T]) SortWith(cmpFn func(a, b T) int) Slice[T] {
	out := slices.Clone([]T(s))
	slices.SortFunc(out, cmpFn)
	return Slice[T](out)
}

// Tail returns all elements except the first.
func (s Slice[T]) Tail() Slice[T] { return s.Skip(1) }

// InsertAt returns a new Slice with v inserted before index.
func (s Slice[T]) InsertAt(index int, v T) Slice[T] {
	out := make(Slice[T], len(s)+1)
	copy(out, s[:index])
	out[index] = v
	copy(out[index+1:], s[index:])
	return out
}

// RemoveAt returns a new Slice with the element at index removed.
func (s Slice[T]) RemoveAt(index int) Slice[T] {
	out := make(Slice[T], len(s)-1)
	copy(out, s[:index])
	copy(out[index:], s[index+1:])
	return out
}

// UpdateAt returns a new Slice with the element at index replaced by v.
func (s Slice[T]) UpdateAt(index int, v T) Slice[T] {
	out := slices.Clone([]T(s))
	out[index] = v
	return Slice[T](out)
}

// -- terminal / query methods --------------------------------------------------

// Iter calls fn for each element.
func (s Slice[T]) Iter(fn func(T)) {
	for _, v := range s {
		fn(v)
	}
}

// Iteri calls fn with the index and each element.
func (s Slice[T]) Iteri(fn func(int, T)) {
	for i, v := range s {
		fn(i, v)
	}
}

// IsEmpty reports whether the slice has no elements.
func (s Slice[T]) IsEmpty() bool { return len(s) == 0 }

// Length returns the number of elements.
func (s Slice[T]) Length() int { return len(s) }

// Head returns the first element and true, or the zero value and false if empty.
func (s Slice[T]) Head() (T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, false
	}
	return s[0], true
}

// TryHead returns the first element as an Option.
func (s Slice[T]) TryHead() option.Option[T] {
	if len(s) == 0 {
		return option.None[T]()
	}
	return option.Some(s[0])
}

// Last returns the last element and true, or the zero value and false if empty.
func (s Slice[T]) Last() (T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, false
	}
	return s[len(s)-1], true
}

// TryLast returns the last element as an Option.
func (s Slice[T]) TryLast() option.Option[T] {
	if len(s) == 0 {
		return option.None[T]()
	}
	return option.Some(s[len(s)-1])
}

// Item returns the element at index as an Option (None if out of bounds).
func (s Slice[T]) Item(index int) option.Option[T] {
	if index < 0 || index >= len(s) {
		return option.None[T]()
	}
	return option.Some(s[index])
}

// CountBy returns the number of elements satisfying fn.
func (s Slice[T]) CountBy(fn func(T) bool) int {
	n := 0
	for _, v := range s {
		if fn(v) {
			n++
		}
	}
	return n
}

// Exists reports whether any element satisfies fn.
func (s Slice[T]) Exists(fn func(T) bool) bool {
	for _, v := range s {
		if fn(v) {
			return true
		}
	}
	return false
}

// ForAll reports whether every element satisfies fn.
func (s Slice[T]) ForAll(fn func(T) bool) bool {
	for _, v := range s {
		if !fn(v) {
			return false
		}
	}
	return true
}

// TryFind returns the first element satisfying fn, or None.
func (s Slice[T]) TryFind(fn func(T) bool) option.Option[T] {
	for _, v := range s {
		if fn(v) {
			return option.Some(v)
		}
	}
	return option.None[T]()
}

// TryFindBack returns the last element satisfying fn, or None.
func (s Slice[T]) TryFindBack(fn func(T) bool) option.Option[T] {
	for i := len(s) - 1; i >= 0; i-- {
		if fn(s[i]) {
			return option.Some(s[i])
		}
	}
	return option.None[T]()
}

// Reduce reduces the slice left-to-right with fn; returns (zero, false) if empty.
func (s Slice[T]) Reduce(fn func(T, T) T) (T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, false
	}
	acc := s[0]
	for _, v := range s[1:] {
		acc = fn(acc, v)
	}
	return acc, true
}

// ToSlice returns the underlying []T.
func (s Slice[T]) ToSlice() []T { return []T(s) }

// -- generic methods (Go 1.27+) -----------------------------------------------

// Map returns a new Slice with each element transformed by fn.
func (s Slice[T]) Map[R any](fn func(T) R) Slice[R] {
	out := make(Slice[R], len(s))
	for i, v := range s {
		out[i] = fn(v)
	}
	return out
}

// Mapi returns a new Slice with each (index, element) pair transformed by fn.
func (s Slice[T]) Mapi[R any](fn func(int, T) R) Slice[R] {
	out := make(Slice[R], len(s))
	for i, v := range s {
		out[i] = fn(i, v)
	}
	return out
}

// Collect applies fn to each element and flattens the results.
func (s Slice[T]) Collect[R any](fn func(T) []R) Slice[R] {
	var out Slice[R]
	for _, v := range s {
		out = append(out, fn(v)...)
	}
	return out
}

// Choose applies fn to each element and collects the Some values.
func (s Slice[T]) Choose[R any](fn func(T) option.Option[R]) Slice[R] {
	var out Slice[R]
	for _, v := range s {
		if r := fn(v); r.IsSome() {
			out = append(out, r.Unwrap())
		}
	}
	return out
}

// Fold accumulates a result by applying fn left-to-right starting from initial.
func (s Slice[T]) Fold[A any](initial A, fn func(A, T) A) A {
	acc := initial
	for _, v := range s {
		acc = fn(acc, v)
	}
	return acc
}

// FoldBack accumulates a result by applying fn right-to-left starting from initial.
func (s Slice[T]) FoldBack[A any](initial A, fn func(T, A) A) A {
	acc := initial
	for i := len(s) - 1; i >= 0; i-- {
		acc = fn(s[i], acc)
	}
	return acc
}

// GroupBy groups elements by the key produced by key.
func (s Slice[T]) GroupBy[K comparable](key func(T) K) map[K][]T {
	m := make(map[K][]T)
	for _, v := range s {
		k := key(v)
		m[k] = append(m[k], v)
	}
	return m
}

// SortBy returns a new Slice sorted ascending by the key produced by key.
func (s Slice[T]) SortBy[K cmp.Ordered](key func(T) K) Slice[T] {
	out := slices.Clone([]T(s))
	slices.SortFunc(out, func(a, b T) int { return cmp.Compare(key(a), key(b)) })
	return Slice[T](out)
}

// SortByDescending returns a new Slice sorted descending by the key produced by key.
func (s Slice[T]) SortByDescending[K cmp.Ordered](key func(T) K) Slice[T] {
	out := slices.Clone([]T(s))
	slices.SortFunc(out, func(a, b T) int { return cmp.Compare(key(b), key(a)) })
	return Slice[T](out)
}

// DistinctBy returns a new Slice keeping only the first element with each key.
func (s Slice[T]) DistinctBy[K comparable](key func(T) K) Slice[T] {
	seen := make(map[K]struct{})
	var out Slice[T]
	for _, v := range s {
		k := key(v)
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}

// TryPick applies fn to each element, returning the first Some result or None.
func (s Slice[T]) TryPick[R any](fn func(T) option.Option[R]) option.Option[R] {
	for _, v := range s {
		if r := fn(v); r.IsSome() {
			return r
		}
	}
	return option.None[R]()
}

// Scan returns a Slice of accumulated states applying fn left-to-right.
// The first element is initial; length is len(s)+1.
func (s Slice[T]) Scan[S any](initial S, fn func(S, T) S) Slice[S] {
	out := make(Slice[S], len(s)+1)
	out[0] = initial
	for i, v := range s {
		out[i+1] = fn(out[i], v)
	}
	return out
}

// ScanBack returns a Slice of accumulated states applying fn right-to-left.
// The last element is initial; length is len(s)+1.
func (s Slice[T]) ScanBack[S any](initial S, fn func(T, S) S) Slice[S] {
	out := make(Slice[S], len(s)+1)
	out[len(s)] = initial
	for i := len(s) - 1; i >= 0; i-- {
		out[i] = fn(s[i], out[i+1])
	}
	return out
}

// MapFold transforms elements and threads state left-to-right, returning the
// transformed elements and the final state.
func (s Slice[T]) MapFold[S, R any](initial S, fn func(S, T) (R, S)) (Slice[R], S) {
	out := make(Slice[R], len(s))
	acc := initial
	for i, v := range s {
		out[i], acc = fn(acc, v)
	}
	return out, acc
}

// -- package-level functions (constraints prevent methods) ----------------------

// Zip pairs corresponding elements from a and b. Stops at the shorter slice.
//
// Note: Zip cannot be a method because returning Slice[option.Pair[T,U]] from a
// Slice[T] method creates an instantiation cycle in the Go 1.27 type checker.
func Zip[T, U any](a Slice[T], b Slice[U]) Slice[option.Pair[T, U]] {
	n := min(len(a), len(b))
	out := make(Slice[option.Pair[T, U]], n)
	for i := range n {
		out[i] = option.Pair[T, U]{First: a[i], Second: b[i]}
	}
	return out
}

// Distinct returns a new Slice with duplicate elements removed, keeping the
// first occurrence.
func Distinct[T comparable](s Slice[T]) Slice[T] {
	seen := make(map[T]struct{}, len(s))
	var out Slice[T]
	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}

// Contains reports whether s contains value.
func Contains[T comparable](s Slice[T], value T) bool {
	for _, v := range s {
		if v == value {
			return true
		}
	}
	return false
}

// Except returns a new Slice with all elements present in exclusion removed.
func Except[T comparable](s, exclusion Slice[T]) Slice[T] {
	excl := make(map[T]struct{}, len(exclusion))
	for _, v := range exclusion {
		excl[v] = struct{}{}
	}
	var out Slice[T]
	for _, v := range s {
		if _, ok := excl[v]; !ok {
			out = append(out, v)
		}
	}
	return out
}
