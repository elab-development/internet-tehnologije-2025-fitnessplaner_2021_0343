-- Quick fix - dodavanje kolona
USE app_db;

-- Dodavanje calories_burned u workouts
ALTER TABLE workouts ADD COLUMN calories_burned DECIMAL(10, 2) DEFAULT 0 AFTER duration;

-- Dodavanje progress_date u progress
ALTER TABLE progress ADD COLUMN progress_date DATE NOT NULL DEFAULT (CURDATE()) AFTER notes;


