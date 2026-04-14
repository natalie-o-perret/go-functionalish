package pseq_test

import (
	"math"
	"testing"

	lop "github.com/samber/lo/parallel"

	"github.com/natalie-o-perret/gof/pseq"
	"github.com/natalie-o-perret/gof/seq"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

var vsSink any

func makeInts(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = i
	}
	return out
}

// CPU-heavy: ~500 sqrt iterations per element
func cpuHeavy(n int) float64 {
	x := float64(n)
	for range 500 {
		x = math.Sqrt(x + 1)
	}
	return x
}

// Lightweight: trivial arithmetic
func lightweight(n int) int { return n*3 + 1 }

// ---------------------------------------------------------------------------
// Map — CPU-heavy workload
// ---------------------------------------------------------------------------

func BenchmarkVs_Map_Heavy_1k_Seq(b *testing.B) {
	data := makeInts(1_000)
	for b.Loop() {
		vsSink = seq.Map(seq.OfSlice(data), cpuHeavy).ToSlice()
	}
}

func BenchmarkVs_Map_Heavy_1k_PSeq(b *testing.B) {
	data := makeInts(1_000)
	for b.Loop() {
		vsSink = pseq.Map(seq.OfSlice(data), cpuHeavy).ToSlice()
	}
}

func BenchmarkVs_Map_Heavy_1k_Lo(b *testing.B) {
	data := makeInts(1_000)
	for b.Loop() {
		vsSink = lop.Map(data, func(n int, _ int) float64 { return cpuHeavy(n) })
	}
}

func BenchmarkVs_Map_Heavy_10k_Seq(b *testing.B) {
	data := makeInts(10_000)
	for b.Loop() {
		vsSink = seq.Map(seq.OfSlice(data), cpuHeavy).ToSlice()
	}
}

func BenchmarkVs_Map_Heavy_10k_PSeq(b *testing.B) {
	data := makeInts(10_000)
	for b.Loop() {
		vsSink = pseq.Map(seq.OfSlice(data), cpuHeavy).ToSlice()
	}
}

func BenchmarkVs_Map_Heavy_10k_Lo(b *testing.B) {
	data := makeInts(10_000)
	for b.Loop() {
		vsSink = lop.Map(data, func(n int, _ int) float64 { return cpuHeavy(n) })
	}
}

func BenchmarkVs_Map_Heavy_100k_Seq(b *testing.B) {
	data := makeInts(100_000)
	for b.Loop() {
		vsSink = seq.Map(seq.OfSlice(data), cpuHeavy).ToSlice()
	}
}

func BenchmarkVs_Map_Heavy_100k_PSeq(b *testing.B) {
	data := makeInts(100_000)
	for b.Loop() {
		vsSink = pseq.Map(seq.OfSlice(data), cpuHeavy).ToSlice()
	}
}

func BenchmarkVs_Map_Heavy_100k_Lo(b *testing.B) {
	data := makeInts(100_000)
	for b.Loop() {
		vsSink = lop.Map(data, func(n int, _ int) float64 { return cpuHeavy(n) })
	}
}

// ---------------------------------------------------------------------------
// Map — Lightweight workload (shows goroutine overhead)
// ---------------------------------------------------------------------------

func BenchmarkVs_Map_Light_1k_Seq(b *testing.B) {
	data := makeInts(1_000)
	for b.Loop() {
		vsSink = seq.Map(seq.OfSlice(data), lightweight).ToSlice()
	}
}

func BenchmarkVs_Map_Light_1k_PSeq(b *testing.B) {
	data := makeInts(1_000)
	for b.Loop() {
		vsSink = pseq.Map(seq.OfSlice(data), lightweight).ToSlice()
	}
}

func BenchmarkVs_Map_Light_1k_Lo(b *testing.B) {
	data := makeInts(1_000)
	for b.Loop() {
		vsSink = lop.Map(data, func(n int, _ int) int { return lightweight(n) })
	}
}

func BenchmarkVs_Map_Light_10k_Seq(b *testing.B) {
	data := makeInts(10_000)
	for b.Loop() {
		vsSink = seq.Map(seq.OfSlice(data), lightweight).ToSlice()
	}
}

func BenchmarkVs_Map_Light_10k_PSeq(b *testing.B) {
	data := makeInts(10_000)
	for b.Loop() {
		vsSink = pseq.Map(seq.OfSlice(data), lightweight).ToSlice()
	}
}

func BenchmarkVs_Map_Light_10k_Lo(b *testing.B) {
	data := makeInts(10_000)
	for b.Loop() {
		vsSink = lop.Map(data, func(n int, _ int) int { return lightweight(n) })
	}
}

func BenchmarkVs_Map_Light_100k_Seq(b *testing.B) {
	data := makeInts(100_000)
	for b.Loop() {
		vsSink = seq.Map(seq.OfSlice(data), lightweight).ToSlice()
	}
}

func BenchmarkVs_Map_Light_100k_PSeq(b *testing.B) {
	data := makeInts(100_000)
	for b.Loop() {
		vsSink = pseq.Map(seq.OfSlice(data), lightweight).ToSlice()
	}
}

func BenchmarkVs_Map_Light_100k_Lo(b *testing.B) {
	data := makeInts(100_000)
	for b.Loop() {
		vsSink = lop.Map(data, func(n int, _ int) int { return lightweight(n) })
	}
}

// ---------------------------------------------------------------------------
// ForEach — CPU-heavy
// ---------------------------------------------------------------------------

func BenchmarkVs_ForEach_Heavy_10k_PSeq(b *testing.B) {
	data := makeInts(10_000)
	for b.Loop() {
		pseq.ForEach(seq.OfSlice(data), func(n int) { cpuHeavy(n) })
	}
}

func BenchmarkVs_ForEach_Heavy_10k_Lo(b *testing.B) {
	data := makeInts(10_000)
	for b.Loop() {
		lop.ForEach(data, func(n int, _ int) { cpuHeavy(n) })
	}
}

// ---------------------------------------------------------------------------
// GroupBy — CPU-heavy key function
// ---------------------------------------------------------------------------

func BenchmarkVs_GroupBy_Heavy_10k_PSeq(b *testing.B) {
	data := makeInts(10_000)
	for b.Loop() {
		vsSink = pseq.GroupBy(seq.OfSlice(data), func(n int) int {
			return int(cpuHeavy(n)*10) % 10
		})
	}
}

func BenchmarkVs_GroupBy_Heavy_10k_Lo(b *testing.B) {
	data := makeInts(10_000)
	for b.Loop() {
		vsSink = lop.GroupBy(data, func(n int) int {
			return int(cpuHeavy(n)*10) % 10
		})
	}
}

