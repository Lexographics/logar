package api

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"reflect"
	"strconv"

	"github.com/Lexographics/logar"
)

type HandlerConfig struct {
	BasePath       string
	ApiURL         string
	WebClientFiles fs.FS
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

	mux.HandleFunc("GET /feature-flags", h.AuthMiddleware(h.GetFeatureFlags))
	mux.HandleFunc("PUT /feature-flags", h.AuthMiddleware(h.UpdateFeatureFlag))
	mux.HandleFunc("POST /feature-flags", h.AuthMiddleware(h.CreateFeatureFlag))
	mux.HandleFunc("DELETE /feature-flags", h.AuthMiddleware(h.DeleteFeatureFlag))

	mux.HandleFunc("GET /globals", h.AuthMiddleware(h.GetGlobals))
	mux.HandleFunc("PUT /globals", h.AuthMiddleware(h.UpdateGlobal))
	mux.HandleFunc("DELETE /globals", h.AuthMiddleware(h.DeleteGlobal))

	if h.cfg.WebClientFiles != nil && !dev {
		sub, err := fs.Sub(h.cfg.WebClientFiles, "build")
		if err != nil {
			h.logger.GetLogger().Error("logar-errors", fmt.Sprintf("Failed to create subdirectory: %v", err), "api")
			return
		}
		mux.Handle("/", h.SetMetadataMiddleware(http.FileServer(http.FS(sub))))
	} else {
		mux.Handle("/", h.SetMetadataMiddleware(http.FileServer(http.Dir("webclient/build"))))
	}
}

func (h *Handler) GetLanguage(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, h.logger.GetWebPanel().GetDefaultLanguage()))
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
