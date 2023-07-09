<template>

  <div class="py-2 bookmark-area z-10 " v-contextmenu:contextmenu>

    <div class="p-2 z-10">
      <ul class="rounded-box bg-base-100 menu category menu-horizontal bg-base-100 shadow z-10" style="width: auto">
        <li :class="{'active':isHome}">
          <span @click="toHome">
            ⭐
          </span>
        </li>
        <li v-for="item in cateTree.children" :key="item.cateId"
            :class="{'active':!isHome && curCate.includes(item.cateId)}">
          <a-dropdown trigger="hover">
            <template #overlay>
              <a-menu v-if="item.children" >
                <template v-for="subItem in item.children" :key="subItem.cateId">
                  <template v-if="!subItem.children">
                    <a-menu-item :key="item.cateId +'/'+subItem.cateId">
                      <div  @click="handleMenuClick([item.cateId,subItem.cateId])">{{ subItem.title }}</div>
                    </a-menu-item>
                  </template>
                  <template v-else>
                    <SubMemu @click="e=>handleMenuClick([subItem.parentId,...e])" :key="subItem.cateId" :menu-info="subItem"/>
                  </template>
                </template>
              </a-menu>
            </template>

            <div @click="handleMenuClick([item.cateId])">
              <span>{{ item.title }}</span>
            </div>
          </a-dropdown>
        </li>
        <li @click="cateContextmenu">
          <span>⚙️</span>
        </li>
      </ul>

      <!-- content-->
      <div v-if="isHome">
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
          <a-breadcrumb-item :key="curCateInfo.cateId">{{ curCateInfo.title }}</a-breadcrumb-item>
        </a-breadcrumb>


        <div class="flex">
          <transition-group appear name="slide-fade" tag="div" class="flex flex-wrap justify-start">
            <bookmark v-for="item in bookmarkList" :key="item.bmId" :item="item" @edit="edit(item)"/>
          </transition-group>
        </div>

        <div class="mt-4">
          <div className="tabs">
            <a :className="'tab tab-sm tab-lifted '+(curSubCateId === item.cateId?'tab-active':'')"
               v-for="item in curCateInfo.children" :key="item.cateId"
               @click="clickSubCate(item.cateId)">{{ item.title }}</a>
          </div>
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
    <search ref="BookmarkSearchModalRef"/>


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
      <a-float-button tooltip="toggle to note" @click="toggleTo('note')">
        <template #icon>
          <block-outlined />
        </template>
      </a-float-button>
    </a-float-button-group>

  </div>


  <v-contextmenu ref="contextmenu">
    <v-contextmenu-item @click="openBookmarkModal({})">➕ 新增书签</v-contextmenu-item>
    <v-contextmenu-item @click="cateContextmenu({})">⚙️ 管理分类</v-contextmenu-item>
    <v-contextmenu-item @click="logout({})">🔘 退出登录</v-contextmenu-item>
  </v-contextmenu>
</template>

<script lang="ts" setup>

import { ref, onMounted } from 'vue'
import SubMemu from '@/views/bookmark/components/SubMemu.vue'
import { useBookmark } from './hook/bookmark'
import Setting from './components/Setting.vue'
import BookmarkEditModal from '@/views/bookmark/components/BookmarkEditModal.vue'
import Bookmark from '@/views/bookmark/components/Bookmark.vue'
import { PlusOutlined,SearchOutlined,BlockOutlined } from '@ant-design/icons-vue'
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

const emit = defineEmits(["toggleTo"])
function toggleTo(e){
  emit("toggleTo",e)
}

const curCateInfo = ref({})

function clickSubCate (cateId: string) {
  curSubCateId.value = cateId
  loadSubCateBookmark()
}

loadGroup().then(() => {
  loadLatest()
  loadCate()
})


const settingRef = ref(null)

function cateContextmenu () {
  settingRef.value && settingRef.value.open()
}


// 当前是否⭐页面
const isHome = ref(true)

function toHome () {
  isHome.value = true
  loadLatest()
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
  console.log("tree",tree)
  for (const treeElement of tree) {
    if (treeElement.cateId === keys[keys.length - 1]) {
      curCateInfo.value = {
        ...treeElement,
        parents
      }

      console.log(curCateInfo.value,"curCateInfo")

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


  console.log("keys",keys)

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

function searchBookmark(item){
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


</style>
