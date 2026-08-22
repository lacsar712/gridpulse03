package ctxwait

import (
	"context"
	"testing"
	"time"
)

func TestUntilHonorsCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(30 * time.Millisecond)
		cancel()
	}()
	start := time.Now()
	err := Until(ctx, 600*time.Millisecond)
	elapsed := time.Since(start)
	if err == nil {
		t.Fatal("expected cancel error")
	}
	if elapsed > 250*time.Millisecond {
		t.Fatalf("Until ignored cancel, elapsed=%s", elapsed)
	}
}

func TestUntilCompletes(t *testing.T) {
	ctx := context.Background()
	if err := Until(ctx, 20*time.Millisecond); err != nil {
		t.Fatal(err)
	}
}
