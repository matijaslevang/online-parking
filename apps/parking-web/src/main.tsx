import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { AuthProvider } from 'react-oidc-context'
import './index.css'
import App from './App.tsx'
import { RequireAuth } from './auth/RequireAuth.tsx'
import { removeLoginParamsFromUrl, userManager } from './auth/userManager.ts'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <AuthProvider userManager={userManager} onSigninCallback={removeLoginParamsFromUrl}>
      <RequireAuth roles={['user']}>
        <App />
      </RequireAuth>
    </AuthProvider>
  </StrictMode>,
)
