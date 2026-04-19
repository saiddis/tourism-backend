package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	middleware "tourism-backend/internal/middlewares"
	"tourism-backend/internal/service"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	service      *service.UserService
	cookieDomain string
}

func NewUserHandler(service *service.UserService, cookieDomain string) *UserHandler {
	return &UserHandler{service: service, cookieDomain: cookieDomain}
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
	user, err := h.service.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := middleware.GenerateTokenPair(user.ID, user.Email, string(user.Role), user.Name, user.AvatarURL, user.Balance)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	setRefreshTokenCookie(w, tokens.RefreshToken, h.cookieDomain)

	resp := map[string]interface{}{
		"user":               user,
		"access_token":       tokens.AccessToken,
		"token_type":         "Bearer",
		"access_expires_in":  tokens.AccessExpiresIn,
		"refresh_expires_in": tokens.RefreshExpiresIn,
	}
	respondJSON(w, http.StatusCreated, resp)
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

	user, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	tokens, err := middleware.GenerateTokenPair(user.ID, user.Email, string(user.Role), user.Name, user.AvatarURL, user.Balance)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	setRefreshTokenCookie(w, tokens.RefreshToken, h.cookieDomain)

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

	user, err := h.service.GetByID(r.Context(), claims.UserID)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	tokens, err := middleware.GenerateTokenPair(user.ID, user.Email, string(user.Role), user.Name, user.AvatarURL, user.Balance)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	setRefreshTokenCookie(w, tokens.RefreshToken, h.cookieDomain)

	resp := map[string]interface{}{
		"access_token":       tokens.AccessToken,
		"token_type":         "Bearer",
		"access_expires_in":  tokens.AccessExpiresIn,
		"refresh_expires_in": tokens.RefreshExpiresIn,
	}
	respondJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	clearRefreshTokenCookie(w, h.cookieDomain)
	respondJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	user, err := h.service.GetByID(r.Context(), claims.UserID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, user)
}

func setRefreshTokenCookie(w http.ResponseWriter, token string, domain string) {
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Now().Add(middleware.RefreshTokenTTL),
	}
	if domain != "" {
		cookie.Domain = domain
	}
	http.SetCookie(w, cookie)
}

func getRefreshTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func clearRefreshTokenCookie(w http.ResponseWriter, domain string) {
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   -1,
	}
	if domain != "" {
		cookie.Domain = domain
	}
	http.SetCookie(w, cookie)
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll(r.Context())
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

	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.service.Update(r.Context(), claims.UserID, req.Name, req.Email)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) SetAvatarURL(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.AvatarURL == "" {
		respondError(w, http.StatusBadRequest, "avatar_url is required")
		return
	}

	if err := h.service.UpdateAvatarURL(r.Context(), claims.UserID, req.AvatarURL); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := h.service.GetByID(r.Context(), claims.UserID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "failed to parse form")
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		respondError(w, http.StatusBadRequest, "avatar file is required")
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}
	if !allowedExts[ext] {
		respondError(w, http.StatusBadRequest, "invalid file type. allowed: jpg, jpeg, png, webp, gif")
		return
	}

	if header.Size > 5*1024*1024 {
		respondError(w, http.StatusBadRequest, "file too large. max 5MB")
		return
	}

	avatarDir := getAvatarUploadDir()
	if err := os.MkdirAll(avatarDir, 0755); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create upload directory")
		return
	}

	filename := fmt.Sprintf("avatar_%d%s", claims.UserID, ext)
	avatarPath := filepath.Join(avatarDir, filename)

	dst, err := os.Create(avatarPath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to save file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to save file")
		return
	}

	avatarURL := fmt.Sprintf("/uploads/avatars/%s", filename)
	if err := h.service.UpdateAvatarURL(r.Context(), claims.UserID, avatarURL); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := h.service.GetByID(r.Context(), claims.UserID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Amount <= 0 {
		respondError(w, http.StatusBadRequest, "amount must be positive")
		return
	}

	user, err := h.service.Deposit(r.Context(), claims.UserID, req.Amount)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := middleware.GenerateTokenPair(user.ID, user.Email, string(user.Role), user.Name, user.AvatarURL, user.Balance)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	resp := map[string]interface{}{
		"user":         user,
		"access_token": tokens.AccessToken,
	}
	respondJSON(w, http.StatusOK, resp)
}

var uploadDir string

func SetUploadDir(dir string) {
	uploadDir = dir
}

func getAvatarUploadDir() string {
	if uploadDir != "" {
		return uploadDir + "/avatars"
	}
	return "./uploads/avatars"
}
