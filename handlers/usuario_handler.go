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

type UsuarioHandler struct {
    usuarioService *services.UsuarioService
}

func NewUsuarioHandler() *UsuarioHandler {
    return &UsuarioHandler{
        usuarioService: services.NewUsuarioService(),
    }
}

// GetAll obtiene todos los usuarios
func (h *UsuarioHandler) GetAll(w http.ResponseWriter, r *http.Request) {
    usuarios, err := h.usuarioService.GetAll()
    if err != nil {
        respondError(w, http.StatusInternalServerError, "Error al obtener usuarios")
        return
    }

    respondJSON(w, http.StatusOK, usuarios)
}

// GetByID obtiene un usuario por ID
func (h *UsuarioHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    usuario, err := h.usuarioService.GetByID(id)
    if err != nil {
        respondError(w, http.StatusNotFound, err.Error())
        return
    }

    respondJSON(w, http.StatusOK, usuario)
}

// Update actualiza un usuario
func (h *UsuarioHandler) Update(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)

    // Un usuario solo puede actualizar su propio perfil
    if claims.UserID != id {
        respondError(w, http.StatusForbidden, "No tienes permiso para actualizar este usuario")
        return
    }

    var usuario models.Usuario
    if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
        respondError(w, http.StatusBadRequest, "Datos inválidos")
        return
    }

    updatedUsuario, err := h.usuarioService.Update(id, &usuario)
    if err != nil {
        respondError(w, http.StatusBadRequest, err.Error())
        return
    }

    respondJSON(w, http.StatusOK, map[string]interface{}{
        "message": "Usuario actualizado exitosamente",
        "usuario": updatedUsuario,
    })
}

// Delete elimina un usuario
func (h *UsuarioHandler) Delete(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)

    // Un usuario solo puede eliminar su propio perfil
    if claims.UserID != id {
        respondError(w, http.StatusForbidden, "No tienes permiso para eliminar este usuario")
        return
    }

    err = h.usuarioService.Delete(id)
    if err != nil {
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondJSON(w, http.StatusOK, map[string]string{
        "message": "Usuario eliminado exitosamente",
    })
}

// ChangePassword cambia la contraseña del usuario
func (h *UsuarioHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
    claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)

    var req struct {
        OldPassword string `json:"old_password"`
        NewPassword string `json:"new_password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "Datos inválidos")
        return
    }

    err := h.usuarioService.ChangePassword(claims.UserID, req.OldPassword, req.NewPassword)
    if err != nil {
        respondError(w, http.StatusBadRequest, err.Error())
        return
    }

    respondJSON(w, http.StatusOK, map[string]string{
        "message": "Contraseña actualizada exitosamente",
    })
}
