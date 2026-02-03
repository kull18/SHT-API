package services

import (
    "cursos-api/models"
    "cursos-api/repository"
    "errors"
)

type CursoService struct {
    cursoRepo   *repository.CursoRepository
    usuarioRepo *repository.UsuarioRepository
}

func NewCursoService() *CursoService {
    return &CursoService{
        cursoRepo:   repository.NewCursoRepository(),
        usuarioRepo: repository.NewUsuarioRepository(),
    }
}

// Create crea un nuevo curso
func (s *CursoService) Create(curso *models.Curso, userID int, userRol string) (*models.Curso, error) {
    // Validaciones
    if curso.Nombre == "" {
        return nil, errors.New("el nombre del curso es requerido")
    }

    if curso.DuracionHoras <= 0 {
        return nil, errors.New("la duración debe ser mayor a 0")
    }

    // Solo instructores pueden crear cursos
    if userRol != "instructor" {
        return nil, errors.New("solo los instructores pueden crear cursos")
    }

    // Verificar que el instructor existe
    instructor, err := s.usuarioRepo.FindByID(curso.InstructorID)
    if err != nil {
        return nil, errors.New("instructor no encontrado")
    }

    // Verificar que el instructor es realmente un instructor
    if instructor.Rol != "instructor" {
        return nil, errors.New("el usuario especificado no es un instructor")
    }

    // El instructor solo puede crear cursos para sí mismo (a menos que sea admin en el futuro)
    if curso.InstructorID != userID {
        return nil, errors.New("no puedes crear cursos para otros instructores")
    }

    // Establecer activo por defecto
    curso.Activo = true

    // Crear curso
    err = s.cursoRepo.Create(curso)
    if err != nil {
        return nil, err
    }

    return curso, nil
}

// GetAll obtiene todos los cursos
func (s *CursoService) GetAll(userRol string, userID int) ([]models.Curso, error) {
    // Los instructores solo ven sus propios cursos
    if userRol == "instructor" {
        return s.cursoRepo.GetByInstructor(userID)
    }

    // Los alumnos ven todos los cursos activos
    return s.cursoRepo.GetActivos()
}

// GetByID obtiene un curso por ID
func (s *CursoService) GetByID(id int, userRol string, userID int) (*models.Curso, error) {
    curso, err := s.cursoRepo.FindByID(id)
    if err != nil {
        return nil, err
    }

    // Los instructores solo pueden ver sus propios cursos
    if userRol == "instructor" && curso.InstructorID != userID {
        return nil, errors.New("no tienes permiso para ver este curso")
    }

    // Los alumnos solo pueden ver cursos activos
    if userRol == "alumno" && !curso.Activo {
        return nil, errors.New("curso no disponible")
    }

    return curso, nil
}

// GetMyCursos obtiene los cursos de un instructor
func (s *CursoService) GetMyCursos(instructorID int) ([]models.Curso, error) {
    return s.cursoRepo.GetByInstructor(instructorID)
}

// Update actualiza un curso
func (s *CursoService) Update(id int, curso *models.Curso, userID int, userRol string) (*models.Curso, error) {
    // Validaciones
    if curso.Nombre == "" {
        return nil, errors.New("el nombre del curso es requerido")
    }

    if curso.DuracionHoras <= 0 {
        return nil, errors.New("la duración debe ser mayor a 0")
    }

    // Solo instructores pueden actualizar cursos
    if userRol != "instructor" {
        return nil, errors.New("solo los instructores pueden actualizar cursos")
    }

    // Verificar que el curso existe y pertenece al instructor
    exists, err := s.cursoRepo.VerifyInstructor(id, userID)
    if err != nil {
        return nil, err
    }

    if !exists {
        return nil, errors.New("curso no encontrado o no tienes permiso para modificarlo")
    }

    // El instructor_id no puede cambiar
    curso.InstructorID = userID

    // Actualizar
    err = s.cursoRepo.Update(id, curso)
    if err != nil {
        return nil, err
    }

    return curso, nil
}

// Delete elimina un curso
func (s *CursoService) Delete(id int, userID int, userRol string) error {
    // Solo instructores pueden eliminar cursos
    if userRol != "instructor" {
        return errors.New("solo los instructores pueden eliminar cursos")
    }

    // Verificar que el curso existe y pertenece al instructor
    exists, err := s.cursoRepo.VerifyInstructor(id, userID)
    if err != nil {
        return err
    }

    if !exists {
        return errors.New("curso no encontrado o no tienes permiso para eliminarlo")
    }

    return s.cursoRepo.Delete(id)
}

// ToggleActivo activa o desactiva un curso
func (s *CursoService) ToggleActivo(id int, userID int, userRol string) (*models.Curso, error) {
    // Solo instructores pueden cambiar el estado
    if userRol != "instructor" {
        return nil, errors.New("solo los instructores pueden cambiar el estado del curso")
    }

    // Verificar que el curso existe y pertenece al instructor
    exists, err := s.cursoRepo.VerifyInstructor(id, userID)
    if err != nil {
        return nil, err
    }

    if !exists {
        return nil, errors.New("curso no encontrado o no tienes permiso")
    }

    // Obtener curso actual
    curso, err := s.cursoRepo.FindByID(id)
    if err != nil {
        return nil, err
    }

    // Cambiar estado
    curso.Activo = !curso.Activo

    // Actualizar
    err = s.cursoRepo.Update(id, curso)
    if err != nil {
        return nil, err
    }

    return curso, nil
}
