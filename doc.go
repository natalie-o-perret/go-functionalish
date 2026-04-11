// Package gof is a cohesive, type-safe functional programming library for Go 1.24+.
//
// Inspired by C# LINQ, F# sequences, Option and Result types.
// No reflection. No interface{}. Pure generics and lazy by default.
//
// It provides four packages:
//
//   - [github.com/natalie-o-perret/gof/seq]: lazy Seq[T] sequence pipelines (filter, map, sort, group, zip, ...)
//   - [github.com/natalie-o-perret/gof/option]: Option[T] for explicit presence/absence, no nil
//   - [github.com/natalie-o-perret/gof/result]: Result[T,E] for railway-oriented error handling
//   - [github.com/natalie-o-perret/gof/validation]: Validation[T,E] for applicative error accumulation
//   - [github.com/natalie-o-perret/gof/pipe]: Pipe2 through Pipe8 plus Compose for F#-style |> threading
package gof
