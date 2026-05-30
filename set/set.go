// Package set provides an immutable Set[T comparable] type with fluent generic
// methods (Go 1.27+). All operations that "modify" the set return a new Set;
// the receiver is never changed.
//
// Generic methods allow type-changing operations (Map, Collect, Choose, Fold)
// to chain directly:
//
//	set.Of(1, 2, 3, 4, 5).
//	    Filter(func(n int) bool { return n%2 == 0 }).
//	    Map(func(n int) string { return strconv.Itoa(n) }).
//	    Contains("2") // true
//
// Iteration order is unspecified (Go map randomisation). Use Fold only for
// operations that are associative and commutative (sum, count, etc.).
package set

import (
	"iter"

	"github.com/natalie-o-perret/go-functionalish/option"
)

// Set is an immutable set of comparable values backed by a map.
// The zero value is an empty, ready-to-use Set.
type Set[T comparable] struct {
	m map[T]struct{}
}

// -- constructors --------------------------------------------------------------

// Of creates a Set from variadic elements. Duplicates are silently dropped.
func Of[T comparable](items ...T) Set[T] {
	m := make(map[T]struct{}, len(items))
	for _, v := range items {
		m[v] = struct{}{}
	}
	return Set[T]{m: m}
}

// FromSlice creates a Set from a slice. Duplicates are silently dropped.
func FromSlice[T comparable](items []T) Set[T] { return Of(items...) }

// Empty returns an empty Set.
func Empty[T comparable]() Set[T] { return Set[T]{m: make(map[T]struct{})} }

// -- modification (all return a new Set) ---------------------------------------

// Add returns a new Set with the given items added.
func (s Set[T]) Add(items ...T) Set[T] {
	m := make(map[T]struct{}, len(s.m)+len(items))
	for v := range s.m {
		m[v] = struct{}{}
	}
	for _, v := range items {
		m[v] = struct{}{}
	}
	return Set[T]{m: m}
}

// Remove returns a new Set with the given items removed (absent items are ignored).
func (s Set[T]) Remove(items ...T) Set[T] {
	m := make(map[T]struct{}, len(s.m))
	for v := range s.m {
		m[v] = struct{}{}
	}
	for _, v := range items {
		delete(m, v)
	}
	return Set[T]{m: m}
}

// -- set algebra ---------------------------------------------------------------

// Union returns a new Set containing all elements from s and other (s ∪ other).
func (s Set[T]) Union(other Set[T]) Set[T] {
	m := make(map[T]struct{}, len(s.m)+len(other.m))
	for v := range s.m {
		m[v] = struct{}{}
	}
	for v := range other.m {
		m[v] = struct{}{}
	}
	return Set[T]{m: m}
}

// Intersect returns a new Set containing elements present in both s and other (s ∩ other).
func (s Set[T]) Intersect(other Set[T]) Set[T] {
	m := make(map[T]struct{})
	for v := range s.m {
		if _, ok := other.m[v]; ok {
			m[v] = struct{}{}
		}
	}
	return Set[T]{m: m}
}

// Difference returns elements in s not present in other (s ∖ other).
func (s Set[T]) Difference(other Set[T]) Set[T] {
	m := make(map[T]struct{})
	for v := range s.m {
		if _, ok := other.m[v]; !ok {
			m[v] = struct{}{}
		}
	}
	return Set[T]{m: m}
}

// SymmetricDifference returns elements in either s or other but not both (s △ other).
func (s Set[T]) SymmetricDifference(other Set[T]) Set[T] {
	return s.Difference(other).Union(other.Difference(s))
}

// -- query methods -------------------------------------------------------------

// Contains reports whether the set contains v.
func (s Set[T]) Contains(v T) bool {
	_, ok := s.m[v]
	return ok
}

// Len returns the number of elements.
func (s Set[T]) Len() int { return len(s.m) }

// IsEmpty reports whether the set has no elements.
func (s Set[T]) IsEmpty() bool { return len(s.m) == 0 }

// IsSubset reports whether every element of s is also in other (s ⊆ other).
func (s Set[T]) IsSubset(other Set[T]) bool {
	for v := range s.m {
		if _, ok := other.m[v]; !ok {
			return false
		}
	}
	return true
}

// IsSuperset reports whether every element of other is also in s (s ⊇ other).
func (s Set[T]) IsSuperset(other Set[T]) bool { return other.IsSubset(s) }

// Filter returns a new Set containing only elements for which fn returns true.
func (s Set[T]) Filter(fn func(T) bool) Set[T] {
	m := make(map[T]struct{})
	for v := range s.m {
		if fn(v) {
			m[v] = struct{}{}
		}
	}
	return Set[T]{m: m}
}

// Exclude returns a new Set keeping only elements for which fn returns false.
func (s Set[T]) Exclude(fn func(T) bool) Set[T] {
	return s.Filter(func(v T) bool { return !fn(v) })
}

// ForAll reports whether every element satisfies fn.
func (s Set[T]) ForAll(fn func(T) bool) bool {
	for v := range s.m {
		if !fn(v) {
			return false
		}
	}
	return true
}

// Exists reports whether any element satisfies fn.
func (s Set[T]) Exists(fn func(T) bool) bool {
	for v := range s.m {
		if fn(v) {
			return true
		}
	}
	return false
}

// CountBy returns the number of elements satisfying fn.
func (s Set[T]) CountBy(fn func(T) bool) int {
	n := 0
	for v := range s.m {
		if fn(v) {
			n++
		}
	}
	return n
}

// -- iteration / materialise ---------------------------------------------------

// Iter calls fn for every element. Iteration order is unspecified.
func (s Set[T]) Iter(fn func(T)) {
	for v := range s.m {
		fn(v)
	}
}

// All returns an iter.Seq[T] over all elements. Iteration order is unspecified.
func (s Set[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range s.m {
			if !yield(v) {
				return
			}
		}
	}
}

// ToSlice materialises the set into a []T. Order is unspecified.
func (s Set[T]) ToSlice() []T {
	out := make([]T, 0, len(s.m))
	for v := range s.m {
		out = append(out, v)
	}
	return out
}

// -- generic methods (Go 1.27+) ------------------------------------------------

// Map returns a new Set with each element transformed by fn.
// If fn maps two distinct elements to the same value, the duplicate is silently dropped.
func (s Set[T]) Map[R comparable](fn func(T) R) Set[R] {
	m := make(map[R]struct{}, len(s.m))
	for v := range s.m {
		m[fn(v)] = struct{}{}
	}
	return Set[R]{m: m}
}

// Collect applies fn to each element and merges all resulting slices into a new Set.
func (s Set[T]) Collect[R comparable](fn func(T) []R) Set[R] {
	m := make(map[R]struct{})
	for v := range s.m {
		for _, r := range fn(v) {
			m[r] = struct{}{}
		}
	}
	return Set[R]{m: m}
}

// Choose applies fn to each element, keeping Some values in a new Set.
func (s Set[T]) Choose[R comparable](fn func(T) option.Option[R]) Set[R] {
	m := make(map[R]struct{})
	for v := range s.m {
		if r := fn(v); r.IsSome() {
			m[r.Unwrap()] = struct{}{}
		}
	}
	return Set[R]{m: m}
}

// Fold accumulates a result over all elements. Because iteration order is
// non-deterministic, fn should be associative and commutative.
func (s Set[T]) Fold[A any](initial A, fn func(A, T) A) A {
	acc := initial
	for v := range s.m {
		acc = fn(acc, v)
	}
	return acc
}
