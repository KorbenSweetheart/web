import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/auth': 'http://backend:8080',
      '/users': 'http://backend:8080',
      '/me': 'http://backend:8080',
      '/recommendations': 'http://backend:8080',
      '/connections': 'http://backend:8080',
      '/chats': 'http://backend:8080',
      '/ws': {
        target: 'http://backend:8080',
        ws: true,
      },
    },
  },
})