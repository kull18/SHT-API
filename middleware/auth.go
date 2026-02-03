package middleware

import (
    "context"
    "cursos-api/utils"
    "net/http"
    "strings"
)

type contextKey string

const UserContextKey contextKey = "user"

// AuthMiddleware verifica el token JWT
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, `{"error":"Token no proporcionado"}`, http.StatusUnauthorized)
            return
        }

        // Formato: "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, `{"error":"Formato de token inválido"}`, http.StatusUnauthorized)
            return
        }

        claims, err := utils.ValidateJWT(parts[1])
        if err != nil {
            http.Error(w, `{"error":"Token inválido o expirado"}`, http.StatusUnauthorized)
            return
        }

        // Agregar claims al contexto
        ctx := context.WithValue(r.Context(), UserContextKey, claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}

// RoleMiddleware verifica que el usuario tenga un rol específico
func RoleMiddleware(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
    return AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
        claims := r.Context().Value(UserContextKey).(*utils.Claims)

        if claims.Rol != requiredRole {
            http.Error(w, `{"error":"Acceso no autorizado"}`, http.StatusForbidden)
            return
        }

        next.ServeHTTP(w, r)
    })
}

// CORS Middleware
func CORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}
