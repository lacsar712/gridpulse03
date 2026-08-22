package pathsafe

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestJoinUnderRejectsEscape(t *testing.T) {
	root := filepath.Join("data", "store")
	_, err := JoinUnder(root, filepath.Join("..", "secret.txt"))
	if err == nil {
		t.Fatal("expected escape error")
	}
}

func TestJoinUnderOK(t *testing.T) {
	root := filepath.Join("data", "store")
	got, err := JoinUnder(root, filepath.Join("a", "b.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "b.txt") {
		t.Fatalf("got=%q", got)
	}
}
