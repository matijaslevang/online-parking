import { UserManager } from 'oidc-client-ts'

const appUrl = window.location.origin + '/'

export const userManager = new UserManager({
  authority: import.meta.env.VITE_KEYCLOAK_URL,
  client_id: import.meta.env.VITE_KEYCLOAK_CLIENT_ID,
  redirect_uri: appUrl,
  post_logout_redirect_uri: appUrl,
  scope: 'openid profile email',
  // Background renewal would keep the session alive while the user is idle. Renewal is planned in docs/todo.md.
  automaticSilentRenew: false,
})

export function removeLoginParamsFromUrl() {
  window.history.replaceState({}, document.title, window.location.pathname)
}
