package seq_test

import (
	"cmp"
	"testing"

	"github.com/natalie-o-perret/gof/pipe"
	"github.com/natalie-o-perret/gof/seq"
)

// ── data ──────────────────────────────────────────────────────────────────────

var benchData = func() []int {
	out := make([]int, 10_000)
	for i := range out {
		out[i] = i
	}
	return out
}()

// sink prevents dead-code elimination.
var sink any

// ── small pipeline: filter → map → truncate → toSlice ─────────────────────────

// Baseline: direct method / function calls, no curried wrappers.
func BenchmarkSmall_Direct(b *testing.B) {
	for b.Loop() {
		sink = seq.Map(
			seq.OfSlice(benchData).Filter(func(n int) bool { return n%2 == 0 }),
			func(n int) int { return n * 3 },
		).Truncate(100).ToSlice()
	}
}

// Pipe + curried Fn helpers (each step wrapped in a closure).
func BenchmarkSmall_Pipe(b *testing.B) {
	for b.Loop() {
		sink = pipe.Pipe3(
			seq.OfSlice(benchData),
			seq.FilterFn(func(n int) bool { return n%2 == 0 }),
			seq.MapFn(func(n int) int { return n * 3 }),
			pipe.Compose2(
				seq.TruncateFn[int](100),
				seq.ToSliceFn[int](),
			),
		)
	}
}

// ── medium pipeline: filter → exclude → skip → truncate → map → sort → toSlice

func BenchmarkMedium_Direct(b *testing.B) {
	for b.Loop() {
		sink = seq.Map(
			seq.OfSlice(benchData).
				Filter(func(n int) bool { return n%2 == 0 }).
				Exclude(func(n int) bool { return n%10 == 0 }).
				Skip(10).
				Truncate(200),
			func(n int) int { return n * 2 },
		).SortWith(cmp.Compare).ToSlice()
	}
}

func BenchmarkMedium_Pipe(b *testing.B) {
	for b.Loop() {
		sink = pipe.Pipe2(
			seq.OfSlice(benchData),
			pipe.ComposeN(
				seq.FilterFn(func(n int) bool { return n%2 == 0 }),
				seq.ExcludeFn(func(n int) bool { return n%10 == 0 }),
				seq.SkipFn[int](10),
				seq.TruncateFn[int](200),
			),
			pipe.Compose3(
				seq.MapFn(func(n int) int { return n * 2 }),
				seq.SortWithFn[int](cmp.Compare),
				seq.ToSliceFn[int](),
			),
		)
	}
}

// ── large pipeline: 10 same-type steps ────────────────────────────────────────

func BenchmarkLarge_Direct(b *testing.B) {
	for b.Loop() {
		sink = seq.OfSlice(benchData).
			Filter(func(n int) bool { return n%2 == 0 }).
			Exclude(func(n int) bool { return n%10 == 0 }).
			Skip(5).
			TakeWhile(func(n int) bool { return n < 8000 }).
			Truncate(500).
			Append(seq.OfSlice([]int{99999})).
			Rev().
			Skip(1).
			Truncate(100).
			SortWith(cmp.Compare).
			ToSlice()
	}
}

func BenchmarkLarge_Pipe(b *testing.B) {
	for b.Loop() {
		sink = pipe.Pipe2(
			seq.OfSlice(benchData),
			pipe.ComposeN(
				seq.FilterFn(func(n int) bool { return n%2 == 0 }),
				seq.ExcludeFn(func(n int) bool { return n%10 == 0 }),
				seq.SkipFn[int](5),
				seq.TakeWhileFn(func(n int) bool { return n < 8000 }),
				seq.TruncateFn[int](500),
				seq.AppendFn(seq.OfSlice([]int{99999})),
				seq.RevFn[int](),
				seq.SkipFn[int](1),
				seq.TruncateFn[int](100),
				seq.SortWithFn[int](cmp.Compare),
			),
			seq.ToSliceFn[int](),
		)
	}
}

// ── compose overhead in isolation ─────────────────────────────────────────────

func BenchmarkComposeN_Overhead(b *testing.B) {
	// 5 composed same-type steps vs 5 direct method calls
	b.Run("direct", func(b *testing.B) {
		for b.Loop() {
			sink = seq.OfSlice(benchData).
				Filter(func(n int) bool { return n > 100 }).
				Skip(10).
				Truncate(500).
				Filter(func(n int) bool { return n%3 == 0 }).
				Skip(5).
				ToSlice()
		}
	})
	b.Run("compose", func(b *testing.B) {
		for b.Loop() {
			sink = pipe.Pipe2(
				seq.OfSlice(benchData),
				pipe.ComposeN(
					seq.FilterFn(func(n int) bool { return n > 100 }),
					seq.SkipFn[int](10),
					seq.TruncateFn[int](500),
					seq.FilterFn(func(n int) bool { return n%3 == 0 }),
					seq.SkipFn[int](5),
				),
				seq.ToSliceFn[int](),
			)
		}
	})
}
