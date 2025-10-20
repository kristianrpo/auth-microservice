package http

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/handler"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/handler/shared"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/middleware"
	"github.com/kristianrpo/auth-microservice/internal/application/services"
)

const version = "1.0.0"

// NewRouter creates and configures the main router
func NewRouter(
	authService *services.AuthService,
	oauth2Service *services.OAuth2Service,
	db *sql.DB,
	redisClient *redis.Client,
	logger *zap.Logger,
) *mux.Router {
	router := mux.NewRouter()

	// Handlers
	authHandler := shared.NewAuthHandler(authService, logger)
	oauth2Handler := shared.NewOAuth2Handler(oauth2Service, logger)
	adminOAuthHandler := shared.NewAdminOAuthClientsHandler(oauth2Service, logger)
	healthHandler := handler.NewHealthHandler(db, redisClient, logger, version)
	
	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(authService, logger)
	roleMiddleware := middleware.NewRoleMiddleware(logger)

	// Global middleware
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.RecoveryMiddleware(logger))

	// API v1 routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Public routes - Authentication routes
	api.HandleFunc("/auth/register", handler.Register(authHandler)).Methods(http.MethodPost)
	api.HandleFunc("/auth/login", handler.Login(authHandler)).Methods(http.MethodPost)
	api.HandleFunc("/auth/refresh", handler.Refresh(authHandler)).Methods(http.MethodPost)
	
	// OAuth2 Client Credentials endpoint
	api.HandleFunc("/auth/token", handler.Token(oauth2Handler)).Methods(http.MethodPost)

	// Protected routes - Authentication required routes
	protected := api.PathPrefix("/auth").Subrouter()
	protected.Use(authMiddleware.Authenticate)
	protected.HandleFunc("/logout", handler.Logout(authHandler)).Methods(http.MethodPost)
	protected.HandleFunc("/me", handler.GetMe(authHandler)).Methods(http.MethodGet)

	// Health checks
	api.HandleFunc("/health", healthHandler.Health).Methods(http.MethodGet)
	api.HandleFunc("/health/ready", healthHandler.Ready).Methods(http.MethodGet)
	api.HandleFunc("/health/live", healthHandler.Live).Methods(http.MethodGet)

	// Metrics (Prometheus)
	api.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	// Admin routes (require ADMIN role)
	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(authMiddleware.Authenticate)
	admin.Use(roleMiddleware.RequireAdmin)
	admin.HandleFunc("/oauth-clients", handler.CreateOAuthClient(adminOAuthHandler)).Methods(http.MethodPost)
	admin.HandleFunc("/oauth-clients", handler.ListOAuthClients(adminOAuthHandler)).Methods(http.MethodGet)

	// Swagger documentation route
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	// Root endpoint route
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"service": "auth-microservice",
			"version": version,
			"status":  "running",
		})
	}).Methods(http.MethodGet)

	return router
}


