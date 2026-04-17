package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"tourism-backend/internal/domain"
	middleware "tourism-backend/internal/middlewares"
	"tourism-backend/internal/service"

	"github.com/go-chi/chi/v5"
)

type ProviderApplicationHandler struct {
	service *service.ProviderApplicationService
}

func NewProviderApplicationHandler(service *service.ProviderApplicationService) *ProviderApplicationHandler {
	return &ProviderApplicationHandler{service: service}
}

func (h *ProviderApplicationHandler) Submit(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req domain.ProviderApplication
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	app, err := h.service.Submit(r.Context(), claims.UserID, &req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, app)
}

func (h *ProviderApplicationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	apps, err := h.service.GetAll(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, apps)
}

func (h *ProviderApplicationHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	app, err := h.service.GetByUserID(r.Context(), claims.UserID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if app == nil {
		respondJSON(w, http.StatusOK, nil)
		return
	}
	respondJSON(w, http.StatusOK, app)
}

func (h *ProviderApplicationHandler) Accept(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req struct {
		AdminNote string `json:"admin_note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	app, err := h.service.Accept(r.Context(), id, req.AdminNote)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, app)
}

func (h *ProviderApplicationHandler) Reject(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req struct {
		AdminNote string `json:"admin_note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	app, err := h.service.Reject(r.Context(), id, req.AdminNote)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, app)
}

func (h *ProviderApplicationHandler) AcceptByToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token == "" {
		http.Redirect(w, r, "/admin/provider-applications?result=error&message=invalid+token", http.StatusFound)
		return
	}

	app, err := h.service.AcceptByToken(r.Context(), token)
	if err != nil {
		http.Redirect(w, r, "/admin/provider-applications?result=error&message="+err.Error(), http.StatusFound)
		return
	}

	http.Redirect(w, r, "/admin/provider-applications?result=accepted", http.StatusFound)
	_ = app
}

func (h *ProviderApplicationHandler) RejectByToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token == "" {
		http.Redirect(w, r, "/admin/provider-applications?result=error&message=invalid+token", http.StatusFound)
		return
	}

	app, err := h.service.RejectByToken(r.Context(), token)
	if err != nil {
		http.Redirect(w, r, "/admin/provider-applications?result=error&message="+err.Error(), http.StatusFound)
		return
	}

	http.Redirect(w, r, "/admin/provider-applications?result=rejected", http.StatusFound)
	_ = app
}
