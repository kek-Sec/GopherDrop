/**
 * Main application entry point.
 * Initializes Vue, Vuetify, and the Router.
 */
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// Import Vuetify and styles
import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import '@mdi/font/css/materialdesignicons.css'

// Create a custom light theme
const customLightTheme = {
  dark: false,
  colors: {
    primary: '#6750A4', // A modern purple
    secondary: '#E8DEF8', // A light, complementary purple
    background: '#FFFFFF',
    surface: '#FEF7FF', // Off-white for cards and surfaces
    info: '#6750A4', // Using primary color for info state for consistency
    error: '#B3261E',
    success: '#4CAF50',
    warning: '#FB8C00',
  }
};

// Create a custom dark theme
const customDarkTheme = {
  dark: true,
  colors: {
    primary: '#D0BCFF', // A lighter purple for dark mode contrast
    secondary: '#4A4458',
    background: '#141218',
    surface: '#262329',
    info: '#D0BCFF', // Lighter purple for info state
    error: '#F2B8B5',
    success: '#B7F3B9',
    warning: '#FFD6A8',
  }
}

// Create Vuetify instance
const vuetify = createVuetify({
  theme: {
    defaultTheme: 'customLightTheme',
    themes: {
      customLightTheme,
      customDarkTheme
    },
  },
  icons: {
    defaultSet: 'mdi',
  },
})

// Create and mount the Vue application
const app = createApp(App)
app.use(router)
app.use(vuetify)
app.mount('#app')