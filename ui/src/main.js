/**
 * Main application entry point.
 * Initializes Vue, Vuetify, and the Router.
 */
import { createApp } from 'vue'
import App from './App.vue'
import router from './router' // This now imports from the /router directory

// Import Vuetify and styles
import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import '@mdi/font/css/materialdesignicons.css'

// Create Vuetify instance
const vuetify = createVuetify({
  icons: {
    defaultSet: 'mdi',
  },
})

// Create and mount the Vue application
const app = createApp(App)
app.use(router)
app.use(vuetify)
app.mount('#app')
