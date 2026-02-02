package utils

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// InitDB inicijalizuje konekciju sa bazom
func InitDB() error {
	dsn := "root:Vojislav123!@tcp(127.0.0.1:3306)/app_db?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connected successfully")
	return nil
}

// CloseDB zaustavlja konekciju sa bazom
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
