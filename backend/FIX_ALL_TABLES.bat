@echo off
echo ========================================
echo Fixing All Database Tables
echo ========================================
echo.

echo Fixing workouts and progress tables...
mysql -u root -pVojislav123! app_db < migrations\003_fix_all_tables.sql

if %errorlevel% neq 0 (
    echo.
    echo ❌ ERROR: Failed to fix tables!
    echo.
    echo Try running the SQL manually in MySQL Workbench or command line.
    echo.
    pause
    exit /b 1
)

echo.
echo ✅ All tables fixed!
echo.
pause

