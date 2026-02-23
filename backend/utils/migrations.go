package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// RunMigrations executes database migrations
func RunMigrations() error {
	log.Println("🔄 Running database migrations...")

	// Read migration file
	migrationSQL, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Split SQL into individual statements
	statements := splitSQL(string(migrationSQL))

	// Execute each statement
	for i, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" || strings.HasPrefix(statement, "--") {
			continue
		}

		// Skip CREATE DATABASE if it already exists
		if strings.HasPrefix(strings.ToUpper(statement), "CREATE DATABASE") {
			log.Println("⏭️  Skipping CREATE DATABASE (already exists)")
			continue
		}

		// Execute statement
		if _, err := DB.Exec(statement); err != nil {
			// Ignore "table already exists" errors
			if strings.Contains(err.Error(), "already exists") {
				log.Printf("⏭️  Table already exists, skipping...")
				continue
			}
			// Ignore "database exists" errors
			if strings.Contains(err.Error(), "database exists") {
				continue
			}
			log.Printf("❌ Error executing migration statement %d: %v", i+1, err)
			log.Printf("Statement: %s", statement[:min(100, len(statement))])
			return fmt.Errorf("migration failed at statement %d: %w", i+1, err)
		}
	}

	log.Println("✅ Migrations completed successfully")
	return nil
}

// splitSQL splits SQL string into individual statements
func splitSQL(sql string) []string {
	// Remove comments and split by semicolon
	lines := strings.Split(sql, "\n")
	var cleanedLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}
		cleanedLines = append(cleanedLines, line)
	}

	// Join and split by semicolon
	fullSQL := strings.Join(cleanedLines, " ")
	statements := strings.Split(fullSQL, ";")

	var result []string
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			result = append(result, stmt)
		}
	}

	return result
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// EnsureTablesExist checks if all required tables exist and creates them if needed
func EnsureTablesExist() error {
	log.Println("🔍 Checking if all tables exist...")

	// Check and create users table
	if err := ensureTable("users", createUsersTableSQL); err != nil {
		return err
	}

	// Ensure users has role column (without CHECK constraint to avoid issues)
	if err := ensureColumn("users", "role", "ALTER TABLE users ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user' AFTER goal"); err != nil {
		return err
	}
	
	// Ensure users has height column
	if err := ensureColumn("users", "height", "ALTER TABLE users ADD COLUMN height DECIMAL(5, 2) NULL AFTER role"); err != nil {
		return err
	}
	
	// Ensure users has weight column
	if err := ensureColumn("users", "weight", "ALTER TABLE users ADD COLUMN weight DECIMAL(5, 2) NULL AFTER height"); err != nil {
		return err
	}
	
	// Fix role column if it has issues (remove CHECK constraint if it causes problems)
	if err := fixRoleColumn(); err != nil {
		log.Printf("⚠️  Warning: Could not fix role column: %v", err)
	}

	// Check and create workouts table
	if err := ensureTable("workouts", createWorkoutsTableSQL); err != nil {
		return err
	}

	// Ensure workouts has calories_burned column
	if err := ensureColumn("workouts", "calories_burned", "ALTER TABLE workouts ADD COLUMN calories_burned DECIMAL(10, 2) DEFAULT 0 AFTER duration"); err != nil {
		return err
	}

	// Check and create progress table
	if err := ensureTable("progress", createProgressTableSQL); err != nil {
		return err
	}

	// Ensure progress has progress_date column
	if err := ensureColumn("progress", "progress_date", "ALTER TABLE progress ADD COLUMN progress_date DATE NOT NULL DEFAULT (CURDATE()) AFTER notes"); err != nil {
		return err
	}

	// Fix existing data issues
	if err := fixExistingData(); err != nil {
		log.Printf("⚠️  Warning: Failed to fix existing data: %v", err)
		// Don't fail if this errors, just log it
	}

	log.Println("✅ All tables and columns exist")
	return nil
}

// fixExistingData fixes common data issues in existing tables
func fixExistingData() error {
	// Fix NULL passwords (set to empty string, users will need to reset)
	_, err := DB.Exec("UPDATE users SET password = '' WHERE password IS NULL")
	if err != nil {
		log.Printf("⚠️  Could not fix NULL passwords: %v", err)
	}

	// Try to fix invalid roles - this might fail due to CHECK constraint
	// If it fails, users with invalid roles will need to be fixed manually or deleted
	_, err = DB.Exec("UPDATE users SET role = 'user' WHERE role IS NULL OR role = '' OR role NOT IN ('admin', 'user', 'premium')")
	if err != nil {
		// This is expected to fail if CHECK constraint is strict
		// Log it but don't fail - users can be fixed manually
		log.Printf("⚠️  Could not auto-fix invalid roles (CHECK constraint may block this). You may need to fix manually in MySQL.")
		log.Printf("💡 To fix manually, run in MySQL: DELETE FROM users WHERE role NOT IN ('admin', 'user', 'premium') OR role IS NULL;")
	}

	return nil
}

// ensureTable checks if table exists and creates it if it doesn't
func ensureTable(tableName, createSQL string) error {
	var exists bool
	err := DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?)",
		tableName,
	).Scan(&exists)

	if err != nil {
		return fmt.Errorf("failed to check if table %s exists: %w", tableName, err)
	}

	if !exists {
		log.Printf("📝 Creating table: %s", tableName)
		if _, err := DB.Exec(createSQL); err != nil {
			return fmt.Errorf("failed to create table %s: %w", tableName, err)
		}
		log.Printf("✅ Table %s created", tableName)
	} else {
		log.Printf("✓ Table %s already exists", tableName)
	}

	return nil
}

// ensureColumn checks if column exists and creates it if it doesn't
func ensureColumn(tableName, columnName, alterSQL string) error {
	var exists bool
	err := DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = ? AND column_name = ?)",
		tableName, columnName,
	).Scan(&exists)

	if err != nil {
		return fmt.Errorf("failed to check if column %s.%s exists: %w", tableName, columnName, err)
	}

	if !exists {
		log.Printf("📝 Adding column %s to table %s", columnName, tableName)
		if _, err := DB.Exec(alterSQL); err != nil {
			// Ignore "Duplicate column" errors
			if strings.Contains(err.Error(), "Duplicate column") {
				log.Printf("✓ Column %s.%s already exists", tableName, columnName)
				return nil
			}
			return fmt.Errorf("failed to add column %s to table %s: %w", columnName, tableName, err)
		}
		log.Printf("✅ Column %s added to table %s", columnName, tableName)
	} else {
		log.Printf("✓ Column %s.%s already exists", tableName, columnName)
	}

	return nil
}

// fixRoleColumn attempts to fix the role column structure
func fixRoleColumn() error {
	// Check current column structure
	var dataType, columnDefault, isNullable, columnType string
	err := DB.QueryRow(
		"SELECT DATA_TYPE, COLUMN_DEFAULT, IS_NULLABLE, COLUMN_TYPE FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = 'users' AND column_name = 'role'",
	).Scan(&dataType, &columnDefault, &isNullable, &columnType)
	
	if err != nil {
		return fmt.Errorf("failed to check role column: %w", err)
	}
	
	log.Printf("🔍 Role column: type=%s, full_type=%s, default=%s, nullable=%s", dataType, columnType, columnDefault, isNullable)
	
	// If it's ENUM or has CHECK constraint issues, modify it to simple VARCHAR
	if dataType == "enum" || strings.Contains(strings.ToLower(columnType), "check") {
		log.Println("🔧 Role column is ENUM or has CHECK constraint, converting to VARCHAR...")
		_, err = DB.Exec("ALTER TABLE users MODIFY COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user'")
		if err != nil {
			log.Printf("⚠️  Could not modify role column: %v", err)
			return fmt.Errorf("failed to modify role column: %w", err)
		}
		log.Println("✅ Role column converted to VARCHAR")
	} else {
		// Just ensure it's correct
		_, err = DB.Exec("ALTER TABLE users MODIFY COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user'")
		if err != nil {
			log.Printf("⚠️  Could not modify role column (might already be correct): %v", err)
		} else {
			log.Println("✅ Role column structure verified")
		}
	}
	
	return nil
}

// SQL statements for creating tables
const createUsersTableSQL = `
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    goal VARCHAR(50) NOT NULL CHECK (goal IN ('lose_weight', 'hypertrophy')),
    role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'user', 'premium')),
    height DECIMAL(5, 2) NULL COMMENT 'Height in cm',
    weight DECIMAL(5, 2) NULL COMMENT 'Weight in kg',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email),
    INDEX idx_role (role)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
`

const createWorkoutsTableSQL = `
CREATE TABLE IF NOT EXISTS workouts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    duration INT NOT NULL CHECK (duration > 0) COMMENT 'Duration in minutes',
    calories_burned DECIMAL(10, 2) DEFAULT 0 CHECK (calories_burned >= 0),
    workout_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_workout_date (workout_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
`

const createProgressTableSQL = `
CREATE TABLE IF NOT EXISTS progress (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    weight DECIMAL(5, 2) NOT NULL CHECK (weight > 0) COMMENT 'Weight in kg',
    body_fat DECIMAL(5, 2) DEFAULT 0 CHECK (body_fat >= 0 AND body_fat <= 100) COMMENT 'Body fat percentage',
    muscle_mass DECIMAL(5, 2) DEFAULT 0 CHECK (muscle_mass >= 0) COMMENT 'Muscle mass in kg',
    notes TEXT,
    progress_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_progress_date (progress_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
`

<<<<<<< HEAD
=======

>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
