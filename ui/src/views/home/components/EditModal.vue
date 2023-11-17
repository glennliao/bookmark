<script setup lang="ts">
import { apiJson } from '@/api'

const emit = defineEmits("refresh")
const openModel = ref(false)
const editInfo = ref({})


const open = () => {
  editInfo.value = {}
  openModel.value = true
}

const submit = ()=>{
  console.log("subm",editInfo, editInfo.value)
  apiJson.post({
    "AppList":{
      ...toRaw(editInfo.value)
    },
    "tag":"AppList"
  }).then(data=>{
    emit("refresh")
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
    <a-form>
      {{editInfo}}
      <a-form-item label="title">
        <a-input v-model:value="editInfo.title"></a-input>
      </a-form-item>
      <a-form-item label="url">
        <a-input v-model:value="editInfo.url"></a-input>
      </a-form-item>
      <a-form-item label="icon">
        <a-input v-model:value="editInfo.icon"></a-input>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<style scoped>

</style>
