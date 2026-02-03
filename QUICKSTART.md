# 🚀 Guía Rápida de Inicio

## Pasos para ejecutar la API

### 1. Preparar Base de Datos

```bash
# Crear la base de datos
createdb cursos_db

# O usando psql
psql -U postgres
CREATE DATABASE cursos_db;
\q

# Ejecutar el script de inicialización
psql -U postgres -d cursos_db -f database/init.sql
```

### 2. Configurar Variables de Entorno

```bash
# Copiar el archivo de ejemplo
cp .env.example .env

# Editar .env con tus credenciales
nano .env
```

### 3. Instalar Dependencias

```bash
go mod download
```

### 4. Ejecutar la API

```bash
go run main.go
```

La API estará disponible en: `http://localhost:8080`

## 🧪 Probar la API

### Usando cURL

#### 1. Registrar un instructor
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Juan Instructor",
    "email": "juan@example.com",
    "password": "password123",
    "rol": "instructor"
  }'
```

#### 2. Login (guarda el token)
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "juan@example.com",
    "password": "password123"
  }'
```

Copia el token de la respuesta.

#### 3. Crear un curso
```bash
curl -X POST http://localhost:8080/api/cursos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TU_TOKEN_AQUI" \
  -d '{
    "nombre": "Introducción a Go",
    "descripcion": "Aprende Go desde cero",
    "duracion_horas": 40,
    "instructor_id": 1
  }'
```

#### 4. Listar cursos
```bash
curl -X GET http://localhost:8080/api/cursos \
  -H "Authorization: Bearer TU_TOKEN_AQUI"
```

### Usando Postman

1. Importa el archivo `postman_collection.json`
2. Ejecuta "Login" para obtener el token automáticamente
3. El token se guardará en las variables de colección
4. Ejecuta los demás endpoints

## 📝 Flujo Típico de Uso

### Para Instructores:

1. **Registrarse** → `POST /api/auth/register` (rol: instructor)
2. **Login** → `POST /api/auth/login`
3. **Crear cursos** → `POST /api/cursos`
4. **Ver mis cursos** → `GET /api/cursos/my-cursos`
5. **Editar curso** → `PUT /api/cursos/{id}`
6. **Activar/Desactivar** → `PATCH /api/cursos/{id}/toggle-activo`

### Para Alumnos:

1. **Registrarse** → `POST /api/auth/register` (rol: alumno)
2. **Login** → `POST /api/auth/login`
3. **Ver cursos disponibles** → `GET /api/cursos`
4. **Ver detalles de un curso** → `GET /api/cursos/{id}`

## ⚠️ Problemas Comunes

### Error de conexión a la base de datos
```
Error al conectar a la base de datos
```
**Solución:** Verifica que PostgreSQL esté corriendo y las credenciales en `.env` sean correctas.

### Token inválido
```
{"error":"Token inválido o expirado"}
```
**Solución:** Haz login nuevamente para obtener un nuevo token.

### No puedo crear cursos
```
{"error":"solo los instructores pueden crear cursos"}
```
**Solución:** Asegúrate de estar registrado como instructor.

## 📊 Estructura de la Base de Datos

```
usuarios
├── id (PK)
├── nombre
├── email (UNIQUE)
├── password_hash
├── rol (instructor/alumno)
├── created_at
└── updated_at

cursos
├── id (PK)
├── nombre
├── descripcion
├── duracion_horas
├── instructor_id (FK → usuarios)
├── activo
├── created_at
└── updated_at
```

## 🔒 Seguridad

- Todas las contraseñas se hashean con bcrypt
- Los tokens JWT expiran en 24 horas
- Los instructores solo pueden modificar sus propios cursos
- Los alumnos solo pueden ver cursos activos

## 📚 Endpoints Principales

| Método | Endpoint | Descripción | Requiere Auth | Rol |
|--------|----------|-------------|---------------|-----|
| POST | `/api/auth/register` | Registrar usuario | No | - |
| POST | `/api/auth/login` | Iniciar sesión | No | - |
| GET | `/api/auth/profile` | Ver perfil | Sí | Todos |
| GET | `/api/cursos` | Listar cursos | Sí | Todos |
| POST | `/api/cursos` | Crear curso | Sí | Instructor |
| PUT | `/api/cursos/{id}` | Actualizar curso | Sí | Instructor |
| DELETE | `/api/cursos/{id}` | Eliminar curso | Sí | Instructor |
| PATCH | `/api/cursos/{id}/toggle-activo` | Activar/Desactivar | Sí | Instructor |

## 💡 Tips

- Usa variables de entorno para el token en Postman
- Guarda el token después del login
- Los instructores solo ven/editan sus propios cursos
- Los alumnos solo ven cursos activos
- Cambia el `JWT_SECRET` en producción

## 🎯 Próximos Pasos

Una vez que domines estos endpoints básicos, puedes:
- Agregar inscripciones de alumnos
- Implementar lecciones por curso
- Agregar sistema de evaluaciones
- Implementar progreso del alumno

---

¡Ya estás listo para usar la API! 🎉
