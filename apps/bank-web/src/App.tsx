import { useEffect, useState } from 'react'

function App() {
  const [apiStatus, setApiStatus] = useState('checking...')

  useEffect(() => {
    fetch('/api/health')
      .then((response) => response.json())
      .then((health) => setApiStatus(health.status))
      .catch(() => setApiStatus('unreachable'))
  }, [])

  return (
    <main className="flex min-h-screen flex-col items-center justify-center gap-2 bg-slate-50 text-slate-800">
      <h1 className="text-3xl font-semibold">Bank</h1>
      <p className="text-sm text-slate-500">API status: {apiStatus}</p>
    </main>
  )
}

export default App
