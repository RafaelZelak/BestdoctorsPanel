// vite.config.js
import vue from '@vitejs/plugin-vue'
import path from 'path'
import { defineConfig } from 'vite'

export default defineConfig({
  // Remove base path for independent operation
  base: '/',
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    host: '0.0.0.0', // Allow access from Docker network
    port: 5173,
    proxy: {
      '/bestdoctors': {
        target: process.env.VITE_API_PROXY || 'http://backend:9002',
        changeOrigin: true,
        secure: false,
      },
    },
  },
  preview: {
    host: '0.0.0.0',
    port: 5173,
  },
})
