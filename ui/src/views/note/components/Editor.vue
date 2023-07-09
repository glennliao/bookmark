<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import Vditor from 'vditor'
import 'vditor/dist/index.css'

const vditor = ref<Vditor | null>(null)

const emit = defineEmits(['update:modelValue'])
const props = defineProps({
  preview: {},
  modelValue: {
    type: String,
    default: ''
  }
})

const value = ref('')


onMounted(() => {
  vditor.value = new Vditor('vditor', {
    mode: 'ir',
    height: '400px',
    toolbar: [],
    "tab": "\t",
    placeholder: 'write something',
    after: () => {
      vditor.value!.setValue(props.modelValue)
    },
    input (e) {
      value.value = e
      emit('update:modelValue', e)
    }
  })
})

watch(()=>props.modelValue, (e)=>{
  if (e !== value.value){
    vditor.value!.setValue(e)
    value.value = e
  }
})

defineExpose({
  getContent:()=>{
    return vditor.value.getValue()
  }
})


</script>

<template>
  <div id="vditor"/>
  <div style="font-size: 13px;color:rgba(159,159,159,0.78)">
    support markdown
  </div>
</template>

<style scoped lang="scss">

:deep(.vditor-toolbar) {
  padding-left: 0 !important;
}

:deep(.vditor){
  border-radius: 6px;
  border: 2px solid #ccc !important;
}

:deep(.vditor-reset) {
  padding: 6px !important;

  p{
    margin-bottom: 2px;
  }
}
</style>
