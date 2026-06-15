# Deployment Guide — Surya Photography

## Pre-deploy checklist

- [ ] Change admin password if needed (default `surya@1995`)
- [ ] Set strong `JWT_SECRET` (32+ random characters)
- [ ] Configure `SMTP_*` for inquiry email notifications
- [ ] Set real contact URLs in backend `.env`
- [ ] Point `PUBLIC_MEDIA_URL` to your public API uploads URL
- [ ] Set `CORS_ORIGINS` and `FRONTEND_URL` to your live domain
- [ ] Persist `backend/data/` and `uploads/` on disk (JSON + media)

## Environment variables

### Backend (`backend/.env`)

| Variable | Example |
|----------|---------|
| `PORT` | `8080` |
| `APP_ENV` | `production` |
| `FRONTEND_URL` | `https://suryaphotography.com` |
| `CORS_ORIGINS` | `https://suryaphotography.com` |
| `JWT_SECRET` | long random string |
| `DATA_DIR` | `/app/data` (Docker) or `data` (local) |
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
# Edit backend/.env with production secrets

export NEXT_PUBLIC_API_URL=https://api.yoursite.com
export NEXT_PUBLIC_MEDIA_URL=https://api.yoursite.com/uploads
export NEXT_PUBLIC_SITE_URL=https://yoursite.com

docker compose up -d --build
```

## Option C — AWS EC2 (t3.micro VM)

Recommended setup: run both the Go API and Next.js app with Docker Compose, and persist JSON + uploads on the EC2 disk.

1. Create an EC2 instance (t3.micro), attach an EBS volume (or use the root volume) with enough space for `uploads/`.
2. Security Group inbound:
   - 22 (SSH) from your IP
   - 80 (HTTP) and 443 (HTTPS) from the internet (or from a load balancer)
3. Install Docker and Docker Compose on the instance.
4. Copy this repo to the instance (git clone or upload a zip).
5. Create and edit `backend/.env` (copy from `.env.example`) and set:
   - `APP_ENV=production`
   - `JWT_SECRET=...`
   - `FRONTEND_URL=https://your-domain`
   - `CORS_ORIGINS=https://your-domain`
   - `PUBLIC_MEDIA_URL=https://api.your-domain/uploads` (or your reverse-proxy path)
6. Start:

```bash
docker compose up -d --build
```

Persisted data:
- `backend/data/` stores JSON (admins, settings, categories, media, inquiries, functions)
- `uploads/` stores uploaded files

## Post-deploy verification

1. `GET https://api.yoursite.com/health` → `{"status":"ok"}`
2. Public site loads home + portfolio
3. Contact form saves inquiry
4. Admin login at `/admin/login`
5. Upload test image in admin portfolio
6. Swagger: `/swagger/index.html` (restrict in production firewall if needed)

## Security notes

- Use HTTPS everywhere
- Remove or protect Swagger in production
- Rotate `JWT_SECRET` if compromised
- Back up `backend/data/` and `uploads/` regularly
