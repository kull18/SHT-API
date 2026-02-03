package models

import "time"

type Usuario struct {
    ID           int       `json:"id"`
    Nombre       string    `json:"nombre"`
    Email        string    `json:"email"`
    Password     string    `json:"password,omitempty"`
    PasswordHash string    `json:"-"`
    Rol          string    `json:"rol"` // "instructor" o "alumno"
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

type Curso struct {
    ID            int       `json:"id"`
    Nombre        string    `json:"nombre"`
    Descripcion   string    `json:"descripcion"`
    DuracionHoras int       `json:"duracion_horas"`
    InstructorID  int       `json:"instructor_id"`
    Instructor    *Usuario  `json:"instructor,omitempty"`
    Activo        bool      `json:"activo"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

type Leccion struct {
    ID              int       `json:"id"`
    CursoID         int       `json:"curso_id"`
    Titulo          string    `json:"titulo"`
    Contenido       string    `json:"contenido"`
    Orden           int       `json:"orden"`
    DuracionMinutos int       `json:"duracion_minutos"`
    CreatedAt       time.Time `json:"created_at"`
}

type Inscripcion struct {
    ID                 int       `json:"id"`
    UsuarioID          int       `json:"usuario_id"`
    CursoID            int       `json:"curso_id"`
    FechaInscripcion   time.Time `json:"fecha_inscripcion"`
    Estado             string    `json:"estado"` // "activo", "completado", "cancelado"
    ProgresoPorcentaje float64   `json:"progreso_porcentaje"`
    Curso              *Curso    `json:"curso,omitempty"`
}

type ProgresoLeccion struct {
    ID               int        `json:"id"`
    InscripcionID    int        `json:"inscripcion_id"`
    LeccionID        int        `json:"leccion_id"`
    Completada       bool       `json:"completada"`
    FechaCompletado  *time.Time `json:"fecha_completado,omitempty"`
}

type Evaluacion struct {
    ID                 int       `json:"id"`
    CursoID            int       `json:"curso_id"`
    Titulo             string    `json:"titulo"`
    Descripcion        string    `json:"descripcion"`
    CalificacionMinima float64   `json:"calificacion_minima"`
    CreatedAt          time.Time `json:"created_at"`
}

type ResultadoEvaluacion struct {
    ID              int       `json:"id"`
    EvaluacionID    int       `json:"evaluacion_id"`
    UsuarioID       int       `json:"usuario_id"`
    Calificacion    float64   `json:"calificacion"`
    Aprobado        bool      `json:"aprobado"`
    FechaEvaluacion time.Time `json:"fecha_evaluacion"`
}

type Certificado struct {
    ID                 int       `json:"id"`
    InscripcionID      int       `json:"inscripcion_id"`
    CodigoCertificado  string    `json:"codigo_certificado"`
    FechaEmision       time.Time `json:"fecha_emision"`
    URLPDF             string    `json:"url_pdf,omitempty"`
}

// DTOs para requests
type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type RegisterRequest struct {
    Nombre   string `json:"nombre"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Rol      string `json:"rol"`
}

type LoginResponse struct {
    Token   string   `json:"token"`
    Usuario *Usuario `json:"usuario"`
}
