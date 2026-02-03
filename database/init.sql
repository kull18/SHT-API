-- ============================================
-- SCHEMA: Gestión de Cursos
-- Script de inicialización de base de datos
-- ============================================

-- Eliminar tablas si existen (para desarrollo)
DROP TABLE IF EXISTS cursos CASCADE;
DROP TABLE IF EXISTS usuarios CASCADE;

-- ============================================
-- TABLA: usuarios
-- ============================================
CREATE TABLE usuarios (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    rol VARCHAR(20) NOT NULL CHECK (rol IN ('instructor', 'alumno')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- TABLA: cursos
-- ============================================
CREATE TABLE cursos (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(200) NOT NULL,
    descripcion TEXT,
    duracion_horas INTEGER NOT NULL,
    instructor_id INTEGER NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    activo BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- ÍNDICES para mejorar rendimiento
-- ============================================
CREATE INDEX idx_usuarios_email ON usuarios(email);
CREATE INDEX idx_usuarios_rol ON usuarios(rol);
CREATE INDEX idx_cursos_instructor ON cursos(instructor_id);
CREATE INDEX idx_cursos_activo ON cursos(activo);

-- ============================================
-- DATOS DE PRUEBA (opcional)
-- ============================================

-- Insertar usuarios de prueba
-- Nota: Las contraseñas deben ser hasheadas en la aplicación
-- Estos son solo ejemplos, usar el endpoint /api/auth/register en producción

INSERT INTO usuarios (nombre, email, password_hash, rol) VALUES
('Juan Pérez', 'juan.instructor@example.com', '$2a$10$ejemplo_hash_contraseña', 'instructor'),
('María García', 'maria.instructor@example.com', '$2a$10$ejemplo_hash_contraseña', 'instructor'),
('Carlos López', 'carlos.alumno@example.com', '$2a$10$ejemplo_hash_contraseña', 'alumno'),
('Ana Martínez', 'ana.alumno@example.com', '$2a$10$ejemplo_hash_contraseña', 'alumno');

-- Insertar cursos de prueba
INSERT INTO cursos (nombre, descripcion, duracion_horas, instructor_id, activo) VALUES
('Introducción a Go', 'Aprende los fundamentos del lenguaje de programación Go', 40, 1, true),
('Desarrollo Web con React', 'Construcción de aplicaciones web modernas con React', 60, 1, true),
('Bases de Datos PostgreSQL', 'Diseño y administración de bases de datos relacionales', 30, 2, true),
('Arquitectura de Software', 'Patrones y mejores prácticas en arquitectura de software', 50, 2, false);

-- ============================================
-- FUNCIÓN para actualizar updated_at automáticamente
-- ============================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers para actualizar updated_at
CREATE TRIGGER update_usuarios_updated_at
    BEFORE UPDATE ON usuarios
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_cursos_updated_at
    BEFORE UPDATE ON cursos
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- VERIFICACIÓN
-- ============================================
SELECT 'Base de datos inicializada correctamente' AS status;
