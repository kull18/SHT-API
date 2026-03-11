package services

import (
	"cursos-api/models"
	"cursos-api/repository"
	"errors"
	"fmt"
)

type InscripcionService struct {
	inscripcionRepo *repository.InscripcionRepository
	cursoRepo       *repository.CursoRepository
	leccionRepo     *repository.LeccionRepository
}

func NewInscripcionService() *InscripcionService {
	return &InscripcionService{
		inscripcionRepo: repository.NewInscripcionRepository(),
		cursoRepo:       repository.NewCursoRepository(),
		leccionRepo:     repository.NewLeccionRepository(),
	}
}

// Inscribir inscribe a un alumno en un curso
func (s *InscripcionService) Inscribir(cursoID int, userID int, userRol string) (*models.Inscripcion, error) {
	// Solo alumnos pueden inscribirse
	if userRol != "alumno" {
		return nil, errors.New("solo los alumnos pueden inscribirse a cursos")
	}

	// Verificar que el curso existe y está activo
	curso, err := s.cursoRepo.FindByID(cursoID)
	if err != nil {
		return nil, errors.New("curso no encontrado")
	}
	if !curso.Activo {
		return nil, errors.New("el curso no está disponible")
	}

	// Verificar que el alumno no esté ya inscrito
	existente, err := s.inscripcionRepo.FindByUsuarioYCurso(userID, cursoID)
	if err != nil {
		return nil, err
	}
	if existente != nil {
		return nil, errors.New("ya estás inscrito en este curso")
	}

	// Crear inscripción
	inscripcion := &models.Inscripcion{
		UsuarioID: userID,
		CursoID:   cursoID,
		Estado:    "activo",
	}

	err = s.inscripcionRepo.Create(inscripcion)
	if err != nil {
		return nil, err
	}

	inscripcion.Curso = curso
	return inscripcion, nil
}

// GetMisInscripciones obtiene todas las inscripciones del alumno autenticado
func (s *InscripcionService) GetMisInscripciones(userID int, userRol string) ([]models.Inscripcion, error) {
	// Solo alumnos pueden ver sus inscripciones
	if userRol != "alumno" {
		return nil, errors.New("solo los alumnos pueden ver sus inscripciones")
	}

	return s.inscripcionRepo.FindByUsuario(userID)
}

// GetByID obtiene una inscripción por ID verificando que pertenece al alumno
func (s *InscripcionService) GetByID(inscripcionID int, userID int, userRol string) (*models.Inscripcion, error) {
	inscripcion, err := s.inscripcionRepo.FindByID(inscripcionID)
	if err != nil {
		return nil, err
	}

	// Verificar que la inscripción pertenece al alumno autenticado
	if userRol == "alumno" && inscripcion.UsuarioID != userID {
		return nil, errors.New("no tienes permiso para ver esta inscripción")
	}

	return inscripcion, nil
}

// MarcarLeccion marca o desmarca una lección como completada y recalcula el progreso
func (s *InscripcionService) MarcarLeccion(inscripcionID, leccionID int, completada bool, userID int, userRol string) (*models.ProgresoLeccion, error) {
	// Solo alumnos pueden marcar lecciones
	if userRol != "alumno" {
		return nil, errors.New("solo los alumnos pueden marcar lecciones como completadas")
	}

	// Verificar que la inscripción existe y pertenece al alumno
	inscripcion, err := s.inscripcionRepo.FindByID(inscripcionID)
	if err != nil {
		return nil, err
	}
	if inscripcion.UsuarioID != userID {
		return nil, errors.New("no tienes permiso para modificar esta inscripción")
	}

	// Verificar que la lección pertenece al curso de la inscripción
	belongs, err := s.leccionRepo.VerifyCurso(leccionID, inscripcion.CursoID)
	if err != nil {
		return nil, err
	}
	if !belongs {
		return nil, errors.New("la lección no pertenece al curso de esta inscripción")
	}

	// Registrar progreso de la lección
	progreso := &models.ProgresoLeccion{
		InscripcionID: inscripcionID,
		LeccionID:     leccionID,
		Completada:    completada,
	}

	err = s.inscripcionRepo.CreateProgresoLeccion(progreso)
	if err != nil {
		return nil, err
	}

	// Recalcular progreso total
	err = s.recalcularProgreso(inscripcionID, inscripcion.CursoID)
	if err != nil {
		return nil, err
	}

	return progreso, nil
}

// GetProgreso obtiene el progreso de todas las lecciones de una inscripción
func (s *InscripcionService) GetProgreso(inscripcionID int, userID int, userRol string) ([]models.ProgresoLeccion, error) {
	// Verificar que la inscripción existe y pertenece al alumno
	inscripcion, err := s.inscripcionRepo.FindByID(inscripcionID)
	if err != nil {
		return nil, err
	}

	if userRol == "alumno" && inscripcion.UsuarioID != userID {
		return nil, errors.New("no tienes permiso para ver este progreso")
	}

	return s.inscripcionRepo.GetProgresoLecciones(inscripcionID)
}

// recalcularProgreso recalcula y actualiza el porcentaje de progreso de una inscripción
func (s *InscripcionService) recalcularProgreso(inscripcionID, cursoID int) error {
	// Total de lecciones del curso
	totalLecciones, err := s.leccionRepo.CountByCurso(cursoID)
	if err != nil {
		return err
	}

	if totalLecciones == 0 {
		return nil
	}

	// Lecciones completadas
	completadas, err := s.inscripcionRepo.CountLeccionesCompletadas(inscripcionID)
	if err != nil {
		return err
	}

	// Calcular porcentaje
	porcentaje := (float64(completadas) / float64(totalLecciones)) * 100

	// Actualizar progreso
	err = s.inscripcionRepo.UpdateProgreso(inscripcionID, porcentaje)
	if err != nil {
		return err
	}

	// Si el progreso es 100% marcar como completado
	if porcentaje == 100 {
		err = s.inscripcionRepo.UpdateEstado(inscripcionID, "completado")
		if err != nil {
			return fmt.Errorf("error al actualizar estado: %w", err)
		}
	}

	return nil
}