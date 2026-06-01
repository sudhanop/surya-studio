@echo off
cd /d "%~dp0backend"
:: SQLEXPRESS on port 1433 (SQL Browser can stay off). Change if your instance differs.
set DB_SERVER=localhost,1433
echo Starting Surya Photography API on http://localhost:8080
echo Database: %DB_SERVER% / SuryaPhotography
go run ./cmd/server
pause
