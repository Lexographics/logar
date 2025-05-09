package logarweb

import (
	"net/http"

	"github.com/Lexographics/logar"
	"github.com/Lexographics/logar/internal/api"
)

// ServeHTTP serves the logarweb handler
//
// url: url of the logarweb server. e.g. "http://localhost:3000". should not end with a /
//
// basePath: should either start with / or be empty (""). e.g. "/logger" or ""
//
// returns: http.Handler
//
// example:
//
// logarweb.ServeHTTP("http://localhost:3000", "/logger", logger)
//
// logarweb.ServeHTTP("https://example.com", "/logs", logger)
func ServeHTTP(url, basePath string, l logar.App) http.Handler {
	router := http.NewServeMux()

	handler := api.NewHandler(l.(*logar.AppImpl), api.HandlerConfig{
		BasePath: basePath,
		ApiURL:   url + basePath,
	})
	handler.Router(router)

	mux := http.NewServeMux()
	mux.Handle(basePath+"/", http.StripPrefix(basePath, router))
	return mux
}
