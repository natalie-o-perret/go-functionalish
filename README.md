# gof — Go Functional

A cohesive, type-safe functional programming library for Go 1.23+.  
Inspired by **F# sequences, Option, Result** and **C# LINQ**.  
No reflection. No `interface{}`. Pure generics.

## Packages

| Package | Description |
|---|---|
| `enum` | Lazy `Enumerable[T]` — LINQ-style sequence pipelines |
| `option` | `Option[T]` — explicit presence/absence, no nil |
| `result` | `Result[T,E]` — railway-oriented error handling |
| `pipe` | `Pipe2`…`Pipe8` — F#-style `\|>` operator equivalent |

## Quick start

```go
import (
    "github.com/natalie/gof/enum"
    "github.com/natalie/gof/option"
    "github.com/natalie/gof/result"
    "github.com/natalie/gof/pipe"
)
```

### enum — lazy sequences

```go
type Car struct { Year int; Owner, Model string }

cars := []Car{
    {2012, "Alice", "Toyota"}, {2016, "Bob", "Honda"},
    {2018, "Charlie", "Ford"}, {2015, "Diana", "BMW"},
}

// Filter + Map (Map is pkg-level due to Go type system)
owners := enum.Map(
    enum.From(cars).Filter(func(c Car) bool { return c.Year >= 2015 }),
    func(c Car) string { return c.Owner },
).ToSlice()
// → ["Bob", "Charlie", "Diana"]

// Sort, group, distinct
byCar := enum.GroupBy(enum.From(cars), func(c Car) string { return c.Model })

unique := enum.UniqBy(
    enum.From(cars).SortBy(func(a, b Car) int { return cmp.Compare(a.Model, b.Model) }),
    func(c Car) string { return c.Model },
).ToSlice()

// Short-circuiting terminals
first := enum.From(cars).Filter(...).FirstOption() // → option.Option[Car]
count := enum.From(cars).CountBy(func(c Car) bool { return c.Year >= 2015 })

// Generators
squares := enum.Map(enum.Range(1, 6), func(n int) int { return n * n }).ToSlice()
// → [1 4 9 16 25]

// Zip
pairs := enum.Zip(enum.From([]int{1, 2, 3}), enum.From([]string{"a", "b", "c"})).ToSlice()
// → [{1 a} {2 b} {3 c}]

// Cycle + Take (infinite sequences)
pattern := enum.From([]string{"ping", "pong"}).Cycle().Take(5).ToSlice()
// → ["ping", "pong", "ping", "pong", "ping"]
```

### option — explicit optionality

```go
// Instead of (T, bool) or *T
name := option.Some("Alice")
none := option.None[string]()

upper := option.Map(name, strings.ToUpper) // → Some("ALICE")
option.Map(none, strings.ToUpper)           // → None

name.UnwrapOr("anonymous") // → "Alice"
none.UnwrapOr("anonymous") // → "anonymous"

// Chain optional operations (FlatMap flattens)
result := option.FlatMap(findUser(id), func(u User) option.Option[Profile] {
    return findProfile(u.ProfileID)
})

// Integrates with enum
firstModern := enum.From(cars).
    Filter(func(c Car) bool { return c.Year >= 2015 }).
    FirstOption() // → option.Option[Car]
```

### result — railway-oriented error handling

```go
// Instead of (T, error) with checks everywhere
res := result.Try(func() (User, error) { return db.FindUser(id) })

// Chain operations — Err short-circuits, no explicit checks
final := result.FlatMap(
    result.Map(res, enrichUser),
    func(u User) result.Result[Response, error] { return serialize(u) },
)

if final.IsOk() {
    fmt.Println(final.Unwrap())
} else {
    fmt.Println("error:", final.UnwrapErr())
}

// Interop with option
opt := res.ToOption()                           // Ok → Some, Err → None
res2 := result.FromOption(opt, errors.New("not found"))
```

### pipe — threading values

```go
// F#: input |> trim |> toLower |> validate
validated := pipe.Pipe3(
    rawInput,
    strings.TrimSpace,
    strings.ToLower,
    validate,
)
```

## Design notes

### Why are `Map`, `FlatMap`, `GroupBy` package-level functions?

Go does not allow methods to introduce new type parameters. A method on
`Enumerable[Car]` cannot return `Enumerable[string]` because that would require
`func (e Enumerable[T]) Map[R any](fn func(T) R) Enumerable[R]` — which the
compiler rejects.

The workaround: **type-transforming operations are package-level functions**,
same-type operations are methods:

```go
// ✅ method — stays Enumerable[Car]
From(cars).Filter(fn).SortBy(less).Take(10)

// ✅ package-level — changes type
enum.Map(From(cars).Filter(fn), Car.Owner)
//       ↑ sub-chain            ↑ transform
```

### Lazy vs eager

All pipeline operations (`Filter`, `Map`, `Take`, …) are **lazy** — they wrap
the previous iterator and produce no output until a terminal is called. Only
`SortBy`, `SortAscBy`, `Reverse` must **materialise** (you can't sort
a stream you haven't fully read).

## Dependency graph

```
pipe    →  ∅
option  →  ∅
result  →  option
enum    →  option, result, go-functional/v2
```

No circular imports. `go-functional/v2` is the only external dependency.

