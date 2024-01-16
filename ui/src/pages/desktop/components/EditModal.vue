<script setup lang="ts">
import { apiJson } from '@/api'
import { ref } from 'vue'

const emit = defineEmits(['refresh'])
const openModel = ref(false)
const formModel = ref({})

const open = (item = {}) => {
  formModel.value = item
  formModel.value.sorts = formModel.value.sorts || 10
  formModel.value.openType = formModel.value.openType || ''
  formModel.value.openParam = formModel.value.openParam || JSON.stringify({
    width: 450,
    height: 300
  })
  openModel.value = true
}

const submit = () => {

  let api = apiJson.post

  if (formModel.value.id) {
    api = apiJson.put
  }

  api({
    'AppList': {
      ...toRaw(formModel.value)
    },
    'tag': 'AppList'
  }).then(data => {
    emit('refresh')
    openModel.value = false
  })
}

defineExpose({
  open
})
</script>

<template>
  <a-modal title="edit"
           v-model:open="openModel"
           :mask-closable="false"
           @ok="submit"
  >
    <a-form layout="vertical">
      <a-form-item label="title">
        <a-input v-model:value="formModel.title"></a-input>
      </a-form-item>
      <a-form-item label="url">
        <a-input v-model:value="formModel.url"></a-input>
      </a-form-item>
      <a-form-item label="openType" help="新浏览器标签适用于无法使用iframe打开的页面">
        <a-radio-group  v-model:value="formModel.openType" :options="[
          {
            label:'内部打开',
            value:''
          },
          {
            label:'新浏览器标签',
            value:'_blank'
          },
        ]"></a-radio-group>
      </a-form-item>

      <a-form-item label="openParam">
        <a-textarea  v-model:value="formModel.openParam" ></a-textarea>
      </a-form-item>
      <a-form-item label="sorts">
        <a-input-number  v-model:value="formModel.sorts" ></a-input-number>
      </a-form-item>

      <a-form-item label="icon">
        <a-input v-model:value="formModel.icon"></a-input>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<style scoped>
:deep(.ant-form-item){
  margin-bottom: 6px;
}
</style>
