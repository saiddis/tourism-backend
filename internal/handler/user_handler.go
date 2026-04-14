package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	middleware "tourism-backend/internal/middlewares"
	"tourism-backend/internal/service"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	user, err := h.service.Register(req.Name, req.Email, req.Password)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	tokens, err := middleware.GenerateTokenPair(user.ID, user.Email, string(user.Role))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	setRefreshTokenCookie(w, tokens.RefreshToken)

	resp := map[string]interface{}{
		"access_token":       tokens.AccessToken,
		"token_type":         "Bearer",
		"access_expires_in":  tokens.AccessExpiresIn,
		"refresh_expires_in": tokens.RefreshExpiresIn,
	}
	respondJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken := getRefreshTokenFromCookie(r)
	if refreshToken == "" {
		respondError(w, http.StatusUnauthorized, "refresh token required")
		return
	}

	claims, err := middleware.ValidateRefreshToken(refreshToken)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	user, err := h.service.GetByID(claims.UserID)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	tokens, err := middleware.GenerateTokenPair(user.ID, user.Email, string(user.Role))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	setRefreshTokenCookie(w, tokens.RefreshToken)

	resp := map[string]interface{}{
		"access_token":       tokens.AccessToken,
		"token_type":         "Bearer",
		"access_expires_in":  tokens.AccessExpiresIn,
		"refresh_expires_in": tokens.RefreshExpiresIn,
	}
	respondJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	clearRefreshTokenCookie(w)
	respondJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	user, err := h.service.GetByID(claims.UserID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, user)
}

func setRefreshTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		HttpOnly: true,
		Secure:   false, // NOTE: secure when https
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(middleware.RefreshTokenTTL),
	})
}

func getRefreshTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func clearRefreshTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   -1,
	})
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, user)
}
