<<<<<<< HEAD
-- Fix progress table - add progress_date column
USE app_db;

-- Try to add the column (will fail if it already exists, which is OK)
ALTER TABLE progress 
ADD COLUMN progress_date DATE NOT NULL DEFAULT (CURDATE()) AFTER notes;

-- Update existing rows
=======
Dodavanje progress_date column
USE app_db;

ALTER TABLE progress 
ADD COLUMN progress_date DATE NOT NULL DEFAULT (CURDATE()) AFTER notes;

>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
UPDATE progress 
SET progress_date = CURDATE() 
WHERE progress_date IS NULL OR progress_date = '0000-00-00';

<<<<<<< HEAD
-- Add index
=======
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
CREATE INDEX idx_progress_date ON progress(progress_date);
