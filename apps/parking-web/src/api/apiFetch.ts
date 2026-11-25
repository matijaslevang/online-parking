import { userManager } from '../auth/userManager'

type ApiFetchOptions = {
  method?: string
  body?: unknown
}

// All backend calls go through this function. It adds the access token and handles errors in one place.
export async function apiFetch<T>(path: string, options: ApiFetchOptions = {}): Promise<T> {
  const user = await userManager.getUser()
  if (!user || user.expired) {
    // Without token renewal (docs/todo.md) an expired token needs a new login.
    // Keycloak skips the login form while its own session is still active.
    await userManager.signinRedirect()
    throw new Error('Not authenticated')
  }

  const response = await fetch(`/api${path}`, {
    method: options.method ?? 'GET',
    headers: {
      Authorization: `Bearer ${user.access_token}`,
      'Content-Type': 'application/json',
    },
    body: options.body === undefined ? undefined : JSON.stringify(options.body),
  })

  if (response.status === 401) {
    await userManager.signinRedirect()
    throw new Error('Not authenticated')
  }

  if (!response.ok) {
    const errorBody = await response.json().catch(() => null)
    throw new Error(errorBody?.error ?? `Request failed with status ${response.status}`)
  }

  if (response.status === 204) {
    return undefined as T
  }
  return response.json()
}
