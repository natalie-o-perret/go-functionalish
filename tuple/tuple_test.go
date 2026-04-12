package tuple_test

import (
	"strings"
	"testing"

	"github.com/natalie-o-perret/gof/tuple"
)

func TestOf2(t *testing.T) {
	p := tuple.Of2(1, "hello")
	if p.First != 1 || p.Second != "hello" {
		t.Fatalf("unexpected T2: %v", p)
	}
}

func TestOf3(t *testing.T) {
	tr := tuple.Of3(1, "hello", true)
	if tr.First != 1 || tr.Second != "hello" || tr.Third != true {
		t.Fatalf("unexpected T3: %v", tr)
	}
}

func TestOf4(t *testing.T) {
	q := tuple.Of4(1, "hello", true, 3.14)
	if q.First != 1 || q.Second != "hello" || q.Third != true || q.Fourth != 3.14 {
		t.Fatalf("unexpected T4: %v", q)
	}
}

func TestT2Unpack(t *testing.T) {
	a, b := tuple.Of2(42, "world").Unpack()
	if a != 42 || b != "world" {
		t.Fatalf("got %v %v", a, b)
	}
}

func TestT2Swap(t *testing.T) {
	s := tuple.Of2("x", 99).Swap()
	if s.First != 99 || s.Second != "x" {
		t.Fatalf("unexpected swap: %v", s)
	}
}

func TestT3Unpack(t *testing.T) {
	a, b, c := tuple.Of3(1, 2, 3).Unpack()
	if a != 1 || b != 2 || c != 3 {
		t.Fatalf("got %v %v %v", a, b, c)
	}
}

func TestT4Unpack(t *testing.T) {
	a, b, c, d := tuple.Of4(1, 2, 3, 4).Unpack()
	if a != 1 || b != 2 || c != 3 || d != 4 {
		t.Fatalf("got %v %v %v %v", a, b, c, d)
	}
}

func TestApply(t *testing.T) {
	got := tuple.Apply(tuple.Of2(3, 4), func(a, b int) int { return a + b })
	if got != 7 {
		t.Fatalf("expected 7, got %d", got)
	}
}

func TestApply3(t *testing.T) {
	got := tuple.Apply3(tuple.Of3(1, 2, 3), func(a, b, c int) int { return a + b + c })
	if got != 6 {
		t.Fatalf("expected 6, got %d", got)
	}
}

func TestApply4(t *testing.T) {
	got := tuple.Apply4(tuple.Of4(1, 2, 3, 4), func(a, b, c, d int) int { return a + b + c + d })
	if got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestMapFirst(t *testing.T) {
	got := tuple.MapFirst(tuple.Of2("hello", 42), strings.ToUpper)
	if got.First != "HELLO" || got.Second != 42 {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestMapSecond(t *testing.T) {
	got := tuple.MapSecond(tuple.Of2(42, "hello"), strings.ToUpper)
	if got.First != 42 || got.Second != "HELLO" {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestMap(t *testing.T) {
	got := tuple.Map(tuple.Of2("hello", "world"), strings.ToUpper, strings.ToUpper)
	if got.First != "HELLO" || got.Second != "WORLD" {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestCurryUncurry(t *testing.T) {
	add := func(a, b int) int { return a + b }
	curried := tuple.Curry(add)
	if curried(3)(4) != 7 {
		t.Fatal("curried(3)(4) should be 7")
	}
	back := tuple.Uncurry(curried)
	if back(3, 4) != 7 {
		t.Fatal("uncurried(3,4) should be 7")
	}
}

func TestFromFunc2_ToFunc2(t *testing.T) {
	add := func(a, b int) int { return a + b }
	tupledAdd := tuple.FromFunc2(add)
	if tupledAdd(tuple.Of2(2, 3)) != 5 {
		t.Fatal("expected 5")
	}
	backToAdd := tuple.ToFunc2(tupledAdd)
	if backToAdd(2, 3) != 5 {
		t.Fatal("expected 5")
	}
}
