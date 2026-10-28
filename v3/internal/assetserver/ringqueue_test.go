package assetserver

import "testing"

func TestRingqueueReleasesRemovedValues(t *testing.T) {
	q := newRingqueue[*int](4)
	first, second := 1, 2
	q.Add(&first)
	q.Add(&second)
	if value, ok := q.Remove(); !ok || value != &first {
		t.Fatalf("Remove = %v, %v; want first value", value, ok)
	}
	if q.nodes[0] != nil {
		t.Error("removed slot still holds a reference")
	}
	if value, ok := q.Peek(); !ok || value != &second {
		t.Fatalf("Peek = %v, %v; want second value", value, ok)
	}
	q.Remove()
	for _, value := range q.nodes {
		if value != nil {
			t.Error("empty queue still holds a reference")
		}
	}
	if value, ok := q.Remove(); ok || value != nil {
		t.Errorf("empty Remove = %v, %v; want nil, false", value, ok)
	}
}

func TestRingqueueFIFOAfterWrapAndResize(t *testing.T) {
	q := newRingqueue[int](2)
	for i := 0; i < 32; i++ {
		q.Add(i)
	}
	for want := 0; want < 7; want++ {
		if value, ok := q.Remove(); !ok || value != want {
			t.Fatalf("Remove = %d, %v; want %d", value, ok, want)
		}
	}
	for i := 32; i < 64; i++ {
		q.Add(i)
	}
	for want := 7; want < 64; want++ {
		if value, ok := q.Remove(); !ok || value != want {
			t.Fatalf("Remove = %d, %v; want %d", value, ok, want)
		}
	}
	if q.Len() != 0 {
		t.Errorf("Len = %d, want 0", q.Len())
	}
	if _, ok := q.Peek(); ok {
		t.Error("Peek succeeded on an empty queue")
	}
}
