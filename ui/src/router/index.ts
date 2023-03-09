import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Index',
    component: ()=>import("../views/index.vue")
  },
  {
    path: '/login',
    name:  'Login',
    component: ()=>import("../views/login.vue")
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
