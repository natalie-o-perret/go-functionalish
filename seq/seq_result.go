package seq

import "github.com/natalie-o-perret/go-functionalish/result"

// OfResult returns a Seq containing the Ok value, or an empty Seq if Err.
// Symmetric with OfOption.
func OfResult[T, E any](r result.Result[T, E]) Seq[T] {
	if r.IsOk() {
		return Singleton(r.Unwrap())
	}
	return Empty[T]()
}

