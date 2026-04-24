// Package gof is a cohesive, type-safe functional programming library for Go 1.24+.
//
// No reflection. No interface{}. Pure generics and lazy by default.
//
// It provides six packages:
//
//   - [github.com/natalie-o-perret/go-functional-ish/seq]: lazy Seq[T] sequence pipelines (filter, map, sort, group, zip, ...)
//   - [github.com/natalie-o-perret/go-functional-ish/pseq]: parallel sequence operations — Map, Filter, Reduce, GroupBy, ... using goroutine-per-chunk
//   - [github.com/natalie-o-perret/go-functional-ish/option]: Option[T] for explicit presence/absence, no nil
//   - [github.com/natalie-o-perret/go-functional-ish/result]: Result[T,E] for railway-oriented error handling
//   - [github.com/natalie-o-perret/go-functional-ish/validation]: Validation[T,E] for applicative error accumulation
//   - [github.com/natalie-o-perret/go-functional-ish/pipe]: Pipe2-Pipe8 / PipeEndoN and Compose2-Compose4 / ComposeEndoN for left-to-right value threading
package gofunctionalish
