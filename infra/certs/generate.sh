#!/bin/sh
# Generates a local development certificate for all *.localhost domains.
# Uses mkcert when installed (trusted by the browser), otherwise falls back to a self-signed certificate.
set -e

cd "$(dirname "$0")"

DOMAINS="localhost parking.localhost psp.localhost bank.localhost auth.localhost traefik.localhost"

if command -v mkcert >/dev/null 2>&1; then
  mkcert -install
  mkcert -cert-file local-cert.pem -key-file local-key.pem $DOMAINS
else
  echo "mkcert not found, generating a self-signed certificate (the browser will show a warning)."
  SAN=""
  for domain in $DOMAINS; do
    SAN="${SAN}DNS:${domain},"
  done
  openssl req -x509 -newkey rsa:2048 -nodes -days 825 \
    -keyout local-key.pem -out local-cert.pem \
    -subj "/CN=localhost" \
    -addext "subjectAltName=${SAN%,}"
fi

echo "Certificate written to infra/certs/local-cert.pem"
