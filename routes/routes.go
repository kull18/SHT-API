package routes

import (
	"cursos-api/handlers"
	"cursos-api/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
    router := mux.NewRouter()

    // Handlers
    authHandler := handlers.NewAuthHandler()
    usuarioHandler := handlers.NewUsuarioHandler()
    cursoHandler := handlers.NewCursoHandler()

    // API prefix
    api := router.PathPrefix("/api").Subrouter()

    // ============================================
    // RUTAS PÚBLICAS (sin autenticación)
    // ============================================
    api.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
    api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

    // ============================================
    // RUTAS PROTEGIDAS (requieren autenticación)
    // ============================================

    // --- Perfil de usuario ---
    api.HandleFunc("/auth/profile", middleware.AuthMiddleware(authHandler.GetProfile)).Methods("GET")

    // --- Usuarios ---
    api.HandleFunc("/usuarios", middleware.AuthMiddleware(usuarioHandler.GetAll)).Methods("GET")
    api.HandleFunc("/usuarios/{id}", middleware.AuthMiddleware(usuarioHandler.GetByID)).Methods("GET")
    api.HandleFunc("/usuarios/{id}", middleware.AuthMiddleware(usuarioHandler.Update)).Methods("PUT")
    api.HandleFunc("/usuarios/{id}", middleware.AuthMiddleware(usuarioHandler.Delete)).Methods("DELETE")
    api.HandleFunc("/usuarios/change-password", middleware.AuthMiddleware(usuarioHandler.ChangePassword)).Methods("POST")

    // --- Cursos ---
    // Rutas para instructores
    api.HandleFunc("/cursos", middleware.RoleMiddleware("instructor", cursoHandler.Create)).Methods("POST")
    api.HandleFunc("/cursos/my-cursos", middleware.RoleMiddleware("instructor", cursoHandler.GetMyCursos)).Methods("GET")
    api.HandleFunc("/cursos/{id}", middleware.RoleMiddleware("instructor", cursoHandler.Update)).Methods("PUT")
    api.HandleFunc("/cursos/{id}", middleware.RoleMiddleware("instructor", cursoHandler.Delete)).Methods("DELETE")
    api.HandleFunc("/cursos/{id}/toggle-activo", middleware.RoleMiddleware("instructor", cursoHandler.ToggleActivo)).Methods("PATCH")

    // Rutas disponibles para todos los usuarios autenticados
    api.HandleFunc("/cursos", middleware.AuthMiddleware(cursoHandler.GetAll)).Methods("GET")
    api.HandleFunc("/cursos/{id}", middleware.AuthMiddleware(cursoHandler.GetByID)).Methods("GET")

    // ============================================
    // RUTA DE SALUD
    // ============================================
    router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(200)
        w.Write([]byte(`{"status":"OK","message":"API de Cursos funcionando correctamente"}`))
    }).Methods("GET")

    return router
}
