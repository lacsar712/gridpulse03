package kvbag

import "testing"

func TestSetGet(t *testing.T) {
	b := New()
	b.Set("k", "v")
	got, ok := b.Get("k")
	if !ok || got != "v" {
		t.Fatalf("Get=%q ok=%v", got, ok)
	}
	if b.Len() != 1 {
		t.Fatalf("len=%d", b.Len())
	}
}
