@echo off
echo ========================================
echo Recreating All Tables
echo ========================================
echo.
echo WARNING: This will DROP and RECREATE all tables!
echo All data will be lost!
echo.
set /p confirm="Are you sure? (yes/no): "
if /i not "%confirm%"=="yes" (
    echo Cancelled.
    pause
    exit /b 0
)

echo.
echo Dropping existing tables...
mysql -u root -pVojislav123! app_db -e "DROP TABLE IF EXISTS progress; DROP TABLE IF EXISTS workouts; DROP TABLE IF EXISTS users;"

echo.
echo Creating tables from migration...
mysql -u root -pVojislav123! app_db < migrations\001_init.sql

if %errorlevel% neq 0 (
    echo.
    echo ❌ ERROR: Failed to recreate tables!
    echo.
    pause
    exit /b 1
)

echo.
echo ✅ All tables recreated successfully!
echo.
pause

