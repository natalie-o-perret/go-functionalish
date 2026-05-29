package seq_test

import (
	"cmp"
	"fmt"
	"testing"

	"github.com/natalie-o-perret/go-functionalish/seq"
)

func TestFilter(t *testing.T) {
	got := seq.OfSlice([]int{1, 2, 3, 4, 5}).
		Filter(func(n int) bool { return n%2 == 0 }).
		ToSlice()
	want := []int{2, 4}
	for i, v := range got {
		if v != want[i] {
			t.Fatalf("pos %d: got %d want %d", i, v, want[i])
		}
	}
}

func TestMap(t *testing.T) {
	got := seq.Map(seq.Range(1, 4), func(n int) int { return n * n }).ToSlice()
	want := []int{1, 4, 9}
	for i, v := range got {
		if v != want[i] {
			t.Fatalf("pos %d: got %d", i, v)
		}
	}
}

func TestGroupBy(t *testing.T) {
	g := seq.GroupBy(seq.Range(1, 6), func(n int) string {
		if n%2 == 0 {
			return "even"
		}
		return "odd"
	})
	if len(g["even"]) != 2 {
		t.Fatalf("expected 2 even, got %v", g["even"])
	}
	if len(g["odd"]) != 3 {
		t.Fatalf("expected 3 odd, got %v", g["odd"])
	}
}

func TestFold(t *testing.T) {
	sum := seq.Fold(seq.Range(1, 6), 0, func(acc, v int) int { return acc + v })
	if sum != 15 {
		t.Fatalf("got %d", sum)
	}
}

func TestDistinctBy(t *testing.T) {
	got := seq.DistinctBy(seq.OfSlice([]int{1, 2, 1, 3, 2}), func(n int) int { return n }).ToSlice()
	if len(got) != 3 {
		t.Fatalf("expected 3 unique, got %v", got)
	}
}

func TestZip(t *testing.T) {
	pairs := seq.Zip(seq.Range(1, 4), seq.OfSlice([]string{"a", "b", "c"})).ToSlice()
	if len(pairs) != 3 {
		t.Fatalf("expected 3 pairs")
	}
	if pairs[0].First != 1 || pairs[0].Second != "a" {
		t.Fatalf("unexpected %v", pairs[0])
	}
}

func TestCycle(t *testing.T) {
	got := seq.OfSlice([]int{1, 2}).Cycle().Truncate(5).ToSlice()
	want := []int{1, 2, 1, 2, 1}
	for i, v := range got {
		if v != want[i] {
			t.Fatalf("pos %d: got %d", i, v)
		}
	}
}

func TestSortBy(t *testing.T) {
	got := seq.OfSlice([]int{3, 1, 2}).SortBy(func(n int) int { return n }).ToSlice()
	for i, v := range got {
		if v != i+1 {
			t.Fatalf("pos %d: got %d", i, v)
		}
	}
}

func TestSortByDescending(t *testing.T) {
	got := seq.OfSlice([]int{1, 3, 2}).SortByDescending(func(n int) int { return n }).ToSlice()
	want := []int{3, 2, 1}
	for i, v := range got {
		if v != want[i] {
			t.Fatalf("pos %d: got %d", i, v)
		}
	}
}

func TestTryHeadTryLast(t *testing.T) {
	s := seq.OfSlice([]int{10, 20, 30})
	if s.TryHead().UnwrapOr(0) != 10 {
		t.Fatal("wrong head")
	}
	if s.TryLast().UnwrapOr(0) != 30 {
		t.Fatal("wrong last")
	}
	if seq.Empty[int]().TryHead().IsNone() == false {
		t.Fatal("expected None")
	}
}

func TestExistsForAll(t *testing.T) {
	s := seq.Range(1, 6)
	if !s.Exists(func(n int) bool { return n > 4 }) {
		t.Fatal("exists > 4 should be true")
	}
	if s.ForAll(func(n int) bool { return n > 4 }) {
		t.Fatal("forall > 4 should be false")
	}
}

func TestAppend(t *testing.T) {
	got := seq.OfSlice([]int{1, 2}).Append(seq.OfSlice([]int{3, 4})).Length()
	if got != 4 {
		t.Fatalf("got %d", got)
	}
}

func TestSortWith(t *testing.T) {
	got := seq.OfSlice([]int{5, 3, 1, 4, 2}).
		SortWith(cmp.Compare).
		ToSlice()
	for i, v := range got {
		if v != i+1 {
			t.Fatalf("pos %d: got %d", i, v)
		}
	}
}

func TestZipWith(t *testing.T) {
	got := seq.OfSlice([]int{1, 2, 3}).
		ZipWith(seq.OfSlice([]string{"a", "b", "c"}), func(n int, s string) string {
			return fmt.Sprintf("%d%s", n, s)
		}).ToSlice()
	want := []string{"1a", "2b", "3c"}
	for i, v := range got {
		if v != want[i] {
			t.Fatalf("ZipWith[%d]: got %q want %q", i, v, want[i])
		}
	}
	// stops at shorter
	short := seq.OfSlice([]int{1, 2}).
		ZipWith(seq.OfSlice([]string{"x", "y", "z"}), func(n int, s string) string { return s }).
		ToSlice()
	if len(short) != 2 {
		t.Fatalf("ZipWith truncation: len %d", len(short))
	}
}
