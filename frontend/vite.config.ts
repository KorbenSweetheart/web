import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
     proxy: {
      '/auth': 'http://localhost:8080',
      '/users': 'http://localhost:8080',
      '/me': 'http://localhost:8080',
      '/recommendations': 'http://localhost:8080',
      '/connections': 'http://localhost:8080',
    },
  },
})