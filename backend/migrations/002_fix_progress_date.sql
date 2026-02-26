-- Dodavanje progress_date kolone u progress tabelu
-- Ova migracija dodaje kolonu progress_date ako ne postoji
-- Greške "Duplicate column" će biti ignorisane

-- Dodavanje kolone
ALTER TABLE progress 
ADD COLUMN progress_date DATE NOT NULL DEFAULT (CURDATE()) AFTER notes;

-- Ažuriranje postojećih redova
UPDATE progress 
SET progress_date = CURDATE() 
WHERE progress_date IS NULL OR progress_date = '0000-00-00';

-- Kreiranje indeksa
CREATE INDEX idx_progress_date ON progress(progress_date);
