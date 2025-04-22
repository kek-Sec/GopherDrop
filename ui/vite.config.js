import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vuetify from 'vite-plugin-vuetify'

// Create HTML plugin instance with template parameters
const htmlPlugin = () => {
  return {
    name: 'html-transform',
    transformIndexHtml(html) {
      return html.replace(
        /__TITLE__/g, 
        process.env.VITE_APP_TITLE || 'GopherDrop'
      ).replace(
        /__DESCRIPTION__/g,
        process.env.VITE_APP_DESCRIPTION || 'Secure one-time secret and file sharing'
      )
    }
  }
}

export default defineConfig({
  plugins: [vue(), vuetify(), htmlPlugin()],
  define: {
    'process.env.VITE_APP_TITLE': JSON.stringify(process.env.VITE_APP_TITLE || 'GopherDrop'),
    'process.env.VITE_APP_DESCRIPTION': JSON.stringify(process.env.VITE_APP_DESCRIPTION || 'Secure one-time secret and file sharing')
  }
})
