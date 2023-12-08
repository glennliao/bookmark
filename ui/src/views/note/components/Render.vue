<script setup lang="ts">
import { ref, onMounted, watch, nextTick, onUnmounted } from 'vue'
import { MdPreview, MdCatalog } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'
import * as MdEditor from 'md-editor-v3';

MdEditor.config({
  editorExtensions:{
    highlight:{
      js:"//cdn.bootcdn.net/ajax/libs/highlight.js/11.7.0/highlight.min.js",
      css:{
        atom:{
          light:"//cdn.bootcdn.net/ajax/libs/highlight.js/11.7.0/styles/atom-one-light.min.css",
          dark:"//cdn.bootcdn.net/ajax/libs/highlight.js/11.7.0/styles/atom-one-dark.min.css"
        },
      }

    },
    mermaid:{
      js:"//cdn.bootcdn.net/ajax/libs/mermaid/10.1.0/mermaidAPI-67f627de.js",
    },
    katex:{
      js:"//cdn.bootcdn.net/ajax/libs/KaTeX/0.16.3/katex.min.js",
      css:"//cdn.bootcdn.net/ajax/libs/KaTeX/0.16.3/katex.min.css"
    },

  }
})


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

  <MdPreview showCodeRowNumber :editorId="id" :modelValue="text" />
<!--  <div>-->
<!--    <MdPreview showCodeRowNumber :editorId="id" :modelValue="text" />-->
<!--&lt;!&ndash;    <MdCatalog :editorId="id" :scrollElement="scrollElement" />&ndash;&gt;-->
<!--  </div>-->
</template>

<style lang="scss">

</style>
