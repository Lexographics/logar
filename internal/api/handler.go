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
	"github.com/mileusna/useragent"
)

type HandlerConfig struct {
	BasePath string
	ApiURL   string
}

type InvokeActionRequest struct {
	Path string   `json:"path"`
	Args []string `json:"args"`
}

type InvokeActionResponse struct {
	Result any    `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

type Handler struct {
	logger  *logar.AppImpl
	service *Service
	cfg     HandlerConfig
}

func NewHandler(logger *logar.AppImpl, cfg HandlerConfig) *Handler {
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
	mux.HandleFunc("GET /language", h.GetLanguage)
	mux.HandleFunc("POST /auth/login", h.Login)
	mux.HandleFunc("POST /auth/logout", h.AuthMiddleware(h.Logout))
	mux.HandleFunc("GET /auth/sessions", h.AuthMiddleware(h.GetActiveSessions))
	mux.HandleFunc("POST /auth/revoke-session", h.AuthMiddleware(h.RevokeSession))
	mux.HandleFunc("GET /models", h.AuthMiddleware(h.ListModels))
	mux.HandleFunc("GET /logs/{model}", h.AuthMiddleware(h.GetLogs))
	mux.HandleFunc("GET /logs/{model}/sse", h.AuthMiddleware(h.GetLogsSSE))
	mux.HandleFunc("GET /actions", h.AuthMiddleware(h.ListActions))
	mux.HandleFunc("POST /actions/invoke", h.AuthMiddleware(h.InvokeActionHandler))
	mux.HandleFunc("PUT /user", h.AuthMiddleware(h.UpdateUser))
	mux.HandleFunc("POST /user", h.AuthMiddleware(h.CreateUser))
	mux.HandleFunc("GET /user", h.AuthMiddleware(h.GetAllUsers))
	mux.HandleFunc("GET /analytics", h.AuthMiddleware(h.GetAnalytics))

	if dev {
		mux.Handle("/", h.SetMetadataMiddleware(http.FileServer(http.Dir("webclient/build"))))
	} else {
		sub, err := fs.Sub(staticFiles, "build")
		if err != nil {
			h.logger.Error("logar-errors", fmt.Sprintf("Failed to create subdirectory: %v", err), "api")
			return
		}
		mux.Handle("/", h.SetMetadataMiddleware(http.FileServer(http.FS(sub))))
	}
}

func (h *Handler) SetMetadataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		basePath, err := r.Cookie("base-path")
		if err != nil || basePath.Value != h.cfg.BasePath {
			fmt.Println("Setting base-path cookie to", h.cfg.BasePath)
			http.SetCookie(w, &http.Cookie{
				Name:   "base-path",
				Value:  h.cfg.BasePath,
				Path:   "/",
				MaxAge: 86400,
			})
		}

		apiUrl, err := r.Cookie("api-url")
		if err != nil || apiUrl.Value != h.cfg.ApiURL {
			fmt.Println("Setting api-url cookie to", h.cfg.ApiURL)
			http.SetCookie(w, &http.Cookie{
				Name:   "api-url",
				Value:  h.cfg.ApiURL,
				Path:   "/",
				MaxAge: 86400,
			})
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			authorization = r.URL.Query().Get("token")
		}

		if authorization == "" {
			WriteSessionExpired(w)
			return
		}

		token := strings.TrimPrefix(authorization, "Bearer ")

		_, err := h.logger.GetSession(token)
		if err != nil {
			WriteSessionExpired(w)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (h *Handler) GetLanguage(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, h.logger.GetDefaultLanguage()))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if !h.logger.Auth(r) {
		WriteSessionExpired(w)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.logger.LoginUser(username, password)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, err.Error()))
		return
	}

	ua := useragent.Parse(r.UserAgent())
	device := ua.Name + "/" + ua.OS
	token, err := h.logger.CreateSession(user, device)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, map[string]any{
		"token": token,
		"user":  user,
	}))
}

func (h *Handler) RevokeSession(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing authorization header"))
		return
	}

	sessionId := r.FormValue("session_id")
	if sessionId == "" {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing 'session_id' in request body"))
		return
	}

	h.logger.DeleteSession(sessionId)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, nil))
}

func (h *Handler) GetActiveSessions(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing authorization header"))
		return
	}

	token := strings.TrimPrefix(authorization, "Bearer ")
	session, err := h.logger.GetSession(token)
	if err != nil {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing authorization header"))
		return
	}

	activeSessions, err := h.logger.GetActiveSessions(session.UserID)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, "Failed to get active sessions"))
		return
	}

	sessionData := []SessionData{}
	for _, session := range activeSessions {
		sessionData = append(sessionData, SessionData{
			Device:       session.Device,
			LastActivity: session.LastActivity.Format(time.DateTime),
			CreatedAt:    session.CreatedAt.Format(time.DateTime),
			IsCurrent:    session.Token == token,
			Token:        session.Token,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, sessionData))
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing authorization header"))
		return
	}

	token := strings.TrimPrefix(authorization, "Bearer ")
	session, err := h.logger.GetSession(token)
	if err != nil {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing authorization header"))
		return
	}

	user, err := h.logger.GetUser(session.UserID)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, "Failed to get user"))
		return
	}

	if user.ID == 0 {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Main user cannot be updated"))
		return
	}

	displayName := r.FormValue("display_name")
	if displayName != "" {
		user.DisplayName = displayName
	}

	err = h.logger.UpdateUser(user)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, "Failed to update user"))
		return
	}

	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, user))
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing authorization header"))
		return
	}

	session, err := h.logger.GetSession(strings.TrimPrefix(authorization, "Bearer "))
	if err != nil {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing authorization header"))
		return
	}

	h.logger.DeleteSession(session.Token)
	w.WriteHeader(http.StatusOK)
}

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

func (h *Handler) ListActions(w http.ResponseWriter, r *http.Request) {
	actionsMap := h.logger.GetActionsMap()
	details := []ActionDetails{}

	for _, action := range actionsMap {
		argTypes, err := h.logger.GetActionArgTypes(action.Path)
		if err != nil {
			h.logger.Error("logar-errors", fmt.Sprintf("Error getting arg types for action %s: %v", action.Path, err), "api")
			continue
		}

		argTypeData := make([]ArgType, len(argTypes))
		for i, t := range argTypes {
			kind, ok := h.logger.GetTypeKind(t)
			if !ok {
				kind = logar.TypeKind_Text
			}

			argTypeData[i] = ArgType{
				Type: t.String(),
				Kind: string(kind),
			}
		}

		details = append(details, ActionDetails{
			Path:        action.Path,
			Args:        argTypeData,
			Description: action.Description,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, details))
}

func (h *Handler) InvokeActionHandler(w http.ResponseWriter, r *http.Request) {
	var req InvokeActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, fmt.Sprintf("Invalid request body: %v", err)))
		return
	}

	if req.Path == "" {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing 'path' in request body"))
		return
	}

	expectedTypes, err := h.logger.GetActionArgTypes(req.Path)
	if err != nil {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, fmt.Sprintf("Error finding action '%s': %v", req.Path, err)))
		return
	}

	if len(req.Args) != len(expectedTypes) {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, fmt.Sprintf("Action '%s' expects %d arguments, but received %d", req.Path, len(expectedTypes), len(req.Args))))
		return
	}

	parsedArgs := make([]any, len(req.Args))
	for i, argStr := range req.Args {
		expectedType := expectedTypes[i]
		val, err := parseStringArg(argStr, expectedType)
		if err != nil {
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, fmt.Sprintf("Error parsing argument %d for action '%s': expected type %s, error: %v", i+1, req.Path, expectedType.String(), err)))
			return
		}
		parsedArgs[i] = val
	}

	result, err := h.logger.InvokeAction(req.Path, parsedArgs...)

	w.Header().Set("Content-Type", "application/json")
	resp := InvokeActionResponse{}
	if err != nil {
		resp.Error = err.Error()
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, err.Error()))
	} else {
		if len(result) == 1 {
			resp.Result = result[0]
		} else {
			resp.Result = result
		}
	}

	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, resp))
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.logger.GetAllUsers()
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, "Failed to get all users"))
		return
	}

	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, users))
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing authorization header"))
		return
	}

	token := strings.TrimPrefix(authorization, "Bearer ")
	session, err := h.logger.GetSession(token)
	if err != nil {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing authorization header"))
		return
	}

	selfUser, err := h.logger.GetUser(session.UserID)
	if err != nil {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing authorization header"))
		return
	}
	if !selfUser.IsAdmin {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Only admins can create users"))
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	isAdmin := r.FormValue("is_admin") == "true"
	display_name := r.FormValue("display_name")
	if username == "" || password == "" {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing 'username' or 'password' in request body"))
		return
	}

	user, err := h.logger.CreateUser(username, display_name, password, isAdmin)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, "Failed to create user"))
		return
	}

	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, user))
}

func parseStringArg(argStr string, targetType reflect.Type) (any, error) {
	switch targetType.Kind() {
	case reflect.String:
		p := reflect.New(targetType)
		p.Elem().SetString(argStr)
		return p.Elem().Interface(), nil
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
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, err.Error()))
		return
	}

	lastLog, err := h.logger.GetLogs(
		func(qo *logar.QueryOptions) {
			qo.Limit = 1
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

func (h *Handler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	analytics, err := h.logger.GetStatistics(time.Now().Add(-time.Hour*24*7), time.Now())
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, err.Error()))
		return
	}

	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, analytics))
}
