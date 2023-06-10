<template>
  <!--  <HeaderBar/>-->


  <div class="px-4 py-2 bookmark-area z-10 absolute " v-contextmenu:contextmenu>

    <div class="toggle-bar" @click="groupsVisible = !groupsVisible">
      <div class="n-layout-toggle-bar">
        <div class="n-layout-toggle-bar__top"></div>
        <div class="n-layout-toggle-bar__bottom"></div>
      </div>
    </div>

    <div class="p-2 z-10">
      <ul class="rounded-box bg-base-100 menu menu-horizontal bg-base-100 shadow z-10" style="width: auto">
        <li :class="{'active':isHome}">
          <span @click="toHome">
            ⭐
          </span>
        </li>
        <li v-for="item in cateTree.children" :key="item.cateId"
            :class="{'active':!isHome && curCate.includes(item.cateId)}">
          <a-dropdown trigger="hover">
            <template #overlay>
              <a-menu v-if="item.children" @click="handleMenuClick">
                <template v-for="subItem in item.children" :key="subItem.cateId">
                  <template v-if="!subItem.children">
                    <a-menu-item :key="item.cateId+'/'+subItem.cateId">
                      {{ subItem.title }}
                    </a-menu-item>
                  </template>
                  <template v-else>
                    <SubMemu :key="item.cateId" :menu-info="item"/>
                  </template>
                </template>
              </a-menu>
            </template>

            <div @click="handleMenuClick({keyPath:[item.cateId]})">
              <span>
              {{ item.title }}
          </span>

            </div>
          </a-dropdown>
        </li>
        <li @click="cateContextmenu">
          <span>⚙️</span>
        </li>
      </ul>
      <div v-if="isHome">
        <div>🕐 最近访问</div>
        <div class="flex mt-1">
          <transition-group appear name="slide-fade" tag="div" class="flex flex-wrap justify-start">
            <bookmark :simple="true" v-for="item in latestVisitList" :key="item.bmId" :item="item" @edit="edit(item)"/>
          </transition-group>
        </div>
      </div>
      <div v-else>

        <div class="text-sm breadcrumbs" v-if="curCateInfo.parents.length > 0">
          <ul>
            <li v-for="item in curCateInfo.parents" :key="item.cateId">{{ item.title }}</li>
          </ul>
        </div>

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

    <cate-manage ref="cateManageRef"/>
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
    </a-float-button-group>

  </div>

  <a-drawer
    title="Groups"
    placement="left"
    :closable="true"
    v-model:open="groupsVisible"
    width="180px"
  >
    <div @click="changeGroup(item)" v-for="item in groups" style="cursor: pointer"
         class="shadow p-2 rounded bg-orange-500 text-white" :key="item.groupId">
      {{ item.title }}
    </div>
  </a-drawer>


  <v-contextmenu ref="contextmenu">
    <v-contextmenu-item @click="openBookmarkModal({})">➕ 新增书签</v-contextmenu-item>
    <v-contextmenu-item @click="cateContextmenu({})">⚙️ 管理分类</v-contextmenu-item>
    <v-contextmenu-item @click="logout({})">🔘 退出登录</v-contextmenu-item>
  </v-contextmenu>
</template>

<script lang="ts" setup>
// import HeaderBar from './layout/Header.vue'
import { ref, onMounted } from 'vue'
import SubMemu from '@/views/components/SubMemu.vue'
import { useBookmark } from './hook/bookmark'
import CateManage from '@/views/components/CateManage.vue'
import BookmarkEditModal from '@/views/components/BookmarkEditModal.vue'
import Bookmark from '@/views/components/Bookmark.vue'
import { PlusOutlined,SearchOutlined } from '@ant-design/icons-vue'
import { apiJson } from '@/api'
import { useRoute } from 'vue-router'
import Search from '@/views/components/Search.vue'

const {
  loadCate,
  curCate,
  clickCate,
  cateTree,
  bookmarkList,
  loadGroup,
  curGroupId,
  groups,
  curSubCateId,
  loadSubCateBookmark,
  curSubCateBookmark
} = useBookmark()

const curCateInfo = ref({})

function clickSubCate (cateId: string) {
  curSubCateId.value = cateId
  loadSubCateBookmark()
}

loadGroup().then(() => loadCate())

function changeGroup (item) {
  curGroupId.value = item.groupId
  loadCate()
  groupsVisible.value = false
}

const cateManageRef = ref(null)

function cateContextmenu (action, cate) {
  cateManageRef.value && cateManageRef.value.open()
}

const groupsVisible = ref(false)

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
    console.debug(data)
    latestVisitList.value = data['[]'].filter((item: any) => item.Bookmark).map((item: { Bookmark: any; }) => {
      return item.Bookmark
    })
  })
}

loadLatest()

function foundCurCateInfo (keys: string[], tree: any[], parents: any[]) {
  console.log(tree)
  for (const treeElement of tree) {
    if (treeElement.cateId === keys[keys.length - 1]) {
      curCateInfo.value = {
        ...treeElement,
        parents
      }
      if (treeElement.children && treeElement.children.length > 0) {
        console.log('fkc', treeElement.children)
        clickSubCate(treeElement.children[0].cateId)
      }

      return
    }
    if (keys.includes(treeElement.cateId) && treeElement.children) {
      foundCurCateInfo(keys, treeElement.children, parents.concat(treeElement))
    }
  }
  console.log(curCateInfo.value)
}

function handleMenuClick (e: { keyPath: string[] }) {
  const keys = e.keyPath[0].split('/').filter(item => item)
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

function logout () {
  localStorage.removeItem('token')
  location.reload()
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

.bookmark-area {
  height: 100vh;
  width: 100vw;
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

.n-layout-toggle-bar:hover .n-layout-toggle-bar__top {
  transform: rotate(-12deg) scale(1.15) translateY(-2px);
}

.n-layout-toggle-bar:hover .n-layout-toggle-bar__bottom {
  transform: rotate(12deg) scale(1.15) translateY(2px);
}

.n-layout-toggle-bar {

  //left: -28px;
  transform: rotate(180deg);

}

.n-layout-toggle-bar:hover .n-layout-toggle-bar__top {
  transform: rotate(12deg) scale(1.15) translateY(-2px);
}

.n-layout-toggle-bar:hover .n-layout-toggle-bar__bottom {
  transform: rotate(-12deg) scale(1.15) translateY(2px);
}

.n-layout-toggle-bar:hover .n-layout-toggle-bar__top {
  transform: rotate(-12deg) scale(1.15) translateY(-2px);
}

.n-layout-toggle-bar:hover .n-layout-toggle-bar__bottom {
  transform: rotate(12deg) scale(1.15) translateY(2px);
}

.toggle-bar {
  --n-bezier: cubic-bezier(0.4, 0, 0.2, 1);
  --n-toggle-button-color: #FFF;
  --n-toggle-button-border: 1px solid rgb(239, 239, 245);
  --n-toggle-bar-color: rgba(191, 191, 191, 1);
  --n-toggle-bar-color-hover: rgba(153, 153, 153, 1);
  --n-color: #fff;
  --n-text-color: rgb(51, 54, 57);
  --n-border-color: rgb(239, 239, 245);
  --n-toggle-button-icon-color: rgb(51, 54, 57);
  position: absolute;
  left: -4px;
  top: 0;
  width: 12px;
  bottom: 0;
}

.n-layout-toggle-bar {

  cursor: pointer;
  height: 72px;
  width: 32px;
  position: absolute;
  top: calc(50% - 36px);
  //right: -28px;

}

.n-layout-toggle-bar .n-layout-toggle-bar__top, .n-layout-toggle-bar .n-layout-toggle-bar__bottom {

  position: absolute;
  width: 4px;
  border-radius: 2px;
  height: 38px;
  left: 14px;
  transition: background-color .3s var(--n-bezier),
  transform .3s var(--n-bezier);

}

.n-layout-toggle-bar .n-layout-toggle-bar__bottom {

  position: absolute;
  top: 34px;

}

.n-layout-toggle-bar:hover .n-layout-toggle-bar__top {
  transform: rotate(12deg) scale(1.15) translateY(-2px);
}

.n-layout-toggle-bar:hover .n-layout-toggle-bar__bottom {
  transform: rotate(-12deg) scale(1.15) translateY(2px);
}

.n-layout-toggle-bar .n-layout-toggle-bar__top, .n-layout-toggle-bar .n-layout-toggle-bar__bottom {
  background-color: var(--n-toggle-bar-color);
}

.n-layout-toggle-bar:hover .n-layout-toggle-bar__top, .n-layout-toggle-bar:hover .n-layout-toggle-bar__bottom {
  background-color: var(--n-toggle-bar-color-hover);
}

</style>
