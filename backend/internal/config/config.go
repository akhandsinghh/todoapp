package config

import (
	"bufio"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type AppConfig struct {
	environment string
	database    *DatabaseConfig
	server      *ServerConfig
	auth        *AuthConfig
}

type DatabaseConfig struct {
	host         string
	port         string
	username     string
	password     string
	database     string
	maxOpenConns int
	maxIdleConns int
}

type ServerConfig struct {
	port       string
	host       string
	appEnv     string
	corsOrigin string
}

type AuthConfig struct {
	jwtSecret string
}

func (c *AppConfig) Database() *DatabaseConfig {
	if c == nil {
		return nil
	}
	return c.database
}

func (c *AppConfig) Server() *ServerConfig {
	if c == nil {
		return nil
	}
	return c.server
}

func (c *AppConfig) Auth() *AuthConfig {
	if c == nil {
		return nil
	}
	return c.auth
}

func (c *DatabaseConfig) Host() string      { return c.host }
func (c *DatabaseConfig) Port() string      { return c.port }
func (c *DatabaseConfig) Username() string  { return c.username }
func (c *DatabaseConfig) Password() string  { return c.password }
func (c *DatabaseConfig) Name() string      { return c.database }
func (c *DatabaseConfig) MaxOpenConns() int { return c.maxOpenConns }
func (c *DatabaseConfig) MaxIdleConns() int { return c.maxIdleConns }

func (c *ServerConfig) Port() string       { return c.port }
func (c *ServerConfig) Host() string       { return c.host }
func (c *ServerConfig) AppEnv() string     { return c.appEnv }
func (c *ServerConfig) CORSOrigin() string { return c.corsOrigin }

func (c *AuthConfig) JWTSecret() string { return c.jwtSecret }

// LoadConfig initializes and returns the centralized application configuration.
func LoadConfig() *AppConfig {
	return &AppConfig{
		environment: getEnv("APP_ENV", "production"),
		database:    NewDatabaseConfig(),
		server:      NewServerConfig(),
		auth:        NewAuthConfig(),
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

// NewDatabaseConfig creates and initializes database configuration from environment variables.
func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		host:         getEnv("DB_HOST", "127.0.0.1"),
		port:         getEnv("DB_PORT", "3306"),
		username:     getEnv("DB_USER", "root"),
		password:     getEnv("DB_PASSWORD", "password"),
		database:     getEnv("DB_NAME", "todo_app"),
		maxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 20),
		maxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 5),
	}
}

// NewServerConfig creates and initializes server configuration from environment variables.
func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		port:       getEnv("APP_PORT", "8080"),
		host:       getEnv("APP_HOST", "0.0.0.0"),
		appEnv:     getEnv("APP_ENV", "production"),
		corsOrigin: getEnv("CORS_ALLOWED_ORIGIN", "http://localhost:3000"),
	}
}

// NewAuthConfig creates and initializes authentication configuration from environment variables.
func NewAuthConfig() *AuthConfig {
	return &AuthConfig{
		jwtSecret: getEnv("JWT_SECRET", "change-this-secret"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.Atoi(value)
		if err == nil {
			return parsed
		}
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
