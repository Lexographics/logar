package api

import (
	"encoding/json"
	"net/http"
)

type StatusCode int

// Unknown values
const (
	StatusCode_Unknown StatusCode = iota
)

// Success values
const (
	StatusCode_Success StatusCode = iota + 1
)

// Error values
const (
	StatusCode_Error StatusCode = iota + 1000
	StatusCode_SessionExpired
	StatusCode_InvalidRequest
)

type Status int

const (
	StatusSuccess Status = iota
	StatusError
)

type Response struct {
	StatusCode StatusCode `json:"status_code"`
	Data       any        `json:"data,omitempty"`
}

func NewResponse(statusCode StatusCode, data any) *Response {
	return &Response{
		StatusCode: statusCode,
		Data:       data,
	}
}

var SessionExpiredResponse = NewResponse(StatusCode_SessionExpired, nil)

func WriteSessionExpired(w http.ResponseWriter) {
	w.WriteHeader(401)
	json.NewEncoder(w).Encode(SessionExpiredResponse)
}
