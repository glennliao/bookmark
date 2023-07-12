<template>
  <a-upload-dragger
    v-model:fileList="fileList"
    name="file"
    action=""
    :maxCount="1"
    @change="handleChange"
    :beforeUpload="()=>false"
  >
    <p class="ant-upload-drag-icon">
      <inbox-outlined></inbox-outlined>
    </p>
    <p class="ant-upload-text">Click or drag .html file to this area to upload</p>
    <p class="ant-upload-hint">
      Support for a single or bulk upload. Strictly prohibit from uploading company data or other
      band files
    </p>
  </a-upload-dragger>

  <a-alert class="mt-2">
    <template #message>
      上传完毕后将自动在后台同步书签的信息,可刷新页面后查看。
    </template>
  </a-alert>

  <div class="mt-2">
    <a-button @click="fetchMeta">手动触发书签图标更新</a-button>
    <div style="margin-top: 4px;font-size: 13px;color:#818181">
      尝试更新 @import 标记的书签， 如果网络获取失败，则暂保留该标记
    </div>
  </div>
</template>
<script lang="ts" setup>
import { InboxOutlined } from '@ant-design/icons-vue'
import { ref } from 'vue'
import type { UploadChangeParam } from 'ant-design-vue'
import { message } from 'ant-design-vue'
import { apiJson } from '@/api'
import { useBookmark } from '@/views/bookmark/hook/bookmark'
const bookmark = useBookmark()
const handleChange = (info: UploadChangeParam) => {
  console.log(info)
  const reader = new FileReader()

  reader.addEventListener('load', () => {
    const text = reader.result

    apiJson.get({
      data: text,
      'import()': 'import(data)'
    }).then(data => {
      console.log('import data', data)
      bookmark.loadCate()
      bookmark.loadBookmarkList()
      bookmark.fetchMetaBatch()
      message.success('导入成功')
    })
  })

  reader.readAsText(info.file)
}

function fetchMeta(){
  bookmark.fetchMetaBatch()
}
const fileList = ref([])

</script>
