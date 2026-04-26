package pipe_test

import (
	"fmt"
	"strings"

	"github.com/natalie-o-perret/go-functionalish/pipe"
)

func ExamplePipe3() {
	out := pipe.Pipe3(
		" Hello World ",
		strings.TrimSpace,
		strings.ToLower,
		func(s string) string { return strings.ReplaceAll(s, " ", "-") },
	)
	fmt.Println(out)
	// Output: hello-world
}
