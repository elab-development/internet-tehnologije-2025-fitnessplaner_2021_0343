@echo off
echo ========================================
echo RECREATING DATABASE FROM SCRATCH
echo ========================================
echo.
echo WARNING: This will DELETE all existing data!
echo.
set /p confirm="Type 'yes' to continue: "
if /i not "%confirm%"=="yes" (
    echo Cancelled.
    pause
    exit /b 0
)

echo.
echo Dropping existing tables...
mysql -u root -pVojislav123! app_db -e "DROP TABLE IF EXISTS progress; DROP TABLE IF EXISTS workouts; DROP TABLE IF EXISTS users;" 2>nul

echo.
echo Creating fresh tables...
mysql -u root -pVojislav123! app_db < migrations\001_init.sql

if %errorlevel% neq 0 (
    echo.
    echo ❌ ERROR: Failed to create tables!
    echo.
    pause
    exit /b 1
)

echo.
echo ✅ Database recreated successfully!
echo.
echo Tables created:
mysql -u root -pVojislav123! app_db -e "SHOW TABLES;" 2>nul

echo.
echo Workouts table structure:
mysql -u root -pVojislav123! app_db -e "DESCRIBE workouts;" 2>nul

echo.
echo Progress table structure:
mysql -u root -pVojislav123! app_db -e "DESCRIBE progress;" 2>nul

echo.
echo ✅ Now restart your backend server!
echo.
pause

