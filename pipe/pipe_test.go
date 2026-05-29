package pipe_test

import (
	"strings"
	"testing"

	"github.com/natalie-o-perret/go-functionalish/pipe"
)

func TestPipe2(t *testing.T) {
	got := pipe.Pipe2("hello", strings.ToUpper, func(s string) int { return len(s) })
	if got != 5 {
		t.Fatalf("got %d", got)
	}
}

func TestPipe3(t *testing.T) {
	got := pipe.Pipe3(" hello ", strings.TrimSpace, strings.ToUpper, func(s string) string { return s + "!" })
	if got != "HELLO!" {
		t.Fatalf("got %s", got)
	}
}

func TestPipeEndoN(t *testing.T) {
	got := pipe.PipeEndoN(
		" hello world ",
		strings.TrimSpace,
		strings.ToUpper,
		func(s string) string { return s + "!" },
		func(s string) string { return "[" + s + "]" },
		func(s string) string { return s + "!" },
	)
	if got != "[HELLO WORLD!]!" {
		t.Fatalf("got %s", got)
	}

	// zero steps: value passes through unchanged
	if pipe.PipeEndoN(42) != 42 {
		t.Fatal("zero-step PipeEndoN should return value unchanged")
	}
}

func TestPipe4(t *testing.T) {
	got := pipe.Pipe4(
		" hello ",
		strings.TrimSpace,
		strings.ToUpper,
		func(s string) string { return s + "!" },
		func(s string) int { return len(s) },
	)
	if got != 6 {
		t.Fatalf("got %d", got)
	}
}

func TestPipe5(t *testing.T) {
	got := pipe.Pipe5(
		1,
		func(n int) int { return n + 1 },
		func(n int) int { return n * 2 },
		func(n int) int { return n + 3 },
		func(n int) int { return n * 4 },
		func(n int) string { return strings.Repeat("x", n) },
	)
	// ((1+1)*2+3)*4 = 28 x's
	if len(got) != 28 {
		t.Fatalf("got len %d", len(got))
	}
}

func TestPipe6(t *testing.T) {
	got := pipe.Pipe6(
		0,
		func(n int) int { return n + 1 },
		func(n int) int { return n + 1 },
		func(n int) int { return n + 1 },
		func(n int) int { return n + 1 },
		func(n int) int { return n + 1 },
		func(n int) int { return n * 10 },
	)
	if got != 50 {
		t.Fatalf("got %d", got)
	}
}

func TestPipe7(t *testing.T) {
	got := pipe.Pipe7(
		0,
		func(n int) int { return n + 1 },
		func(n int) int { return n + 1 },
		func(n int) int { return n + 1 },
		func(n int) int { return n + 1 },
		func(n int) int { return n + 1 },
		func(n int) int { return n + 1 },
		func(n int) int { return n * 10 },
	)
	if got != 60 {
		t.Fatalf("got %d", got)
	}
}

func TestPipe8(t *testing.T) {
	got := pipe.Pipe8(
		0,
		func(n int) int { return n + 1 },
		func(n int) int { return n + 1 },
		func(n int) int { return n + 1 },
		func(n int) int { return n + 1 },
		func(n int) int { return n + 1 },
		func(n int) int { return n + 1 },
		func(n int) int { return n + 1 },
		func(n int) int { return n * 10 },
	)
	if got != 70 {
		t.Fatalf("got %d", got)
	}
}

func TestComposeEndoN(t *testing.T) {
	f := pipe.ComposeEndoN(
		strings.TrimSpace,
		strings.ToUpper,
		func(s string) string { return s + "!" },
	)
	if got := f(" hello "); got != "HELLO!" {
		t.Fatalf("got %s", got)
	}
	// zero fns: identity
	id := pipe.ComposeEndoN[string]()
	if got := id("x"); got != "x" {
		t.Fatal("zero-fn ComposeEndoN should be identity")
	}
}

func TestCompose2(t *testing.T) {
	f := pipe.Compose2(strings.TrimSpace, strings.ToUpper)
	if got := f(" hello "); got != "HELLO" {
		t.Fatalf("got %s", got)
	}
}

func TestCompose3(t *testing.T) {
	f := pipe.Compose3(
		strings.TrimSpace,
		strings.ToUpper,
		func(s string) int { return len(s) },
	)
	if got := f(" hello "); got != 5 {
		t.Fatalf("got %d", got)
	}
}

func TestCompose4(t *testing.T) {
	f := pipe.Compose4(
		strings.TrimSpace,
		strings.ToUpper,
		func(s string) string { return s + "!" },
		func(s string) int { return len(s) },
	)
	if got := f(" hello "); got != 6 {
		t.Fatalf("got %d", got)
	}
}

func TestTap(t *testing.T) {
	var seen string
	got := pipe.Pipe3(
		" hello ",
		strings.TrimSpace,
		pipe.Tap(func(s string) { seen = s }),
		strings.ToUpper,
	)
	if got != "HELLO" {
		t.Fatalf("got %s", got)
	}
	if seen != "hello" {
		t.Fatalf("tap saw %q, want %q", seen, "hello")
	}
}

func TestIdentity(t *testing.T) {
	if pipe.Identity(42) != 42 {
		t.Fatal("Identity(42) != 42")
	}
	if pipe.Identity("hi") != "hi" {
		t.Fatal("Identity(\"hi\") != \"hi\"")
	}
}

func TestPredicates(t *testing.T) {
	// bool predicates
	if !pipe.IsTrue(true) || pipe.IsTrue(false) {
		t.Fatal("IsTrue")
	}
	if !pipe.IsFalse(false) || pipe.IsFalse(true) {
		t.Fatal("IsFalse")
	}

	// zero / non-zero
	if !pipe.IsZero(0) || pipe.IsZero(1) {
		t.Fatal("IsZero int")
	}
	if !pipe.IsZero("") || pipe.IsZero("x") {
		t.Fatal("IsZero string")
	}
	if pipe.IsNonZero(0) || !pipe.IsNonZero(1) {
		t.Fatal("IsNonZero")
	}

	// even / odd — integers only
	for _, n := range []int{-4, 0, 2, 100} {
		if !pipe.IsEven(n) {
			t.Fatalf("IsEven(%d)", n)
		}
	}
	for _, n := range []int{-3, 1, 7, 99} {
		if !pipe.IsOdd(n) {
			t.Fatalf("IsOdd(%d)", n)
		}
	}
	// unsigned integers
	if !pipe.IsEven(uint(4)) || pipe.IsEven(uint(3)) {
		t.Fatal("IsEven uint")
	}

	// IsPositive: integers (signed and unsigned) and floats
	if !pipe.IsPositive(1) || pipe.IsPositive(0) || pipe.IsPositive(-1) {
		t.Fatal("IsPositive int")
	}
	if !pipe.IsPositive(uint(1)) || pipe.IsPositive(uint(0)) {
		t.Fatal("IsPositive uint")
	}
	if !pipe.IsPositive(1.5) || pipe.IsPositive(0.0) || pipe.IsPositive(-1.5) {
		t.Fatal("IsPositive float64")
	}
	if !pipe.IsPositive(float32(0.1)) || pipe.IsPositive(float32(-0.1)) {
		t.Fatal("IsPositive float32")
	}

	// IsNegative: signed integers and floats only (unsigned excluded by constraint)
	if !pipe.IsNegative(-1) || pipe.IsNegative(0) || pipe.IsNegative(1) {
		t.Fatal("IsNegative int")
	}
	if !pipe.IsNegative(int8(-1)) || pipe.IsNegative(int8(0)) {
		t.Fatal("IsNegative int8")
	}
	if !pipe.IsNegative(-1.5) || pipe.IsNegative(0.0) || pipe.IsNegative(1.5) {
		t.Fatal("IsNegative float64")
	}
	if !pipe.IsNegative(float32(-0.1)) || pipe.IsNegative(float32(0.0)) {
		t.Fatal("IsNegative float32")
	}

	// Not
	isZero := pipe.Not(pipe.IsNonZero[int])
	if !isZero(0) || isZero(1) {
		t.Fatal("Not")
	}
}
