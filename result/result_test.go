package result_test

import (
	"errors"
	"testing"

	"github.com/natalie-o-perret/gof/result"
)

func TestOk(t *testing.T) {
	r := result.Ok[int, error](42)
	if !r.IsOk() { t.Fatal("expected Ok") }
	if r.Unwrap() != 42 { t.Fatalf("got %d", r.Unwrap()) }
}

func TestErr(t *testing.T) {
	r := result.Err[int, error](errors.New("boom"))
	if !r.IsErr() { t.Fatal("expected Err") }
	if r.UnwrapOr(99) != 99 { t.Fatal("expected default") }
}

func TestTry(t *testing.T) {
	ok := result.Try(func() (int, error) { return 1, nil })
	if !ok.IsOk() { t.Fatal("expected Ok") }

	fail := result.Try(func() (int, error) { return 0, errors.New("bad") })
	if !fail.IsErr() { t.Fatal("expected Err") }
}

func TestMap(t *testing.T) {
	r := result.Map(result.Ok[int, error](3), func(n int) int { return n * 2 })
	if r.Unwrap() != 6 { t.Fatalf("got %d", r.Unwrap()) }
}

func TestMapErr(t *testing.T) {
	r := result.MapErr(
		result.Err[int, string]("oops"),
		func(s string) int { return len(s) },
	)
	if r.UnwrapErr() != 4 { t.Fatalf("got %d", r.UnwrapErr()) }
}

func TestFlatMap(t *testing.T) {
	double := func(n int) result.Result[int, error] { return result.Ok[int, error](n * 2) }
	got := result.FlatMap(result.Ok[int, error](5), double).Unwrap()
	if got != 10 { t.Fatalf("got %d", got) }
}

func TestToOption(t *testing.T) {
	opt := result.Ok[int, error](7).ToOption()
	if !opt.IsSome() || opt.Unwrap() != 7 { t.Fatal("expected Some(7)") }
	none := result.Err[int, error](errors.New("x")).ToOption()
	if !none.IsNone() { t.Fatal("expected None") }
}

