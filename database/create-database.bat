@echo off
setlocal EnableExtensions EnableDelayedExpansion

:: =============================================================================
::  Surya Photography — Create SQL Server Database
::  Double-click this file OR run from Command Prompt in the database folder.
::
::  Prerequisites:
::    - SQL Server installed (Express / Developer / Standard)
::    - sqlcmd in PATH (installed with SQL Server tools)
::    - SSMS optional (you can also open schema.sql in SSMS and press F5)
:: =============================================================================

title Surya Photography - Database Setup

:: ----- CONFIGURE THESE VALUES -----------------------------------------------
set "SQL_SERVER=localhost,1433"
set "SQL_INSTANCE="
:: Examples (pick ONE that works on your PC):
::   set "SQL_SERVER=localhost,1433          (SQLEXPRESS on port 1433 — use if Browser is off)
::   set "SQL_SERVER=.\SQLEXPRESS            (works in SSMS / sqlcmd)
::   set "SQL_SERVER=localhost\SQLEXPRESS    (needs SQL Server Browser running)

:: Authentication: WINDOWS or SQL
set "AUTH_MODE=WINDOWS"

:: Only used when AUTH_MODE=SQL
set "SQL_USER=sa"
set "SQL_PASSWORD=YourStrongPassword123"

:: Database name (must match schema.sql)
set "DATABASE_NAME=SuryaPhotography"
:: ----------------------------------------------------------------------------

set "SCRIPT_DIR=%~dp0"
set "SCHEMA_FILE=%SCRIPT_DIR%schema.sql"

if not exist "%SCHEMA_FILE%" (
    echo [ERROR] schema.sql not found at:
    echo   %SCHEMA_FILE%
    pause
    exit /b 1
)

where sqlcmd >nul 2>&1
if errorlevel 1 (
    echo [ERROR] sqlcmd not found in PATH.
    echo Install "SQL Server Command Line Utilities" or SQL Server Management Studio tools.
    echo Download: https://learn.microsoft.com/sql/tools/sqlcmd/sqlcmd-utility
    pause
    exit /b 1
)

if defined SQL_INSTANCE (
    set "SERVER_TARGET=%SQL_SERVER%\%SQL_INSTANCE%"
) else (
    set "SERVER_TARGET=%SQL_SERVER%"
)

echo.
echo ============================================================
echo   Surya Photography - Database Setup
echo ============================================================
echo   Server:   %SERVER_TARGET%
echo   Database: %DATABASE_NAME%
echo   Auth:     %AUTH_MODE%
echo   Script:   %SCHEMA_FILE%
echo ============================================================
echo.

if /I "%AUTH_MODE%"=="SQL" (
  if "%SQL_PASSWORD%"=="YourStrongPassword123" (
    echo [WARNING] You are using the default SQL password placeholder.
    echo           Edit create-database.bat and set SQL_PASSWORD before continuing.
    echo.
    choice /C YN /M "Continue anyway"
    if errorlevel 2 exit /b 1
  )
  set "SQLCMD_AUTH=-S %SERVER_TARGET% -U %SQL_USER% -P %SQL_PASSWORD%"
) else (
  set "SQLCMD_AUTH=-S %SERVER_TARGET% -E"
)

echo Running schema.sql ...
echo.

sqlcmd %SQLCMD_AUTH% -b -i "%SCHEMA_FILE%"

if errorlevel 1 (
    echo.
    echo [FAILED] Database setup encountered errors.
    echo.
    echo Troubleshooting:
    echo   1. Verify SQL Server is running (services.msc - SQL Server service)
    echo   2. Set SQL_SERVER to your instance, e.g. localhost\SQLEXPRESS
    echo   3. For named instance, use AUTH_MODE=WINDOWS if you use Windows login in SSMS
    echo   4. Or set AUTH_MODE=SQL with correct SQL_USER and SQL_PASSWORD
    echo   5. Run SSMS as Administrator if permission errors occur
    echo.
    pause
    exit /b 1
)

echo.
echo [SUCCESS] Database "%DATABASE_NAME%" is ready.
echo.
echo Tables created:
echo   - admin
echo   - categories          (9 portfolio categories seeded)
echo   - portfolio_media
echo   - inquiries
echo   - functions
echo.
echo Default admin login:
echo   Username: surya@admin.com
echo   Password: surya@1995
echo.
echo Connect in SSMS:
echo   Server: %SERVER_TARGET%
echo   Database: %DATABASE_NAME%
echo.
echo Backend connection string example:
echo   server=%SERVER_TARGET%;database=%DATABASE_NAME%;trusted_connection=true;
echo.

pause
exit /b 0
