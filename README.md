# Surya Photography

Premium cinematic photography portfolio and internal studio management system.

## What's included

| Layer | Stack |
|-------|--------|
| **Public site** | Home, portfolio (9 categories), about, contact, cinematic UI |
| **Admin** | Dashboard, inquiries, functions tracker, media uploads |
| **API** | Go Fiber, JWT, Swagger, JSON file storage |
| **Storage** | JSON files (`backend/data/`) + local uploads (`uploads/`) |

## Quick start (local)

### 1. Backend

```bat
cd backend
copy .env.example .env
go run ./cmd/server
```

- API: http://localhost:8080  
- Swagger: http://localhost:8080/swagger/index.html  
- Uploads: http://localhost:8080/uploads  

### 2. Frontend

```bat
cd frontend
npm install
npm run dev
```

If PowerShell blocks `npm`/`next` scripts, run `start-frontend.bat` (recommended) or use `npm.cmd run dev`.

- Site: http://localhost:3000  
- Admin: http://localhost:3000/admin/login  

**Default admin:** `surya@admin.com` / `surya@1995`

Or use root scripts: `start-backend.bat` and `start-frontend.bat`.

## Production deploy

See **[DEPLOY.md](DEPLOY.md)** for:

- Environment variables
- Docker Compose
- Windows / AWS deployment
- Security checklist

```bash
docker compose up -d --build
```

## Project structure

```
backend/           Go Fiber API
frontend/          Next.js 16 app
uploads/           Media files (photos/videos per category)
backend/data/      JSON data files (admins, categories, media, inquiries, functions, settings)
docker-compose.yml Optional production stack
```

## Features

- Cinematic hero slideshow, masonry gallery, fullscreen lightbox
- Category pages with auto-hidden video section when empty
- Contact inquiries + optional SMTP email to admin
- Single admin (JWT + bcrypt, rate-limited login)
- Function/production tracker with payment and delivery status
- Local uploads with storage abstraction for future S3 migration

## Build verification

```bat
cd backend && go build ./cmd/server
cd frontend && npm.cmd run build
```

Both must succeed before deploy.
