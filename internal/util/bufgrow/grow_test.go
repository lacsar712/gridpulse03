package bufgrow

import "testing"

func TestExtendNoWriteThrough(t *testing.T) {
	src := make([]byte, 3, 16)
	copy(src, []byte{1, 2, 3})
	got := Extend(src, 9)
	if len(got) != 4 || got[3] != 9 {
		t.Fatalf("got=%v", got)
	}
	got[0] = 77
	if src[0] != 1 {
		t.Fatalf("Extend wrote through to src: %v", src)
	}
}

func TestConcatIndependent(t *testing.T) {
	a := []byte{1, 2}
	b := []byte{3}
	got := Concat(a, b)
	got[0] = 9
	if a[0] != 1 {
		t.Fatalf("Concat aliased a: %v", a)
	}
}
