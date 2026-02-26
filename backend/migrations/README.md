# Sistem Migracija Baze Podataka

## Pregled

Sistem migracija omogućava automatsko upravljanje promenama u strukturi baze podataka. Migracije se izvršavaju automatski pri pokretanju backend servera.

## Kako funkcioniše

1. **Tabela za praćenje**: Sistem automatski kreira tabelu `schema_migrations` koja prati koje migracije su već primenjene.

2. **Automatsko izvršavanje**: Pri pokretanju servera, sistem:
   - Proverava koje migracije već postoje u bazi
   - Pronalazi sve `.sql` fajlove u `migrations/` folderu
   - Sortira ih po imenu (001, 002, 003...)
   - Izvršava samo one migracije koje još nisu primenjene

3. **Idempotentnost**: Migracije se izvršavaju samo jednom. Ako je migracija već primenjena, preskače se.

## Struktura migracija

Migracije se čuvaju u `backend/migrations/` folderu sa imenima:
- `001_init.sql` - Inicijalna migracija (kreiranje osnovnih tabela)
- `002_add_feature.sql` - Primer dodatne migracije
- `003_update_schema.sql` - Primer ažuriranja šeme

**Važno**: Imena migracija moraju biti numerisana (001, 002, 003...) da bi se izvršavale u ispravnom redosledu.

## Kreiranje nove migracije

1. **Kreiraj novi SQL fajl** u `migrations/` folderu:
   ```bash
   # Primer: 002_add_meal_plans.sql
   ```

2. **Dodaj SQL naredbe**:
   ```sql
   -- Dodavanje nove tabele
   CREATE TABLE IF NOT EXISTS meal_plans (
       id INT AUTO_INCREMENT PRIMARY KEY,
       user_id INT NOT NULL,
       name VARCHAR(255) NOT NULL,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
   ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
   ```

3. **Restartuj backend server** - migracija će se automatski primeniti

## Provera statusa migracija

Možeš proveriti koje migracije su primenjene direktno u bazi:

```sql
SELECT * FROM schema_migrations ORDER BY applied_at;
```

## Ručno izvršavanje migracija

Ako želiš da ručno pokreneš migracije (npr. za testiranje):

```go
// U Go kodu
if err := utils.RunMigrations(); err != nil {
    log.Fatal(err)
}
```

## Napomene

- **Backup**: Uvek napravi backup baze pre primene migracija u produkciji
- **Testiranje**: Testiraj migracije na development okruženju pre produkcije
- **Redosled**: Migracije se izvršavaju po redosledu imena fajlova
- **Idempotentnost**: Koristi `IF NOT EXISTS` i `IF EXISTS` u SQL naredbama gde je moguće

## Trenutne migracije

- `001_init.sql` - Kreiranje osnovnih tabela (users, workouts, progress)
- `002_fix_progress_date.sql` - Dodavanje `progress_date` kolone u `progress` tabelu
- `003_fix_all_tables.sql` - Dodavanje nedostajućih kolona (`calories_burned` u `workouts`, `progress_date` u `progress`) i kreiranje indeksa

## Napomene o greškama

Sistem automatski ignorise sledeće greške koje su očekivane:
- `Duplicate column` - kolona već postoji
- `Duplicate key name` - indeks već postoji  
- `already exists` - tabela/baza već postoji

Ove greške su normalne kada se migracije pokreću više puta ili kada se primenjuju na bazu koja već ima neke strukture.

