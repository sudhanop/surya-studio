# Surya Photography — Database Setup

## Quick start (Windows)

1. Edit `create-database.bat` — set `SQL_SERVER` (e.g. `localhost\SQLEXPRESS`).
2. Double-click `create-database.bat` or run from terminal:

```bat
cd database
create-database.bat
```

## Alternative: run in SSMS

1. Open **SQL Server Management Studio**.
2. Connect to your instance.
3. **File → Open → File** → select `schema.sql`.
4. Press **F5** to execute.

## Schema overview

| Table | Purpose |
|-------|---------|
| `admin` | Single admin account (bcrypt password) |
| `categories` | Portfolio categories (9 seeded) |
| `portfolio_media` | Photos/videos per category |
| `inquiries` | Contact form submissions |
| `functions` | Studio production tracker |

## Default admin

- **Username:** `surya@admin.com`
- **Password:** `surya@1995`

To update credentials on an existing database, run `update-admin-credentials.sql`.

## Connection string (Go backend)

If **SQL Server Browser** is disabled (common on Express), use the TCP port instead of `localhost\SQLEXPRESS`:

```
server=localhost,1433;database=SuryaPhotography;trusted_connection=yes;encrypt=disable
```

In SSMS / sqlcmd, `.\SQLEXPRESS` still works. Check your port in SSMS: Server Properties → Connection → or registry `TcpDynamicPorts` under `MSSQL17.SQLEXPRESS`.

Or with SQL auth:

```
server=localhost;user id=sa;password=YOUR_PASSWORD;database=SuryaPhotography;encrypt=disable
```
