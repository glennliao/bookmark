<template>

  <div class="py-2 content-area z-10 ">
    <div class="flex mx-auto" >

      <a-row class="w-full">
        <a-col :sm="24" :md="0" >
          <div class="mt-8 ml-10 flex flex-wrap" style="align-content: start">
            <div class="cursor-pointer p-2 m-1 rounded shadow-sm bg-white" v-for="item in tagList" :key="item"
                 @click="filterTag(item)"
                 :class="{'active-tag': activeTag == item}"
            >
              #{{ item }}
            </div>
          </div>
        </a-col>
        <a-col :sm="24" :md="18" class="w-full">
          <div class="note-item shadow-sm rounded-sm bg-white w-full"  v-for="item in list" :key="item.noteId">
            <render :tags="item.tags" :content="item.content.markdown" style="display: block"/>
            <div class="flex justify-end" style="border-top: 1px solid rgba(213,213,213,0.53);padding-top: 4px">
              <a-button size="small" @click="edit(item)">Edit</a-button>
            </div>
          </div>
        </a-col>
        <a-col v-if="!smallerThanSm" :sm="0" :md="6">
          <div class="mt-8 ml-4 flex flex-wrap" style="align-content: start">
          <div class="cursor-pointer p-2 m-1 rounded shadow-sm bg-white" v-for="item in tagList" :key="item"
               @click="filterTag(item)"
               :class="{'active-tag': activeTag == item}"
          >
            #{{ item }}
          </div>
        </div>
        </a-col>
      </a-row>


    </div>

    <edit-modal ref="editModalRef" @ok="load"/>


    <a-float-button-group shape="square" :style="{ right: '24px' }">
      <a-float-button type="primary" @click="edit({})">
        <template #icon>
          <PlusOutlined/>
        </template>
      </a-float-button>
    </a-float-button-group>

  </div>


</template>

<script lang="ts" setup>


import { ref, watch } from 'vue'

import EditModal from './components/EditModal.vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { apiJson } from '@/api'
import Render from '@/views/note/components/Render.vue'
import { useRouter } from 'vue-router'
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core'

const router = useRouter()

const list = ref([])

const breakpoints = useBreakpoints(breakpointsTailwind)

const smallerThanSm = breakpoints.smallerOrEqual('sm')

function loadList () {
  let p = {}
  if(activeTag.value){
    p.tag = activeTag.value
  }
  apiJson.get({
    'Note[]': {
      '@alias': 'list',
      '@order': 'id desc',
      ...p
    }
  }).then(data => {
    list.value = data.list.map(item => {
      item.content = JSON.parse(item.content)
      item.tags = JSON.parse(item.tags || '[]')
      return item
    })
  })

}

const tagList = ref([])

function loadTagList () {
  apiJson.get({
    'tags()': 'noteTags()'
  }).then(data => {
    tagList.value = data.tags
  })
}

const activeTag = ref('')
watch(() => router.currentRoute.value.query, (e) => {
  activeTag.value = e.tag
  loadList()
},{
  immediate:true
})

function filterTag (tag) {

  if(tag === activeTag.value){
    tag = ""
  }
  router.push({
    path: router.currentRoute.value.path,
    query: {
      tag
    }
  })
}

function load () {

  loadList()
  loadTagList()

}

load()

const editModalRef = ref(null)

function edit (item) {
  editModalRef.value.open(item)
}


</script>
<style scoped lang="scss">

.note-item {
  margin: 6px auto;
  padding: 6px;
}

.active-tag{
  background: #46a0fc;
  color:white;
}

</style>
