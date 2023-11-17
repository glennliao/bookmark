<script setup lang="ts">
import Footer from '@/layout/Footer.vue'
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core'
import {
  FormOutlined,
  BookOutlined,
  UserOutlined,
  AppstoreOutlined
} from '@ant-design/icons-vue'
import useUserStore from '@/store/modules/user'

const breakpoints = useBreakpoints(breakpointsTailwind)
const smallerThanSm = breakpoints.smallerOrEqual('sm')

const router = useRouter()
const collapsed = ref<boolean>(true)

let key = router.currentRoute.value.path

const selectedKeys = ref<string[]>([key])
const route = useRoute()
const defaultKey = 'bookmark'

watch(() => route.path, (v:any) => {
  let path = v
  if (path === '/') {
    path = defaultKey
  }
  selectedKeys.value = [path]
})

if (key === '/') {
  key = defaultKey
}

const userStore = useUserStore()
const goto = (e: { key: string }) => {
  if (e.key.startsWith('@')) {
    userStore.logout()
  } else {
    router.push(e.key)
  }
}

</script>

<template>
  <a-layout style="min-height:100vh;">
    <a-layout-sider v-model:collapsed="collapsed"
                    :breakpoint="smallerThanSm?'sm':undefined"
                    :collapsed-width="smallerThanSm?0:undefined"
                    :class="{'side-fixed':!smallerThanSm}"
                    :style="{'background':'#ffffff1f'}"
    >
      <a-menu @click="goto" :selectedKeys="selectedKeys" theme="dark" mode="inline">
        <a-menu-item key="/home">
          <AppstoreOutlined/>
          <span>Home</span>
        </a-menu-item>
        <a-menu-item key="/bookmark">
          <book-outlined/>
          <span>Bookmark</span>
        </a-menu-item>
        <a-menu-item key="/note">
          <form-outlined/>
          <span>Note</span>
        </a-menu-item>
        <a-sub-menu key="sub1" style="position: absolute;bottom:12px;width:100%">
          <template #title>
            <span>
              <user-outlined/>
              <span>User</span>
            </span>
          </template>
          <a-menu-item key="@logout">Logout</a-menu-item>
        </a-sub-menu>
      </a-menu>
    </a-layout-sider>
    <a-layout-content :style="{'margin-left':smallerThanSm?'2px':'86px'}" style="max-height: 100vh;overflow: auto">
      <div class="body z-10 w-full">
        <router-view/>
      </div>
      <Footer/>
    </a-layout-content>
  </a-layout>

</template>

<style>
.body {
  min-height: calc(100vh - 26px);
  padding-bottom: 40px;
  font-size: 14px;
  position: relative;
}

.side-fixed {
  overflow: auto !important;
  height: 100vh !important;
  position: fixed !important;
  left: 0;
  top: 0;
  bottom: 0
}

.ant-layout-sider-zero-width-trigger-left {
  top: 44px !important;
  z-index: 1000 !important;
}

</style>
