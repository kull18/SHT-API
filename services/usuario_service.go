package services

import (
    "cursos-api/models"
    "cursos-api/repository"
    "cursos-api/utils"
    "errors"
)

type UsuarioService struct {
    usuarioRepo *repository.UsuarioRepository
}

func NewUsuarioService() *UsuarioService {
    return &UsuarioService{
        usuarioRepo: repository.NewUsuarioRepository(),
    }
}

// GetAll obtiene todos los usuarios
func (s *UsuarioService) GetAll() ([]models.Usuario, error) {
    usuarios, err := s.usuarioRepo.GetAll()
    if err != nil {
        return nil, err
    }

    // Limpiar password hashes
    for i := range usuarios {
        usuarios[i].PasswordHash = ""
    }

    return usuarios, nil
}

// GetByID obtiene un usuario por ID
func (s *UsuarioService) GetByID(id int) (*models.Usuario, error) {
    usuario, err := s.usuarioRepo.FindByID(id)
    if err != nil {
        return nil, err
    }

    usuario.PasswordHash = ""
    return usuario, nil
}

// Update actualiza un usuario
func (s *UsuarioService) Update(id int, usuario *models.Usuario) (*models.Usuario, error) {
    // Validaciones
    if usuario.Nombre == "" || usuario.Email == "" {
        return nil, errors.New("nombre y email son requeridos")
    }

    if usuario.Rol != "instructor" && usuario.Rol != "alumno" {
        return nil, errors.New("rol inválido")
    }

    // Verificar que el usuario existe
    existing, err := s.usuarioRepo.FindByID(id)
    if err != nil {
        return nil, err
    }

    // Verificar si el email cambió y si ya existe
    if usuario.Email != existing.Email {
        emailExists, _ := s.usuarioRepo.FindByEmail(usuario.Email)
        if emailExists != nil {
            return nil, errors.New("el email ya está registrado")
        }
    }

    // Actualizar
    err = s.usuarioRepo.Update(id, usuario)
    if err != nil {
        return nil, err
    }

    usuario.PasswordHash = ""
    return usuario, nil
}

// Delete elimina un usuario
func (s *UsuarioService) Delete(id int) error {
    return s.usuarioRepo.Delete(id)
}

// ChangePassword cambia la contraseña de un usuario
func (s *UsuarioService) ChangePassword(id int, oldPassword, newPassword string) error {
    // Validaciones
    if oldPassword == "" || newPassword == "" {
        return errors.New("contraseñas requeridas")
    }

    if len(newPassword) < 6 {
        return errors.New("la nueva contraseña debe tener al menos 6 caracteres")
    }

    // Obtener usuario
    usuario, err := s.usuarioRepo.FindByID(id)
    if err != nil {
        return err
    }

    // Verificar contraseña actual
    if !utils.CheckPasswordHash(oldPassword, usuario.PasswordHash) {
        return errors.New("contraseña actual incorrecta")
    }

    // Hash de la nueva contraseña
    newHash, err := utils.HashPassword(newPassword)
    if err != nil {
        return errors.New("error al procesar la nueva contraseña")
    }

    // Actualizar
    return s.usuarioRepo.UpdatePassword(id, newHash)
}
