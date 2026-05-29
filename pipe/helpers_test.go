package pipe_test

import (
	"cmp"
	"strings"
	"testing"

	"github.com/natalie-o-perret/go-functionalish/pipe"
)

// -- arithmetic ---------------------------------------------------------------

func TestIncDec(t *testing.T) {
	if pipe.Inc(4) != 5 {
		t.Fatal("Inc int")
	}
	if pipe.Inc(4.5) != 5.5 {
		t.Fatal("Inc float64")
	}
	if pipe.Dec(4) != 3 {
		t.Fatal("Dec int")
	}
	if pipe.Dec(uint(1)) != 0 {
		t.Fatal("Dec uint")
	}
}

func TestAddSubMulDiv(t *testing.T) {
	add10 := pipe.Add(10)
	if add10(5) != 15 {
		t.Fatal("Add")
	}
	sub3 := pipe.Sub(3)
	if sub3(10) != 7 {
		t.Fatal("Sub")
	}
	double := pipe.Mul(2)
	if double(6) != 12 {
		t.Fatal("Mul")
	}
	half := pipe.Div(2.0)
	if half(10.0) != 5.0 {
		t.Fatal("Div float64")
	}
	// integer Div
	divBy3 := pipe.Div(3)
	if divBy3(9) != 3 {
		t.Fatal("Div int")
	}
}

func TestNegate(t *testing.T) {
	if pipe.Negate(-5) != 5 {
		t.Fatal("Negate int")
	}
	if pipe.Negate(float64(3.14)) != -3.14 {
		t.Fatal("Negate float64")
	}
}

func TestAbs(t *testing.T) {
	if pipe.Abs(-7) != 7 {
		t.Fatal("Abs negative int")
	}
	if pipe.Abs(7) != 7 {
		t.Fatal("Abs positive int")
	}
	if pipe.Abs(0) != 0 {
		t.Fatal("Abs zero")
	}
	if pipe.Abs(float64(-3.14)) != 3.14 {
		t.Fatal("Abs float64")
	}
}

func TestClamp(t *testing.T) {
	clamp := pipe.Clamp(0, 100)
	if clamp(-5) != 0 {
		t.Fatalf("Clamp below lo: got %d", clamp(-5))
	}
	if clamp(50) != 50 {
		t.Fatalf("Clamp in range: got %d", clamp(50))
	}
	if clamp(150) != 100 {
		t.Fatalf("Clamp above hi: got %d", clamp(150))
	}
	// float
	clampF := pipe.Clamp(0.0, 1.0)
	if clampF(1.5) != 1.0 {
		t.Fatal("Clamp float64 hi")
	}
	if clampF(-0.5) != 0.0 {
		t.Fatal("Clamp float64 lo")
	}
}

// -- predicate combinators ----------------------------------------------------

func TestAnd(t *testing.T) {
	isSmallEven := pipe.And(pipe.IsEven[int], func(n int) bool { return n < 10 })
	if !isSmallEven(4) {
		t.Fatal("And: 4 should match")
	}
	if isSmallEven(12) {
		t.Fatal("And: 12 too large")
	}
	if isSmallEven(3) {
		t.Fatal("And: 3 not even")
	}
	// vacuous truth
	if !pipe.And[int]()(42) {
		t.Fatal("And() should be vacuously true")
	}
}

func TestOr(t *testing.T) {
	isZeroOrHundred := pipe.Or(pipe.IsZero[int], func(n int) bool { return n == 100 })
	if !isZeroOrHundred(0) {
		t.Fatal("Or: 0 should match")
	}
	if !isZeroOrHundred(100) {
		t.Fatal("Or: 100 should match")
	}
	if isZeroOrHundred(42) {
		t.Fatal("Or: 42 should not match")
	}
	// vacuous false
	if pipe.Or[int]()(42) {
		t.Fatal("Or() should be vacuously false")
	}
}

// -- general combinators ------------------------------------------------------

func TestConst(t *testing.T) {
	always42 := pipe.Const[int, string](42)
	if always42("anything") != 42 {
		t.Fatal("Const")
	}
	if always42("") != 42 {
		t.Fatal("Const empty")
	}
}

func TestFlip(t *testing.T) {
	divide := func(a, b float64) float64 { return a / b }
	divideInto := pipe.Flip(divide)
	if divideInto(2.0, 10.0) != 5.0 {
		t.Fatalf("Flip: got %v", divideInto(2.0, 10.0))
	}
	// flip strings.HasPrefix: (prefix, s) -> (s, prefix)
	flipped := pipe.Flip(strings.HasPrefix)
	if !flipped("go", "golang") {
		t.Fatal("Flip HasPrefix")
	}
}

func TestOn(t *testing.T) {
	byLen := pipe.On(cmp.Compare, func(s string) int { return len(s) })
	if byLen("hi", "hello") >= 0 {
		t.Fatal("On: 'hi' should sort before 'hello'")
	}
	if byLen("hello", "hi") <= 0 {
		t.Fatal("On: 'hello' should sort after 'hi'")
	}
	if byLen("ab", "cd") != 0 {
		t.Fatal("On: equal-length strings should compare as 0")
	}
}

func TestCurry2Uncurry2(t *testing.T) {
	add := func(a, b int) int { return a + b }
	curriedAdd := pipe.Curry2(add)
	add5 := curriedAdd(5)
	if add5(3) != 8 {
		t.Fatal("Curry2")
	}
	uncurriedAdd := pipe.Uncurry2(curriedAdd)
	if uncurriedAdd(5, 3) != 8 {
		t.Fatal("Uncurry2")
	}
	// strings.HasPrefix(s, prefix) — curry fixes s first, then prefix
	// Use Flip to get HasPrefix(prefix, s) for the natural "fix the prefix" pattern
	hasPrefix := pipe.Curry2(pipe.Flip(strings.HasPrefix))
	startsWithGo := hasPrefix("go")
	if !startsWithGo("golang") || startsWithGo("python") {
		t.Fatal("Curry2 HasPrefix")
	}
}

func TestPartial1Partial2(t *testing.T) {
	sub := func(a, b int) int { return a - b }
	sub10From := pipe.Partial1(sub, 10)
	if sub10From(3) != 7 {
		t.Fatalf("Partial1: got %d", sub10From(3))
	}
	subFrom10 := pipe.Partial2(sub, 10)
	if subFrom10(3) != -7 {
		t.Fatalf("Partial2: got %d", subFrom10(3))
	}

	contains := pipe.Partial2(strings.Contains, "go")
	if !contains("golang") || contains("python") {
		t.Fatal("Partial2 strings.Contains")
	}
}
