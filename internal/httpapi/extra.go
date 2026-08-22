package httpapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"example.com/gridpulse/internal/backup"
	"example.com/gridpulse/internal/model"
	"example.com/gridpulse/internal/policy"
	"example.com/gridpulse/internal/report"
	"example.com/gridpulse/internal/validate"
	"example.com/gridpulse/internal/workflow"
)

func (s *Server) handleReport(w http.ResponseWriter, r *http.Request) {
	recs, err := s.cat.List(r.Context(), model.ListFilter{Limit: 200})
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	tag := r.URL.Query().Get("tag")
	recs = report.FilterTag(recs, tag)
	recs = report.SortByTitle(recs)
	switch r.URL.Query().Get("fmt") {
	case "md":
		w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		_, _ = w.Write([]byte(report.Build(recs).Markdown()))
	case "csv":
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		_, _ = w.Write([]byte(report.CSV(recs)))
	case "json":
		b, err := report.JSON(recs)
		if err != nil {
			writeJSON(w, 500, map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write(b)
	default:
		writeJSON(w, 200, report.Build(recs))
	}
}

func (s *Server) handleBackup(w http.ResponseWriter, r *http.Request) {
	recs, err := s.cat.List(r.Context(), model.ListFilter{Limit: 500})
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	b, err := backup.Bundle(recs)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	if _, err := backup.ReadManifest(bytes.NewReader(b), int64(len(b))); err != nil {
		writeJSON(w, 500, map[string]string{"error": "manifest: " + err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=gridpulse.zip")
	_, _ = w.Write(b)
}

func (s *Server) handlePlan(w http.ResponseWriter, r *http.Request) {
	steps := policy.Steps()
	p := workflow.Plan(steps)
	recs, _ := s.cat.List(r.Context(), model.ListFilter{Limit: 50})
	err := p.Run(func(name string) error {
		if strings.Contains(name, "check") || name == "validate" || name == "ingest" {
			for _, rec := range recs {
				if err := validate.Title(rec.Title); err != nil {
					return err
				}
				if err := policy.Enforce(rec.Title, rec.Body, rec.Tags); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(p.Report()))
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	n, err := s.st.CountRecords(r.Context())
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{
		"records": n,
		"page":    strconv.Itoa(s.cfg.PageSize),
		"steps":   policy.Steps(),
	})
}

func decodeJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	return json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(v)
}

func (s *Server) handleBySlug(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/api/by-slug/")
	rec, err := s.cat.GetBySlug(r.Context(), slug)
	if err != nil {
		writeJSON(w, 404, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, 200, rec)
}

func (s *Server) handleImport(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, s.cfg.MaxBodyBytes))
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	n, err := s.cat.ImportBody(r.Context(), body, 1)
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]int{"imported": n})
}

func (s *Server) handleSnapshot(w http.ResponseWriter, r *http.Request) {
	if err := s.cat.WriteSnapshot(r.Context(), s.cfg.DataDir); err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]string{"ok": "1"})
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	last, ok := s.cat.Events.Last()
	if !ok {
		writeJSON(w, 200, map[string]any{"events": []any{}})
		return
	}
	writeJSON(w, 200, map[string]any{"last": last, "errors": s.cat.Events.Filter("error")})
}

func (s *Server) handleAttachment(w http.ResponseWriter, r *http.Request, id int64) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var in struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &in); err != nil || in.Name == "" {
		writeJSON(w, 400, map[string]string{"error": "name required"})
		return
	}
	a, err := s.cat.AddFile(r.Context(), id, in.Name)
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 201, a)
}
