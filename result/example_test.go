package result_test

import (
	"fmt"

	"github.com/natalie-o-perret/go-functionalish/result"
)

func ExampleTry() {
	res := result.Try(func() (int, error) { return 42, nil })
	doubled := result.Map(res, func(n int) int { return n * 2 })
	fmt.Println(doubled.UnwrapOr(0))
	// Output: 84
}

func ExampleResult_ToOption() {
	res := result.Try(func() (int, error) { return 42, nil })
	doubled := result.Map(res, func(n int) int { return n * 2 })
	fmt.Println(doubled.ToOption().UnwrapOr(0))
	// Output: 84
}
