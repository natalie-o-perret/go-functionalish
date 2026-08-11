package taskseq_test

import (
	"context"
	"fmt"

	"github.com/natalie-o-perret/go-functionalish/seq"
	"github.com/natalie-o-perret/go-functionalish/taskseq"
)

func ExampleSeq_MapAsync() {
	values, err := taskseq.FromSeq(seq.Range(1, 6)).
		MapAsync(func(_ context.Context, value int) (int, error) {
			return value * 2, nil
		}).
		FilterAsync(func(_ context.Context, value int) (bool, error) {
			return value%4 == 0, nil
		}).
		Take(2).
		ToSlice(context.Background())

	fmt.Println(values, err)
	// Output: [4 8] <nil>
}
