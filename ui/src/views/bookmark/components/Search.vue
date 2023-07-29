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
        <div v-for="item in list" :key="item.bmId" class="m-2 mr-4">
          <Bookmark :item="item" width="100%" @drop="onSearch" @edit="edit"/>
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

import { apiJson } from '../../../api'
import { message } from 'ant-design-vue'
import { nextTick, ref } from 'vue'
import { useBookmark } from '@/views/bookmark/hook/bookmark'
import Bookmark from '@/views/bookmark/components/Bookmark.vue'

const visible = ref(false)

const list = ref([])
const info = ref({})

const fetchLoading = ref(false)
const hasSearch = ref(false)

const {
  curGroupId,
  cateList
} = useBookmark()

const emit = defineEmits(["edit"])

function edit(e){
  emit("edit",e)
}

window.addEventListener('keydown', function (event) {
  if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === 'f' && !visible.value) {
    open({})

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
  const beginAt = new Date().getTime()

  apiJson.get({
    'Bookmark[]': {
      q,
      '@alias': 'list',
      '@order': 'createdAt desc',
      '@column': 'bmId,title,remark,url,icon'
    },
    groupId: curGroupId.value

  }).then(data => {
    list.value = data.list.map(item => {
      const cateInfo = getCateInfo(item.cateId)
      item.cateInfo = cateInfo
      return item
    })
    hasSearch.value = true
  }).finally(() => {
    setTimeout(() => {
      fetchLoading.value = false
    }, (new Date().getTime() - beginAt) > (512 * 1000) ? 0 : 512)
  })
}

const input = ref(null)
function open (_info = {}) {

  hasSearch.value = false
  list.value = []
  info.value = {
    q:""
  }
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

function getCateInfo (cateId) {
  const list = []
  let parentId = cateId
  let i = 10
  do {
    cateList.value.forEach(item => {
      if (parentId === item.cateId) {
        list.unshift({
          cateId: item.cateId,
          title: item.title
        })
        parentId = item.parentId
      }
    })
    i--
    if(i < 0){
      console.log(parentId,cateId)
      break
    }

  } while (parentId !== 'root' && parentId !== '')

  return list
}

defineExpose({
  open
})

</script>

<style  scoped>
:deep(.ant-breadcrumb .ant-breadcrumb-separator){
  margin-inline: 2px;
}

:deep(.ant-breadcrumb){
  font-size: 12px;
}

.logo {
  width: 24px;
  min-width: 24px;
  height: 24px;
  line-height: 24px;
  border-radius: 100%;
  text-transform: uppercase;
  font-weight: bold;
}

.title {
  font-size: 16px;
  color: #000;
  line-height: 24px;
  font-weight: bold;
}
</style>
