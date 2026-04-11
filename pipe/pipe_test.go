package pipe_test

import (
	"strings"
	"testing"

	"github.com/natalie-o-perret/gof/pipe"
)

func TestPipe2(t *testing.T) {
	got := pipe.Pipe2("hello", strings.ToUpper, func(s string) int { return len(s) })
	if got != 5 {
		t.Fatalf("got %d", got)
	}
}

func TestPipe3(t *testing.T) {
	got := pipe.Pipe3(" hello ", strings.TrimSpace, strings.ToUpper, func(s string) string { return s + "!" })
	if got != "HELLO!" {
		t.Fatalf("got %s", got)
	}
}

func TestPipeN(t *testing.T) {
	// arbitrary number of same-type steps
	got := pipe.PipeN(
		" hello world ",
		strings.TrimSpace,
		strings.ToUpper,
		func(s string) string { return s + "!" },
		func(s string) string { return "[" + s + "]" },
		func(s string) string { return s + "!" },
	)
	if got != "[HELLO WORLD!]!" {
		t.Fatalf("got %s", got)
	}

	// zero steps: value passes through unchanged
	if pipe.PipeN(42) != 42 {
		t.Fatal("zero-step PipeN should return value unchanged")
	}
}
