-- Korišćenje baze podataka
USE app_db;

-- Popravka workouts tabele - dodavanje calories_burned
ALTER TABLE workouts 
ADD COLUMN calories_burned DECIMAL(10, 2) DEFAULT 0 AFTER duration;

-- Popravka progress tabele - dodavanje progress_date
ALTER TABLE progress 
ADD COLUMN progress_date DATE NOT NULL DEFAULT (CURDATE()) AFTER notes;

-- Ažuriranje postojećih redova u progress tabeli
UPDATE progress 
SET progress_date = CURDATE() 
WHERE progress_date IS NULL OR progress_date = '0000-00-00';

-- Kreiranje indeksa za bržu pretragu
CREATE INDEX idx_workout_date ON workouts(workout_date);
CREATE INDEX idx_progress_date ON progress(progress_date);
