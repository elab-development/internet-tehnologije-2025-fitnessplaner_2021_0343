
USE app_db;

-- Dropovanje postojećih tabela
DROP TABLE IF EXISTS progress;
DROP TABLE IF EXISTS workouts;
DROP TABLE IF EXISTS users;

-- Kreiranje users tabele
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
);

-- Kreiranje workouts tabele
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
);

-- Kreiranje progress tabele
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
);

-- Verifikovanje da su tabele kreirane
SELECT 'Tables created successfully!' AS Status;
SHOW TABLES;

-- Prikaz strukture workouts tabele
DESCRIBE workouts;

-- Prikaz strukture progress tabele
DESCRIBE progress;


