import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  base: '/admin/',
  build: {
    outDir: 'dist',
    emptyOutDir: true,
    rollupOptions: {
      output: {
        manualChunks: undefined,
      }
    }
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      },
      '/admin/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      },
      '/admin/login': {
        target: 'http://localhost:8080',
        changeOrigin: true
      },
      '/admin/logout': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})