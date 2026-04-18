package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"tourism-backend/internal/domain"
	middleware "tourism-backend/internal/middlewares"
	"tourism-backend/internal/service"

	"github.com/go-chi/chi/v5"
)

func decodeJSON(body io.Reader, target interface{}) error {
	return json.NewDecoder(body).Decode(target)
}

type TourHighlightHandler struct {
	service *service.TourHighlightService
}

func NewTourHighlightHandler(service *service.TourHighlightService) *TourHighlightHandler {
	return &TourHighlightHandler{service: service}
}

func (h *TourHighlightHandler) Create(w http.ResponseWriter, r *http.Request) {
	tourID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid tour id")
		return
	}

	// Check content type for multipart
	contentType := r.Header.Get("Content-Type")
	if contentType != "" && strings.HasPrefix(contentType, "multipart/form-data") {
		h.createWithFile(w, r, tourID)
		return
	}

	// JSON fallback
	var highlight domain.TourHighlight
	if err := decodeJSON(r.Body, &highlight); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	highlight.TourID = tourID

	if err := h.service.Create(r.Context(), &highlight); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, highlight)
}

func (h *TourHighlightHandler) createWithFile(w http.ResponseWriter, r *http.Request, tourID int) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "failed to parse form")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		respondError(w, http.StatusBadRequest, "image file is required")
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

	uploadDir := getUploadDir()
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create upload directory")
		return
	}

	filename := fmt.Sprintf("highlight_%d_%d%s", tourID, claims.UserID, ext)
	highlightPath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(highlightPath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to save file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to save file")
		return
	}

	imageURL := fmt.Sprintf("/uploads/%s", filename)

	highlight := &domain.TourHighlight{
		TourID:   tourID,
		ImageURL: imageURL,
	}

	if err := h.service.Create(r.Context(), highlight); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, highlight)
}

func (h *TourHighlightHandler) GetByTourID(w http.ResponseWriter, r *http.Request) {
	tourID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid tour id")
		return
	}

	highlights, err := h.service.GetByTourID(r.Context(), tourID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, highlights)
}

func (h *TourHighlightHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "highlightId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid highlight id")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "highlight deleted"})
}
