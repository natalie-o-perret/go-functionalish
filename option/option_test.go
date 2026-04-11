package option_test

import (
	"strings"
	"testing"

	"github.com/natalie-o-perret/gof/option"
)

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
		if len(s) == 0 {
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
