## Surya Photography Deployment

### Prerequisites
- Docker Desktop (or Docker Engine + Docker Compose v2)
- A domain (optional but recommended) and HTTPS handled by your reverse proxy / hosting platform

### 1) Configure environment variables
Create a file named `.env` in the project root (same folder as `docker-compose.yml`).

Start from the template:
- Copy `.env.example` to `.env`

Set at minimum:
- `JWT_SECRET` to a long random secret
- `NEXT_PUBLIC_SITE_URL` to your public website URL (example: `https://example.com`)
- `NEXT_PUBLIC_API_URL` to your public API URL (example: `https://api.example.com` or `https://example.com/api` if you proxy it)
- `PUBLIC_MEDIA_URL` / `NEXT_PUBLIC_MEDIA_URL` to your public uploads URL (example: `https://api.example.com/uploads`)
- `FRONTEND_URL` / `CORS_ORIGINS` to your public website URL

### 2) Start in production mode
From the project root:

```bash
docker compose -f docker-compose.prod.yml up -d --build
```

If you are deploying locally (single server, localhost URLs), you can use:

```bash
docker compose up -d --build
```

### 3) Verify
- Website: `http://SERVER_IP:3000`
- API health: `http://SERVER_IP:8080/health`
- Admin: `http://SERVER_IP:3000/admin/login`

### 4) Data persistence and backups
Docker volumes are used for:
- `uploads` (uploaded images/files)
- `data` (backend JSON storage)

Back up these volumes regularly (your hosting provider’s snapshot tools or `docker run --rm -v <volume>:/data ...`).

### 5) Reverse proxy (recommended)
For real production (HTTPS + domain), run a reverse proxy and route:
- `/` to the frontend container port `3000`
- `/uploads` and `/health` and `/api/*` to the backend container port `8080`

This lets you keep only one public domain and avoids exposing the API port directly.
