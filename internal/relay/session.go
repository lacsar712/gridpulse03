package relay

import "sync"

type LiveSession struct {
	mu    sync.Mutex
	label string
	n     int
}

func NewLiveSession(label string) *LiveSession {
	return &LiveSession{label: label}
}

func (s *LiveSession) Inc() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.n++
}

func (s *LiveSession) Value() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.n
}

func CloneSession(s *LiveSession) *LiveSession {
	s.mu.Lock()
	defer s.mu.Unlock()
	return &LiveSession{label: s.label, n: s.n}
}
