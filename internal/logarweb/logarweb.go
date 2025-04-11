package logarweb

import (
	"net/http"

	"github.com/Lexographics/logar"
	"github.com/Lexographics/logar/internal/api"
)

func ServeHTTP(basePath string, l *logar.Logger) http.Handler {
	router := http.NewServeMux()

	handler := api.NewHandler(l)
	handler.Router(router)

	// router.HandleFunc("/auth", handler.Auth)
	// router.HandleFunc("/", handler.AuthMiddleware(handler.Index))
	// router.HandleFunc("/{model}", handler.AuthMiddleware(handler.Index))
	// router.HandleFunc("/{model}/logs", handler.AuthMiddleware(handler.Logs))
	// router.HandleFunc("/{model}/json", handler.AuthMiddleware(handler.LogsJSON))

	mux := http.NewServeMux()
	mux.Handle(basePath+"/", http.StripPrefix(basePath, router))
	return mux
}
