package precopy

import "testing"

func TestClonePrefixIndependent(t *testing.T) {
	src := []byte("abcdef")
	got := ClonePrefix(src, 3)
	if string(got) != "abc" {
		t.Fatalf("prefix=%q", got)
	}
	got[0] = 'Z'
	if src[0] != 'a' {
		t.Fatalf("ClonePrefix must not alias src, src became %q", src)
	}
}

func TestCloneAllIndependent(t *testing.T) {
	src := []byte{1, 2, 3}
	got := CloneAll(src)
	got[1] = 9
	if src[1] != 2 {
		t.Fatalf("CloneAll aliased src: %v", src)
	}
}
