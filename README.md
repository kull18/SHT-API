# 📚 API de Gestión de Cursos

API REST desarrollada en Go para la gestión de cursos educativos con autenticación JWT.

## 🚀 Características

- ✅ Autenticación JWT
- ✅ Control de roles (Instructor/Alumno)
- ✅ CRUD completo de Usuarios
- ✅ CRUD completo de Cursos
- ✅ Validaciones de seguridad
- ✅ PostgreSQL como base de datos

## 📋 Requisitos Previos

- Go 1.21 o superior
- PostgreSQL 12 o superior
- Git

## 🛠️ Instalación

### 1. Clonar el repositorio

```bash
git clone <tu-repositorio>
cd cursos-api
```

### 2. Instalar dependencias

```bash
go mod download
```

### 3. Configurar la base de datos

Crear la base de datos en PostgreSQL:

```sql
CREATE DATABASE cursos_db;
```

Ejecutar el script de inicialización:

```bash
psql -U postgres -d cursos_db -f database/init.sql
```

### 4. Configurar variables de entorno

Copiar el archivo de ejemplo y configurarlo:

```bash
cp .env.example .env
```

Editar `.env` con tus credenciales:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=tu_password
DB_NAME=cursos_db
DB_SSLMODE=disable

JWT_SECRET=tu_clave_secreta_super_segura
PORT=8080
```

### 5. Ejecutar la aplicación

```bash
go run main.go
```

La API estará disponible en `http://localhost:8080`

## 📖 Endpoints de la API

### 🔐 Autenticación

#### Registro de Usuario
```http
POST /api/auth/register
Content-Type: application/json

{
  "nombre": "Juan Pérez",
  "email": "juan@example.com",
  "password": "password123",
  "rol": "instructor"  // o "alumno"
}
```

**Respuesta exitosa (201):**
```json
{
  "message": "Usuario registrado exitosamente",
  "usuario": {
    "id": 1,
    "nombre": "Juan Pérez",
    "email": "juan@example.com",
    "rol": "instructor",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "juan@example.com",
  "password": "password123"
}
```

**Respuesta exitosa (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "usuario": {
    "id": 1,
    "nombre": "Juan Pérez",
    "email": "juan@example.com",
    "rol": "instructor",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

#### Obtener Perfil
```http
GET /api/auth/profile
Authorization: Bearer {token}
```

### 👥 Usuarios

#### Listar Todos los Usuarios
```http
GET /api/usuarios
Authorization: Bearer {token}
```

#### Obtener Usuario por ID
```http
GET /api/usuarios/{id}
Authorization: Bearer {token}
```

#### Actualizar Usuario
```http
PUT /api/usuarios/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "nombre": "Juan Pérez Actualizado",
  "email": "juan.nuevo@example.com",
  "rol": "instructor"
}
```

**Nota:** Un usuario solo puede actualizar su propio perfil.

#### Eliminar Usuario
```http
DELETE /api/usuarios/{id}
Authorization: Bearer {token}
```

**Nota:** Un usuario solo puede eliminar su propio perfil.

#### Cambiar Contraseña
```http
POST /api/usuarios/change-password
Authorization: Bearer {token}
Content-Type: application/json

{
  "old_password": "password123",
  "new_password": "newpassword456"
}
```

### 📚 Cursos

#### Crear Curso (Solo Instructores)
```http
POST /api/cursos
Authorization: Bearer {token}
Content-Type: application/json

{
  "nombre": "Desarrollo Web con Go",
  "descripcion": "Aprende a crear aplicaciones web con Go",
  "duracion_horas": 40,
  "instructor_id": 1
}
```

**Respuesta exitosa (201):**
```json
{
  "message": "Curso creado exitosamente",
  "curso": {
    "id": 1,
    "nombre": "Desarrollo Web con Go",
    "descripcion": "Aprende a crear aplicaciones web con Go",
    "duracion_horas": 40,
    "instructor_id": 1,
    "activo": true,
    "created_at": "2024-01-15T11:00:00Z",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

#### Listar Cursos
```http
GET /api/cursos
Authorization: Bearer {token}
```

**Comportamiento:**
- **Instructores:** Solo ven sus propios cursos
- **Alumnos:** Ven todos los cursos activos

#### Obtener Mis Cursos (Solo Instructores)
```http
GET /api/cursos/my-cursos
Authorization: Bearer {token}
```

#### Obtener Curso por ID
```http
GET /api/cursos/{id}
Authorization: Bearer {token}
```

**Comportamiento:**
- **Instructores:** Solo pueden ver sus propios cursos
- **Alumnos:** Solo pueden ver cursos activos

#### Actualizar Curso (Solo Instructores)
```http
PUT /api/cursos/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "nombre": "Desarrollo Web Avanzado con Go",
  "descripcion": "Aprende técnicas avanzadas",
  "duracion_horas": 50,
  "instructor_id": 1,
  "activo": true
}
```

**Nota:** Un instructor solo puede actualizar sus propios cursos.

#### Eliminar Curso (Solo Instructores)
```http
DELETE /api/cursos/{id}
Authorization: Bearer {token}
```

**Nota:** Un instructor solo puede eliminar sus propios cursos.

#### Activar/Desactivar Curso (Solo Instructores)
```http
PATCH /api/cursos/{id}/toggle-activo
Authorization: Bearer {token}
```

### 🏥 Salud del Servidor

#### Health Check
```http
GET /health
```

**Respuesta:**
```json
{
  "status": "OK",
  "message": "API de Cursos funcionando correctamente"
}
```

## 🔒 Seguridad

### JWT (JSON Web Tokens)

Todos los endpoints protegidos requieren un token JWT en el header:

```
Authorization: Bearer {tu_token_jwt}
```

El token se obtiene al hacer login y expira después de 24 horas.

### Control de Roles

- **Instructor:** Puede crear, ver, editar y eliminar sus propios cursos
- **Alumno:** Puede ver cursos activos e inscribirse (futura funcionalidad)

### Validaciones

- Contraseñas hasheadas con bcrypt
- Validación de emails únicos
- Validación de roles
- Control de permisos por usuario
- Sanitización de entradas

## 🗂️ Estructura del Proyecto

```
cursos-api/
├── config/           # Configuración de BD
├── database/         # Scripts SQL
├── handlers/         # Controladores HTTP
├── middleware/       # Middlewares (Auth, CORS)
├── models/          # Modelos de datos
├── repository/      # Capa de acceso a datos
├── routes/          # Definición de rutas
├── services/        # Lógica de negocio
├── utils/           # Utilidades (JWT, Hash)
├── .env.example     # Ejemplo de variables de entorno
├── go.mod           # Dependencias
└── main.go          # Punto de entrada
```

## 🧪 Pruebas con cURL

### Registrar un instructor
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

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "juan@example.com",
    "password": "password123"
  }'
```

### Crear un curso
```bash
curl -X POST http://localhost:8080/api/cursos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {TU_TOKEN}" \
  -d '{
    "nombre": "Curso de Go",
    "descripcion": "Aprende Go desde cero",
    "duracion_horas": 40,
    "instructor_id": 1
  }'
```

### Listar cursos
```bash
curl -X GET http://localhost:8080/api/cursos \
  -H "Authorization: Bearer {TU_TOKEN}"
```

## 📝 Códigos de Estado HTTP

- `200 OK` - Solicitud exitosa
- `201 Created` - Recurso creado exitosamente
- `400 Bad Request` - Datos inválidos
- `401 Unauthorized` - No autenticado
- `403 Forbidden` - Sin permisos
- `404 Not Found` - Recurso no encontrado
- `500 Internal Server Error` - Error del servidor

## 🐛 Solución de Problemas

### Error de conexión a la base de datos
- Verificar que PostgreSQL esté corriendo
- Revisar credenciales en `.env`
- Verificar que la base de datos existe

### Token inválido o expirado
- Hacer login nuevamente para obtener un nuevo token
- Verificar que el token se envíe correctamente en el header

### Error al crear curso
- Verificar que el usuario tenga rol de instructor
- Verificar que el instructor_id corresponda al usuario autenticado

## 📄 Licencia

Este proyecto es de código abierto para propósitos educativos.

## 👨‍💻 Autor

Desarrollo para el curso de Gestión de Cursos

---

**¡Feliz codificación! 🚀**
