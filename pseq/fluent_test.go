package pseq_test

import (
	"sort"
	"sync/atomic"
	"testing"

	"github.com/natalie-o-perret/go-functionalish/option"
	"github.com/natalie-o-perret/go-functionalish/pseq"
	"github.com/natalie-o-perret/go-functionalish/seq"
)

// ---------------------------------------------------------------------------
// PSeq construction
// ---------------------------------------------------------------------------

func TestPSeq_OfSlice_Empty(t *testing.T) {
	got := pseq.OfSlice([]int{}).ToSlice()
	if len(got) != 0 {
		t.Fatalf("want [], got %v", got)
	}
}

func TestPSeq_OfSlice_ToSlice(t *testing.T) {
	got := pseq.OfSlice([]int{1, 2, 3}).ToSlice()
	assertSliceEqual(t, got, []int{1, 2, 3})
}

func TestPSeq_Of_FromSeq(t *testing.T) {
	got := pseq.Of(seq.OfSlice([]string{"a", "b"})).ToSlice()
	assertSliceEqual(t, got, []string{"a", "b"})
}

func TestPSeq_Seq_RoundTrip(t *testing.T) {
	p := pseq.OfSlice([]int{10, 20, 30})
	got := p.Seq().ToSlice()
	assertSliceEqual(t, got, []int{10, 20, 30})
}

// ---------------------------------------------------------------------------
// Filter
// ---------------------------------------------------------------------------

func TestPSeq_Filter_Empty(t *testing.T) {
	got := pseq.OfSlice([]int{}).Filter(func(n int) bool { return true }).ToSlice()
	if len(got) != 0 {
		t.Fatalf("want [], got %v", got)
	}
}

func TestPSeq_Filter_KeepsEvens(t *testing.T) {
	got := pseq.OfSlice([]int{1, 2, 3, 4, 5, 6}).
		Filter(func(n int) bool { return n%2 == 0 }).
		ToSlice()
	assertSliceEqual(t, got, []int{2, 4, 6})
}

func TestPSeq_Filter_OrderPreserved(t *testing.T) {
	input := ints(100)
	got := pseq.OfSlice(input).
		Filter(func(n int) bool { return n%3 == 0 }).
		ToSlice()
	want := seq.OfSlice(input).Filter(func(n int) bool { return n%3 == 0 }).ToSlice()
	assertSliceEqual(t, got, want)
}

// ---------------------------------------------------------------------------
// Map (type-changing generic method)
// ---------------------------------------------------------------------------

func TestPSeq_Map_Double(t *testing.T) {
	got := pseq.OfSlice([]int{1, 2, 3}).
		Map(func(n int) int { return n * 2 }).
		ToSlice()
	assertSliceEqual(t, got, []int{2, 4, 6})
}

func TestPSeq_Map_TypeChange(t *testing.T) {
	got := pseq.OfSlice([]int{1, 2, 3}).
		Map(func(n int) string {
			if n%2 == 0 {
				return "even"
			}
			return "odd"
		}).
		ToSlice()
	assertSliceEqual(t, got, []string{"odd", "even", "odd"})
}

func TestPSeq_Map_OrderPreserved(t *testing.T) {
	input := ints(200)
	got := pseq.OfSlice(input).
		Map(func(n int) int { return n * n }).
		ToSlice()
	want := seq.OfSlice(input).Map(func(n int) int { return n * n }).ToSlice()
	assertSliceEqual(t, got, want)
}

// ---------------------------------------------------------------------------
// Collect
// ---------------------------------------------------------------------------

func TestPSeq_Collect_Empty(t *testing.T) {
	got := pseq.OfSlice([]int{}).
		Collect(func(n int) []int { return []int{n, n} }).
		ToSlice()
	if len(got) != 0 {
		t.Fatalf("want [], got %v", got)
	}
}

func TestPSeq_Collect_Expand(t *testing.T) {
	got := pseq.OfSlice([]int{1, 2, 3}).
		Collect(func(n int) []int { return []int{n, -n} }).
		ToSlice()
	assertSliceEqual(t, got, []int{1, -1, 2, -2, 3, -3})
}

// ---------------------------------------------------------------------------
// Choose
// ---------------------------------------------------------------------------

func TestPSeq_Choose_FilterAndMap(t *testing.T) {
	got := pseq.OfSlice([]int{1, 2, 3, 4, 5}).
		Choose(func(n int) option.Option[string] {
			if n%2 == 0 {
				return option.Some("even")
			}
			return option.None[string]()
		}).
		ToSlice()
	assertSliceEqual(t, got, []string{"even", "even"})
}

// ---------------------------------------------------------------------------
// Chained pipeline
// ---------------------------------------------------------------------------

func TestPSeq_Chain_FilterMapCollect(t *testing.T) {
	got := pseq.OfSlice(ints(20)).
		Filter(func(n int) bool { return n%2 == 0 }).
		Map(func(n int) int { return n * 10 }).
		Filter(func(n int) bool { return n > 50 }).
		ToSlice()
	want := []int{60, 80, 100, 120, 140, 160, 180}
	assertSliceEqual(t, got, want)
}

func TestPSeq_Chain_WithWorkers(t *testing.T) {
	// Workers configured once, flow through each step.
	got := pseq.OfSlice(ints(50), pseq.WithWorkers(2)).
		Filter(func(n int) bool { return n%2 == 0 }).
		Map(func(n int) int { return n + 1 }).
		ToSlice()
	want := seq.OfSlice(ints(50)).
		Filter(func(n int) bool { return n%2 == 0 }).
		Map(func(n int) int { return n + 1 }).
		ToSlice()
	assertSliceEqual(t, got, want)
}

// ---------------------------------------------------------------------------
// GroupBy
// ---------------------------------------------------------------------------

func TestPSeq_GroupBy_Empty(t *testing.T) {
	got := pseq.OfSlice([]int{}).GroupBy(func(n int) int { return n % 2 })
	if len(got) != 0 {
		t.Fatalf("want empty map, got %v", got)
	}
}

func TestPSeq_GroupBy_OddEven(t *testing.T) {
	got := pseq.OfSlice([]int{1, 2, 3, 4, 5, 6}).
		GroupBy(func(n int) string {
			if n%2 == 0 {
				return "even"
			}
			return "odd"
		})
	sort.Ints(got["even"])
	sort.Ints(got["odd"])
	assertSliceEqual(t, got["even"], []int{2, 4, 6})
	assertSliceEqual(t, got["odd"], []int{1, 3, 5})
}

// ---------------------------------------------------------------------------
// CountByKey
// ---------------------------------------------------------------------------

func TestPSeq_CountByKey(t *testing.T) {
	got := pseq.OfSlice([]string{"a", "b", "a", "c", "b", "a"}).
		CountByKey(func(s string) string { return s })
	if got["a"] != 3 || got["b"] != 2 || got["c"] != 1 {
		t.Fatalf("unexpected counts: %v", got)
	}
}

// ---------------------------------------------------------------------------
// ForEach
// ---------------------------------------------------------------------------

func TestPSeq_ForEach_AllVisited(t *testing.T) {
	var count atomic.Int64
	pseq.OfSlice(ints(100)).ForEach(func(int) { count.Add(1) })
	if got := count.Load(); got != 100 {
		t.Fatalf("want 100, got %d", got)
	}
}

// ---------------------------------------------------------------------------
// Reduce
// ---------------------------------------------------------------------------

func TestPSeq_Reduce_Empty(t *testing.T) {
	_, ok := pseq.OfSlice([]int{}).Reduce(func(a, b int) int { return a + b })
	if ok {
		t.Fatal("want false for empty")
	}
}

func TestPSeq_Reduce_Sum(t *testing.T) {
	got, ok := pseq.OfSlice([]int{1, 2, 3, 4, 5}).
		Reduce(func(a, b int) int { return a + b })
	if !ok || got != 15 {
		t.Fatalf("want (15, true), got (%d, %v)", got, ok)
	}
}

// ---------------------------------------------------------------------------
// Partition
// ---------------------------------------------------------------------------

func TestPSeq_Partition(t *testing.T) {
	yes, no := pseq.OfSlice([]int{1, 2, 3, 4, 5, 6}).
		Partition(func(n int) bool { return n%2 == 0 })
	sort.Ints(yes)
	sort.Ints(no)
	assertSliceEqual(t, yes, []int{2, 4, 6})
	assertSliceEqual(t, no, []int{1, 3, 5})
}

// ---------------------------------------------------------------------------
// Exists / ForAll
// ---------------------------------------------------------------------------

func TestPSeq_Exists_Found(t *testing.T) {
	if !pseq.OfSlice([]int{1, 2, 3}).Exists(func(n int) bool { return n == 2 }) {
		t.Fatal("want true")
	}
}

func TestPSeq_Exists_NotFound(t *testing.T) {
	if pseq.OfSlice([]int{1, 3, 5}).Exists(func(n int) bool { return n%2 == 0 }) {
		t.Fatal("want false")
	}
}

func TestPSeq_ForAll_AllMatch(t *testing.T) {
	if !pseq.OfSlice([]int{2, 4, 6}).ForAll(func(n int) bool { return n%2 == 0 }) {
		t.Fatal("want true")
	}
}

func TestPSeq_ForAll_OneFails(t *testing.T) {
	if pseq.OfSlice([]int{2, 3, 6}).ForAll(func(n int) bool { return n%2 == 0 }) {
		t.Fatal("want false")
	}
}

// ---------------------------------------------------------------------------
// SumBy (generic method) / SumOf (package-level)
// ---------------------------------------------------------------------------

func TestPSeq_SumBy(t *testing.T) {
	got := pseq.OfSlice([]string{"hello", "world", "!"}).
		SumBy(func(s string) int { return len(s) })
	if got != 11 {
		t.Fatalf("want 11, got %d", got)
	}
}

func TestPSeq_SumOf(t *testing.T) {
	got := pseq.SumOf(pseq.OfSlice([]int{1, 2, 3, 4, 5}))
	if got != 15 {
		t.Fatalf("want 15, got %d", got)
	}
}
