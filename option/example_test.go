package option_test

import (
	"fmt"
	"strings"

	"github.com/natalie-o-perret/go-functionalish/option"
)

func ExampleSome() {
	name := option.Some("alice")
	upper := option.Map(name, strings.ToUpper)
	fmt.Println(upper.UnwrapOr("?"))
	// Output: ALICE
}

func ExampleNone() {
	none := option.None[string]()
	upper := option.Map(none, strings.ToUpper)
	fmt.Println(upper.UnwrapOr("none"))
	// Output: none
}

func ExampleMap() {
	result := option.Map(option.Some(21), func(n int) int { return n * 2 })
	fmt.Println(result.UnwrapOr(0))
	// Output: 42
}

