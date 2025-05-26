package api

import (
	"encoding/json"
	"net/http"
	"time"
)

func (h *Handler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	analytics, err := h.logger.GetAnalytics().GetStatistics(time.Now().Add(-time.Hour*24*30), time.Now())
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, err.Error()))
		return
	}

	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, analytics))
}
