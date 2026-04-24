package seq_test

import (
	"cmp"
	"fmt"
	"strings"
	"testing"

	"github.com/natalie-o-perret/go-functionalish/option"
	"github.com/natalie-o-perret/go-functionalish/pipe"
	"github.com/natalie-o-perret/go-functionalish/seq"
)

type person struct {
	Name string
	Age  int
}

func TestThenMap(t *testing.T) {
	got := seq.Then(
		seq.OfSlice([]int{1, 2, 3}),
		seq.MapFn(func(n int) int { return n * 10 }),
	).ToSlice()
	assertSlice(t, got, []int{10, 20, 30})
}

func TestThenChained(t *testing.T) {
	people := []person{{"Alice", 30}, {"Bob", 17}, {"Charlie", 25}}

	// Filter (method) => Map (Then) => terminal
	got := seq.Then(
		seq.OfSlice(people).Filter(func(p person) bool { return p.Age >= 18 }),
		seq.MapFn(func(p person) string { return p.Name }),
	).ToSlice()
	assertSlice(t, got, []string{"Alice", "Charlie"})
}

func TestThenWithPipe(t *testing.T) {
	got := pipe.Pipe2(
		seq.OfSlice([]string{" Alice ", " Bob "}),
		seq.MapFn(strings.TrimSpace),
		seq.MapFn(strings.ToUpper),
	).ToSlice()
	assertSlice(t, got, []string{"ALICE", "BOB"})
}

func TestPipeMultipleMaps(t *testing.T) {
	// pipe chains multiple type-changing transforms without nesting
	got := pipe.Pipe3(
		seq.OfSlice([]person{{"Alice", 30}, {"Bob", 17}, {"Charlie", 25}}),
		seq.FilterFn(func(p person) bool { return p.Age >= 18 }),
		seq.MapFn(func(p person) string { return p.Name }),
		seq.MapFn(strings.ToUpper),
	).ToSlice()
	assertSlice(t, got, []string{"ALICE", "CHARLIE"})
}

func TestPipeFullPipeline(t *testing.T) {
	// Mix same-type and type-changing operations all in one pipe
	got := pipe.Pipe6(
		seq.OfSlice([]person{
			{"Alice", 30}, {"Bob", 17}, {"Charlie", 25},
			{"Diana", 30}, {"Eve", 22}, {"Alice", 28},
		}),
		seq.FilterFn(func(p person) bool { return p.Age >= 18 }),
		seq.MapFn(func(p person) string { return p.Name }),
		seq.DistinctFn[string](),
		seq.SortWithFn[string](cmp.Compare),
		seq.TruncateFn[string](3),
		seq.ToSliceFn[string](),
	)
	assertSlice(t, got, []string{"Alice", "Charlie", "Diana"})
}

func TestPipeSameTypeChain(t *testing.T) {
	// All same-type operations via pipe
	got := pipe.Pipe5(
		seq.Range(1, 20),
		seq.FilterFn(func(n int) bool { return n%2 == 0 }),
		seq.SkipFn[int](2),
		seq.TruncateFn[int](4),
		seq.RevFn[int](),
		seq.ToSliceFn[int](),
	)
	assertSlice(t, got, []int{12, 10, 8, 6})
}

func TestSortByFn(t *testing.T) {
	got := seq.Then(
		seq.OfSlice([]person{{"Charlie", 25}, {"Alice", 30}, {"Bob", 17}}),
		seq.SortByFn(func(p person) string { return p.Name }),
	).ToSlice()
	if got[0].Name != "Alice" || got[2].Name != "Charlie" {
		t.Fatalf("got %v", got)
	}
}

func TestDistinctFn(t *testing.T) {
	got := seq.Then(
		seq.OfSlice([]int{1, 2, 1, 3, 2}),
		seq.DistinctFn[int](),
	).ToSlice()
	assertSlice(t, got, []int{1, 2, 3})
}

func TestChooseFn(t *testing.T) {
	got := seq.Then(
		seq.OfSlice([]int{1, 2, 3, 4, 5}),
		seq.ChooseFn(func(n int) option.Option[int] {
			if n%2 == 0 {
				return option.Some(n * 10)
			}
			return option.None[int]()
		}),
	).ToSlice()
	assertSlice(t, got, []int{20, 40})
}

func TestFoldFn(t *testing.T) {
	sum := seq.Then(
		seq.Range(1, 6),
		seq.FoldFn(0, func(acc, v int) int { return acc + v }),
	)
	if sum != 15 {
		t.Fatalf("got %d", sum)
	}
}

func TestGroupByFn(t *testing.T) {
	groups := seq.Then(
		seq.Range(1, 6),
		seq.GroupByFn(func(n int) string {
			if n%2 == 0 {
				return "even"
			}
			return "odd"
		}),
	)
	if len(groups["even"]) != 2 || len(groups["odd"]) != 3 {
		t.Fatalf("got %v", groups)
	}
}

func TestPipelineComposition(t *testing.T) {
	// Full pipeline: filter => map type => sort => distinct => take
	people := []person{
		{"Alice", 30}, {"Bob", 17}, {"Charlie", 25},
		{"Diana", 30}, {"Eve", 22}, {"Alice", 28},
	}

	got := seq.Then(
		seq.Then(
			seq.OfSlice(people).Filter(func(p person) bool { return p.Age >= 18 }),
			seq.MapFn(func(p person) string { return p.Name }),
		),
		seq.DistinctFn[string](),
	).SortWith(cmp.Compare).Truncate(3).ToSlice()

	assertSlice(t, got, []string{"Alice", "Charlie", "Diana"})
}

// -- long pipeline examples ----------------------------------------------------

func TestPipeline10Steps(t *testing.T) {
	type employee struct {
		Name       string
		Department string
		Salary     int
		Active     bool
		YearsExp   int
	}

	data := []employee{
		{"Alice", "Engineering", 95000, true, 8},
		{"Bob", "Marketing", 72000, true, 5},
		{"Charlie", "Engineering", 88000, false, 6},
		{"Diana", "Sales", 65000, true, 3},
		{"Eve", "Engineering", 110000, true, 12},
		{"Frank", "Marketing", 72000, true, 4},
		{"Grace", "Sales", 58000, true, 2},
		{"Hank", "Engineering", 92000, true, 7},
		{"Ivy", "Marketing", 81000, true, 6},
		{"Jack", "Sales", 45000, true, 1},
	}

	// 10-step transformation pipeline:
	//  1. Filter      - keep only active employees
	//  2. Exclude     - remove < 3 years experience
	//  3. SortWith    - sort by salary descending
	//  4. Skip        - drop the highest earner
	//  5. Truncate    - keep the next 4
	//  6. Map         - extract name (employee -> string)
	//  7. Distinct    - deduplicate names
	//  8. Map         - uppercase each name
	//  9. SortWith    - sort alphabetically
	// 10. Fold        - join into comma-separated string

	// Steps 1-5: same-type method chain on Seq[employee]
	top4 := seq.OfSlice(data).
		Filter(func(e employee) bool { return e.Active }).        // 1
		Exclude(func(e employee) bool { return e.YearsExp < 3 }). // 2
		SortWith(func(a, b employee) int {                        // 3
			return cmp.Compare(b.Salary, a.Salary)
		}).
		Skip(1).    // 4
		Truncate(4) // 5

	// Steps 6-9: type-changing operations via Then, then method chaining
	names := seq.Then(
		seq.Then(
			top4,
			seq.MapFn(func(e employee) string { return e.Name }), // 6
		),
		seq.DistinctFn[string](), // 7
	)
	sorted := seq.Then(names, seq.MapFn(strings.ToUpper)). // 8
								SortWith(cmp.Compare) // 9

	// Step 10: terminal
	got := seq.Fold(sorted, "", func(acc, name string) string { // 10
		if acc == "" {
			return name
		}
		return acc + ", " + name
	})

	// Trace:
	//   Active + >=3yr: Alice(95k), Bob(72k), Diana(65k), Eve(110k),
	//     Frank(72k), Hank(92k), Ivy(81k)
	//   Sorted desc: Eve, Alice, Hank, Ivy, Bob, Frank, Diana
	//   Skip(1) + Truncate(4): Alice, Hank, Ivy, Bob
	//   Map+Distinct+Upper+Sort: ALICE, BOB, HANK, IVY
	expected := "ALICE, BOB, HANK, IVY"
	if got != expected {
		t.Fatalf("10-step pipeline:\n got: %s\nwant: %s", got, expected)
	}
}

func TestPipeline20Steps(t *testing.T) {
	// 20-step transformation of Range(1, 100)  - all in one pipe expression.
	//
	//  1. Filter      - keep even numbers
	//  2. Exclude     - drop multiples of 10
	//  3. Skip        - skip first 3
	//  4. TakeWhile   - stop at 80
	//  5. Truncate    - keep 12
	//  6. Append      - add sentinel values [1000, 2000]
	//  7. Rev         - reverse the sequence
	//  8. Skip        - drop the first element (2000)
	//  9. Truncate    - keep 10
	// 10. SortWith    - sort ascending
	// 11. Filter      - drop values >= 100
	// 12. Map         - double each value
	// 13. Scan        - running sum
	// 14. Skip        - drop initial accumulator (0)
	// 15. Filter      - keep running sums > 100
	// 16. Indexed     - pair each value with its index
	// 17. Map         - format as "#index=value"
	// 18. Rev         - reverse order
	// 19. Truncate    - keep top 5
	// 20. Fold        - join with " | "

	got := pipe.Pipe4(
		seq.Range(1, 100),

		// Steps 1-15: all Seq[int] -> Seq[int]
		pipe.ComposeEndoN(
			seq.FilterFn(func(n int) bool { return n%2 == 0 }),   // 1
			seq.ExcludeFn(func(n int) bool { return n%10 == 0 }), // 2
			seq.SkipFn[int](3), // 3
			seq.TakeWhileFn(func(n int) bool { return n < 80 }),    // 4
			seq.TruncateFn[int](12),                                // 5
			seq.AppendFn(seq.OfSlice([]int{1000, 2000})),           // 6
			seq.RevFn[int](),                                       // 7
			seq.SkipFn[int](1),                                     // 8
			seq.TruncateFn[int](10),                                // 9
			seq.SortWithFn[int](cmp.Compare),                       // 10
			seq.FilterFn(func(n int) bool { return n < 100 }),      // 11
			seq.MapFn(func(n int) int { return n * 2 }),            // 12
			seq.ScanFn(0, func(acc, n int) int { return acc + n }), // 13
			seq.SkipFn[int](1),                                     // 14
			seq.FilterFn(func(n int) bool { return n > 100 }),      // 15
		),

		// Steps 16-17: Seq[int] -> Seq[Pair[int,int]] -> Seq[string]
		pipe.Compose2(
			seq.IndexedFn[int](), // 16
			seq.MapFn(func(p seq.Pair[int, int]) string { // 17
				return fmt.Sprintf("#%d=%d", p.First, p.Second)
			}),
		),

		// Steps 18-19: Seq[string] -> Seq[string]
		pipe.ComposeEndoN(
			seq.RevFn[string](),       // 18
			seq.TruncateFn[string](5), // 19
		),

		// Step 20: Seq[string] -> string
		seq.FoldFn("", func(acc, s string) string { // 20
			if acc == "" {
				return s
			}
			return acc + " | " + s
		}),
	)

	// Trace:
	//   [1..99] -> evens -> drop ×10 -> skip 3 -> <80 -> first 12
	//   -> [8..36] + [1000,2000] -> rev -> skip 2000 -> first 10 -> sort
	//   -> [16,18,22,24,26,28,32,34,36,1000]
	//   -> drop >=100 -> double -> [32,36,44,48,52,56,64,68,72]
	//   -> scan(+) -> [0,32,68,112,160,212,268,332,400,472] -> skip 0
	//   -> keep >100 -> indexed -> format -> rev -> top 5 -> join
	expected := "#6=472 | #5=400 | #4=332 | #3=268 | #2=212"
	if got != expected {
		t.Fatalf("20-step pipeline:\n got: %s\nwant: %s", got, expected)
	}
}
