package logar

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Route struct {
	Name string
	ID   string
}

type LogsData struct {
	Model  string
	Logs   []Log
	LastID uint
}

type IndexData struct {
	AppName      string
	Routes       []Route
	CurrentRoute string
	Logs         LogsData
}

func (l *Logger) ServeHTTP() http.Handler {
	handler := NewHandler(l)

	mux := http.NewServeMux()
	mux.HandleFunc("/logger/auth", handler.Auth)
	mux.HandleFunc("/logger/", handler.AuthMiddleware(handler.Index))
	mux.HandleFunc("/logger/{model}", handler.AuthMiddleware(handler.Index))
	mux.HandleFunc("/logger/{model}/logs", handler.AuthMiddleware(handler.Logs))
	return mux
}

type Handler struct {
	logger   *Logger
	template *template.Template
}

func NewHandler(logger *Logger) *Handler {
	return &Handler{logger: logger}
}

func (h *Handler) GetLogs(r *http.Request) (model string, logs []Log, lastLogId uint, err error) {
	model = r.PathValue("model")
	if model == "all" {
		model = ""
	}

	cursor, _ := strconv.Atoi(r.URL.Query().Get("cursor"))
	count := 50
	severity, _ := strconv.Atoi(r.URL.Query().Get("severity"))
	filter := r.URL.Query().Get("filter")

	logs, err = h.logger.GetLogs(
		WithCursorPagination(cursor, count),
		WithModel(model),
		WithSeverity(Severity(severity)),
		WithFilter(filter),
	)
	if err != nil {
		return model, nil, 0, err
	}

	lastId := uint(0)
	if len(logs) > 0 && len(logs) == count {
		lastId = logs[len(logs)-1].ID
	}
	return model, logs, lastId, nil
}

func (h *Handler) Logs(w http.ResponseWriter, r *http.Request) {
	h.LoadTemplate()

	model, logs, lastId, err := h.GetLogs(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.template.ExecuteTemplate(w, "logs", LogsData{
		Model:  model,
		Logs:   logs,
		LastID: lastId,
	})
}

func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	if h.logger.config.AuthFunc != nil && h.logger.config.AuthFunc(r) {
		http.SetCookie(w, &http.Cookie{
			Name:     "logger-auth",
			Value:    authToken,
			Expires:  time.Now().Add(time.Hour * 24 * 7),
			MaxAge:   86400 * 7,
			Secure:   true,
			HttpOnly: true,
		})
		http.Redirect(w, r, "/logger", http.StatusFound)
		return
	}
	http.NotFound(w, r)
}

func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	if !h.logger.config.RequireAuth {
		return next.ServeHTTP
	}

	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("logger-auth")
		if err != nil || cookie.Value != authToken {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	h.LoadTemplate()

	model, logs, lastId, err := h.GetLogs(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.template.ExecuteTemplate(w, "index", IndexData{
		AppName:      h.logger.config.AppName,
		Routes:       h.logger.config.Models.ToRoutes(),
		CurrentRoute: model,
		Logs: LogsData{
			Model:  model,
			Logs:   logs,
			LastID: lastId,
		},
	})
}

func (h *Handler) LoadTemplate() {
	template := template.New("logger").Funcs(template.FuncMap{
		"Severity_String": func(s Severity) string {
			return s.String()
		},
		"ToUpper": func(s string) string {
			return strings.ToUpper(s)
		},
		"ToLower": func(s string) string {
			return strings.ToLower(s)
		},
		"Minus": func(a, b int) int {
			return a - b
		},
	})

	template, err := template.Parse(index_html)
	if err != nil {
		fmt.Printf("load template err: %v\n", err)
		return
	}

	h.template = template
}

func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// TODO: Generate random token for each session
var authToken string = generateRandomString(66)

func (l LogModels) ToRoutes() []Route {
	routes := make([]Route, len(l))
	for i, m := range l {
		routes[i] = Route{Name: m.DisplayName, ID: m.ModelId}
	}
	return routes
}
