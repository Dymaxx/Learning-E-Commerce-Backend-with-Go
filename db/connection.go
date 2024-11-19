package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Import the Postgres driver
)

// DB struct to hold the SQL connection
type DB struct {
	Conn *sqlx.DB // Holds the database connection
}

// InitializeDB initializes and returns a new database connection
func InitializeDB() *DB {
	// Load environment variables
	host := getEnvOrFail("DB_HOST")
	port := getEnvOrFail("DB_PORT")
	user := getEnvOrFail("DB_USER")
	password := getEnvOrFail("DB_PASSWORD")
	dbname := getEnvOrFail("DB_NAME")

	// Create a new DB instance
	db, err := NewDB(host, port, user, password, dbname)
	if err != nil {
		log.Fatalf("Could not initialize the database: %v", err)
	}
	return db
}

// NewDB creates and returns a new DB instance
func NewDB(host, port, user, password, dbname string) (*DB, error) {
	// Build the connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	// Connect to the database
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database is unreachable: %w", err)
	}
	log.Println("Database connection established")
	return &DB{Conn: db}, nil
}

// getEnvOrFail fetches an environment variable or logs a fatal error if it's missing
func getEnvOrFail(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}
