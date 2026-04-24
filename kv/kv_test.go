package kv_test
import (
"sort"
"testing"
"github.com/natalie-o-perret/go-functionalish/kv"
"github.com/natalie-o-perret/go-functionalish/seq"
)
func TestOf_Collect(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	got := kv.Collect(kv.Of(m))
	if len(got) != 3 || got["a"] != 1 || got["b"] != 2 || got["c"] != 3 {
		t.Fatalf("got %v", got)
	}
}
func TestKeys(t *testing.T) {
	m := map[string]int{"x": 10, "y": 20}
	keys := kv.Keys(kv.Of(m)).ToSlice()
	sort.Strings(keys)
	if len(keys) != 2 || keys[0] != "x" || keys[1] != "y" {
		t.Fatalf("got %v", keys)
	}
}
func TestValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	vals := kv.Values(kv.Of(m)).ToSlice()
	sort.Ints(vals)
	if len(vals) != 2 || vals[0] != 1 || vals[1] != 2 {
		t.Fatalf("got %v", vals)
	}
}
func TestToSeq_FromSeq(t *testing.T) {
	m := map[string]int{"p": 7, "q": 8}
	pairs := kv.ToSeq(kv.Of(m)).ToSlice()
	if len(pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(pairs))
	}
	back := kv.Collect(kv.FromSeq(seq.OfSlice(pairs)))
	if back["p"] != 7 || back["q"] != 8 {
		t.Fatalf("round-trip failed: %v", back)
	}
}
func TestMapValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	doubled := kv.Collect(kv.MapValues(kv.Of(m), func(v int) int { return v * 2 }))
	if doubled["a"] != 2 || doubled["b"] != 4 {
		t.Fatalf("got %v", doubled)
	}
}
func TestMapKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	upper := kv.Collect(kv.MapKeys(kv.Of(m), func(k string) string { return k + k }))
	if upper["aa"] != 1 || upper["bb"] != 2 {
		t.Fatalf("got %v", upper)
	}
}
func TestFilter(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	got := kv.Collect(kv.Filter(kv.Of(m), func(_ string, v int) bool { return v > 1 }))
	if len(got) != 2 || got["b"] != 2 || got["c"] != 3 {
		t.Fatalf("got %v", got)
	}
}
func TestFold(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	sum := kv.Fold(kv.Of(m), 0, func(acc int, _ string, v int) int { return acc + v })
	if sum != 6 {
		t.Fatalf("got %d", sum)
	}
}
func TestContainsKey(t *testing.T) {
	m := map[string]int{"hello": 1}
	if !kv.ContainsKey(kv.Of(m), "hello") {
		t.Fatal("expected true")
	}
	if kv.ContainsKey(kv.Of(m), "world") {
		t.Fatal("expected false")
	}
}
