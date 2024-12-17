import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// Import Vuetify and styles
import 'vuetify/styles' 
import { createVuetify } from 'vuetify'
import '@mdi/font/css/materialdesignicons.css' 

const vuetify = createVuetify({
  icons: {
    defaultSet: 'mdi', 
  },
})

const app = createApp(App)
app.use(router)
app.use(vuetify)
app.mount('#app')
