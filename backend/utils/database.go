package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// InitDB inicijalizuje konekciju sa bazom podataka
func InitDB() error {
	// Uzimanje kredencijala iz environment promenljivih ili korišćenje podrazumevanih vrednosti
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "Vojislav123!")
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "app_db")

	log.Printf("🔌 Pokušaj konekcije na MySQL server: %s@%s:%s", dbUser, dbHost, dbPort)

	// Prvo, konekcija na MySQL server bez specificiranja baze
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort)

	tempDB, err := sql.Open("mysql", dsnWithoutDB)
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL server: %w\n💡 Check if MySQL is running and credentials are correct", err)
	}
	defer tempDB.Close()

	// Testiranje konekcije na MySQL server
	if err := tempDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping MySQL server: %w\n💡 Possible issues:\n   - MySQL server is not running\n   - Wrong username/password\n   - Wrong host/port", err)
	}

	log.Println("✅ Konektovano na MySQL server")

	// Kreiranje baze ako ne postoji
	log.Printf("📦 Provera da li baza '%s' postoji...", dbName)
	_, err = tempDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName))
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	log.Printf("✅ Baza '%s' spremna", dbName)

	// Sada konekcija na specifičnu bazu
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Testiranje konekcije
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Podešavanje connection pool-a
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	log.Println("✅ Baza podataka uspešno konektovana")

	// Automatsko osiguravanje da sve tabele postoje (fallback pristup)
	if err := EnsureTablesExist(); err != nil {
		log.Printf("⚠️  EnsureTablesExist failed: %v", err)
		log.Println("🔄 Pokušavanje sa RunMigrations...")
		// Pokušaj sa migracijama ako EnsureTablesExist ne uspe
		if err := RunMigrations(); err != nil {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
	} else {
		// Ako EnsureTablesExist uspe, pokreni migracije za dodatne izmene
		if err := RunMigrations(); err != nil {
			log.Printf("⚠️  RunMigrations failed: %v (nastavlja se sa postojećim tabelama)", err)
		}
	}

	return nil
}

// getEnv uzima environment promenljivu ili vraća podrazumevanu vrednost
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// CloseDB zatvara konekciju sa bazom podataka
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
