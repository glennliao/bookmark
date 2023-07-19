import router from '@/router'
import nProgress from 'nprogress'
import 'nprogress/nprogress.css'
import { storeToRefs } from 'pinia'
import useUserStore from '@/store/modules/user'

nProgress.configure({
  showSpinner: false
})

// 全局前置守卫
router.beforeEach((to, from,next) => {
  nProgress.start()

  const demoStore = useUserStore()
  const { token } = storeToRefs(demoStore)

  if (token.value === '') {
    console.log("??")
  }

  return next()
})

// 全局后置钩子
router.afterEach(() => {
  nProgress.done(true)
})
