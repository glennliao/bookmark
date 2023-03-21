<template>
  <a-modal
    title="书签"
    v-model:visible="visible"
    width="60%"
  >
    <div>
      <a-form>
        <div v-if="!hasFetchURL">
          <a-input-search
            v-model:value="info.url"
            placeholder="url"
            enter-button="Search"
            size="large"
            @search="onSearch"
            :loading="fetchLoading"
          />

          <div v-if="equalBm.bmId">
            <a-card>
              <div>
                已存在相同url
              </div>
              <div>
                {{ equalBm.title }}
              </div>
              <div>
                {{ equalBm.url }}
              </div>
              <div>
                {{ equalBm.cate }}
              </div>
            </a-card>
          </div>
        </div>
        <div v-else>
          <a-form-item label="目录" v-bind="validateInfos.parentId">
            <a-tree-select v-model:value="info.parentId" :tree-data="parentIdTreeData"/>
          </a-form-item>
          <a-alert show-icon v-if="info.title.length > 20"
                   style="width:60%;text-align: center;margin-left: 10%;margin-bottom: 4px">
            <template #message>标题过长, 不利于查看</template>
          </a-alert>
          <a-form-item label="名称" v-bind="validateInfos.title">
            <a-input v-model:value="info.title"/>
          </a-form-item>
          <a-form-item label="地址" v-bind="validateInfos.url">
            <a-input v-model:value="info.url"/>
          </a-form-item>
          <a-form-item label="图标" v-bind="validateInfos.icon">
            <a-input v-model:value="info.icon">
              <template #addonAfter>
                <div>
                  <img :src="info.icon" style="width: 20px;height: 20px;max-width: 20px"/>
                </div>
              </template>
            </a-input>
          </a-form-item>
          <a-form-item label="描述" v-bind="validateInfos.description">
            <a-textarea :rows="3" showCount v-model:value="info.description" type="textarea"/>
          </a-form-item>
          <a-form-item label="备注" v-bind="validateInfos.remark">
            <a-textarea :rows="3" showCount v-model:value="info.remark" type="textarea"/>
          </a-form-item>
        </div>
      </a-form>
    </div>

    <template #footer>
      <a-button key="back" @click="visible = false">Close</a-button>
      <a-button key="submitAndNext" v-if="hasFetchURL && !info.bmId" type="primary" :loading="addLoading" @click="handleAdd(true)">
        Submit & Next
      </a-button>
      <a-button key="submit" v-if="hasFetchURL" type="primary" :loading="addLoading" @click="handleAdd(false)">Submit
      </a-button>
    </template>
  </a-modal>

</template>

<script setup>
import {reactive, ref, toRaw} from 'vue'
import {apiJson} from '../../api'
import {toCateTree} from '@/utils/tree'
import {useBookmark} from '@/views/hook/bookmark'
import {useForm} from 'ant-design-vue/es/form/index.js'
import {message} from 'ant-design-vue'

const visible = ref(false)
const addVisible = ref(false)
const treeData = ref([])
const parentIdTreeData = ref([])
const bookmark = useBookmark()
const info = ref({
  parentId: '',
  title: '',
  groupId: '',
  url: ''
})

const hasFetchURL = ref(false)
const fetchLoading = ref(false)

const equalBm = ref({})


const {resetFields, validate, validateInfos} = useForm(
  info,
  reactive({
    parentId: [
      {
        required: true
      }
    ],
    title: [
      {
        required: true
      }
    ],
    url: [
      {
        required: true
      }
    ]
  })
)

function onSearch() {
  let url = info.value.url.trim()

  if (!url) {
    message.warn('url 不能为空')
    return
  }
  fetchLoading.value = true


  apiJson.get({
    url: url,
    'meta()': 'fetchURL(url)',
    "equal": {
      "Bookmark": {
        url // or suffix /
      },
      "GroupBookmark": {
        "bmId@": "/Bookmark/bmId"
      },
      "BookmarkCate": {
        "cateId@": "/GroupBookmark/cateId"
      }
    },
    "domain": {
      "Bookmark[]": {
        "url$": "%" + (new URL(url)).host + "%"
      }
    },
  }).then(data => {
    let equal = data.equal.Bookmark
    if (equal.bmId) {
      equal.cate = data.equal.BookmarkCate.title
      equalBm.value = equal
    } else {
      hasFetchURL.value = true
    }
    if (bookmark.curSubCateId.value) {
      info.value.parentId = bookmark.curSubCateId.value || ''
    } else {
      if(bookmark.curCate.value.length){
        info.value.parentId = bookmark.curCate.value[bookmark.curCate.value.length - 1] || ''
      }
    }
    info.value.url = data.meta.url
    info.value.title = data.meta.title
    info.value.icon = data.meta.icon
    info.value.description = data.meta.description

  }).finally(() => {
    fetchLoading.value = false
  })
}

const addLoading = ref(false)

function handleAdd(next) {
  addLoading.value = true
  const _info = toRaw(info.value)
  validate().then(data => {

    let api = apiJson.post
    let GroupBookmark = {
      cateId: _info.parentId,
      groupId: bookmark.curGroupId.value,
    }
    if(_info.bmId){
      api = apiJson.put
      GroupBookmark.id = groupBookmarkId.value
    }else{
      GroupBookmark["bmId@"] = "Bookmark/bmId"
    }

    api({
      Bookmark: _info,
      GroupBookmark,
      tag: 'Bookmark'
    }).then(data => {
      console.log(data)
      hasFetchURL.value = false
      info.value = {}
      if (!next) {
        visible.value = false
      }
      bookmark.loadCate()
      bookmark.loadBookmarkList()
      if (bookmark.curSubCateId.value){
        bookmark.loadSubCateBookmark()
      }
    })
  }).finally(() => {
    addLoading.value = false
  })
}

const groupBookmarkId = ref('')
function open(_info = {}) {
  resetFields()
  equalBm.value = {}
  visible.value = true
  info.value = {
    ..._info
  }

  if(_info.bmId){
    apiJson.get({
      "GroupBookmark":{
        groupId: bookmark.curGroupId.value,
        "bmId":_info.bmId
      }
    }).then(data=>{
      groupBookmarkId.value = data.GroupBookmark.id
      info.value = {
        ...info.value,
        parentId:data.GroupBookmark.cateId
      }
      console.log(data)
    })
  }
  hasFetchURL.value = _info.bmId
  loadCate()
}

function loadCate() {
  apiJson.get({
    'BookmarkCate[]': {
      count: 0
    }
  }).then(data => {
    treeData.value = toCateTree((data['BookmarkCate[]'].map(item => {
      return {
        ...item,
        label: item.title,
        value: item.cateId
      }
    })))
    parentIdTreeData.value = toCateTree([].concat(data['BookmarkCate[]'].map(item => {
      return {
        ...item,
        label: item.title,
        value: item.cateId
      }
    })))
  })
}

function add() {
  addVisible.value = true
}

defineExpose({
  open
})

</script>

<style scoped>

</style>
