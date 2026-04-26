// Package gofunctionalish is a cohesive, type-safe functional programming library for Go 1.24+.
//
// No reflection. No interface{}. Pure generics and lazy by default.
//
// It provides seven packages:
//
//   - [github.com/natalie-o-perret/go-functionalish/seq]: lazy Seq[T] sequence pipelines (filter, map, sort, group, zip, ...)
//   - [github.com/natalie-o-perret/go-functionalish/pseq]: parallel sequence operations — Map, Filter, Reduce, GroupBy, ... using goroutine-per-chunk
//   - [github.com/natalie-o-perret/go-functionalish/option]: Option[T] for explicit presence/absence, no nil
//   - [github.com/natalie-o-perret/go-functionalish/result]: Result[T,E] for railway-oriented error handling
//   - [github.com/natalie-o-perret/go-functionalish/validation]: Validation[T,E] for applicative error accumulation
//   - [github.com/natalie-o-perret/go-functionalish/pipe]: Pipe2-Pipe8 / PipeEndoN and Compose2-Compose4 / ComposeEndoN for left-to-right value threading
//   - [github.com/natalie-o-perret/go-functionalish/kv]: lazy Seq2[K,V] key-value pipelines over iter.Seq2 and Go maps
package gofunctionalish
