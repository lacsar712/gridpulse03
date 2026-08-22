package svc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"example.com/gridpulse/internal/model"
	"example.com/gridpulse/internal/util/filedump"
	"example.com/gridpulse/internal/util/jsonx"
)

type Snapshot struct {
	Title   string            `json:"title"`
	Count   int               `json:"count"`
	Records []model.Record    `json:"records"`
	Export  model.ExportFile  `json:"export"`
	Meta    map[string]string `json:"meta"`
}

func (c *Catalog) ExportJSON(ctx context.Context) ([]byte, error) {
	recs, err := c.List(ctx, model.ListFilter{Limit: 200})
	if err != nil {
		return nil, err
	}
	exp := model.ExportFile{Name: "GridPulse Relay", Age: len(recs)}
	raw, err := jsonx.EncodeExport(exp)
	if err != nil {
		return nil, err
	}
	decoded, err := jsonx.DecodeExport(raw)
	if err != nil {
		return nil, err
	}
	snap := Snapshot{
		Title:   "GridPulse Relay",
		Count:   len(recs),
		Records: recs,
		Export:  decoded,
		Meta:    map[string]string{"app": "gridpulse"},
	}
	return json.MarshalIndent(snap, "", "  ")
}

func (c *Catalog) WriteSnapshot(ctx context.Context, dir string) error {
	b, err := c.ExportJSON(ctx)
	if err != nil {
		return err
	}
	path := filepath.Join(dir, "snapshot.json")
	return filedump.WriteAll(path, string(b))
}

func (c *Catalog) ImportBody(ctx context.Context, body []byte, owner int64) (int, error) {
	var items []model.Record
	if err := json.Unmarshal(body, &items); err != nil {
		dec := json.NewDecoder(bytes.NewReader(body))
		if err2 := dec.Decode(&items); err2 != nil {
			return 0, fmt.Errorf("import decode: %w", err)
		}
	}
	n := 0
	for _, it := range items {
		title := strings.TrimSpace(it.Title)
		if title == "" {
			continue
		}
		if _, err := c.Create(ctx, title, it.Body, it.Tags, owner); err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}
