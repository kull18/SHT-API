package repository

import (
	"cursos-api/config"
	"cursos-api/models"
	"database/sql"
	"errors"
	"time"
)

type LeccionRepository struct{}

func NewLeccionRepository() *LeccionRepository {
	return &LeccionRepository{}
}

// Create crea una nueva lección
func (r *LeccionRepository) Create(leccion *models.Leccion) error {
	query := `
		INSERT INTO lecciones (curso_id, titulo, contenido, imagen_url, orden, duracion_minutos, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`

	now := time.Now()
	err := config.DB.QueryRow(
		query,
		leccion.CursoID,
		leccion.Titulo,
		leccion.Contenido,
		leccion.ImagenURL,
		leccion.Orden,
		leccion.DuracionMinutos,
		now,
	).Scan(&leccion.ID, &leccion.CreatedAt)

	return err
}

// FindByID busca una lección por ID
func (r *LeccionRepository) FindByID(id int) (*models.Leccion, error) {
	query := `
		SELECT id, curso_id, titulo, contenido, imagen_url, orden, duracion_minutos, created_at
		FROM lecciones
		WHERE id = $1
	`

	leccion := &models.Leccion{}
	var imagenURL sql.NullString

	err := config.DB.QueryRow(query, id).Scan(
		&leccion.ID,
		&leccion.CursoID,
		&leccion.Titulo,
		&leccion.Contenido,
		&imagenURL,
		&leccion.Orden,
		&leccion.DuracionMinutos,
		&leccion.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("lección no encontrada")
	}

	if imagenURL.Valid {
		leccion.ImagenURL = imagenURL.String
	}

	return leccion, err
}

// GetByCurso obtiene todas las lecciones de un curso ordenadas
func (r *LeccionRepository) GetByCurso(cursoID int) ([]models.Leccion, error) {
	query := `
		SELECT id, curso_id, titulo, contenido, imagen_url, orden, duracion_minutos, created_at
		FROM lecciones
		WHERE curso_id = $1
		ORDER BY orden ASC
	`

	rows, err := config.DB.Query(query, cursoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lecciones []models.Leccion
	for rows.Next() {
		var leccion models.Leccion
		var imagenURL sql.NullString

		err := rows.Scan(
			&leccion.ID,
			&leccion.CursoID,
			&leccion.Titulo,
			&leccion.Contenido,
			&imagenURL,
			&leccion.Orden,
			&leccion.DuracionMinutos,
			&leccion.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if imagenURL.Valid {
			leccion.ImagenURL = imagenURL.String
		}

		lecciones = append(lecciones, leccion)
	}

	return lecciones, nil
}

// Update actualiza una lección
func (r *LeccionRepository) Update(id int, leccion *models.Leccion) error {
	query := `
		UPDATE lecciones
		SET titulo = $1, contenido = $2, imagen_url = $3, orden = $4, duracion_minutos = $5, updated_at = $6
		WHERE id = $7
		RETURNING updated_at
	`

	now := time.Now()
	err := config.DB.QueryRow(
		query,
		leccion.Titulo,
		leccion.Contenido,
		leccion.ImagenURL,
		leccion.Orden,
		leccion.DuracionMinutos,
		now,
		id,
	).Scan(&leccion.UpdatedAt)

	if err == sql.ErrNoRows {
		return errors.New("lección no encontrada")
	}

	leccion.ID = id
	return err
}

// Delete elimina una lección
func (r *LeccionRepository) Delete(id int) error {
	query := `DELETE FROM lecciones WHERE id = $1`

	result, err := config.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("lección no encontrada")
	}

	return nil
}

// VerifyCurso verifica que una lección pertenece a un curso
func (r *LeccionRepository) VerifyCurso(leccionID, cursoID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM lecciones WHERE id = $1 AND curso_id = $2)`

	var exists bool
	err := config.DB.QueryRow(query, leccionID, cursoID).Scan(&exists)

	return exists, err
}

// CountByCurso cuenta el número de lecciones de un curso
func (r *LeccionRepository) CountByCurso(cursoID int) (int, error) {
	query := `SELECT COUNT(*) FROM lecciones WHERE curso_id = $1`

	var count int
	err := config.DB.QueryRow(query, cursoID).Scan(&count)

	return count, err
}