import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/bookmark',
    name: 'Bookmark',
    component: ()=>import("../views/bookmark/index.vue")
  },
  {
    path: '/note',
    name: 'Note',
    component: ()=>import("../views/note/index.vue")
  },
  {
    path: '/',
    name: 'Index',
    component: ()=>import("../views/index.vue")
  },
  {
    path: '/login',
    name:  'Login',
    component: ()=>import("../views/login/index.vue")
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
