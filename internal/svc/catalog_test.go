package svc

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"example.com/gridpulse/internal/config"
	"example.com/gridpulse/internal/db"
	"example.com/gridpulse/internal/model"
	"example.com/gridpulse/internal/policy"
	"example.com/gridpulse/internal/store"
)

func TestCatalogCreateList(t *testing.T) {
	dir := t.TempDir()
	sqlDB, err := db.Open(filepath.Join(dir, "t.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	st := store.New(sqlDB)
	cfg := config.Config{DataDir: filepath.Join(dir, "data"), AdminToken: "dev-token", PageSize: 20}.Normalized()
	_ = os.MkdirAll(cfg.DataDir, 0o755)
	c := NewCatalog(st, cfg)
	ctx := context.Background()
	if err := c.BootstrapAdmin(ctx); err != nil {
		t.Fatal(err)
	}
	s := policy.Sample()
	rec, err := c.Create(ctx, s.Title, s.Body, s.Tags, 1)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Title != s.Title {
		t.Fatalf("title=%q", rec.Title)
	}
	list, err := c.List(ctx, model.ListFilter{Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("len=%d", len(list))
	}
	got, err := c.GetBySlug(ctx, rec.Slug)
	if err != nil || got.ID != rec.ID {
		t.Fatalf("by slug: %+v %v", got, err)
	}
}

func TestCatalogPolicyReject(t *testing.T) {
	dir := t.TempDir()
	sqlDB, err := db.Open(filepath.Join(dir, "t.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	c := NewCatalog(store.New(sqlDB), config.Config{DataDir: dir, AdminToken: "dev-token", PageSize: 20}.Normalized())
	_, err = c.Create(context.Background(), "x", "not-valid-payload", nil, 1)
	if err == nil {
		t.Fatal("expected policy reject")
	}
}
