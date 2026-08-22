package store

import (
	"context"
	"time"

	"example.com/gridpulse/internal/model"
)

func (s *Store) AddAudit(ctx context.Context, actor, action, detail string) error {
	_, err := s.DB.ExecContext(ctx, "INSERT INTO audits(actor,action,detail,created_at) VALUES(?,?,?,?)",
		actor, action, detail, time.Now().UTC().Format(time.RFC3339))
	return err
}

func (s *Store) ListAudit(ctx context.Context, limit int) ([]model.Audit, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := s.DB.QueryContext(ctx, "SELECT id,actor,action,detail,created_at FROM audits ORDER BY id DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Audit
	for rows.Next() {
		var a model.Audit
		var created string
		if err := rows.Scan(&a.ID, &a.Actor, &a.Action, &a.Detail, &created); err != nil {
			return nil, err
		}
		a.CreatedAt, _ = time.Parse(time.RFC3339, created)
		out = append(out, a)
	}
	return out, rows.Err()
}

func (s *Store) SetSetting(ctx context.Context, k, v string) error {
	_, err := s.DB.ExecContext(ctx, "INSERT INTO settings(k,v) VALUES(?,?) ON CONFLICT(k) DO UPDATE SET v=excluded.v", k, v)
	return err
}

func (s *Store) GetSetting(ctx context.Context, k string) (string, error) {
	var v string
	err := s.DB.QueryRowContext(ctx, "SELECT v FROM settings WHERE k=?", k).Scan(&v)
	return v, err
}
