-- Korišćenje baze podataka
USE app_db;

-- Dodavanje kolone progress_date u tabelu progress
ALTER TABLE progress 
ADD COLUMN progress_date DATE NOT NULL DEFAULT (CURDATE()) AFTER notes;

-- Ažuriranje postojećih redova koji imaju nevalidne datume
UPDATE progress 
SET progress_date = CURDATE() 
WHERE progress_date IS NULL OR progress_date = '0000-00-00';

-- Kreiranje indeksa za bržu pretragu po datumu progresa
CREATE INDEX idx_progress_date ON progress(progress_date);
