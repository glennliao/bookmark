import { computed, ref } from 'vue'
import { apiJson } from '@/api'
import { toCateTree, treeEach } from '@/utils/tree'

const cateList = ref([])
const curCate = ref([] as string[])
const curGroupId = ref('')
const cateBookmarkNumMap = ref({})
const cateTree = computed(() => {
  let treeList = toCateTree(cateList.value)
  treeEach(treeList, (item) => {
    item.count = cateBookmarkNumMap.value[item.cateId] || 0
    if(item.children){
      item.children.forEach(subItem=>{
        item.count += subItem.count
      })
    }
  })
  return {
    children: treeList
  }
})

function loadCate () {
  apiJson.get({
    'BookmarkCate[]': {
      count: 0
    },
    'cateBmNum()': 'cateBmNum()'
  }).then((data) => {
    cateList.value = data['BookmarkCate[]']
    loadCateBmNum()
  })
}

function loadCateBmNum () {
  apiJson.get({
    'cateBmNum()': 'cateBmNum()'
  }).then((data) => {
    let list = data.cateBmNum
    let m = {}
    list.forEach(item => {
      m[item.cateId] = item.cnt
    })
    cateBookmarkNumMap.value = m
  })
}

function clickCate (cateIds: string[]) {
  curCate.value = cateIds
  loadBookmarkList()
}

const groups = ref([] as { groupId: string }[])

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
      '@order': 'createdAt desc'
      // 'cateId()': 'cateIdByBmId(bmId,groupId)'
    }
  }).then(data => {
    bookmarkList.value = data['Bookmark[]']
  })
}

const curSubCateBookmark = ref([])
const curSubCateId = ref('') // 下方分类
function loadSubCateBookmark () {
  apiJson.get({
    'Bookmark[]': {
      count: 0,
      cateId: curSubCateId.value,
      '@order': 'createdAt desc'
    }
  }).then(data => {
    curSubCateBookmark.value = data['Bookmark[]']
  })
}

function fetchMetaBatch () {
  apiJson.get({
    'fetchMetaBatch()': 'fetchMetaBatch()'
  }).then((data) => {

  })
}

export function useBookmark () {
  return {
    fetchMetaBatch,
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
    loadSubCateBookmark,
    loadCateBmNum
  }
}
