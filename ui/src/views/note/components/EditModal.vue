<template>
  <a-modal
    title="Note"
    v-model:open="visible"
    destroy-on-close
    :mask-closable="false"
    width="100%"
    :style="{'max-width':'780px'}"
  >
    <div>
      <editor ref="editor" :model-value="info.content.markdown" />
    </div>

    <template #footer>
      <div>

        <a-button key="back" @click="visible = false" class="mr-2">Close</a-button>

        <a-dropdown-button type="primary" :loading="addLoading" @click="handleAdd(false)">
          Submit
<!--          <template #overlay>-->
<!--            <a-menu>-->
<!--              <a-menu-item v-if="hasFetchURL && !info.bmId" key="1" @click="handleAdd(true)">Submit & Next-->
<!--              </a-menu-item>-->
<!--            </a-menu>-->
<!--          </template>-->
<!--          <template #icon>-->
<!--            <DownOutlined/>-->
<!--          </template>-->
        </a-dropdown-button>
      </div>
    </template>
  </a-modal>

</template>

<script setup>
import {  reactive, ref, toRaw } from 'vue'
import { apiJson } from '../../../api/index'
import { useBookmark } from '@/views/bookmark/hook/bookmark'
import { useForm } from 'ant-design-vue/es/form/index'
import { message } from 'ant-design-vue'
import {  DownOutlined } from '@ant-design/icons-vue'
import Editor from '@/views/note/components/Editor.vue'

const visible = ref(false)
const addVisible = ref(false)

const bookmark = useBookmark()
const info = ref({
  content: {
    markdown:""
  },
})

const emit = defineEmits(["ok"])

const {
  resetFields,
  validate,
  validateInfos
} = useForm(
  info,
  reactive({
    parentId: [
      {
        required: true
      }
    ],
    title: [
      {
        required: true
      }
    ],
    url: [
      {
        required: true
      }
    ]
  })
)

const addLoading = ref(false)
const editor = ref(null)
function handleAdd (next) {
  addLoading.value = true
  const _info = toRaw(info.value)

  let markdownText = editor.value.getContent()


  if (!markdownText){
    message.warn('请输入内容')
    addLoading.value = false
    return
  }

  let api = apiJson.post

  const param = {
    "Note":{
      content:{
        markdown:markdownText,
      }
    },
    "tag":"Note"
  }

  const tagReg = /#([^#\s]+)/g

  const tags = new Set((markdownText.match(tagReg)||[]).map(item=>item.slice(1)))

  param.Note.tags = Array.from(tags)


  if (_info.noteId) {
    api = apiJson.put
    param.Note.noteId = _info.noteId
  } else {

  }

  api(param).then(data => {
    visible.value = false

    console.log(data)
    emit("ok")
  }).finally(()=>{
    addLoading.value = false
  })
}

const groupBookmarkId = ref('')

function open (_info = {}) {
  resetFields()
  visible.value = true
  info.value = {
    content:{
      markdown:""
    },
    ..._info
  }

  console.log(info)
}

function add () {
  addVisible.value = true
}

defineExpose({
  open
})

</script>

<style scoped>

:deep(.ant-form-item) {
  margin-bottom: 16px;
}

:deep(.ant-upload-wrapper .ant-upload-drag) {
  border: none;
}

:deep(.ant-upload-wrapper .ant-upload-drag :hover) {
  border: 1px dashed #d9d9d9;
  border-radius: 8px;
}

:deep(.ant-upload-wrapper .ant-upload-drag .ant-upload) {
  padding: 0;
}

@media (max-width: 575px) {

  :deep(.ant-form-item .ant-form-item-label) {
    padding: 0;
  }
}
</style>
