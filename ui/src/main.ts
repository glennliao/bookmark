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
  FloatButton,
  Upload,
  Breadcrumb,
  Badge,
  Tag, Avatar, Segmented,
  Layout, Row, Col,Pagination,Select,
} from 'ant-design-vue'
import 'ant-design-vue/dist/reset.css'

import lazyPlugin from 'vue3-lazy'

const app = createApp(App)

app.use(lazyPlugin, {})
app.use(router)
app.use(pinia)

app.use(Drawer).use(Button).use(Tree)
  .use(Dropdown).use(Menu).use(Modal)
  .use(Form).use(Alert).use(Input)
  .use(Card).use(TreeSelect).use(Popover)
  .use(Tabs).use(FloatButton).use(Empty).use(Upload).use(Breadcrumb)
  .use(Badge).use(Tag).use(Avatar).use(Layout).use(Segmented)
  .use(Row).use(Col).use(Pagination).use(Select)
app.mount('#app')
