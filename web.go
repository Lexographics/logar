package logar

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Lexographics/logar/internal/domain/models"
	"github.com/Lexographics/logar/internal/options/config"
)

type Route struct {
	Name string
	ID   string
}

type LogsData struct {
	Model  string
	Logs   []models.Log
	LastID uint
}

type IndexData struct {
	AppName      string
	Routes       []Route
	CurrentRoute string
	Logs         LogsData
}

func (l *Logger) ServeHTTP() http.Handler {
	basePath := "/logger"
	handler := NewHandler(l)

	router := http.NewServeMux()
	router.HandleFunc("/auth", handler.Auth)
	router.HandleFunc("/", handler.AuthMiddleware(handler.Index))
	router.HandleFunc("/{model}", handler.AuthMiddleware(handler.Index))
	router.HandleFunc("/{model}/logs", handler.AuthMiddleware(handler.Logs))

	mux := http.NewServeMux()
	mux.Handle(basePath+"/", http.StripPrefix(basePath, router))
	return mux
}

type Handler struct {
	logger   *Logger
	template *template.Template
}

func NewHandler(logger *Logger) *Handler {
	return &Handler{logger: logger}
}

func (h *Handler) GetLogs(r *http.Request) (model string, logs []models.Log, lastLogId uint, err error) {
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
		WithSeverity(models.Severity(severity)),
		WithFilter(filter),
	)
	if err != nil {
		return model, nil, 0, err
	}

	if filter != "" {
		replacer := NewCaseInsensitiveReplacer(filter, "[mark]"+filter+"[/mark]")
		for i, log := range logs {
			logs[i].Message = replacer.Replace(log.Message)
		}
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
		Routes:       LogModelToRoutes(&h.logger.config.Models),
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
		"Severity_String": func(s models.Severity) string {
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
		"Escape": func(s string) template.HTML {
			// Converts all special characters to their entity equivalents
			// TODO: Replace all at the same time
			s = strings.Replace(s, "\\u003c", "<", -1)
			s = strings.Replace(s, "\\u003e", ">", -1)
			s = strings.Replace(s, "&", "&amp;", -1)
			s = strings.Replace(s, "<", "&lt;", -1)
			s = strings.Replace(s, ">", "&gt;", -1)
			s = strings.Replace(s, "\"", "&quot;", -1)
			s = strings.Replace(s, "'", "&#39;", -1)
			s = strings.Replace(s, "\n", "<br>", -1)
			s = strings.Replace(s, "\\n", "<br>", -1)
			s = strings.Replace(s, "\t", "&nbsp;&nbsp;&nbsp;&nbsp;", -1)
			s = strings.Replace(s, " ", "&nbsp;", -1)

			s = strings.Replace(s, "[mark]", "<mark>", -1)
			s = strings.Replace(s, "[/mark]", "</mark>", -1)

			return template.HTML(s)
		},
	})

	template, err := template.Parse(index_html)
	// template, err := template.ParseFiles("static/index.html")
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

func LogModelToRoutes(logs *config.LogModels) []Route {
	routes := make([]Route, len(*logs))
	for i, m := range *logs {
		routes[i] = Route{Name: m.DisplayName, ID: m.ModelId}
	}
	return routes
}

type CaseInsensitiveReplacer struct {
	toReplace      *regexp.Regexp
	replaceWith    string
	originalSearch string
}

func NewCaseInsensitiveReplacer(toReplace, replaceWith string) *CaseInsensitiveReplacer {
	return &CaseInsensitiveReplacer{
		toReplace:      regexp.MustCompile("(?i)(" + regexp.QuoteMeta(toReplace) + ")"),
		replaceWith:    replaceWith,
		originalSearch: toReplace,
	}
}

func (cir *CaseInsensitiveReplacer) Replace(str string) string {
	return cir.toReplace.ReplaceAllStringFunc(str, func(found string) string {
		return strings.Replace(cir.replaceWith, cir.originalSearch, found, 1)
	})
}
