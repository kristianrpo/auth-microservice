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

## 📋 Tabla de Contenidos

- [Arquitectura](#arquitectura)
- [Tecnologías](#tecnologías)
- [Requisitos](#requisitos)
- [Instalación](#instalación)
- [Uso](#uso)
- [API Endpoints](#api-endpoints)
- [Autenticación JWT](#autenticación-jwt)
- [Deployment](#deployment)
- [Monitoreo](#monitoreo)

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

## ⚙️ Configuración

### Variables de Entorno

| Variable | Descripción | Default | Requerido |
|----------|-------------|---------|-----------|
| `SERVER_HOST` | Host del servidor | `0.0.0.0` | No |
| `SERVER_PORT` | Puerto del servidor | `8080` | No |
| `DB_HOST` | Host de PostgreSQL | `localhost` | Sí |
| `DB_PORT` | Puerto de PostgreSQL | `5432` | No |
| `DB_USER` | Usuario de PostgreSQL | `authuser` | Sí |
| `DB_PASSWORD` | Contraseña de PostgreSQL | - | **Sí** |
| `DB_NAME` | Nombre de la base de datos | `authdb` | No |
| `DB_SSL_MODE` | Modo SSL de PostgreSQL | `disable` | No |
| `REDIS_HOST` | Host de Redis | `localhost` | Sí |
| `REDIS_PORT` | Puerto de Redis | `6379` | No |
| `REDIS_PASSWORD` | Contraseña de Redis | - | No |
| `REDIS_DB` | Número de base de datos Redis | `0` | No |
| `JWT_SECRET` | Clave secreta para JWT (min 32 chars) | - | **Sí** |
| `JWT_ACCESS_TOKEN_DURATION` | Duración del access token | `15m` | No |
| `JWT_REFRESH_TOKEN_DURATION` | Duración del refresh token | `168h` | No |
| `APP_ENV` | Entorno de la aplicación | `development` | No |
| `LOG_LEVEL` | Nivel de logging | `info` | No |

## 🎯 Uso

### Desarrollo Local

```bash
# Iniciar servicios
make docker-up

# Ver logs
make docker-logs

# Ejecutar tests
make test

# Ejecutar con coverage
make test-coverage

# Linting
make lint

# Formatear código
make fmt
```

### Sin Docker

```bash
# Asegúrate de tener PostgreSQL y Redis corriendo

# Ejecutar la aplicación
make run
# o
go run ./cmd/server/main.go
```

## 📡 API Endpoints

### Base URL

```
http://localhost:8080/api/v1
```

### Autenticación

#### 1. Registro de Usuario

```http
POST /api/v1/auth/register
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
POST /api/v1/auth/login
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
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 4. Logout (Requiere autenticación)

```http
POST /api/v1/auth/logout
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 5. Obtener Usuario Actual (Requiere autenticación)

```http
GET /api/v1/auth/me
Authorization: Bearer {access_token}
```

### Health Checks

```http
GET /api/v1/health              # Health check completo
GET /api/v1/health/ready        # Readiness probe
GET /api/v1/health/live         # Liveness probe
```

### Métricas

```http
GET /api/v1/metrics             # Métricas Prometheus
```

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

## 🚢 Deployment

### Kubernetes con Kustomize

```bash
# Development
kubectl apply -k k8s/overlays/dev

# Staging
kubectl apply -k k8s/overlays/staging

# Production
kubectl apply -k k8s/overlays/production
```

Ver más detalles en [`k8s/README.md`](k8s/README.md)

### Terraform (AWS)

```bash
cd infra/terraform/aws

# Inicializar
terraform init

# Planificar
terraform plan -var-file=terraform.tfvars

# Aplicar
terraform apply -var-file=terraform.tfvars
```

Ver más detalles en [`infra/terraform/aws/README.md`](infra/terraform/aws/README.md)

## 📊 Monitoreo

### Prometheus

Métricas disponibles en `http://localhost:9090`

### Grafana

Dashboards en `http://localhost:3000`
- **Usuario**: `admin`
- **Contraseña**: `admin`

### Logs

```bash
# Ver logs locales
make docker-logs

# Ver logs en Kubernetes
kubectl logs -l app=auth-service -f
```

## 🧪 Testing

```bash
# Ejecutar todos los tests
make test

# Con coverage
make test-coverage

# Ver coverage en el navegador
open coverage.html
```

## 🔒 Seguridad

- ✅ **Passwords**: Hasheadas con bcrypt (cost 10)
- ✅ **JWT**: Firmados con HS256
- ✅ **HTTPS**: Recomendado en producción
- ✅ **Rate Limiting**: Configurado en Ingress
- ✅ **SQL Injection**: Protección con prepared statements
- ✅ **Secrets**: Gestionados con External Secrets Operator

## 🤝 Contribución

1. Fork el proyecto
2. Crea tu feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push al branch (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📝 Licencia

MIT License - Ver [LICENSE](LICENSE) para más detalles

## 👤 Autor

**Kristian Restrepo**
- GitHub: [@kristianrpo](https://github.com/kristianrpo)
- Repositorios relacionados:
  - [Infrastructure Shared](https://github.com/kristianrpo/infrastructure-shared)
  - [Documents Management](https://github.com/kristianrpo/documents-management-microservice)

## 🙏 Agradecimientos

- Clean Architecture principles
- JWT best practices
- Go microservices patterns

