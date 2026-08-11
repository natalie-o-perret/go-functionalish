package taskseq_test

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/natalie-o-perret/go-functionalish/taskseq"
)

func TestPipelineIsLazyAndOrdered(t *testing.T) {
	mapped := 0
	pipeline := taskseq.FromSeq(slices.Values([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})).
		MapAsync(func(_ context.Context, value int) (int, error) {
			mapped++
			return value * 2, nil
		}).
		FilterAsync(func(_ context.Context, value int) (bool, error) {
			return value%4 == 0, nil
		}).
		Take(2)

	if mapped != 0 {
		t.Fatalf("pipeline ran before consumption: mapped %d elements", mapped)
	}

	for run := 1; run <= 2; run++ {
		got, err := pipeline.ToSlice(context.Background())
		if err != nil {
			t.Fatalf("run %d: ToSlice returned an error: %v", run, err)
		}
		if want := []int{4, 8}; !slices.Equal(got, want) {
			t.Fatalf("run %d: got %v, want %v", run, got, want)
		}
	}
	if mapped != 8 {
		t.Fatalf("mapped %d elements, want 8", mapped)
	}
}

func TestMapAsyncStopsAtError(t *testing.T) {
	wantErr := errors.New("map failed")
	got, err := taskseq.FromSeq(slices.Values([]int{1, 2, 3, 4, 5})).
		MapAsync(func(_ context.Context, value int) (int, error) {
			if value == 3 {
				return 0, wantErr
			}
			return value * 10, nil
		}).
		ToSlice(context.Background())

	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
	if want := []int{10, 20}; !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestFilterAsyncStopsAtError(t *testing.T) {
	wantErr := errors.New("filter failed")
	got, err := taskseq.FromSeq(slices.Values([]int{1, 2, 3, 4, 5})).
		FilterAsync(func(_ context.Context, value int) (bool, error) {
			if value == 3 {
				return false, wantErr
			}
			return value%2 == 0, nil
		}).
		ToSlice(context.Background())

	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
	if want := []int{2}; !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestTakeDoesNotPullAnExtraElement(t *testing.T) {
	pulled := 0
	source := taskseq.TaskSeq[int](func(_ context.Context, yield func(int) bool) error {
		for value := 1; value <= 5; value++ {
			pulled++
			if !yield(value) {
				return nil
			}
		}
		return nil
	})

	got, err := source.Take(2).ToSlice(context.Background())
	if err != nil {
		t.Fatalf("ToSlice returned an error: %v", err)
	}
	if want := []int{1, 2}; !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if pulled != 2 {
		t.Fatalf("pulled %d elements, want 2", pulled)
	}
}

func TestTakeNonPositiveDoesNotConsumeSource(t *testing.T) {
	for _, n := range []int{-1, 0} {
		t.Run("n", func(t *testing.T) {
			pulled := 0
			source := taskseq.TaskSeq[int](func(_ context.Context, _ func(int) bool) error {
				pulled++
				return nil
			})

			got, err := source.Take(n).ToSlice(context.Background())
			if err != nil {
				t.Fatalf("ToSlice returned an error: %v", err)
			}
			if len(got) != 0 {
				t.Fatalf("got %v, want no values", got)
			}
			if pulled != 0 {
				t.Fatalf("source ran %d times, want 0", pulled)
			}
		})
	}
}

func TestForEachStopsAtError(t *testing.T) {
	wantErr := errors.New("action failed")
	pulled := 0
	source := taskseq.TaskSeq[int](func(_ context.Context, yield func(int) bool) error {
		for value := 1; value <= 5; value++ {
			pulled++
			if !yield(value) {
				return nil
			}
		}
		return nil
	})

	err := source.ForEach(context.Background(), func(_ context.Context, value int) error {
		if value == 2 {
			return wantErr
		}
		return nil
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
	if pulled != 2 {
		t.Fatalf("pulled %d elements, want 2", pulled)
	}
}

func TestCancellationOnFinalElementIsReturned(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	err := taskseq.FromSeq(slices.Values([]int{1})).ForEach(ctx, func(_ context.Context, _ int) error {
		calls++
		cancel()
		return nil
	})

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("got error %v, want %v", err, context.Canceled)
	}
	if calls != 1 {
		t.Fatalf("callback ran %d times, want 1", calls)
	}
}

func TestCancellationDoesNotPullAnotherElement(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	pulled := 0
	source := taskseq.TaskSeq[int](func(_ context.Context, yield func(int) bool) error {
		for value := 1; value <= 5; value++ {
			pulled++
			if !yield(value) {
				return nil
			}
		}
		return nil
	})

	err := source.ForEach(ctx, func(_ context.Context, _ int) error {
		cancel()
		return nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("got error %v, want %v", err, context.Canceled)
	}
	if pulled != 1 {
		t.Fatalf("pulled %d elements, want 1", pulled)
	}
}

func TestCancellationDuringEmptySourceIsReturned(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	source := taskseq.TaskSeq[int](func(_ context.Context, _ func(int) bool) error {
		cancel()
		return nil
	})

	_, err := source.ToSlice(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("got error %v, want %v", err, context.Canceled)
	}
}

func TestCancellationDuringSourceCleanupIsReturned(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	source := taskseq.TaskSeq[int](func(_ context.Context, yield func(int) bool) error {
		if !yield(1) {
			return nil
		}
		cancel()
		return nil
	})

	got, err := source.ToSlice(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("got error %v, want %v", err, context.Canceled)
	}
	if want := []int{1}; !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestOperationAndCleanupErrorsAreJoined(t *testing.T) {
	operationErr := errors.New("operation failed")
	cleanupErr := errors.New("cleanup failed")
	source := taskseq.TaskSeq[int](func(_ context.Context, yield func(int) bool) error {
		_ = yield(1)
		return cleanupErr
	})

	_, err := source.
		MapAsync(func(context.Context, int) (int, error) { return 0, operationErr }).
		ToSlice(context.Background())
	if !errors.Is(err, operationErr) {
		t.Fatalf("got error %v, want joined operation error", err)
	}
	if !errors.Is(err, cleanupErr) {
		t.Fatalf("got error %v, want joined cleanup error", err)
	}
}

func TestToSliceReturnsValuesBeforeSourceError(t *testing.T) {
	wantErr := errors.New("source failed")
	source := taskseq.TaskSeq[int](func(_ context.Context, yield func(int) bool) error {
		if !yield(1) {
			return nil
		}
		return wantErr
	})

	got, err := source.ToSlice(context.Background())
	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
	if want := []int{1}; !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
