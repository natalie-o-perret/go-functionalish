// Package kv provides functional operations over lazy key-value sequences (iter.Seq2[K,V]),
// complementing the seq package for map-shaped data.
//
// The central type is [Seq2], a named type over iter.Seq2[K,V]. K must be
// [comparable] so that terminal operations like [Seq2.Collect] and
// [Seq2.ContainsKey] can use it as a map key.
//
// Pipelines are built by chaining methods left-to-right:
//
//	kv.Of(m).
//	    Filter(func(k string, v int) bool { return v > 0 }).
//	    MapValues(func(v int) string { return strconv.Itoa(v) }).
//	    Collect()
//
// Generic methods require Go 1.27+.
package kv

import (
	"iter"

	"github.com/natalie-o-perret/go-functionalish/seq"
)

// Seq2 is a lazy key-value sequence, a named type over iter.Seq2[K,V].
// K must be comparable so that [Seq2.Collect] and [Seq2.ContainsKey] work.
type Seq2[K comparable, V any] iter.Seq2[K, V]

// ---------------------------------------------------------------------------
// Constructors
// ---------------------------------------------------------------------------

// Of wraps a Go map into a Seq2. Iteration order is not guaranteed.
func Of[K comparable, V any](m map[K]V) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m {
			if !yield(k, v) {
				return
			}
		}
	}
}

// FromSeq converts a seq.Seq of seq.Pair[K,V] into a Seq2.
func FromSeq[K comparable, V any](s seq.Seq[seq.Pair[K, V]]) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for p := range s {
			if !yield(p.First, p.Second) {
				return
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Same-type pipeline methods
// ---------------------------------------------------------------------------

// Filter yields only pairs for which fn returns true.
func (s Seq2[K, V]) Filter(fn func(K, V) bool) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range iter.Seq2[K, V](s) {
			if fn(k, v) {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Type-changing pipeline methods (Go 1.27 generic methods)
// ---------------------------------------------------------------------------

// MapValues transforms values lazily, keeping keys unchanged.
func (s Seq2[K, V]) MapValues[R any](fn func(V) R) Seq2[K, R] {
	return func(yield func(K, R) bool) {
		for k, v := range iter.Seq2[K, V](s) {
			if !yield(k, fn(v)) {
				return
			}
		}
	}
}

// MapKeys transforms keys lazily. Duplicate keys in the output are not resolved.
func (s Seq2[K, V]) MapKeys[K2 comparable](fn func(K) K2) Seq2[K2, V] {
	return func(yield func(K2, V) bool) {
		for k, v := range iter.Seq2[K, V](s) {
			if !yield(fn(k), v) {
				return
			}
		}
	}
}

// Fold reduces the sequence into a single value using fn, starting from initial.
func (s Seq2[K, V]) Fold[A any](initial A, fn func(A, K, V) A) A {
	acc := initial
	for k, v := range iter.Seq2[K, V](s) {
		acc = fn(acc, k, v)
	}
	return acc
}

// ---------------------------------------------------------------------------
// Terminal / projection methods
// ---------------------------------------------------------------------------

// Keys returns a seq.Seq of all keys.
func (s Seq2[K, V]) Keys() seq.Seq[K] {
	return func(yield func(K) bool) {
		for k := range iter.Seq2[K, V](s) {
			if !yield(k) {
				return
			}
		}
	}
}

// Values returns a seq.Seq of all values.
func (s Seq2[K, V]) Values() seq.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range iter.Seq2[K, V](s) {
			if !yield(v) {
				return
			}
		}
	}
}

// ToSeq converts the Seq2 into a seq.Seq of seq.Pair[K,V].
func (s Seq2[K, V]) ToSeq() seq.Seq[seq.Pair[K, V]] {
	return func(yield func(seq.Pair[K, V]) bool) {
		for k, v := range iter.Seq2[K, V](s) {
			if !yield(seq.Pair[K, V]{First: k, Second: v}) {
				return
			}
		}
	}
}

// Collect materialises the Seq2 into a map. Later pairs overwrite on duplicate keys.
func (s Seq2[K, V]) Collect() map[K]V {
	m := make(map[K]V)
	for k, v := range iter.Seq2[K, V](s) {
		m[k] = v
	}
	return m
}

// ContainsKey reports whether k appears as a key in s. Short-circuits on first match.
func (s Seq2[K, V]) ContainsKey(k K) bool {
	for key := range iter.Seq2[K, V](s) {
		if key == k {
			return true
		}
	}
	return false
}
