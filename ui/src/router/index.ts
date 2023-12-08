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
    path: '/home',
    name: 'Home',
    component: () => import('../views/home/index.vue')
  },
  {
    path: '/',
    name: 'Index',
    redirect: location.href.includes("?url=")?'/bookmark':"/home"
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
