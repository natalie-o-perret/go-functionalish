package option_test

import (
	"strings"
	"testing"

	"github.com/natalie-o-perret/go-functionalish/option"
)

func TestFlatten(t *testing.T) {
	if option.Flatten(option.Some(option.Some(42))).Unwrap() != 42 {
		t.Fatal("expected 42")
	}
	if option.Flatten(option.None[option.Option[int]]()).IsSome() {
		t.Fatal("expected None")
	}
}

func TestOrElse(t *testing.T) {
	got := option.OrElse(option.Some(1), func() option.Option[int] { return option.Some(2) })
	if got.Unwrap() != 1 {
		t.Fatal("expected Some(1) to win")
	}
	got = option.OrElse(option.None[int](), func() option.Option[int] { return option.Some(2) })
	if got.Unwrap() != 2 {
		t.Fatal("expected fallback Some(2)")
	}
}

func TestContains(t *testing.T) {
	if !option.Contains(option.Some(42), 42) {
		t.Fatal("expected true")
	}
	if option.Contains(option.Some(42), 99) {
		t.Fatal("expected false")
	}
	if option.Contains(option.None[int](), 42) {
		t.Fatal("expected false on None")
	}
}

func TestDefaultWith(t *testing.T) {
	called := false
	got := option.DefaultWith(option.Some(1), func() int { called = true; return 99 })
	if got != 1 || called {
		t.Fatal("fn should not be called when Some")
	}
	got = option.DefaultWith(option.None[int](), func() int { return 99 })
	if got != 99 {
		t.Fatal("expected 99")
	}
}

func TestSome(t *testing.T) {
	o := option.Some(42)
	if !o.IsSome() {
		t.Fatal("expected Some")
	}
	if o.Unwrap() != 42 {
		t.Fatalf("got %d", o.Unwrap())
	}
}

func TestNone(t *testing.T) {
	o := option.None[int]()
	if !o.IsNone() {
		t.Fatal("expected None")
	}
	if o.UnwrapOr(99) != 99 {
		t.Fatal("expected default")
	}
}

func TestMap(t *testing.T) {
	got := option.Map(option.Some("hello"), strings.ToUpper).UnwrapOr("")
	if got != "HELLO" {
		t.Fatalf("got %s", got)
	}
}

func TestMapNone(t *testing.T) {
	got := option.Map(option.None[string](), strings.ToUpper)
	if !got.IsNone() {
		t.Fatal("expected None")
	}
}

func TestBind(t *testing.T) {
	first := func(s string) option.Option[byte] {
		if s == "" {
			return option.None[byte]()
		}
		return option.Some(s[0])
	}
	got := option.Bind(option.Some("abc"), first).UnwrapOr(0)
	if got != 'a' {
		t.Fatalf("got %c", got)
	}
}

func TestFilter(t *testing.T) {
	even := option.Some(4).Filter(func(n int) bool { return n%2 == 0 })
	if !even.IsSome() {
		t.Fatal("expected Some")
	}
	odd := option.Some(3).Filter(func(n int) bool { return n%2 == 0 })
	if !odd.IsNone() {
		t.Fatal("expected None")
	}
}

func TestToSlice(t *testing.T) {
	if len(option.Some(1).ToSlice()) != 1 {
		t.Fatal("expected 1 element")
	}
	if len(option.None[int]().ToSlice()) != 0 {
		t.Fatal("expected empty")
	}
}
func TestMap2(t *testing.T) {
	got := option.Map2(option.Some(2), option.Some(3), func(a, b int) int { return a + b })
	if got.Unwrap() != 5 {
		t.Fatalf("expected 5, got %d", got.Unwrap())
	}
	if option.Map2(option.None[int](), option.Some(3), func(a, b int) int { return a + b }).IsSome() {
		t.Fatal("expected None when first is None")
	}
	if option.Map2(option.Some(2), option.None[int](), func(a, b int) int { return a + b }).IsSome() {
		t.Fatal("expected None when second is None")
	}
}
func TestTee(t *testing.T) {
	var seen int
	out := option.Tee(option.Some(7), func(v int) { seen = v })
	if seen != 7 || out.Unwrap() != 7 {
		t.Fatal("Tee should call fn and return Some unchanged")
	}
	seen = 0
	out2 := option.Tee(option.None[int](), func(v int) { seen = v })
	if seen != 0 || !out2.IsNone() {
		t.Fatal("Tee should not call fn on None")
	}
}
func TestTeeNone(t *testing.T) {
	called := false
	out := option.TeeNone(option.None[int](), func() { called = true })
	if !called || !out.IsNone() {
		t.Fatal("TeeNone should call fn on None")
	}
	called = false
	out2 := option.TeeNone(option.Some(1), func() { called = true })
	if called || out2.UnwrapOr(0) != 1 {
		t.Fatal("TeeNone should not call fn on Some")
	}
}
