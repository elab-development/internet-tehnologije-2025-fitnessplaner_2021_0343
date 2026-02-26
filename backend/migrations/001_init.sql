-- Kompletna inicijalizacija baze podataka
-- Pokreni ovaj fajl da postaviš celu bazu podataka

CREATE DATABASE IF NOT EXISTS app_db;
USE app_db;

-- Kreiranje tabele usera sa kolonom role
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    goal VARCHAR(50) NOT NULL CHECK (goal IN ('lose_weight', 'hypertrophy')),
    role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'user', 'premium')),
    height DECIMAL(5, 2) NULL COMMENT 'Visina u cm',
    weight DECIMAL(5, 2) NULL COMMENT 'Težina u kg',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email),
    INDEX idx_role (role)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Kreiranje tabele za vezbe
CREATE TABLE IF NOT EXISTS workouts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    duration INT NOT NULL CHECK (duration > 0) COMMENT 'Trajanje u minutama',
    calories_burned DECIMAL(10, 2) DEFAULT 0 CHECK (calories_burned >= 0),
    workout_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_workout_date (workout_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Kreiranje tabele za progres
CREATE TABLE IF NOT EXISTS progress (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    weight DECIMAL(5, 2) NOT NULL CHECK (weight > 0) COMMENT 'Težina u kg',
    body_fat DECIMAL(5, 2) DEFAULT 0 CHECK (body_fat >= 0 AND body_fat <= 100) COMMENT 'Procenat masti u telu',
    muscle_mass DECIMAL(5, 2) DEFAULT 0 CHECK (muscle_mass >= 0) COMMENT 'Mišićna masa u kg',
    notes TEXT,
    progress_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_progress_date (progress_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
