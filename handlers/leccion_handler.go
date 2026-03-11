package handlers

import (
	"cursos-api/middleware"
	"cursos-api/models"
	"cursos-api/services"
	"cursos-api/utils"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type LeccionHandler struct {
	leccionService *services.LeccionService
}

func NewLeccionHandler() *LeccionHandler {
	return &LeccionHandler{
		leccionService: services.NewLeccionService(),
	}
}

// Create crea una nueva lección con imagen opcional (multipart/form-data)
func (h *LeccionHandler) Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cursoID, err := strconv.Atoi(vars["cursoId"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID de curso inválido")
		return
	}

	claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)

	// Parsear multipart form (máximo 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "Error al procesar el formulario")
		return
	}

	// Obtener campos del form
	ordenStr := r.FormValue("orden")
	duracionStr := r.FormValue("duracion_minutos")

	orden, err := strconv.Atoi(ordenStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "El orden debe ser un número válido")
		return
	}

	duracion, err := strconv.Atoi(duracionStr)
	if err != nil {
		duracion = 0
	}

	req := &models.CreateLeccionRequest{
		Titulo:          r.FormValue("titulo"),
		Contenido:       r.FormValue("contenido"),
		Orden:           orden,
		DuracionMinutos: duracion,
	}

	// Obtener imagen si viene en el request (opcional)
	var imagen multipart.File
	var imagenHeader *multipart.FileHeader
	imagen, imagenHeader, err = r.FormFile("imagen")
	if err != nil && err != http.ErrMissingFile {
		respondError(w, http.StatusBadRequest, "Error al leer la imagen")
		return
	}
	if imagen != nil {
		defer imagen.Close()
	}

	leccion, err := h.leccionService.Create(cursoID, req, imagen, imagenHeader, claims.UserID, claims.Rol)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Lección creada exitosamente",
		"leccion": leccion,
	})
}

// GetByCurso obtiene todas las lecciones de un curso
func (h *LeccionHandler) GetByCurso(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cursoID, err := strconv.Atoi(vars["cursoId"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID de curso inválido")
		return
	}

	claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)

	lecciones, err := h.leccionService.GetByCurso(cursoID, claims.UserID, claims.Rol)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, lecciones)
}

// GetByID obtiene una lección por ID
func (h *LeccionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cursoID, err := strconv.Atoi(vars["cursoId"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID de curso inválido")
		return
	}

	leccionID, err := strconv.Atoi(vars["leccionId"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID de lección inválido")
		return
	}

	claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)

	leccion, err := h.leccionService.GetByID(cursoID, leccionID, claims.UserID, claims.Rol)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, leccion)
}

// Update actualiza una lección
func (h *LeccionHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cursoID, err := strconv.Atoi(vars["cursoId"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID de curso inválido")
		return
	}

	leccionID, err := strconv.Atoi(vars["leccionId"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID de lección inválido")
		return
	}

	claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)

	var req models.UpdateLeccionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Datos inválidos")
		return
	}

	leccion, err := h.leccionService.Update(cursoID, leccionID, &req, claims.UserID, claims.Rol)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Lección actualizada exitosamente",
		"leccion": leccion,
	})
}

// Delete elimina una lección
func (h *LeccionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cursoID, err := strconv.Atoi(vars["cursoId"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID de curso inválido")
		return
	}

	leccionID, err := strconv.Atoi(vars["leccionId"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID de lección inválido")
		return
	}

	claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)

	err = h.leccionService.Delete(cursoID, leccionID, claims.UserID, claims.Rol)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Lección eliminada exitosamente",
	})
}