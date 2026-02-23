<<<<<<< HEAD
-- ============================================
-- COMPLETE DATABASE SETUP FOR FITNESS APP
-- ============================================
-- Kopirajte i nalepite ceo ovaj kod u MySQL Workbench
-- i pokrenite ga (Execute button)
-- ============================================
=======
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608

-- Kreiraj bazu ako ne postoji
CREATE DATABASE IF NOT EXISTS app_db;
USE app_db;

<<<<<<< HEAD
-- ============================================
-- 1. BRISANJE POSTOJEĆIH TABELA (sa svim zavisnostima)
-- ============================================
=======
-- 1. BRISANJE POSTOJEĆIH TABELA (sa svim zavisnostima)
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
-- Prvo obriši sve tabele koje imaju foreign key-ove
DROP TABLE IF EXISTS workout_exercises;
DROP TABLE IF EXISTS progress;
DROP TABLE IF EXISTS workouts;
DROP TABLE IF EXISTS users;

<<<<<<< HEAD
-- ============================================
-- 2. KREIRANJE USERS TABELE
-- ============================================
=======

-- 2. KREIRANJE USERS TABELE
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    goal VARCHAR(50) NOT NULL CHECK (goal IN ('lose_weight', 'hypertrophy')),
    role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'user', 'premium')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email),
    INDEX idx_role (role)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

<<<<<<< HEAD
-- ============================================
-- 3. KREIRANJE WORKOUTS TABELE
-- ============================================
=======

-- 3. KREIRANJE WORKOUTS TABELE
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
CREATE TABLE workouts (
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

<<<<<<< HEAD
-- ============================================
-- 4. KREIRANJE PROGRESS TABELE
-- ============================================
=======
-- 4. KREIRANJE PROGRESS TABELE
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
CREATE TABLE progress (
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

<<<<<<< HEAD
-- ============================================
-- 5. PROVERA - PRIKAZ KREIRANIH TABELA
-- ============================================
=======

-- 5. PROVERA - PRIKAZ KREIRANIH TABELA
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
SELECT '✅ Database setup completed successfully!' AS Status;

-- Prikaži sve tabele
SHOW TABLES;

-- Prikaži strukturu users tabele
SELECT '=== USERS TABLE STRUCTURE ===' AS Info;
DESCRIBE users;

-- Prikaži strukturu workouts tabele
SELECT '=== WORKOUTS TABLE STRUCTURE ===' AS Info;
DESCRIBE workouts;

-- Prikaži strukturu progress tabele
SELECT '=== PROGRESS TABLE STRUCTURE ===' AS Info;
DESCRIBE progress;

<<<<<<< HEAD
-- ============================================
-- 6. FINALNA PROVERA - BROJ TABELA
-- ============================================
=======
-- 6. FINALNA PROVERA - BROJ TABELA
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
SELECT COUNT(*) AS 'Ukupno tabela' FROM information_schema.tables 
WHERE table_schema = 'app_db';

SELECT '✅ Setup complete! You can now start your backend server.' AS Final_Message;
