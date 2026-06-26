package config

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// AppConfig is the central configuration structure for the application.
type AppConfig struct {
	Environment string
	Database    *DatabaseConfig
	Server      *ServerConfig
	Auth        *AuthConfig
}

// DatabaseConfig holds database connection configuration.
type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

// ServerConfig holds server-related configuration.
type ServerConfig struct {
	Port       string // Server port
	Host       string // Server host
	AppEnv     string // Application environment (development, production, etc.)
	CORSOrigin string // CORS allowed origin
}

// AuthConfig holds authentication-related configuration.
type AuthConfig struct {
	JWTSecret string // Secret key for JWT token signing and verification
}

// LoadConfig initializes and returns the centralized application configuration.
func LoadConfig() *AppConfig {
	return &AppConfig{
		Environment: getEnv("APP_ENV", "production"),
		Database:    NewDatabaseConfig(),
		Server:      NewServerConfig(),
		Auth:        NewAuthConfig(),
	}
}

// LoadEnv loads environment variables from a .env file.
func LoadEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if os.Getenv(parts[0]) == "" {
			_ = os.Setenv(parts[0], parts[1])
		}
	}
	if err := scanner.Err(); err != nil {
		return
	}
}

// Connect establishes a database connection using the provided config.
func Connect(cfg DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	database.SetMaxOpenConns(20)
	database.SetMaxIdleConns(5)
	database.SetConnMaxLifetime(time.Hour)
	return database, database.Ping()
}

// RunMigrations executes all SQL migration files in the specified directory.
func RunMigrations(database *sql.DB, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		content, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return err
		}
		if _, err := database.Exec(string(content)); err != nil {
			return fmt.Errorf("%s: %w", e.Name(), err)
		}
	}
	return nil
}

// NewDatabaseConfig creates and initializes database configuration from environment variables.
func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("DB_HOST", "127.0.0.1"),
		Port:     getEnv("DB_PORT", "3306"),
		Username: getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "password"),
		Database: getEnv("DB_NAME", "todo_app"),
	}
}

// NewServerConfig creates and initializes server configuration from environment variables.
func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:       getEnv("APP_PORT", "8080"),
		Host:       getEnv("APP_HOST", "0.0.0.0"),
		AppEnv:     getEnv("APP_ENV", "production"),
		CORSOrigin: getEnv("CORS_ALLOWED_ORIGIN", "http://localhost:3000"),
	}
}

// NewAuthConfig creates and initializes authentication configuration from environment variables.
func NewAuthConfig() *AuthConfig {
	return &AuthConfig{
		JWTSecret: getEnv("JWT_SECRET", "change-this-secret"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// BackendRoot returns the root directory of the backend project.
func BackendRoot() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "."
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}
