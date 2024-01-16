import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'
import Layout from '../layout/index.vue'
const routes: Array<RouteRecordRaw> = [

  {
    path: '/',
    name: 'Index',
    redirect: location.href.includes("?url=")?'/bookmark':"/home",
    component: Layout,
    children:[
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
        name: 'Desktop',
        component: () => import('../pages/desktop/index.vue')
      },
    ]
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
