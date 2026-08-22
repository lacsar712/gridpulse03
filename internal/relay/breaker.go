package relay

import "sync"

type TripBreaker struct {
	mu       sync.Mutex
	failures int
	open     bool
	limit    int
}

func NewTripBreaker(limit int) *TripBreaker {
	if limit <= 0 {
		limit = 2
	}
	return &TripBreaker{limit: limit}
}

func (b *TripBreaker) Fail() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.limit {
		b.open = true
	}
}

func (b *TripBreaker) Success() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.open = false
}

func (b *TripBreaker) Open() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.open
}

func RecordOutcome(b *TripBreaker, err error) {
	if err != nil {
		b.Fail()
		return
	}
	b.Success()
}
