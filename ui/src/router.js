/**
 * Sets up the Vue Router with basic routes for create, view, and error pages.
 */
import { createRouter, createWebHistory } from 'vue-router'
import Create from './pages/Create.vue'
import View from './pages/View.vue'
import Error404 from './pages/Error404.vue'
import ErrorGeneral from './pages/ErrorGeneral.vue'

const routes = [
  { path: '/', component: Create },
  { path: '/view/:hash', component: View },
  { path: '/error', component: ErrorGeneral },
  { path: '/:pathMatch(.*)*', component: Error404 }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
