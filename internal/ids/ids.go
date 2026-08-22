package ids

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

func NewSlug(prefix string) string {
	prefix = strings.TrimSpace(strings.ToLower(prefix))
	if prefix == "" {
		prefix = "id"
	}
	var buf [4]byte
	_, _ = rand.Read(buf[:])
	return fmt.Sprintf("%s-%d-%s", prefix, time.Now().Unix()%100000, hex.EncodeToString(buf[:]))
}

func ShortID(n int) string {
	if n < 4 {
		n = 4
	}
	if n > 32 {
		n = 32
	}
	buf := make([]byte, n)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)[:n]
}

func Key(parts ...string) string {
	return strings.Join(parts, ":")
}
