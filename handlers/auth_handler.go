package handlers

import (
    "cursos-api/middleware"
    "cursos-api/models"
    "cursos-api/services"
    "cursos-api/utils"
    "encoding/json"
    "net/http"
)

type AuthHandler struct {
    authService *services.AuthService
}

func NewAuthHandler() *AuthHandler {
    return &AuthHandler{
        authService: services.NewAuthService(),
    }
}

// Register maneja el registro de usuarios
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req models.RegisterRequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "Datos inválidos")
        return
    }

    usuario, err := h.authService.Register(&req)
    if err != nil {
        respondError(w, http.StatusBadRequest, err.Error())
        return
    }

    respondJSON(w, http.StatusCreated, map[string]interface{}{
        "message": "Usuario registrado exitosamente",
        "usuario": usuario,
    })
}

// Login maneja el inicio de sesión
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req models.LoginRequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "Datos inválidos")
        return
    }

    response, err := h.authService.Login(&req)
    if err != nil {
        respondError(w, http.StatusUnauthorized, err.Error())
        return
    }

    respondJSON(w, http.StatusOK, response)
}

// GetProfile obtiene el perfil del usuario autenticado
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
    claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)

    usuario, err := h.authService.GetProfile(claims.UserID)
    if err != nil {
        respondError(w, http.StatusNotFound, "Usuario no encontrado")
        return
    }

    respondJSON(w, http.StatusOK, usuario)
}

// Utilidades para respuestas JSON
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
    respondJSON(w, status, map[string]string{"error": message})
}
