<template>
  <a-modal
    title="书签搜索"
    v-model:open="visible"
  >
    <div>
      <a-input-search
        v-model:value="info.q"
        placeholder="url or title"
        enter-button="Search"
        size="default"
        ref="input"
        @search="onSearch"
        :loading="fetchLoading"
      />

      <div v-if="hasSearch" class="mt-4 p-1 rounded" style="background: #efefef;overflow: auto;max-height: 72vh">
        <a-empty v-if="list.length === 0" />
        <div v-for="item in list" :key="item.bmId" @click="toURL(item)" class="p-2 m-2 bg-white rounded shadow cursor-pointer">

          <div class="font-bold">
            {{item.title.length > 30 ? item.title.substring(0, 30) + '...' : item.title}}
          </div>
          <div>
            {{item.url}}
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex justify-between">
        <div>
          <span style="font-size: 13px;color:#acacac">
              也可以直接ctrl+f打开搜索框
            </span>
        </div>
        <div>
          <a-button key="back" @click="visible = false">Close</a-button>
        </div>
      </div>
    </template>
  </a-modal>

</template>

<script setup>

import { apiJson } from '../../api'
import { message } from 'ant-design-vue'
import { nextTick, ref } from 'vue'

const visible = ref(false)

const list = ref([])
const info = ref({})

const fetchLoading = ref(false)
const hasSearch = ref(false)

const equalBm = ref({})

window.addEventListener('keydown', function (event) {
  if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === 'f' && !visible.value) {
    visible.value = true

    event.preventDefault()
    nextTick(() => {

      input.value.focus()
    })
  }
})


function onSearch () {
  const q = info.value.q.trim()

  if (!q) {
    message.warn('不能为空')
    return
  }


  fetchLoading.value = true
  let beginAt = new Date().getTime()

  apiJson.get({
    'Bookmark[]': {
      q,
      '@alias': 'list',
    }
  }).then(data => {
    list.value = data.list
    hasSearch.value = true

  }).finally(() => {
    setTimeout(() => {
      fetchLoading.value = false
    }, (new Date().getTime() - beginAt) > (512 * 1000) ? 0 : 512)
  })
}

const input = ref(null)
function open (_info = {}) {
  visible.value = true
  nextTick(() => {
    input.value.focus()
  })
}

function toURL (item) {
  apiJson.put({
    tag: 'BookmarkUse',
    BookmarkUse: {
      bmId: item.bmId
    }
  })
  window.open(item.url, '_blank')
}

defineExpose({
  open
})

</script>

<style scoped>

</style>
