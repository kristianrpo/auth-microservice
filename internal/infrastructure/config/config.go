package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config contains all the application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
    RabbitMQ RabbitMQConfig
	App      AppConfig
}

// ServerConfig contains the HTTP server configuration
type ServerConfig struct {
	Host string
	Port int
}

// DatabaseConfig contains the PostgreSQL database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// RedisConfig contains the Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// JWTConfig contains the JWT configuration
type JWTConfig struct {
	Secret               string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

// RabbitMQConfig holds RabbitMQ configuration
type RabbitMQConfig struct {
    URL string

    // Consumer queue configuration
    ConsumerQueue string

    // Queue settings
    Durable       bool
    PrefetchCount int
    AutoAck       bool
}

// AppConfig contains the general application configuration
type AppConfig struct {
	Environment string
	LogLevel    string
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env if it exists (useful for local development)
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvAsInt("SERVER_PORT", 8080),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "authuser"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "authdb"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:               getEnv("JWT_SECRET", ""),
			AccessTokenDuration:  getEnvAsDuration("JWT_ACCESS_TOKEN_DURATION", 15*time.Minute),
			RefreshTokenDuration: getEnvAsDuration("JWT_REFRESH_TOKEN_DURATION", 7*24*time.Hour),
		},
    RabbitMQ: RabbitMQConfig{
        URL:                        getEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"),
        ConsumerQueue:              getEnv("RABBITMQ_CONSUMER_QUEUE", "auth_user_transferred"),
        Durable:                    true,
        PrefetchCount:              getEnvAsInt("RABBITMQ_PREFETCH_COUNT", 1),
        AutoAck:                    getEnv("RABBITMQ_AUTO_ACK", "false") == "true",
    },
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate validates that the configuration is correct
func (c *Config) Validate() error {
	if c.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if len(c.JWT.Secret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}
	return nil
}

// DatabaseConnectionString returns the connection string for PostgreSQL
func (c *Config) DatabaseConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// RedisAddress returns the Redis address
func (c *Config) RedisAddress() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

// ServerAddress returns the server address
func (c *Config) ServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// IsProd returns true if the environment is production
func (c *Config) IsProd() bool {
	return c.App.Environment == "production"
}

// Helper functions
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}