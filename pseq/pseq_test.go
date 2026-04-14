package pseq_test

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/natalie-o-perret/gof/option"
	"github.com/natalie-o-perret/gof/pipe"
	"github.com/natalie-o-perret/gof/pseq"
	"github.com/natalie-o-perret/gof/seq"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func ints(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = i
	}
	return out
}

func assertSliceEqual[T comparable](t *testing.T, got, want []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d, want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("index %d: got %v, want %v", i, got[i], want[i])
		}
	}
}

func assertMapEqual[K comparable, V comparable](t *testing.T, got, want map[K]V) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("map length mismatch: got %d, want %d", len(got), len(want))
	}
	for k, wv := range want {
		gv, ok := got[k]
		if !ok {
			t.Fatalf("missing key %v", k)
		}
		if gv != wv {
			t.Fatalf("key %v: got %v, want %v", k, gv, wv)
		}
	}
}

// ---------------------------------------------------------------------------
// Map
// ---------------------------------------------------------------------------

func TestMap_Empty(t *testing.T) {
	got := pseq.Map(seq.Empty[int](), func(n int) int { return n * 2 }).ToSlice()
	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

func TestMap_Single(t *testing.T) {
	got := pseq.Map(seq.OfSlice([]int{5}), func(n int) int { return n * 3 }).ToSlice()
	assertSliceEqual(t, got, []int{15})
}

func TestMap_OrderPreserved(t *testing.T) {
	data := ints(1000)
	got := pseq.Map(seq.OfSlice(data), func(n int) int { return n * 2 }).ToSlice()
	want := make([]int, 1000)
	for i := range want {
		want[i] = i * 2
	}
	assertSliceEqual(t, got, want)
}

func TestMap_TypeChanging(t *testing.T) {
	got := pseq.Map(seq.OfSlice([]int{1, 2, 3}), func(n int) string {
		return fmt.Sprintf("n=%d", n)
	}).ToSlice()
	assertSliceEqual(t, got, []string{"n=1", "n=2", "n=3"})
}

func TestMap_WithWorkers1(t *testing.T) {
	data := ints(50)
	got := pseq.Map(seq.OfSlice(data), func(n int) int { return n + 1 }, pseq.WithWorkers(1)).ToSlice()
	want := make([]int, 50)
	for i := range want {
		want[i] = i + 1
	}
	assertSliceEqual(t, got, want)
}

func TestMap_WithWorkers16(t *testing.T) {
	data := ints(100)
	got := pseq.Map(seq.OfSlice(data), func(n int) int { return n * n }, pseq.WithWorkers(16)).ToSlice()
	want := make([]int, 100)
	for i := range want {
		want[i] = i * i
	}
	assertSliceEqual(t, got, want)
}

func TestMap_MoreWorkersThanElements(t *testing.T) {
	got := pseq.Map(seq.OfSlice([]int{1, 2, 3}), func(n int) int { return n * 10 }, pseq.WithWorkers(100)).ToSlice()
	assertSliceEqual(t, got, []int{10, 20, 30})
}

func TestMap_MatchesSeqMap(t *testing.T) {
	data := ints(500)
	fn := func(n int) int { return n*3 + 1 }
	want := seq.Map(seq.OfSlice(data), fn).ToSlice()
	got := pseq.Map(seq.OfSlice(data), fn, pseq.WithWorkers(4)).ToSlice()
	assertSliceEqual(t, got, want)
}

// ---------------------------------------------------------------------------
// Filter
// ---------------------------------------------------------------------------

func TestFilter_Empty(t *testing.T) {
	got := pseq.Filter(seq.Empty[int](), func(n int) bool { return true }).ToSlice()
	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

func TestFilter_AllMatch(t *testing.T) {
	data := ints(10)
	got := pseq.Filter(seq.OfSlice(data), func(n int) bool { return true }).ToSlice()
	assertSliceEqual(t, got, data)
}

func TestFilter_NoneMatch(t *testing.T) {
	got := pseq.Filter(seq.OfSlice(ints(10)), func(n int) bool { return false }).ToSlice()
	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

func TestFilter_OrderPreserved(t *testing.T) {
	data := ints(100)
	got := pseq.Filter(seq.OfSlice(data), func(n int) bool { return n%2 == 0 }, pseq.WithWorkers(8)).ToSlice()
	want := seq.OfSlice(data).Filter(func(n int) bool { return n%2 == 0 }).ToSlice()
	assertSliceEqual(t, got, want)
}

func TestFilter_MatchesSeqFilter(t *testing.T) {
	data := ints(500)
	pred := func(n int) bool { return n%3 == 0 }
	want := seq.OfSlice(data).Filter(pred).ToSlice()
	got := pseq.Filter(seq.OfSlice(data), pred, pseq.WithWorkers(4)).ToSlice()
	assertSliceEqual(t, got, want)
}

// ---------------------------------------------------------------------------
// Collect (flatMap)
// ---------------------------------------------------------------------------

func TestCollect_Empty(t *testing.T) {
	got := pseq.Collect(seq.Empty[int](), func(n int) []int { return []int{n, n} }).ToSlice()
	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

func TestCollect_OrderPreserved(t *testing.T) {
	data := []int{1, 2, 3}
	fn := func(n int) []string { return []string{fmt.Sprintf("%d", n), fmt.Sprintf("%d0", n)} }
	got := pseq.Collect(seq.OfSlice(data), fn, pseq.WithWorkers(3)).ToSlice()
	assertSliceEqual(t, got, []string{"1", "10", "2", "20", "3", "30"})
}

func TestCollect_MatchesSeqCollect(t *testing.T) {
	data := ints(100)
	fn := func(n int) []int { return []int{n, n * 10} }
	want := seq.Collect(seq.OfSlice(data), fn).ToSlice()
	got := pseq.Collect(seq.OfSlice(data), fn, pseq.WithWorkers(4)).ToSlice()
	assertSliceEqual(t, got, want)
}

// ---------------------------------------------------------------------------
// Choose
// ---------------------------------------------------------------------------

func TestChoose_Empty(t *testing.T) {
	got := pseq.Choose(seq.Empty[int](), func(n int) option.Option[int] {
		return option.Some(n)
	}).ToSlice()
	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

func TestChoose_OrderPreserved(t *testing.T) {
	data := ints(100)
	fn := func(n int) option.Option[string] {
		if n%2 == 0 {
			return option.Some(fmt.Sprintf("even:%d", n))
		}
		return option.None[string]()
	}
	got := pseq.Choose(seq.OfSlice(data), fn, pseq.WithWorkers(4)).ToSlice()
	want := seq.Choose(seq.OfSlice(data), fn).ToSlice()
	assertSliceEqual(t, got, want)
}

// ---------------------------------------------------------------------------
// ForEach
// ---------------------------------------------------------------------------

func TestForEach_Empty(t *testing.T) {
	called := false
	pseq.ForEach(seq.Empty[int](), func(n int) { called = true })
	if called {
		t.Fatal("fn should not be called for empty seq")
	}
}

func TestForEach_AllElementsVisited(t *testing.T) {
	data := ints(100)
	var count atomic.Int64
	pseq.ForEach(seq.OfSlice(data), func(n int) { count.Add(1) }, pseq.WithWorkers(8))
	if count.Load() != 100 {
		t.Fatalf("expected 100 calls, got %d", count.Load())
	}
}

// ---------------------------------------------------------------------------
// Reduce
// ---------------------------------------------------------------------------

func TestReduce_Empty(t *testing.T) {
	_, ok := pseq.Reduce(seq.Empty[int](), func(a, b int) int { return a + b })
	if ok {
		t.Fatal("expected false for empty seq")
	}
}

func TestReduce_Single(t *testing.T) {
	v, ok := pseq.Reduce(seq.OfSlice([]int{42}), func(a, b int) int { return a + b })
	if !ok || v != 42 {
		t.Fatalf("expected (42, true), got (%d, %v)", v, ok)
	}
}

func TestReduce_Sum(t *testing.T) {
	data := ints(1000)
	got, ok := pseq.Reduce(seq.OfSlice(data), func(a, b int) int { return a + b }, pseq.WithWorkers(4))
	want := 0
	for _, v := range data {
		want += v
	}
	if !ok || got != want {
		t.Fatalf("expected (%d, true), got (%d, %v)", want, got, ok)
	}
}

func TestReduce_MatchesSeqReduce(t *testing.T) {
	data := ints(500)
	fn := func(a, b int) int { return a + b }
	want, _ := seq.OfSlice(data).Reduce(fn)
	got, ok := pseq.Reduce(seq.OfSlice(data), fn, pseq.WithWorkers(4))
	if !ok || got != want {
		t.Fatalf("expected (%d, true), got (%d, %v)", want, got, ok)
	}
}

// ---------------------------------------------------------------------------
// GroupBy
// ---------------------------------------------------------------------------

func TestGroupBy_Empty(t *testing.T) {
	got := pseq.GroupBy(seq.Empty[int](), func(n int) int { return n })
	if len(got) != 0 {
		t.Fatalf("expected empty map, got %v", got)
	}
}

func TestGroupBy_OrderPreservedInGroups(t *testing.T) {
	data := ints(100)
	key := func(n int) int { return n % 3 }

	got := pseq.GroupBy(seq.OfSlice(data), key, pseq.WithWorkers(4))
	want := seq.GroupBy(seq.OfSlice(data), key)

	if len(got) != len(want) {
		t.Fatalf("group count mismatch: got %d, want %d", len(got), len(want))
	}
	for k, wv := range want {
		gv, ok := got[k]
		if !ok {
			t.Fatalf("missing group %d", k)
		}
		assertSliceEqual(t, gv, wv)
	}
}

// ---------------------------------------------------------------------------
// CountByKey
// ---------------------------------------------------------------------------

func TestCountByKey_Empty(t *testing.T) {
	got := pseq.CountByKey(seq.Empty[string](), func(s string) string { return s })
	if len(got) != 0 {
		t.Fatalf("expected empty map, got %v", got)
	}
}

func TestCountByKey_MatchesSeq(t *testing.T) {
	data := []string{"a", "b", "a", "c", "a", "b"}
	key := func(s string) string { return s }
	want := seq.CountByKey(seq.OfSlice(data), key)
	got := pseq.CountByKey(seq.OfSlice(data), key, pseq.WithWorkers(3))
	assertMapEqual(t, got, want)
}

// ---------------------------------------------------------------------------
// Partition
// ---------------------------------------------------------------------------

func TestPartition_Empty(t *testing.T) {
	yes, no := pseq.Partition(seq.Empty[int](), func(n int) bool { return true })
	if len(yes) != 0 || len(no) != 0 {
		t.Fatalf("expected empty, got yes=%v no=%v", yes, no)
	}
}

func TestPartition_OrderPreserved(t *testing.T) {
	data := ints(100)
	pred := func(n int) bool { return n%2 == 0 }
	wantYes, wantNo := seq.Partition(seq.OfSlice(data), pred)
	gotYes, gotNo := pseq.Partition(seq.OfSlice(data), pred, pseq.WithWorkers(4))
	assertSliceEqual(t, gotYes, wantYes)
	assertSliceEqual(t, gotNo, wantNo)
}

// ---------------------------------------------------------------------------
// Sum / SumBy
// ---------------------------------------------------------------------------

func TestSum_Empty(t *testing.T) {
	got := pseq.Sum(seq.Empty[int]())
	if got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestSum_MatchesSeq(t *testing.T) {
	data := ints(1000)
	want := seq.Sum(seq.OfSlice(data))
	got := pseq.Sum(seq.OfSlice(data), pseq.WithWorkers(4))
	if got != want {
		t.Fatalf("expected %d, got %d", want, got)
	}
}

func TestSumBy_MatchesSeq(t *testing.T) {
	type item struct{ val int }
	data := make([]item, 500)
	for i := range data {
		data[i] = item{val: i}
	}
	fn := func(it item) int { return it.val }
	want := seq.SumBy(seq.OfSlice(data), fn)
	got := pseq.SumBy(seq.OfSlice(data), fn, pseq.WithWorkers(4))
	if got != want {
		t.Fatalf("expected %d, got %d", want, got)
	}
}

// ---------------------------------------------------------------------------
// Exists / ForAll
// ---------------------------------------------------------------------------

func TestExists_Empty(t *testing.T) {
	if pseq.Exists(seq.Empty[int](), func(n int) bool { return true }) {
		t.Fatal("expected false for empty seq")
	}
}

func TestExists_Found(t *testing.T) {
	data := ints(10000)
	if !pseq.Exists(seq.OfSlice(data), func(n int) bool { return n == 9999 }, pseq.WithWorkers(4)) {
		t.Fatal("expected true")
	}
}

func TestExists_NotFound(t *testing.T) {
	data := ints(100)
	if pseq.Exists(seq.OfSlice(data), func(n int) bool { return n == 999 }, pseq.WithWorkers(4)) {
		t.Fatal("expected false")
	}
}

func TestForAll_Empty(t *testing.T) {
	if !pseq.ForAll(seq.Empty[int](), func(n int) bool { return false }) {
		t.Fatal("expected true for empty seq")
	}
}

func TestForAll_AllTrue(t *testing.T) {
	data := ints(100)
	if !pseq.ForAll(seq.OfSlice(data), func(n int) bool { return n >= 0 }, pseq.WithWorkers(4)) {
		t.Fatal("expected true")
	}
}

func TestForAll_OneFails(t *testing.T) {
	data := ints(10000)
	if pseq.ForAll(seq.OfSlice(data), func(n int) bool { return n < 5000 }, pseq.WithWorkers(4)) {
		t.Fatal("expected false")
	}
}

// ---------------------------------------------------------------------------
// Fluent / pipe integration
// ---------------------------------------------------------------------------

func TestMapFn_Pipe(t *testing.T) {
	data := ints(50)
	got := pipe.Pipe2(
		seq.OfSlice(data),
		pseq.MapFn(func(n int) int { return n * 2 }, pseq.WithWorkers(4)),
		seq.ToSliceFn[int](),
	)
	want := make([]int, 50)
	for i := range want {
		want[i] = i * 2
	}
	assertSliceEqual(t, got, want)
}

func TestFilterFn_Pipe(t *testing.T) {
	data := ints(50)
	got := pipe.Pipe2(
		seq.OfSlice(data),
		pseq.FilterFn(func(n int) bool { return n%2 == 0 }, pseq.WithWorkers(4)),
		seq.ToSliceFn[int](),
	)
	want := seq.OfSlice(data).Filter(func(n int) bool { return n%2 == 0 }).ToSlice()
	assertSliceEqual(t, got, want)
}

func TestCollectFn_Pipe(t *testing.T) {
	got := pipe.Pipe2(
		seq.OfSlice([]int{1, 2, 3}),
		pseq.CollectFn(func(n int) []string {
			return []string{fmt.Sprintf("%d", n)}
		}),
		seq.ToSliceFn[string](),
	)
	assertSliceEqual(t, got, []string{"1", "2", "3"})
}

func TestGroupByFn_Pipe(t *testing.T) {
	data := []string{"apple", "avocado", "banana", "blueberry", "cherry"}
	got := pipe.Pipe2(
		seq.OfSlice(data),
		pseq.FilterFn(func(s string) bool { return len(s) > 5 }),
		pseq.GroupByFn(func(s string) string { return string(s[0]) }),
	)
	if len(got) != 3 {
		t.Fatalf("expected 3 groups, got %d: %v", len(got), got)
	}
	assertSliceEqual(t, got["a"], []string{"avocado"})
	assertSliceEqual(t, got["b"], []string{"banana", "blueberry"})
	assertSliceEqual(t, got["c"], []string{"cherry"})
}

// ---------------------------------------------------------------------------
// WithWorkers edge cases
// ---------------------------------------------------------------------------

func TestWithWorkers_Zero_DefaultsToGOMAXPROCS(t *testing.T) {
	// WithWorkers(0) should be ignored, falling back to GOMAXPROCS
	data := ints(10)
	got := pseq.Map(seq.OfSlice(data), func(n int) int { return n + 1 }, pseq.WithWorkers(0)).ToSlice()
	want := make([]int, 10)
	for i := range want {
		want[i] = i + 1
	}
	assertSliceEqual(t, got, want)
}

func TestWithWorkers_Negative_DefaultsToGOMAXPROCS(t *testing.T) {
	data := ints(10)
	got := pseq.Map(seq.OfSlice(data), func(n int) int { return n + 1 }, pseq.WithWorkers(-5)).ToSlice()
	want := make([]int, 10)
	for i := range want {
		want[i] = i + 1
	}
	assertSliceEqual(t, got, want)
}

// ---------------------------------------------------------------------------
// CPU-heavy scenario (validates actual parallelism gains under -race)
// ---------------------------------------------------------------------------

func TestMap_CPUHeavy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping CPU-heavy test in short mode")
	}
	data := ints(200)
	heavy := func(n int) float64 {
		x := float64(n)
		for range 1000 {
			x = math.Sqrt(x + 1)
		}
		return x
	}
	got := pseq.Map(seq.OfSlice(data), heavy, pseq.WithWorkers(4)).ToSlice()
	want := seq.Map(seq.OfSlice(data), heavy).ToSlice()
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d, want %d", len(got), len(want))
	}
	for i := range got {
		if math.Abs(got[i]-want[i]) > 1e-10 {
			t.Fatalf("index %d: got %v, want %v", i, got[i], want[i])
		}
	}
}

// ---------------------------------------------------------------------------
// Mixed pipeline: parallel + sequential steps
// ---------------------------------------------------------------------------

func TestMixedPipeline(t *testing.T) {
	data := ints(200)

	// parallel filter → sequential sort → parallel map → collect
	result := pseq.Map(
		pseq.Filter(seq.OfSlice(data), func(n int) bool { return n%2 == 0 }, pseq.WithWorkers(4)).
			SortWith(func(a, b int) int { return b - a }), // descending
		func(n int) string { return fmt.Sprintf("%d", n) },
		pseq.WithWorkers(4),
	).ToSlice()

	// verify only evens and correct count
	wantEvens := seq.OfSlice(data).Filter(func(n int) bool { return n%2 == 0 }).ToSlice()
	if len(result) != len(wantEvens) {
		t.Fatalf("expected %d elements, got %d", len(wantEvens), len(result))
	}
	for _, s := range result {
		if !strings.HasSuffix(s, "0") && !strings.HasSuffix(s, "2") &&
			!strings.HasSuffix(s, "4") && !strings.HasSuffix(s, "6") &&
			!strings.HasSuffix(s, "8") {
			t.Fatalf("expected even number string, got %s", s)
		}
	}
}

// ---------------------------------------------------------------------------
// Reduce with non-trivial associative operation
// ---------------------------------------------------------------------------

func TestReduce_Concat(t *testing.T) {
	// String concatenation is associative but not commutative
	data := []string{"a", "b", "c", "d", "e", "f"}
	got, ok := pseq.Reduce(seq.OfSlice(data), func(a, b string) string { return a + b }, pseq.WithWorkers(3))
	if !ok || got != "abcdef" {
		t.Fatalf("expected (abcdef, true), got (%s, %v)", got, ok)
	}
}

// ---------------------------------------------------------------------------
// Large data set - correctness with many workers
// ---------------------------------------------------------------------------

func TestFilter_LargeDataManyWorkers(t *testing.T) {
	data := ints(10000)
	pred := func(n int) bool { return n%7 == 0 }
	want := seq.OfSlice(data).Filter(pred).ToSlice()
	got := pseq.Filter(seq.OfSlice(data), pred, pseq.WithWorkers(32)).ToSlice()
	assertSliceEqual(t, got, want)
}

func TestGroupBy_LargeData(t *testing.T) {
	data := ints(10000)
	key := func(n int) int { return n % 10 }
	want := seq.GroupBy(seq.OfSlice(data), key)
	got := pseq.GroupBy(seq.OfSlice(data), key, pseq.WithWorkers(8))

	if len(got) != len(want) {
		t.Fatalf("group count: got %d, want %d", len(got), len(want))
	}
	// Sort both sides for comparison since intra-group order must match
	for k := range want {
		sort.Ints(want[k])
		gv, ok := got[k]
		if !ok {
			t.Fatalf("missing group %d", k)
		}
		sort.Ints(gv)
		assertSliceEqual(t, gv, want[k])
	}
}
