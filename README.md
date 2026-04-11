# gof: Go Functional

A cohesive, opinionated, type-safe functional programming library for Go 1.23+.  
Inspired by C# LINQ, F# sequences, Option, Result**.  
No reflection. No `interface{}`. Pure generics and lazy by default.

> [!NOTE]
> Unapologetically lamely vibe-coded with excruciatingly expensive Claude Opus 4.6.
> Unapologetically not "idiomatic Go."
>
> Go gave us (at last) generics 17 years after C# and 18 after Java (the  latter still erases them at runtime).
> We're using them, for `Option[T]`, `Result[T,E]`, and lazy pipelines
> instead of `if err != nil` sixty times with  per file.

## Packages

| Package  | Description                                           |
|----------|-------------------------------------------------------|
| `seq`    | Lazy `Seq[T]`: F#-style sequence pipelines           |
| `option` | `Option[T]`: explicit presence/absence, no nil       |
| `result` | `Result[T,E]`: railway-oriented error handling       |
| `pipe`   | `Pipe2`...`Pipe8`: F#-style `\|>` operator equivalent |

## Quick start

```go
import (
    "github.com/natalie/gof/seq"
    "github.com/natalie/gof/option"
    "github.com/natalie/gof/result"
    "github.com/natalie/gof/pipe"
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

// Chain optional operations (FlatMap flattens)
result := option.FlatMap(findUser(id), func(u User) option.Option[Profile] {
    return findProfile(u.ProfileID)
})

// Integrates with seq
firstModern := seq.OfSlice(cars).
    Filter(func(c Car) bool { return c.Year >= 2015 }).
    TryHead() // => option.Option[Car]
```

### result : railway-oriented error handling

```go
// Instead of (T, error) with checks everywhere
res := result.Try(func() (User, error) { return db.FindUser(id) })

// Chain operations : Err short-circuits, no explicit checks
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
opt := res.ToOption()                           // Ok => Some, Err => None
res2 := result.FromOption(opt, errors.New("not found"))
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
```

#### Composing seq pipelines with pipe

Instead of nesting `seq.Then` calls, use `pipe.Pipe*` with curried `*Fn` helpers
for a flat, linear style  - even when the element type changes mid-pipeline:

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
This lets you express arbitrarily long pipelines within `Pipe2`-`Pipe8`.

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
pipe    =>  (none)
option  =>  (none)
result  =>  option
seq     =>  option
```

No circular imports. No external dependencies.
