import { createWebHistory, createRouter } from 'vue-router'

import Home from '@/pages/home/ui/Home.vue'
import Auth from '@/pages/auth/ui/Auth.vue'

const routes = [
  {path: "/", component: Home},
  {path: "/auth", component: Auth},
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})