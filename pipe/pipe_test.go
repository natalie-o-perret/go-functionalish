package pipe_test

import (
	"strings"
	"testing"

	"github.com/natalie-o-perret/gof/pipe"
)

func TestPipe2(t *testing.T) {
	got := pipe.Pipe2("hello", strings.ToUpper, func(s string) int { return len(s) })
	if got != 5 { t.Fatalf("got %d", got) }
}

func TestPipe3(t *testing.T) {
	got := pipe.Pipe3(" hello ", strings.TrimSpace, strings.ToUpper, func(s string) string { return s + "!" })
	if got != "HELLO!" { t.Fatalf("got %s", got) }
}

