package errwrap

import (
	"errors"
	"testing"
)

func TestWrapDeniedIs(t *testing.T) {
	err := WrapDenied("update")
	if !errors.Is(err, ErrDenied) {
		t.Fatalf("errors.Is lost sentinel: %v", err)
	}
	if !IsDenied(err) {
		t.Fatalf("IsDenied=false for %v", err)
	}
}

func TestWrapfNil(t *testing.T) {
	if Wrapf(nil, "x") != nil {
		t.Fatal("nil should stay nil")
	}
}
