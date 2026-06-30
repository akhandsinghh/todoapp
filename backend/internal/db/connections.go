package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"todo-app/backend/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

// DatabaseConnector abstracts database opening and migration execution.
type DatabaseConnector interface {
	Connect() (*sql.DB, error)
	Migrate(db *sql.DB, dir string) error
}

type SQLDatabaseConnector struct {
	cfg config.DatabaseConfig
}

// NewDatabaseConnector creates a connector from the provided database config.
func NewDatabaseConnector(cfg config.DatabaseConfig) *SQLDatabaseConnector {
	return &SQLDatabaseConnector{cfg: cfg}
}

// Connect opens the SQL database and verifies the connection.
func (c *SQLDatabaseConnector) Connect() (*sql.DB, error) {
	database, err := sql.Open("mysql", buildDSN(c.cfg))
	if err != nil {
		return nil, err
	}

	database.SetMaxOpenConns(c.cfg.MaxOpenConns())
	database.SetMaxIdleConns(c.cfg.MaxIdleConns())
	database.SetConnMaxLifetime(time.Hour)

	return database, database.Ping()
}

// Migrate runs all SQL migrations from the provided folder.
func (c *SQLDatabaseConnector) Migrate(db *sql.DB, dir string) error {
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
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("%s: %w", e.Name(), err)
		}
	}
	return nil
}

func buildDSN(cfg config.DatabaseConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", cfg.Username(), cfg.Password(), cfg.Host(), cfg.Port(), cfg.Name())
}
