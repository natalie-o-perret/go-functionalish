# gof: Go Functional

[![CI](https://github.com/natalie-o-perret/gof/actions/workflows/ci.yml/badge.svg)](https://github.com/natalie-o-perret/gof/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/natalie-o-perret/gof.svg)](https://pkg.go.dev/github.com/natalie-o-perret/gof)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Contributing](https://img.shields.io/badge/contributing-guide-blue)](CONTRIBUTING.md)

A cohesive, opinionated, type-safe functional programming library for Go 1.24+.  
Inspired by C# LINQ, F# sequences, Option, Result.  
No reflection. No `interface{}`. Pure generics and lazy by default.

> [!NOTE]
> Unapologetically lamely vibe-coded with excruciatingly expensive Claude Opus 4.6.
> 
> Unapologetically not "idiomatic Go."
>
> Go gave us (at last) generics 17 years after C# and 18 after Java (the  latter still erases them at runtime).
> We're using them, for `Option[T]`, `Result[T,E]`, and lazy pipelines
> instead of `if err != nil` sixty times per file.

## Packages

| Package      | Description                                           |
|--------------|-------------------------------------------------------|
| `seq`        | Lazy `Seq[T]`: F#-style sequence pipelines            |
| `option`     | `Option[T]`: explicit presence/absence, no nil        |
| `result`     | `Result[T,E]`: railway-oriented error handling        |
| `validation` | `Validation[T,E]`: applicative error accumulation     |
| `pipe`       | `Pipe2`...`Pipe8`: F#-style `\|>` operator equivalent |

## Quick start

```go
import (
    "github.com/natalie-o-perret/gof/seq"
    "github.com/natalie-o-perret/gof/option"
    "github.com/natalie-o-perret/gof/result"
    "github.com/natalie-o-perret/gof/validation"
    "github.com/natalie-o-perret/gof/pipe"
)
```

### seq: lazy sequences

```go
type Car struct { Year int; Owner, Model string }

cars := []Car{
    {2012, "Alice", "Toyota"}, {2016, "Bob", "Honda"},
    {2018, "Charlie", "Ford"}, {2015, "Diana", "BMW"},
}

// Filter + Map (Map is pkg-level due to Go type system)
owners := seq.Map(
    seq.OfSlice(cars).Filter(func(c Car) bool { return c.Year >= 2015 }),
    func(c Car) string { return c.Owner },
).ToSlice()
// => ["Bob", "Charlie", "Diana"]

// Sort, group, distinct
byCar := seq.GroupBy(seq.OfSlice(cars), func(c Car) string { return c.Model })

unique := seq.DistinctBy(
    seq.OfSlice(cars).SortWith(func(a, b Car) int { return cmp.Compare(a.Model, b.Model) }),
    func(c Car) string { return c.Model },
).ToSlice()

// Short-circuiting terminals
first := seq.OfSlice(cars).Filter(...).TryHead() // => option.Option[Car]
count := seq.OfSlice(cars).CountBy(func(c Car) bool { return c.Year >= 2015 })

// Generators
squares := seq.Map(seq.Range(1, 6), func(n int) int { return n * n }).ToSlice()
// => [1 4 9 16 25]

// Zip
pairs := seq.Zip(seq.OfSlice([]int{1, 2, 3}), seq.OfSlice([]string{"a", "b", "c"})).ToSlice()
// => [{1 a} {2 b} {3 c}]

// Cycle + Truncate (infinite sequences)
pattern := seq.OfSlice([]string{"ping", "pong"}).Cycle().Truncate(5).ToSlice()
// => ["ping", "pong", "ping", "pong", "ping"]
```

### option: explicit optionality

```go
// Instead of (T, bool) or *T
name := option.Some("Alice")
none := option.None[string]()

upper := option.Map(name, strings.ToUpper) // => Some("ALICE")
option.Map(none, strings.ToUpper)           // => None

name.UnwrapOr("anonymous") // => "Alice"
none.UnwrapOr("anonymous") // => "anonymous"

// Chain optional lookups with Bind (F# Option.bind)
profile := option.Bind(findUser(id), func(u User) option.Option[Profile] {
    return findProfile(u.ProfileID)
})

// Integrates with seq
firstModern := seq.OfSlice(cars).
    Filter(func(c Car) bool { return c.Year >= 2015 }).
    TryHead() // => option.Option[Car]
```

### result : railway-oriented error handling

```go
// Wrap Go's (T, error) convention
res := result.Try(func() (User, error) { return db.FindUser(id) })

// Railway pipeline: chain Bind (may fail) and Map (pure transforms).
// Once on the error track, every subsequent step is skipped.
final := result.Map(
    result.Bind(
        result.Map(res, enrichUser),           // Map: pure transform
        func(u User) result.Result[Response, error] {
            return serialize(u)                // Bind: may fail
        },
    ),
    addHeaders,                                // Map: pure transform
)

// Or write it as a linear pipeline (reads top-to-bottom):
r1 := parseRequest(raw)                       // step 1: Bind
r2 := result.Bind(r1, authenticate)           // step 2: Bind
r3 := result.Map(r2, normalize)               // step 3: Map (pure)
r4 := result.Bind(r3, save)                   // step 4: Bind
r5 := result.Map(r4, formatResponse)          // step 5: Map (pure)
// r5 is either Ok(response) or Err from whichever step failed first.

// MapErr adds context to errors
wrapped := result.MapErr(r5, func(e string) string {
    return "request failed: " + e
})

// Interop with option
opt := res.ToOption()                          // Ok => Some, Err => None
res2 := result.FromOption(opt, errors.New("not found"))
```

### validation : applicative error accumulation

```go
// Unlike Result which short-circuits, Validation runs ALL checks
// and collects every error.

validateName := func(name string) validation.Validation[string, string] {
    if name == "" {
        return validation.Failure[string, string]("name is required")
    }
    return validation.Success[string, string](name)
}

validateAge := func(age int) validation.Validation[int, string] {
    if age < 18 {
        return validation.Failure[int, string]("must be 18 or older")
    }
    return validation.Success[int, string](age)
}

// Tip: the trailing type arg can be inferred from the value, so these are equivalent:
//   validation.Failure[int, string]("bad")  ==  validation.Failure[int]("bad")
//   validation.Success[int, string](42)     ==  validation.Success[string](42)

// Map2-Map5: combine N validations, accumulating ALL errors
user := validation.Map2(
    validateName(""),       // Failure
    validateAge(12),        // Failure
    func(name string, age int) User { return User{name, age} },
)
// => Failure(["name is required", "must be 18 or older"])
// Both errors reported, nothing short-circuited!

// Sequence: []Validation => Validation[[]T]
results := validation.Sequence([]validation.Validation[int, string]{
    validation.Success[int, string](1),
    validation.Failure[int, string]("bad"),
    validation.Failure[int, string]("worse"),
})
// => Failure(["bad", "worse"])

// Traverse: map + sequence in one step
parsed := validation.Traverse([]string{"1", "bad", "3"}, parseIntV)
// => Failure(["not a number: bad"])

// Interop: Result <=> Validation
v := validation.FromResult(result.Ok[int, string](42))  // => Success(42)
r := v.ToResult()                                         // => Ok(42)
```

### pipe : threading values

```go
// F#: input |> trim |> toLower |> validate
validated := pipe.Pipe3(
    rawInput,
    strings.TrimSpace,
    strings.ToLower,
    validate,
)

// PipeN: variadic, unlimited steps - but all functions must share the same type T→T.
// Use this when every step is a transform within the same type.
processed := pipe.PipeN(
    rawInput,
    strings.TrimSpace,
    strings.ToLower,
    strings.Title,
    sanitize,
    normalize,
)
```

#### Pipelines where the type changes between steps

When steps change the type (e.g. `string` to `int` to `MyStruct`), Go's type system
requires each type to be a distinct type parameter - so variadic isn't possible.
`Pipe2`-`Pipe8` handle up to 8 type-changing steps. To go beyond 8, use
`Compose`/`Compose2`-`Compose4` to collapse multiple steps into one slot:

```go
// Nested Then (works, but awkward):
seq.Then(
    seq.Then(
        seq.OfSlice(people).Filter(adult),
        seq.MapFn(getName),
    ),
    seq.DistinctFn[string](),
).SortWith(cmp.Compare).Truncate(3).ToSlice()

// Flat pipe (same result):
pipe.Pipe5(
    seq.OfSlice(people),
    seq.FilterFn(adult),                       // Seq[Person] → Seq[Person]
    seq.MapFn(getName),                        // Seq[Person] → Seq[string]
    seq.DistinctFn[string](),                  // Seq[string] → Seq[string]
    pipe.Compose(                              // group same-type steps
        seq.SortWithFn[string](cmp.Compare),
        seq.TruncateFn[string](3),
    ),
    seq.ToSliceFn[string](),                   // Seq[string] → []string
)
```

`pipe.Compose` merges consecutive same-type steps (`T→T`) into a single pipe
slot. `pipe.Compose2`-`Compose4` do the same for type-changing steps (`A→B→C`).
Together with `PipeN` for same-type chains, this covers pipelines of any length.

## Design notes

### Why are `Map`, `Collect`, `GroupBy` package-level functions?

Go does not allow methods to introduce new type parameters. A method on
`Seq[Car]` cannot return `Seq[string]` because that would require
`func (s Seq[T]) Map[R any](fn func(T) R) Seq[R]` : which the
compiler rejects.

The workaround: **type-transforming operations are package-level functions**,
same-type operations are methods:

```go
//  method : stays Seq[Car]
OfSlice(cars).Filter(fn).SortWith(less).Truncate(10)

//  package-level : changes type
seq.Map(OfSlice(cars).Filter(fn), Car.Owner)
//      ↑ sub-chain               ↑ transform
```

### Lazy vs eager

All pipeline operations (`Filter`, `Map`, `Truncate`, ...) are **lazy** : they wrap
the previous iterator and produce no output until a terminal is called. Only
`SortWith`, `SortBy`, `Rev` must **materialise** (you can't sort
a stream you haven't fully read).

## Performance

### Direct method chaining vs pipe + compose

The `*Fn` curried helpers and `pipe.Compose` add thin closure wrappers around
the same underlying methods. Benchmarks on a 10,000-element `[]int` pipeline
(Intel Core Ultra 7):

| Pipeline                              | Style  |   ns/op |   B/op | allocs |
|---------------------------------------|--------|--------:|-------:|-------:|
| **Small** (filter => map => take 100) | Direct |  ~1,780 |  2,064 |      9 |
|                                       | Pipe   |  ~1,620 |  2,064 |      9 |
| **Medium** (7 steps incl. sort)       | Direct | ~13,100 |  8,376 |     26 |
|                                       | Pipe   | ~13,900 |  8,912 |     44 |
| **Large** (10 steps incl. rev+sort)   | Direct | ~27,100 | 12,784 |     46 |
|                                       | Pipe   | ~30,400 | 13,592 |     73 |

**~6-13 % wall-clock overhead**, all from one-time closure allocations when the
pipeline is *built*, not per element. The hot iteration loop is identical either
way. For any real workload (I/O, serialisation, business logic in the lambdas)
this is noise: choose whichever style reads better.

Run the benchmarks yourself:

```sh
go test ./seq/ -bench=. -benchmem
```

## Dependency graph

```
pipe       =>  (none)
option     =>  (none)
result     =>  option
validation =>  option, result
seq        =>  option
```

No circular imports. No external dependencies.
