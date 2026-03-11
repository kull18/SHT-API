package repository

import (
	"cursos-api/config"
	"cursos-api/models"
	"database/sql"
	"errors"
	"time"
)

type InscripcionRepository struct{}

func NewInscripcionRepository() *InscripcionRepository {
	return &InscripcionRepository{}
}

// Create crea una nueva inscripción
func (r *InscripcionRepository) Create(inscripcion *models.Inscripcion) error {
	query := `
		INSERT INTO inscripciones (usuario_id, curso_id, fecha_inscripcion, estado, progreso_porcentaje)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, fecha_inscripcion
	`

	now := time.Now()
	err := config.DB.QueryRow(
		query,
		inscripcion.UsuarioID,
		inscripcion.CursoID,
		now,
		"activo",
		0.00,
	).Scan(&inscripcion.ID, &inscripcion.FechaInscripcion)

	return err
}

// FindByID busca una inscripción por ID
func (r *InscripcionRepository) FindByID(id int) (*models.Inscripcion, error) {
	query := `
		SELECT i.id, i.usuario_id, i.curso_id, i.fecha_inscripcion, i.estado, i.progreso_porcentaje,
		       c.id, c.nombre, c.descripcion, c.duracion_horas, c.instructor_id, c.activo
		FROM inscripciones i
		INNER JOIN cursos c ON i.curso_id = c.id
		WHERE i.id = $1
	`

	inscripcion := &models.Inscripcion{Curso: &models.Curso{}}
	err := config.DB.QueryRow(query, id).Scan(
		&inscripcion.ID,
		&inscripcion.UsuarioID,
		&inscripcion.CursoID,
		&inscripcion.FechaInscripcion,
		&inscripcion.Estado,
		&inscripcion.ProgresoPorcentaje,
		&inscripcion.Curso.ID,
		&inscripcion.Curso.Nombre,
		&inscripcion.Curso.Descripcion,
		&inscripcion.Curso.DuracionHoras,
		&inscripcion.Curso.InstructorID,
		&inscripcion.Curso.Activo,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("inscripción no encontrada")
	}

	return inscripcion, err
}

// FindByUsuario obtiene todas las inscripciones de un usuario
func (r *InscripcionRepository) FindByUsuario(usuarioID int) ([]models.Inscripcion, error) {
	query := `
		SELECT i.id, i.usuario_id, i.curso_id, i.fecha_inscripcion, i.estado, i.progreso_porcentaje,
		       c.id, c.nombre, c.descripcion, c.duracion_horas, c.instructor_id, c.activo
		FROM inscripciones i
		INNER JOIN cursos c ON i.curso_id = c.id
		WHERE i.usuario_id = $1
		ORDER BY i.fecha_inscripcion DESC
	`

	rows, err := config.DB.Query(query, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inscripciones []models.Inscripcion
	for rows.Next() {
		var inscripcion models.Inscripcion
		inscripcion.Curso = &models.Curso{}

		err := rows.Scan(
			&inscripcion.ID,
			&inscripcion.UsuarioID,
			&inscripcion.CursoID,
			&inscripcion.FechaInscripcion,
			&inscripcion.Estado,
			&inscripcion.ProgresoPorcentaje,
			&inscripcion.Curso.ID,
			&inscripcion.Curso.Nombre,
			&inscripcion.Curso.Descripcion,
			&inscripcion.Curso.DuracionHoras,
			&inscripcion.Curso.InstructorID,
			&inscripcion.Curso.Activo,
		)
		if err != nil {
			return nil, err
		}

		inscripciones = append(inscripciones, inscripcion)
	}

	return inscripciones, nil
}

// FindByUsuarioYCurso verifica si un usuario ya está inscrito en un curso
func (r *InscripcionRepository) FindByUsuarioYCurso(usuarioID, cursoID int) (*models.Inscripcion, error) {
	query := `
		SELECT id, usuario_id, curso_id, fecha_inscripcion, estado, progreso_porcentaje
		FROM inscripciones
		WHERE usuario_id = $1 AND curso_id = $2
	`

	inscripcion := &models.Inscripcion{}
	err := config.DB.QueryRow(query, usuarioID, cursoID).Scan(
		&inscripcion.ID,
		&inscripcion.UsuarioID,
		&inscripcion.CursoID,
		&inscripcion.FechaInscripcion,
		&inscripcion.Estado,
		&inscripcion.ProgresoPorcentaje,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No está inscrito, no es error
	}

	return inscripcion, err
}

// UpdateProgreso actualiza el progreso de una inscripción
func (r *InscripcionRepository) UpdateProgreso(inscripcionID int, progreso float64) error {
	query := `
		UPDATE inscripciones
		SET progreso_porcentaje = $1
		WHERE id = $2
	`

	result, err := config.DB.Exec(query, progreso, inscripcionID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("inscripción no encontrada")
	}

	return nil
}

// UpdateEstado actualiza el estado de una inscripción
func (r *InscripcionRepository) UpdateEstado(inscripcionID int, estado string) error {
	query := `
		UPDATE inscripciones
		SET estado = $1
		WHERE id = $2
	`

	result, err := config.DB.Exec(query, estado, inscripcionID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("inscripción no encontrada")
	}

	return nil
}

// CreateProgresLeccion registra o actualiza el progreso de una lección
func (r *InscripcionRepository) CreateProgresoLeccion(progreso *models.ProgresoLeccion) error {
	query := `
		INSERT INTO progreso_lecciones (inscripcion_id, leccion_id, completada, fecha_completado)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (inscripcion_id, leccion_id)
		DO UPDATE SET completada = $3, fecha_completado = $4
		RETURNING id
	`

	var fechaCompletado interface{}
	if progreso.Completada {
		now := time.Now()
		fechaCompletado = now
		progreso.FechaCompletado = &now
	} else {
		fechaCompletado = nil
		progreso.FechaCompletado = nil
	}

	err := config.DB.QueryRow(
		query,
		progreso.InscripcionID,
		progreso.LeccionID,
		progreso.Completada,
		fechaCompletado,
	).Scan(&progreso.ID)

	return err
}

// GetProgresoLecciones obtiene el progreso de todas las lecciones de una inscripción
func (r *InscripcionRepository) GetProgresoLecciones(inscripcionID int) ([]models.ProgresoLeccion, error) {
	query := `
		SELECT id, inscripcion_id, leccion_id, completada, fecha_completado
		FROM progreso_lecciones
		WHERE inscripcion_id = $1
	`

	rows, err := config.DB.Query(query, inscripcionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progresos []models.ProgresoLeccion
	for rows.Next() {
		var progreso models.ProgresoLeccion
		err := rows.Scan(
			&progreso.ID,
			&progreso.InscripcionID,
			&progreso.LeccionID,
			&progreso.Completada,
			&progreso.FechaCompletado,
		)
		if err != nil {
			return nil, err
		}
		progresos = append(progresos, progreso)
	}

	return progresos, nil
}

// CountLeccionesCompletadas cuenta las lecciones completadas de una inscripción
func (r *InscripcionRepository) CountLeccionesCompletadas(inscripcionID int) (int, error) {
	query := `
		SELECT COUNT(*) FROM progreso_lecciones
		WHERE inscripcion_id = $1 AND completada = true
	`

	var count int
	err := config.DB.QueryRow(query, inscripcionID).Scan(&count)
	return count, err
}