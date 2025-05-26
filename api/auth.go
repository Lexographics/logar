package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/mileusna/useragent"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if !h.logger.GetWebPanel().Auth(r) {
		WriteSessionExpired(w)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.logger.GetWebPanel().LoginUser(username, password)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, err.Error()))
		return
	}

	ua := useragent.Parse(r.UserAgent())
	device := ua.Name + "/" + ua.OS
	token, err := h.logger.GetWebPanel().CreateSession(user, device)
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

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing authorization header"))
		return
	}

	session, err := h.logger.GetWebPanel().GetSession(strings.TrimPrefix(authorization, "Bearer "))
	if err != nil {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing authorization header"))
		return
	}

	h.logger.GetWebPanel().DeleteSession(session.Token)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetActiveSessions(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing authorization header"))
		return
	}

	token := strings.TrimPrefix(authorization, "Bearer ")
	session, err := h.logger.GetWebPanel().GetSession(token)
	if err != nil {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing authorization header"))
		return
	}

	activeSessions, err := h.logger.GetWebPanel().GetActiveSessions(session.UserID)
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

	h.logger.GetWebPanel().DeleteSession(sessionId)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, nil))
}
