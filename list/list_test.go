package list_test

import (
	"cmp"
	"strconv"
	"testing"

	"github.com/natalie-o-perret/go-functionalish/list"
	"github.com/natalie-o-perret/go-functionalish/option"
)

func assertList[T comparable](t *testing.T, got list.List[T], want ...T) {
	t.Helper()
	gs := got.ToSlice()
	if len(gs) != len(want) {
		t.Fatalf("len: got %d (%v), want %d (%v)", len(gs), gs, len(want), want)
	}
	for i := range gs {
		if gs[i] != want[i] {
			t.Fatalf("index %d: got %v, want %v", i, gs[i], want[i])
		}
	}
}

// ---------------------------------------------------------------------------
// Constructors
// ---------------------------------------------------------------------------

func TestOf(t *testing.T) {
	assertList(t, list.Of(1, 2, 3), 1, 2, 3)
}

func TestFromSlice(t *testing.T) {
	src := []int{4, 5, 6}
	l := list.FromSlice(src)
	// mutation of source must not affect list
	src[0] = 99
	assertList(t, l, 4, 5, 6)
}

func TestReplicate(t *testing.T) {
	assertList(t, list.Replicate("x", 3), "x", "x", "x")
	if !list.Replicate("x", 0).IsEmpty() {
		t.Fatal("expected empty for n=0")
	}
}

func TestInit(t *testing.T) {
	assertList(t, list.Init(4, func(i int) int { return i * 2 }), 0, 2, 4, 6)
}

// ---------------------------------------------------------------------------
// Immutability guarantee
// ---------------------------------------------------------------------------

func TestImmutability_ToSlice(t *testing.T) {
	l := list.Of(1, 2, 3)
	s := l.ToSlice()
	s[0] = 99
	assertList(t, l, 1, 2, 3)
}

// ---------------------------------------------------------------------------
// Same-type pipeline
// ---------------------------------------------------------------------------

func TestFilter(t *testing.T) {
	assertList(t, list.Of(1, 2, 3, 4).Filter(func(n int) bool { return n%2 == 0 }), 2, 4)
}

func TestExclude(t *testing.T) {
	assertList(t, list.Of(1, 2, 3, 4).Exclude(func(n int) bool { return n%2 == 0 }), 1, 3)
}

func TestAppend(t *testing.T) {
	assertList(t, list.Of(1, 2).Append(3, 4), 1, 2, 3, 4)
}

func TestConcat(t *testing.T) {
	assertList(t, list.Of(1, 2).Concat(list.Of(3, 4), list.Of(5)), 1, 2, 3, 4, 5)
}

func TestTruncate(t *testing.T) {
	assertList(t, list.Of(1, 2, 3, 4).Truncate(2), 1, 2)
	assertList(t, list.Of(1, 2).Truncate(10), 1, 2)
	if !list.Of(1, 2).Truncate(0).IsEmpty() {
		t.Fatal("truncate(0) should be empty")
	}
}

func TestSkip(t *testing.T) {
	assertList(t, list.Of(1, 2, 3, 4).Skip(2), 3, 4)
	if !list.Of(1, 2).Skip(10).IsEmpty() {
		t.Fatal("skip past end should be empty")
	}
}

func TestTakeWhile(t *testing.T) {
	assertList(t, list.Of(1, 2, 3, 4).TakeWhile(func(n int) bool { return n < 3 }), 1, 2)
}

func TestSkipWhile(t *testing.T) {
	assertList(t, list.Of(1, 2, 3, 4).SkipWhile(func(n int) bool { return n < 3 }), 3, 4)
}

func TestRev(t *testing.T) {
	assertList(t, list.Of(1, 2, 3).Rev(), 3, 2, 1)
}

func TestTail(t *testing.T) {
	assertList(t, list.Of(1, 2, 3).Tail(), 2, 3)
}

func TestSortWith(t *testing.T) {
	assertList(t, list.Of(3, 1, 2).SortWith(cmp.Compare), 1, 2, 3)
}

func TestInsertAt(t *testing.T) {
	assertList(t, list.Of(1, 3).InsertAt(1, 2), 1, 2, 3)
}

func TestRemoveAt(t *testing.T) {
	assertList(t, list.Of(1, 2, 3).RemoveAt(1), 1, 3)
}

func TestUpdateAt(t *testing.T) {
	assertList(t, list.Of(1, 2, 3).UpdateAt(1, 99), 1, 99, 3)
}

// ---------------------------------------------------------------------------
// Query / terminal
// ---------------------------------------------------------------------------

func TestAt(t *testing.T) {
	l := list.Of(10, 20, 30)
	if l.At(1).UnwrapOr(0) != 20 {
		t.Fatal("expected 20")
	}
	if l.At(5).IsSome() {
		t.Fatal("expected None for out-of-bounds")
	}
}

func TestHead(t *testing.T) {
	v, ok := list.Of(1, 2, 3).Head()
	if !ok || v != 1 {
		t.Fatalf("got (%v,%v)", v, ok)
	}
	_, ok = list.Of[int]().Head()
	if ok {
		t.Fatal("expected false for empty")
	}
}

func TestLast(t *testing.T) {
	v, ok := list.Of(1, 2, 3).Last()
	if !ok || v != 3 {
		t.Fatalf("got (%v,%v)", v, ok)
	}
}

func TestAll(t *testing.T) {
	l := list.Of(1, 2, 3)
	sum := 0
	for v := range l.All() {
		sum += v
	}
	if sum != 6 {
		t.Fatalf("expected 6, got %d", sum)
	}
}

func TestCountBy(t *testing.T) {
	if list.Of(1, 2, 3, 4).CountBy(func(n int) bool { return n > 2 }) != 2 {
		t.Fatal("expected 2")
	}
}

func TestExists(t *testing.T) {
	if !list.Of(1, 2, 3).Exists(func(n int) bool { return n == 2 }) {
		t.Fatal("expected true")
	}
}

func TestForAll(t *testing.T) {
	if !list.Of(2, 4, 6).ForAll(func(n int) bool { return n%2 == 0 }) {
		t.Fatal("expected true")
	}
}

func TestTryFind(t *testing.T) {
	r := list.Of(1, 2, 3).TryFind(func(n int) bool { return n > 1 })
	if r.UnwrapOr(0) != 2 {
		t.Fatal("expected 2")
	}
}

func TestReduce(t *testing.T) {
	v, ok := list.Of(1, 2, 3, 4).Reduce(func(a, b int) int { return a + b })
	if !ok || v != 10 {
		t.Fatalf("expected (10, true), got (%d, %v)", v, ok)
	}
}

// ---------------------------------------------------------------------------
// Generic methods
// ---------------------------------------------------------------------------

func TestMap(t *testing.T) {
	assertList(t, list.Of(1, 2, 3).Map(func(n int) string { return strconv.Itoa(n) }), "1", "2", "3")
}

func TestMapi(t *testing.T) {
	assertList(t, list.Of(10, 20, 30).Mapi(func(i, n int) int { return i + n }), 10, 21, 32)
}

func TestCollect(t *testing.T) {
	assertList(t, list.Of(1, 2, 3).Collect(func(n int) []int { return []int{n, n * 10} }),
		1, 10, 2, 20, 3, 30)
}

func TestChoose(t *testing.T) {
	assertList(t,
		list.Of(1, 2, 3, 4).Choose(func(n int) option.Option[string] {
			if n%2 == 0 {
				return option.Some(strconv.Itoa(n))
			}
			return option.None[string]()
		}),
		"2", "4",
	)
}

func TestFold(t *testing.T) {
	if list.Of(1, 2, 3, 4).Fold(0, func(acc, v int) int { return acc + v }) != 10 {
		t.Fatal("expected 10")
	}
}

func TestFoldBack(t *testing.T) {
	// right-fold: (1-(2-(3-0))) = 2
	got := list.Of(1, 2, 3).FoldBack(0, func(v, acc int) int { return v - acc })
	if got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestGroupBy(t *testing.T) {
	m := list.Of(1, 2, 3, 4, 5).GroupBy(func(n int) string {
		if n%2 == 0 {
			return "even"
		}
		return "odd"
	})
	if len(m["even"]) != 2 || len(m["odd"]) != 3 {
		t.Fatalf("got %v", m)
	}
}

func TestSortBy(t *testing.T) {
	assertList(t, list.Of(3, 1, 2).SortBy(func(n int) int { return n }), 1, 2, 3)
}

func TestSortByDescending(t *testing.T) {
	assertList(t, list.Of(1, 3, 2).SortByDescending(func(n int) int { return n }), 3, 2, 1)
}

func TestDistinctBy(t *testing.T) {
	assertList(t,
		list.Of(1, 2, 3, 4).DistinctBy(func(n int) int { return n % 2 }),
		1, 2,
	)
}

func TestScan(t *testing.T) {
	assertList(t, list.Of(1, 2, 3).Scan(0, func(s, v int) int { return s + v }), 0, 1, 3, 6)
}

func TestScanBack(t *testing.T) {
	assertList(t, list.Of(1, 2, 3).ScanBack(0, func(v, s int) int { return v + s }), 6, 5, 3, 0)
}

func TestMapFold(t *testing.T) {
	result, state := list.Of(1, 2, 3).MapFold(0, func(s, v int) (string, int) {
		return strconv.Itoa(v), s + v
	})
	assertList(t, result, "1", "2", "3")
	if state != 6 {
		t.Fatalf("expected state 6, got %d", state)
	}
}

func TestZipWith(t *testing.T) {
	assertList(t,
		list.Of(1, 2, 3).ZipWith(list.Of(10, 20, 30), func(a, b int) int { return a + b }),
		11, 22, 33,
	)
}

// ---------------------------------------------------------------------------
// Package-level functions
// ---------------------------------------------------------------------------

func TestZip(t *testing.T) {
	pairs := list.Zip(list.Of(1, 2), list.Of("a", "b")).ToSlice()
	if len(pairs) != 2 || pairs[0].First != 1 || pairs[0].Second != "a" {
		t.Fatalf("got %v", pairs)
	}
}

func TestDistinct(t *testing.T) {
	assertList(t, list.Distinct(list.Of(1, 2, 1, 3, 2)), 1, 2, 3)
}

func TestContains(t *testing.T) {
	if !list.Contains(list.Of(1, 2, 3), 2) {
		t.Fatal("expected true")
	}
	if list.Contains(list.Of(1, 2, 3), 9) {
		t.Fatal("expected false")
	}
}

func TestExcept(t *testing.T) {
	assertList(t, list.Except(list.Of(1, 2, 3, 4), list.Of(2, 4)), 1, 3)
}

// ---------------------------------------------------------------------------
// Chaining
// ---------------------------------------------------------------------------

func TestChain(t *testing.T) {
	got := list.Of(1, 2, 3, 4, 5, 6).
		Filter(func(n int) bool { return n%2 == 0 }).
		Map(func(n int) int { return n * n }).
		SortByDescending(func(n int) int { return n }).
		ToSlice()
	if len(got) != 3 || got[0] != 36 || got[1] != 16 || got[2] != 4 {
		t.Fatalf("got %v", got)
	}
}
