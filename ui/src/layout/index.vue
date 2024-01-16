<script setup lang="ts">

import Footer from '@/layout/Footer.vue'
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core'
import {
  FormOutlined,
  BookOutlined,
  UserOutlined,
  AppstoreOutlined, RightOutlined
} from '@ant-design/icons-vue'
import useUserStore from '@/store/modules/user'
import { computed, ref } from 'vue'
import { apiJson } from '@/api'
import { compareVersion } from '@/utils/str-utils'

const breakpoints = useBreakpoints(breakpointsTailwind)
const smallerThanSm = breakpoints.smallerOrEqual('sm')

const router = useRouter()
const collapsed = ref<boolean>(true)

let key = router.currentRoute.value.path

const selectedKeys = ref<string[]>([key])
const route = useRoute()
const defaultKey = 'bookmark'

watch(() => route.path, (v: any) => {
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

const version = import.meta.env.VITE_app_version.trim()
const latestVersion = ref('v0.0.0')

function checkNewVersion () {
  apiJson.get({
    'version()': 'latestVersion()'
  }).then(data => {
    if (typeof data.version === 'string') {
      latestVersion.value = data.version
      // latestVersion.value = "v0.20.0"
    }
  })
}

const canUpdate = computed(() => {
  return compareVersion(latestVersion.value || '', version) > 0
})

setTimeout(()=>{
  checkNewVersion()
},256)
</script>

<template>
  <a-config-provider>
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
          <a-sub-menu key="sub1" style="position: absolute;bottom:62px;width:100%">
            <template #title>
            <span>
              <user-outlined/>
              <span>User</span>
            </span>
            </template>
            <a-menu-item key="@logout">Logout</a-menu-item>
          </a-sub-menu>
        </a-menu>
        <div class="version-info" style="position: absolute;bottom:20px;width:100%;height: 30px;text-align: center;">
          <a href="https://github.com/glennliao/bookmark" target="_blank" style="display: block;height: 30px">
            <div>
              <div class="text-white">
                <svg t="1686376535417" class="icon" viewBox="0 0 1024 1024" version="1.1"
                     xmlns="http://www.w3.org/2000/svg"
                     p-id="2370" width="16" height="16">
                  <path
                    d="M512 12.64c-282.752 0-512 229.216-512 512 0 226.208 146.72 418.144 350.144 485.824 25.6 4.736 35.008-11.104 35.008-24.64 0-12.192-0.48-52.544-0.704-95.328-142.464 30.976-172.512-60.416-172.512-60.416-23.296-59.168-56.832-74.912-56.832-74.912-46.464-31.776 3.52-31.136 3.52-31.136 51.392 3.616 78.464 52.768 78.464 52.768 45.664 78.272 119.776 55.648 148.992 42.56 4.576-33.088 17.856-55.68 32.512-68.48-113.728-12.928-233.28-56.864-233.28-253.024 0-55.904 20-101.568 52.768-137.44-5.312-12.896-22.848-64.96 4.96-135.488 0 0 43.008-13.76 140.832 52.48 40.832-11.36 84.64-17.024 128.16-17.248 43.488 0.192 87.328 5.888 128.256 17.248 97.728-66.24 140.64-52.48 140.64-52.48 27.872 70.528 10.336 122.592 5.024 135.488 32.832 35.84 52.704 81.536 52.704 137.44 0 196.64-119.776 239.936-233.792 252.64 18.368 15.904 34.72 47.04 34.72 94.816 0 68.512-0.608 123.648-0.608 140.512 0 13.632 9.216 29.6 35.168 24.576 203.328-67.776 349.856-259.616 349.856-485.76 0-282.784-229.248-512-512-512z"
                    fill="#fff" p-id="2371"></path>
                </svg>
              </div>
              <div class="mt-2 text-white">
                  <span v-if="canUpdate">
                    <a-badge count="new" style="padding: 0 2px;height: 14px;line-height: 14px">
                      <a target="_blank" href="https://github.com/glennliao/bookmark" class="ml-1">{{ latestVersion }}</a>
                    </a-badge>
                  </span>
                <span v-else>
                    {{latestVersion}}
                </span>
              </div>
            </div>
          </a>

        </div>
      </a-layout-sider>
      <a-layout-content :style="{'margin-left':smallerThanSm?'2px':'86px'}" style="max-height: 100vh;overflow: auto">
        <div class="body z-10 w-full">
          <router-view/>
        </div>
        <!--        <Footer/>-->
      </a-layout-content>
    </a-layout>

  </a-config-provider>
</template>

<style scoped>
.body {
  //min-height: calc(100vh - 26px);
  //padding-bottom: 40px;
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

:deep(.ant-layout-sider-zero-width-trigger){
  z-index: 10002;
}

:deep(.ant-layout-sider-collapsed.ant-layout-sider-zero-width){
  .version-info{
    display: none;
  }
}
</style>
