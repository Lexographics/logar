package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Lexographics/logar"
	"github.com/Lexographics/logar/internal/domain/models"
)

type Handler struct {
	logger  *logar.Logger
	service *Service
}

func NewHandler(logger *logar.Logger) *Handler {
	return &Handler{
		logger:  logger,
		service: NewService(),
	}
}

func (h *Handler) Router(mux *http.ServeMux) {
	mux.HandleFunc("/{model}/json", h.GetLogs)
	mux.HandleFunc("/{model}/sse", h.GetLogsSSE)
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

	for _, filter := range filters.MessageFilters {
		opts = append(opts, logar.WithFilter(filter))
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

			for _, filter := range filters.MessageFilters {
				opts = append(opts, logar.WithFilter(filter))
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
