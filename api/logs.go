package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"sadk.dev/logar"
	"sadk.dev/logar/models"
)

func (h *Handler) ListModels(w http.ResponseWriter, r *http.Request) {
	models := h.logger.GetAllModels()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, models))
}

func (h *Handler) GetLogs(w http.ResponseWriter, r *http.Request) {
	model, cursor, count, severity, filters, err := h.service.ParseLogFilters(r)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, err.Error()))
		return
	}

	query := logar.NewQuery().
		WithCursorPagination(cursor, count).
		WithModel(model).
		WithSeverity(models.Severity(severity))

	for _, filter := range filters {
		query.WithFilter(filter)
	}

	logs, err := h.logger.GetLogs(query)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, err.Error()))
		return
	}

	lastId := uint(0)
	if len(logs) > 0 && len(logs) == count {
		lastId = logs[len(logs)-1].ID
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, map[string]any{
		"Model":  model,
		"Logs":   logs,
		"LastId": lastId,
	}))
}

func (h *Handler) GetLogsSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	model, _, count, severity, filters, err := h.service.ParseLogFilters(r)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, err.Error()))
		return
	}

	lastLog, err := h.logger.GetLogs(
		&logar.Query{
			Options: &logar.QueryOptions{
				Limit: 1,
			},
		},
	)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, err.Error()))
		return
	}
	lastId := uint(0)
	if len(lastLog) > 0 {
		lastId = lastLog[0].ID
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, "Streaming not supported"))
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

			query := logar.NewQuery().
				WithModel(model).
				WithSeverity(models.Severity(severity))

			for _, filter := range filters {
				query.WithFilter(filter)
			}

			query.WithIDGreaterThan(uint(lastId))
			query.Options.Limit = count
			query.Options.PaginationStrategy = logar.PaginationStatus_None

			logs, err := h.logger.GetLogs(query)
			if err != nil {
				h.logger.GetLogger().Error("logar-errors", "Error fetching logs for SSE: "+err.Error(), "sse")
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
					h.logger.GetLogger().Error("logar-errors", "Error marshaling logs for SSE: "+err.Error(), "sse")
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
