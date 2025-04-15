package logarweb

import (
	"net/http"

	"github.com/Lexographics/logar"
	"github.com/Lexographics/logar/internal/api"
)

func ServeHTTP(basePath string, l *logar.Logger) http.Handler {
	router := http.NewServeMux()

	handler := api.NewHandler(l, api.HandlerConfig{})
	handler.Router(router)

	mux := http.NewServeMux()
	mux.Handle(basePath+"/", http.StripPrefix(basePath, router))
	return mux
}
