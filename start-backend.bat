@echo off
cd /d "%~dp0backend"
echo Starting Surya Photography API on http://localhost:8080
go run ./cmd/server
pause
