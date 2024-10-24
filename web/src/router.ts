import { createRouter, createWebHistory } from 'vue-router'

import PersonalizedView from '@/views/PersonalizedView.vue'
import SearchView from '@/views/SearchView.vue'

const routes = [
  { path: '/', component: PersonalizedView },
  { path: '/search', component: SearchView },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})
