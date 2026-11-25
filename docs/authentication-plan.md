# Authentication Plan

Keycloak is the identity provider. Each system has its own realm: `parking`, `psp`, `bank`.
PCI DSS decisions are recorded in [pci-dss.md](pci-dss.md).

## Users and roles

| System | Users | Roles | Login |
| --- | --- | --- | --- |
| Parking | Drivers, self-registration | `user` | Password |
| PSP | Super admin, merchants created by super admin | `super-admin`, `merchant` | Password + TOTP, password change on first login |
| Bank mobile app | Bank customers | `customer` | Password |
| Bank payment page (`bank-web`) | No login | — | Opened via one-time `PAYMENT_URL` |

## Decisions

- The bank has no admin account. The specification does not require one.
- MFA is implemented with TOTP (Google or Microsoft Authenticator), as the specification requires for requirement 8.
- The specification does not say which users need MFA. We require it for PSP users (super admin and merchants). Bank customers log in with a password only.
- `bank-web` has no login. It is the card payment page.
- PSP merchant registration is part of this work.

## Steps

1. **Keycloak realm configuration** (`infra/keycloak/realms`)
   - Roles per realm and one development user per role.
   - Password policy, brute-force lockout, session timeouts, short-lived access tokens.
   - Login and admin events enabled.
   - PSP realm: TOTP required, password change on first login for new merchants.
   - Bank realm: replace the `bank-web` client with a `bank-mobile` public client with app redirect URI.
   - PSP realm: service account client used by the PSP backend to create merchant users.

2. **Token validation inside Docker**
   - Tokens are issued for `https://auth.localhost`, which containers cannot reach.
   - Backends fetch signing keys from `http://keycloak:8080` but still validate the public issuer.
   - New config: `KEYCLOAK_ISSUER_URL`, `KEYCLOAK_JWKS_URL`.

3. **Go backends** (`internal/auth` in each service)
   - Middleware that validates the token signature, issuer, audience and expiry (`github.com/coreos/go-oidc/v3`).
   - Role check, e.g. `RequireRole("super-admin")`.
   - `GET /api/me` returns the logged-in user and roles.

4. **Web apps** (`parking-web`, `psp-web`; `bank-web` stays public)
   - `react-oidc-context`, authorization code flow with PKCE.
   - Login, logout, protected pages, bearer token on API calls, show user name and roles.

5. **Mobile app** (`bank-mobile`)
   - `expo-auth-session` with PKCE, token stored with `expo-secure-store`, logout.
   - The device cannot reach `*.localhost` or trust the self-signed certificate. Testing needs the computer's LAN IP and a trusted certificate (mkcert root CA on the device).

6. **PSP merchant registration**
   - Super admin creates a merchant in `psp-web`.
   - The PSP backend saves the merchant in its database and creates the Keycloak user with the `merchant` role, a temporary password and required actions: update password, configure TOTP.
   - On first login the merchant user must change the password and set up the authenticator app.

7. **Verification**
   - Every login flow works. Requests without a token or with a wrong role are rejected.

## Out of scope (done together with payment flows)

- `MERCHANT_ID` + `API_KEY` authentication of the parking system towards the PSP.
- HMAC or certificate authentication between the PSP and the bank.
- Signed status callbacks from the PSP to the parking system.
