import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// Proxy target: inside Docker the backend is reachable as "backend",
// but when running the dev server on the host it must be "localhost".
// VITE_BACKEND_URL lets docker-compose override the host default.
const backendTarget = process.env.VITE_BACKEND_URL || 'http://localhost:8080'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/auth': backendTarget,
      '/users': backendTarget,
      '/me': backendTarget,
      '/recommendations': backendTarget,
      '/connections': backendTarget,
      '/chats': backendTarget,
      '/ws': {
        target: backendTarget,
        ws: true,
      },
    },
  },
})