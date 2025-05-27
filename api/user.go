package api

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	session, err := h.getSession(r)
	if err != nil {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing authorization header"))
		return
	}

	user, err := h.logger.GetWebPanel().GetUser(session.UserID)
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

	err = h.logger.GetWebPanel().UpdateUser(user)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, "Failed to update user"))
		return
	}

	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, user))
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.logger.GetWebPanel().GetAllUsers()
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, "Failed to get all users"))
		return
	}

	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, users))
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	session, err := h.getSession(r)
	if err != nil {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing authorization header"))
		return
	}

	selfUser, err := h.logger.GetWebPanel().GetUser(session.UserID)
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

	user, err := h.logger.GetWebPanel().CreateUser(username, display_name, password, isAdmin)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, "Failed to create user"))
		return
	}

	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, user))
}
