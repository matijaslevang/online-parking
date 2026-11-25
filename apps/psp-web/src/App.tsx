import { useEffect, useState } from 'react'
import { useAuth } from 'react-oidc-context'
import { apiFetch } from './api/apiFetch'

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

  return (
    <main className="flex min-h-screen flex-col items-center justify-center gap-2 bg-slate-50 text-slate-800">
      <h1 className="text-3xl font-semibold">Payment Service Provider</h1>
      {currentUser && (
        <>
          <p>
            Logged in as {currentUser.name} ({currentUser.username})
          </p>
          <p className="text-sm text-slate-500">Roles: {currentUser.roles.join(', ')}</p>
        </>
      )}
      {errorMessage && <p className="text-sm text-red-600">{errorMessage}</p>}
      <button
        className="mt-4 rounded bg-slate-800 px-4 py-2 text-sm text-white hover:bg-slate-700"
        onClick={() => auth.signoutRedirect()}
      >
        Log out
      </button>
    </main>
  )
}

export default App
