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
- [mkcert](https://github.com/FiloSottile/mkcert), for a trusted local HTTPS certificate

## First-time setup: trusted HTTPS certificate

Login does not work without this step. The browser blocks requests to Keycloak (`https://auth.localhost`) when the certificate is not trusted, and the app shows "Login failed: Failed to fetch".

1. Install mkcert:

   | OS | Command |
   | --- | --- |
   | macOS | `brew install mkcert` (and `brew install nss` for Firefox) |
   | Windows | `choco install mkcert` or `scoop bucket add extras && scoop install mkcert` |
   | Linux (Debian/Ubuntu) | `sudo apt install mkcert libnss3-tools` |

2. Generate the certificate (asks for your password once, to trust mkcert's local certificate authority):

   ```sh
   ./infra/certs/generate.sh
   ```

   On Windows, run it in Git Bash or WSL, or run these commands in `infra/certs`:

   ```sh
   mkcert -install
   mkcert -cert-file local-cert.pem -key-file local-key.pem localhost parking.localhost psp.localhost bank.localhost auth.localhost traefik.localhost
   ```

3. Restart the browser. If Docker is already running, also run `docker compose restart traefik`.

Every team member does this on their own computer. Never share the mkcert private key (`rootCA-key.pem` in the folder shown by `mkcert -CAROOT`).

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

Without mkcert, `generate.sh` creates a self-signed certificate. The browser then shows warnings and login fails (see First-time setup).

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
