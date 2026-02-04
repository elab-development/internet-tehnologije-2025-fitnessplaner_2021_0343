USE app_db;

ALTER TABLE workouts 
ADD COLUMN calories_burned DECIMAL(10, 2) DEFAULT 0 AFTER duration;
 
ALTER TABLE progress 
ADD COLUMN progress_date DATE NOT NULL DEFAULT (CURDATE()) AFTER notes;

UPDATE progress 
SET progress_date = CURDATE() 
WHERE progress_date IS NULL OR progress_date = '0000-00-00';

CREATE INDEX idx_workout_date ON workouts(workout_date);
CREATE INDEX idx_progress_date ON progress(progress_date);
