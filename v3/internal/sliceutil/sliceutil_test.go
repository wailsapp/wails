package sliceutil

import (
	"reflect"
	"testing"
)

func TestUnique(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		want  []int
	}{
		{
			name:  "no duplicates",
			slice: []int{1, 2, 3},
			want:  []int{1, 2, 3},
		},
		{
			name:  "with duplicates",
			slice: []int{1, 2, 2, 3, 3, 3},
			want:  []int{1, 2, 3},
		},
		{
			name:  "all duplicates",
			slice: []int{1, 1, 1},
			want:  []int{1},
		},
		{
			name:  "preserves order",
			slice: []int{3, 1, 2, 1, 3, 2},
			want:  []int{3, 1, 2},
		},
		{
			name:  "single element",
			slice: []int{1},
			want:  []int{1},
		},
		{
			name:  "empty slice",
			slice: []int{},
			want:  []int{},
		},
		{
			name:  "nil slice",
			slice: nil,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Unique(tt.slice)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnique_Strings(t *testing.T) {
	slice := []string{"a", "b", "a", "c", "b"}
	got := Unique(slice)
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Unique() = %v, want %v", got, want)
	}
}

func TestUnique_DoesNotModifyOriginal(t *testing.T) {
	original := []int{1, 2, 2, 3}
	originalCopy := make([]int, len(original))
	copy(originalCopy, original)

	_ = Unique(original)

	if !reflect.DeepEqual(original, originalCopy) {
		t.Errorf("Unique() modified original slice: got %v, want %v", original, originalCopy)
	}
}

func TestFindMapKey(t *testing.T) {
	tests := []struct {
		name      string
		m         map[string]int
		val       int
		wantKey   string
		wantFound bool
	}{
		{
			name:      "find existing value",
			m:         map[string]int{"a": 1, "b": 2, "c": 3},
			val:       2,
			wantKey:   "b",
			wantFound: true,
		},
		{
			name:      "value not found",
			m:         map[string]int{"a": 1, "b": 2},
			val:       3,
			wantKey:   "",
			wantFound: false,
		},
		{
			name:      "empty map",
			m:         map[string]int{},
			val:       1,
			wantKey:   "",
			wantFound: false,
		},
		{
			name:      "nil map",
			m:         nil,
			val:       1,
			wantKey:   "",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotFound := FindMapKey(tt.m, tt.val)
			if gotFound != tt.wantFound {
				t.Errorf("FindMapKey() found = %v, want %v", gotFound, tt.wantFound)
			}
			if gotFound && gotKey != tt.wantKey {
				t.Errorf("FindMapKey() key = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}

func TestFindMapKey_DuplicateValues(t *testing.T) {
	// When multiple keys have the same value, any matching key is acceptable
	m := map[string]int{"a": 1, "b": 1, "c": 2}
	key, found := FindMapKey(m, 1)
	if !found {
		t.Error("FindMapKey() should find a key")
	}
	if key != "a" && key != "b" {
		t.Errorf("FindMapKey() = %v, want 'a' or 'b'", key)
	}
}

func TestFindMapKey_IntKeys(t *testing.T) {
	m := map[int]string{1: "one", 2: "two", 3: "three"}
	key, found := FindMapKey(m, "two")
	if !found || key != 2 {
		t.Errorf("FindMapKey() = (%v, %v), want (2, true)", key, found)
	}
}

// Benchmarks

func BenchmarkUnique(b *testing.B) {
	slice := []int{1, 2, 3, 1, 2, 3, 4, 5, 4, 5, 6, 7, 8, 9, 10}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Unique(slice)
	}
}

func BenchmarkFindMapKey(b *testing.B) {
	m := make(map[string]int)
	for i := 0; i < 100; i++ {
		m[string(rune('a'+i))] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FindMapKey(m, 50)
	}
}
