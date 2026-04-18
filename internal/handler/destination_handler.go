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
	"time"
	"tourism-backend/internal/domain"
	"tourism-backend/internal/service"

	"github.com/go-chi/chi/v5"
)

type DestinationHandler struct {
	service *service.DestinationService
}

func NewDestinationHandler(service *service.DestinationService) *DestinationHandler {
	return &DestinationHandler{service: service}
}

func (h *DestinationHandler) Create(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			respondError(w, http.StatusBadRequest, "failed to parse form")
			return
		}

		name := r.FormValue("name")
		description := r.FormValue("description")

		if name == "" {
			respondError(w, http.StatusBadRequest, "name is required")
			return
		}

		var imageURL string

		if file, header, err := r.FormFile("image"); err == nil {
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

			destDir := getDestinationUploadDir()
			if err := os.MkdirAll(destDir, 0755); err != nil {
				respondError(w, http.StatusInternalServerError, "failed to create upload directory")
				return
			}

			filename := fmt.Sprintf("destination_%d%s", time.Now().UnixNano(), ext)
			destPath := filepath.Join(destDir, filename)

			dst, err := os.Create(destPath)
			if err != nil {
				respondError(w, http.StatusInternalServerError, "failed to save file")
				return
			}
			defer dst.Close()

			if _, err := io.Copy(dst, file); err != nil {
				respondError(w, http.StatusInternalServerError, "failed to save file")
				return
			}

			imageURL = fmt.Sprintf("/uploads/destinations/%s", filename)
		}

		destination := &domain.Destination{
			Name:        name,
			Description: description,
			ImageURL:    imageURL,
		}

		result, err := h.service.Create(r.Context(), destination)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusCreated, result)
		return
	}

	// Fallback to JSON for backwards compatibility
	var destination domain.Destination
	if err := json.NewDecoder(r.Body).Decode(&destination); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	result, err := h.service.Create(r.Context(), &destination)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, result)
}

func (h *DestinationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	destinations, err := h.service.GetAll(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, destinations)
}

func (h *DestinationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	destination, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, destination)
}

func (h *DestinationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "destination deleted"})
}

func getDestinationUploadDir() string {
	return "./uploads/destinations"
}
