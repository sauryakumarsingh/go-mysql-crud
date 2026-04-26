package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql" // registers the mysql driver
)

// Connect opens a MySQL connection pool and verifies connectivity.
// Call this once at startup and pass *sql.DB around.
func Connect() (*sql.DB, error) {
	dsn := buildDSN()

	database, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	// Connection pool tuning — important for interviews!
	database.SetMaxOpenConns(25)              // max simultaneous connections
	database.SetMaxIdleConns(10)              // kept alive when idle
	database.SetConnMaxLifetime(5 * time.Minute) // recycle after 5 min

	// Verify the connection is actually reachable
	if err := database.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	return database, nil
}

// buildDSN constructs the Data Source Name from env vars.
// Format: user:password@tcp(host:port)/dbname?parseTime=true
func buildDSN() string {
	host     := getEnv("DB_HOST", "localhost")
	port     := getEnv("DB_PORT", "3306")
	user     := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "root")
	dbname   := getEnv("DB_NAME", "productsdb")

	// parseTime=true → MySQL DATETIME columns scan directly into time.Time
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		user, password, host, port, dbname)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
