<template>
  <div id="bg"></div>

  <a-layout style="min-height:100vh;">
    <a-layout-sider v-model:collapsed="collapsed"
                    :style="{ overflow: 'auto', height: '100vh', position: 'fixed', left: 0, top: 0, bottom: 0 }"
    >
      <div class="logo" />
      <a-menu @click="to" :selectedKeys="selectedKeys" theme="dark" mode="inline">
        <a-menu-item key="/bookmark">
          <book-outlined />
          <span>Bookmark</span>
        </a-menu-item>
        <a-menu-item key="/note">
          <form-outlined />
          <span>Note</span>
        </a-menu-item>
<!--        <a-sub-menu key="sub1">-->
<!--          <template #title>-->
<!--            <span>-->
<!--              <user-outlined />-->
<!--              <span>User</span>-->
<!--            </span>-->
<!--          </template>-->
<!--          <a-menu-item key="3">Tom</a-menu-item>-->
<!--          <a-menu-item key="4">Bill</a-menu-item>-->
<!--          <a-menu-item key="5">Alex</a-menu-item>-->
<!--        </a-sub-menu>-->
      </a-menu>
    </a-layout-sider>
    <a-layout>
      <a-layout-content style="margin: 2px 4px 2px 82px">
        <div :style="{ padding: '12px', minHeight: '360px' }">
          <div class="body z-10 w-full">
            <router-view/>
          </div>
        </div>
      </a-layout-content>
      <a-layout-footer style="text-align: center;padding: 6px 50px">
        <Footer/>
      </a-layout-footer>
    </a-layout>
  </a-layout>



</template>
<script setup lang="ts">
import Footer from '@/layout/Footer.vue'
import { ref, watch } from 'vue'
import {
  PieChartOutlined,
  DesktopOutlined,
  FormOutlined,
  BookOutlined,
  TeamOutlined,
  FileOutlined,
  UserOutlined, VideoCameraOutlined, UploadOutlined,MenuFoldOutlined,MenuUnfoldOutlined } from '@ant-design/icons-vue';
import { useRoute, useRouter } from 'vue-router'


const router = useRouter()
const collapsed = ref<boolean>(true);
let key = router.currentRoute.value.path

const selectedKeys = ref<string[]>([key]);
const route = useRoute()

watch(() => route.path,(v)=>{
  console.log(v)
  let path = v
  if (path === "/"){
    path = "/bookmark"
  }
  console.log(path)
  selectedKeys.value = [path]
})

console.log(router.currentRoute.value.path,key,router.currentRoute.value)
if (key == "/") {
  key = "/bookmark"
}

function to(e){
  router.push(e.key)
}
</script>

<style>
.body{
  min-height: calc(100vh - 26px - 40px);
  font-size: 14px;
  position: relative;
}

</style>
