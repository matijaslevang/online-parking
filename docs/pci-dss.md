# PCI DSS v4.0.1 Decisions

Requirements the project must cover (from the project specification): 3, 4, 7, 8 (MFA with Google or Microsoft Authenticator) and 10.

Status: **Planned** = decided, not implemented yet. **Done** = implemented.

Requirement numbers must match the official PCI DSS v4.0.1 document. Check each new entry against it.

## Requirement 4 — Protect cardholder data with strong cryptography during transmission

| Requirement | Decision | Status |
| --- | --- | --- |
| 4.2.1 | All browser traffic goes over HTTPS through Traefik. HTTP is redirected to HTTPS. | Done |

## Requirement 7 — Restrict access by business need to know

| Requirement | Decision | Status |
| --- | --- | --- |
| 7.2.1 | Each system has its own Keycloak realm with its own roles. | Done |
| 7.2.2 | Users get only the role they need (`user`, `merchant`, `super-admin`, `customer`). Backends check the role on every protected endpoint. | Done |

## Requirement 8 — Identify users and authenticate access

| Requirement | Decision | Status |
| --- | --- | --- |
| 8.2.1 | Every user has a unique Keycloak account. No shared accounts. | Done |
| 8.2.8 | Users must log in again after 15 minutes of inactivity. | Done |
| 8.3.1 | Every request to a protected endpoint needs a valid access token. Backends check its signature, issuer, audience and expiry. | Done |
| 8.3.2 | Passwords are hashed by Keycloak and sent only over HTTPS. | Done |
| 8.3.4 | The account is locked after 10 failed logins, for 30 minutes. | Done |
| 8.3.5 | New PSP merchants get a unique random temporary password and must change it on first login. The PSP does not store it. | Done |
| 8.3.6 | Passwords have at least 12 characters, with letters and digits. | Done |
| 8.3.7 | A new password cannot be one of the last 4 passwords. | Done |
| 8.4.1, 8.4.2 | PSP users must use TOTP (Google or Microsoft Authenticator) in addition to the password. New merchants must set it up on first login. | Done |
| 8.5.1 | A TOTP code can be used only once, so a stolen code cannot be replayed. | Done |

## Requirement 10 — Log and monitor all access

| Requirement | Decision | Status |
| --- | --- | --- |
| 10.2.1 | Keycloak login and admin events are enabled in all realms. | Done |
| 10.2.1.4 | Failed login attempts are recorded as Keycloak events. | Done |
| 10.2.1.5 | Password and MFA changes are recorded as Keycloak events. | Done |
| 10.2.2 | Merchant registration records which super admin created the merchant (in the PSP database and log), because Keycloak only sees the PSP service account. | Done |
| 10.5.1 | Keycloak events are kept for 12 months. | Done |

## Requirement 3 — Protect stored account data

No decisions yet. Will be added with card payment.
