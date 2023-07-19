import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/bookmark',
    name: 'Bookmark',
    component: () => import('../views/bookmark/index.vue')
  },
  {
    path: '/note',
    name: 'Note',
    component: () => import('../views/note/index.vue')
  },
  {
    path: '/',
    name: 'Index',
    redirect: '/bookmark'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/login/index.vue')
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
