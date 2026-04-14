// Package pseq provides parallel counterparts to the most parallelism-friendly
// operations from [github.com/natalie-o-perret/gof/seq].
//
// Inspired by F#'s [FSharp.Collections.ParallelSeq] (PSeq) and Go's [lo/parallel].
//
// # Architecture
//
// Every function materialises its input [seq.Seq] into a slice, partitions
// it into roughly-equal chunks (one per CPU core by default), dispatches one
// goroutine per chunk via [sync.WaitGroup], and collects results.
// Order is always preserved.
//
// This is the same chunked strategy used by .NET's PLINQ (which F#'s PSeq
// wraps). It avoids the per-element goroutine spawning used by [lo/parallel],
// resulting in ~300× fewer allocations and 1.2-9× better throughput depending
// on workload. See the README for full benchmark tables.
//
// # Configuration
//
// Parallelism is configurable via [WithWorkers]; the default is
// [runtime.GOMAXPROCS](0). For lightweight lambdas (field access, arithmetic),
// sequential [seq] operations are faster — parallelism pays off when the
// per-element work is CPU-heavy (parsing, math, serialisation).
//
// # Omitted operations
//
// Operations that are inherently sequential (Scan, Unfold, TakeWhile, SkipWhile,
// Pairwise, Windowed, Zip, Cycle, Distinct) are intentionally omitted — use
// the [seq] package for those.
//
// [FSharp.Collections.ParallelSeq]: https://github.com/fsprojects/FSharp.Collections.ParallelSeq
// [lo/parallel]: https://github.com/samber/lo
package pseq

import (
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/natalie-o-perret/gof/option"
	"github.com/natalie-o-perret/gof/seq"
)

// ---------------------------------------------------------------------------
// Configuration
// ---------------------------------------------------------------------------

type config struct {
	workers int
}

// Option configures a parallel operation.
type Option func(*config)

// WithWorkers sets the number of goroutines used for parallel execution.
// Defaults to [runtime.GOMAXPROCS](0) when unset or when n ≤ 0.
func WithWorkers(n int) Option {
	return func(c *config) {
		if n > 0 {
			c.workers = n
		}
	}
}

func resolveConfig(opts []Option) config {
	c := config{workers: runtime.GOMAXPROCS(0)}
	for _, o := range opts {
		o(&c)
	}
	if c.workers <= 0 {
		c.workers = 1
	}
	return c
}

// ---------------------------------------------------------------------------
// Internal chunking helpers
// ---------------------------------------------------------------------------

type chunk struct{ offset, size int }

// computeChunks divides n items into count roughly-equal chunks.
// Mirrors the distribution logic in [seq.SplitInto].
func computeChunks(n, count int) []chunk {
	if count <= 0 || n <= 0 {
		return nil
	}
	if count > n {
		count = n
	}
	base, rem := n/count, n%count
	specs := make([]chunk, 0, count)
	offset := 0
	for i := range count {
		sz := base
		if i < rem {
			sz++
		}
		if sz == 0 {
			break
		}
		specs = append(specs, chunk{offset: offset, size: sz})
		offset += sz
	}
	return specs
}

// ---------------------------------------------------------------------------
// Transforms  (seq.Seq[T] → seq.Seq[R])
// ---------------------------------------------------------------------------

// Map applies fn to every element in parallel and returns the results in the
// original order. Each chunk of the input is processed by a separate goroutine.
func Map[T, R any](s seq.Seq[T], fn func(T) R, opts ...Option) seq.Seq[R] {
	cfg := resolveConfig(opts)
	items := s.ToSlice()
	n := len(items)
	if n == 0 {
		return seq.Empty[R]()
	}

	chunks := computeChunks(n, cfg.workers)
	result := make([]R, n)

	var wg sync.WaitGroup
	for _, c := range chunks {
		wg.Add(1)
		go func(off, sz int) {
			defer wg.Done()
			for i := range sz {
				result[off+i] = fn(items[off+i])
			}
		}(c.offset, c.size)
	}
	wg.Wait()

	return seq.OfSlice(result)
}

// Filter keeps elements where fn returns true, processing chunks in parallel.
// Output order matches input order.
func Filter[T any](s seq.Seq[T], fn func(T) bool, opts ...Option) seq.Seq[T] {
	cfg := resolveConfig(opts)
	items := s.ToSlice()
	n := len(items)
	if n == 0 {
		return seq.Empty[T]()
	}

	chunks := computeChunks(n, cfg.workers)
	perChunk := make([][]T, len(chunks))

	var wg sync.WaitGroup
	for ci, c := range chunks {
		wg.Add(1)
		go func(idx, off, sz int) {
			defer wg.Done()
			var local []T
			for i := range sz {
				v := items[off+i]
				if fn(v) {
					local = append(local, v)
				}
			}
			perChunk[idx] = local
		}(ci, c.offset, c.size)
	}
	wg.Wait()

	var merged []T
	for _, pc := range perChunk {
		merged = append(merged, pc...)
	}
	return seq.OfSlice(merged)
}

// Collect (flatMap) applies fn to every element in parallel, concatenating
// the resulting slices in input order.
func Collect[T, R any](s seq.Seq[T], fn func(T) []R, opts ...Option) seq.Seq[R] {
	cfg := resolveConfig(opts)
	items := s.ToSlice()
	n := len(items)
	if n == 0 {
		return seq.Empty[R]()
	}

	chunks := computeChunks(n, cfg.workers)
	perChunk := make([][]R, len(chunks))

	var wg sync.WaitGroup
	for ci, c := range chunks {
		wg.Add(1)
		go func(idx, off, sz int) {
			defer wg.Done()
			var local []R
			for i := range sz {
				local = append(local, fn(items[off+i])...)
			}
			perChunk[idx] = local
		}(ci, c.offset, c.size)
	}
	wg.Wait()

	var merged []R
	for _, pc := range perChunk {
		merged = append(merged, pc...)
	}
	return seq.OfSlice(merged)
}

// Choose applies fn to every element in parallel, keeping [option.Some] values
// and discarding [option.None]. Output order is preserved.
func Choose[T, R any](s seq.Seq[T], fn func(T) option.Option[R], opts ...Option) seq.Seq[R] {
	cfg := resolveConfig(opts)
	items := s.ToSlice()
	n := len(items)
	if n == 0 {
		return seq.Empty[R]()
	}

	chunks := computeChunks(n, cfg.workers)
	perChunk := make([][]R, len(chunks))

	var wg sync.WaitGroup
	for ci, c := range chunks {
		wg.Add(1)
		go func(idx, off, sz int) {
			defer wg.Done()
			var local []R
			for i := range sz {
				if r := fn(items[off+i]); r.IsSome() {
					local = append(local, r.Unwrap())
				}
			}
			perChunk[idx] = local
		}(ci, c.offset, c.size)
	}
	wg.Wait()

	var merged []R
	for _, pc := range perChunk {
		merged = append(merged, pc...)
	}
	return seq.OfSlice(merged)
}

// ---------------------------------------------------------------------------
// Terminal operations
// ---------------------------------------------------------------------------

// ForEach calls fn for every element in parallel. No ordering guarantee on
// the invocations of fn; use Map if you need ordered results.
func ForEach[T any](s seq.Seq[T], fn func(T), opts ...Option) {
	cfg := resolveConfig(opts)
	items := s.ToSlice()
	n := len(items)
	if n == 0 {
		return
	}

	chunks := computeChunks(n, cfg.workers)
	var wg sync.WaitGroup
	for _, c := range chunks {
		wg.Add(1)
		go func(off, sz int) {
			defer wg.Done()
			for i := range sz {
				fn(items[off+i])
			}
		}(c.offset, c.size)
	}
	wg.Wait()
}

// Reduce combines elements using fn, processing chunks in parallel. Each chunk
// is reduced independently, then the partial results are reduced sequentially.
//
// The reducer fn must be associative: fn(fn(a,b),c) == fn(a,fn(b,c)).
// Commutativity is NOT required; the original element order is respected
// within and across chunks.
//
// Returns (zero, false) for an empty sequence.
func Reduce[T any](s seq.Seq[T], fn func(T, T) T, opts ...Option) (T, bool) {
	cfg := resolveConfig(opts)
	items := s.ToSlice()
	n := len(items)
	if n == 0 {
		var zero T
		return zero, false
	}
	if n == 1 {
		return items[0], true
	}

	chunks := computeChunks(n, cfg.workers)
	partials := make([]T, len(chunks))

	var wg sync.WaitGroup
	for ci, c := range chunks {
		wg.Add(1)
		go func(idx, off, sz int) {
			defer wg.Done()
			acc := items[off]
			for i := 1; i < sz; i++ {
				acc = fn(acc, items[off+i])
			}
			partials[idx] = acc
		}(ci, c.offset, c.size)
	}
	wg.Wait()

	// sequential merge of partial results
	acc := partials[0]
	for i := 1; i < len(partials); i++ {
		acc = fn(acc, partials[i])
	}
	return acc, true
}

// GroupBy groups elements by key, processing chunks in parallel and merging
// local maps. Order within each group is preserved.
func GroupBy[T any, K comparable](s seq.Seq[T], key func(T) K, opts ...Option) map[K][]T {
	cfg := resolveConfig(opts)
	items := s.ToSlice()
	n := len(items)
	if n == 0 {
		return make(map[K][]T)
	}

	chunks := computeChunks(n, cfg.workers)
	localMaps := make([]map[K][]T, len(chunks))

	var wg sync.WaitGroup
	for ci, c := range chunks {
		wg.Add(1)
		go func(idx, off, sz int) {
			defer wg.Done()
			m := make(map[K][]T)
			for i := range sz {
				v := items[off+i]
				k := key(v)
				m[k] = append(m[k], v)
			}
			localMaps[idx] = m
		}(ci, c.offset, c.size)
	}
	wg.Wait()

	// merge in chunk order to preserve intra-group order
	merged := make(map[K][]T)
	for _, m := range localMaps {
		for k, vs := range m {
			merged[k] = append(merged[k], vs...)
		}
	}
	return merged
}

// CountByKey counts occurrences of each key in parallel.
func CountByKey[T any, K comparable](s seq.Seq[T], fn func(T) K, opts ...Option) map[K]int {
	cfg := resolveConfig(opts)
	items := s.ToSlice()
	n := len(items)
	if n == 0 {
		return make(map[K]int)
	}

	chunks := computeChunks(n, cfg.workers)
	localMaps := make([]map[K]int, len(chunks))

	var wg sync.WaitGroup
	for ci, c := range chunks {
		wg.Add(1)
		go func(idx, off, sz int) {
			defer wg.Done()
			m := make(map[K]int)
			for i := range sz {
				m[fn(items[off+i])]++
			}
			localMaps[idx] = m
		}(ci, c.offset, c.size)
	}
	wg.Wait()

	merged := make(map[K]int)
	for _, m := range localMaps {
		for k, c := range m {
			merged[k] += c
		}
	}
	return merged
}

// Partition splits elements into two slices in parallel: the first where fn
// returns true, the second where fn returns false. Order is preserved in each.
func Partition[T any](s seq.Seq[T], fn func(T) bool, opts ...Option) ([]T, []T) {
	cfg := resolveConfig(opts)
	items := s.ToSlice()
	n := len(items)
	if n == 0 {
		return nil, nil
	}

	chunks := computeChunks(n, cfg.workers)
	yesChunks := make([][]T, len(chunks))
	noChunks := make([][]T, len(chunks))

	var wg sync.WaitGroup
	for ci, c := range chunks {
		wg.Add(1)
		go func(idx, off, sz int) {
			defer wg.Done()
			var yes, no []T
			for i := range sz {
				v := items[off+i]
				if fn(v) {
					yes = append(yes, v)
				} else {
					no = append(no, v)
				}
			}
			yesChunks[idx] = yes
			noChunks[idx] = no
		}(ci, c.offset, c.size)
	}
	wg.Wait()

	var allYes, allNo []T
	for i := range chunks {
		allYes = append(allYes, yesChunks[i]...)
		allNo = append(allNo, noChunks[i]...)
	}
	return allYes, allNo
}

// Sum returns the sum of all numeric elements, computed in parallel.
func Sum[T seq.Numeric](s seq.Seq[T], opts ...Option) T {
	cfg := resolveConfig(opts)
	items := s.ToSlice()
	n := len(items)
	if n == 0 {
		var zero T
		return zero
	}

	chunks := computeChunks(n, cfg.workers)
	partials := make([]T, len(chunks))

	var wg sync.WaitGroup
	for ci, c := range chunks {
		wg.Add(1)
		go func(idx, off, sz int) {
			defer wg.Done()
			var sum T
			for i := range sz {
				sum += items[off+i]
			}
			partials[idx] = sum
		}(ci, c.offset, c.size)
	}
	wg.Wait()

	var total T
	for _, p := range partials {
		total += p
	}
	return total
}

// SumBy returns the sum of fn(element) for all elements, computed in parallel.
func SumBy[T any, N seq.Numeric](s seq.Seq[T], fn func(T) N, opts ...Option) N {
	cfg := resolveConfig(opts)
	items := s.ToSlice()
	n := len(items)
	if n == 0 {
		var zero N
		return zero
	}

	chunks := computeChunks(n, cfg.workers)
	partials := make([]N, len(chunks))

	var wg sync.WaitGroup
	for ci, c := range chunks {
		wg.Add(1)
		go func(idx, off, sz int) {
			defer wg.Done()
			var sum N
			for i := range sz {
				sum += fn(items[off+i])
			}
			partials[idx] = sum
		}(ci, c.offset, c.size)
	}
	wg.Wait()

	var total N
	for _, p := range partials {
		total += p
	}
	return total
}

// Exists returns true if at least one element satisfies fn. Goroutines
// short-circuit as soon as a match is found.
func Exists[T any](s seq.Seq[T], fn func(T) bool, opts ...Option) bool {
	cfg := resolveConfig(opts)
	items := s.ToSlice()
	n := len(items)
	if n == 0 {
		return false
	}

	chunks := computeChunks(n, cfg.workers)
	var found atomic.Bool

	var wg sync.WaitGroup
	for _, c := range chunks {
		wg.Add(1)
		go func(off, sz int) {
			defer wg.Done()
			for i := range sz {
				if found.Load() {
					return
				}
				if fn(items[off+i]) {
					found.Store(true)
					return
				}
			}
		}(c.offset, c.size)
	}
	wg.Wait()

	return found.Load()
}

// ForAll returns true if all elements satisfy fn. Goroutines short-circuit as
// soon as a non-matching element is found.
func ForAll[T any](s seq.Seq[T], fn func(T) bool, opts ...Option) bool {
	cfg := resolveConfig(opts)
	items := s.ToSlice()
	n := len(items)
	if n == 0 {
		return true
	}

	chunks := computeChunks(n, cfg.workers)
	var failed atomic.Bool

	var wg sync.WaitGroup
	for _, c := range chunks {
		wg.Add(1)
		go func(off, sz int) {
			defer wg.Done()
			for i := range sz {
				if failed.Load() {
					return
				}
				if !fn(items[off+i]) {
					failed.Store(true)
					return
				}
			}
		}(c.offset, c.size)
	}
	wg.Wait()

	return !failed.Load()
}

