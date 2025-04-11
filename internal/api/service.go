package api

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

type Filters struct {
	MessageFilters []string
}

func (s *Service) ParseLogFilters(r *http.Request) (model string, cursor int, count int, severity int, filters Filters, error error) {
	model = r.PathValue("model")
	cursor, _ = strconv.Atoi(r.URL.Query().Get("cursor"))
	count = 20
	severity, _ = strconv.Atoi(r.URL.Query().Get("severity"))

	queryString := strings.Split(r.URL.RawQuery, "&")
	for _, filter := range queryString {
		values := strings.Split(filter, "=")
		key := values[0]
		value := ""
		if len(values) > 1 {
			value = strings.Join(values[1:], "=")
		}

		if key == "message" {
			value, _ = url.QueryUnescape(value)
			filters.MessageFilters = append(filters.MessageFilters, value)
		}
	}

	return model, cursor, count, severity, filters, nil
}
