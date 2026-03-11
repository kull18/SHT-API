package services

import (
	"cursos-api/models"
	"cursos-api/repository"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type LeccionService struct {
	leccionRepo *repository.LeccionRepository
	cursoRepo   *repository.CursoRepository
}

func NewLeccionService() *LeccionService {
	return &LeccionService{
		leccionRepo: repository.NewLeccionRepository(),
		cursoRepo:   repository.NewCursoRepository(),
	}
}

// Create crea una nueva lección con imagen opcional
func (s *LeccionService) Create(cursoID int, req *models.CreateLeccionRequest, imagen multipart.File, imagenHeader *multipart.FileHeader, userID int, userRol string) (*models.Leccion, error) {
	// Solo instructores pueden crear lecciones
	if userRol != "instructor" {
		return nil, errors.New("solo los instructores pueden crear lecciones")
	}

	// Verificar que el curso existe y pertenece al instructor
	exists, err := s.cursoRepo.VerifyInstructor(cursoID, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("curso no encontrado o no tienes permiso para modificarlo")
	}

	// Validaciones
	if req.Titulo == "" {
		return nil, errors.New("el título de la lección es requerido")
	}
	if req.Orden <= 0 {
		return nil, errors.New("el orden debe ser mayor a 0")
	}
	if req.DuracionMinutos < 0 {
		return nil, errors.New("la duración no puede ser negativa")
	}

	leccion := &models.Leccion{
		CursoID:         cursoID,
		Titulo:          req.Titulo,
		Contenido:       req.Contenido,
		Orden:           req.Orden,
		DuracionMinutos: req.DuracionMinutos,
	}

	// Guardar imagen si viene en el request
	if imagen != nil && imagenHeader != nil {
		imagenURL, err := s.guardarImagen(imagen, imagenHeader)
		if err != nil {
			return nil, fmt.Errorf("error al guardar la imagen: %w", err)
		}
		leccion.ImagenURL = imagenURL
	}

	// Crear lección en BD
	err = s.leccionRepo.Create(leccion)
	if err != nil {
		return nil, err
	}

	return leccion, nil
}

// GetByCurso obtiene todas las lecciones de un curso
func (s *LeccionService) GetByCurso(cursoID int, userID int, userRol string) ([]models.Leccion, error) {
	// Verificar que el curso existe
	curso, err := s.cursoRepo.FindByID(cursoID)
	if err != nil {
		return nil, errors.New("curso no encontrado")
	}

	// Instructor solo puede ver lecciones de sus propios cursos
	if userRol == "instructor" && curso.InstructorID != userID {
		return nil, errors.New("no tienes permiso para ver las lecciones de este curso")
	}

	// Alumno solo puede ver lecciones de cursos activos
	if userRol == "alumno" && !curso.Activo {
		return nil, errors.New("curso no disponible")
	}

	return s.leccionRepo.GetByCurso(cursoID)
}

// GetByID obtiene una lección por ID
func (s *LeccionService) GetByID(cursoID, leccionID int, userID int, userRol string) (*models.Leccion, error) {
	// Verificar que la lección pertenece al curso
	exists, err := s.leccionRepo.VerifyCurso(leccionID, cursoID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("lección no encontrada en este curso")
	}

	// Verificar permisos sobre el curso
	curso, err := s.cursoRepo.FindByID(cursoID)
	if err != nil {
		return nil, errors.New("curso no encontrado")
	}

	if userRol == "instructor" && curso.InstructorID != userID {
		return nil, errors.New("no tienes permiso para ver esta lección")
	}

	if userRol == "alumno" && !curso.Activo {
		return nil, errors.New("curso no disponible")
	}

	return s.leccionRepo.FindByID(leccionID)
}

// Update actualiza una lección
func (s *LeccionService) Update(cursoID, leccionID int, req *models.UpdateLeccionRequest, userID int, userRol string) (*models.Leccion, error) {
	// Solo instructores pueden actualizar lecciones
	if userRol != "instructor" {
		return nil, errors.New("solo los instructores pueden actualizar lecciones")
	}

	// Verificar que el curso pertenece al instructor
	exists, err := s.cursoRepo.VerifyInstructor(cursoID, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("curso no encontrado o no tienes permiso para modificarlo")
	}

	// Verificar que la lección pertenece al curso
	belongs, err := s.leccionRepo.VerifyCurso(leccionID, cursoID)
	if err != nil {
		return nil, err
	}
	if !belongs {
		return nil, errors.New("lección no encontrada en este curso")
	}

	// Validaciones
	if req.Titulo == "" {
		return nil, errors.New("el título de la lección es requerido")
	}
	if req.Orden <= 0 {
		return nil, errors.New("el orden debe ser mayor a 0")
	}

	leccion := &models.Leccion{
		Titulo:          req.Titulo,
		Contenido:       req.Contenido,
		Orden:           req.Orden,
		DuracionMinutos: req.DuracionMinutos,
	}

	err = s.leccionRepo.Update(leccionID, leccion)
	if err != nil {
		return nil, err
	}

	return leccion, nil
}

// Delete elimina una lección
func (s *LeccionService) Delete(cursoID, leccionID int, userID int, userRol string) error {
	// Solo instructores pueden eliminar lecciones
	if userRol != "instructor" {
		return errors.New("solo los instructores pueden eliminar lecciones")
	}

	// Verificar que el curso pertenece al instructor
	exists, err := s.cursoRepo.VerifyInstructor(cursoID, userID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("curso no encontrado o no tienes permiso para modificarlo")
	}

	// Verificar que la lección pertenece al curso
	belongs, err := s.leccionRepo.VerifyCurso(leccionID, cursoID)
	if err != nil {
		return err
	}
	if !belongs {
		return errors.New("lección no encontrada en este curso")
	}

	return s.leccionRepo.Delete(leccionID)
}

// guardarImagen guarda la imagen en la carpeta uploads y retorna la URL
func (s *LeccionService) guardarImagen(imagen multipart.File, header *multipart.FileHeader) (string, error) {
	// Crear carpeta uploads si no existe
	uploadDir := "./uploads/lecciones"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("error al crear directorio: %w", err)
	}

	// Generar nombre único para la imagen
	ext := filepath.Ext(header.Filename)
	nombreArchivo := fmt.Sprintf("leccion_%d%s", time.Now().UnixNano(), ext)
	rutaArchivo := filepath.Join(uploadDir, nombreArchivo)

	// Crear archivo destino
	destino, err := os.Create(rutaArchivo)
	if err != nil {
		return "", fmt.Errorf("error al crear archivo: %w", err)
	}
	defer destino.Close()

	// Copiar contenido
	if _, err := io.Copy(destino, imagen); err != nil {
		return "", fmt.Errorf("error al guardar imagen: %w", err)
	}

	// Retornar URL relativa
	return fmt.Sprintf("/uploads/lecciones/%s", nombreArchivo), nil
}