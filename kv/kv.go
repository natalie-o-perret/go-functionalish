// Package kv provides functional operations over lazy key-value sequences (iter.Seq2[K,V]),
// complementing the seq package for map-shaped data.
//
// The central type is [Seq2], a named alias over iter.Seq2[K,V] that supports
// lazy map, filter, fold and conversion to/from Go maps and [seq.Seq] pairs.
package kv

import (
	"iter"

	"github.com/natalie-o-perret/go-functionalish/seq"
)

// Seq2[K,V] is a lazy key-value sequence, a named type over iter.Seq2[K,V].
type Seq2[K, V any] iter.Seq2[K, V]

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

// Keys returns a seq.Seq of all keys from s.
func Keys[K, V any](s Seq2[K, V]) seq.Seq[K] {
	return func(yield func(K) bool) {
		for k, _ := range iter.Seq2[K, V](s) {
			if !yield(k) {
				return
			}
		}
	}
}

// Values returns a seq.Seq of all values from s.
func Values[K, V any](s Seq2[K, V]) seq.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range iter.Seq2[K, V](s) {
			if !yield(v) {
				return
			}
		}
	}
}

// ToSeq converts a Seq2 into a seq.Seq of seq.Pair[K,V].
func ToSeq[K, V any](s Seq2[K, V]) seq.Seq[seq.Pair[K, V]] {
	return func(yield func(seq.Pair[K, V]) bool) {
		for k, v := range iter.Seq2[K, V](s) {
			if !yield(seq.Pair[K, V]{First: k, Second: v}) {
				return
			}
		}
	}
}

// FromSeq converts a seq.Seq of seq.Pair[K,V] into a Seq2.
func FromSeq[K, V any](s seq.Seq[seq.Pair[K, V]]) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for p := range s {
			if !yield(p.First, p.Second) {
				return
			}
		}
	}
}

// MapValues transforms values lazily, keeping keys unchanged.
func MapValues[K, V, R any](s Seq2[K, V], fn func(V) R) Seq2[K, R] {
	return func(yield func(K, R) bool) {
		for k, v := range iter.Seq2[K, V](s) {
			if !yield(k, fn(v)) {
				return
			}
		}
	}
}

// MapKeys transforms keys lazily. Duplicate keys in the output are not resolved.
func MapKeys[K1, K2, V any](s Seq2[K1, V], fn func(K1) K2) Seq2[K2, V] {
	return func(yield func(K2, V) bool) {
		for k, v := range iter.Seq2[K1, V](s) {
			if !yield(fn(k), v) {
				return
			}
		}
	}
}

// Filter yields only pairs for which fn returns true.
func Filter[K, V any](s Seq2[K, V], fn func(K, V) bool) Seq2[K, V] {
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

// Fold reduces a Seq2 into a single value using fn, starting from initial.
func Fold[K, V, A any](s Seq2[K, V], initial A, fn func(A, K, V) A) A {
	acc := initial
	for k, v := range iter.Seq2[K, V](s) {
		acc = fn(acc, k, v)
	}
	return acc
}

// Collect materialises a Seq2 into a map. Later pairs overwrite on duplicate keys.
func Collect[K comparable, V any](s Seq2[K, V]) map[K]V {
	m := make(map[K]V)
	for k, v := range iter.Seq2[K, V](s) {
		m[k] = v
	}
	return m
}

// ContainsKey reports whether k appears as a key in s. Short-circuits.
func ContainsKey[K comparable, V any](s Seq2[K, V], k K) bool {
	for key, _ := range iter.Seq2[K, V](s) {
		if key == k {
			return true
		}
	}
	return false
}
