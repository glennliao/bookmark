<template>
  <a-modal
    v-model:open="visible"
    width="460px"
  >
    <template #footer>
      <a-button @click="close">关闭</a-button>
    </template>

    <a-tabs destroy-inactive-tab-pane>
      <a-tab-pane key="1" tab="目录管理">
        <CateManage/>
      </a-tab-pane>
      <a-tab-pane key="2" tab="导入书签">
        <Import/>
      </a-tab-pane>
    </a-tabs>
  </a-modal>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import CateManage from './CateManage.vue'
import Import from './Import.vue'

const visible = ref(false)

function open () {
  visible.value = true
}

function close () {
  visible.value = false
}

defineExpose({
  open
})

</script>

<style scoped>

</style>
