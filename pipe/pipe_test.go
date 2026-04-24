package pipe_test

import (
	"strings"
	"testing"

	"github.com/natalie-o-perret/go-functionalish/pipe"
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

func TestPipeEndoN(t *testing.T) {
	got := pipe.PipeEndoN(
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
	if pipe.PipeEndoN(42) != 42 {
		t.Fatal("zero-step PipeEndoN should return value unchanged")
	}
}

func TestTap(t *testing.T) {
	var seen string
	got := pipe.Pipe3(
		" hello ",
		strings.TrimSpace,
		pipe.Tap(func(s string) { seen = s }),
		strings.ToUpper,
	)
	if got != "HELLO" {
		t.Fatalf("got %s", got)
	}
	if seen != "hello" {
		t.Fatalf("tap saw %q, want %q", seen, "hello")
	}
}
