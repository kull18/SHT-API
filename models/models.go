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
	ImagenURL       string    `json:"imagen_url,omitempty"` // URL de imagen subida al servidor
	Orden           int       `json:"orden"`
	DuracionMinutos int       `json:"duracion_minutos"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
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
	ID              int        `json:"id"`
	InscripcionID   int        `json:"inscripcion_id"`
	LeccionID       int        `json:"leccion_id"`
	Completada      bool       `json:"completada"`
	FechaCompletado *time.Time `json:"fecha_completado,omitempty"`
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

type CreateLeccionRequest struct {
	Titulo          string `json:"titulo"`
	Contenido       string `json:"contenido"`
	Orden           int    `json:"orden"`
	DuracionMinutos int    `json:"duracion_minutos"`
}

type UpdateLeccionRequest struct {
	Titulo          string `json:"titulo"`
	Contenido       string `json:"contenido"`
	Orden           int    `json:"orden"`
	DuracionMinutos int    `json:"duracion_minutos"`
}

type InscripcionRequest struct {
	CursoID int `json:"curso_id"`
}

type ProgresoRequest struct {
	LeccionID  int  `json:"leccion_id"`
	Completada bool `json:"completada"`
}