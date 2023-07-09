<script setup lang="ts">
import { ref, onMounted, watch, nextTick, onUnmounted } from 'vue'
import markdownIt from 'markdown-it'
import Clipboard from 'clipboard'
import "github-markdown-css/github-markdown.css"

import hljs from 'highlight.js' // https://highlightjs.org
import "highlight.js/styles/atom-one-dark.css"
import go from 'highlight.js/lib/languages/go'
import shell from 'highlight.js/lib/languages/shell'
import bash from 'highlight.js/lib/languages/bash'
import json from 'highlight.js/lib/languages/json'
import yaml from 'highlight.js/lib/languages/yaml'

hljs.registerLanguage("go",go)
hljs.registerLanguage("shell",bash)
hljs.registerLanguage("bash",bash)
hljs.registerLanguage("json",json)
hljs.registerLanguage("yaml",yaml)
import { message } from 'ant-design-vue'
const codeIndex = parseInt(Date.now()) + Math.floor(Math.random() * 10000000)

const md = new markdownIt({
  highlight: function (str, lang) {
    // 当前时间加随机数生成唯一的id标识

    // 复制功能主要使用的是 clipboard.js
    let html = `<button class="copy-btn" id="copy-btn-${codeIndex}" type="button" data-clipboard-action="copy" data-clipboard-target="#copy${codeIndex}">复制</button>`
    const linesLength = str.split(/\n/).length - 1
    // 生成行号
    let linesNum = '<span aria-hidden="true" class="line-numbers-rows">'
    for (let index = 0; index < linesLength; index++) {
      linesNum = linesNum + '<span></span>'
    }
    linesNum += '</span>'

    console.log("land",lang)
    if (lang && hljs.getLanguage(lang)) {
      try {
        // highlight.js 高亮代码
        const preCode = hljs.highlight( str, {
          language:lang,
          ignoreIllegals:true,

        }).value
        html = html + preCode
        if (linesLength) {
          html += '<b class="name">' + lang + '</b>'
        }
        // 将代码包裹在 textarea 中，由于防止textarea渲染出现问题，这里将 "<" 用 "&lt;" 代替，不影响复制功能
        return `<pre class="hljs"><code>${html}</code>${linesNum}</pre><textarea style="position: absolute;top: -9999px;left: -9999px;z-index: -9999;" id="copy${codeIndex}">${str.replace(/<\/textarea>/g, '&lt;/textarea>')}</textarea>`
      } catch (error) {
        console.log(error)
      }
    }

    const preCode = md.utils.escapeHtml(str)
    html = html + preCode
    return `<pre class="hljs"><code>${html}</code>${linesNum}</pre><textarea style="position: absolute;top: -9999px;left: -9999px;z-index: -9999;" id="copy${codeIndex}">${str.replace(/<\/textarea>/g, '&lt;/textarea>')}</textarea>`

  }
})

const props = defineProps({

  content: {
    type: String,
    default: ''
  },
  tags:{
    type:Array,
    default:()=>[]
  }
})


let clipboard = null

onMounted(() => {
  nextTick(()=>{
    clipboard = new Clipboard('#copy-btn-'+codeIndex,{

    })
    clipboard.on("success",()=>{
      message.success("success")
    })
  })
})

onUnmounted(()=>{
  clipboard && clipboard.destroy()
})

const ele = ref(null)
const html = ref('')
watch(() => [props.content,props.tags], (e) => {


  let _html = md.render(props.content)
  props.tags?.forEach(item=>{
    _html = _html.replaceAll("#"+item,`<a href="#/note?tag=${item}" >#${item}</a>`)
  })
  html.value = _html
}, { immediate: true })
</script>

<template>
  <div class="render-area markdown-body" ref="ele" v-html="html"/>
</template>

<style lang="scss">

:deep(.vditor-toolbar){
  padding-left: 0 !important;
}

:deep(.vditor-reset){
  padding: 6px !important;
}

.markdown-body {
  box-sizing: border-box;
  min-width: 200px;
  max-width: 980px;
  margin: 0 auto;
  padding: 15px;

  p {
    margin-bottom:6px;
  }
}

@media (max-width: 767px) {
  .markdown-body {
    padding: 15px;
  }
}

.render-area{
  pre.hljs {
    --color-canvas-subtle: black;
    padding: 12px 2px 12px 40px !important;
    border-radius: 5px !important;
    position: relative;
    font-size: 14px !important;
    line-height: 22px !important;
    overflow: hidden !important;
    background: #282c34;
    code {
      display: block !important;
      margin: 0 10px !important;
      overflow-x: auto !important;
      &::-webkit-scrollbar {
        z-index: 11;
        width: 6px;
      }
      &::-webkit-scrollbar:horizontal {
        height: 6px;
      }
      &::-webkit-scrollbar-thumb {
        border-radius: 5px;
        width: 6px;
        background: #666;
      }
      &::-webkit-scrollbar-corner,&::-webkit-scrollbar-track {
        background: #1E1E1E;
      }
      &::-webkit-scrollbar-track-piece {
        background: #1E1E1E;
        width: 6px
      }
    }
    .line-numbers-rows {
      position: absolute;
      pointer-events: none;
      top: 12px;
      bottom: 12px;
      left: 0;
      font-size: 100%;
      width: 40px;
      text-align: center;
      letter-spacing: -1px;
      border-right: 1px solid rgba(0, 0, 0, .66);
      user-select: none;
      counter-reset: linenumber;
      span {
        pointer-events: none;
        display: block;
        counter-increment: linenumber;
        &:before {
          content: counter(linenumber);
          color: #999;
          display: block;
          text-align: center;
        }
      }
    }
    b.name {
      position: absolute;
      top: 2px;
      right: 50px;
      z-index: 10;
      color: #999;
      pointer-events: none;
    }
    .copy-btn {
      position: absolute;
      top: 2px;
      right: 4px;
      font-size:12px;
      z-index: 10;
      color: #333;
      cursor: pointer;
      background-color: rgba(255, 255, 255, 0.74);
      border: 0;
      border-radius: 2px;
    }
  }
}


</style>
