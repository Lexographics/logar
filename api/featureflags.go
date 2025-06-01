package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sadk.dev/logar/models"
)

func (h *Handler) GetFeatureFlags(w http.ResponseWriter, r *http.Request) {
	flags, err := h.logger.GetFeatureFlags().GetFeatureFlags()
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, err.Error()))
		return
	}

	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, flags))
}

func (h *Handler) UpdateFeatureFlag(w http.ResponseWriter, r *http.Request) {
	flagID := r.FormValue("id")
	if flagID == "" {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing 'id' in request body"))
		return
	}

	flagIDUint, err := strconv.ParseUint(flagID, 10, 64)
	if err != nil {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Invalid 'id' in request body"))
		return
	}

	flag, err := h.logger.GetFeatureFlags().GetFeatureFlag(uint(flagIDUint))
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, err.Error()))
		return
	}

	flag.Name = r.FormValue("name")
	flag.Enabled = r.FormValue("enabled") == "true"
	flag.Condition = r.FormValue("condition")

	err = h.logger.GetFeatureFlags().UpdateFeatureFlag(&flag)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, err.Error()))
		return
	}

	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, flag))
}

func (h *Handler) CreateFeatureFlag(w http.ResponseWriter, r *http.Request) {
	flagName := r.FormValue("name")
	if flagName == "" {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing 'name' in request body"))
		return
	}

	flag := models.FeatureFlag{
		Name:      flagName,
		Enabled:   r.FormValue("enabled") == "true",
		Condition: r.FormValue("condition"),
	}

	err := h.logger.GetFeatureFlags().CreateFeatureFlag(&flag)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, err.Error()))
		return
	}

	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, flag))
}

func (h *Handler) DeleteFeatureFlag(w http.ResponseWriter, r *http.Request) {
	flagID := r.URL.Query().Get("id")
	if flagID == "" {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing 'id' in request body"))
		return
	}

	flagIDUint, err := strconv.ParseUint(flagID, 10, 64)
	if err != nil {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Invalid 'id' in request body"))
		return
	}

	err = h.logger.GetFeatureFlags().DeleteFeatureFlag(uint(flagIDUint))
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, err.Error()))
		return
	}

	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, "Feature flag deleted"))
}
