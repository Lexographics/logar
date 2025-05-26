package api

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) GetGlobals(w http.ResponseWriter, r *http.Request) {
	globals, err := h.logger.GetAllGlobals()
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, err.Error()))
		return
	}

	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, globals))
}

func (h *Handler) UpdateGlobal(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing 'key' in request body"))
		return
	}

	var data any
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Invalid request body"))
		return
	}

	err := h.logger.SetGlobal(key, data)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, err.Error()))
		return
	}

	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, nil))
}

func (h *Handler) DeleteGlobal(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing 'key' in request body"))
		return
	}

	err := h.logger.DeleteGlobal(key)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, err.Error()))
		return
	}

	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, nil))
}
