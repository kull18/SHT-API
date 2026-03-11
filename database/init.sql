-- ============================================
-- SCHEMA: Gestión de Cursos
-- Script de inicialización de base de datos
-- ============================================

-- Eliminar tablas si existen (para desarrollo)
DROP TABLE IF EXISTS progreso_lecciones CASCADE;
DROP TABLE IF EXISTS inscripciones CASCADE;
DROP TABLE IF EXISTS lecciones CASCADE;
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
-- TABLA: lecciones
-- ============================================
CREATE TABLE lecciones (
    id SERIAL PRIMARY KEY,
    curso_id INTEGER NOT NULL REFERENCES cursos(id) ON DELETE CASCADE,
    titulo VARCHAR(200) NOT NULL,
    contenido TEXT,
    imagen_url VARCHAR(500),
    orden INTEGER NOT NULL DEFAULT 1,
    duracion_minutos INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- TABLA: inscripciones
-- ============================================
CREATE TABLE inscripciones (
    id SERIAL PRIMARY KEY,
    usuario_id INTEGER NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    curso_id INTEGER NOT NULL REFERENCES cursos(id) ON DELETE CASCADE,
    fecha_inscripcion TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    estado VARCHAR(20) NOT NULL DEFAULT 'activo' CHECK (estado IN ('activo', 'completado', 'cancelado')),
    progreso_porcentaje DECIMAL(5,2) DEFAULT 0.00,
    UNIQUE(usuario_id, curso_id)
);

-- ============================================
-- TABLA: progreso_lecciones
-- ============================================
CREATE TABLE progreso_lecciones (
    id SERIAL PRIMARY KEY,
    inscripcion_id INTEGER NOT NULL REFERENCES inscripciones(id) ON DELETE CASCADE,
    leccion_id INTEGER NOT NULL REFERENCES lecciones(id) ON DELETE CASCADE,
    completada BOOLEAN DEFAULT false,
    fecha_completado TIMESTAMP,
    UNIQUE(inscripcion_id, leccion_id)
);

-- ============================================
-- ÍNDICES para mejorar rendimiento
-- ============================================
CREATE INDEX idx_usuarios_email ON usuarios(email);
CREATE INDEX idx_usuarios_rol ON usuarios(rol);
CREATE INDEX idx_cursos_instructor ON cursos(instructor_id);
CREATE INDEX idx_cursos_activo ON cursos(activo);
CREATE INDEX idx_lecciones_curso ON lecciones(curso_id);
CREATE INDEX idx_lecciones_orden ON lecciones(curso_id, orden);
CREATE INDEX idx_inscripciones_usuario ON inscripciones(usuario_id);
CREATE INDEX idx_inscripciones_curso ON inscripciones(curso_id);
CREATE INDEX idx_progreso_inscripcion ON progreso_lecciones(inscripcion_id);

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

CREATE TRIGGER update_lecciones_updated_at
    BEFORE UPDATE ON lecciones
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- VERIFICACIÓN
-- ============================================
SELECT 'Base de datos inicializada correctamente' AS status;