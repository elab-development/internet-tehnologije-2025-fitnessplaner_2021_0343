@echo off
echo ========================================
echo Checking and Fixing Database Tables
echo ========================================
echo.

echo Step 1: Checking if tables exist...
mysql -u root -pVojislav123! app_db -e "SHOW TABLES;" 2>nul
if %errorlevel% neq 0 (
    echo.
    echo ❌ ERROR: Cannot connect to database or database doesn't exist!
    echo.
    echo Creating database and tables from scratch...
    mysql -u root -pVojislav123! < migrations\001_init.sql
    if %errorlevel% equ 0 (
        echo ✅ Database and tables created!
    ) else (
        echo ❌ Failed to create tables!
        pause
        exit /b 1
    )
) else (
    echo ✅ Tables exist, checking structure...
)

echo.
echo Step 2: Checking workouts table structure...
mysql -u root -pVojislav123! app_db -e "DESCRIBE workouts;" 2>nul

echo.
echo Step 3: Adding missing columns to workouts...
mysql -u root -pVojislav123! app_db -e "ALTER TABLE workouts ADD COLUMN IF NOT EXISTS calories_burned DECIMAL(10, 2) DEFAULT 0 AFTER duration;" 2>nul
if %errorlevel% neq 0 (
    echo Trying alternative method...
    mysql -u root -pVojislav123! app_db -e "ALTER TABLE workouts ADD COLUMN calories_burned DECIMAL(10, 2) DEFAULT 0 AFTER duration;" 2>nul
)

echo.
echo Step 4: Checking progress table structure...
mysql -u root -pVojislav123! app_db -e "DESCRIBE progress;" 2>nul

echo.
echo Step 5: Adding missing columns to progress...
mysql -u root -pVojislav123! app_db -e "ALTER TABLE progress ADD COLUMN IF NOT EXISTS progress_date DATE NOT NULL DEFAULT (CURDATE()) AFTER notes;" 2>nul
if %errorlevel% neq 0 (
    echo Trying alternative method...
    mysql -u root -pVojislav123! app_db -e "ALTER TABLE progress ADD COLUMN progress_date DATE NOT NULL DEFAULT (CURDATE()) AFTER notes;" 2>nul
)

echo.
echo Step 6: Final check - showing table structures...
echo.
echo === WORKOUTS TABLE ===
mysql -u root -pVojislav123! app_db -e "DESCRIBE workouts;" 2>nul
echo.
echo === PROGRESS TABLE ===
mysql -u root -pVojislav123! app_db -e "DESCRIBE progress;" 2>nul

echo.
echo ✅ Fix completed! Please restart your backend.
echo.
pause

