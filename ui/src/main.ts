import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'
import pinia from '@/store'
import 'virtual:svg-icons-register'
import '@/router/guard/index'

import 'ant-design-vue/es/message/style/css'
// @ts-ignore
import contextmenu from 'v-contextmenu'
import 'v-contextmenu/dist/themes/default.css'

import lazyPlugin from 'vue3-lazy'

const app = createApp(App)
app.use(contextmenu)
app.use(lazyPlugin, {})
app.use(router)
app.use(pinia)
app.mount('#app')
