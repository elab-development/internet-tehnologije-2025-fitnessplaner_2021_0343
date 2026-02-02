-- Quick fix - add missing columns (ignore errors if columns already exist)
USE app_db;

-- Add calories_burned to workouts
ALTER TABLE workouts ADD COLUMN calories_burned DECIMAL(10, 2) DEFAULT 0 AFTER duration;

-- Add progress_date to progress
ALTER TABLE progress ADD COLUMN progress_date DATE NOT NULL DEFAULT (CURDATE()) AFTER notes;

