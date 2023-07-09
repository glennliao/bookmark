import { ref } from 'vue'
import { apiJson } from '@/api'
import { toCateTree } from '@/utils/tree'

const cateTree = ref({ children: [] } as { children: any[] })
const cateList = ref([])
const curCate = ref([] as string[])
const curGroupId = ref('')

function loadCate () {
  apiJson.get({
    'BookmarkCate[]': {
      count: 0
    }
  }).then((data) => {
    cateList.value = data['BookmarkCate[]']
    cateTree.value = {
      children: toCateTree(data['BookmarkCate[]'])
    }
  })
}

function clickCate (cateIds: string[]) {
  curCate.value = cateIds
  loadBookmarkList()
}

const groups = ref([] as {groupId:string}[])

function loadGroup () {
  return apiJson.get({
    'Groups[]': {
      count: 0
    }
  }).then(data => {
    groups.value = data['Groups[]']
    if (groups.value.length) {
      curGroupId.value = groups.value[0].groupId
    }
  })
}

const bookmarkList = ref([] as any[])
function loadBookmarkList () {
  apiJson.get({
    'Bookmark[]': {
      cateId: curCate.value[curCate.value.length - 1],
      count: 0,
      '@order': 'id desc',
    }
  }).then(data => {
    bookmarkList.value = data['Bookmark[]']
  })
}

const curSubCateBookmark = ref([])
const curSubCateId = ref('') // 下方分类
function loadSubCateBookmark(){
  apiJson.get({
    'Bookmark[]': {
      count: 0,
      cateId:curSubCateId.value
    }
  }).then(data => {
    curSubCateBookmark.value = data['Bookmark[]']
  })
}

export function useBookmark () {
  return {
    loadCate,
    clickCate,
    curCate,
    cateTree,
    cateList,
    bookmarkList,
    curGroupId,
    loadGroup,
    loadBookmarkList,
    groups,
    curSubCateId,
    curSubCateBookmark,
    loadSubCateBookmark
  }
}
