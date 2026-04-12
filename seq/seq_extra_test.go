package seq_test

import (
	"fmt"
	"testing"

	"github.com/natalie-o-perret/gof/option"
	"github.com/natalie-o-perret/gof/seq"
)

func assertSlice[T comparable](t *testing.T, got, want []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("length: got %d want %d (%v vs %v)", len(got), len(want), got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("pos %d: got %v want %v", i, got[i], want[i])
		}
	}
}
func TestInit(t *testing.T) {
	assertSlice(t, seq.Init(5, func(i int) int { return i * 2 }).ToSlice(), []int{0, 2, 4, 6, 8})
}
func TestInitInfinite(t *testing.T) {
	assertSlice(t, seq.InitInfinite(func(i int) int { return i * i }).Truncate(4).ToSlice(), []int{0, 1, 4, 9})
}
func TestSingleton(t *testing.T) { assertSlice(t, seq.Singleton(42).ToSlice(), []int{42}) }
func TestTail(t *testing.T) {
	assertSlice(t, seq.OfSlice([]int{1, 2, 3}).Tail().ToSlice(), []int{2, 3})
}
func TestInsertAt(t *testing.T) {
	assertSlice(t, seq.OfSlice([]int{1, 3, 4}).InsertAt(1, 2).ToSlice(), []int{1, 2, 3, 4})
	assertSlice(t, seq.OfSlice([]int{1, 2}).InsertAt(2, 3).ToSlice(), []int{1, 2, 3})
}
func TestRemoveAt(t *testing.T) {
	assertSlice(t, seq.OfSlice([]int{1, 2, 3}).RemoveAt(1).ToSlice(), []int{1, 3})
}
func TestUpdateAt(t *testing.T) {
	assertSlice(t, seq.OfSlice([]int{1, 9, 3}).UpdateAt(1, 2).ToSlice(), []int{1, 2, 3})
}
func TestPermute(t *testing.T) {
	assertSlice(t, seq.OfSlice([]int{1, 2, 3}).Permute(func(i int) int { return 2 - i }).ToSlice(), []int{3, 2, 1})
}
func TestIteri(t *testing.T) {
	var idx []int
	seq.OfSlice([]string{"a", "b", "c"}).Iteri(func(i int, _ string) { idx = append(idx, i) })
	assertSlice(t, idx, []int{0, 1, 2})
}
func TestIsEmpty(t *testing.T) {
	if !seq.Empty[int]().IsEmpty() {
		t.Fatal("expected empty")
	}
	if seq.OfSlice([]int{1}).IsEmpty() {
		t.Fatal("expected non-empty")
	}
}
func TestItem(t *testing.T) {
	v, ok := seq.OfSlice([]int{10, 20, 30}).Item(1)
	if !ok || v != 20 {
		t.Fatalf("got %d %v", v, ok)
	}
	_, ok = seq.OfSlice([]int{1}).Item(5)
	if ok {
		t.Fatal("expected false")
	}
}
func TestTryFind(t *testing.T) {
	r := seq.OfSlice([]int{1, 2, 3, 4}).TryFind(func(n int) bool { return n > 2 })
	if r.UnwrapOr(0) != 3 {
		t.Fatalf("got %d", r.UnwrapOr(0))
	}
	if seq.OfSlice([]int{1}).TryFind(func(n int) bool { return n > 5 }).IsSome() {
		t.Fatal("expected None")
	}
}
func TestTryFindIndex(t *testing.T) {
	r := seq.OfSlice([]int{10, 20, 30}).TryFindIndex(func(n int) bool { return n == 20 })
	if r.UnwrapOr(-1) != 1 {
		t.Fatalf("got %d", r.UnwrapOr(-1))
	}
}
func TestTryFindBack(t *testing.T) {
	r := seq.OfSlice([]int{1, 2, 3, 2, 1}).TryFindBack(func(n int) bool { return n == 2 })
	if r.UnwrapOr(0) != 2 {
		t.Fatalf("got %d", r.UnwrapOr(0))
	}
}
func TestTryFindIndexBack(t *testing.T) {
	r := seq.OfSlice([]int{1, 2, 3, 2, 1}).TryFindIndexBack(func(n int) bool { return n == 2 })
	if r.UnwrapOr(-1) != 3 {
		t.Fatalf("got %d", r.UnwrapOr(-1))
	}
}
func TestReduceMethod(t *testing.T) {
	v, ok := seq.OfSlice([]int{1, 2, 3, 4}).Reduce(func(a, b int) int { return a + b })
	if !ok || v != 10 {
		t.Fatalf("got %d %v", v, ok)
	}
	_, ok = seq.Empty[int]().Reduce(func(a, b int) int { return a + b })
	if ok {
		t.Fatal("expected false")
	}
}
func TestReduceBackMethod(t *testing.T) {
	v, _ := seq.OfSlice([]int{1, 2, 3, 4}).ReduceBack(func(a, b int) int { return a - b })
	if v != -2 {
		t.Fatalf("got %d", v)
	}
}
func TestExactlyOne(t *testing.T) {
	v, ok := seq.Singleton(42).ExactlyOne()
	if !ok || v != 42 {
		t.Fatalf("got %d %v", v, ok)
	}
	_, ok = seq.OfSlice([]int{1, 2}).ExactlyOne()
	if ok {
		t.Fatal("expected false")
	}
	_, ok = seq.Empty[int]().ExactlyOne()
	if ok {
		t.Fatal("expected false")
	}
}
func TestTryExactlyOne(t *testing.T) {
	if seq.Singleton(42).TryExactlyOne().UnwrapOr(0) != 42 {
		t.Fatal("wrong")
	}
}
func TestMapi(t *testing.T) {
	got := seq.Mapi(seq.OfSlice([]string{"a", "b"}), func(i int, s string) string {
		return fmt.Sprintf("%d:%s", i, s)
	}).ToSlice()
	assertSlice(t, got, []string{"0:a", "1:b"})
}
func TestIndexed(t *testing.T) {
	got := seq.Indexed(seq.OfSlice([]string{"x", "y"})).ToSlice()
	if got[0].First != 0 || got[1].First != 1 || got[1].Second != "y" {
		t.Fatalf("got %v", got)
	}
}
func TestPairwise(t *testing.T) {
	got := seq.Pairwise(seq.OfSlice([]int{1, 2, 3, 4})).ToSlice()
	if len(got) != 3 {
		t.Fatalf("got %d", len(got))
	}
	if got[0].First != 1 || got[0].Second != 2 {
		t.Fatalf("got %v", got[0])
	}
	if got[2].First != 3 || got[2].Second != 4 {
		t.Fatalf("got %v", got[2])
	}
}
func TestWindowed(t *testing.T) {
	got := seq.Windowed(seq.OfSlice([]int{1, 2, 3, 4, 5}), 3).ToSlice()
	if len(got) != 3 {
		t.Fatalf("got %d", len(got))
	}
	assertSlice(t, got[0], []int{1, 2, 3})
	assertSlice(t, got[1], []int{2, 3, 4})
	assertSlice(t, got[2], []int{3, 4, 5})
}
func TestChunkBySize(t *testing.T) {
	got := seq.ChunkBySize(seq.OfSlice([]int{1, 2, 3, 4, 5}), 2).ToSlice()
	if len(got) != 3 {
		t.Fatalf("got %d", len(got))
	}
	assertSlice(t, got[0], []int{1, 2})
	assertSlice(t, got[1], []int{3, 4})
	assertSlice(t, got[2], []int{5})
}
func TestSplitInto(t *testing.T) {
	got := seq.SplitInto(seq.OfSlice([]int{1, 2, 3, 4, 5}), 3).ToSlice()
	if len(got) != 3 {
		t.Fatalf("got %d", len(got))
	}
	assertSlice(t, got[0], []int{1, 2})
	assertSlice(t, got[1], []int{3, 4})
	assertSlice(t, got[2], []int{5})
}
func TestScan(t *testing.T) {
	assertSlice(t, seq.Scan(seq.OfSlice([]int{1, 2, 3}), 0, func(a, v int) int { return a + v }).ToSlice(), []int{0, 1, 3, 6})
}
func TestScanBack(t *testing.T) {
	assertSlice(t, seq.ScanBack(seq.OfSlice([]int{1, 2, 3}), 0, func(v, a int) int { return v + a }).ToSlice(), []int{6, 5, 3, 0})
}
func TestTryPick(t *testing.T) {
	r := seq.TryPick(seq.OfSlice([]int{1, 2, 3}), func(n int) option.Option[string] {
		if n == 2 {
			return option.Some("found")
		}
		return option.None[string]()
	})
	if r.UnwrapOr("") != "found" {
		t.Fatalf("got %s", r.UnwrapOr(""))
	}
}
func TestContains(t *testing.T) {
	if !seq.Contains(seq.OfSlice([]int{1, 2, 3}), 2) {
		t.Fatal("expected true")
	}
	if seq.Contains(seq.OfSlice([]int{1, 2, 3}), 5) {
		t.Fatal("expected false")
	}
}
func TestExcept(t *testing.T) {
	assertSlice(t, seq.Except(seq.OfSlice([]int{1, 2, 3, 4, 5}), seq.OfSlice([]int{2, 4})).ToSlice(), []int{1, 3, 5})
}
func TestSum(t *testing.T) {
	if seq.Sum(seq.OfSlice([]int{1, 2, 3})) != 6 {
		t.Fatal("wrong")
	}
}
func TestSumByExtra(t *testing.T) {
	type item struct{ v int }
	if seq.SumBy(seq.OfSlice([]item{{1}, {2}, {3}}), func(i item) int { return i.v }) != 6 {
		t.Fatal("wrong")
	}
}
func TestMinMax(t *testing.T) {
	v, ok := seq.Min(seq.OfSlice([]int{3, 1, 2}))
	if !ok || v != 1 {
		t.Fatalf("min: %d", v)
	}
	v, ok = seq.Max(seq.OfSlice([]int{3, 1, 2}))
	if !ok || v != 3 {
		t.Fatalf("max: %d", v)
	}
}
func TestMinByMaxBy(t *testing.T) {
	type item struct {
		n string
		v int
	}
	items := []item{{"a", 3}, {"b", 1}, {"c", 2}}
	v, _ := seq.MinBy(seq.OfSlice(items), func(i item) int { return i.v })
	if v.n != "b" {
		t.Fatalf("minBy: %v", v)
	}
	v, _ = seq.MaxBy(seq.OfSlice(items), func(i item) int { return i.v })
	if v.n != "a" {
		t.Fatalf("maxBy: %v", v)
	}
}
func TestAverage(t *testing.T) {
	avg, ok := seq.Average(seq.OfSlice([]int{2, 4, 6}))
	if !ok || avg != 4.0 {
		t.Fatalf("got %f", avg)
	}
	_, ok = seq.Average(seq.Empty[int]())
	if ok {
		t.Fatal("expected false")
	}
}
func TestZip3(t *testing.T) {
	got := seq.Zip3(seq.OfSlice([]int{1, 2}), seq.OfSlice([]string{"a", "b"}), seq.OfSlice([]float64{1.0, 2.0})).ToSlice()
	if len(got) != 2 {
		t.Fatalf("got %d", len(got))
	}
	if got[0].First != 1 || got[0].Second != "a" || got[0].Third != 1.0 {
		t.Fatalf("got %v", got[0])
	}
}
func TestUnzip(t *testing.T) {
	ints, strs := seq.Unzip(seq.Zip(seq.OfSlice([]int{1, 2}), seq.OfSlice([]string{"a", "b"})))
	assertSlice(t, ints, []int{1, 2})
	assertSlice(t, strs, []string{"a", "b"})
}
func TestConcat(t *testing.T) {
	got := seq.Concat(seq.OfSlice([]seq.Seq[int]{seq.OfSlice([]int{1, 2}), seq.OfSlice([]int{3})})).ToSlice()
	assertSlice(t, got, []int{1, 2, 3})
}
func TestAllPairs(t *testing.T) {
	got := seq.AllPairs(seq.OfSlice([]int{1, 2}), seq.OfSlice([]string{"a", "b"})).ToSlice()
	if len(got) != 4 {
		t.Fatalf("got %d", len(got))
	}
	if got[0].First != 1 || got[0].Second != "a" {
		t.Fatalf("got %v", got[0])
	}
	if got[3].First != 2 || got[3].Second != "b" {
		t.Fatalf("got %v", got[3])
	}
}
func TestMap2(t *testing.T) {
	assertSlice(t, seq.Map2(seq.OfSlice([]int{1, 2}), seq.OfSlice([]int{10, 20}), func(a, b int) int { return a + b }).ToSlice(), []int{11, 22})
}
func TestIter2(t *testing.T) {
	var sums []int
	seq.Iter2(seq.OfSlice([]int{1, 2}), seq.OfSlice([]int{10, 20}), func(a, b int) { sums = append(sums, a+b) })
	assertSlice(t, sums, []int{11, 22})
}
func TestExists2(t *testing.T) {
	if !seq.Exists2(seq.OfSlice([]int{1, 2, 3}), seq.OfSlice([]int{3, 2, 1}), func(a, b int) bool { return a == b }) {
		t.Fatal("expected true")
	}
	if seq.Exists2(seq.OfSlice([]int{1, 2}), seq.OfSlice([]int{3, 4}), func(a, b int) bool { return a == b }) {
		t.Fatal("expected false")
	}
}
func TestForAll2(t *testing.T) {
	if !seq.ForAll2(seq.OfSlice([]int{1, 2}), seq.OfSlice([]int{2, 3}), func(a, b int) bool { return a < b }) {
		t.Fatal("expected true")
	}
	if seq.ForAll2(seq.OfSlice([]int{1, 5}), seq.OfSlice([]int{2, 3}), func(a, b int) bool { return a < b }) {
		t.Fatal("expected false")
	}
}
func TestCompareWith(t *testing.T) {
	c := func(a, b int) int { return a - b }
	if seq.CompareWith(seq.OfSlice([]int{1, 2}), seq.OfSlice([]int{1, 3}), c) >= 0 {
		t.Fatal("expected <")
	}
	if seq.CompareWith(seq.OfSlice([]int{1, 2}), seq.OfSlice([]int{1, 2}), c) != 0 {
		t.Fatal("expected ==")
	}
	if seq.CompareWith(seq.OfSlice([]int{1, 2, 3}), seq.OfSlice([]int{1, 2}), c) <= 0 {
		t.Fatal("expected >")
	}
}
func TestFoldBackExtra(t *testing.T) {
	got := seq.FoldBack(seq.OfSlice([]int{1, 2, 3}), 0, func(v, acc int) int { return v - acc })
	if got != 2 {
		t.Fatalf("got %d", got)
	}
}
func TestTranspose(t *testing.T) {
	rows := seq.OfSlice([]seq.Seq[int]{seq.OfSlice([]int{1, 2}), seq.OfSlice([]int{3, 4})})
	cols := seq.Transpose(rows).ToSlice()
	if len(cols) != 2 {
		t.Fatalf("got %d", len(cols))
	}
	assertSlice(t, cols[0].ToSlice(), []int{1, 3})
	assertSlice(t, cols[1].ToSlice(), []int{2, 4})
}
func TestMapFold(t *testing.T) {
	results, state := seq.MapFold(seq.OfSlice([]int{1, 2, 3}), 0, func(s, v int) (int, int) {
		ns := s + v
		return ns, ns
	})
	assertSlice(t, results, []int{1, 3, 6})
	if state != 6 {
		t.Fatalf("state: %d", state)
	}
}
func TestUnfold(t *testing.T) {
	// first 6 fibonacci numbers
	fibs := seq.Unfold([2]int{0, 1}, func(s [2]int) option.Option[seq.Pair[int, [2]int]] {
		if s[0] > 8 {
			return option.None[seq.Pair[int, [2]int]]()
		}
		return option.Some(seq.Pair[int, [2]int]{First: s[0], Second: [2]int{s[1], s[0] + s[1]}})
	}).ToSlice()
	want := []int{0, 1, 1, 2, 3, 5, 8}
	if len(fibs) != len(want) {
		t.Fatalf("got %v", fibs)
	}
	for i, v := range fibs {
		if v != want[i] {
			t.Fatalf("fibs[%d] = %d, want %d", i, v, want[i])
		}
	}
}
func TestRangeStep(t *testing.T) {
	got := seq.RangeStep(0, 10, 2).ToSlice()
	want := []int{0, 2, 4, 6, 8}
	if len(got) != len(want) {
		t.Fatalf("got %v", got)
	}
	for i, v := range got {
		if v != want[i] {
			t.Fatalf("[%d] got %d want %d", i, v, want[i])
		}
	}
	// descending
	desc := seq.RangeStep(5, 0, -1).ToSlice()
	if len(desc) != 5 || desc[0] != 5 || desc[4] != 1 {
		t.Fatalf("descending got %v", desc)
	}
}
func TestOfMap(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	counts := seq.CountByKey(seq.OfMap(m), func(p seq.Pair[string, int]) string { return p.First })
	if counts["a"] != 1 || counts["b"] != 1 {
		t.Fatalf("got %v", counts)
	}
}
func TestOfOption(t *testing.T) {
	got := seq.OfOption(option.Some(42)).ToSlice()
	if len(got) != 1 || got[0] != 42 {
		t.Fatalf("got %v", got)
	}
	if len(seq.OfOption(option.None[int]()).ToSlice()) != 0 {
		t.Fatal("expected empty")
	}
}
func TestPartition(t *testing.T) {
	even, odd := seq.Partition(seq.Range(1, 7), func(n int) bool { return n%2 == 0 })
	if len(even) != 3 || even[0] != 2 {
		t.Fatalf("even: %v", even)
	}
	if len(odd) != 3 || odd[0] != 1 {
		t.Fatalf("odd: %v", odd)
	}
}
func TestCountByKey(t *testing.T) {
	words := seq.OfSlice([]string{"a", "b", "a", "c", "b", "a"})
	counts := seq.CountByKey(words, func(s string) string { return s })
	if counts["a"] != 3 || counts["b"] != 2 || counts["c"] != 1 {
		t.Fatalf("got %v", counts)
	}
}
func TestInterleave(t *testing.T) {
	a := seq.OfSlice([]int{1, 3, 5})
	b := seq.OfSlice([]int{2, 4, 6})
	got := seq.Interleave(a, b).ToSlice()
	want := []int{1, 2, 3, 4, 5, 6}
	if len(got) != len(want) {
		t.Fatalf("got %v", got)
	}
	for i, v := range got {
		if v != want[i] {
			t.Fatalf("[%d] got %d want %d", i, v, want[i])
		}
	}
}
