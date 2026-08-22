package httpapi

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"example.com/gridpulse/internal/auth"
	"example.com/gridpulse/internal/model"
	"example.com/gridpulse/internal/search"
	"example.com/gridpulse/internal/textfmt"
)

type recordIn struct {
	Title string   `json:"title"`
	Body  string   `json:"body"`
	Tags  []string `json:"tags"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if err := s.st.Ping(r.Context()); err != nil {
		writeJSON(w, 500, map[string]string{"status": "down", "error": err.Error()})
		return
	}
	n, _ := s.st.CountRecords(r.Context())
	writeJSON(w, 200, map[string]any{"status": "ok", "records": n, "app": "GridPulse Relay"})
}

func (s *Server) handleMeta(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{
		"title":  "GridPulse Relay",
		"config": s.cfg.PublicMeta(),
	})
}

func (s *Server) handleRecords(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		f := model.ListFilter{
			Query: r.URL.Query().Get("q"),
			Tag:   r.URL.Query().Get("tag"),
		}
		f.Limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
		recs, err := s.cat.List(r.Context(), f)
		if err != nil {
			writeJSON(w, 500, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 200, recs)
	case http.MethodPost:
		var in recordIn
		if err := json.NewDecoder(io.LimitReader(r.Body, s.cfg.MaxBodyBytes)).Decode(&in); err != nil {
			writeJSON(w, 400, map[string]string{"error": "bad json"})
			return
		}
		rec, err := s.cat.Create(r.Context(), in.Title, in.Body, in.Tags, 1)
		if err != nil {
			writeJSON(w, 400, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 201, rec)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleRecordOne(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/api/records/")
	parts := strings.Split(rest, "/")
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, 400, map[string]string{"error": "bad id"})
		return
	}
	if len(parts) > 1 && parts[1] == "attachments" {
		s.handleAttachment(w, r, id)
		return
	}
	switch r.Method {
	case http.MethodGet:
		rec, revs, atts, err := s.cat.Detail(r.Context(), id)
		if err != nil {
			writeJSON(w, 404, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, 200, map[string]any{"record": rec, "revisions": revs, "attachments": atts, "tag_preview": s.cat.PreviewTags(rec.Tags, 8)})
	case http.MethodPut:
		var in recordIn
		if err := decodeJSON(r, &in); err != nil {
			writeJSON(w, 400, map[string]string{"error": "bad json"})
			return
		}
		editor := "web"
		if tok := auth.ParseBearer(r.Header.Get("Authorization")); tok != "" {
			if u, err := s.st.UserByTokenHash(r.Context(), auth.HashToken(tok)); err == nil {
				editor = u.Name
			}
		}
		rec, err := s.cat.Update(r.Context(), id, in.Title, in.Body, in.Tags, editor)
		if err != nil {
			writeJSON(w, 400, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 200, rec)
	case http.MethodDelete:
		if err := s.st.DeleteRecord(r.Context(), id); err != nil {
			writeJSON(w, 500, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 200, map[string]string{"ok": "1"})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleExport(w http.ResponseWriter, r *http.Request) {
	b, err := s.cat.ExportJSON(r.Context())
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(b)
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	recs, err := s.cat.List(r.Context(), model.ListFilter{Query: q, Limit: 50})
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	idx := search.Build(recs)
	hits := idx.Find(q)
	type hit struct {
		ID      int64  `json:"id"`
		Title   string `json:"title"`
		Snippet string `json:"snippet"`
	}
	out := make([]hit, 0, len(hits))
	for _, rec := range hits {
		out = append(out, hit{ID: rec.ID, Title: rec.Title, Snippet: textfmt.Excerpt(rec.Body, 80)})
	}
	writeJSON(w, 200, out)
}

func (s *Server) handleAudit(w http.ResponseWriter, r *http.Request) {
	items, err := s.st.ListAudit(r.Context(), 40)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, items)
}
