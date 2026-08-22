package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

type Session struct {
	Name string
	Role string
	ID   int64
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (s *Session) DisplayName() string {
	if s == nil {
		return ""
	}
	return s.Name
}

func (s *Session) IsAdmin() bool {
	if s == nil {
		return false
	}
	return s.Role == "admin"
}

func ParseBearer(h string) string {
	h = strings.TrimSpace(h)
	if strings.HasPrefix(strings.ToLower(h), "bearer ") {
		return strings.TrimSpace(h[7:])
	}
	return h
}
