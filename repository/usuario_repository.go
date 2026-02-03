package repository

import (
    "cursos-api/config"
    "cursos-api/models"
    "database/sql"
    "errors"
    "time"
)

type UsuarioRepository struct{}

func NewUsuarioRepository() *UsuarioRepository {
    return &UsuarioRepository{}
}

// Create crea un nuevo usuario
func (r *UsuarioRepository) Create(usuario *models.Usuario) error {
    query := `
        INSERT INTO usuarios (nombre, email, password_hash, rol, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at
    `
    
    now := time.Now()
    err := config.DB.QueryRow(
        query,
        usuario.Nombre,
        usuario.Email,
        usuario.PasswordHash,
        usuario.Rol,
        now,
        now,
    ).Scan(&usuario.ID, &usuario.CreatedAt, &usuario.UpdatedAt)

    return err
}

// FindByEmail busca un usuario por email
func (r *UsuarioRepository) FindByEmail(email string) (*models.Usuario, error) {
    query := `
        SELECT id, nombre, email, password_hash, rol, created_at, updated_at
        FROM usuarios
        WHERE email = $1
    `
    
    usuario := &models.Usuario{}
    err := config.DB.QueryRow(query, email).Scan(
        &usuario.ID,
        &usuario.Nombre,
        &usuario.Email,
        &usuario.PasswordHash,
        &usuario.Rol,
        &usuario.CreatedAt,
        &usuario.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, errors.New("usuario no encontrado")
    }

    return usuario, err
}

// FindByID busca un usuario por ID
func (r *UsuarioRepository) FindByID(id int) (*models.Usuario, error) {
    query := `
        SELECT id, nombre, email, password_hash, rol, created_at, updated_at
        FROM usuarios
        WHERE id = $1
    `
    
    usuario := &models.Usuario{}
    err := config.DB.QueryRow(query, id).Scan(
        &usuario.ID,
        &usuario.Nombre,
        &usuario.Email,
        &usuario.PasswordHash,
        &usuario.Rol,
        &usuario.CreatedAt,
        &usuario.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, errors.New("usuario no encontrado")
    }

    return usuario, err
}

// GetAll obtiene todos los usuarios
func (r *UsuarioRepository) GetAll() ([]models.Usuario, error) {
    query := `
        SELECT id, nombre, email, rol, created_at, updated_at
        FROM usuarios
        ORDER BY created_at DESC
    `
    
    rows, err := config.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var usuarios []models.Usuario
    for rows.Next() {
        var usuario models.Usuario
        err := rows.Scan(
            &usuario.ID,
            &usuario.Nombre,
            &usuario.Email,
            &usuario.Rol,
            &usuario.CreatedAt,
            &usuario.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        usuarios = append(usuarios, usuario)
    }

    return usuarios, nil
}

// Update actualiza un usuario
func (r *UsuarioRepository) Update(id int, usuario *models.Usuario) error {
    query := `
        UPDATE usuarios
        SET nombre = $1, email = $2, rol = $3, updated_at = $4
        WHERE id = $5
        RETURNING updated_at
    `
    
    now := time.Now()
    err := config.DB.QueryRow(
        query,
        usuario.Nombre,
        usuario.Email,
        usuario.Rol,
        now,
        id,
    ).Scan(&usuario.UpdatedAt)

    if err == sql.ErrNoRows {
        return errors.New("usuario no encontrado")
    }

    usuario.ID = id
    return err
}

// Delete elimina un usuario
func (r *UsuarioRepository) Delete(id int) error {
    query := `DELETE FROM usuarios WHERE id = $1`
    
    result, err := config.DB.Exec(query, id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return errors.New("usuario no encontrado")
    }

    return nil
}

// UpdatePassword actualiza la contraseña de un usuario
func (r *UsuarioRepository) UpdatePassword(id int, newPasswordHash string) error {
    query := `
        UPDATE usuarios
        SET password_hash = $1, updated_at = $2
        WHERE id = $3
    `
    
    result, err := config.DB.Exec(query, newPasswordHash, time.Now(), id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return errors.New("usuario no encontrado")
    }

    return nil
}
