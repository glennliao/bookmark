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
</template>
<script lang="ts" setup>
import { InboxOutlined } from '@ant-design/icons-vue'
import { ref } from 'vue'
import type { UploadChangeParam } from 'ant-design-vue'
import { message } from 'ant-design-vue'
import { apiJson } from '@/api'
import { useBookmark } from '@/views/hook/bookmark'
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
      message.success('导入成功')
    })
  })

  reader.readAsText(info.file)
}
const fileList = ref([])

</script>
