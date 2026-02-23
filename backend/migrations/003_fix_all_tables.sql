<<<<<<< HEAD
-- Fix all tables - add missing columns
USE app_db;

-- Fix workouts table - add calories_burned
ALTER TABLE workouts 
ADD COLUMN calories_burned DECIMAL(10, 2) DEFAULT 0 AFTER duration;

-- Fix progress table - add progress_date  
ALTER TABLE progress 
ADD COLUMN progress_date DATE NOT NULL DEFAULT (CURDATE()) AFTER notes;

-- Update existing rows in progress
=======
USE app_db;

ALTER TABLE workouts 
ADD COLUMN calories_burned DECIMAL(10, 2) DEFAULT 0 AFTER duration;
 
ALTER TABLE progress 
ADD COLUMN progress_date DATE NOT NULL DEFAULT (CURDATE()) AFTER notes;

>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
UPDATE progress 
SET progress_date = CURDATE() 
WHERE progress_date IS NULL OR progress_date = '0000-00-00';

<<<<<<< HEAD
-- Create indexes
=======
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
CREATE INDEX idx_workout_date ON workouts(workout_date);
CREATE INDEX idx_progress_date ON progress(progress_date);
