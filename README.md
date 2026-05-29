# go-functionalish

[![CI](https://github.com/natalie-o-perret/go-functionalish/actions/workflows/ci.yml/badge.svg)](https://github.com/natalie-o-perret/go-functionalish/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/natalie-o-perret/go-functionalish.svg)](https://pkg.go.dev/github.com/natalie-o-perret/go-functionalish)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Contributing](https://img.shields.io/badge/contributing-guide-blue)](CONTRIBUTING.md)

A cohesive, opinionated, type-safe functional programming library for Go 1.27+.
No reflection. No `interface{}`. Pure generics, lazy by default, and fully
fluent pipelines via Go 1.27 generic methods.

> [!NOTE]
> Unapologetically vibe-coded with GitHub Copilot & Claude Sonnet 4.6.
>
> Unapologetically not "idiomatic Go."
>
> Go gave us generics 17 years after C# and 18 after Java (the latter still erases them at runtime).
> We're using them for `Option[T]`, `Result[T,E]`, and lazy pipelines
> instead of `if err != nil` sixty times per file.

## Packages

| Package      | Description                                                            |
| ------------ | ---------------------------------------------------------------------- |
| `seq`        | Lazy `Seq[T]`: F#-style sequence pipelines with fully fluent methods   |
| `slice`      | Eager `Slice[T]`: same fluent API as `seq` over in-memory slices       |
| `pseq`       | Parallel `Seq[T]`: goroutine-per-chunk Map, Filter, Reduce, ...        |
| `option`     | `Option[T]`: explicit presence/absence, no nil                         |
| `result`     | `Result[T,E]`: railway-oriented error handling                         |
| `tuple`      | `T2 / T3 / T4`: typed tuples with fluent `Apply`, `MapFirst`, `Map`... |
| `validation` | `Validation[T,E]`: applicative error accumulation                      |
| `pipe`       | `Pipe2`...`Pipe8`: F#-style `\|>` operator equivalent                  |
| `kv`         | Lazy `Seq2[K,V]`: functional pipelines over `iter.Seq2` / maps         |

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

// Go 1.27 generic methods: type-changing operations are now methods,
// enabling fully left-to-right fluent pipelines.
owners := seq.OfSlice(cars).
    Filter(func(c Car) bool { return c.Year >= 2015 }).
    Map(func(c Car) string { return c.Owner }).
    ToSlice()
// => ["Bob", "Charlie", "Diana"]

// Sort, group, distinct -- all chainable
byCar := seq.OfSlice(cars).GroupBy(func(c Car) string { return c.Model })

unique := seq.OfSlice(cars).
    SortWith(func(a, b Car) int { return cmp.Compare(a.Model, b.Model) }).
    DistinctBy(func(c Car) string { return c.Model }).
    ToSlice()

// Short-circuiting terminals
first := seq.OfSlice(cars).Filter(...).TryHead() // => option.Option[Car]
count := seq.OfSlice(cars).CountBy(func(c Car) bool { return c.Year >= 2015 })

// Generators
squares := seq.Range(1, 6).Map(func(n int) int { return n * n }).ToSlice()
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

// CountByKey: occurrence counts
freq := seq.CountByKey(seq.OfSlice([]string{"a", "b", "a", "c", "a", "b"}),
    func(s string) string { return s })
// => map[a:3 b:2 c:1]

// OfOption: lift an Option into a Seq
seq.OfOption(option.Some(42)).ToSlice() // => [42]
seq.OfOption(option.None[int]()).ToSlice() // => []

// OfResult: lift a Result into a Seq (Ok => singleton, Err => empty)
seq.OfResult(result.Ok[int, string](7)).ToSlice()  // => [7]
seq.OfResult(result.Err[int, string]("e")).ToSlice() // => []

// StepBy: yield every n-th element starting from the first
seq.Range(0, 10).StepBy(3).ToSlice() // => [0 3 6 9]

// Scan: emit running accumulator (like a running sum)
seq.Range(1, 5).Scan(0, func(acc, n int) int { return acc + n }).ToSlice()
// => [0 1 3 6 10]

// ToMap / ToMapBy: materialise into a map
m := seq.ToMap(seq.OfSlice([]seq.Pair[string, int]{{"a", 1}, {"b", 2}}))
// => map[a:1 b:2]

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
dispatches one goroutine per chunk, and collects results. **Order is always preserved.**
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
(`n*2`, field access), the goroutine overhead dominates, so stick with `seq`.

### option: explicit optionality

```go
// Instead of (T, bool) or *T
name := option.Some("Alice")
none := option.None[string]()

// Map and Bind are now methods (Go 1.27 generic methods)
upper := name.Map(strings.ToUpper) // => Some("ALICE")
none.Map(strings.ToUpper)          // => None

name.UnwrapOr("anonymous") // => "Alice"
none.UnwrapOr("anonymous") // => "anonymous"

// Chain with Bind
profile := findUser(id).
    Bind(func(u User) option.Option[Profile] { return findProfile(u.ProfileID) })

// DefaultWith: lazy default -- fn is only called when None
val := none.DefaultWith(func() string { return expensiveDefault() })

// Contains: value equality check
option.Contains(option.Some(42), 42) // => true

// Tee / TeeNone: side-effects without breaking the chain
opt := option.Some(42).Tee(func(v int) { log.Println("got", v) }) // => Some(42)

// ZipWith: combine two Options with a function
option.Some(2).ZipWith(option.Some(3), func(a, b int) int { return a + b }) // => Some(5)

// OrElse: fallback if None
resolved := lookupCache(key).OrElse(func() option.Option[string] { return lookupDB(key) })

// Flatten: unwrap Option[Option[T]]
option.Flatten(option.Some(option.Some(42))) // => Some(42)

// Integrates with seq
firstModern := seq.OfSlice(cars).
    Filter(func(c Car) bool { return c.Year >= 2015 }).
    TryHead() // => option.Option[Car]

```go
// Wrap Go's (T, error) convention
res := result.Try(func() (User, error) { return db.FindUser(id) })

// Railway pipeline -- Map and Bind are now methods (Go 1.27 generic methods).
// Once on the error track, every subsequent step is skipped.
r := parseRequest(raw).          // Result[Request, string]
    Bind(authenticate).           // step 2: may fail
    Map(normalize).               // step 3: pure transform
    Bind(save).                   // step 4: may fail
    Map(formatResponse).          // step 5: pure transform
    MapErr(func(e string) string { return "request failed: " + e })
// r is either Ok(response) or Err from whichever step failed first.

// Tee / TeeErr: side-effects without breaking the chain
r.Tee(func(v Response) { log.Printf("ok: %v", v) })
r.TeeErr(func(e string) { log.Printf("err: %s", e) })

// OrElse: try a fallback on Err
user := lookupPrimary(id).OrElse(func(e error) result.Result[User, error] {
})

// Flatten: unwrap Result[Result[T,E],E]
result.Flatten(result.Ok[result.Result[int, string], string](result.Ok[int, string](42)))
// => Ok(42)

// Zip: combine two Results into a pair (first Err wins)
result.Zip(result.Ok[int, string](1), result.Ok[string, string]("hi"))
// => Ok({1, "hi"})

// ZipWith: combine two Results with a function (short-circuits on first Err)
result.Ok[int, string](3).ZipWith(result.Ok[int, string](4),
    func(a, b int) int { return a + b }) // => Ok(7)

// Sequence: []Result => Result[[]T] (short-circuits on first Err)
result.Sequence([]result.Result[int, string]{
    result.Ok[int, string](1), result.Ok[int, string](2),
}) // => Ok([1 2])

// Traverse: map + sequence in one pass
result.Traverse([]string{"1", "2", "3"}, func(s string) result.Result[int, string] {
    n, err := strconv.Atoi(s)
    if err != nil { return result.Err[int, string](err.Error()) }
    return result.Ok[int, string](n)
}) // => Ok([1 2 3])

// Interop with option
opt := res.ToOption() // Ok => Some, Err => None
res2 := result.FromOption(opt, errors.New("not found"))
```

### slice: eager in-memory sequences

`Slice[T]` is a named type over `[]T` with the same fluent method set as `Seq[T]`
but operating eagerly on in-memory data. Use it when you already have a slice and
don't need lazy evaluation.

```go
import goslice "github.com/natalie-o-perret/go-functionalish/slice"

type User struct { Name string; Active bool; Score int }

users := goslice.Of(
    User{"alice", true, 80}, User{"bob", false, 90},
    User{"carol", true, 70}, User{"dave", true, 85},
)

// Fully fluent -- type-changing methods chain left-to-right
names := users.
    Filter(func(u User) bool { return u.Active }).
    SortByDescending(func(u User) int { return u.Score }).
    Map(func(u User) string { return u.Name })
// => Slice["dave", "alice", "carol"]

// GroupBy, Choose, Fold -- all methods
byActive := users.GroupBy(func(u User) bool { return u.Active })
// => map[true:[...] false:[...]]

total := users.Fold(0, func(acc int, u User) int { return acc + u.Score })
// => 325

// Scan: running totals
users.Map(func(u User) int { return u.Score }).
    Scan(0, func(acc, n int) int { return acc + n }).
    ToSlice()
// => [0 80 170 240 325]

// Constructors
goslice.Replicate(0, 5)           // => [0 0 0 0 0]
goslice.Init(4, func(i int) int { return i * i }) // => [0 1 4 9]

// Package-level (constraint prevents method form)
goslice.Zip(goslice.Of(1, 2, 3), goslice.Of("a", "b", "c"))
// => [{1 a} {2 b} {3 c}]

goslice.Distinct(goslice.Of(1, 2, 1, 3, 2)) // => [1 2 3]
goslice.Except(goslice.Of(1,2,3,4), goslice.Of(2,4)) // => [1 3]
```

### tuple: typed tuples

```go
import "github.com/natalie-o-perret/go-functionalish/tuple"

t := tuple.New2("alice", 30)        // T2[string, int]
t3 := tuple.New3("x", 1, true)      // T3[string, int, bool]

// Apply: call a function with the tuple's fields
t.Apply(func(name string, age int) string {
    return fmt.Sprintf("%s is %d", name, age)
}) // => "alice is 30"

// MapFirst / MapSecond: transform individual fields
t.MapFirst(strings.ToUpper)   // => T2["ALICE", 30]
t.MapSecond(func(n int) int { return n + 1 }) // => T2["alice", 31]

// Map: transform both fields independently
t.Map(strings.ToUpper, func(n int) float64 { return float64(n) * 1.5 })
// => T2["ALICE", 45.0]

// Chain transforms
tuple.New2("hello", []int{1, 2, 3}).
    MapFirst(strings.ToUpper).
    MapSecond(func(s []int) int { return len(s) })
// => T2["HELLO", 3]

// Unpack into individual variables
name, age := t.Unpack()
```

### kv: key-value pipelines

Functional operations over `iter.Seq2[K,V]` - the lazy, composable counterpart to
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

## Design notes

### Fluent pipelines via Go 1.27 generic methods

Go 1.27 ([#77273](https://github.com/golang/go/issues/77273)) lifted the
restriction that methods cannot introduce new type parameters. Operations that
transform the element type (`Map`, `Collect`, `Choose`, `Fold`, `GroupBy`,
`SortBy`, etc.) are now **methods**, enabling fully left-to-right pipelines:

```go
// Before Go 1.27: inside-out function calls
seq.Map(seq.OfSlice(cars).Filter(isActive), toName)

// Go 1.27+: fluent methods
seq.OfSlice(cars).Filter(isActive).Map(toName).SortBy(strings.ToLower).ToSlice()
```

One exception: `Zip` cannot be a method on `Seq[T]`. If it were, its return type
`Seq[Pair[T,U]]` would be another instantiation of `Seq`, whose own method set
would need to be checked — including `Zip`, returning `Seq[Pair[Pair[T,U],V]]`,
and so on forever. Go's type checker rejects this infinite expansion, so `Zip`
stays a package-level function.

### Lazy (`seq`) vs eager (`slice`)

`Seq[T]` is a `func(yield func(T) bool)`, a lazy pull iterator. Every
pipeline operation wraps the previous iterator; **nothing runs** until a terminal
(`ToSlice`, `Head`, `Length`, ...) is called. Only operations that require the
full sequence (`SortWith`, `SortBy`, `Rev`) must materialise early.

`Slice[T]` is a named type over `[]T`. Every operation allocates a new slice
immediately. Use it when you already have a slice, need random access, or want
predictable allocation behaviour.

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
| ------------------------------------- | ------ | ------: | -----: | -----: |
| **Small** (filter => map => take 100) | Direct |  ~1,780 |  2,064 |      9 |
|                                       | Pipe   |  ~1,620 |  2,064 |      9 |
| **Medium** (7 steps incl. sort)       | Direct | ~13,100 |  8,376 |     26 |
|                                       | Pipe   | ~13,900 |  8,912 |     44 |
| **Large** (10 steps incl. rev+sort)   | Direct | ~27,100 | 12,784 |     46 |
|                                       | Pipe   | ~30,400 | 13,592 |     73 |

**~6-13 % wall-clock overhead**, all from one-time closure allocations when the
pipeline is _built_, not per element. The hot iteration loop is identical either
way. For any real workload (I/O, serialisation, business logic in the lambdas)
this is noise: choose whichever style reads better.

### pseq vs lo/parallel

[lo/parallel](https://github.com/samber/lo) spawns **one goroutine per element**,
simple, but O(n) scheduling overhead. `pseq` partitions into `GOMAXPROCS` chunks
and runs **one goroutine per chunk** (the same strategy as .NET's PLINQ / F#'s `PSeq`).

Benchmarks on `[]int` pipelines (Intel Core Ultra 7, 8 cores):

#### CPU-heavy workload (500 sqrt iterations per element)

| Operation       | `seq` (sequential) | `pseq` (chunked) | `lo/parallel` (per-element) | pseq vs lo       |
| --------------- | -----------------: | ---------------: | --------------------------: | ---------------- |
| **Map 1k**      |           1,968 µs |       **717 µs** |                      724 µs | ≈ tied           |
| **Map 10k**     |          19,243 µs |     **5,509 µs** |                    6,815 µs | **1.24× faster** |
| **Map 100k**    |         192,687 µs |    **44,555 µs** |                   57,882 µs | **1.30× faster** |
| **ForEach 10k** |                N/A |       **328 µs** |                    2,904 µs | **8.9× faster**  |
| **GroupBy 10k** |                N/A |     **4,587 µs** |                    5,170 µs | **1.13× faster** |

#### Lightweight workload (`n*3+1` - exposes overhead)

| Operation    | `seq` (sequential) | `pseq` (chunked) | `lo/parallel` (per-element) | pseq vs lo       |
| ------------ | -----------------: | ---------------: | --------------------------: | ---------------- |
| **Map 1k**   |           **6 µs** |            25 µs |                      316 µs | **12.7× faster** |
| **Map 10k**  |         **122 µs** |           357 µs |                    3,473 µs | **9.7× faster**  |
| **Map 100k** |       **1,427 µs** |         3,635 µs |                   34,398 µs | **9.5× faster**  |

#### Memory (10k Map)

|               | `pseq` | `lo/parallel` | ratio                              |
| ------------- | ------ | ------------- | ---------------------------------- |
| **B/op**      | 798 KB | 1,067 KB      | lo uses 1.3× more memory           |
| **allocs/op** | **66** | 20,051        | lo allocates **303× more objects** |

**Why?** `lo` does `go func(...)` inside a `for i, item := range`, spawning 10k goroutines
for 10k items. `pseq` splits into ~8 chunks. Goroutine spawn+schedule is ~2-4 µs each,
so `lo` pays ~20-40 ms in scheduling alone for 10k items, while `pseq` pays ~16-32 µs.
When the per-element work is heavy enough, both approaches saturate the CPUs and converge.
When it isn't, `lo` is 10-13x slower than `pseq`, and even slower than sequential `seq`.

**Rule of thumb:** for lightweight lambdas, don't parallelize at all - use `seq`.
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
tuple      =>  (none)
result     =>  option
validation =>  option, result
seq        =>  option, result
slice      =>  option
pseq       =>  seq, option
kv         =>  seq
```
