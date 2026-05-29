package slice_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/natalie-o-perret/go-functionalish/option"
	goslice "github.com/natalie-o-perret/go-functionalish/slice"
)

// helpers

func eq[T comparable](a, b goslice.Slice[T]) bool {
	return slices.Equal([]T(a), []T(b))
}

// -- constructors --

func TestOf(t *testing.T) {
	s := goslice.Of(1, 2, 3)
	if !eq(s, goslice.Of(1, 2, 3)) {
		t.Fatal("Of failed")
	}
}

func TestFromSlice(t *testing.T) {
	s := goslice.FromSlice([]string{"a", "b"})
	if !eq(s, goslice.Of("a", "b")) {
		t.Fatal("FromSlice failed")
	}
}

func TestReplicate(t *testing.T) {
	s := goslice.Replicate(7, 3)
	if !eq(s, goslice.Of(7, 7, 7)) {
		t.Fatal("Replicate failed")
	}
}

func TestInit(t *testing.T) {
	s := goslice.Init(4, func(i int) int { return i * i })
	if !eq(s, goslice.Of(0, 1, 4, 9)) {
		t.Fatal("Init failed")
	}
}

// -- same-type methods --

func TestFilter(t *testing.T) {
	got := goslice.Of(1, 2, 3, 4, 5).Filter(func(n int) bool { return n%2 == 0 })
	if !eq(got, goslice.Of(2, 4)) {
		t.Fatalf("Filter: got %v", got)
	}
}

func TestExclude(t *testing.T) {
	got := goslice.Of(1, 2, 3, 4).Exclude(func(n int) bool { return n%2 == 0 })
	if !eq(got, goslice.Of(1, 3)) {
		t.Fatalf("Exclude: got %v", got)
	}
}

func TestAppend(t *testing.T) {
	got := goslice.Of(1, 2).Append(3, 4)
	if !eq(got, goslice.Of(1, 2, 3, 4)) {
		t.Fatalf("Append: got %v", got)
	}
}

func TestTruncate(t *testing.T) {
	if !eq(goslice.Of(1, 2, 3, 4).Truncate(2), goslice.Of(1, 2)) {
		t.Fatal("Truncate within bounds")
	}
	if !eq(goslice.Of(1, 2).Truncate(10), goslice.Of(1, 2)) {
		t.Fatal("Truncate beyond bounds")
	}
}

func TestSkip(t *testing.T) {
	if !eq(goslice.Of(1, 2, 3, 4).Skip(2), goslice.Of(3, 4)) {
		t.Fatal("Skip within bounds")
	}
	if !eq(goslice.Of(1, 2).Skip(5), goslice.Of[int]()) {
		t.Fatal("Skip beyond bounds")
	}
}

func TestTakeWhile(t *testing.T) {
	got := goslice.Of(1, 2, 3, 4).TakeWhile(func(n int) bool { return n < 3 })
	if !eq(got, goslice.Of(1, 2)) {
		t.Fatalf("TakeWhile: got %v", got)
	}
}

func TestSkipWhile(t *testing.T) {
	got := goslice.Of(1, 2, 3, 4).SkipWhile(func(n int) bool { return n < 3 })
	if !eq(got, goslice.Of(3, 4)) {
		t.Fatalf("SkipWhile: got %v", got)
	}
}

func TestRev(t *testing.T) {
	if !eq(goslice.Of(1, 2, 3).Rev(), goslice.Of(3, 2, 1)) {
		t.Fatal("Rev failed")
	}
}

func TestSortWith(t *testing.T) {
	got := goslice.Of(3, 1, 2).SortWith(func(a, b int) int { return a - b })
	if !eq(got, goslice.Of(1, 2, 3)) {
		t.Fatalf("SortWith: got %v", got)
	}
}

func TestTail(t *testing.T) {
	if !eq(goslice.Of(1, 2, 3).Tail(), goslice.Of(2, 3)) {
		t.Fatal("Tail failed")
	}
}

func TestInsertAt(t *testing.T) {
	got := goslice.Of(1, 2, 4).InsertAt(2, 3)
	if !eq(got, goslice.Of(1, 2, 3, 4)) {
		t.Fatalf("InsertAt: got %v", got)
	}
}

func TestRemoveAt(t *testing.T) {
	got := goslice.Of(1, 2, 99, 3).RemoveAt(2)
	if !eq(got, goslice.Of(1, 2, 3)) {
		t.Fatalf("RemoveAt: got %v", got)
	}
}

func TestUpdateAt(t *testing.T) {
	got := goslice.Of(1, 2, 0, 4).UpdateAt(2, 3)
	if !eq(got, goslice.Of(1, 2, 3, 4)) {
		t.Fatalf("UpdateAt: got %v", got)
	}
}

// -- terminal / query methods --

func TestIterAndIteri(t *testing.T) {
	sum := 0
	goslice.Of(1, 2, 3).Iter(func(n int) { sum += n })
	if sum != 6 {
		t.Fatalf("Iter: sum %d", sum)
	}

	indexed := 0
	goslice.Of(10, 20, 30).Iteri(func(i, n int) { indexed += i * n })
	if indexed != 0*10+1*20+2*30 {
		t.Fatalf("Iteri: got %d", indexed)
	}
}

func TestIsEmptyAndLength(t *testing.T) {
	if !goslice.Of[int]().IsEmpty() {
		t.Fatal("IsEmpty on empty")
	}
	if goslice.Of(1).IsEmpty() {
		t.Fatal("IsEmpty on non-empty")
	}
	if goslice.Of(1, 2, 3).Length() != 3 {
		t.Fatal("Length")
	}
}

func TestHeadLast(t *testing.T) {
	s := goslice.Of(10, 20, 30)
	if h, ok := s.Head(); !ok || h != 10 {
		t.Fatal("Head")
	}
	if _, ok := goslice.Of[int]().Head(); ok {
		t.Fatal("Head on empty should be false")
	}
	if l, ok := s.Last(); !ok || l != 30 {
		t.Fatal("Last")
	}
	if s.TryHead().Unwrap() != 10 {
		t.Fatal("TryHead")
	}
	if s.TryLast().Unwrap() != 30 {
		t.Fatal("TryLast")
	}
	if goslice.Of[int]().TryHead().IsSome() {
		t.Fatal("TryHead empty")
	}
}

func TestItem(t *testing.T) {
	s := goslice.Of("a", "b", "c")
	if s.Item(1).Unwrap() != "b" {
		t.Fatal("Item(1)")
	}
	if s.Item(-1).IsSome() || s.Item(3).IsSome() {
		t.Fatal("Item out of bounds")
	}
}

func TestCountByExistsForAll(t *testing.T) {
	s := goslice.Of(1, 2, 3, 4)
	if s.CountBy(func(n int) bool { return n > 2 }) != 2 {
		t.Fatal("CountBy")
	}
	if !s.Exists(func(n int) bool { return n == 3 }) {
		t.Fatal("Exists true")
	}
	if s.Exists(func(n int) bool { return n == 99 }) {
		t.Fatal("Exists false")
	}
	if !s.ForAll(func(n int) bool { return n > 0 }) {
		t.Fatal("ForAll true")
	}
	if s.ForAll(func(n int) bool { return n > 2 }) {
		t.Fatal("ForAll false")
	}
}

func TestTryFind(t *testing.T) {
	s := goslice.Of(1, 2, 3, 4)
	if s.TryFind(func(n int) bool { return n > 2 }).Unwrap() != 3 {
		t.Fatal("TryFind")
	}
	if s.TryFindBack(func(n int) bool { return n < 3 }).Unwrap() != 2 {
		t.Fatal("TryFindBack")
	}
	if s.TryFind(func(n int) bool { return n > 99 }).IsSome() {
		t.Fatal("TryFind none")
	}
}

func TestReduce(t *testing.T) {
	sum, ok := goslice.Of(1, 2, 3, 4).Reduce(func(a, b int) int { return a + b })
	if !ok || sum != 10 {
		t.Fatal("Reduce")
	}
	_, ok = goslice.Of[int]().Reduce(func(a, b int) int { return a + b })
	if ok {
		t.Fatal("Reduce empty should be false")
	}
}

func TestToSlice(t *testing.T) {
	s := goslice.Of(1, 2, 3)
	raw := s.ToSlice()
	if !slices.Equal(raw, []int{1, 2, 3}) {
		t.Fatal("ToSlice")
	}
}

// -- generic methods --

func TestMap(t *testing.T) {
	got := goslice.Of(1, 2, 3).Map(func(n int) string { return strings.Repeat("x", n) })
	if !eq(got, goslice.Of("x", "xx", "xxx")) {
		t.Fatalf("Map: got %v", got)
	}
}

func TestMapi(t *testing.T) {
	got := goslice.Of("a", "b", "c").Mapi(func(i int, s string) string { return strings.Repeat(s, i+1) })
	if !eq(got, goslice.Of("a", "bb", "ccc")) {
		t.Fatalf("Mapi: got %v", got)
	}
}

func TestCollect(t *testing.T) {
	got := goslice.Of(1, 2, 3).Collect(func(n int) []int { return goslice.Replicate(n, n).ToSlice() })
	if !eq(got, goslice.Of(1, 2, 2, 3, 3, 3)) {
		t.Fatalf("Collect: got %v", got)
	}
}

func TestChoose(t *testing.T) {
	got := goslice.Of(1, 2, 3, 4).Choose(func(n int) option.Option[string] {
		if n%2 == 0 {
			return option.Some(strings.Repeat("*", n))
		}
		return option.None[string]()
	})
	if !eq(got, goslice.Of("**", "****")) {
		t.Fatalf("Choose: got %v", got)
	}
}

func TestFold(t *testing.T) {
	sum := goslice.Of(1, 2, 3, 4).Fold(0, func(acc, n int) int { return acc + n })
	if sum != 10 {
		t.Fatalf("Fold: got %d", sum)
	}
}

func TestFoldBack(t *testing.T) {
	// FoldBack on [1,2,3] with (n, acc) -> n::acc from initial [] reproduces the slice
	got := goslice.Of(1, 2, 3).FoldBack([]int{}, func(n int, acc []int) []int {
		return append([]int{n}, acc...)
	})
	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("FoldBack: got %v", got)
	}
}

func TestGroupBy(t *testing.T) {
	m := goslice.Of(1, 2, 3, 4, 5, 6).GroupBy(func(n int) int { return n % 2 })
	if len(m) != 2 {
		t.Fatal("GroupBy len")
	}
	if !slices.Equal(m[0], []int{2, 4, 6}) {
		t.Fatalf("GroupBy even: got %v", m[0])
	}
	if !slices.Equal(m[1], []int{1, 3, 5}) {
		t.Fatalf("GroupBy odd: got %v", m[1])
	}
}

func TestSortBy(t *testing.T) {
	got := goslice.Of("banana", "apple", "cherry").SortBy(func(s string) string { return s })
	if !eq(got, goslice.Of("apple", "banana", "cherry")) {
		t.Fatalf("SortBy: got %v", got)
	}
}

func TestSortByDescending(t *testing.T) {
	got := goslice.Of(3, 1, 4, 1, 5).SortByDescending(func(n int) int { return n })
	if !eq(got, goslice.Of(5, 4, 3, 1, 1)) {
		t.Fatalf("SortByDescending: got %v", got)
	}
}

func TestDistinctBy(t *testing.T) {
	got := goslice.Of("foo", "bar", "FOO", "baz", "BAR").DistinctBy(strings.ToLower)
	if !eq(got, goslice.Of("foo", "bar", "baz")) {
		t.Fatalf("DistinctBy: got %v", got)
	}
}

func TestTryPick(t *testing.T) {
	got := goslice.Of(1, 2, 3, 4).TryPick(func(n int) option.Option[string] {
		if n > 2 {
			return option.Some(strings.Repeat("!", n))
		}
		return option.None[string]()
	})
	if got.Unwrap() != "!!!" {
		t.Fatalf("TryPick: got %v", got)
	}
	none := goslice.Of(1, 2).TryPick(func(_ int) option.Option[int] { return option.None[int]() })
	if none.IsSome() {
		t.Fatal("TryPick none")
	}
}

func TestScan(t *testing.T) {
	// running sums: [0, 1, 3, 6, 10]
	got := goslice.Of(1, 2, 3, 4).Scan(0, func(acc, n int) int { return acc + n })
	if !eq(got, goslice.Of(0, 1, 3, 6, 10)) {
		t.Fatalf("Scan: got %v", got)
	}
}

func TestScanBack(t *testing.T) {
	// running sums from right: [6, 5, 3, 0] (last is initial 0)
	got := goslice.Of(1, 2, 3).ScanBack(0, func(n, acc int) int { return n + acc })
	if !eq(got, goslice.Of(6, 5, 3, 0)) {
		t.Fatalf("ScanBack: got %v", got)
	}
}

func TestMapFold(t *testing.T) {
	doubled, sum := goslice.Of(1, 2, 3).MapFold(0, func(acc, n int) (int, int) {
		return n * 2, acc + n
	})
	if !eq(doubled, goslice.Of(2, 4, 6)) {
		t.Fatalf("MapFold items: got %v", doubled)
	}
	if sum != 6 {
		t.Fatalf("MapFold state: got %d", sum)
	}
}

// -- package-level functions --

func TestZip(t *testing.T) {
	a := goslice.Of(1, 2, 3)
	b := goslice.Of("a", "b", "c")
	got := goslice.Zip(a, b)
	want := []option.Pair[int, string]{{First: 1, Second: "a"}, {First: 2, Second: "b"}, {First: 3, Second: "c"}}
	if len(got) != len(want) {
		t.Fatalf("Zip len: %d", len(got))
	}
	for i, p := range got {
		if p != want[i] {
			t.Fatalf("Zip[%d]: got %v want %v", i, p, want[i])
		}
	}

	// shorter slice wins
	short := goslice.Zip(goslice.Of(1, 2), goslice.Of("x", "y", "z"))
	if len(short) != 2 {
		t.Fatalf("Zip truncation: len %d", len(short))
	}
}

func TestDistinct(t *testing.T) {
	got := goslice.Distinct(goslice.Of(1, 2, 1, 3, 2, 4))
	if !eq(got, goslice.Of(1, 2, 3, 4)) {
		t.Fatalf("Distinct: got %v", got)
	}
}

func TestContains(t *testing.T) {
	s := goslice.Of(1, 2, 3)
	if !goslice.Contains(s, 2) {
		t.Fatal("Contains true")
	}
	if goslice.Contains(s, 99) {
		t.Fatal("Contains false")
	}
}

func TestExcept(t *testing.T) {
	got := goslice.Except(goslice.Of(1, 2, 3, 4, 5), goslice.Of(2, 4))
	if !eq(got, goslice.Of(1, 3, 5)) {
		t.Fatalf("Except: got %v", got)
	}
}

// -- fluent pipeline smoke test --

func TestFluentPipeline(t *testing.T) {
	type user struct {
		name   string
		active bool
		score  int
	}
	users := goslice.Of(
		user{"alice", true, 80},
		user{"bob", false, 90},
		user{"carol", true, 70},
		user{"dave", true, 85},
	)

	names := users.
		Filter(func(u user) bool { return u.active }).
		SortByDescending(func(u user) int { return u.score }).
		Map(func(u user) string { return u.name })

	if !eq(names, goslice.Of("dave", "alice", "carol")) {
		t.Fatalf("fluent pipeline: got %v", names)
	}
}
