import { createRouter, createWebHistory } from 'vue-router'
import UploadView from '../views/UploadView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'upload',
      component: UploadView,
    },
    {
      path: '/process/:id',
      name: 'Process',
      component: () => import('../views/ProcessPage.vue'), // Lazy-loaded
    },
    {
      path: '/download/:id',
      name: 'Download',
      component: () => import('../views/DownloadPage.vue'), // Lazy-loaded
    },
  ],
})

export default router
