package api

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Lexographics/logar"
	"github.com/Lexographics/logar/internal/domain/models"
)

type HandlerConfig struct {
}

type InvokeActionRequest struct {
	Path string   `json:"path"`
	Args []string `json:"args"`
}

type InvokeActionResponse struct {
	Result any    `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

type ActionDetails struct {
	Path        string   `json:"path"`
	Args        []string `json:"args"`
	Description string   `json:"description"`
}

type Handler struct {
	logger  *logar.Logger
	service *Service
	cfg     HandlerConfig
}

func NewHandler(logger *logar.Logger, cfg HandlerConfig) *Handler {
	return &Handler{
		logger:  logger,
		service: NewService(),
		cfg:     cfg,
	}
}

var dev = false

func init() {
	dev = os.Getenv("LOGAR_DEV") == "true"
}

//go:embed build/*
var staticFiles embed.FS

func (h *Handler) Router(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/login", h.Login)
	mux.HandleFunc("GET /models", h.AuthMiddleware(h.ListModels))
	mux.HandleFunc("GET /logs/{model}", h.AuthMiddleware(h.GetLogs))
	mux.HandleFunc("GET /logs/{model}/sse", h.AuthMiddleware(h.GetLogsSSE))
	mux.HandleFunc("GET /actions", h.AuthMiddleware(h.ListActions))
	mux.HandleFunc("POST /actions/invoke", h.AuthMiddleware(h.InvokeActionHandler))

	if dev {
		mux.Handle("/", http.FileServer(http.Dir("webclient/build")))
	} else {
		sub, err := fs.Sub(staticFiles, "build")
		if err != nil {
			h.logger.Error("logar-errors", fmt.Sprintf("Failed to create subdirectory: %v", err), "api")
			return
		}
		mux.Handle("/", http.FileServer(http.FS(sub)))
	}
}

func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			authorization = r.URL.Query().Get("token")
		}

		if authorization == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authorization, "Bearer ")

		_, err := h.logger.GetSession(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if !h.logger.Auth(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if !h.logger.IsAuthCredentialsCorrect(username, password) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := h.logger.CreateSession(username)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"token":    token,
		"username": username,
	})
}

func (h *Handler) ListModels(w http.ResponseWriter, r *http.Request) {
	models := h.logger.GetAllModels()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models)
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

func (h *Handler) ListActions(w http.ResponseWriter, r *http.Request) {
	actionsMap := h.logger.GetActionsMap()
	details := []ActionDetails{}

	for _, action := range actionsMap {
		argTypes, err := h.logger.GetActionArgTypes(action.Path)
		if err != nil {
			h.logger.Error("logar-errors", fmt.Sprintf("Error getting arg types for action %s: %v", action.Path, err), "api")
			continue
		}

		argTypeStrings := make([]string, len(argTypes))
		for i, t := range argTypes {
			argTypeStrings[i] = t.String()
		}

		details = append(details, ActionDetails{
			Path:        action.Path,
			Args:        argTypeStrings,
			Description: action.Description,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(details); err != nil {
		h.logger.Error("logar-errors", fmt.Sprintf("Failed to encode action details: %v", err), "api")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) InvokeActionHandler(w http.ResponseWriter, r *http.Request) {
	var req InvokeActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if req.Path == "" {
		http.Error(w, "Missing 'path' in request body", http.StatusBadRequest)
		return
	}

	expectedTypes, err := h.logger.GetActionArgTypes(req.Path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error finding action '%s': %v", req.Path, err), http.StatusNotFound)
		return
	}

	if len(req.Args) != len(expectedTypes) {
		http.Error(w, fmt.Sprintf("Action '%s' expects %d arguments, but received %d", req.Path, len(expectedTypes), len(req.Args)), http.StatusBadRequest)
		return
	}

	parsedArgs := make([]any, len(req.Args))
	for i, argStr := range req.Args {
		expectedType := expectedTypes[i]
		val, err := parseStringArg(argStr, expectedType)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing argument %d for action '%s': expected type %s, error: %v", i+1, req.Path, expectedType.String(), err), http.StatusBadRequest)
			return
		}
		parsedArgs[i] = val
	}

	result, err := h.logger.InvokeAction(req.Path, parsedArgs...)

	w.Header().Set("Content-Type", "application/json")
	resp := InvokeActionResponse{}
	if err != nil {
		resp.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		if len(result) == 1 {
			resp.Result = result[0]
		} else {
			resp.Result = result
		}
	}

	if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
		h.logger.Error("logar-errors", fmt.Sprintf("Failed to encode action response: %v", encodeErr), "api")
	}
}

func parseStringArg(argStr string, targetType reflect.Type) (any, error) {
	switch targetType.Kind() {
	case reflect.String:
		return argStr, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(argStr, 10, targetType.Bits())
		if err != nil {
			return nil, fmt.Errorf("cannot parse '%s' as %s: %w", argStr, targetType.Kind(), err)
		}
		p := reflect.New(targetType)
		p.Elem().SetInt(intVal)
		return p.Elem().Interface(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(argStr, 10, targetType.Bits())
		if err != nil {
			return nil, fmt.Errorf("cannot parse '%s' as %s: %w", argStr, targetType.Kind(), err)
		}
		p := reflect.New(targetType)
		p.Elem().SetUint(uintVal)
		return p.Elem().Interface(), nil
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(argStr, targetType.Bits())
		if err != nil {
			return nil, fmt.Errorf("cannot parse '%s' as %s: %w", argStr, targetType.Kind(), err)
		}
		p := reflect.New(targetType)
		p.Elem().SetFloat(floatVal)
		return p.Elem().Interface(), nil
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(argStr)
		if err != nil {
			return nil, fmt.Errorf("cannot parse '%s' as bool: %w", argStr, err)
		}
		return boolVal, nil
	default:
		p := reflect.New(targetType)
		err := json.Unmarshal([]byte(argStr), p.Interface())
		if err != nil {
			return nil, fmt.Errorf("cannot unmarshal '%s' as JSON into %s: %w", argStr, targetType.String(), err)
		}
		return p.Elem().Interface(), nil
	}
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
