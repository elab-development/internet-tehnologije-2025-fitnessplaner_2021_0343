package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB() error {
	// Get database credentials from environment variables or use defaults
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "Vojislav123!")
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "app_db")

	log.Printf("🔌 Attempting to connect to MySQL server: %s@%s:%s", dbUser, dbHost, dbPort)

	// First, connect to MySQL server without specifying database
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort)

	tempDB, err := sql.Open("mysql", dsnWithoutDB)
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL server: %w\n💡 Check if MySQL is running and credentials are correct", err)
	}
	defer tempDB.Close()

	// Test connection to MySQL server
	if err := tempDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping MySQL server: %w\n💡 Possible issues:\n   - MySQL server is not running\n   - Wrong username/password\n   - Wrong host/port", err)
	}

	log.Println("✅ Connected to MySQL server")

	// Create database if it doesn't exist
	log.Printf("📦 Checking if database '%s' exists...", dbName)
	_, err = tempDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName))
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	log.Printf("✅ Database '%s' ready", dbName)

	// Now connect to the specific database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	log.Println("✅ Database connected successfully")

	// Automatically ensure all tables exist
	if err := EnsureTablesExist(); err != nil {
		return fmt.Errorf("failed to ensure tables exist: %w", err)
	}

	return nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

