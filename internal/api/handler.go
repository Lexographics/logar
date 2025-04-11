package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Lexographics/logar"
	"github.com/Lexographics/logar/internal/domain/models"
)

type Handler struct {
	logger *logar.Logger
}

func NewHandler(logger *logar.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) Router(mux *http.ServeMux) {
	mux.HandleFunc("/{model}/json", h.GetLogs)
}

func (h *Handler) GetLogs(w http.ResponseWriter, r *http.Request) {
	model := r.PathValue("model")
	cursor, _ := strconv.Atoi(r.URL.Query().Get("cursor"))
	count := 20
	severity, _ := strconv.Atoi(r.URL.Query().Get("severity"))
	filter := r.URL.Query().Get("filter")

	opts := []logar.QueryOptFunc{
		logar.WithCursorPagination(cursor, count),
		logar.WithModel(model),
		logar.WithSeverity(models.Severity(severity)),
		logar.WithFilter(filter),
	}

	filters := strings.Split(r.URL.RawQuery, "&")
	for _, filter := range filters {
		values := strings.Split(filter, "=")
		key := values[0]
		value := ""
		if len(values) > 1 {
			value = strings.Join(values[1:], "=")
		}

		if key == "message" {
			value, _ = url.QueryUnescape(value)
			opts = append(opts, logar.WithFilter(value))
		}
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
