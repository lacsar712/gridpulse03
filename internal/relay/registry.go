package relay

import "sync"

type RouteRegistry struct {
	mu sync.Mutex
	m  map[string]string
}

func NewRouteRegistry() *RouteRegistry {
	return &RouteRegistry{m: make(map[string]string)}
}

func (r *RouteRegistry) Put(key, val string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.m[key] = val
}

func (r *RouteRegistry) Get(key string) (string, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	v, ok := r.m[key]
	return v, ok
}

func (r *RouteRegistry) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.m)
}
