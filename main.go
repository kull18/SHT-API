package main

import (
    "cursos-api/config"
    "cursos-api/middleware"
    "cursos-api/routes"
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"
)

func main() {
    // Cargar variables de entorno
    err := godotenv.Load()
    if err != nil {
        log.Println("⚠️  No se encontró archivo .env, usando variables de entorno del sistema")
    }

    // Conectar a la base de datos
    config.ConnectDB()
    defer config.CloseDB()

    // Configurar rutas
    router := routes.SetupRoutes()

    // Aplicar middleware CORS
    handler := middleware.CORS(router)

    // Obtener puerto
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // Iniciar servidor
    log.Printf("🚀 Servidor iniciado en http://localhost:%s\n", port)
    log.Printf("📚 API de Gestión de Cursos\n")
    log.Printf("📖 Documentación: http://localhost:%s/health\n", port)
    
    if err := http.ListenAndServe(":"+port, handler); err != nil {
        log.Fatal("Error al iniciar el servidor:", err)
    }
}
