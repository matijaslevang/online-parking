import { useEffect, useRef, type ReactNode } from 'react'
import { hasAuthParams, useAuth } from 'react-oidc-context'

type RequireAuthProps = {
  children: ReactNode
  // The user needs at least one of these roles. The backend checks roles again on every request.
  roles: string[]
}

export function RequireAuth({ children, roles }: RequireAuthProps) {
  const auth = useAuth()
  const hasTriedLogin = useRef(false)

  const shouldRedirectToLogin =
    !auth.isLoading && !auth.isAuthenticated && !auth.activeNavigator && !auth.error && !hasAuthParams()

  useEffect(() => {
    if (shouldRedirectToLogin && !hasTriedLogin.current) {
      hasTriedLogin.current = true
      auth.signinRedirect()
    }
  }, [shouldRedirectToLogin, auth])

  if (auth.error) {
    return <StatusMessage text={`Login failed: ${auth.error.message}`} />
  }
  if (!auth.user || !auth.isAuthenticated) {
    return <StatusMessage text="Loading..." />
  }

  const userRoles = getRealmRoles(auth.user.access_token)
  const hasRequiredRole = roles.some((role) => userRoles.includes(role))
  if (!hasRequiredRole) {
    return <StatusMessage text="Access denied" />
  }

  return children
}

function getRealmRoles(accessToken: string): string[] {
  const payload = accessToken.split('.')[1]
  const payloadJson = atob(payload.replace(/-/g, '+').replace(/_/g, '/'))
  return JSON.parse(payloadJson).realm_access?.roles ?? []
}

function StatusMessage({ text }: { text: string }) {
  return <main className="flex min-h-screen items-center justify-center text-slate-500">{text}</main>
}
