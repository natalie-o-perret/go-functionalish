package validation_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/natalie-o-perret/go-functionalish/option"
	"github.com/natalie-o-perret/go-functionalish/result"
	"github.com/natalie-o-perret/go-functionalish/validation"
)

// -- helpers ------------------------------------------------------------------

func assertSuccess[T comparable, E any](t *testing.T, v validation.Validation[T, E], want T) {
	t.Helper()
	if !v.IsSuccess() {
		t.Fatalf("expected Success(%v), got Failure(%v)", want, v.UnwrapErrors())
	}
	if v.Unwrap() != want {
		t.Fatalf("got %v, want %v", v.Unwrap(), want)
	}
}

func assertFailureN[T, E any](t *testing.T, v validation.Validation[T, E], wantN int) []E {
	t.Helper()
	if !v.IsFailure() {
		t.Fatalf("expected Failure with %d errors, got Success", wantN)
	}
	errs := v.UnwrapErrors()
	if len(errs) != wantN {
		t.Fatalf("expected %d errors, got %d: %v", wantN, len(errs), errs)
	}
	return errs
}

// -- constructors -------------------------------------------------------------

func TestSuccess(t *testing.T) {
	v := validation.Success[string](42) // E=string, T=int inferred
	assertSuccess(t, v, 42)
}

func TestFailure(t *testing.T) {
	v := validation.Failure[int]("bad", "worse") // T=int, E=string inferred
	errs := assertFailureN(t, v, 2)
	if errs[0] != "bad" || errs[1] != "worse" {
		t.Fatalf("got %v", errs)
	}
}

func TestFailurePanicsOnEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	validation.Failure[int, string]()
}

// -- methods ------------------------------------------------------------------

func TestUnwrapOr(t *testing.T) {
	s := validation.Success[string](10) // Validation[int, string]
	f := validation.Failure[int]("err")
	if s.UnwrapOr(0) != 10 {
		t.Fatal()
	}
	if f.UnwrapOr(99) != 99 {
		t.Fatal()
	}
}

func TestUnwrapOrElse(t *testing.T) {
	f := validation.Failure[int]("a", "b")
	got := f.UnwrapOrElse(func(errs []string) int { return len(errs) })
	if got != 2 {
		t.Fatalf("got %d", got)
	}
}

func TestUnwrapPanicsOnFailure(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	validation.Failure[int]("err").Unwrap()
}

func TestUnwrapErrorsPanicsOnSuccess(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	validation.Success[string](1).UnwrapErrors()
}

// -- Map / MapError / Bind ----------------------------------------------------

func TestMap(t *testing.T) {
	s := validation.Map(validation.Success[string](5), func(n int) string {
		return fmt.Sprintf("%d!", n)
	})
	assertSuccess(t, s, "5!")

	f := validation.Map(validation.Failure[int]("err"), func(n int) string {
		return "nope"
	})
	assertFailureN(t, f, 1)
}

func TestMapError(t *testing.T) {
	f := validation.MapError(
		validation.Failure[int]("bad"),
		strings.ToUpper,
	)
	errs := assertFailureN(t, f, 1)
	if errs[0] != "BAD" {
		t.Fatalf("got %v", errs)
	}

	s := validation.MapError(
		validation.Success[string](42),
		strings.ToUpper,
	)
	assertSuccess(t, s, 42)
}

func TestBind(t *testing.T) {
	half := func(n int) validation.Validation[int, string] {
		if n%2 != 0 {
			return validation.Failure[int]("odd")
		}
		return validation.Success[string](n / 2)
	}
	assertSuccess(t, validation.Bind(validation.Success[string](10), half), 5)
	assertFailureN(t, validation.Bind(validation.Success[string](7), half), 1)
	assertFailureN(t, validation.Bind(validation.Failure[int]("first"), half), 1)
}

// -- Apply --------------------------------------------------------------------

func TestApply(t *testing.T) {
	fnV := validation.Success[string](func(n int) int { return n * 2 })
	valV := validation.Success[string](21)
	assertSuccess(t, validation.Apply(fnV, valV), 42)
}

func TestApplyAccumulatesErrors(t *testing.T) {
	fnF := validation.Failure[func(int) int]("fn error")
	valF := validation.Failure[int]("val error")
	errs := assertFailureN(t, validation.Apply(fnF, valF), 2)
	if errs[0] != "fn error" || errs[1] != "val error" {
		t.Fatalf("got %v", errs)
	}
}

// -- Map2 - Map5 --------------------------------------------------------------

func TestMap2Success(t *testing.T) {
	v := validation.Map2(
		validation.Success[string]("Alice"),
		validation.Success[string](30),
		func(name string, age int) string { return fmt.Sprintf("%s:%d", name, age) },
	)
	assertSuccess(t, v, "Alice:30")
}

func TestMap2AccumulatesErrors(t *testing.T) {
	v := validation.Map2(
		validation.Failure[string]("name required"),
		validation.Failure[int]("age invalid"),
		func(name string, age int) string { return "" },
	)
	errs := assertFailureN(t, v, 2)
	if errs[0] != "name required" || errs[1] != "age invalid" {
		t.Fatalf("got %v", errs)
	}
}

func TestMap3(t *testing.T) {
	type User struct {
		Name, Email string
		Age         int
	}

	v := validation.Map3(
		validation.Failure[string]("name required"),
		validation.Failure[int]("age must be positive"),
		validation.Failure[string]("email invalid"),
		func(n string, a int, e string) User { return User{Name: n, Email: e, Age: a} },
	)
	assertFailureN(t, v, 3)

	v2 := validation.Map3(
		validation.Success[string]("Alice"),
		validation.Success[string](30),
		validation.Success[string]("alice@example.com"),
		func(n string, a int, e string) User { return User{Name: n, Email: e, Age: a} },
	)
	if !v2.IsSuccess() {
		t.Fatal("expected success")
	}
	u := v2.Unwrap()
	if u.Name != "Alice" || u.Age != 30 || u.Email != "alice@example.com" {
		t.Fatalf("got %+v", u)
	}
}

func TestMap4(t *testing.T) {
	v := validation.Map4(
		validation.Success[string](1),
		validation.Failure[int]("b"),
		validation.Success[string](3),
		validation.Failure[int]("d"),
		func(a, b, c, d int) int { return a + b + c + d },
	)
	assertFailureN(t, v, 2) // only b and d fail
}

func TestMap5(t *testing.T) {
	v := validation.Map5(
		validation.Failure[int]("a"),
		validation.Failure[int]("b"),
		validation.Failure[int]("c"),
		validation.Failure[int]("d"),
		validation.Failure[int]("e"),
		func(a, b, c, d, e int) int { return 0 },
	)
	assertFailureN(t, v, 5)
}

// -- Sequence -----------------------------------------------------------------

func TestSequenceAllSuccess(t *testing.T) {
	vs := []validation.Validation[int, string]{
		validation.Success[string](1),
		validation.Success[string](2),
		validation.Success[string](3),
	}
	v := validation.Sequence(vs)
	if !v.IsSuccess() {
		t.Fatal("expected success")
	}
	got := v.Unwrap()
	if len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
		t.Fatalf("got %v", got)
	}
}

func TestSequenceAccumulatesErrors(t *testing.T) {
	vs := []validation.Validation[int, string]{
		validation.Success[string](1),
		validation.Failure[int]("err1"),
		validation.Success[string](3),
		validation.Failure[int]("err2", "err3"),
	}
	errs := assertFailureN(t, validation.Sequence(vs), 3)
	if errs[0] != "err1" || errs[1] != "err2" || errs[2] != "err3" {
		t.Fatalf("got %v", errs)
	}
}

// -- Traverse -----------------------------------------------------------------

func TestTraverseAllValid(t *testing.T) {
	v := validation.Traverse([]string{"1", "2", "3"}, func(s string) validation.Validation[int, string] {
		n, err := strconv.Atoi(s)
		if err != nil {
			return validation.Failure[int](fmt.Sprintf("not a number: %s", s))
		}
		return validation.Success[string](n)
	})
	if !v.IsSuccess() {
		t.Fatalf("expected success, got %v", v.UnwrapErrors())
	}
	got := v.Unwrap()
	if len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
		t.Fatalf("got %v", got)
	}
}

func TestTraverseAccumulatesErrors(t *testing.T) {
	v := validation.Traverse([]string{"1", "bad", "3", "nope"}, func(s string) validation.Validation[int, string] {
		n, err := strconv.Atoi(s)
		if err != nil {
			return validation.Failure[int](fmt.Sprintf("not a number: %s", s))
		}
		return validation.Success[string](n)
	})
	assertFailureN(t, v, 2)
}

// -- interop ------------------------------------------------------------------

func TestFromResult(t *testing.T) {
	ok := validation.FromResult(result.Ok[int, string](42))
	assertSuccess(t, ok, 42)

	fail := validation.FromResult(result.Err[int, string]("boom"))
	errs := assertFailureN(t, fail, 1)
	if errs[0] != "boom" {
		t.Fatalf("got %v", errs)
	}
}

func TestFromOption(t *testing.T) {
	some := validation.FromOption(option.Some(42), "missing")
	assertSuccess(t, some, 42)

	none := validation.FromOption(option.None[int](), "missing")
	errs := assertFailureN(t, none, 1)
	if errs[0] != "missing" {
		t.Fatalf("got %v", errs)
	}
}

func TestToResult(t *testing.T) {
	s := validation.Success[string](42).ToResult()
	if !s.IsOk() || s.Unwrap() != 42 {
		t.Fatal()
	}

	f := validation.Failure[int]("a", "b").ToResult()
	if !f.IsErr() {
		t.Fatal()
	}
	errs := f.UnwrapErr()
	if len(errs) != 2 {
		t.Fatalf("got %v", errs)
	}
}

func TestToOption(t *testing.T) {
	s := validation.Success[string](42).ToOption()
	if s.IsNone() || s.Unwrap() != 42 {
		t.Fatal()
	}

	f := validation.Failure[int]("err").ToOption()
	if f.IsSome() {
		t.Fatal()
	}
}

// -- realistic integration test -----------------------------------------------

func TestFormValidation(t *testing.T) {
	type RegistrationForm struct {
		Name  string
		Email string
		Age   int
	}

	validateName := func(name string) validation.Validation[string, string] {
		if strings.TrimSpace(name) == "" {
			return validation.Failure[string]("name is required")
		}
		if len(name) < 2 {
			return validation.Failure[string]("name must be at least 2 characters")
		}
		return validation.Success[string](name)
	}

	validateEmail := func(email string) validation.Validation[string, string] {
		if !strings.Contains(email, "@") {
			return validation.Failure[string]("email must contain @")
		}
		return validation.Success[string](email)
	}

	validateAge := func(age int) validation.Validation[int, string] {
		if age < 18 {
			return validation.Failure[int]("must be 18 or older")
		}
		if age > 150 {
			return validation.Failure[int]("age seems unrealistic")
		}
		return validation.Success[string](age)
	}

	// All invalid
	v1 := validation.Map3(
		validateName(""),
		validateEmail("not-an-email"),
		validateAge(12),
		func(name, email string, age int) RegistrationForm {
			return RegistrationForm{name, email, age}
		},
	)
	errs := assertFailureN(t, v1, 3)
	if errs[0] != "name is required" ||
		errs[1] != "email must contain @" ||
		errs[2] != "must be 18 or older" {
		t.Fatalf("got %v", errs)
	}

	// All valid
	v2 := validation.Map3(
		validateName("Alice"),
		validateEmail("alice@example.com"),
		validateAge(30),
		func(name, email string, age int) RegistrationForm {
			return RegistrationForm{name, email, age}
		},
	)
	if !v2.IsSuccess() {
		t.Fatalf("expected success, got %v", v2.UnwrapErrors())
	}
	form := v2.Unwrap()
	if form.Name != "Alice" || form.Email != "alice@example.com" || form.Age != 30 {
		t.Fatalf("got %+v", form)
	}
}
