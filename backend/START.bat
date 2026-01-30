@echo off
echo ========================================
echo Starting Application
echo ========================================
echo.
echo Backend: http://localhost:8080
echo Frontend: http://localhost:5173
echo.

start "Backend" cmd /k "cd backend && go mod tidy && go run main.go"
timeout /t 2 /nobreak >nul
start "Frontend" cmd /k "cd frontend && npm install && npm run dev"

echo Servers starting...
pause

