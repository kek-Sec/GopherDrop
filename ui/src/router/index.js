/**
 * Sets up the Vue Router.
 * The import paths are now relative to the `src/router/` directory.
 */
import { createRouter, createWebHistory } from 'vue-router';
import Create from '../pages/Create.vue'; // Corrected Path
import View from '../pages/View.vue'; // Corrected Path
import Error404 from '../pages/Error404.vue'; // Corrected Path
import ErrorGeneral from '../pages/ErrorGeneral.vue'; // Corrected Path

const routes = [
  { path: '/', component: Create, name: 'create' },
  { path: '/view/:hash', component: View, name: 'view' },
  { path: '/error', component: ErrorGeneral, name: 'error' },
  { path: '/:pathMatch(.*)*', component: Error404, name: 'not-found' }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

export default router;
