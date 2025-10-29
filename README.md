# 🔐 Auth Microservice

Microservicio de autenticación y autorización basado en JWT para aplicaciones de microservicios. Proporciona registro, login, gestión de tokens y autenticación segura.

## 🚀 Características

- ✅ **Registro y autenticación de usuarios** con bcrypt
- ✅ **JWT (JSON Web Tokens)** para access y refresh tokens
- ✅ **Cache de tokens** con Redis
- ✅ **Lista negra de tokens** (logout/revocación)
- ✅ **PostgreSQL** para almacenamiento de usuarios
- ✅ **Clean Architecture** (domain, service, repository, handler)
- ✅ **Métricas con Prometheus**
- ✅ **Health checks** (liveness y readiness)
- ✅ **Logging estructurado** con Zap
- ✅ **Dockerizado y Kubernetes-ready**
- ✅ **Kustomize** para múltiples entornos
- ✅ **External Secrets Operator** para gestión de secrets
- ✅ **Horizontal Pod Autoscaler (HPA)**
- ✅ **CI/CD con GitHub Actions**


## 🏗️ Arquitectura

```
┌─────────────────────────────────────────────────────┐
│                   API Gateway                        │
└──────────────────┬──────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────┐
│              Auth Microservice                       │
│  ┌────────────────────────────────────────────┐    │
│  │          HTTP Handlers                      │    │
│  │  (Register, Login, Refresh, Logout, Me)    │    │
│  └─────────────────┬───────────────────────────┘    │
│                    │                                 │
│  ┌─────────────────▼───────────────────────────┐    │
│  │          Service Layer                      │    │
│  │  (AuthService, JWTService)                  │    │
│  └─────────────────┬───────────────────────────┘    │
│                    │                                 │
│  ┌─────────────────▼───────────────────────────┐    │
│  │        Repository Layer                     │    │
│  │  (UserRepo, TokenRepo)                      │    │
│  └─────────┬───────────────┬───────────────────┘    │
└────────────┼───────────────┼──────────────────────────┘
             │               │
      ┌──────▼──────┐  ┌────▼──────┐
      │  PostgreSQL │  │   Redis   │
      │   (Users)   │  │ (Tokens)  │
      └─────────────┘  └───────────┘
```

### Capas de la Arquitectura

1. **Domain**: Entidades de negocio (User, Token, Errors)
2. **Repository**: Interfaces de acceso a datos
3. **Infrastructure**: Implementaciones concretas (PostgreSQL, Redis)
4. **Service**: Lógica de negocio (AuthService, JWTService)
5. **Handler**: HTTP handlers y middleware
6. **Config**: Configuración de la aplicación

## 🛠️ Tecnologías

- **Lenguaje**: Go 1.21+
- **Base de datos**: PostgreSQL 16
- **Cache**: Redis 7
- **Framework HTTP**: Gorilla Mux
- **JWT**: golang-jwt/jwt/v5
- **Password Hashing**: bcrypt
- **Logging**: Uber Zap
- **Métricas**: Prometheus
- **Containerización**: Docker
- **Orquestación**: Kubernetes + Kustomize
- **CI/CD**: GitHub Actions

## 📦 Requisitos

### Para desarrollo local:
- Go 1.21+
- Docker y Docker Compose
- Make (opcional)
- PostgreSQL 16+ (o usar docker-compose)
- Redis 7+ (o usar docker-compose)

### Para producción:
- Kubernetes cluster (EKS recomendado)
- PostgreSQL 16+ (RDS recomendado)
- Redis 7+ (ElastiCache recomendado)
- External Secrets Operator instalado
- Prometheus Operator (opcional)

## 🚀 Instalación

### 1. Clonar el repositorio

```bash
git clone https://github.com/kristianrpo/auth-microservice.git
cd auth-microservice
```

### 2. Instalar dependencias

```bash
make tidy
# o
go mod download
```

### 3. Configurar variables de entorno

```bash
cp .env.example .env
# Editar .env con tus configuraciones
```

### 4. Iniciar servicios con Docker Compose

```bash
make docker-up
```

Esto iniciará:
- PostgreSQL en `localhost:5432`
- Redis en `localhost:6379`
- Auth Service en `localhost:8080`
- Prometheus en `localhost:9090`
- Grafana en `localhost:3000`


## 📡 API Endpoints

### Base URL

```
http://localhost:8080/api/auth
```

### Autenticación

#### 1. Registro de Usuario

```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "usuario@ejemplo.com",
  "password": "Password123!",
  "name": "Usuario Ejemplo"
}
```

**Respuesta exitosa (201):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "usuario@ejemplo.com",
  "name": "Usuario Ejemplo",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

#### 2. Login

```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "usuario@ejemplo.com",
  "password": "Password123!"
}
```

**Respuesta exitosa (200):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

#### 3. Refresh Token

```http
POST /api/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 4. Logout (Requiere autenticación)

```http
POST /api/auth/logout
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 5. Obtener Usuario Actual (Requiere autenticación)

```http
GET /api/auth/me
Authorization: Bearer {access_token}
```

<!-- Health and metrics details consolidated in the 'Endpoints adicionales y notas de desarrollo' section below -->

## 🔐 Autenticación JWT

### ¿Cómo funciona?

1. **Login**: El usuario envía email + password
2. **Tokens**: El servidor genera un **access token** (15 min) y un **refresh token** (7 días)
3. **Acceso**: El cliente incluye el access token en el header: `Authorization: Bearer {token}`
4. **Renovación**: Cuando el access token expira, usa el refresh token para obtener uno nuevo
5. **Logout**: Invalida ambos tokens (lista negra en Redis)

### Estructura del JWT

```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "usuario@ejemplo.com",
  "type": "access",
  "exp": 1704123456,
  "iat": 1704122556,
  "iss": "auth-microservice"
}
```

### Usando los tokens en otros microservicios

Los otros microservicios pueden validar el JWT sin consultar este servicio:

```go
// Verificar firma con la misma clave secreta
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    return []byte(jwtSecret), nil
})
```

## 📊 Monitoreo

### Prometheus

Métricas disponibles en `http://localhost:9090`

### Grafana

Dashboards en `http://localhost:3000`
- **Usuario**: `admin`
- **Contraseña**: `admin`

## 🤝 Contribución

1. Fork el proyecto
2. Crea tu feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push al branch (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📝 Licencia

MIT License - Ver [LICENSE](LICENSE) para más detalles

## 📚 Endpoints adicionales y notas de desarrollo

### Admin (gestión de OAuth clients)

Estos endpoints están pensados para administración (service-to-service) y requieren credenciales adecuadas o token admin.

- POST /api/auth/admin/oauth-clients
  - Crea un nuevo cliente OAuth
  - Body (JSON):
    {
      "name": "My Client",
      "redirect_uris": ["https://app.example.com/callback"],
      "scopes": ["read","write"]
    }
  - Respuesta (201): información del cliente (client_id, client_secret solo al crear, scopes, active)

- GET /api/auth/admin/oauth-clients
  - Lista los OAuth clients registrados (soporta paginación)

- POST /api/auth/auth/token
  - Emite un token por client-credentials (uso administrativo)

### OAuth2 — Client Credentials

Flujo para auth máquina a máquina.

- POST /oauth2/token
  - Content-Type: application/x-www-form-urlencoded
  - Body: grant_type=client_credentials&client_id={id}&client_secret={secret}&scope={scopes}
  - Respuesta (200): access_token, token_type, expires_in

### Métricas y monitoring

- GET /api/auth/metrics
  - Endpoint compatible con Prometheus que expone métricas de requests, latencias, errores en repositorios y generación de tokens.

### Health checks (detalles)

- GET /api/auth/health
  - Comprueba liveness y readiness y los dependientes (Postgres, Redis, RabbitMQ si está configurado). Devuelve 200 si todo OK.

- GET /api/auth/health/ready
  - Readiness: pruebas más completas (ping a la DB y cache).

- GET /api/auth/health/live
  - Liveness: verificación ligera de que el proceso está arriba.

### Eventos / Webhooks (RabbitMQ)

El servicio publica eventos en RabbitMQ cuando ocurren acciones relevantes (ej. user.registered). Los consumidores pueden suscribirse al exchange/queue configurado.

Ejemplo de evento: `user.registered`

Payload:

```json
{
  "user_id": "user-123",
  "email": "user@example.com",
  "created_at": "2025-10-20T12:34:56Z"
}
```

### Variables de entorno clave

- APP_PORT: puerto donde corre el servicio (por defecto 8080)
- DATABASE_URL: URL de conexión a Postgres (ej: `postgres://user:pass@host:5432/dbname?sslmode=disable`)
- REDIS_ADDR: dirección de Redis (ej: `localhost:6379`)
- JWT_SECRET: secreto para firmar JWTs (debe ser >= 32 caracteres en prod)
- LOG_LEVEL: nivel de logging (debug, info, warn, error)

### Ejecutar tests localmente

Para ejecutar todos los tests del repositorio:

```bash
go test ./...
```

Para ejecutar sólo los tests en `internal` y generar coverage (excluye paquetes `/tests`):

```bash
go test ./internal/.../tests -coverpkg=$(go list ./internal/... | grep -v '/tests' | tr '\n' ',' | sed 's/,$//') -coverprofile=coverage.out -covermode=atomic
go tool cover -func=coverage.out | tail -1
```

### Desarrollo local rápido

1. Levantar infra mínima con Docker Compose (Postgres + Redis + RabbitMQ):

```bash
docker-compose up -d
```

2. Exportar variables de ejemplo:

```bash
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable
export REDIS_ADDR=localhost:6379
export JWT_SECRET=test-secret-key-at-least-32-chars-long
```

3. Ejecutar el servicio en modo desarrollo:

```bash
go run ./cmd/server
```

4. Revisa logs y endpoints en `http://localhost:8080/api/auth`.

