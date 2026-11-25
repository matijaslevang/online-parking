import { useEffect, useState } from 'react'
import { useAuth } from 'react-oidc-context'
import { apiFetch } from './api/apiFetch'
import { MerchantsPage } from './merchants/MerchantsPage'

type CurrentUser = {
  id: string
  username: string
  email: string
  name: string
  roles: string[]
}

function App() {
  const auth = useAuth()
  const [currentUser, setCurrentUser] = useState<CurrentUser | null>(null)
  const [errorMessage, setErrorMessage] = useState('')

  useEffect(() => {
    apiFetch<CurrentUser>('/me')
      .then(setCurrentUser)
      .catch((error: Error) => setErrorMessage(error.message))
  }, [])

  const isSuperAdmin = currentUser?.roles.includes('super-admin') ?? false

  return (
    <main className="flex min-h-screen flex-col items-center gap-6 bg-slate-50 p-8 text-slate-800">
      <header className="flex w-full max-w-3xl items-center justify-between">
        <h1 className="text-2xl font-semibold">Payment Service Provider</h1>
        <div className="flex items-center gap-4 text-sm">
          {currentUser && (
            <span>
              {currentUser.name} ({currentUser.username})
            </span>
          )}
          <button
            className="rounded bg-slate-800 px-4 py-2 text-white hover:bg-slate-700"
            onClick={() => auth.signoutRedirect()}
          >
            Log out
          </button>
        </div>
      </header>
      {errorMessage && <p className="text-sm text-red-600">{errorMessage}</p>}
      {isSuperAdmin && <MerchantsPage />}
    </main>
  )
}

export default App
