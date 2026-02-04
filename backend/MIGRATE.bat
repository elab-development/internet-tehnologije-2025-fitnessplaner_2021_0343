@echo off
echo ========================================
echo Database Migration
echo ========================================
echo.

set /p DB_USER="MySQL username (default: root): "
if "%DB_USER%"=="" set DB_USER=root

set /p DB_PASS="MySQL password: "
set /p DB_NAME="Database name (default: app_db): "
if "%DB_NAME%"=="" set DB_NAME=app_db

echo.
echo Running migration...
mysql -u %DB_USER% -p%DB_PASS% < backend\migrations\001_init.sql

echo.
echo Migration completed!
pause
