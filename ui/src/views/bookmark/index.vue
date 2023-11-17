<template>

  <div class="bookmark-area z-10 ">

    <div class="z-10">

      <div class="p-2">

        <a-menu inlineCollapsed class="mb-2" :selected-keys="openKeys" mode="horizontal" selectable
                @select="onCateClick"
                @click="onCateClick" :items="cateItems" style="border-radius: 24px"/>
        <span v-if="smallerThanSm">如果点击某个父级目录，移动端下请点击点两次</span>
      </div>

      <!-- content-->
      <div class="mt-4" v-if="isHome">

        <div class="ml-2 mb-2 mt-4 section-title">🕐 <span >最近访问</span></div>
        <div class="flex mt-1">
          <transition-group appear name="slide-fade" tag="div" class="flex flex-wrap justify-start">
            <bookmark :simple="true" v-for="item in latestVisitList" :key="item.bmId" :item="item" @edit="edit(item)"/>
          </transition-group>
        </div>

        <template v-for="(tag,index) in tagInHome" :key="index">
          <a-popover :open="popupOpen && selectedTagIndex === index" title="Tag" trigger="click"
                     @click="getTagOptions(tag.tag, index)">
            <template #content>
              <a-select
                  v-model:value="selectedTag"
                  style="width: 200px"
                  :options="tagsOptions"
                  allow-clear
              ></a-select>
              <a-button @click="submitTag" class="ml-2">确定</a-button>
            </template>
            <div class="ml-2 mb-2 mt-8 cursor-pointer section-title" style="font-weight: 500;display: inline-block">⭐ <span>{{ tag.tag }}</span> <EditOutlined /></div>
          </a-popover>

          <div class="flex mt-2">
            <transition-group appear name="slide-fade" tag="div" class="flex flex-wrap justify-start">
              <bookmark :simple="true" v-for="item in tag.bookmarkList" :key="item.bmId" :item="item"
                        @edit="edit(item)"/>
            </transition-group>
          </div>
        </template>

        <div class="mt-2 cursor-pointer" style="color:#ccc">
          <a-popover :open="popupOpen && selectedTagIndex === -1" title="Tag" trigger="click"
                     @click="getTagOptions('', -1)">
            <template #content>
              <a-select
                  v-model:value="selectedTag"
                  style="width: 200px"
                  :options="tagsOptions"
                  allow-clear
              ></a-select>
              <a-button @click="submitTag" class="ml-2">确定</a-button>
            </template>
            <div class="ml-2 mb-2 mt-8 cursor-pointer" style="color: #ccc"> + add new tag</div>
          </a-popover>

        </div>
      </div>

      <div v-else>
        <a-breadcrumb class="mb-2 px-2" v-if="curCateInfo.parents && curCateInfo.parents.length > 0">
          <a-breadcrumb-item v-for="item in curCateInfo.parents" :key="item.cateId">{{ item.title }}</a-breadcrumb-item>
          <a-breadcrumb-item :key="curCateInfo.cateId">{{ curCateInfo.title }} ({{ curCateInfo.count }})
          </a-breadcrumb-item>
        </a-breadcrumb>

        <div class="flex">
          <transition-group appear name="slide-fade" tag="div" class="flex flex-wrap justify-start">
            <bookmark v-for="item in bookmarkList" :key="item.bmId" :item="item" @edit="edit(item)"/>
          </transition-group>
        </div>

        <div class="mt-2">

          <a-segmented v-if="subCateList.length" class="mt-2 mx-2" @change="clickSubCate" v-model:value="curSubCateId"
                       :options="subCateList"
                       style="background: rgb(255 255 255 / 44%);border-radius: 20px;max-width: 98%;overflow-x: auto;padding: 2px 0.5rem">
            <template #label="{ payload,title }">
              {{ title }} <span v-if="payload.count"
                                style="font-size: 12px;margin-left: 2px">({{ payload.count }})</span>
            </template>
          </a-segmented>

          <div class="flex mt-2">
            <transition-group appear name="slide-fade" tag="div" class="flex flex-wrap justify-start">
              <bookmark v-for="item in curSubCateBookmark" :key="item.bmId" :item="item" @edit="edit(item)"></bookmark>
            </transition-group>
          </div>
        </div>
      </div>

    </div>

    <!--    <cate-manage ref="cateManageRef"/>-->
    <search ref="BookmarkSearchModalRef" @edit="edit"/>

    <Setting ref="settingRef"/>
    <bookmark-edit-modal ref="BookmarkModalRef"/>

    <a-float-button-group shape="square" :style="{ right: '24px' }">
      <a-float-button @click="searchBookmark({})">
        <template #icon>
          <SearchOutlined/>
        </template>
      </a-float-button>

      <a-float-button type="primary" @click="openBookmarkModal({})">
        <template #icon>
          <PlusOutlined/>
        </template>
      </a-float-button>
    </a-float-button-group>

  </div>
</template>

<script lang="ts" setup>
import { treeEach } from '@/utils/tree'

import { ref, onMounted, h, computed, toRaw, getCurrentInstance } from 'vue'
import { useBookmark } from './hook/bookmark'
import Setting from './components/Setting.vue'
import BookmarkEditModal from '@/views/bookmark/components/BookmarkEditModal.vue'
import Bookmark from '@/views/bookmark/components/Bookmark.vue'
import { PlusOutlined, SearchOutlined, EditOutlined } from '@ant-design/icons-vue'
import { apiJson } from '@/api'
import { useRoute } from 'vue-router'
import Search from '@/views/bookmark/components/Search.vue'
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core'

const {
  loadCate,
  clickCate,
  cateTree,
  bookmarkList,
  loadGroup,
  curSubCateId,
  loadSubCateBookmark,
  curSubCateBookmark
} = useBookmark()

const breakpoints = useBreakpoints(breakpointsTailwind)
const smallerThanSm = breakpoints.smallerOrEqual('sm')

const curCateInfo = ref({})
const openKeys = ref(['@home'])

const subCateList = computed(() => {
  return (curCateInfo.value.children || []).map(item => {
    return {
      title: `${item.title}`,
      value: item.cateId,
      payload: {
        count: item.count
      }
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

  const parentMap: Record<string, string> = {}

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
      // class: item.cateId === openKeys.value[0] ? 'menu-active' : '',
      onClick: () => {
        const keys = [item.cateId]
        let curCateId = item.cateId
        while (parentMap[curCateId] && parentMap[curCateId] !== 'root') {
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
    item.isMe = item.key === openKeys.value[0]
    list.push(item)
  })

  list.push({
    key: '@setting',
    label: ' ⚙️ ',
    title: '⚙️'
  })
  return list
})

function onCateClick(e) {
  console.log("onCate",e)
  switch (e.key) {
    case '@home':
      isHome.value = true
      openKeys.value = ['@home']
      loadLatest()
      break
    case '@setting':
      settingRef.value && settingRef.value.open()
      break
  }
}

function clickSubCate(cateId: string) {
  console.log("clickSubCate",cateId)
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

const latestVisitList = ref([])

function loadLatest() {
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
      item.Bookmark.tags = JSON.parse(item.Bookmark.tags || '[]')
      return item.Bookmark
    })
  })
}

function foundCurCateInfo(keys: string[], tree: any[], parents: any[]) {
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

const BookmarkModalRef = ref()
const BookmarkSearchModalRef = ref()

function openBookmarkModal(item) {
  BookmarkModalRef.value.open(item)
}

function searchBookmark(item) {
  BookmarkSearchModalRef.value.open(item)
}

function edit(item) {
  BookmarkModalRef.value.open(item)
}

// add from url
const route = useRoute()
if (route.query.url) {
  onMounted(() => {
    openBookmarkModal({ url: route.query.url })
  })
}

const popupOpen = ref(false)
const tagInHome = ref([])
const tagsOptions = ref([])
const selectedTag = ref('')
const selectedTagIndex = ref(-2)

function getTagOptions(tag: string, index: number) {
  selectedTag.value = tag
  selectedTagIndex.value = index
  apiJson.get({
    'tags()': 'bmTags()'
  }).then(data => {
    tagsOptions.value = data.tags.map(item => {
      return {
        value: item
      }
    })
    popupOpen.value = true
  })
}

const currentInstance = getCurrentInstance()

function loadTagInHome() {
  apiJson.get({
    Config: {
      key: 'tagInHome'
    }
  }).then(data => {
    console.log(data)
    tagInHome.value = JSON.parse(data.Config.value || '[]')
    if (tagInHome.value.length === 0) {
      tagInHome.value = [
        { tag: '选择书签标签用以首页显示', placeholder: true }
      ]
    } else {
      tagInHome.value = tagInHome.value.map(item => {
        item.bookmarkList = ref([])

        apiJson.get({
          'Bookmark[]': {
            '@alias': 'list',
            '@order': 'createdAt desc',
            tags: item.tag
          }
        }).then(data => {
          item.bookmarkList = data.list.map(item => {
            item.tags = JSON.parse(item.tags)
            return item
          })
          currentInstance?.proxy.$forceUpdate()
        })
        return item
      })
    }
  })
}

loadTagInHome()

function submitTag() {

  let value = tagInHome.value

  if (selectedTagIndex.value === -1) {
    if (selectedTag.value) {
      value.push({
        tag: selectedTag.value
      })
    }
  } else {
    if (!selectedTag.value) {
      value = []
    } else {
      value[selectedTagIndex.value].tag = selectedTag.value
    }
  }


  apiJson.put({
    Config: {
      key: 'tagInHome',
      value
    },
    tag: 'Config'
  }).then(data => {
    popupOpen.value = false
    loadTagInHome()
  })
}

</script>

<style lang="scss">

li[isme="true"]{
  color: #1677ff;
  font-weight: bold;
}

li[isme="true"]:not(.ant-menu-submenu-active).ant-menu-submenu:after {
  position: absolute;
  inset-inline: 16px;
  bottom: 0;
  border-bottom: 2px solid #1677ff;
  transition: border-color 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
  content: "";

}
</style>

<style scoped lang="scss">

.section-title{
  font-weight: 500;
  span:nth-child(1){
    filter: invert(100%)  grayscale(1) contrast(9);

  }

  span:nth-child(2){
    color: #ccc;
  }
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

:deep(.ant-menu-horizontal >.ant-menu-item) {
  min-width: 66px;
  text-align: center;
  padding-inline: 12px !important;
}

:deep(.ant-menu-horizontal .ant-menu-submenu) {
  padding-inline: 12px !important;
}

:deep(.ant-segmented .ant-segmented-group) {
  padding: 2px;
}

:deep(.ant-segmented) {
  &::-webkit-scrollbar {
    height: 6px;
  }

  &::-webkit-scrollbar-thumb {
    background: #ccc; // 滑块颜色
    border-radius: 5px; // 滑块圆角
  }

  &::-webkit-scrollbar-thumb:hover {
    background: rgba(63, 170, 231, 0.73); // 鼠标移入滑块变红
  }

  .ant-segmented-item {
    border-radius: 12px;
  }
}

</style>
