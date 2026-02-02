@echo off
echo ========================================
echo Quick Fix - Adding Missing Columns
echo ========================================
echo.

echo Adding calories_burned to workouts table...
mysql -u root -pVojislav123! app_db -e "ALTER TABLE workouts ADD COLUMN calories_burned DECIMAL(10, 2) DEFAULT 0 AFTER duration;" 2>nul
if %errorlevel% equ 0 (
    echo ✅ Added calories_burned column
) else (
    echo ⚠️  Column may already exist (this is OK)
)

echo.
echo Adding progress_date to progress table...
mysql -u root -pVojislav123! app_db -e "ALTER TABLE progress ADD COLUMN progress_date DATE NOT NULL DEFAULT (CURDATE()) AFTER notes;" 2>nul
if %errorlevel% equ 0 (
    echo ✅ Added progress_date column
) else (
    echo ⚠️  Column may already exist (this is OK)
)

echo.
echo ✅ Quick fix completed!
echo.
echo Now restart your backend server.
echo.
pause

