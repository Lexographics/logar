package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Lexographics/logar/models"
)

func (h *Handler) getSession(r *http.Request) (*models.Session, error) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return nil, errors.New("missing authorization header")
	}

	token := strings.TrimPrefix(authorization, "Bearer ")
	session, err := h.logger.GetWebPanel().GetSession(token)
	if err != nil {
		return nil, err
	}

	return session, nil
}
