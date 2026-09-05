package lo

import (
	"errors"
	"sort"
	"testing"
)

// ----- Associate -----

func TestAssociate_basic(t *testing.T) {
	in := []string{"a", "bb", "ccc"}
	got := Associate(in, func(s string) (string, int) { return s, len(s) })
	want := map[string]int{"a": 1, "bb": 2, "ccc": 3}
	if len(got) != len(want) {
		t.Fatalf("got len %d want %d", len(got), len(want))
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("key %q: got %d want %d", k, got[k], v)
		}
	}
}

func TestAssociate_empty(t *testing.T) {
	got := Associate([]int{}, func(n int) (int, int) { return n, n })
	if len(got) != 0 {
		t.Fatalf("expected empty map, got %v", got)
	}
}

func TestAssociate_nil(t *testing.T) {
	var in []int
	got := Associate(in, func(n int) (int, int) { return n, n })
	if len(got) != 0 {
		t.Fatalf("expected empty map for nil slice, got %v", got)
	}
}

func TestAssociate_duplicateKey(t *testing.T) {
	// Later key wins
	in := []int{1, 2, 3}
	got := Associate(in, func(n int) (string, int) { return "k", n })
	if got["k"] != 3 {
		t.Errorf("expected last value 3, got %d", got["k"])
	}
}

// ----- Contains -----

func TestContains_found(t *testing.T) {
	if !Contains([]int{1, 2, 3}, 2) {
		t.Error("expected true")
	}
}

func TestContains_notFound(t *testing.T) {
	if Contains([]int{1, 2, 3}, 99) {
		t.Error("expected false")
	}
}

func TestContains_empty(t *testing.T) {
	if Contains([]string{}, "x") {
		t.Error("expected false on empty slice")
	}
}

func TestContains_nil(t *testing.T) {
	var s []string
	if Contains(s, "x") {
		t.Error("expected false on nil slice")
	}
}

func TestContains_strings(t *testing.T) {
	if !Contains([]string{"foo", "bar"}, "bar") {
		t.Error("expected true")
	}
}

// ----- ContainsBy -----

func TestContainsBy_found(t *testing.T) {
	if !ContainsBy([]int{1, 2, 3}, func(n int) bool { return n > 2 }) {
		t.Error("expected true")
	}
}

func TestContainsBy_notFound(t *testing.T) {
	if ContainsBy([]int{1, 2, 3}, func(n int) bool { return n > 10 }) {
		t.Error("expected false")
	}
}

func TestContainsBy_empty(t *testing.T) {
	if ContainsBy([]int{}, func(n int) bool { return true }) {
		t.Error("expected false on empty slice")
	}
}

func TestContainsBy_nil(t *testing.T) {
	var s []string
	if ContainsBy(s, func(x string) bool { return true }) {
		t.Error("expected false on nil slice")
	}
}

// ----- Find -----

func TestFind_found(t *testing.T) {
	val, ok := Find([]int{10, 20, 30}, func(n int) bool { return n == 20 })
	if !ok || val != 20 {
		t.Errorf("got (%d, %v) want (20, true)", val, ok)
	}
}

func TestFind_firstMatch(t *testing.T) {
	val, ok := Find([]int{1, 2, 3, 2}, func(n int) bool { return n == 2 })
	if !ok || val != 2 {
		t.Errorf("got (%d, %v)", val, ok)
	}
}

func TestFind_notFound(t *testing.T) {
	val, ok := Find([]int{1, 2, 3}, func(n int) bool { return n == 99 })
	if ok || val != 0 {
		t.Errorf("expected (0, false), got (%d, %v)", val, ok)
	}
}

func TestFind_empty(t *testing.T) {
	_, ok := Find([]string{}, func(s string) bool { return true })
	if ok {
		t.Error("expected false on empty slice")
	}
}

func TestFind_nil(t *testing.T) {
	var s []string
	_, ok := Find(s, func(x string) bool { return true })
	if ok {
		t.Error("expected false on nil slice")
	}
}

func TestFind_zeroValueOnMiss(t *testing.T) {
	val, ok := Find([]string{"a", "b"}, func(s string) bool { return s == "z" })
	if ok || val != "" {
		t.Errorf("expected (\"\", false), got (%q, %v)", val, ok)
	}
}

// ----- Keys -----

func TestKeys_basic(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keys := Keys(m)
	sort.Strings(keys)
	want := []string{"a", "b", "c"}
	if len(keys) != len(want) {
		t.Fatalf("got %v want %v", keys, want)
	}
	for i, k := range want {
		if keys[i] != k {
			t.Errorf("index %d: got %q want %q", i, keys[i], k)
		}
	}
}

func TestKeys_empty(t *testing.T) {
	keys := Keys(map[int]string{})
	if len(keys) != 0 {
		t.Fatalf("expected no keys, got %v", keys)
	}
}

func TestKeys_nil(t *testing.T) {
	var m map[string]bool
	keys := Keys(m)
	if len(keys) != 0 {
		t.Fatalf("expected no keys for nil map, got %v", keys)
	}
}

func TestKeys_singleEntry(t *testing.T) {
	keys := Keys(map[int]int{42: 99})
	if len(keys) != 1 || keys[0] != 42 {
		t.Errorf("expected [42], got %v", keys)
	}
}

// ----- Must -----

func TestMust_noError(t *testing.T) {
	val := Must(42, nil)
	if val != 42 {
		t.Errorf("expected 42, got %d", val)
	}
}

func TestMust_panicOnError(t *testing.T) {
	sentinel := errors.New("boom")
	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic, got none")
		}
		if r != sentinel {
			t.Errorf("expected sentinel error in panic, got %v", r)
		}
	}()
	Must(0, sentinel)
}

func TestMust_stringType(t *testing.T) {
	val := Must("hello", nil)
	if val != "hello" {
		t.Errorf("expected \"hello\", got %q", val)
	}
}

// ----- Ternary -----

func TestTernary_true(t *testing.T) {
	if Ternary(true, "yes", "no") != "yes" {
		t.Error("expected yes")
	}
}

func TestTernary_false(t *testing.T) {
	if Ternary(false, "yes", "no") != "no" {
		t.Error("expected no")
	}
}

func TestTernary_int(t *testing.T) {
	if Ternary(1 > 2, 100, 200) != 200 {
		t.Error("expected 200")
	}
}

func TestTernary_bool(t *testing.T) {
	if !Ternary(true, true, false) {
		t.Error("expected true")
	}
}

// ----- Without -----

func TestWithout_basic(t *testing.T) {
	got := Without([]int{1, 2, 3, 4, 5}, 2, 4)
	want := []int{1, 3, 5}
	if len(got) != len(want) {
		t.Fatalf("got %v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %d want %d", i, got[i], want[i])
		}
	}
}

func TestWithout_empty(t *testing.T) {
	got := Without([]int{}, 1)
	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

func TestWithout_nil(t *testing.T) {
	var s []int
	got := Without(s, 1)
	if len(got) != 0 {
		t.Fatalf("expected empty for nil slice, got %v", got)
	}
}

func TestWithout_noExcludes(t *testing.T) {
	in := []int{1, 2, 3}
	got := Without(in)
	if len(got) != len(in) {
		t.Fatalf("got %v want %v", got, in)
	}
	for i := range in {
		if got[i] != in[i] {
			t.Errorf("index %d: got %d want %d", i, got[i], in[i])
		}
	}
}

func TestWithout_excludeAll(t *testing.T) {
	got := Without([]string{"a", "b", "c"}, "a", "b", "c")
	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

func TestWithout_duplicatesInCollection(t *testing.T) {
	got := Without([]int{1, 1, 2, 1, 3}, 1)
	want := []int{2, 3}
	if len(got) != len(want) {
		t.Fatalf("got %v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %d want %d", i, got[i], want[i])
		}
	}
}

func TestWithout_excludeNotPresent(t *testing.T) {
	in := []int{1, 2, 3}
	got := Without(in, 99)
	if len(got) != len(in) {
		t.Fatalf("got %v want %v", got, in)
	}
	for i := range in {
		if got[i] != in[i] {
			t.Errorf("index %d: got %d want %d", i, got[i], in[i])
		}
	}
}

func TestWithout_multipleExcludes(t *testing.T) {
	got := Without([]string{"x", "y", "z", "y", "x"}, "x", "y")
	want := []string{"z"}
	if len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("got %v want %v", got, want)
	}
}
