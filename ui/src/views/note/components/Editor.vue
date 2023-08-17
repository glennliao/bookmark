<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { MdEditor } from 'md-editor-v3';
import 'md-editor-v3/lib/style.css';
import sanitizeHtml from 'sanitize-html';

const sanitize = (html) => sanitizeHtml(html);
const text = ref('');

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

})

watch(()=>props.modelValue, (e)=>{
  if (e !== value.value){
    text.value = e
  }
},{immediate:true})

defineExpose({
  getContent:()=>{
    return text.value
  }
})


</script>

<template>
  <MdEditor placeholder="Hello Editor!" v-model="text" :sanitize="sanitize" autoFocus :toolbars="[
    // 'bold',
  // 'underline',
  // 'italic',
  // '-',
  // 'title',
  // 'strikeThrough',
  // 'sub',
  // 'sup',
  // 'quote',
  // 'unorderedList',
  // 'orderedList',
  // 'task',
  // '-',
  // 'codeRow',
  // 'code',
  // 'link',
  // 'image',
  // 'table',
  // 'mermaid',
  // 'katex',
  // '-',
  // 'revoke',
  // 'next',
  // 'save',
  // '=',
  // 'pageFullscreen',
  // 'fullscreen',
  'preview',
  // 'htmlPreview',
  // 'catalog',
  // 'github'
  ]"/>
  <div style="font-size: 13px;color:rgba(159,159,159,0.78)">
    markdown
  </div>
</template>

<style scoped lang="scss">

#vditor {
  border-radius: 6px;
  border: 2px solid #e7e7e7 !important;
}

:deep(.vditor-toolbar) {
  padding-left: 0 !important;
  display: none;
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
