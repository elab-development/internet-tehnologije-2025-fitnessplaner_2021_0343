package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// RunMigrations izvršava migracije baze podataka
func RunMigrations() error {
	log.Println("🔄 Pokretanje migracija baze podataka...")

	// Čitanje migration fajla
	migrationSQL, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Podela SQL-a na pojedinačne naredbe
	statements := splitSQL(string(migrationSQL))

	// Izvršavanje svake naredbe
	for i, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" || strings.HasPrefix(statement, "--") {
			continue
		}

		// Preskakanje CREATE DATABASE ako već postoji
		if strings.HasPrefix(strings.ToUpper(statement), "CREATE DATABASE") {
			log.Println("⏭️  Preskakanje CREATE DATABASE (već postoji)")
			continue
		}

		// Izvršavanje naredbe
		if _, err := DB.Exec(statement); err != nil {
			// Ignorisanje grešaka "tabela već postoji"
			if strings.Contains(err.Error(), "already exists") {
				log.Printf("⏭️  Tabela već postoji, preskače se...")
				continue
			}
			// Ignorisanje grešaka "baza postoji"
			if strings.Contains(err.Error(), "database exists") {
				continue
			}
			log.Printf("❌ Error executing migration statement %d: %v", i+1, err)
			log.Printf("Statement: %s", statement[:min(100, len(statement))])
			return fmt.Errorf("migration failed at statement %d: %w", i+1, err)
		}
	}

	log.Println("✅ Migracije uspešno završene")
	return nil
}

// splitSQL deli SQL string na pojedinačne naredbe
func splitSQL(sql string) []string {
	// Uklanjanje komentara i podela po tački-zapeti
	lines := strings.Split(sql, "\n")
	var cleanedLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Preskakanje praznih linija i komentara
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}
		cleanedLines = append(cleanedLines, line)
	}

	// Spajanje i podela po tački-zapeti
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

// min vraća minimum od dva cela broja
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// EnsureTablesExist proverava da li sve potrebne tabele postoje i kreira ih ako je potrebno
func EnsureTablesExist() error {
	log.Println("🔍 Provera da li sve tabele postoje...")

	// Provera i kreiranje users tabele
	if err := ensureTable("users", createUsersTableSQL); err != nil {
		return err
	}

	// Osiguravanje da users ima role kolonu (bez CHECK constraint-a da bi se izbegli problemi)
	if err := ensureColumn("users", "role", "ALTER TABLE users ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user' AFTER goal"); err != nil {
		return err
	}
	
	// Osiguravanje da users ima height kolonu
	if err := ensureColumn("users", "height", "ALTER TABLE users ADD COLUMN height DECIMAL(5, 2) NULL AFTER role"); err != nil {
		return err
	}
	
	// Osiguravanje da users ima weight kolonu
	if err := ensureColumn("users", "weight", "ALTER TABLE users ADD COLUMN weight DECIMAL(5, 2) NULL AFTER height"); err != nil {
		return err
	}
	
	// Popravka role kolone ako ima problema (uklanjanje CHECK constraint-a ako pravi probleme)
	if err := fixRoleColumn(); err != nil {
		log.Printf("⚠️  Warning: Could not fix role column: %v", err)
	}

	// Provera i kreiranje workouts tabele
	if err := ensureTable("workouts", createWorkoutsTableSQL); err != nil {
		return err
	}

	// Osiguravanje da workouts ima calories_burned kolonu
	if err := ensureColumn("workouts", "calories_burned", "ALTER TABLE workouts ADD COLUMN calories_burned DECIMAL(10, 2) DEFAULT 0 AFTER duration"); err != nil {
		return err
	}

	// Provera i kreiranje progress tabele
	if err := ensureTable("progress", createProgressTableSQL); err != nil {
		return err
	}

	// Osiguravanje da progress ima progress_date kolonu
	if err := ensureColumn("progress", "progress_date", "ALTER TABLE progress ADD COLUMN progress_date DATE NOT NULL DEFAULT (CURDATE()) AFTER notes"); err != nil {
		return err
	}

	// Popravka postojećih problema sa podacima
	if err := fixExistingData(); err != nil {
		log.Printf("⚠️  Warning: Failed to fix existing data: %v", err)
		// Ne propadati ako ovo baca grešku, samo logovati
	}

	log.Println("✅ Sve tabele i kolone postoje")
	return nil
}

// fixExistingData popravlja uobičajene probleme sa podacima u postojećim tabelama
func fixExistingData() error {
	// Popravka NULL lozinki (postavljanje na prazan string, korisnici će morati da resetuju)
	_, err := DB.Exec("UPDATE users SET password = '' WHERE password IS NULL")
	if err != nil {
		log.Printf("⚠️  Could not fix NULL passwords: %v", err)
	}

	// Pokušaj popravke nevažećih uloga - ovo može da ne uspe zbog CHECK constraint-a
	// Ako ne uspe, korisnici sa nevažećim ulogama će morati da se poprave ručno ili obrišu
	_, err = DB.Exec("UPDATE users SET role = 'user' WHERE role IS NULL OR role = '' OR role NOT IN ('admin', 'user', 'premium')")
	if err != nil {
		// Ovo je očekivano da ne uspe ako je CHECK constraint strog
		// Logovati ali ne propadati - korisnici mogu biti popravljeni ručno
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
			// Ignorisanje grešaka "Duplikat kolone"
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

// fixRoleColumn pokušava da popravi strukturu role kolone
func fixRoleColumn() error {
	// Provera trenutne strukture kolone
	var dataType, columnDefault, isNullable, columnType string
	err := DB.QueryRow(
		"SELECT DATA_TYPE, COLUMN_DEFAULT, IS_NULLABLE, COLUMN_TYPE FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = 'users' AND column_name = 'role'",
	).Scan(&dataType, &columnDefault, &isNullable, &columnType)
	
	if err != nil {
		return fmt.Errorf("failed to check role column: %w", err)
	}
	
	log.Printf("🔍 Role column: type=%s, full_type=%s, default=%s, nullable=%s", dataType, columnType, columnDefault, isNullable)
	
	// Ako je ENUM ili ima probleme sa CHECK constraint-om, modifikovati na jednostavan VARCHAR
	if dataType == "enum" || strings.Contains(strings.ToLower(columnType), "check") {
		log.Println("🔧 Role kolona je ENUM ili ima CHECK constraint, konvertovanje u VARCHAR...")
		_, err = DB.Exec("ALTER TABLE users MODIFY COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user'")
		if err != nil {
			log.Printf("⚠️  Could not modify role column: %v", err)
			return fmt.Errorf("failed to modify role column: %w", err)
		}
		log.Println("✅ Role kolona konvertovana u VARCHAR")
	} else {
		// Samo osigurati da je ispravna
		_, err = DB.Exec("ALTER TABLE users MODIFY COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user'")
		if err != nil {
			log.Printf("⚠️  Could not modify role column (might already be correct): %v", err)
		} else {
			log.Println("✅ Struktura role kolone verifikovana")
		}
	}
	
	return nil
}

// SQL naredbe za kreiranje tabela
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

