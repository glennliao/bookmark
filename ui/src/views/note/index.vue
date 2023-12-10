<template>

  <div class="py-2 content-area z-10 ">

    <div class="flex mx-auto">

      <a-row class="w-full">
        <a-col :sm="24" :md="0" class="w-full">
          <div class="py-2 mx-2 my-2 mt-1">
            <a-input-search allow-clear autofocus placeholder="search something..." v-model:value="searchKey"
                            @search="search"/>
          </div>
        </a-col>
        <a-col :sm="24" :md="0">
          <div class="mt-2 ml-2 flex flex-wrap" style="align-content: start">
            <div class="cursor-pointer p-2 m-1 rounded shadow-sm bg-white" v-for="item in tagList" :key="item"
                 @click="filterTag(item)"
                 :class="{'active-tag': activeTag == item}"
            >
              #{{ item }}
            </div>
          </div>
        </a-col>
        <a-col :sm="24" :md="18" class="w-full">

          <div class="note-item shadow rounded  w-full " v-for="(item,index) in list" :key="item.noteId">
            <div class="max-note-height relative"  :style="{
              '--note-max-height': item.openFlag?'auto':'240px',
              paddingBottom:noteHeightMap[item.noteId]>240?'20px':0
            }">
              <render :id="item.noteId" @onHtmlChanged="onHtmlChange(item)" :tags="item.tags" :content="item.content.markdown" style="display: block"/>
              <div v-if="noteHeightMap[item.noteId]>240" @click="item.openFlag = !item.openFlag" class="cursor-pointer text-white" style="position: absolute;background: rgb(139 139 139 / 66%);bottom:0;width:100%;text-align: center;height: 20px;line-height: 20px">
                  展开 / 收起
              </div>
            </div>

            <div class="flex justify-end" style="border-top: 1px solid rgba(213,213,213,0.53);padding: 6px 12px;background: #e0e0e069">
              <a-button size="small" danger class="mr-2" @click="del(item)">del</a-button>
              <a-button size="small" @click="edit(item)">Edit</a-button>
            </div>
          </div>

          <div class="p-2 text-center">
            <a-pagination v-model:current="pageNum" :default-page-size="pageSize" :total="total" show-less-items @change="onPageChange"/>
          </div>
        </a-col>
        <a-col v-if="!smallerThanSm" :sm="0" :md="6" >
          <div class="mt-8 mx-4 ">
            <a-input-search ref="searchInput" allow-clear autofocus placeholder="search something..." v-model:value="searchKey"
                            @search="search"/>
          </div>
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

import { createVNode, ref, watch } from 'vue'

import EditModal from './components/EditModal.vue'
import { ExclamationCircleOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { apiJson } from '@/api'
import Render from '@/views/note/components/Render.vue'
import { useRouter } from 'vue-router'
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core'
import { Modal } from 'ant-design-vue'

const router = useRouter()

const list = ref([])

const breakpoints = useBreakpoints(breakpointsTailwind)

const smallerThanSm = breakpoints.smallerOrEqual('sm')

const pageNum = ref(1)
const pageSize = ref(6)
const total = ref(0)

function loadList(searchKey = '') {
  const p = {}
  if (activeTag.value) {
    p.tag = activeTag.value
  }
  if (searchKey) {
    p.q = searchKey
  }
  apiJson.get({
    'Note[]': {
      '@alias': 'list',
      '@order': 'id desc',
      'count': pageSize.value,
      'page': pageNum.value,
      ...p
    },
    'total@': 'Note[]/total',
  }).then(data => {
    total.value = data.total
    list.value = data.list.map(item => {
      item.content = JSON.parse(item.content)
      item.tags = JSON.parse(item.tags || '[]')
      item.openFlag = false
      return item
    })
  })
}

function onPageChange(_pageNum,_pageSize){
  pageNum.value = _pageNum
  loadList(searchKey.value)
}

const tagList = ref([])

function loadTagList() {
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
}, {
  immediate: true
})

function filterTag(tag) {
  if (tag === activeTag.value) {
    tag = ''
  }
  router.push({
    path: router.currentRoute.value.path,
    query: {
      tag
    }
  })
}

function load() {
  loadList()
  loadTagList()
}

load()

const editModalRef = ref(null)

function edit(item) {
  editModalRef.value.open(item)
}

function del(item){
  Modal.confirm({
    title: 'Do you Want to delete this item?',
    icon: createVNode(ExclamationCircleOutlined),
    content: createVNode('div', { style: 'color:red;' }, ''),
    onOk() {
      apiJson.delete({
        "Note":{
          noteId:item.noteId
        },
        "tag":"Note"
      }).then(()=>{
        loadList(searchKey.value)
      })
    },
    onCancel() {
      console.log('Cancel');
    },
  });
}

const searchKey = ref('')
function search() {
  loadList(searchKey.value)
}

const searchInput = ref(null)
onMounted(()=>{
  searchInput.value && searchInput.value.focus()
})


const onHtmlChange = (item)=>{
  const noteId = item.noteId
  const h = document.getElementById(noteId).clientHeight
  noteHeightMap.value[noteId] = h
}


const noteHeightMap = ref({})

</script>



<style scoped lang="scss">

.note-item {
  margin: 0  auto 6px;
  padding: 2px;
}

.active-tag {
  background: #46a0fc;
  color: white;
}

.max-note-height {
  max-height: var(--note-max-height);
  overflow: hidden;
}

</style>
