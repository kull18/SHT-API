package services

import (
    "cursos-api/models"
    "cursos-api/repository"
    "cursos-api/utils"
    "errors"
)

type AuthService struct {
    usuarioRepo *repository.UsuarioRepository
}

func NewAuthService() *AuthService {
    return &AuthService{
        usuarioRepo: repository.NewUsuarioRepository(),
    }
}

// Register registra un nuevo usuario
func (s *AuthService) Register(req *models.RegisterRequest) (*models.Usuario, error) {
    // Validaciones
    if req.Nombre == "" || req.Email == "" || req.Password == "" {
        return nil, errors.New("todos los campos son requeridos")
    }

    if req.Rol != "instructor" && req.Rol != "alumno" {
        return nil, errors.New("rol inválido, debe ser 'instructor' o 'alumno'")
    }

    if len(req.Password) < 6 {
        return nil, errors.New("la contraseña debe tener al menos 6 caracteres")
    }

    // Verificar si el email ya existe
    existingUser, _ := s.usuarioRepo.FindByEmail(req.Email)
    if existingUser != nil {
        return nil, errors.New("el email ya está registrado")
    }

    // Hash de la contraseña
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        return nil, errors.New("error al procesar la contraseña")
    }

    // Crear usuario
    usuario := &models.Usuario{
        Nombre:       req.Nombre,
        Email:        req.Email,
        PasswordHash: hashedPassword,
        Rol:          req.Rol,
    }

    err = s.usuarioRepo.Create(usuario)
    if err != nil {
        return nil, err
    }

    return usuario, nil
}

// Login autentica a un usuario
func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
    // Validaciones
    if req.Email == "" || req.Password == "" {
        return nil, errors.New("email y contraseña son requeridos")
    }

    // Buscar usuario
    usuario, err := s.usuarioRepo.FindByEmail(req.Email)
    if err != nil {
        return nil, errors.New("credenciales inválidas")
    }

    // Verificar contraseña
    if !utils.CheckPasswordHash(req.Password, usuario.PasswordHash) {
        return nil, errors.New("credenciales inválidas")
    }

    // Generar token
    token, err := utils.GenerateJWT(usuario.ID, usuario.Email, usuario.Rol)
    if err != nil {
        return nil, errors.New("error al generar el token")
    }

    // Limpiar el password hash antes de devolver
    usuario.PasswordHash = ""

    return &models.LoginResponse{
        Token:   token,
        Usuario: usuario,
    }, nil
}

// GetProfile obtiene el perfil de un usuario
func (s *AuthService) GetProfile(userID int) (*models.Usuario, error) {
    usuario, err := s.usuarioRepo.FindByID(userID)
    if err != nil {
        return nil, err
    }

    // Limpiar el password hash
    usuario.PasswordHash = ""

    return usuario, nil
}
