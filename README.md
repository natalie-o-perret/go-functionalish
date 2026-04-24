# go-functionalish

[![CI](https://github.com/natalie-o-perret/go-functionalish/actions/workflows/ci.yml/badge.svg)](https://github.com/natalie-o-perret/go-functionalish/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/natalie-o-perret/go-functionalish.svg)](https://pkg.go.dev/github.com/natalie-o-perret/go-functionalish)
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

| Package      | Description                                                     |
|--------------|-----------------------------------------------------------------|
| `seq`        | Lazy `Seq[T]`: F#-style sequence pipelines                      |
| `pseq`       | Parallel `Seq[T]`: goroutine-per-chunk Map, Filter, Reduce, ... |
| `option`     | `Option[T]`: explicit presence/absence, no nil                  |
| `result`     | `Result[T,E]`: railway-oriented error handling                  |
| `validation` | `Validation[T,E]`: applicative error accumulation               |
| `pipe`       | `Pipe2`...`Pipe8`: F#-style `\|>` operator equivalent           |
| `kv`         | Lazy `Seq2[K,V]`: functional pipelines over `iter.Seq2` / maps  |

## Quick start

```go
import (
"github.com/natalie-o-perret/go-functionalish/seq"
"github.com/natalie-o-perret/go-functionalish/pseq"
"github.com/natalie-o-perret/go-functionalish/option"
"github.com/natalie-o-perret/go-functionalish/result"
"github.com/natalie-o-perret/go-functionalish/validation"
"github.com/natalie-o-perret/go-functionalish/pipe"
"github.com/natalie-o-perret/go-functionalish/kv"
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
seq.OfSlice(cars).Filter(func (c Car) bool { return c.Year >= 2015 }),
func (c Car) string { return c.Owner },
).ToSlice()
// => ["Bob", "Charlie", "Diana"]

// Sort, group, distinct
byCar := seq.GroupBy(seq.OfSlice(cars), func (c Car) string { return c.Model })

unique := seq.DistinctBy(
seq.OfSlice(cars).SortWith(func (a, b Car) int { return cmp.Compare(a.Model, b.Model) }),
func (c Car) string { return c.Model },
).ToSlice()

// Short-circuiting terminals
first := seq.OfSlice(cars).Filter(...).TryHead() // => option.Option[Car]
count := seq.OfSlice(cars).CountBy(func (c Car) bool { return c.Year >= 2015 })

// Generators
squares := seq.Map(seq.Range(1, 6), func (n int) int { return n * n }).ToSlice()
// => [1 4 9 16 25]

// RangeStep: custom step, supports descending
evens := seq.RangeStep(0, 10, 2).ToSlice() // => [0 2 4 6 8]
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
fibs := seq.Unfold([2]int{0, 1}, func (s [2]int) option.Option[seq.Pair[int, [2]int]] {
if s[0] > 20 {
return option.None[seq.Pair[int, [2]int]]()
}
return option.Some(seq.Pair[int, [2]int]{First: s[0], Second: [2]int{s[1], s[0] + s[1]}})
}).ToSlice()
// => [0 1 1 2 3 5 8 13]

// Partition: one pass, two slices
evens, odds := seq.Partition(seq.Range(1, 7), func (n int) bool { return n%2 == 0 })
// evens => [2 4 6],  odds => [1 3 5]

// OfMap: iterate over a map
counts := seq.CountByKey(seq.OfMap(map[string]int{"a": 1, "b": 2}),
func (p seq.Pair[string, int]) string { return p.First })
// => map[a:1 b:1]

// CountByKey: occurrence counts
freq := seq.CountByKey(seq.OfSlice([]string{"a", "b", "a", "c", "a", "b"}),
func (s string) string { return s })
// => map[a:3 b:2 c:1]

// OfOption: lift an Option into a Seq
seq.OfOption(option.Some(42)).ToSlice() // => [42]
seq.OfOption(option.None[int]()).ToSlice() // => []

// OfResult: lift a Result into a Seq (Ok => singleton, Err => empty)
seq.OfResult(result.Ok[int, string](7)).ToSlice()  // => [7]
seq.OfResult(result.Err[int, string]("e")).ToSlice() // => []

// Intersperse: insert separator between elements
seq.OfSlice([]int{1, 2, 3}).Intersperse(0).ToSlice() // => [1 0 2 0 3]

// StepBy: yield every n-th element starting from the first
seq.Range(0, 10).StepBy(3).ToSlice() // => [0 3 6 9]

// ToMap / ToMapBy: materialise into a map
m := seq.ToMap(seq.OfSlice([]seq.Pair[string, int]{{"a", 1}, {"b", 2}}))
// => map[a:1 b:2]

type Car struct { Year int; Owner, Model string }
byOwner := seq.ToMapBy(seq.OfSlice(cars),
    func(c Car) string { return c.Owner },
    func(c Car) int    { return c.Year },
)
// => map[Alice:2012 Bob:2016 ...]
```

### pseq: parallel sequences

Parallel counterparts to the most parallelism-friendly `seq` operations,
inspired by [FSharp.Collections.ParallelSeq](https://github.com/fsprojects/FSharp.Collections.ParallelSeq)
and Go's [lo/lop](https://github.com/samber/lo) parallel helpers.

Every function materialises the input `Seq[T]`, partitions it into chunks,
dispatches one goroutine per chunk, and collects results — **order is always preserved**.
Parallelism defaults to `runtime.GOMAXPROCS(0)` and is tunable via `WithWorkers`.

```go
data := seq.OfSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

// Parallel Map (order preserved)
doubled := pseq.Map(data, func (n int) int { return n * 2 }).ToSlice()
// => [2 4 6 8 10 12 14 16 18 20]

// Parallel Filter
evens := pseq.Filter(data, func (n int) bool { return n%2 == 0 }).ToSlice()
// => [2 4 6 8 10]

// Parallel Reduce (fn must be associative)
sum, _ := pseq.Reduce(data, func (a, b int) int { return a + b })
// => 55

// Parallel GroupBy
groups := pseq.GroupBy(data, func (n int) string {
if n%2 == 0 { return "even" }
return "odd"
})
// => map[even:[2 4 6 8 10] odd:[1 3 5 7 9]]

// Configure workers
pseq.Map(data, heavyFn, pseq.WithWorkers(8))

// Parallel Exists / ForAll (short-circuit across goroutines)
pseq.Exists(data, func (n int) bool { return n > 9 }) // => true
pseq.ForAll(data, func (n int) bool { return n > 0 }) // => true

// Parallel Sum / SumBy
pseq.Sum(data) // => 55

// Parallel Partition
yes, no := pseq.Partition(data, func(n int) bool { return n <= 5 })
// yes => [1 2 3 4 5], no => [6 7 8 9 10]

// Parallel Choose (filter+map with Option)
pseq.Choose(data, func (n int) option.Option[string] {
if n%3 == 0 { return option.Some(fmt.Sprintf("fizz:%d", n)) }
return option.None[string]()
}).ToSlice()
// => ["fizz:3" "fizz:6" "fizz:9"]

// Pipe integration via curried helpers
pipe.Pipe3(
seq.OfSlice(bigData),
pseq.FilterFn(isValid, pseq.WithWorkers(8)),
pseq.MapFn(transform, pseq.WithWorkers(8)),
seq.ToSliceFn[Result](),
)
```

**When to use `pseq` vs `seq`:** parallel execution pays off when the per-element
work is CPU-heavy (parsing, math, serialisation). For lightweight lambdas
(`n*2`, field access), the goroutine overhead dominates — stick with `seq`.

### option: explicit optionality

```go
// Instead of (T, bool) or *T
name := option.Some("Alice")
none := option.None[string]()

upper := option.Map(name, strings.ToUpper) // => Some("ALICE")
option.Map(none, strings.ToUpper) // => None

name.UnwrapOr("anonymous") // => "Alice"
none.UnwrapOr("anonymous") // => "anonymous"

// DefaultWith: lazy default - fn is only called when None
val := option.DefaultWith(none, func () string { return expensiveDefault() })

// Contains: value equality check
option.Contains(option.Some(42), 42) // => true
option.Contains(option.None[int](), 42) // => false

// Tee / TeeNone: side-effects without breaking the chain
opt := option.Tee(option.Some(42), func(v int) { log.Println("got", v) }) // => Some(42)
opt  = option.TeeNone(option.None[int](), func() { log.Println("missing") }) // => None

// Map2: combine two Options
option.Map2(option.Some(2), option.Some(3), func (a, b int) int { return a + b }) // => Some(5)
option.Map2(option.None[int](), option.Some(3), func (a, b int) int { return a + b }) // => None

// OrElse: fallback if None
resolved := option.OrElse(lookupCache(key), func () option.Option[string] {
return lookupDB(key)
})

// Flatten: unwrap Option[Option[T]]
option.Flatten(option.Some(option.Some(42))) // => Some(42)
option.Flatten(option.None[option.Option[int]]()) // => None

// Chain optional lookups with Bind
profile := option.Bind(findUser(id), func(u User) option.Option[Profile] {
return findProfile(u.ProfileID)
})

// Integrates with seq
firstModern := seq.OfSlice(cars).
Filter(func (c Car) bool { return c.Year >= 2015 }).
TryHead() // => option.Option[Car]
```

### result : railway-oriented error handling

```go
// Wrap Go's (T, error) convention
res := result.Try(func () (User, error) { return db.FindUser(id) })

// Railway pipeline: chain Bind (may fail) and Map (pure transforms).
// Once on the error track, every subsequent step is skipped.
r1 := parseRequest(raw) // step 1: Bind
r2 := result.Bind(r1, authenticate) // step 2: Bind
r3 := result.Map(r2, normalize)               // step 3: Map (pure)
r4 := result.Bind(r3, save) // step 4: Bind
r5 := result.Map(r4, formatResponse) // step 5: Map (pure)
// r5 is either Ok(response) or Err from whichever step failed first.

// Tee / TeeErr: side-effects without breaking the chain - great for logging
r := result.Tee(r3, func (v Request) { log.Printf("normalised: %v", v) })
r = result.TeeErr(r, func (e string) { log.Printf("failed: %s", e) })

// OrElse: try a fallback on Err
user := result.OrElse(lookupPrimary(id), func (e error) result.Result[User, error] {
return lookupReplica(id)
})

// Flatten: unwrap Result[Result[T,E],E]
result.Flatten(result.Ok[result.Result[int, string], string](result.Ok[int, string](42)))
// => Ok(42)

// Zip: combine two Results into a pair (first Err wins)
result.Zip(result.Ok[int, string](1), result.Ok[string, string]("hi"))
// => Ok({1, "hi"})

// Map2: combine two Results (short-circuits on first Err)
result.Map2(
    result.Ok[int, string](3),
    result.Ok[int, string](4),
    func(a, b int) int { return a + b },
) // => Ok(7)

// Contains: value equality check
result.Contains(result.Ok[int, string](42), 42) // => true

// Sequence: []Result => Result[[]T] (short-circuits on first Err)
result.Sequence([]result.Result[int, string]{
    result.Ok[int, string](1),
    result.Ok[int, string](2),
}) // => Ok([1 2])

// Traverse: map + sequence in one pass (short-circuits on first Err)
result.Traverse([]string{"1", "2", "3"}, func(s string) result.Result[int, string] {
    n, err := strconv.Atoi(s)
    if err != nil { return result.Err[int, string](err.Error()) }
    return result.Ok[int, string](n)
}) // => Ok([1 2 3])

// MapErr adds context to errors
wrapped := result.MapErr(r5, func (e string) string {
return "request failed: " + e
})

// Interop with option
opt := res.ToOption() // Ok => Some, Err => None
res2 := result.FromOption(opt, errors.New("not found"))
```

### kv: key-value pipelines

Functional operations over `iter.Seq2[K,V]` — the lazy, composable counterpart to
Go's `maps` package.

```go
inventory := map[string]int{
    "apple": 50, "banana": 3, "cherry": 120, "date": 0,
}

// Of wraps a map into a lazy Seq2
s := kv.Of(inventory)

// Filter: keep only non-zero stock
inStock := kv.Filter(s, func(_ string, qty int) bool { return qty > 0 })

// MapValues: apply a discount
discounted := kv.MapValues(inStock, func(qty int) int { return qty * 9 / 10 })

// Collect: materialise back to a map
result := kv.Collect(discounted)
// => map[apple:45 cherry:108]

// Keys / Values: extract as seq.Seq
keys := kv.Keys(s).SortWith(cmp.Compare).ToSlice()
// => [apple banana cherry date]

// MapKeys: transform keys
upper := kv.Collect(kv.MapKeys(s, strings.ToUpper))
// => map[APPLE:50 BANANA:3 ...]

// Fold: reduce to a single value
total := kv.Fold(s, 0, func(acc int, _ string, qty int) int { return acc + qty })
// => 173

// ContainsKey: short-circuiting membership test
kv.ContainsKey(s, "apple") // => true
kv.ContainsKey(s, "mango") // => false

// ToSeq / FromSeq: bridge to seq.Seq[seq.Pair[K,V]]
pairs := kv.ToSeq(s).Filter(func(p seq.Pair[string, int]) bool { return p.Second > 10 }).ToSlice()
back  := kv.Collect(kv.FromSeq(seq.OfSlice(pairs)))
```

### pipe : threading values

**Naming convention** - the suffix tells you the type contract:

| Suffix  | Variant                                | Type contract                                          | Example                 |
|---------|----------------------------------------|--------------------------------------------------------|-------------------------|
| `EndoN` | `PipeEndoN`, `ComposeEndoN`            | variadic, all steps `T => T` (same type, endomorphic)  | string transforms       |
| `2..8`  | `Pipe2`-`Pipe8`, `Compose2`-`Compose4` | fixed-arity, each step may change type (`A => B => C`) | parse + validate + save |

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
rawInput,          // string
strings.TrimSpace, // string => string
parse,             // string => int
validate,          // int    => error
)

// Tap: side-effect (logging, metrics) without changing the value
result := pipe.Pipe4(
rawInput,
strings.TrimSpace,
pipe.Tap(func (s string) { log.Println("trimmed:", s) }),
strings.ToUpper,
validate,
)
```

#### Pipelines longer than 8 type-changing steps

Use `ComposeEndoN`/`Compose2`-`Compose4` to collapse multiple steps into one slot:

```go
pipe.Pipe5(
seq.OfSlice(people),
seq.FilterFn(adult), // Seq[Person] => Seq[Person]
seq.MapFn(getName), // Seq[Person] => Seq[string]
seq.DistinctFn[string](), // Seq[string] => Seq[string]
pipe.ComposeEndoN(        // group same-type steps into one slot
seq.SortWithFn[string](cmp.Compare),
seq.TruncateFn[string](3),
),
seq.ToSliceFn[string](), // Seq[string] => []string
)
```

`pipe.ComposeEndoN` merges consecutive same-type steps (`T=>T`) into one pipe slot.
`pipe.Compose2`-`Compose4` do the same for type-changing steps (`A=>B=>C`).
Together they cover pipelines of any length.

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

### Why does `pseq` use chunking instead of per-element goroutines?

Spawning one goroutine per element (as `lo/parallel` does) is simple but scales
poorly: 100k items = 100k goroutines = ~300k allocations and ~200–400 ms of pure
scheduler overhead before any real work begins.

`pseq` splits the input into `GOMAXPROCS` chunks (default 8 on a typical machine)
and runs one goroutine per chunk. This is the same strategy .NET's PLINQ uses
under the hood (which F#'s `PSeq` wraps). It means:

- **8 goroutines** instead of 100k → ~300× fewer allocs
- Each goroutine processes a contiguous slice → **cache-friendly** sequential access
- Overhead is constant regardless of input size → **O(workers)**, not **O(n)**
- Users can tune it via `WithWorkers(n)` when the default doesn't fit

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

### pseq vs lo/parallel

[lo/parallel](https://github.com/samber/lo) spawns **one goroutine per element** —
simple, but O(n) scheduling overhead. `pseq` partitions into `GOMAXPROCS` chunks
and runs **one goroutine per chunk** (the same strategy as .NET's PLINQ / F#'s `PSeq`).

Benchmarks on `[]int` pipelines (Intel Core Ultra 7, 8 cores):

#### CPU-heavy workload (500 sqrt iterations per element)

| Operation       | `seq` (sequential) | `pseq` (chunked) | `lo/parallel` (per-element) | pseq vs lo       |
|-----------------|-------------------:|-----------------:|----------------------------:|------------------|
| **Map 1k**      |           1,968 µs |       **717 µs** |                      724 µs | ≈ tied           |
| **Map 10k**     |          19,243 µs |     **5,509 µs** |                    6,815 µs | **1.24× faster** |
| **Map 100k**    |         192,687 µs |    **44,555 µs** |                   57,882 µs | **1.30× faster** |
| **ForEach 10k** |                  — |       **328 µs** |                    2,904 µs | **8.9× faster**  |
| **GroupBy 10k** |                  — |     **4,587 µs** |                    5,170 µs | **1.13× faster** |

#### Lightweight workload (`n*3+1` — exposes overhead)

| Operation    | `seq` (sequential) | `pseq` (chunked) | `lo/parallel` (per-element) | pseq vs lo       |
|--------------|-------------------:|-----------------:|----------------------------:|------------------|
| **Map 1k**   |           **6 µs** |            25 µs |                      316 µs | **12.7× faster** |
| **Map 10k**  |         **122 µs** |           357 µs |                    3,473 µs | **9.7× faster**  |
| **Map 100k** |       **1,427 µs** |         3,635 µs |                   34,398 µs | **9.5× faster**  |

#### Memory (10k Map)

|               | `pseq` | `lo/parallel` | ratio                              |
|---------------|--------|---------------|------------------------------------|
| **B/op**      | 798 KB | 1,067 KB      | lo uses 1.3× more memory           |
| **allocs/op** | **66** | 20,051        | lo allocates **303× more objects** |

**Why?** `lo` does `go func(...)` inside a `for i, item := range` — 10k goroutines
for 10k items. `pseq` splits into ~8 chunks. Goroutine spawn+schedule is ~2-4 µs each,
so `lo` pays ~20-40 ms in scheduling alone for 10k items, while `pseq` pays ~16-32 µs.
When the per-element work is heavy enough, both approaches saturate the CPUs and converge.
When it isn't, `lo` is 10-13× slower than `pseq` — and even slower than sequential `seq`.

**Rule of thumb:** for lightweight lambdas, don't parallelize at all — use `seq`.
For CPU-heavy work (parsing, crypto, compression, complex transforms), `pseq` gives
the parallel speedup with a fraction of the scheduling cost.

Run the benchmarks yourself:

```sh
# seq pipeline benchmarks
go test ./seq/ -bench=. -benchmem

# pseq benchmarks (includes vs-lo comparison)
go test ./pseq/ -bench=. -benchmem
```

## Dependency graph

```text
pipe       =>  (none)
option     =>  (none)
result     =>  option
validation =>  option, result
seq        =>  option, result
pseq       =>  seq, option
kv         =>  seq
```
