<template>
  <div class="py-2 content-area z-10 pt-10 ">


    <div class="flex">

      <div class="flex-1">
        <div class="note-item shadow-sm rounded-sm bg-white" v-for="item in list" :key="item.noteId">
          <render :tags="item.tags" :content="item.content.markdown"/>
          <div class="flex justify-between" style="border-top: 1px solid rgba(213,213,213,0.53);padding-top: 4px">
            <div>

            </div>

            <a-button size="small" @click="edit(item)">Edit</a-button>
          </div>
        </div>
      </div>

      <div class="mt-8 flex flex-wrap" style="width:200px;align-content: start">
        <div class="cursor-pointer p-2 m-1 rounded shadow-sm bg-white" v-for="item in tagList" :key="item"
             @click="filterTag(item)"
             :class="{'active-tag': activeTag == item}"
        >
          #{{ item }}
        </div>
      </div>

    </div>

    <edit-modal ref="editModalRef" @ok="load"/>


    <a-float-button-group shape="square" :style="{ right: '24px' }">
      <a-float-button type="primary" @click="edit({})">
        <template #icon>
          <PlusOutlined/>
        </template>
      </a-float-button>
      <a-float-button tooltip="toggle to bookmark" @click="toggleTo('bookmark')">
        <template #icon>
          <block-outlined />
        </template>
      </a-float-button>
    </a-float-button-group>

  </div>


</template>

<script lang="ts" setup>


import { ref, watch } from 'vue'

import EditModal from './components/EditModal.vue'
import { BlockOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { apiJson } from '@/api'
import Render from '@/views/note/components/Render.vue'
import { useRoute, useRouter } from 'vue-router'

const emit = defineEmits(["toggleTo"])
function toggleTo(e){
  emit("toggleTo",e)
}

const list = ref([])

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
    console.log(data)
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

const router = useRouter()
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
  width: 100%;
  max-width: 560px;
  margin: 6px auto;
  padding: 6px;
}

.active-tag{
  background: #46a0fc;
  color:white;
}

</style>
