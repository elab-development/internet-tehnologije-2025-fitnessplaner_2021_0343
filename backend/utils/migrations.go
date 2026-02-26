package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// initMigrationsTable kreira tabelu za praćenje migracija ako ne postoji
func initMigrationsTable() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	if _, err := DB.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	return nil
}

// isMigrationApplied proverava da li je migracija već primenjena
func isMigrationApplied(version string) (bool, error) {
	var count int
	err := DB.QueryRow(
		"SELECT COUNT(*) FROM schema_migrations WHERE version = ?",
		version,
	).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("failed to check migration status: %w", err)
	}

	return count > 0, nil
}

// markMigrationApplied označava migraciju kao primenjenu
func markMigrationApplied(version string) error {
	_, err := DB.Exec(
		"INSERT INTO schema_migrations (version) VALUES (?)",
		version,
	)
	if err != nil {
		return fmt.Errorf("failed to mark migration as applied: %w", err)
	}
	return nil
}

// RunMigrations izvršava sve migracije iz migrations/ foldera
func RunMigrations() error {
	log.Println("🔄 Pokretanje migracija baze podataka...")

	// Kreiranje tabele za praćenje migracija
	if err := initMigrationsTable(); err != nil {
		return fmt.Errorf("failed to initialize migrations table: %w", err)
	}

	// Pronalaženje svih SQL fajlova u migrations/ folderu
	migrationsDir := "migrations"
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Filtriranje i sortiranje migracija
	var migrationFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			migrationFiles = append(migrationFiles, entry.Name())
		}
	}

	if len(migrationFiles) == 0 {
		log.Println("⚠️  Nema migracija za izvršavanje")
		return nil
	}

	// Sortiranje po imenu (001_init.sql, 002_add_column.sql, itd.)
	sort.Strings(migrationFiles)

	// Izvršavanje svake migracije
	for _, filename := range migrationFiles {
		version := strings.TrimSuffix(filename, ".sql")

		// Provera da li je migracija već primenjena
		applied, err := isMigrationApplied(version)
		if err != nil {
			return fmt.Errorf("failed to check migration %s: %w", version, err)
		}

		if applied {
			log.Printf("⏭️  Migracija %s već je primenjena, preskače se", version)
			continue
		}

		log.Printf("📝 Primena migracije: %s", filename)

		// Čitanje migration fajla
		migrationPath := filepath.Join(migrationsDir, filename)
		migrationSQL, err := os.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
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
				errStr := strings.ToLower(err.Error())
				// Ignorisanje grešaka koje su očekivane (tabela/kolona/indeks već postoji)
				if strings.Contains(errStr, "already exists") ||
					strings.Contains(errStr, "duplicate column") ||
					strings.Contains(errStr, "duplicate key name") ||
					strings.Contains(errStr, "database exists") {
					log.Printf("⏭️  Već postoji, preskače se: %s", statement[:min(50, len(statement))])
					continue
				}
				log.Printf("❌ Greška pri izvršavanju migracije %s, naredba %d: %v", version, i+1, err)
				log.Printf("Naredba: %s", statement[:min(100, len(statement))])
				return fmt.Errorf("migration %s failed at statement %d: %w", version, i+1, err)
			}
		}

		// Označavanje migracije kao primenjene
		if err := markMigrationApplied(version); err != nil {
			return fmt.Errorf("failed to mark migration %s as applied: %w", version, err)
		}

		log.Printf("✅ Migracija %s uspešno primenjena", version)
	}

	log.Println("✅ Sve migracije uspešno završene")
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
