package api

import (
	"net/http"
	"strings"
)

func (h *Handler) SetMetadataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		basePath, err := r.Cookie("base-path")
		if err != nil || basePath.Value != h.cfg.BasePath {
			http.SetCookie(w, &http.Cookie{
				Name:   "base-path",
				Value:  h.cfg.BasePath,
				Path:   "/",
				MaxAge: 86400,
			})
		}

		apiUrl, err := r.Cookie("api-url")
		if err != nil || apiUrl.Value != h.cfg.ApiURL {
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

		_, err := h.logger.GetWebPanel().GetSession(token)
		if err != nil {
			WriteSessionExpired(w)
			return
		}

		next.ServeHTTP(w, r)
	}
}
