package handlers

import (
	"encoding/json"
	"net/http"

	"ecochain-victoria/internal/utils"
)

func (h *HandlerSet) LoginPage(w http.ResponseWriter, r *http.Request) {
	h.Render.Page(w, r, "login", nil)
}

func (h *HandlerSet) RegisterPage(w http.ResponseWriter, r *http.Request) {
	h.Render.Page(w, r, "register", nil)
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (h *HandlerSet) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "invalid request body")
		return
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		utils.BadRequest(w, "name, email and password are required")
		return
	}

	user, err := h.Auth.Register(req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		utils.Error(w, http.StatusConflict, err.Error())
		return
	}

	token, err := h.Auth.Token(user)
	if err != nil {
		utils.InternalError(w, "could not issue token")
		return
	}
	setAuthCookie(w, token)
	utils.Created(w, map[string]interface{}{"id": user.ID, "email": user.Email, "role": user.Role})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *HandlerSet) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "invalid request body")
		return
	}

	user, err := h.Auth.Login(req.Email, req.Password)
	if err != nil {
		utils.Unauthorized(w, err.Error())
		return
	}

	token, err := h.Auth.Token(user)
	if err != nil {
		utils.InternalError(w, "could not issue token")
		return
	}
	setAuthCookie(w, token)
	utils.OK(w, map[string]interface{}{"id": user.ID, "email": user.Email, "role": user.Role})
}

func (h *HandlerSet) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	utils.OK(w, map[string]interface{}{"logged_out": true})
}

func setAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60 * 24,
	})
}
