@echo off
echo ========================================
echo Fixing Progress Table
echo ========================================
echo.

echo Adding progress_date column if missing...
mysql -u root -pVojislav123! app_db < migrations\002_fix_progress_date.sql

if %errorlevel% neq 0 (
    echo.
    echo ❌ ERROR: Failed to fix progress table!
    echo.
    pause
    exit /b 1
)

echo.
echo ✅ Progress table fixed!
echo.
pause

