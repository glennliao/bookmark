<script setup lang="ts">
import {ref} from 'vue'
import Bookmark from '@/views/bookmark/index.vue'
import Note from '@/views/note/index.vue'
import { UserOutlined,DownOutlined } from '@ant-design/icons-vue'

const mode = ref('bookmark')

function logout () {
  localStorage.removeItem('token')
  location.reload()
}

</script>

<template>
  <div>
    <bookmark @toggleTo="mode = $event" v-if="mode === 'bookmark'"/>
    <note @toggleTo="mode = $event" v-if="mode === 'note'"/>
  </div>

  <div class="user-area ">
    <a-dropdown>
      <a class="ant-dropdown-link" @click.prevent>
        <a-avatar  :size="32">
          <template #icon><UserOutlined /></template>
        </a-avatar>
      </a>
      <template #overlay>
        <a-menu>
          <a-menu-item @click="logout">
            <a href="javascript:;">logout</a>
          </a-menu-item>
        </a-menu>
      </template>
    </a-dropdown>
  </div>
</template>

<style scoped>
.user-area{
  position: absolute;
  right: 24px;
  top:24px;
}
</style>
