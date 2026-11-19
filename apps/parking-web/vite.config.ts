import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { defineConfig, loadEnv } from 'vite'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')

  return {
    plugins: [react(), tailwindcss()],
    server: {
      host: true,
      port: 5173,
      // Used only when running outside Docker; behind Traefik /api is routed directly to the service.
      proxy: {
        '/api': env.API_PROXY_TARGET || 'http://localhost:8081',
      },
    },
  }
})
