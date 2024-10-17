package connection

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"sync"
)

var (
	CommandDB   *gorm.DB
	QueryDB     *sql.DB
	once        sync.Once // Ensure the connection is established only once
	isConnected bool      // To track if the connection has been established
)

func loadEnv() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func getDBConfig() (string, string, string, string, string, string, string) {
	// Define the PostgreSQL connection settings using environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	dbdriver := os.Getenv("DB_DRIVER")

	return host, port, user, password, dbname, sslmode, dbdriver
}

func connectGORM(host, port, user, password, dbname, sslmode string) (*gorm.DB, error) {
	// Construct the DSN for GORM
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		host, port, user, dbname, password, sslmode)

	// Establish the GORM connection
	ormdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database using GORM: %w", err)
	}

	log.Println("ORM connection established successfully")
	CommandDB = ormdb
	return CommandDB, nil
}

func connectSQL(host, port, user, password, dbname, sslmode, dbdriver string) (*sql.DB, error) {
	// Construct the connection string for the standard library SQL connection
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		host, port, user, dbname, password, sslmode)
	// Opening a connection to the database
	log.Println("Attempting to establish a connection with the database...")
	db, err := sql.Open(dbdriver, connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)

	}

	QueryDB = db
	log.Println("Successfully connected to the database using SQL!")
	return QueryDB, nil
}

func Connect() {
	once.Do(func() {
		loadEnv()

		host, port, user, password, dbname, sslmode, dbdriver := getDBConfig()

		var err error

		// Initialize GORM connection
		CommandDB, err = connectGORM(host, port, user, password, dbname, sslmode)
		if err != nil {
			log.Fatalf("Could not connect to the GORM database: %v", err)
		}

		// Initialize SQL connection
		QueryDB, err = connectSQL(host, port, user, password, dbname, sslmode, dbdriver)
		if err != nil {
			log.Fatalf("Could not connect to the SQL database: %v", err)
		}

		isConnected = true
	})
}

// GetCommandDB returns the GORM database connection if it is established.
// If the connection is not established, it returns an error.
func GetCommandDB() (*gorm.DB, error) {
	if !isConnected {
		return nil, fmt.Errorf("database connection is not established, call Connect() first")
	}
	return CommandDB, nil
}

// GetQueryDB returns the SQL database connection if it is established.
// If the connection is not established, it returns an error.
func GetQueryDB() (*sql.DB, error) {
	if !isConnected {
		return nil, fmt.Errorf("database connection is not established, call Connect() first")
	}
	return QueryDB, nil
}
