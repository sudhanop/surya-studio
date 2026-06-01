@echo off
setlocal EnableExtensions
title Surya Photography - Production
cd /d "%~dp0"
if not exist "logs" mkdir "logs"

if exist ".env" (
  echo Loading environment from .env
  for /f "usebackq eol=# tokens=1,* delims==" %%A in (".env") do set "%%A=%%B"
) else (
  echo No .env found. Using default local settings.
)

if "%APP_ENV%"=="" set "APP_ENV=production"
if "%API_PORT%"=="" set "API_PORT=8080"
if "%WEB_PORT%"=="" set "WEB_PORT=3000"

if "%NEXT_PUBLIC_API_URL%"=="" set "NEXT_PUBLIC_API_URL=http://localhost:%API_PORT%"
if "%NEXT_PUBLIC_SITE_URL%"=="" set "NEXT_PUBLIC_SITE_URL=http://localhost:%WEB_PORT%"
if "%NEXT_PUBLIC_MEDIA_URL%"=="" set "NEXT_PUBLIC_MEDIA_URL=http://localhost:%API_PORT%/uploads"

if "%FRONTEND_URL%"=="" set "FRONTEND_URL=%NEXT_PUBLIC_SITE_URL%"
if "%CORS_ORIGINS%"=="" set "CORS_ORIGINS=%NEXT_PUBLIC_SITE_URL%"
if "%PUBLIC_MEDIA_URL%"=="" set "PUBLIC_MEDIA_URL=%NEXT_PUBLIC_MEDIA_URL%"

echo.
echo [1/2] Starting API...
start "Surya API" cmd /k "cd /d %~dp0backend && set \"APP_ENV=%APP_ENV%\" && set \"PORT=%API_PORT%\" && set \"FRONTEND_URL=%FRONTEND_URL%\" && set \"CORS_ORIGINS=%CORS_ORIGINS%\" && set \"JWT_SECRET=%JWT_SECRET%\" && set \"PUBLIC_MEDIA_URL=%PUBLIC_MEDIA_URL%\" && if not exist server.exe (where /q go || (echo Go is not installed or not in PATH. Install Go and retry.& pause& exit /b 1) && go build -o server.exe ./cmd/server) && server.exe"
timeout /t 3 /nobreak >nul
echo [2/2] Starting website...
start "Surya Web" cmd /k "cd /d %~dp0frontend && where /q npm.cmd || (echo Node.js/npm is not installed or not in PATH. Install Node.js LTS and retry.& pause& exit /b 1) && if not exist node_modules call npm.cmd install && set \"NEXT_PUBLIC_API_URL=%NEXT_PUBLIC_API_URL%\" && set \"NEXT_PUBLIC_SITE_URL=%NEXT_PUBLIC_SITE_URL%\" && set \"NEXT_PUBLIC_MEDIA_URL=%NEXT_PUBLIC_MEDIA_URL%\" && call npm.cmd run build && set \"NODE_ENV=production\" && set \"HOSTNAME=0.0.0.0\" && set \"PORT=%WEB_PORT%\" && node .next\\standalone\\server.js"
echo.
echo API:  http://localhost:%API_PORT%
echo Site: http://localhost:%WEB_PORT%
echo Admin: http://localhost:%WEB_PORT%/admin/login
echo.
pause
