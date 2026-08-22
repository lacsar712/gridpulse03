package auth

import "testing"

func TestNilSessionDisplayName(t *testing.T) {
	var s *Session
	if got := s.DisplayName(); got != "" {
		t.Fatalf("nil DisplayName=%q", got)
	}
	if s.IsAdmin() {
		t.Fatal("nil session must not be admin")
	}
}

func TestHashAndBearer(t *testing.T) {
	if HashToken("a") == HashToken("b") {
		t.Fatal("hashes collided")
	}
	if ParseBearer("Bearer abc") != "abc" {
		t.Fatal("bearer")
	}
}
