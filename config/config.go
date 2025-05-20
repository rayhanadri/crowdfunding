package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect initializes the database connection using GORM and PostgreSQL.
func Connect() {
	// Get database connection details from environment variables
	dsn := "host=" + os.Getenv("POSTGRES_HOST") +
		" user=" + os.Getenv("POSTGRES_USER") +
		" password=" + os.Getenv("POSTGRES_PASSWORD") +
		" dbname=" + os.Getenv("POSTGRES_DB") +
		" port=" + os.Getenv("POSTGRES_PORT") +
		" sslmode=require"
	if os.Getenv("POSTGRES_HOST") == "" || os.Getenv("POSTGRES_USER") == "" || os.Getenv("POSTGRES_PASSWORD") == "" || os.Getenv("POSTGRES_DB") == "" || os.Getenv("POSTGRES_PORT") == "" {
		log.Fatal("One or more required environment variables are not set")
	}

	// Connect to PostgreSQL database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	log.Println("Database connection established")
}

func LoadEnv() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	log.Println("Environment variables loaded successfully")
}
