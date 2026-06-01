@echo off
cd /d "%~dp0frontend"
echo Starting Surya Photography website on http://localhost:3000
if not exist "node_modules" (
  echo Installing frontend dependencies...
  call npm.cmd install
)
call .\node_modules\.bin\next.cmd dev -p 3000
pause
