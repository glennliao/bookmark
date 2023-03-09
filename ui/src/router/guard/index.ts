import router from '@/router'
import nProgress from 'nprogress'
import 'nprogress/nprogress.css'
import useDemoStore from '@/store/modules/demo'
import { storeToRefs } from 'pinia'

nProgress.configure({
  showSpinner: false
})

// 全局前置守卫
router.beforeEach((to, from) => {
  nProgress.start()

  const demoStore = useDemoStore()
  const { counter } = storeToRefs(demoStore)
  // 从 store 中获取其他值，再决定返回值
  // 这里演示获取 store 中 counter 的值
  console.log(`counter：${counter.value}`)
  return true
})

// 全局后置钩子
router.afterEach(() => {
  nProgress.done(true)
})
