import {ref} from "vue";
import {apiJson} from "@/api";
import {toCateTree} from "@/utils/tree";

const cateTree = ref({children: []} as { children: any[] })
const curCate = ref([] as string[])
const curGroupId = ref('')
function loadCate() {
  apiJson.get({
    'BookmarkCate[]': {
      "count":0,
    }
  }).then((data) => {

    cateTree.value = {
      children:toCateTree(data['BookmarkCate[]'])
    }

    console.log(toCateTree(data["BookmarkCate[]"]),data["BookmarkCate[]"])
  })
}

function clickCate(cateIds: string[]) {
  curCate.value = cateIds
  loadBookmarkList()
}

const groups = ref([])

function loadGroup(){
  return apiJson.get({
    "Groups[]":{
    }
  }).then(data=>{
    groups.value = data["Groups[]"]
    curGroupId.value = groups.value[0].groupId
  })
}



const bookmarkList = ref([] as any[])
function loadBookmarkList(){
  apiJson.get({
    "Bookmark[]":{
      "cateId": curCate.value[curCate.value.length-1]
    }
  }).then(data=>{
    bookmarkList.value = data["Bookmark[]"]
  })
}



export function useBookmark(){
  return {
    loadCate,clickCate,curCate,cateTree,
    bookmarkList,curGroupId, loadGroup,loadBookmarkList,groups
  }
}
