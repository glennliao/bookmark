import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'
import pinia from '@/store'
import 'virtual:svg-icons-register'
import '@/router/guard/index'
import {
  Drawer,
  Button,
  Tree,
  Dropdown,
  Menu,
  Modal,
  Form,
  Alert,
  Input,
  Card,
  TreeSelect,
  Popover,
  Tabs,
  Empty,
  FloatButton
} from 'ant-design-vue'
import 'ant-design-vue/dist/reset.css'

// @ts-ignore
import contextmenu from 'v-contextmenu'
import 'v-contextmenu/dist/themes/default.css'

import lazyPlugin from 'vue3-lazy'

const app = createApp(App)
app.use(contextmenu)
app.use(lazyPlugin, {})
app.use(router)
app.use(pinia)
app.use(Drawer).use(Button).use(Tree)
  .use(Dropdown).use(Menu).use(Modal)
  .use(Form).use(Alert).use(Input)
  .use(Card).use(TreeSelect).use(Popover).use(Tabs).use(FloatButton).use(Empty)
app.mount('#app')
