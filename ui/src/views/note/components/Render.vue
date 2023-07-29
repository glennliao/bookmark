<script setup lang="ts">
import { ref, onMounted, watch, nextTick, onUnmounted } from 'vue'
import { MdPreview, MdCatalog } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'

const props = defineProps({

  content: {
    type: String,
    default: ''
  },
  tags: {
    type: Array,
    default: () => []
  }
})

const id = 'preview-only'
const scrollElement = document.documentElement

onMounted(() => {
  nextTick(() => {

  })
})

onUnmounted(() => {

})

const text = ref('')
watch(() => [props.content, props.tags], (e) => {
  let _text = props.content
  props.tags?.forEach(item=>{
    _text = _text.replaceAll("#"+item,`<a href="#/note?tag=${item}" >#${item}</a>`)
  })
  text.value = _text
}, { immediate: true })
</script>

<template>
  <view>
    <MdPreview showCodeRowNumber :editorId="id" :modelValue="text" />
    <MdCatalog :editorId="id" :scrollElement="scrollElement" />
  </view>
</template>

<style lang="scss">

</style>
