# Online Parking and Electronic Payment System

University project (SEP 2026/2027): a parking application, a Payment Service Provider (PSP) and a simulated bank. Requirements: `SEP - Specifikacija 2026_2027.pdf`.

## Structure

```text
services/          Go backends (each an independent Go module, GORM + PostgreSQL)
  parking/         Parking system API
  psp/             Payment Service Provider API
  bank/            Bank API
apps/
  parking-web/     React + TypeScript + Vite + Tailwind
  psp-web/         React + TypeScript + Vite + Tailwind
  bank-web/        React + TypeScript + Vite + Tailwind
  bank-mobile/     React Native (Expo) bank app for QR payments
infra/
  traefik/         Reverse proxy, HTTPS and load balancer config
  certs/           Development certificate script
  postgres/        Database init script (one database per service)
  keycloak/        Realm imports (one realm per system)
docs/              Project documentation
```

## Requirements

- Docker with Docker Compose

## Run everything with Docker

```sh
cp .env.example .env
./infra/certs/generate.sh
docker compose up --build
```

| URL | Description |
| --- | --- |
| https://parking.localhost | Parking web app (`/api` → parking service) |
| https://psp.localhost | PSP web app (`/api` → PSP service, 2 instances) |
| https://bank.localhost | Bank web app (`/api` → bank service, 2 instances) |
| https://auth.localhost | Keycloak (admin credentials in `.env`) |
| https://traefik.localhost | Traefik dashboard |

Health checks: `https://<system>.localhost/api/health`. The `instance` field shows which replica answered.

Without mkcert the certificate is self-signed, so the browser shows a warning and `curl` needs `-k`.

Change the number of instances with `docker compose up --scale psp=3 --scale bank=3`.

## Run a single service or app

Start only the database (and Keycloak if needed):

```sh
docker compose up postgres keycloak traefik
```

Backend service (example: parking, port 8081; psp uses 8082, bank 8083):

```sh
cd services/parking
cp .env.example .env
set -a; . ./.env; set +a
go run .
```

Web app (example: parking-web on port 5173; psp-web 5174, bank-web 5175). `/api` is proxied to the local service:

```sh
cd apps/parking-web
npm install
npm run dev
```

Mobile app:

```sh
cd apps/bank-mobile
npm install
npm start
```
