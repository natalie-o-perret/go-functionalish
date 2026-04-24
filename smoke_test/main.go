package main

import (
	"cmp"
	"fmt"
	"strconv"
	"strings"

	"github.com/natalie-o-perret/go-functional-ish/option"
	"github.com/natalie-o-perret/go-functional-ish/pipe"
	"github.com/natalie-o-perret/go-functional-ish/result"
	"github.com/natalie-o-perret/go-functional-ish/seq"
)

type Car struct {
	Owner, Model string
	Year         int
}

func main() {
	cars := []Car{
		{"Alice", "Toyota", 2012}, {"Bob", "Honda", 2016},
		{"Charlie", "Ford", 2018}, {"Diana", "BMW", 2015},
		{"Eve", "Tesla", 2020}, {"Frank", "Honda", 2013},
	}

	// -- seq (fluent style with Then + curried helpers) ----------------------
	owners := seq.Then(
		seq.OfSlice(cars).Filter(func(c Car) bool { return c.Year >= 2015 }),
		seq.MapFn(func(c Car) string { return c.Owner }),
	).ToSlice()
	fmt.Println("owners >= 2015 :", owners)

	sorted := seq.Then(
		seq.OfSlice(cars),
		seq.SortByFn(func(c Car) string { return c.Model }),
	).TryHead()
	fmt.Println("first by model :", sorted.UnwrapOr(Car{}).Model)

	squares := seq.Map(seq.Range(1, 6), func(n int) int { return n * n }).ToSlice()
	fmt.Println("squares 1-5    :", squares)

	pattern := seq.OfSlice([]string{"ping", "pong"}).Cycle().Truncate(5).ToSlice()
	fmt.Println("cyclex5        :", pattern)

	pairs := seq.Zip(seq.OfSlice([]int{1, 2, 3}), seq.OfSlice([]string{"a", "b", "c"})).ToSlice()
	fmt.Println("zip            :", pairs)

	sum := seq.Fold(seq.Range(1, 6), 0, func(acc, v int) int { return acc + v })
	fmt.Println("sum 1-5        :", sum)

	// -- option ---------------------------------------------------------------
	name := option.Some("alice")
	upper := option.Map(name, strings.ToUpper)
	fmt.Println("option map     :", upper.UnwrapOr("?"))
	fmt.Println("option none    :", option.Map(option.None[string](), strings.ToUpper).UnwrapOr("none"))

	// -- result ---------------------------------------------------------------
	res := result.Try(func() (int, error) { return 42, nil })
	doubled := result.Map(res, func(n int) int { return n * 2 })
	fmt.Println("result map     :", doubled.UnwrapOr(0))
	fmt.Println("result=>option  :", doubled.ToOption().UnwrapOr(0))

	// -- pipe -----------------------------------------------------------------
	out := pipe.Pipe3(" Hello World ", strings.TrimSpace, strings.ToLower,
		func(s string) string { return strings.ReplaceAll(s, " ", "-") })
	fmt.Println("pipe           :", out)

	// -- cross-cutting: Choose (Option-based filtering, F#: Seq.choose) ------
	parsed := seq.Choose(
		seq.OfSlice([]string{"2015", "bad", "2018", "nope", "2020"}),
		func(s string) option.Option[int] {
			n, err := strconv.Atoi(s)
			if err != nil {
				return option.None[int]()
			}
			return option.Some(n)
		},
	).SortWith(cmp.Compare).ToSlice()
	fmt.Println("choose         :", parsed)
}
