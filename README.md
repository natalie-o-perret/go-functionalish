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
| `set`        | `Set[T]`: immutable set with full algebra and fluent generic methods   |
| `list`       | `List[T]`: immutable list with private backing slice; truly immutable  |

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
"github.com/natalie-o-perret/go-functionalish/set"
"github.com/natalie-o-perret/go-functionalish/list"
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

// Windowed / ChunkBySize / SplitInto: return ChunkedSeq[T], a distinct type
// that avoids the Seq[T] → Seq[[]T] instantiation cycle.
// Use ChunkedToSeq to re-enter normal Seq[T] chaining.
windows := seq.Range(1, 6).Windowed(3).ToSlice()
// => [[1 2 3] [2 3 4] [3 4 5]]

chunks := seq.Range(1, 11).ChunkBySize(3).ToSlice()
// => [[1 2 3] [4 5 6] [7 8 9] [10]]

// Bridge back to Seq[[]T] for further chaining
seq.ChunkedToSeq(seq.Range(1, 11).ChunkBySize(3)).
    Filter(func(c []int) bool { return len(c) == 3 }).
    ToSlice()
// => [[1 2 3] [4 5 6] [7 8 9]]

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

````go
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
````

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

### set: immutable sets

`Set[T]` is a `comparable`-constrained set backed by a `map[T]struct{}`. All
operations that "modify" the set return a new `Set`; the receiver is never
changed.

```go
import "github.com/natalie-o-perret/go-functionalish/set"

s  := set.Of(1, 2, 3, 4, 5)
s2 := set.Of(3, 4, 5, 6, 7)

// Set algebra (all return a new Set)
s.Union(s2)                // {1 2 3 4 5 6 7}
s.Intersect(s2)            // {3 4 5}
s.Difference(s2)           // {1 2}
s.SymmetricDifference(s2)  // {1 2 6 7}

// Add / Remove (immutable)
s.Add(6, 7)    // {1 2 3 4 5 6 7}
s.Remove(1, 2) // {3 4 5}

// Queries
s.Contains(3)                               // true
s.IsSubset(set.Of(1, 2, 3, 4, 5, 6))       // true
s.IsSuperset(set.Of(1, 2))                  // true
s.Len()                                     // 5
s.ForAll(func(n int) bool { return n > 0 }) // true
s.Exists(func(n int) bool { return n > 4 }) // true
s.CountBy(func(n int) bool { return n%2 == 0 }) // 2

// Filter / Exclude
evens := s.Filter(func(n int) bool { return n%2 == 0 }) // {2 4}

// Type-changing methods (Go 1.27 generic methods)
set.Of(1, 2, 3, 4).
    Filter(func(n int) bool { return n%2 == 0 }).
    Map(func(n int) string { return strconv.Itoa(n) }).
    ToSlice() // ["2" "4"] (order unspecified)

// Fold: reduce to a single value (use only for commutative operations)
sum := set.Of(1, 2, 3, 4, 5).Fold(0, func(acc, n int) int { return acc + n }) // 15

// Iteration (order unspecified -- Go map randomisation)
for v := range s.All() { fmt.Println(v) }
s.Iter(func(v int) { fmt.Println(v) })
```

> **Note:** Iteration order is unspecified. Only use `Fold` for associative and
> commutative operations (sum, count, etc.).

### list: immutable lists

`List[T]` is backed by a **private** `[]T`. Because the slice is unexported,
direct index-write (`s[i] = v`) is impossible from outside the package; all
operations return a new `List`, giving F#-style `list<'T>` semantics without
the cache penalty of a linked list.

```go
import "github.com/natalie-o-perret/go-functionalish/list"

l := list.Of(1, 2, 3, 4, 5)

// All operations return a new List; the receiver is never mutated
l.Append(6, 7)                                         // [1 2 3 4 5 6 7]
l.Filter(func(n int) bool { return n%2 == 0 })         // [2 4]
l.Rev()                                                // [5 4 3 2 1]
l.SortWith(cmp.Compare)                                // [1 2 3 4 5]
l.Truncate(3)                                          // [1 2 3]
l.Skip(2)                                              // [3 4 5]

// Safe element access via Option
l.At(0)   // Some(1)
l.At(10)  // None
l.Head()  // Some(1)
l.Last()  // Some(5)

// Type-changing methods (Go 1.27 generic methods)
list.Of(1, 2, 3, 4, 5).
    Filter(func(n int) bool { return n%2 == 0 }).
    Map(func(n int) string { return strconv.Itoa(n) }).
    ToSlice() // ["2" "4"]

// Fold / FoldBack
l.Fold(0, func(acc, n int) int { return acc + n })     // 15
l.FoldBack(0, func(n, acc int) int { return acc + n }) // 15

// Scan: running accumulator
list.Of(1, 2, 3, 4, 5).
    Scan(0, func(acc, n int) int { return acc + n }).
    ToSlice() // [0 1 3 6 10 15]

// GroupBy, SortBy, DistinctBy -- all methods
byParity := l.GroupBy(func(n int) string {
    if n%2 == 0 { return "even" }
    return "odd"
}) // map[even:[2 4] odd:[1 3 5]]

// Package-level (comparable constraint prevents method form)
list.Distinct(list.Of(1, 2, 1, 3, 2))              // [1 2 3]
list.Contains(list.Of(1, 2, 3), 2)                  // true
list.Except(list.Of(1, 2, 3, 4), list.Of(2, 4))    // [1 3]
list.Zip(list.Of(1, 2, 3), list.Of("a", "b", "c")) // [{1 a} {2 b} {3 c}]

// Constructors
list.Replicate(0, 5)                         // [0 0 0 0 0]
list.Init(4, func(i int) int { return i * i }) // [0 1 4 9]

// Iteration
for v := range l.All() { fmt.Println(v) }
```

> **vs `slice`:** `Slice[T]` is a named type over `[]T`; `s[i] = v` is valid.
> `List[T]` hides the backing slice, so it's truly immutable from outside the package.
> Use `list` for F#-style `list<'T>` guarantees; use `slice` when you need
> direct indexing or interop with Go's slice APIs.

### tuple: typed tuples

```go
import "github.com/natalie-o-perret/go-functionalish/tuple"

t  := tuple.Of("alice", 30)         // T2[string, int]
t3 := tuple.Of3("x", 1, true)       // T3[string, int, bool]
t4 := tuple.Of4("x", 1, true, 3.14) // T4[string, int, bool, float64]

// Apply: call a function with the tuple's fields
t.Apply(func(name string, age int) string {
    return fmt.Sprintf("%s is %d", name, age)
}) // => "alice is 30"

// T2: MapFirst / MapSecond / Map
t.MapFirst(strings.ToUpper)                                          // T2["ALICE", 30]
t.MapSecond(func(n int) int { return n + 1 })                        // T2["alice", 31]
t.Map(strings.ToUpper, func(n int) float64 { return float64(n) * 1.5 }) // T2["ALICE", 45.0]

// T3: MapFirst / MapSecond / MapThird / Map / Drop*
t3.MapFirst(strings.ToUpper)                           // T3["X", 1, true]
t3.MapSecond(func(n int) int { return n * 10 })        // T3["x", 10, true]
t3.MapThird(func(b bool) bool { return !b })           // T3["x", 1, false]
t3.Map(strings.ToUpper, func(n int) int { return n * 2 }, func(b bool) bool { return !b })
// => T3["X", 2, false]
t3.DropFirst()  // => T2[1, true]
t3.DropSecond() // => T2["x", true]
t3.DropThird()  // => T2["x", 1]

// T4: MapFirst / MapSecond / MapThird / MapFourth / Map / Drop*
t4.MapFourth(func(f float64) float64 { return f * 2 }) // T4["x", 1, true, 6.28]
t4.DropFirst()   // => T3[1, true, 3.14]
t4.DropFourth()  // => T3["x", 1, true]

// Extend: grow a tuple by one element
tuple.Of("alice", 30).Extend(true)         // T3["alice", 30, true]
tuple.Of3("alice", 30, true).Extend(99.9)  // T4["alice", 30, true, 99.9]

// Getter functions -- useful as first-class function values
// (field access t.First works too; these exist for passing to Map/seq pipelines)
tuple.Fst(t)                               // "alice"
tuple.Snd(t)                               // 30
tuple.Thd(t3)                              // true
tuple.Fth(t4)                              // 3.14
// e.g.: seq.OfSlice(pairs).Map(tuple.Fst[string, int]).ToSlice()

// Map2 / Map3 / Map4: single function over homogeneous tuples
// (for heterogeneous tuples, use the .Map method with one fn per field)
tuple.Map2(tuple.Of(2, 3), func(n int) int { return n * 10 })    // T2[20, 30]
tuple.Map3(tuple.Of3(1, 2, 3), func(n int) string { return fmt.Sprintf("%d", n) })
// T3["1", "2", "3"]
tuple.Map4(tuple.Of4(1, 2, 3, 4), func(n int) bool { return n%2 == 0 })
// T4[false, true, false, true]

// Chain transforms
tuple.Of("hello", []int{1, 2, 3}).
    MapFirst(strings.ToUpper).
    MapSecond(func(s []int) int { return len(s) })
// => T2["HELLO", 3]

// Unpack into individual variables
name, age := t.Unpack()
x, n, b   := t3.Unpack()
```

### kv: key-value pipelines

Functional operations over `iter.Seq2[K,V]` - the lazy, composable counterpart to
Go's `maps` package.

```go
inventory := map[string]int{
    "apple": 50, "banana": 3, "cherry": 120, "date": 0,
}

// Fully fluent left-to-right pipeline
result := kv.Of(inventory).
    Filter(func(_ string, qty int) bool { return qty > 0 }).   // drop zeros
    MapValues(func(qty int) int { return qty * 9 / 10 }).       // apply 10% discount
    ToMap()
// => map[apple:45 cherry:108]

// Map: transform both key and value at once
kv.Of(inventory).
    Map(func(k string, v int) (string, string) {
        return strings.ToUpper(k), fmt.Sprintf("%d units", v)
    }).
    ToMap()
// => map[APPLE:"50 units" BANANA:"3 units" ...]

// MapValues / MapKeys independently
kv.Of(inventory).MapValues(func(v int) float64 { return float64(v) * 0.9 }).ToMap()
kv.Of(inventory).MapKeys(func(k string) string { return strings.ToUpper(k) }).ToMap()

// Keys / Values: extract as seq.Seq
kv.Of(inventory).Keys().SortWith(cmp.Compare).ToSlice()
// => [apple banana cherry date]

// Fold: reduce to a single value
total := kv.Of(inventory).Fold(0, func(acc int, _ string, qty int) int { return acc + qty })
// => 173

// ContainsKey: short-circuiting membership test
kv.Of(inventory).ContainsKey("apple") // => true
kv.Of(inventory).ContainsKey("mango") // => false

// ToSeq / FromSeq: bridge to seq.Seq[seq.Pair[K,V]]
pairs := kv.Of(inventory).ToSeq().
    Filter(func(p seq.Pair[string, int]) bool { return p.Second > 10 }).
    ToSlice()
back := kv.FromSeq(seq.OfSlice(pairs)).ToMap()
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
would need to be checked (including `Zip`, returning `Seq[Pair[Pair[T,U],V]]`,)
and so on forever. Go's type checker rejects this infinite expansion, so `Zip`
stays a package-level function.

The same constraint applies to `Windowed`, `ChunkBySize`, and `SplitInto`:
returning `Seq[[]T]` would force the type checker to instantiate `Seq[[]T]`,
then `Seq[[][]T]`, and so on. The solution is `ChunkedSeq[T]`, a distinct
named type whose method set (`ToSlice`, `ForEach`) carries no reference back to
`Seq[...]`. To re-enter normal `Seq[T]` chaining, use the package-level bridge
`seq.ChunkedToSeq[T]`.

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
poorly: 100k items = 100k goroutines = ~300k allocations and ~200-400 ms of pure
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
| **Small** (filter => map => take 100) | Direct |    ~700 |  2,064 |      9 |
|                                       | Pipe   |    ~700 |  2,064 |      9 |
| **Medium** (7 steps incl. sort)       | Direct |  ~6,300 |  8,376 |     26 |
|                                       | Pipe   |  ~6,500 |  8,912 |     44 |
| **Large** (10 steps incl. rev+sort)   | Direct | ~13,400 | 12,784 |     46 |
|                                       | Pipe   | ~13,600 | 13,592 |     73 |

**< 3% wall-clock overhead** (noise for any real workload), all from one-time
closure allocations when the pipeline is _built_, not per element. The hot
iteration loop is identical either way. Choose whichever style reads better.

### pseq vs lo/parallel

[lo/parallel](https://github.com/samber/lo) spawns **one goroutine per element**,
simple, but O(n) scheduling overhead. `pseq` partitions into `GOMAXPROCS` chunks
and runs **one goroutine per chunk** (the same strategy as .NET's PLINQ / F#'s `PSeq`).

Benchmarks on `[]int` pipelines (Intel Core Ultra 7, 8 cores):

#### CPU-heavy workload (500 sqrt iterations per element)

| Operation       | `seq` (sequential) | `pseq` (chunked) | `lo/parallel` (per-element) | pseq vs lo      |
| --------------- | -----------------: | ---------------: | --------------------------: | --------------- |
| **Map 1k**      |           1,836 µs |       **459 µs** |                      524 µs | **1.1× faster** |
| **Map 10k**     |          19,116 µs |     **4,503 µs** |                    4,865 µs | **1.1× faster** |
| **Map 100k**    |         184,837 µs |    **36,925 µs** |                   48,079 µs | **1.3× faster** |
| **ForEach 10k** |                N/A |       **406 µs** |                    2,756 µs | **6.8× faster** |
| **GroupBy 10k** |                N/A |     **4,682 µs** |                    5,058 µs | **1.1× faster** |

#### Lightweight workload (`n*3+1` - exposes overhead)

| Operation    | `seq` (sequential) | `pseq` (chunked) | `lo/parallel` (per-element) | pseq vs lo       |
| ------------ | -----------------: | ---------------: | --------------------------: | ---------------- |
| **Map 1k**   |           **7 µs** |            22 µs |                      278 µs | **12.8× faster** |
| **Map 10k**  |          **90 µs** |           300 µs |                    2,730 µs | **9.1× faster**  |
| **Map 100k** |       **1,315 µs** |         2,440 µs |                   25,608 µs | **10.5× faster** |

#### Memory (10k Map)

|               | `pseq` | `lo/parallel` | ratio                              |
| ------------- | ------ | ------------- | ---------------------------------- |
| **B/op**      | 798 KB | 1,047 KB      | lo uses 1.3× more memory           |
| **allocs/op** | **66** | 20,003        | lo allocates **303× more objects** |

**Why?** `lo` does `go func(...)` inside a `for i, item := range`, spawning 10k goroutines
for 10k items. `pseq` splits into ~8 chunks. Goroutine spawn+schedule is ~2-4 µs each,
so `lo` pays ~20-40 ms in scheduling alone for 10k items, while `pseq` pays ~16-32 µs.
When the per-element work is heavy enough, both approaches saturate the CPUs and converge.
When it isn't, `lo` is **9-13× slower** than `pseq`, and even slower than sequential `seq`.

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
set        =>  option
list       =>  option
```
