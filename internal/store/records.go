package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"example.com/gridpulse/internal/db"
	"example.com/gridpulse/internal/model"
)

func (s *Store) CreateRecord(ctx context.Context, rec model.Record) (model.Record, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	var id int64
	err := db.WithTx(ctx, s.DB, func(tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, "INSERT INTO records(slug,title,body,tags,owner_id,bytes,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?)",
			rec.Slug, rec.Title, rec.Body, model.JoinTags(rec.Tags), rec.OwnerID, len(rec.Body), now, now)
		if err != nil {
			return err
		}
		id, _ = res.LastInsertId()
		_, err = tx.ExecContext(ctx, "INSERT INTO revisions(record_id,body,editor,created_at) VALUES(?,?,?,?)", id, rec.Body, "system", now)
		return err
	})
	if err != nil {
		return model.Record{}, err
	}
	return s.GetRecord(ctx, id)
}

func (s *Store) UpdateRecord(ctx context.Context, rec model.Record, editor string) (model.Record, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := s.DB.ExecContext(ctx, "UPDATE records SET title=?, body=?, tags=?, bytes=?, updated_at=? WHERE id=?",
		rec.Title, rec.Body, model.JoinTags(rec.Tags), len(rec.Body), now, rec.ID)
	if err != nil {
		return model.Record{}, err
	}
	if _, err := s.DB.ExecContext(ctx, "INSERT INTO revisions(record_id,body,editor,created_at) VALUES(?,?,?,?)", rec.ID, rec.Body, editor, now); err != nil {
		return model.Record{}, err
	}
	return s.GetRecord(ctx, rec.ID)
}

func (s *Store) GetRecord(ctx context.Context, id int64) (model.Record, error) {
	row := s.DB.QueryRowContext(ctx, "SELECT id,slug,title,body,tags,owner_id,bytes,created_at,updated_at FROM records WHERE id=?", id)
	return scanRecord(row)
}

func (s *Store) GetBySlug(ctx context.Context, slug string) (model.Record, error) {
	row := s.DB.QueryRowContext(ctx, "SELECT id,slug,title,body,tags,owner_id,bytes,created_at,updated_at FROM records WHERE slug=?", slug)
	return scanRecord(row)
}

func scanRecord(row *sql.Row) (model.Record, error) {
	var rec model.Record
	var tags, created, updated string
	err := row.Scan(&rec.ID, &rec.Slug, &rec.Title, &rec.Body, &tags, &rec.OwnerID, &rec.Bytes, &created, &updated)
	if err == sql.ErrNoRows {
		return model.Record{}, fmt.Errorf("record: %w", err)
	}
	if err != nil {
		return model.Record{}, err
	}
	rec.Tags = model.SplitTags(tags)
	rec.CreatedAt, _ = time.Parse(time.RFC3339, created)
	rec.UpdatedAt, _ = time.Parse(time.RFC3339, updated)
	return rec, nil
}

func (s *Store) ListRecords(ctx context.Context, f model.ListFilter) ([]model.Record, error) {
	q := "SELECT id,slug,title,body,tags,owner_id,bytes,created_at,updated_at FROM records WHERE 1=1"
	args := []any{}
	if f.Query != "" {
		q += " AND (title LIKE ? OR body LIKE ? OR slug LIKE ?)"
		like := likeContains(f.Query)
		args = append(args, like, like, like)
	}
	if f.Tag != "" {
		q += " AND (','||tags||',') LIKE ?"
		args = append(args, "%,"+f.Tag+",%")
	}
	if f.OwnerID > 0 {
		q += " AND owner_id=?"
		args = append(args, f.OwnerID)
	}
	if f.OrderAsc {
		q += " ORDER BY id ASC"
	} else {
		q += " ORDER BY id DESC"
	}
	q += " LIMIT ? OFFSET ?"
	args = append(args, f.Limit, f.Offset)
	rows, err := s.DB.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Record
	for rows.Next() {
		var rec model.Record
		var tags, created, updated string
		if err := rows.Scan(&rec.ID, &rec.Slug, &rec.Title, &rec.Body, &tags, &rec.OwnerID, &rec.Bytes, &created, &updated); err != nil {
			return nil, err
		}
		rec.Tags = model.SplitTags(tags)
		rec.CreatedAt, _ = time.Parse(time.RFC3339, created)
		rec.UpdatedAt, _ = time.Parse(time.RFC3339, updated)
		out = append(out, rec)
	}
	return out, rows.Err()
}

func (s *Store) DeleteRecord(ctx context.Context, id int64) error {
	if _, err := s.DB.ExecContext(ctx, "DELETE FROM attachments WHERE record_id=?", id); err != nil {
		return err
	}
	if _, err := s.DB.ExecContext(ctx, "DELETE FROM revisions WHERE record_id=?", id); err != nil {
		return err
	}
	_, err := s.DB.ExecContext(ctx, "DELETE FROM records WHERE id=?", id)
	return err
}

func (s *Store) Revisions(ctx context.Context, recordID int64) ([]model.Revision, error) {
	rows, err := s.DB.QueryContext(ctx, "SELECT id,record_id,body,editor,created_at FROM revisions WHERE record_id=? ORDER BY id DESC", recordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Revision
	for rows.Next() {
		var r model.Revision
		var created string
		if err := rows.Scan(&r.ID, &r.RecordID, &r.Body, &r.Editor, &created); err != nil {
			return nil, err
		}
		r.CreatedAt, _ = time.Parse(time.RFC3339, created)
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *Store) AddAttachment(ctx context.Context, a model.Attachment) (model.Attachment, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	res, err := s.DB.ExecContext(ctx, "INSERT INTO attachments(record_id,name,sha,size,path,created_at) VALUES(?,?,?,?,?,?)",
		a.RecordID, a.Name, a.SHA, a.Size, a.Path, now)
	if err != nil {
		return model.Attachment{}, err
	}
	id, _ := res.LastInsertId()
	a.ID = id
	a.CreatedAt, _ = time.Parse(time.RFC3339, now)
	return a, nil
}

func (s *Store) Attachments(ctx context.Context, recordID int64) ([]model.Attachment, error) {
	rows, err := s.DB.QueryContext(ctx, "SELECT id,record_id,name,sha,size,path,created_at FROM attachments WHERE record_id=?", recordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Attachment
	for rows.Next() {
		var a model.Attachment
		var created string
		if err := rows.Scan(&a.ID, &a.RecordID, &a.Name, &a.SHA, &a.Size, &a.Path, &created); err != nil {
			return nil, err
		}
		a.CreatedAt, _ = time.Parse(time.RFC3339, created)
		out = append(out, a)
	}
	return out, rows.Err()
}
