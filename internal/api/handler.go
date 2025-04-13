package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Lexographics/logar"
	"github.com/Lexographics/logar/internal/domain/models"
	"github.com/google/uuid"
)

type HandlerConfig struct {
}

type Session struct {
	ID        string
	Username  string
	ExpiresAt time.Time
}

type Handler struct {
	logger   *logar.Logger
	service  *Service
	sessions sync.Map
	cfg      HandlerConfig
}

func NewHandler(logger *logar.Logger, cfg HandlerConfig) *Handler {
	return &Handler{
		logger:   logger,
		service:  NewService(),
		sessions: sync.Map{},
		cfg:      cfg,
	}
}

func (h *Handler) Router(mux *http.ServeMux) {
	mux.HandleFunc("/auth", h.Auth)
	mux.HandleFunc("/{model}/json", h.GetLogs)
	mux.HandleFunc("/{model}/sse", h.GetLogsSSE)
	mux.Handle("/", http.FileServer(http.Dir("webclient/build")))
}

func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		sessionAny, ok := h.sessions.Load(authorization)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		session := sessionAny.(*Session)

		if session.ExpiresAt.Before(time.Now()) {
			h.sessions.Delete(authorization)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	if !h.logger.Auth(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	session := &Session{
		ID:        uuid.New().String(),
		Username:  "admin",
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}
	h.sessions.Store(session.ID, session)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) GetLogs(w http.ResponseWriter, r *http.Request) {
	model, cursor, count, severity, filters, err := h.service.ParseLogFilters(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	opts := []logar.QueryOptFunc{
		logar.WithCursorPagination(cursor, count),
		logar.WithModel(model),
		logar.WithSeverity(models.Severity(severity)),
	}

	for _, filter := range filters {
		opts = append(opts, logar.WithFilterStruct(filter))
	}

	logs, err := h.logger.GetLogs(opts...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastId := uint(0)
	if len(logs) > 0 && len(logs) == count {
		lastId = logs[len(logs)-1].ID
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"Model":  model,
		"Logs":   logs,
		"LastId": lastId,
	})
}

func (h *Handler) GetLogsSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	model, _, count, severity, filters, err := h.service.ParseLogFilters(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastLog, err := h.logger.GetLogs(
		func(qo *logar.QueryOptions) {
			qo.Limit = 1
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	lastId := uint(0)
	if len(lastLog) > 0 {
		lastId = lastLog[0].ID
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "event: ping\ndata: {}\n\n")
	flusher.Flush()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	ctx := r.Context()
	for {
		select {
		case <-ticker.C:

			opts := []logar.QueryOptFunc{
				logar.WithModel(model),
				logar.WithSeverity(models.Severity(severity)),
			}

			for _, filter := range filters {
				opts = append(opts, logar.WithFilterStruct(filter))
			}

			opts = append(opts, logar.WithIDGreaterThan(uint(lastId)))
			opts = append(opts, func(o *logar.QueryOptions) {
				o.Limit = count
				o.PaginationStrategy = logar.PaginationStatus_None
			})

			logs, err := h.logger.GetLogs(opts...)
			if err != nil {
				h.logger.Error("logar-errors", "Error fetching logs for SSE: "+err.Error(), "sse")
				continue
			}

			for _, log := range logs {
				if log.ID > lastId {
					lastId = log.ID
				}
			}

			if len(logs) > 0 {
				data, err := json.Marshal(map[string]any{
					"Model":  model,
					"Logs":   logs,
					"LastId": lastId,
				})

				if err != nil {
					h.logger.Error("logar-errors", "Error marshaling logs for SSE: "+err.Error(), "sse")
					continue
				}

				fmt.Fprintf(w, "event: logs\ndata: %s\n\n", data)
				flusher.Flush()
			}
		case <-ctx.Done():
			return
		}
	}
}
