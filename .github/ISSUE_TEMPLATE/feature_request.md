---
name: Feature request
about: Propose a new function, method, or package
title: "feat(<package>): <short description>"
labels: enhancement
assignees: ""
---

## Package

<!-- Which package would this belong to, or should a new package be created?
     e.g. seq, pseq, option, result, validation, pipe, kv -->

## Motivation

<!-- Why is this needed? What problem does it solve?
     Link to prior art (F#, Haskell, Rust Iterator, etc.) if applicable. -->

## Proposed API

```go
// Signature(s) you have in mind.
// Follow the existing convention: same-type operations are methods,
// type-transforming operations are package-level functions.

// Example:
func Zip[A, B any](a Seq[A], b Seq[B]) Seq[Pair[A, B]]
```

## Example usage

```go
// How would a caller use it?
pairs := seq.Zip(seq.OfSlice([]int{1, 2}), seq.OfSlice([]string{"a", "b"})).ToSlice()
// => [{1 a} {2 b}]
```

## Alternatives considered

<!-- Have you looked at workarounds using existing primitives?
     Why aren't they sufficient? -->

## Additional context

<!-- Benchmarks, links, screenshots, or anything else. -->

