package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"example.com/gridpulse/internal/model"
)

func (s *Store) CreateUser(ctx context.Context, name, tokenHash, role string) (model.User, error) {
	res, err := s.DB.ExecContext(ctx, "INSERT INTO users(name, token_hash, role, created_at) VALUES(?,?,?,?)", name, tokenHash, role, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return model.User{}, err
	}
	id, _ := res.LastInsertId()
	return s.GetUser(ctx, id)
}

func (s *Store) GetUser(ctx context.Context, id int64) (model.User, error) {
	var u model.User
	var created string
	err := s.DB.QueryRowContext(ctx, "SELECT id,name,role,created_at FROM users WHERE id=?", id).Scan(&u.ID, &u.Name, &u.Role, &created)
	if err == sql.ErrNoRows {
		return model.User{}, fmt.Errorf("user %d: %w", id, err)
	}
	if err != nil {
		return model.User{}, err
	}
	u.CreatedAt, _ = time.Parse(time.RFC3339, created)
	return u, nil
}

func (s *Store) UserByTokenHash(ctx context.Context, hash string) (model.User, error) {
	var u model.User
	var created string
	err := s.DB.QueryRowContext(ctx, "SELECT id,name,role,created_at FROM users WHERE token_hash=?", hash).Scan(&u.ID, &u.Name, &u.Role, &created)
	if err != nil {
		return model.User{}, err
	}
	u.CreatedAt, _ = time.Parse(time.RFC3339, created)
	return u, nil
}

func (s *Store) ListUsers(ctx context.Context) ([]model.User, error) {
	rows, err := s.DB.QueryContext(ctx, "SELECT id,name,role,created_at FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.User
	for rows.Next() {
		var u model.User
		var created string
		if err := rows.Scan(&u.ID, &u.Name, &u.Role, &created); err != nil {
			return nil, err
		}
		u.CreatedAt, _ = time.Parse(time.RFC3339, created)
		out = append(out, u)
	}
	return out, rows.Err()
}
