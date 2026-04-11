package enum_test

import (
	"cmp"
	"testing"

	"github.com/natalie-o-perret/gof/enum"
)

func TestFilterMap(t *testing.T) {
	got := enum.From([]int{1, 2, 3, 4, 5}).
		Filter(func(n int) bool { return n%2 == 0 }).
		ToSlice()
	want := []int{2, 4}
	for i, v := range got {
		if v != want[i] { t.Fatalf("pos %d: got %d want %d", i, v, want[i]) }
	}
}

func TestMap(t *testing.T) {
	got := enum.Map(enum.Range(1, 4), func(n int) int { return n * n }).ToSlice()
	want := []int{1, 4, 9}
	for i, v := range got { if v != want[i] { t.Fatalf("pos %d: got %d", i, v) } }
}

func TestGroupBy(t *testing.T) {
	g := enum.GroupBy(enum.Range(1, 6), func(n int) string {
		if n%2 == 0 { return "even" }
		return "odd"
	})
	if len(g["even"]) != 2 { t.Fatalf("expected 2 even, got %v", g["even"]) }
	if len(g["odd"]) != 3 { t.Fatalf("expected 3 odd, got %v", g["odd"]) }
}

func TestReduce(t *testing.T) {
	sum := enum.Reduce(enum.Range(1, 6), 0, func(acc, v int) int { return acc + v })
	if sum != 15 { t.Fatalf("got %d", sum) }
}

func TestUniqBy(t *testing.T) {
	got := enum.UniqBy(enum.From([]int{1, 2, 1, 3, 2}), func(n int) int { return n }).ToSlice()
	if len(got) != 3 { t.Fatalf("expected 3 unique, got %v", got) }
}

func TestZip(t *testing.T) {
	pairs := enum.Zip(enum.Range(1, 4), enum.From([]string{"a", "b", "c"})).ToSlice()
	if len(pairs) != 3 { t.Fatalf("expected 3 pairs") }
	if pairs[0].First != 1 || pairs[0].Second != "a" { t.Fatalf("unexpected %v", pairs[0]) }
}

func TestCycle(t *testing.T) {
	got := enum.From([]int{1, 2}).Cycle().Take(5).ToSlice()
	want := []int{1, 2, 1, 2, 1}
	for i, v := range got { if v != want[i] { t.Fatalf("pos %d: got %d", i, v) } }
}

func TestSortAscBy(t *testing.T) {
	got := enum.SortAscBy(enum.From([]int{3, 1, 2}), func(n int) int { return n }).ToSlice()
	for i, v := range got { if v != i+1 { t.Fatalf("pos %d: got %d", i, v) } }
}

func TestSortDescBy(t *testing.T) {
	got := enum.SortDescBy(enum.From([]int{1, 3, 2}), func(n int) int { return n }).ToSlice()
	want := []int{3, 2, 1}
	for i, v := range got { if v != want[i] { t.Fatalf("pos %d: got %d", i, v) } }
}

func TestFirstLastOption(t *testing.T) {
	e := enum.From([]int{10, 20, 30})
	if e.FirstOption().UnwrapOr(0) != 10 { t.Fatal("wrong first") }
	if e.LastOption().UnwrapOr(0) != 30 { t.Fatal("wrong last") }
	if enum.Empty[int]().FirstOption().IsNone() == false { t.Fatal("expected None") }
}

func TestAnyEvery(t *testing.T) {
	e := enum.Range(1, 6)
	if !e.Any(func(n int) bool { return n > 4 }) { t.Fatal("any > 4 should be true") }
	if e.Every(func(n int) bool { return n > 4 }) { t.Fatal("every > 4 should be false") }
}

func TestChain(t *testing.T) {
	got := enum.From([]int{1, 2}).Chain(enum.From([]int{3, 4})).Count()
	if got != 4 { t.Fatalf("got %d", got) }
}

func TestSortBy(t *testing.T) {
	got := enum.From([]int{5, 3, 1, 4, 2}).
		SortBy(cmp.Compare).
		ToSlice()
	for i, v := range got { if v != i+1 { t.Fatalf("pos %d: got %d", i, v) } }
}

