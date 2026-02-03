package config

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")
    sslmode := os.Getenv("DB_SSLMODE")

    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        host, port, user, password, dbname, sslmode,
    )

    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Error al conectar a la base de datos:", err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatal("Error al hacer ping a la base de datos:", err)
    }

    log.Println("✅ Conexión exitosa a la base de datos")
}

func CloseDB() {
    if DB != nil {
        DB.Close()
    }
}
