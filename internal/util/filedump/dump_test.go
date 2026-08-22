package filedump

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteAllPersists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")
	body := "hello-buffered-body-0123456789"
	if err := WriteAll(path, body); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != body {
		t.Fatalf("file=%q want %q", got, body)
	}
}
