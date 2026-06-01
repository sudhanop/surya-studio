# Deployment Guide — Surya Photography

## Pre-deploy checklist

- [ ] Run `database/create-database.bat` on production SQL Server
- [ ] Change admin password if needed (default `surya@1995`)
- [ ] Set strong `JWT_SECRET` (32+ random characters)
- [ ] Configure `SMTP_*` for inquiry email notifications
- [ ] Set real contact URLs in backend `.env`
- [ ] Point `PUBLIC_MEDIA_URL` to your public API uploads URL
- [ ] Set `CORS_ORIGINS` and `FRONTEND_URL` to your live domain

## Environment variables

### Backend (`backend/.env`)

| Variable | Example |
|----------|---------|
| `PORT` | `8080` |
| `APP_ENV` | `production` |
| `FRONTEND_URL` | `https://suryaphotography.com` |
| `CORS_ORIGINS` | `https://suryaphotography.com` |
| `DB_SERVER` | `your-server.database.windows.net` |
| `DB_DATABASE` | `SuryaPhotography` |
| `DB_USER` / `DB_PASSWORD` | SQL credentials |
| `DB_TRUSTED_CONNECTION` | `false` (cloud) |
| `JWT_SECRET` | long random string |
| `UPLOAD_DIR` | `/app/uploads` (Docker) or `../uploads` |
| `PUBLIC_MEDIA_URL` | `https://api.yoursite.com/uploads` |

### Frontend (`frontend/.env.local` or build args)

| Variable | Example |
|----------|---------|
| `NEXT_PUBLIC_API_URL` | `https://api.yoursite.com` |
| `NEXT_PUBLIC_MEDIA_URL` | `https://api.yoursite.com/uploads` |
| `NEXT_PUBLIC_SITE_URL` | `https://suryaphotography.com` |

## Option A — Windows Server (IIS / manual)

1. Build backend: `cd backend && go build -o bin/server.exe ./cmd/server`
2. Build frontend: `cd frontend && npm run build && npm start`
3. Run API as Windows Service or via `start-backend.bat`
4. Reverse-proxy `/uploads` and API to port 8080; site to port 3000
5. Persist `uploads/` folder on disk

## Option B — Docker

```bash
# Copy and edit env files first
cp backend/.env.example backend/.env
# Edit backend/.env with production SQL + secrets

export NEXT_PUBLIC_API_URL=https://api.yoursite.com
export NEXT_PUBLIC_MEDIA_URL=https://api.yoursite.com/uploads
export NEXT_PUBLIC_SITE_URL=https://yoursite.com

docker compose up -d --build
```

SQL Server must be reachable from the `api` container (not included in compose).

## Option C — AWS (recommended pattern)

| Component | Service |
|-----------|---------|
| Frontend | Amplify / ECS / EC2 + Node |
| API | ECS / EC2 + Go binary |
| Database | RDS SQL Server or existing SSMS host |
| Media | S3 + CloudFront (replace `LocalStorage` later) |

Uploads currently use local disk; migrate `internal/storage` to S3 when scaling.

## Post-deploy verification

1. `GET https://api.yoursite.com/health` → `{"status":"ok"}`
2. Public site loads home + portfolio
3. Contact form saves inquiry
4. Admin login at `/admin/login`
5. Upload test image in admin portfolio
6. Swagger: `/swagger/index.html` (restrict in production firewall if needed)

## Security notes

- Use HTTPS everywhere
- Restrict SQL Server firewall to API host only
- Remove or protect Swagger in production
- Rotate `JWT_SECRET` if compromised
- Back up `uploads/` and database regularly
