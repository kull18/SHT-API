package repository

import (
    "cursos-api/config"
    "cursos-api/models"
    "database/sql"
    "errors"
    "time"
)

type CursoRepository struct{}

func NewCursoRepository() *CursoRepository {
    return &CursoRepository{}
}

// Create crea un nuevo curso
func (r *CursoRepository) Create(curso *models.Curso) error {
    query := `
        INSERT INTO cursos (nombre, descripcion, duracion_horas, instructor_id, activo, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at
    `
    
    now := time.Now()
    err := config.DB.QueryRow(
        query,
        curso.Nombre,
        curso.Descripcion,
        curso.DuracionHoras,
        curso.InstructorID,
        curso.Activo,
        now,
        now,
    ).Scan(&curso.ID, &curso.CreatedAt, &curso.UpdatedAt)

    return err
}

// FindByID busca un curso por ID
func (r *CursoRepository) FindByID(id int) (*models.Curso, error) {
    query := `
        SELECT c.id, c.nombre, c.descripcion, c.duracion_horas, c.instructor_id, c.activo, c.created_at, c.updated_at,
               u.id, u.nombre, u.email, u.rol
        FROM cursos c
        INNER JOIN usuarios u ON c.instructor_id = u.id
        WHERE c.id = $1
    `
    
    curso := &models.Curso{Instructor: &models.Usuario{}}
    err := config.DB.QueryRow(query, id).Scan(
        &curso.ID,
        &curso.Nombre,
        &curso.Descripcion,
        &curso.DuracionHoras,
        &curso.InstructorID,
        &curso.Activo,
        &curso.CreatedAt,
        &curso.UpdatedAt,
        &curso.Instructor.ID,
        &curso.Instructor.Nombre,
        &curso.Instructor.Email,
        &curso.Instructor.Rol,
    )

    if err == sql.ErrNoRows {
        return nil, errors.New("curso no encontrado")
    }

    return curso, err
}

// GetAll obtiene todos los cursos
func (r *CursoRepository) GetAll() ([]models.Curso, error) {
    query := `
        SELECT c.id, c.nombre, c.descripcion, c.duracion_horas, c.instructor_id, c.activo, c.created_at, c.updated_at,
               u.id, u.nombre, u.email, u.rol
        FROM cursos c
        INNER JOIN usuarios u ON c.instructor_id = u.id
        ORDER BY c.created_at DESC
    `
    
    rows, err := config.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var cursos []models.Curso
    for rows.Next() {
        var curso models.Curso
        curso.Instructor = &models.Usuario{}
        
        err := rows.Scan(
            &curso.ID,
            &curso.Nombre,
            &curso.Descripcion,
            &curso.DuracionHoras,
            &curso.InstructorID,
            &curso.Activo,
            &curso.CreatedAt,
            &curso.UpdatedAt,
            &curso.Instructor.ID,
            &curso.Instructor.Nombre,
            &curso.Instructor.Email,
            &curso.Instructor.Rol,
        )
        if err != nil {
            return nil, err
        }
        cursos = append(cursos, curso)
    }

    return cursos, nil
}

// GetByInstructor obtiene todos los cursos de un instructor
func (r *CursoRepository) GetByInstructor(instructorID int) ([]models.Curso, error) {
    query := `
        SELECT c.id, c.nombre, c.descripcion, c.duracion_horas, c.instructor_id, c.activo, c.created_at, c.updated_at,
               u.id, u.nombre, u.email, u.rol
        FROM cursos c
        INNER JOIN usuarios u ON c.instructor_id = u.id
        WHERE c.instructor_id = $1
        ORDER BY c.created_at DESC
    `
    
    rows, err := config.DB.Query(query, instructorID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var cursos []models.Curso
    for rows.Next() {
        var curso models.Curso
        curso.Instructor = &models.Usuario{}
        
        err := rows.Scan(
            &curso.ID,
            &curso.Nombre,
            &curso.Descripcion,
            &curso.DuracionHoras,
            &curso.InstructorID,
            &curso.Activo,
            &curso.CreatedAt,
            &curso.UpdatedAt,
            &curso.Instructor.ID,
            &curso.Instructor.Nombre,
            &curso.Instructor.Email,
            &curso.Instructor.Rol,
        )
        if err != nil {
            return nil, err
        }
        cursos = append(cursos, curso)
    }

    return cursos, nil
}

// GetActivos obtiene todos los cursos activos
func (r *CursoRepository) GetActivos() ([]models.Curso, error) {
    query := `
        SELECT c.id, c.nombre, c.descripcion, c.duracion_horas, c.instructor_id, c.activo, c.created_at, c.updated_at,
               u.id, u.nombre, u.email, u.rol
        FROM cursos c
        INNER JOIN usuarios u ON c.instructor_id = u.id
        WHERE c.activo = true
        ORDER BY c.created_at DESC
    `
    
    rows, err := config.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var cursos []models.Curso
    for rows.Next() {
        var curso models.Curso
        curso.Instructor = &models.Usuario{}
        
        err := rows.Scan(
            &curso.ID,
            &curso.Nombre,
            &curso.Descripcion,
            &curso.DuracionHoras,
            &curso.InstructorID,
            &curso.Activo,
            &curso.CreatedAt,
            &curso.UpdatedAt,
            &curso.Instructor.ID,
            &curso.Instructor.Nombre,
            &curso.Instructor.Email,
            &curso.Instructor.Rol,
        )
        if err != nil {
            return nil, err
        }
        cursos = append(cursos, curso)
    }

    return cursos, nil
}

// Update actualiza un curso
func (r *CursoRepository) Update(id int, curso *models.Curso) error {
    query := `
        UPDATE cursos
        SET nombre = $1, descripcion = $2, duracion_horas = $3, instructor_id = $4, activo = $5, updated_at = $6
        WHERE id = $7
        RETURNING updated_at
    `
    
    now := time.Now()
    err := config.DB.QueryRow(
        query,
        curso.Nombre,
        curso.Descripcion,
        curso.DuracionHoras,
        curso.InstructorID,
        curso.Activo,
        now,
        id,
    ).Scan(&curso.UpdatedAt)

    if err == sql.ErrNoRows {
        return errors.New("curso no encontrado")
    }

    curso.ID = id
    return err
}

// Delete elimina un curso
func (r *CursoRepository) Delete(id int) error {
    query := `DELETE FROM cursos WHERE id = $1`
    
    result, err := config.DB.Exec(query, id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return errors.New("curso no encontrado")
    }

    return nil
}

// VerifyInstructor verifica que un curso pertenece a un instructor
func (r *CursoRepository) VerifyInstructor(cursoID, instructorID int) (bool, error) {
    query := `SELECT EXISTS(SELECT 1 FROM cursos WHERE id = $1 AND instructor_id = $2)`
    
    var exists bool
    err := config.DB.QueryRow(query, cursoID, instructorID).Scan(&exists)
    
    return exists, err
}
