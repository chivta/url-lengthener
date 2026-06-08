import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

const apiPort = process.env.PORT || '8080'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      '/api': `http://api:${apiPort}`,
    },
  },
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/test-setup.ts'],
    globals: true,
  },
})
