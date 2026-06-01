@echo off
title Surya Photography - Production
echo.
echo [1/2] Starting API...
start "Surya API" cmd /k "cd /d %~dp0backend && set APP_ENV=production && go run ./cmd/server"
timeout /t 3 /nobreak >nul
echo [2/2] Starting website...
start "Surya Web" cmd /k "cd /d %~dp0frontend && npm run start"
echo.
echo API:  http://localhost:8080
echo Site: http://localhost:3000
echo Admin: http://localhost:3000/admin/login
echo.
pause
