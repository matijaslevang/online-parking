# TODO

Postponed work, to be done later.

## Authentication

- [ ] **Access token renewal in web apps** (`parking-web`, `psp-web`)
  - Now: access tokens expire after 5 minutes. `apiFetch` then redirects to Keycloak, which logs the user back in without the form while the Keycloak session is active. The page reloads and its state is lost.
  - Plan: in `apiFetch`, renew an expired token with the refresh token (`userManager.signinSilent()`) before the request. Redirect to login only if renewal fails.
  - Keep `automaticSilentRenew: false`. Background renewal would keep the session alive while the user is idle and break the 15-minute idle timeout (PCI DSS 8.2.8). Renewing only on API calls means only real activity keeps the session alive.
