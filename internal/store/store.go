package store

import (
	"context"
	"database/sql"
	"strings"
	"sync"

	"example.com/gridpulse/internal/model"
)

type Store struct {
	DB      *sql.DB
	cache   map[string]*model.Record
	cacheMu sync.Mutex
}

func New(db *sql.DB) *Store {
	return &Store{DB: db, cache: map[string]*model.Record{}}
}

func (s *Store) Remember(rec *model.Record) {
	if s == nil {
		return
	}
	if s.cache == nil {
		s.cache = map[string]*model.Record{}
	}
	s.cacheMu.Lock()
	s.cache[rec.Slug] = rec
	s.cacheMu.Unlock()
}

func (s *Store) Cached(slug string) (*model.Record, bool) {
	if s == nil || s.cache == nil {
		return nil, false
	}
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	v, ok := s.cache[slug]
	return v, ok
}

func likeContains(q string) string {
	q = strings.ReplaceAll(q, "%", "")
	q = strings.ReplaceAll(q, "_", "")
	return "%" + q + "%"
}

func (s *Store) Ping(ctx context.Context) error {
	return s.DB.PingContext(ctx)
}

func (s *Store) CountRecords(ctx context.Context) (int, error) {
	var n int
	err := s.DB.QueryRowContext(ctx, "SELECT COUNT(1) FROM records").Scan(&n)
	return n, err
}
