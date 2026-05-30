package seq_test

import (
	"cmp"
	"fmt"
	"strconv"

	"github.com/natalie-o-perret/go-functionalish/option"
	"github.com/natalie-o-perret/go-functionalish/seq"
)

func ExampleSeq_Map() {
	squares := seq.Range(1, 6).Map(func(n int) int { return n * n }).ToSlice()
	fmt.Println(squares)
	// Output: [1 4 9 16 25]
}

func ExampleSeq_Fold() {
	sum := seq.Range(1, 6).Fold(0, func(acc, v int) int { return acc + v })
	fmt.Println(sum)
	// Output: 15
}

func ExampleZip() {
	pairs := seq.Zip(
		seq.OfSlice([]int{1, 2, 3}),
		seq.OfSlice([]string{"a", "b", "c"}),
	).ToSlice()
	fmt.Println(pairs)
	// Output: [{1 a} {2 b} {3 c}]
}

func ExampleSeq_Cycle() {
	s := seq.OfSlice([]string{"ping", "pong"}).Cycle().Truncate(5).ToSlice()
	fmt.Println(s)
	// Output: [ping pong ping pong ping]
}

func ExampleSeq_Choose() {
	parsed := seq.OfSlice([]string{"2015", "bad", "2018", "nope", "2020"}).Choose(
		func(s string) option.Option[int] {
			n, err := strconv.Atoi(s)
			if err != nil {
				return option.None[int]()
			}
			return option.Some(n)
		},
	).SortWith(cmp.Compare).ToSlice()
	fmt.Println(parsed)
	// Output: [2015 2018 2020]
}

func ExampleThen() {
	type Car struct {
		Owner, Model string
		Year         int
	}
	cars := []Car{
		{"Alice", "Toyota", 2012}, {"Bob", "Honda", 2016},
		{"Charlie", "Ford", 2018}, {"Diana", "BMW", 2015},
		{"Eve", "Tesla", 2020}, {"Frank", "Honda", 2013},
	}

	owners := seq.Then(
		seq.OfSlice(cars).Filter(func(c Car) bool { return c.Year >= 2015 }),
		seq.MapFn(func(c Car) string { return c.Owner }),
	).ToSlice()
	fmt.Println(owners)
	// Output: [Bob Charlie Diana Eve]
}
