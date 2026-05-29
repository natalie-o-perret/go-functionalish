package option_test

import (
	"fmt"
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
	got := option.Some(1).OrElse(func() option.Option[int] { return option.Some(2) })
	if got.Unwrap() != 1 {
		t.Fatal("expected Some(1) to win")
	}
	got = option.None[int]().OrElse(func() option.Option[int] { return option.Some(2) })
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
	got := option.Some(1).DefaultWith(func() int { called = true; return 99 })
	if got != 1 || called {
		t.Fatal("fn should not be called when Some")
	}
	got = option.None[int]().DefaultWith(func() int { return 99 })
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
	got := option.Some("hello").Map(strings.ToUpper).UnwrapOr("")
	if got != "HELLO" {
		t.Fatalf("got %s", got)
	}
}

func TestMapNone(t *testing.T) {
	got := option.None[string]().Map(strings.ToUpper)
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
	got := option.Some("abc").Bind(first).UnwrapOr(0)
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
	got := option.Some(2).ZipWith(option.Some(3), func(a, b int) int { return a + b })
	if got.Unwrap() != 5 {
		t.Fatalf("expected 5, got %d", got.Unwrap())
	}
	if option.None[int]().ZipWith(option.Some(3), func(a, b int) int { return a + b }).IsSome() {
		t.Fatal("expected None when first is None")
	}
	if option.Some(2).ZipWith(option.None[int](), func(a, b int) int { return a + b }).IsSome() {
		t.Fatal("expected None when second is None")
	}
}
func TestTee(t *testing.T) {
	var seen int
	out := option.Some(7).Tee(func(v int) { seen = v })
	if seen != 7 || out.Unwrap() != 7 {
		t.Fatal("Tee should call fn and return Some unchanged")
	}
	seen = 0
	out2 := option.None[int]().Tee(func(v int) { seen = v })
	if seen != 0 || !out2.IsNone() {
		t.Fatal("Tee should not call fn on None")
	}
}
func TestTeeNone(t *testing.T) {
	called := false
	out := option.None[int]().TeeNone(func() { called = true })
	if !called || !out.IsNone() {
		t.Fatal("TeeNone should call fn on None")
	}
	called = false
	out2 := option.Some(1).TeeNone(func() { called = true })
	if called || out2.UnwrapOr(0) != 1 {
		t.Fatal("TeeNone should not call fn on Some")
	}
}

func TestZipWith(t *testing.T) {
	// both Some
	got := option.Some(3).ZipWith(option.Some("px"), func(n int, s string) string {
		return fmt.Sprintf("%d%s", n, s)
	})
	if !got.IsSome() || got.Unwrap() != "3px" {
		t.Fatalf("ZipWith Some+Some: got %v", got)
	}
	// first is None
	none := option.None[int]().ZipWith(option.Some("px"), func(n int, s string) string { return s })
	if !none.IsNone() {
		t.Fatal("ZipWith None+Some should be None")
	}
	// second is None
	none2 := option.Some(3).ZipWith(option.None[string](), func(n int, s string) string { return s })
	if !none2.IsNone() {
		t.Fatal("ZipWith Some+None should be None")
	}
}
