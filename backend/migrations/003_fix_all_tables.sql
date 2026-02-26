<<<<<<< HEAD
-- Popravka svih tabela - dodavanje nedostajućih kolona
-- Ova migracija dodaje kolone koje možda nedostaju u postojećim tabelama
-- Greške "Duplicate column" i "Duplicate key name" će biti ignorisane

-- Popravka workouts tabele - dodavanje calories_burned
ALTER TABLE workouts 
ADD COLUMN calories_burned DECIMAL(10, 2) DEFAULT 0 AFTER duration;

-- Popravka progress tabele - dodavanje progress_date
ALTER TABLE progress 
ADD COLUMN progress_date DATE NOT NULL DEFAULT (CURDATE()) AFTER notes;

-- Ažuriranje postojećih redova u progress
=======
<<<<<<< HEAD
-- Ispravka svih tabela
USE app_db;

-- Isparvka workouts tabele - dodavanje calories_burned
ALTER TABLE workouts 
ADD COLUMN calories_burned DECIMAL(10, 2) DEFAULT 0 AFTER duration;

-- Ispravka progress tabele - dodavanje progress_date  
ALTER TABLE progress 
ADD COLUMN progress_date DATE NOT NULL DEFAULT (CURDATE()) AFTER notes;

-- Azuriranje postojecih redova u progress-u
=======
USE app_db;

ALTER TABLE workouts 
ADD COLUMN calories_burned DECIMAL(10, 2) DEFAULT 0 AFTER duration;
 
ALTER TABLE progress 
ADD COLUMN progress_date DATE NOT NULL DEFAULT (CURDATE()) AFTER notes;

>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
>>>>>>> 7bd8328a2f3e453abbb64671ca0e40634a9cdef7
UPDATE progress 
SET progress_date = CURDATE() 
WHERE progress_date IS NULL OR progress_date = '0000-00-00';

<<<<<<< HEAD
-- Kreiranje indeksa
=======
<<<<<<< HEAD
-- Kreiranje indeksa
=======
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
>>>>>>> 7bd8328a2f3e453abbb64671ca0e40634a9cdef7
CREATE INDEX idx_workout_date ON workouts(workout_date);
CREATE INDEX idx_progress_date ON progress(progress_date);
