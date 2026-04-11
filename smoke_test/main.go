package main

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/natalie-o-perret/gof/enum"
	"github.com/natalie-o-perret/gof/option"
	"github.com/natalie-o-perret/gof/pipe"
	"github.com/natalie-o-perret/gof/result"
)

type Car struct{ Year int; Owner, Model string }

func main() {
	cars := []Car{
		{2012, "Alice", "Toyota"}, {2016, "Bob", "Honda"},
		{2018, "Charlie", "Ford"}, {2015, "Diana", "BMW"},
		{2020, "Eve", "Tesla"},   {2013, "Frank", "Honda"},
	}


	// ── enum ─────────────────────────────────────────────────────────────────
	owners := enum.Map(
		enum.From(cars).Filter(func(c Car) bool { return c.Year >= 2015 }),
		func(c Car) string { return c.Owner },
	).ToSlice()
	fmt.Println("owners >= 2015 :", owners)

	sorted := enum.SortAscBy(enum.From(cars), func(c Car) string { return c.Model }).
		FirstOption()
	fmt.Println("first by model :", sorted.UnwrapOr(Car{}).Model)

	squares := enum.Map(enum.Range(1, 6), func(n int) int { return n * n }).ToSlice()
	fmt.Println("squares 1-5    :", squares)

	pattern := enum.From([]string{"ping", "pong"}).Cycle().Take(5).ToSlice()
	fmt.Println("cycle×5        :", pattern)

	pairs := enum.Zip(enum.From([]int{1, 2, 3}), enum.From([]string{"a", "b", "c"})).ToSlice()
	fmt.Println("zip            :", pairs)

	sum := enum.Reduce(enum.Range(1, 6), 0, func(acc, v int) int { return acc + v })
	fmt.Println("sum 1-5        :", sum)

	// ── option ───────────────────────────────────────────────────────────────
	name := option.Some("alice")
	upper := option.Map(name, strings.ToUpper)
	fmt.Println("option map     :", upper.UnwrapOr("?"))
	fmt.Println("option none    :", option.Map(option.None[string](), strings.ToUpper).UnwrapOr("none"))

	// ── result ───────────────────────────────────────────────────────────────
	res := result.Try(func() (int, error) { return 42, nil })
	doubled := result.Map(res, func(n int) int { return n * 2 })
	fmt.Println("result map     :", doubled.UnwrapOr(0))
	fmt.Println("result→option  :", doubled.ToOption().UnwrapOr(0))

	// ── pipe ─────────────────────────────────────────────────────────────────
	out := pipe.Pipe3(" Hello World ", strings.TrimSpace, strings.ToLower,
		func(s string) string { return strings.ReplaceAll(s, " ", "-") })
	fmt.Println("pipe           :", out)

	// ── cross-cutting: FilterMap (Result-based filtering) ────────────────────
	parsed := enum.FilterMap(
		enum.From([]string{"2015", "bad", "2018", "nope", "2020"}),
		func(s string) result.Result[int, error] {
			var n int
			_, err := fmt.Sscanf(s, "%d", &n)
			return result.Try(func() (int, error) { return n, err })
		},
	).SortBy(cmp.Compare).ToSlice()
	fmt.Println("filtermap      :", parsed)
}

