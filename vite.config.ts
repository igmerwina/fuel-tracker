import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { cpSync } from 'fs'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    {
      name: 'copy-static-assets',
      writeBundle() {
        try {
          cpSync('node_modules/leaflet/dist/images', 'dist/images', { recursive: true, force: true })
          cpSync('node_modules/@fortawesome/fontawesome-free/webfonts', 'dist/webfonts', { recursive: true, force: true })
        } catch {
          // ignore copy errors in environments without fs access
        }
      },
    },
  ],
})
