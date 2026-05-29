package pseq_test

import (
	"math"
	"testing"

	"github.com/natalie-o-perret/go-functionalish/pseq"
	"github.com/natalie-o-perret/go-functionalish/seq"
)

// sink prevents dead-code elimination in benchmarks.
var benchSink any

// heavyFn simulates a CPU-bound transform (sqrt loop) to show the
// crossover point where parallel execution outperforms sequential.
func heavyFn(n int) float64 {
	x := float64(n)
	for range 500 {
		x = math.Sqrt(x + 1)
	}
	return x
}

// ---------------------------------------------------------------------------
// Map benchmarks
// ---------------------------------------------------------------------------

func BenchmarkMap_Seq_10k(b *testing.B) {
	data := ints(10_000)
	for b.Loop() {
		benchSink = seq.OfSlice(data).Map(heavyFn).ToSlice()
	}
}

func BenchmarkMap_PSeq_10k(b *testing.B) {
	data := ints(10_000)
	for b.Loop() {
		benchSink = pseq.Map(seq.OfSlice(data), heavyFn).ToSlice()
	}
}

func BenchmarkMap_PSeq_10k_Workers4(b *testing.B) {
	data := ints(10_000)
	for b.Loop() {
		benchSink = pseq.Map(seq.OfSlice(data), heavyFn, pseq.WithWorkers(4)).ToSlice()
	}
}

func BenchmarkMap_Seq_100k(b *testing.B) {
	data := ints(100_000)
	for b.Loop() {
		benchSink = seq.OfSlice(data).Map(heavyFn).ToSlice()
	}
}

func BenchmarkMap_PSeq_100k(b *testing.B) {
	data := ints(100_000)
	for b.Loop() {
		benchSink = pseq.Map(seq.OfSlice(data), heavyFn).ToSlice()
	}
}

// ---------------------------------------------------------------------------
// Filter benchmarks
// ---------------------------------------------------------------------------

func BenchmarkFilter_Seq_10k(b *testing.B) {
	data := ints(10_000)
	pred := func(n int) bool { return heavyFn(n) > 1.0 }
	for b.Loop() {
		benchSink = seq.OfSlice(data).Filter(pred).ToSlice()
	}
}

func BenchmarkFilter_PSeq_10k(b *testing.B) {
	data := ints(10_000)
	pred := func(n int) bool { return heavyFn(n) > 1.0 }
	for b.Loop() {
		benchSink = pseq.Filter(seq.OfSlice(data), pred).ToSlice()
	}
}

// ---------------------------------------------------------------------------
// Reduce benchmarks
// ---------------------------------------------------------------------------

func BenchmarkReduce_Seq_10k(b *testing.B) {
	data := ints(10_000)
	for b.Loop() {
		benchSink, _ = seq.OfSlice(data).Reduce(func(a, b int) int { return a + b })
	}
}

func BenchmarkReduce_PSeq_10k(b *testing.B) {
	data := ints(10_000)
	for b.Loop() {
		benchSink, _ = pseq.Reduce(seq.OfSlice(data), func(a, b int) int { return a + b })
	}
}

// ---------------------------------------------------------------------------
// Lightweight fn (shows overhead of parallelism for trivial work)
// ---------------------------------------------------------------------------

func BenchmarkMap_Lightweight_Seq_10k(b *testing.B) {
	data := ints(10_000)
	for b.Loop() {
		benchSink = seq.OfSlice(data).Map(func(n int) int { return n * 2 }).ToSlice()
	}
}

func BenchmarkMap_Lightweight_PSeq_10k(b *testing.B) {
	data := ints(10_000)
	for b.Loop() {
		benchSink = pseq.Map(seq.OfSlice(data), func(n int) int { return n * 2 }).ToSlice()
	}
}
