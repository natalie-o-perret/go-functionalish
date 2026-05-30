package tuple_test

import (
	"strings"
	"testing"

	"github.com/natalie-o-perret/go-functionalish/tuple"
)

func TestOf(t *testing.T) {
	p := tuple.Of(1, "hello")
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
	a, b := tuple.Of(42, "world").Unpack()
	if a != 42 || b != "world" {
		t.Fatalf("got %v %v", a, b)
	}
}

func TestT2Swap(t *testing.T) {
	s := tuple.Of("x", 99).Swap()
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
	got := tuple.Of(3, 4).Apply(func(a, b int) int { return a + b })
	if got != 7 {
		t.Fatalf("expected 7, got %d", got)
	}
}

func TestApply3(t *testing.T) {
	got := tuple.Of3(1, 2, 3).Apply(func(a, b, c int) int { return a + b + c })
	if got != 6 {
		t.Fatalf("expected 6, got %d", got)
	}
}

func TestApply4(t *testing.T) {
	got := tuple.Of4(1, 2, 3, 4).Apply(func(a, b, c, d int) int { return a + b + c + d })
	if got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestMapFirst(t *testing.T) {
	got := tuple.Of("hello", 42).MapFirst(strings.ToUpper)
	if got.First != "HELLO" || got.Second != 42 {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestMapSecond(t *testing.T) {
	got := tuple.Of(42, "hello").MapSecond(strings.ToUpper)
	if got.First != 42 || got.Second != "HELLO" {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestMap(t *testing.T) {
	got := tuple.Of("hello", "world").Map(strings.ToUpper, strings.ToUpper)
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
	if tupledAdd(tuple.Of(2, 3)) != 5 {
		t.Fatal("expected 5")
	}
	backToAdd := tuple.ToFunc2(tupledAdd)
	if backToAdd(2, 3) != 5 {
		t.Fatal("expected 5")
	}
}

func TestT2Extend(t *testing.T) {
	got := tuple.Of(1, "x").Extend(true)
	if got.First != 1 || got.Second != "x" || got.Third != true {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestT3Extend(t *testing.T) {
	got := tuple.Of3(1, "x", true).Extend(3.14)
	if got.First != 1 || got.Second != "x" || got.Third != true || got.Fourth != 3.14 {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestFst(t *testing.T) {
	if tuple.Fst(tuple.Of(42, "hello")) != 42 {
		t.Fatal("expected 42")
	}
}

func TestSnd(t *testing.T) {
	if tuple.Snd(tuple.Of(42, "hello")) != "hello" {
		t.Fatal("expected hello")
	}
}

func TestThd(t *testing.T) {
	if tuple.Thd(tuple.Of3(1, 2, 99)) != 99 {
		t.Fatal("expected 99")
	}
}

func TestFth(t *testing.T) {
	if tuple.Fth(tuple.Of4(1, 2, 3, 77)) != 77 {
		t.Fatal("expected 77")
	}
}

func TestMap2(t *testing.T) {
	got := tuple.Map2(tuple.Of(2, 3), func(n int) int { return n * 10 })
	if got.First != 20 || got.Second != 30 {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestMap3(t *testing.T) {
	got := tuple.Map3(tuple.Of3(1, 2, 3), func(n int) int { return n * 10 })
	if got.First != 10 || got.Second != 20 || got.Third != 30 {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestMap4(t *testing.T) {
	got := tuple.Map4(tuple.Of4(1, 2, 3, 4), func(n int) int { return n * 10 })
	if got.First != 10 || got.Second != 20 || got.Third != 30 || got.Fourth != 40 {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestMap2TypeChange(t *testing.T) {
	got := tuple.Map2(tuple.Of(1, 2), func(n int) string { return strings.Repeat("x", n) })
	if got.First != "x" || got.Second != "xx" {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestT3MapFirst(t *testing.T) {
	got := tuple.Of3("hello", 2, true).MapFirst(strings.ToUpper)
	if got.First != "HELLO" || got.Second != 2 || got.Third != true {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestT3MapSecond(t *testing.T) {
	got := tuple.Of3("a", 1, true).MapSecond(func(n int) int { return n * 10 })
	if got.First != "a" || got.Second != 10 || got.Third != true {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestT3MapThird(t *testing.T) {
	got := tuple.Of3("a", 1, false).MapThird(func(b bool) string {
		if b {
			return "yes"
		}
		return "no"
	})
	if got.First != "a" || got.Second != 1 || got.Third != "no" {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestT3Map(t *testing.T) {
	got := tuple.Of3("hello", 2, true).Map(
		strings.ToUpper,
		func(n int) int { return n * 3 },
		func(b bool) bool { return !b },
	)
	if got.First != "HELLO" || got.Second != 6 || got.Third != false {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestT3DropThird(t *testing.T) {
	got := tuple.Of3(1, "x", true).DropThird()
	if got.First != 1 || got.Second != "x" {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestT3DropFirst(t *testing.T) {
	got := tuple.Of3(1, "x", true).DropFirst()
	if got.First != "x" || got.Second != true {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestT3DropSecond(t *testing.T) {
	got := tuple.Of3(1, "x", true).DropSecond()
	if got.First != 1 || got.Second != true {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestT4MapFirst(t *testing.T) {
	got := tuple.Of4("a", 2, true, 3.14).MapFirst(strings.ToUpper)
	if got.First != "A" || got.Second != 2 || got.Third != true || got.Fourth != 3.14 {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestT4MapSecond(t *testing.T) {
	got := tuple.Of4("a", 2, true, 3.14).MapSecond(func(n int) int { return n * 5 })
	if got.First != "a" || got.Second != 10 || got.Third != true || got.Fourth != 3.14 {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestT4MapThird(t *testing.T) {
	got := tuple.Of4("a", 2, false, 3.14).MapThird(func(b bool) bool { return !b })
	if got.First != "a" || got.Second != 2 || got.Third != true || got.Fourth != 3.14 {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestT4MapFourth(t *testing.T) {
	got := tuple.Of4("a", 2, true, 3.14).MapFourth(func(f float64) float64 { return f * 2 })
	if got.First != "a" || got.Second != 2 || got.Third != true || got.Fourth != 6.28 {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestT4Map(t *testing.T) {
	got := tuple.Of4("hello", 2, false, 1.0).Map(
		strings.ToUpper,
		func(n int) int { return n + 8 },
		func(b bool) bool { return !b },
		func(f float64) float64 { return f * 10 },
	)
	if got.First != "HELLO" || got.Second != 10 || got.Third != true || got.Fourth != 10.0 {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestT4DropFourth(t *testing.T) {
	got := tuple.Of4(1, "x", true, 9.9).DropFourth()
	if got.First != 1 || got.Second != "x" || got.Third != true {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestT4DropThird(t *testing.T) {
	got := tuple.Of4(1, "x", true, 9.9).DropThird()
	if got.First != 1 || got.Second != "x" || got.Third != 9.9 {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestT4DropSecond(t *testing.T) {
	got := tuple.Of4(1, "x", true, 9.9).DropSecond()
	if got.First != 1 || got.Second != true || got.Third != 9.9 {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestT4DropFirst(t *testing.T) {
	got := tuple.Of4(1, "x", true, 9.9).DropFirst()
	if got.First != "x" || got.Second != true || got.Third != 9.9 {
		t.Fatalf("unexpected: %v", got)
	}
}
