package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Lexographics/logar/internal/domain/models"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ParseLogFilters(r *http.Request) (model string, cursor int, count int, severity int, filters []models.Filter, error error) {
	model = r.PathValue("model")
	cursor, _ = strconv.Atoi(r.URL.Query().Get("cursor"))
	count = 20
	severity, _ = strconv.Atoi(r.URL.Query().Get("severity"))

	filters = []models.Filter{}
	filtersJSON := r.URL.Query().Get("filters")
	if filtersJSON != "" {
		err := json.Unmarshal([]byte(filtersJSON), &filters)
		if err != nil {
			return "", 0, 0, 0, nil, err
		}
	}

	return model, cursor, count, severity, filters, nil
}
