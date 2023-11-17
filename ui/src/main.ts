import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import pinia from '@/store'
import lazyPlugin from 'vue3-lazy'

import 'virtual:svg-icons-register'
import '@/router/guard/index'
import 'ant-design-vue/dist/reset.css'
import 'virtual:uno.css'
import './style.css'

const app = createApp(App)

app.use(lazyPlugin, {})
app.use(router)
app.use(pinia)

app.mount('#app')
