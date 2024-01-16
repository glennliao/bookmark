<script setup lang="ts">
import VueResizable from 'vue-resizable'

const props = defineProps({
  width: {
    type: Number,
    default: 450
  },
  height: {
    type: Number,
    default: 200
  },
  title: {
    type: String
  },
  windowStyle: {

  }
})

const emit = defineEmits(['action'])
const action = (action:'close') => {
  emit('action', {
    action
  })
}


const resizeIng = ref(false)
const info = ref({})
const onResizeStart = (e)=>{
  resizeIng.value = true
}

const onResizeMove = (e)=>{
  info.value = e
}

const onResizeEnd = (e)=>{
  resizeIng.value = false
}


</script>

<template>
  <vue-resizable class="window" dragSelector=".title-bar"
                 :width="props.windowStyle.width"
                 :height="props.windowStyle.height"
                 :left="props.windowStyle.x"
                 :top="props.windowStyle.y"
                 :style="{'z-index': props.windowStyle.zIndex}"
                 @resize:start="onResizeStart"
                 @resize:move="onResizeMove"
                 @resize:end="onResizeEnd"
  >
    <div class="w-full h-full bg-white rounded shadow">
      <div class="title-bar">
       <div>
         {{title}}
       <span v-if="resizeIng" class="ml-2 text-green font-bold" >
         ({{info.width}} * {{info.height}})
       </span>
       </div>
       <div class="text-right">
         <span class="cursor-pointer mr-2" @click="action('min')">-</span>
         <span class="cursor-pointer" @click="action('close')">x</span>
       </div>
      </div>
      <div class="window-body">
        <slot/>
      </div>
    </div>
  </vue-resizable>
</template>

<style scoped lang="less">
.window {
  position: fixed;

  .title-bar {
    background: white;
    user-select: none;
    height: 26px;
    padding: 6px 12px;
    border-radius: 4px;
    display: flex;
    justify-content: space-between;
  }

  .window-body{
    height: calc(100% - 26px);
  }

  iframe {
    width: 100%;
    height: 100%;
  }
}
</style>
