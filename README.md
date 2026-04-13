# gof: Go Functional

[![CI](https://github.com/natalie-o-perret/gof/actions/workflows/ci.yml/badge.svg)](https://github.com/natalie-o-perret/gof/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/natalie-o-perret/gof.svg)](https://pkg.go.dev/github.com/natalie-o-perret/gof)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Contributing](https://img.shields.io/badge/contributing-guide-blue)](CONTRIBUTING.md)

A cohesive, opinionated, type-safe functional programming library for Go 1.24+.  
No reflection. No `interface{}`. Pure generics and lazy by default.

> [!NOTE]
> Unapologetically vibe-coded with Claude Opus 4.6.
>
> Unapologetically not "idiomatic Go."
>
> Go gave us generics 17 years after C# and 18 after Java (the latter still erases them at runtime).
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
| `tuple`      | `T2`/`T3`/`T4`: immutable generic tuples              |

## Quick start

```go
import (
    "github.com/natalie-o-perret/gof/seq"
    "github.com/natalie-o-perret/gof/option"
    "github.com/natalie-o-perret/gof/result"
    "github.com/natalie-o-perret/gof/validation"
    "github.com/natalie-o-perret/gof/pipe"
    "github.com/natalie-o-perret/gof/tuple"
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

// RangeStep: custom step, supports descending
evens := seq.RangeStep(0, 10, 2).ToSlice()  // => [0 2 4 6 8]
countdown := seq.RangeStep(5, 0, -1).ToSlice() // => [5 4 3 2 1]

// Zip / Interleave
pairs := seq.Zip(seq.OfSlice([]int{1, 2, 3}), seq.OfSlice([]string{"a", "b", "c"})).ToSlice()
// => [{1 a} {2 b} {3 c}]

merged := seq.Interleave(seq.OfSlice([]int{1, 3, 5}), seq.OfSlice([]int{2, 4, 6})).ToSlice()
// => [1 2 3 4 5 6]

// Cycle + Truncate (infinite sequences)
pattern := seq.OfSlice([]string{"ping", "pong"}).Cycle().Truncate(5).ToSlice()
// => ["ping", "pong", "ping", "pong", "ping"]

// Unfold: generate from a seed state (e.g. Fibonacci)
fibs := seq.Unfold([2]int{0, 1}, func(s [2]int) option.Option[seq.Pair[int, [2]int]] {
    if s[0] > 20 {
        return option.None[seq.Pair[int, [2]int]]()
    }
    return option.Some(seq.Pair[int, [2]int]{First: s[0], Second: [2]int{s[1], s[0] + s[1]}})
}).ToSlice()
// => [0 1 1 2 3 5 8 13]

// Partition: one pass, two slices
evens, odds := seq.Partition(seq.Range(1, 7), func(n int) bool { return n%2 == 0 })
// evens => [2 4 6],  odds => [1 3 5]

// OfMap: iterate over a map
counts := seq.CountByKey(seq.OfMap(map[string]int{"a": 1, "b": 2}),
    func(p seq.Pair[string, int]) string { return p.First })
// => map[a:1 b:1]

// CountByKey: occurrence counts
freq := seq.CountByKey(seq.OfSlice([]string{"a", "b", "a", "c", "a", "b"}),
    func(s string) string { return s })
// => map[a:3 b:2 c:1]

// OfOption: lift an Option into a Seq
seq.OfOption(option.Some(42)).ToSlice()   // => [42]
seq.OfOption(option.None[int]()).ToSlice() // => []
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

// DefaultWith: lazy default - fn is only called when None
val := option.DefaultWith(none, func() string { return expensiveDefault() })

// Contains: value equality check
option.Contains(option.Some(42), 42)    // => true
option.Contains(option.None[int](), 42) // => false

// Map2: combine two Options
option.Map2(option.Some(2), option.Some(3), func(a, b int) int { return a + b }) // => Some(5)
option.Map2(option.None[int](), option.Some(3), func(a, b int) int { return a + b }) // => None

// OrElse: fallback if None
resolved := option.OrElse(lookupCache(key), func() option.Option[string] {
    return lookupDB(key)
})

// Flatten: unwrap Option[Option[T]]
option.Flatten(option.Some(option.Some(42)))          // => Some(42)
option.Flatten(option.None[option.Option[int]]())     // => None

// Chain optional lookups with Bind
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
r1 := parseRequest(raw)                       // step 1: Bind
r2 := result.Bind(r1, authenticate)           // step 2: Bind
r3 := result.Map(r2, normalize)               // step 3: Map (pure)
r4 := result.Bind(r3, save)                   // step 4: Bind
r5 := result.Map(r4, formatResponse)          // step 5: Map (pure)
// r5 is either Ok(response) or Err from whichever step failed first.

// Tee / TeeErr: side-effects without breaking the chain - great for logging
r := result.Tee(r3, func(v Request) { log.Printf("normalised: %v", v) })
r = result.TeeErr(r, func(e string) { log.Printf("failed: %s", e) })

// OrElse: try a fallback on Err
user := result.OrElse(lookupPrimary(id), func(e error) result.Result[User, error] {
    return lookupReplica(id)
})

// Flatten: unwrap Result[Result[T,E],E]
result.Flatten(result.Ok[result.Result[int, string], string](result.Ok[int, string](42)))
// => Ok(42)

// Zip: combine two Results into a pair (first Err wins)
result.Zip(result.Ok[int, string](1), result.Ok[string, string]("hi"))
// => Ok({1, "hi"})

// MapErr adds context to errors
wrapped := result.MapErr(r5, func(e string) string {
    return "request failed: " + e
})

// Interop with option
opt := res.ToOption()                          // Ok => Some, Err => None
res2 := result.FromOption(opt, errors.New("not found"))
```

### pipe : threading values

**Naming convention** - the suffix tells you the type contract:

| Suffix  | Variant                                  | Type contract                                          | Example                 |
|---------|------------------------------------------|--------------------------------------------------------|-------------------------|
| `EndoN` | `PipeEndoN`, `ComposeEndoN`              | variadic, all steps `T => T` (same type, endomorphic)  | string transforms       |
| `2..8`  | `Pipe2`-`Pipe8`, `Compose2`-`Compose4`   | fixed-arity, each step may change type (`A => B => C`) | parse + validate + save |

The compiler enforces this: `PipeEndoN` simply cannot accept a function whose output
type differs from its input. `Pipe2`-`Pipe8` each declare distinct type params
`A, B, C, ...` so the change is explicit in the signature itself.

```go
// Same-type chain (T => T): use PipeEndoN - unlimited steps, all string => string
processed := pipe.PipeEndoN(
    rawInput,
    strings.TrimSpace,
    strings.ToLower,
    strings.Title,
    sanitize,
)

// Type-changing chain: use Pipe2-Pipe8
validated := pipe.Pipe3(
    rawInput,            // string
    strings.TrimSpace,   // string => string
    parse,               // string => int
    validate,            // int    => error
)

// Tap: side-effect (logging, metrics) without changing the value
result := pipe.Pipe4(
    rawInput,
    strings.TrimSpace,
    pipe.Tap(func(s string) { log.Println("trimmed:", s) }),
    strings.ToUpper,
    validate,
)
```

#### Pipelines longer than 8 type-changing steps

Use `ComposeEndoN`/`Compose2`-`Compose4` to collapse multiple steps into one slot:

```go
pipe.Pipe5(
    seq.OfSlice(people),
    seq.FilterFn(adult),                          // Seq[Person] => Seq[Person]
    seq.MapFn(getName),                           // Seq[Person] => Seq[string]
    seq.DistinctFn[string](),                     // Seq[string] => Seq[string]
    pipe.ComposeEndoN(                                // group same-type steps into one slot
        seq.SortWithFn[string](cmp.Compare),
        seq.TruncateFn[string](3),
    ),
    seq.ToSliceFn[string](),                      // Seq[string] => []string
)
```

`pipe.ComposeEndoN` merges consecutive same-type steps (`T=>T`) into one pipe slot.
`pipe.Compose2`-`Compose4` do the same for type-changing steps (`A=>B=>C`).
Together they cover pipelines of any length.

### tuple: immutable generic tuples

```go
// Construct
pair  := tuple.Of(1, "hello")       // T2[int, string]
triple := tuple.Of3(1, "hi", true)  // T3[int, string, bool]
quad  := tuple.Of4(1, "hi", true, 3.14) // T4[int, string, bool, float64]

// Access fields
pair.First  // 1
pair.Second // "hello"

// Destructure (closest to F# let (a, b) = t)
a, b := pair.Unpack()
a, b, c := triple.Unpack()

// Swap (T2 only) - returns T2[B, A]
tuple.Of("x", 99).Swap() // T2[int, string]{First: 99, Second: "x"}

// Apply: call a function with the tuple's fields
tuple.Apply(tuple.Of(3, 4), func(a, b int) int { return a + b }) // 7
tuple.Apply3(triple, func(n int, s string, b bool) string { ... })

// Map fields
tuple.MapFirst(tuple.Of("hello", 42), strings.ToUpper)  // T2{"HELLO", 42}
tuple.MapSecond(tuple.Of(42, "hello"), strings.ToUpper) // T2{42, "HELLO"}
tuple.Map(tuple.Of("hello", "world"), strings.ToUpper, strings.ToUpper)

// Curry / Uncurry
add := func(a, b int) int { return a + b }
curried := tuple.Curry(add)   // func(int) func(int) int
curried(3)(4)                 // 7
tuple.Uncurry(curried)(3, 4)  // 7

// Lift binary functions to/from tuple-consuming form
tupledAdd := tuple.FromFunc2(add)     // func(T2[int,int]) int
tupledAdd(tuple.Of(2, 3))            // 5
tuple.ToFunc2(tupledAdd)(2, 3)       // 5
```

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
//      ^ sub-chain               ^ transform
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
tuple      =>  (none)
```
