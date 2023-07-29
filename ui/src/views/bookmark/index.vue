<template>

  <div class="bookmark-area z-10 ">

    <div class="p-2 z-10">

      <a-menu class="mb-2" :open-keys="openKeys" :selectedKeys="openKeys" mode="horizontal" selectable @select="onCateClick" @click="onCateClick" :items="cateItems" style="border-radius: 24px"/>
      <!-- content-->
      <div class="mt-4" v-if="isHome">
        <div>🕐 最近访问</div>
        <div class="flex mt-1">
          <transition-group appear name="slide-fade" tag="div" class="flex flex-wrap justify-start">
            <bookmark :simple="true" v-for="item in latestVisitList" :key="item.bmId" :item="item" @edit="edit(item)"/>
          </transition-group>
        </div>
      </div>
      <div v-else>

        <a-breadcrumb class="mb-2" v-if="curCateInfo.parents && curCateInfo.parents.length > 0">
          <a-breadcrumb-item v-for="item in curCateInfo.parents" :key="item.cateId">{{ item.title }}</a-breadcrumb-item>
          <a-breadcrumb-item :key="curCateInfo.cateId">{{ curCateInfo.title }} ({{curCateInfo.count}})</a-breadcrumb-item>
        </a-breadcrumb>



        <div class="flex">
          <transition-group appear name="slide-fade" tag="div" class="flex flex-wrap justify-start">
            <bookmark v-for="item in bookmarkList" :key="item.bmId" :item="item" @edit="edit(item)"/>
          </transition-group>
        </div>

        <div class="mt-2">
          <a-segmented class="mt-2" v-model:value="curSubCateId" :options="subCateList" style="background: #63636517;border-radius: 5px"/>

          <div class="flex mt-2">
            <transition-group appear name="slide-fade" tag="div" class="flex flex-wrap justify-start">
              <bookmark v-for="item in curSubCateBookmark" :key="item.bmId" :item="item" @edit="edit(item)"></bookmark>
            </transition-group>
          </div>
        </div>
      </div>

    </div>

<!--    <cate-manage ref="cateManageRef"/>-->
    <Setting ref="settingRef"/>
    <bookmark-edit-modal ref="BookmarkModalRef"/>
    <search ref="BookmarkSearchModalRef" @edit="edit"/>

    <a-float-button-group shape="square" :style="{ right: '24px' }">
      <a-float-button @click="searchBookmark({})">
        <template #icon>
          <SearchOutlined />
        </template>
      </a-float-button>

      <a-float-button type="primary" @click="openBookmarkModal({})">
        <template #icon>
          <PlusOutlined  />
        </template>
      </a-float-button>
    </a-float-button-group>

  </div>
</template>

<script lang="ts" setup>
import { treeEach } from '@/utils/tree'

import { ref, onMounted, h, computed, toRaw,reactive } from 'vue'
import { useBookmark } from './hook/bookmark'
import Setting from './components/Setting.vue'
import BookmarkEditModal from '@/views/bookmark/components/BookmarkEditModal.vue'
import Bookmark from '@/views/bookmark/components/Bookmark.vue'
import { PlusOutlined, SearchOutlined } from '@ant-design/icons-vue'
import { apiJson } from '@/api'
import { useRoute } from 'vue-router'
import Search from '@/views/bookmark/components/Search.vue'

const {
  loadCate,
  curCate,
  clickCate,
  cateTree,
  bookmarkList,
  loadGroup,
  curSubCateId,
  loadSubCateBookmark,
  curSubCateBookmark
} = useBookmark()

const data = reactive(['Daily', 'Weekly', 'Monthly', 'Quarterly', 'Yearly']);
const value = ref(data[0]);
const curCateInfo = ref({})
const openKeys = ref([])

const subCateList = computed(()=>{
  return (curCateInfo.value.children||[]).map(item=>{
    return {
      label:item.title,
      value:item.cateId
    }
  })
})

const cateItems = computed(() => {
  const list = [
    {
      key: '@home',
      label: ' ⭐ ',
      title: '⭐'
    }
  ]

  const tree = toRaw(cateTree.value.children || [])

  const parentMap:Record<string, string> = {}

  treeEach(tree, (item) => {
    parentMap[item.cateId] = item.parentId
  })
  treeEach(tree, (item) => {
    const labelChildren = [
      h('span', item.title)
    ]
    if (item.count > 0) {
      labelChildren.push(h('span', {
        style: {
          'font-size': '12px',
          'margin-left': '4px',
          color: '#acacac'
        }
      }, `(${item.count})`))
    }
    item.label = h('div', {
      onClick: () => {

        const keys = [item.cateId]
        let curCateId = item.cateId
        while (parentMap[curCateId] && parentMap[curCateId] !== "root") {
          keys.unshift(parentMap[curCateId])
          curCateId = parentMap[curCateId]
        }
        openKeys.value = keys

        curCateInfo.value = {}
        curSubCateBookmark.value = []
        curSubCateId.value = ''
        foundCurCateInfo(keys, cateTree.value.children, [])
        clickCate(keys)
        isHome.value = false

      },
      onTap: () => {
        console.log('tap', item.cateId)
      }
    }, labelChildren)
    item.key = item.cateId
  })

  tree.forEach(item => {
    list.push(item)
  })

  list.push({
    key: '@setting',
    label: ' ⚙️ ',
    title: '⚙️'
  })
  return list
})

function onCateClick (e) {
  switch (e.key){
    case "@home":
      isHome.value = true
      loadLatest()
      break
    case "@setting":
      settingRef.value && settingRef.value.open()
      break
  }
}


function clickSubCate (cateId: string) {
  curSubCateId.value = cateId
  loadSubCateBookmark()
}

loadGroup().then(() => {
  loadLatest()
  loadCate()
})

const settingRef = ref(null)

// 当前是否⭐页面
const isHome = ref(true)

function toHome () {

}

const latestVisitList = ref([])

function loadLatest () {
  apiJson.get({
    '[]': {
      BookmarkUse: {
        '@order': 'updatedAt desc'
      },
      Bookmark: {
        'bmId@': '/BookmarkUse/bmId'
      },
      count: 10
    }
  }).then(data => {
    latestVisitList.value = data['[]'].filter((item: any) => item.Bookmark).map((item: { Bookmark: any; }) => {
      return item.Bookmark
    })
  })
}

function foundCurCateInfo (keys: string[], tree: any[], parents: any[]) {
  for (const treeElement of tree) {
    if (treeElement.cateId === keys[keys.length - 1]) {
      curCateInfo.value = {
        ...treeElement,
        parents
      }

      if (treeElement.children && treeElement.children.length > 0) {
        clickSubCate(treeElement.children[0].cateId)
      }

      return
    }
    if (keys.includes(treeElement.cateId) && treeElement.children) {
      foundCurCateInfo(keys, treeElement.children, parents.concat(treeElement))
    }
  }
}

function handleMenuClick (keyPath: string[]) {
  const keys = keyPath

  curCateInfo.value = {}
  curSubCateBookmark.value = []
  curSubCateId.value = ''
  foundCurCateInfo(keys, cateTree.value.children, [])
  clickCate(keys)
  isHome.value = false
}

const BookmarkModalRef = ref()
const BookmarkSearchModalRef = ref()

function openBookmarkModal (item) {
  BookmarkModalRef.value.open(item)
}

function searchBookmark (item) {
  BookmarkSearchModalRef.value.open(item)
}

function edit (item) {
  BookmarkModalRef.value.open(item)
}

// add from url
const route = useRoute()
if (route.query.url) {
  onMounted(() => {
    openBookmarkModal({ url: route.query.url })
  })
}
</script>
<style scoped lang="scss">

.menu > li > div {
  //padding: 12px;
}
.active {
  border-bottom: 2px solid;
}

.slide-fade-enter-active {
  transition: all .4s ease;
}

.slide-fade-leave-active {
  transition: all .4s cubic-bezier(1.0, 0.5, 0.8, 1.0);
}

.slide-fade-enter {
  transform: translateX(-50px);
  opacity: .1;
}

.slide-fade-leave-to {
  transform: translateX(-50px);
  opacity: .1;
}

:deep(.ant-menu-horizontal >.ant-menu-item){
  min-width: 66px;
  text-align: center;
  padding-inline: 2px;
}

</style>
