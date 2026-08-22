package httpapi

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"example.com/gridpulse/internal/config"
	"example.com/gridpulse/internal/db"
	"example.com/gridpulse/internal/store"
	"example.com/gridpulse/internal/svc"
	"example.com/gridpulse/internal/webui"
)

type Server struct {
	cfg config.Config
	db  *sql.DB
	st  *store.Store
	cat *svc.Catalog
	mux *http.ServeMux
}

func (s *Server) Close() {
	db.CloseQuietly(s.db)
}

func (s *Server) HTTP() *http.Server {
	h := jsonOnly(logging(s.mux))
	return &http.Server{
		Addr:              s.cfg.Addr,
		Handler:           withTimeout(h, s.cfg.RequestTimeout),
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func NewApp(cfg config.Config) (*Server, error) {
	sqlDB, err := db.Open(cfg.DBPath)
	if err != nil {
		return nil, err
	}
	st := store.New(sqlDB)
	cat := svc.NewCatalog(st, cfg)
	if err := cat.BootstrapAdmin(context.Background()); err != nil {
		db.CloseQuietly(sqlDB)
		return nil, err
	}
	s := &Server{cfg: cfg, db: sqlDB, st: st, cat: cat, mux: http.NewServeMux()}
	s.routes()
	return s, nil
}

func New(cfg config.Config) (*http.Server, error) {
	s, err := NewApp(cfg)
	if err != nil {
		return nil, err
	}
	hs := s.HTTP()
	hs.RegisterOnShutdown(s.Close)
	return hs, nil
}

func (s *Server) routes() {
	s.mux.HandleFunc("/", s.handleIndex)
	s.mux.HandleFunc("/healthz", s.handleHealth)
	s.mux.HandleFunc("/api/meta", s.handleMeta)
	s.mux.HandleFunc("/api/records", s.handleRecords)
	s.mux.HandleFunc("/api/records/", s.handleRecordOne)
	s.mux.HandleFunc("/api/export", s.handleExport)
	s.mux.HandleFunc("/api/search", s.handleSearch)
	s.mux.HandleFunc("/api/audit", s.handleAudit)
	s.mux.HandleFunc("/api/report", s.handleReport)
	s.mux.HandleFunc("/api/backup", s.handleBackup)
	s.mux.HandleFunc("/api/plan", s.handlePlan)
	s.mux.HandleFunc("/api/stats", s.handleStats)
	s.mux.HandleFunc("/api/by-slug/", s.handleBySlug)
	s.mux.HandleFunc("/api/import", s.handleImport)
	s.mux.HandleFunc("/api/snapshot", s.handleSnapshot)
	s.mux.HandleFunc("/api/events", s.handleEvents)
	s.mux.HandleFunc("/static/app.js", s.handleJS)
	s.mux.HandleFunc("/static/app.css", s.handleCSS)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(webui.IndexHTML())
}

func (s *Server) handleJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	_, _ = w.Write(webui.AppJS())
}

func (s *Server) handleCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	_, _ = w.Write(webui.AppCSS())
}
