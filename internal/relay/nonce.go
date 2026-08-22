package relay

import "sync"

type NonceBook struct {
	mu   sync.Mutex
	seen map[string]struct{}
}

func NewNonceBook() *NonceBook {
	return &NonceBook{seen: make(map[string]struct{})}
}

func (b *NonceBook) CheckAndRemember(token string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if token == "" {
		return false
	}
	if _, ok := b.seen[token]; ok {
		return false
	}
	b.seen[token] = struct{}{}
	return true
}
