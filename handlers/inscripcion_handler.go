package handlers

import (
	"cursos-api/middleware"
	"cursos-api/models"
	"cursos-api/services"
	"cursos-api/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type InscripcionHandler struct {
	inscripcionService *services.InscripcionService
}

func NewInscripcionHandler() *InscripcionHandler {
	return &InscripcionHandler{
		inscripcionService: services.NewInscripcionService(),
	}
}

// Inscribir inscribe al alumno autenticado en un curso
func (h *InscripcionHandler) Inscribir(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)

	var req models.InscripcionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Datos inválidos")
		return
	}

	if req.CursoID == 0 {
		respondError(w, http.StatusBadRequest, "El ID del curso es requerido")
		return
	}

	inscripcion, err := h.inscripcionService.Inscribir(req.CursoID, claims.UserID, claims.Rol)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message":     "Inscripción exitosa",
		"inscripcion": inscripcion,
	})
}

// GetMisInscripciones obtiene todas las inscripciones del alumno autenticado
func (h *InscripcionHandler) GetMisInscripciones(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)

	inscripciones, err := h.inscripcionService.GetMisInscripciones(claims.UserID, claims.Rol)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, inscripciones)
}

// GetByID obtiene una inscripción por ID
func (h *InscripcionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	inscripcionID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID de inscripción inválido")
		return
	}

	claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)

	inscripcion, err := h.inscripcionService.GetByID(inscripcionID, claims.UserID, claims.Rol)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, inscripcion)
}

// MarcarLeccion marca o desmarca una lección como completada
func (h *InscripcionHandler) MarcarLeccion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	inscripcionID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID de inscripción inválido")
		return
	}

	claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)

	var req models.ProgresoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Datos inválidos")
		return
	}

	if req.LeccionID == 0 {
		respondError(w, http.StatusBadRequest, "El ID de la lección es requerido")
		return
	}

	progreso, err := h.inscripcionService.MarcarLeccion(inscripcionID, req.LeccionID, req.Completada, claims.UserID, claims.Rol)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	mensaje := "Lección marcada como pendiente"
	if progreso.Completada {
		mensaje = "Lección marcada como completada"
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":  mensaje,
		"progreso": progreso,
	})
}

// GetProgreso obtiene el progreso de todas las lecciones de una inscripción
func (h *InscripcionHandler) GetProgreso(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	inscripcionID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID de inscripción inválido")
		return
	}

	claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)

	progreso, err := h.inscripcionService.GetProgreso(inscripcionID, claims.UserID, claims.Rol)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, progreso)
}