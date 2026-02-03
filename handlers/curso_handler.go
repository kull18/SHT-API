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

type CursoHandler struct {
    cursoService *services.CursoService
}

func NewCursoHandler() *CursoHandler {
    return &CursoHandler{
        cursoService: services.NewCursoService(),
    }
}

// Create crea un nuevo curso
func (h *CursoHandler) Create(w http.ResponseWriter, r *http.Request) {
    claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)

    var curso models.Curso
    if err := json.NewDecoder(r.Body).Decode(&curso); err != nil {
        respondError(w, http.StatusBadRequest, "Datos inválidos")
        return
    }

    createdCurso, err := h.cursoService.Create(&curso, claims.UserID, claims.Rol)
    if err != nil {
        respondError(w, http.StatusBadRequest, err.Error())
        return
    }

    respondJSON(w, http.StatusCreated, map[string]interface{}{
        "message": "Curso creado exitosamente",
        "curso":   createdCurso,
    })
}

// GetAll obtiene todos los cursos
func (h *CursoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
    claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)

    cursos, err := h.cursoService.GetAll(claims.Rol, claims.UserID)
    if err != nil {
        respondError(w, http.StatusInternalServerError, "Error al obtener cursos")
        return
    }

    respondJSON(w, http.StatusOK, cursos)
}

// GetByID obtiene un curso por ID
func (h *CursoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)

    curso, err := h.cursoService.GetByID(id, claims.Rol, claims.UserID)
    if err != nil {
        respondError(w, http.StatusNotFound, err.Error())
        return
    }

    respondJSON(w, http.StatusOK, curso)
}

// GetMyCursos obtiene los cursos del instructor autenticado
func (h *CursoHandler) GetMyCursos(w http.ResponseWriter, r *http.Request) {
    claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)

    if claims.Rol != "instructor" {
        respondError(w, http.StatusForbidden, "Solo los instructores pueden acceder a esta ruta")
        return
    }

    cursos, err := h.cursoService.GetMyCursos(claims.UserID)
    if err != nil {
        respondError(w, http.StatusInternalServerError, "Error al obtener cursos")
        return
    }

    respondJSON(w, http.StatusOK, cursos)
}

// Update actualiza un curso
func (h *CursoHandler) Update(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)

    var curso models.Curso
    if err := json.NewDecoder(r.Body).Decode(&curso); err != nil {
        respondError(w, http.StatusBadRequest, "Datos inválidos")
        return
    }

    updatedCurso, err := h.cursoService.Update(id, &curso, claims.UserID, claims.Rol)
    if err != nil {
        respondError(w, http.StatusBadRequest, err.Error())
        return
    }

    respondJSON(w, http.StatusOK, map[string]interface{}{
        "message": "Curso actualizado exitosamente",
        "curso":   updatedCurso,
    })
}

// Delete elimina un curso
func (h *CursoHandler) Delete(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)

    err = h.cursoService.Delete(id, claims.UserID, claims.Rol)
    if err != nil {
        respondError(w, http.StatusBadRequest, err.Error())
        return
    }

    respondJSON(w, http.StatusOK, map[string]string{
        "message": "Curso eliminado exitosamente",
    })
}

// ToggleActivo activa o desactiva un curso
func (h *CursoHandler) ToggleActivo(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)

    curso, err := h.cursoService.ToggleActivo(id, claims.UserID, claims.Rol)
    if err != nil {
        respondError(w, http.StatusBadRequest, err.Error())
        return
    }

    status := "desactivado"
    if curso.Activo {
        status = "activado"
    }

    respondJSON(w, http.StatusOK, map[string]interface{}{
        "message": "Curso " + status + " exitosamente",
        "curso":   curso,
    })
}
