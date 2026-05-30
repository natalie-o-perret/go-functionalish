package set_test

import (
	"slices"
	"strconv"
	"testing"

	"github.com/natalie-o-perret/go-functionalish/option"
	"github.com/natalie-o-perret/go-functionalish/set"
)

// containsAll asserts that s contains exactly the given values (order-independent).
func containsAll[T comparable](t *testing.T, s set.Set[T], want ...T) {
	t.Helper()
	if s.Len() != len(want) {
		t.Fatalf("len: got %d, want %d", s.Len(), len(want))
	}
	for _, v := range want {
		if !s.Contains(v) {
			t.Fatalf("missing %v", v)
		}
	}
}

// ---------------------------------------------------------------------------
// Constructors
// ---------------------------------------------------------------------------

func TestOf(t *testing.T) {
	s := set.Of(1, 2, 3, 2, 1)
	containsAll(t, s, 1, 2, 3)
}

func TestFromSlice(t *testing.T) {
	s := set.FromSlice([]string{"a", "b", "a"})
	containsAll(t, s, "a", "b")
}

func TestEmpty(t *testing.T) {
	s := set.Empty[int]()
	if !s.IsEmpty() {
		t.Fatal("expected empty")
	}
}

// ---------------------------------------------------------------------------
// Add / Remove
// ---------------------------------------------------------------------------

func TestAdd(t *testing.T) {
	s := set.Of(1, 2).Add(3, 4, 2)
	containsAll(t, s, 1, 2, 3, 4)
}

func TestRemove(t *testing.T) {
	s := set.Of(1, 2, 3).Remove(2, 4)
	containsAll(t, s, 1, 3)
}

func TestAdd_Immutability(t *testing.T) {
	a := set.Of(1, 2)
	_ = a.Add(3)
	if a.Contains(3) {
		t.Fatal("Add must not mutate receiver")
	}
}

// ---------------------------------------------------------------------------
// Set algebra
// ---------------------------------------------------------------------------

func TestUnion(t *testing.T) {
	a := set.Of(1, 2, 3)
	b := set.Of(3, 4, 5)
	containsAll(t, a.Union(b), 1, 2, 3, 4, 5)
}

func TestIntersect(t *testing.T) {
	a := set.Of(1, 2, 3)
	b := set.Of(2, 3, 4)
	containsAll(t, a.Intersect(b), 2, 3)
}

func TestDifference(t *testing.T) {
	a := set.Of(1, 2, 3)
	b := set.Of(2, 3, 4)
	containsAll(t, a.Difference(b), 1)
}

func TestSymmetricDifference(t *testing.T) {
	a := set.Of(1, 2, 3)
	b := set.Of(2, 3, 4)
	containsAll(t, a.SymmetricDifference(b), 1, 4)
}

func TestSymmetricDifference_Commutative(t *testing.T) {
	a := set.Of(1, 2, 3)
	b := set.Of(2, 3, 4)
	ab := a.SymmetricDifference(b)
	ba := b.SymmetricDifference(a)
	if ab.Len() != ba.Len() {
		t.Fatal("symmetric difference must be commutative")
	}
	for v := range ab.All() {
		if !ba.Contains(v) {
			t.Fatalf("missing %v in ba", v)
		}
	}
}

// ---------------------------------------------------------------------------
// Query methods
// ---------------------------------------------------------------------------

func TestContains(t *testing.T) {
	s := set.Of(1, 2, 3)
	if !s.Contains(2) {
		t.Fatal("expected true")
	}
	if s.Contains(9) {
		t.Fatal("expected false")
	}
}

func TestIsSubset(t *testing.T) {
	a := set.Of(1, 2)
	b := set.Of(1, 2, 3)
	if !a.IsSubset(b) {
		t.Fatal("{1,2} ⊆ {1,2,3}")
	}
	if b.IsSubset(a) {
		t.Fatal("{1,2,3} ⊄ {1,2}")
	}
}

func TestIsSuperset(t *testing.T) {
	a := set.Of(1, 2, 3)
	b := set.Of(1, 2)
	if !a.IsSuperset(b) {
		t.Fatal("{1,2,3} ⊇ {1,2}")
	}
}

func TestFilter(t *testing.T) {
	s := set.Of(1, 2, 3, 4, 5)
	evens := s.Filter(func(n int) bool { return n%2 == 0 })
	containsAll(t, evens, 2, 4)
}

func TestExclude(t *testing.T) {
	s := set.Of(1, 2, 3, 4)
	odds := s.Exclude(func(n int) bool { return n%2 == 0 })
	containsAll(t, odds, 1, 3)
}

func TestForAll(t *testing.T) {
	s := set.Of(2, 4, 6)
	if !s.ForAll(func(n int) bool { return n%2 == 0 }) {
		t.Fatal("expected true")
	}
	if s.Add(3).ForAll(func(n int) bool { return n%2 == 0 }) {
		t.Fatal("expected false")
	}
}

func TestExists(t *testing.T) {
	s := set.Of(1, 3, 5)
	if s.Exists(func(n int) bool { return n%2 == 0 }) {
		t.Fatal("expected false")
	}
	if !s.Add(4).Exists(func(n int) bool { return n%2 == 0 }) {
		t.Fatal("expected true")
	}
}

func TestCountBy(t *testing.T) {
	s := set.Of(1, 2, 3, 4, 5)
	if s.CountBy(func(n int) bool { return n > 3 }) != 2 {
		t.Fatal("expected 2")
	}
}

// ---------------------------------------------------------------------------
// Iteration / materialise
// ---------------------------------------------------------------------------

func TestAll(t *testing.T) {
	s := set.Of(1, 2, 3)
	var got []int
	for v := range s.All() {
		got = append(got, v)
	}
	slices.Sort(got)
	if len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
		t.Fatalf("got %v", got)
	}
}

func TestToSlice(t *testing.T) {
	s := set.Of(3, 1, 2)
	got := s.ToSlice()
	slices.Sort(got)
	if len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
		t.Fatalf("got %v", got)
	}
}

func TestIter(t *testing.T) {
	s := set.Of(1, 2, 3)
	sum := 0
	s.Iter(func(n int) { sum += n })
	if sum != 6 {
		t.Fatalf("expected 6, got %d", sum)
	}
}

// ---------------------------------------------------------------------------
// Generic methods
// ---------------------------------------------------------------------------

func TestMap(t *testing.T) {
	s := set.Of(1, 2, 3)
	got := s.Map(strconv.Itoa)
	containsAll(t, got, "1", "2", "3")
}

func TestMap_CollapsesDuplicates(t *testing.T) {
	s := set.Of(1, 2, 3, 4)
	got := s.Map(func(n int) int { return n % 2 }) // both even→0, odd→1
	if got.Len() != 2 {
		t.Fatalf("expected 2, got %d", got.Len())
	}
}

func TestCollect(t *testing.T) {
	s := set.Of(1, 2)
	got := s.Collect(func(n int) []int { return []int{n, n * 10} })
	containsAll(t, got, 1, 2, 10, 20)
}

func TestChoose(t *testing.T) {
	s := set.Of(1, 2, 3, 4)
	got := s.Choose(func(n int) option.Option[string] {
		if n%2 == 0 {
			return option.Some(strconv.Itoa(n))
		}
		return option.None[string]()
	})
	containsAll(t, got, "2", "4")
}

func TestFold(t *testing.T) {
	s := set.Of(1, 2, 3, 4, 5)
	sum := s.Fold(0, func(acc, v int) int { return acc + v })
	if sum != 15 {
		t.Fatalf("expected 15, got %d", sum)
	}
}

// ---------------------------------------------------------------------------
// Chaining
// ---------------------------------------------------------------------------

func TestChain(t *testing.T) {
	got := set.Of(1, 2, 3, 4, 5, 6).
		Filter(func(n int) bool { return n > 1 }).
		Map(func(n int) int { return n * n }).
		Difference(set.Of(4, 9)).
		ToSlice()
	slices.Sort(got)
	want := []int{16, 25, 36}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}
