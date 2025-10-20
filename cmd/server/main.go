package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	httpAdapter "github.com/kristianrpo/auth-microservice/internal/adapters/http"
	"github.com/kristianrpo/auth-microservice/internal/application/services"
	"github.com/kristianrpo/auth-microservice/internal/infrastructure/config"
	"github.com/kristianrpo/auth-microservice/internal/infrastructure/postgres"
	"github.com/kristianrpo/auth-microservice/internal/infrastructure/redis"

	_ "github.com/kristianrpo/auth-microservice/docs" // Swagger docs
)

// @title Auth Microservice API
// @version 1.0
// @description Authentication and user management service with JWT.
// @description This microservice provides endpoints for user registration, login, token refresh, and user management.

// @contact.name Kristian Rodriguez
// @contact.url https://github.com/kristianrpo
// @contact.email kristianrpo@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer <token>" in the field. Example: Bearer eyJhbGciOi...

// @tag.name Authentication
// @tag.description Endpoints related to authentication and authorization

// @tag.name OAuth2
// @tag.description OAuth2 endpoints for service-to-service communication

// @tag.name Admin - OAuth Clients
// @tag.description Admin endpoints for managing OAuth2 clients (requires ADMIN role)

// @tag.name Health
// @tag.description Endpoints for checking the service status

func main() {
	// Inicializar logger
	logger, err := initLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		_ = logger.Sync()
	}()

	logger.Info("Starting auth-microservice")

	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	logger.Info("Configuration loaded successfully",
		zap.String("environment", cfg.App.Environment),
		zap.String("server_address", cfg.ServerAddress()),
	)

	// Inicializar base de datos
	db, err := postgres.NewDB(cfg.DatabaseConnectionString(), logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("Failed to close database connection", zap.Error(err))
		}
	}()

	// Inicializar esquema de base de datos
	if err := postgres.InitSchema(db); err != nil {
		logger.Fatal("Failed to initialize database schema", zap.Error(err))
	}
	logger.Info("Database schema initialized")

	// Inicializar Redis
	redisClient, err := redis.NewRedisClient(cfg.RedisAddress(), cfg.Redis.Password, cfg.Redis.DB, logger)
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.Error("Failed to close Redis connection", zap.Error(err))
		}
	}()

	// Inicializar repositorios
	userRepo := postgres.NewUserRepository(db, logger)
	tokenRepo := redis.NewTokenRepository(redisClient, logger)
	oauthClientRepo := postgres.NewOAuthClientRepository(db, logger)

	// Inicializar servicios
	jwtService := services.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenDuration,
		cfg.JWT.RefreshTokenDuration,
		logger,
	)

	authService := services.NewAuthService(userRepo, tokenRepo, jwtService, logger)

	oauth2Service := services.NewOAuth2Service(
		oauthClientRepo,
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenDuration,
		logger,
	)

	// Inicializar router
	router := httpAdapter.NewRouter(authService, oauth2Service, db, redisClient, logger)

	// Configurar servidor HTTP
	server := &http.Server{
		Addr:         cfg.ServerAddress(),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Canal para errores del servidor
	serverErrors := make(chan error, 1)

	// Iniciar servidor en una goroutine
	go func() {
		logger.Info("Server starting", zap.String("address", server.Addr))
		serverErrors <- server.ListenAndServe()
	}()

	// Canal para señales de sistema
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Esperar por shutdown signal o error del servidor
	select {
	case err := <-serverErrors:
		logger.Fatal("Server error", zap.Error(err))

	case sig := <-shutdown:
		logger.Info("Shutdown signal received", zap.String("signal", sig.String()))

		// Crear contexto con timeout para el shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Intentar graceful shutdown
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("Server shutdown error", zap.Error(err))
			if err := server.Close(); err != nil {
				logger.Fatal("Server close error", zap.Error(err))
			}
		}

		logger.Info("Server stopped gracefully")
	}
}

// initLogger inicializa el logger de Zap
func initLogger() (*zap.Logger, error) {
	env := os.Getenv("APP_ENV")

	var logger *zap.Logger
	var err error

	if env == "production" {
		// Logger de producción (JSON format)
		logger, err = zap.NewProduction()
	} else {
		// Logger de desarrollo (human-friendly)
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, err
	}

	return logger, nil
}
