package option_test

import (
	"fmt"
	"strings"

	"github.com/natalie-o-perret/go-functionalish/option"
)

func ExampleSome() {
	name := option.Some("alice")
	upper := name.Map(strings.ToUpper)
	fmt.Println(upper.UnwrapOr("?"))
	// Output: ALICE
}

func ExampleNone() {
	none := option.None[string]()
	upper := none.Map(strings.ToUpper)
	fmt.Println(upper.UnwrapOr("none"))
	// Output: none
}

func ExampleOption_Map() {
	res := option.Some(21).Map(func(n int) int { return n * 2 })
	fmt.Println(res.UnwrapOr(0))
	// Output: 42
}
