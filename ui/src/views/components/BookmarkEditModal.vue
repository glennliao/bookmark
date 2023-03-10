<template>
  <a-modal
    title="书签"
    v-model:visible="visible"
    width="60%"
  >
    <p>
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
    </div>

    <div v-else>
      <a-form-item label="目录" v-bind="validateInfos.parentId">
        <a-tree-select v-model:value="info.parentId" :tree-data="parentIdTreeData"/>
      </a-form-item>
      <a-alert show-icon v-if="info.title.length > 20" style="width:60%;text-align: center;margin-left: 10%;margin-bottom: 4px">
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
    </p>

    <template #footer>
      <a-button key="back" @click="visible = false">Close</a-button>
      <a-button key="submit" v-if="hasFetchURL" type="primary" :loading="addLoading" @click="handleAdd">Submit</a-button>
    </template>
  </a-modal>


</template>

<script setup>
import {reactive, ref, toRaw} from "vue"
import {apiJson} from "../../api";
import {toCateTree} from "@/utils/tree";
import {useBookmark} from "@/views/hook/bookmark";
import {useForm} from "ant-design-vue/es/form/index.js";
import {message} from "ant-design-vue";

const visible = ref(false)
const addVisible = ref(false)
const treeData = ref([])
const parentIdTreeData = ref([])
const bookmark = useBookmark()
const info = ref({
  parentId: '',
  title: '',
  groupId: '',
  url:'',
})

const hasFetchURL = ref(false)
const fetchLoading = ref(false)
function handleOk() {
  console.log("asd")
  visible.value = false
}

const {resetFields, validate, validateInfos} = useForm(
  info,
  reactive({
    parentId: [
      {
        required: true,
      },
    ],
    'title': [
      {
        required: true,
      },
    ],
    'url': [
  {
    required: true,
  },
],
  }),
);

function onSearch(){
  console.log("")
  if (!info.value.url.trim()){
    message.warn("url 不能为空")
    return
  }
  fetchLoading.value = true

  apiJson.get({
    "url":info.value.url,
    "meta()":`fetchURL(url)`
  }).then(data=>{
    console.log(data)
    info.value.url = data["meta"].url
    info.value.title = data["meta"].title
    info.value.icon = data["meta"].icon
    info.value.description = data["meta"].description
    hasFetchURL.value = true
  }).finally(()=>{
    fetchLoading.value = false
  })
}

const addLoading = ref(false)
function handleAdd() {
  addLoading.value = true
  let _info = toRaw(info.value)
  validate().then(data => {
    apiJson.post({
      Bookmark: _info,
      GroupBookmark: {
        cateId: _info.parentId,
        groupId: bookmark.curGroupId.value,
        "bmId@": "Bookmark/BmId",
      },
      tag: 'Bookmark',
    }).then(data => {
      console.log(data)
      visible.value = false
      bookmark.loadCate()
      bookmark.loadBookmarkList()
    })
  }).finally(()=>{
    addLoading.value = false
  })



}

function open() {
  visible.value = true
  info.value = {}
  hasFetchURL.value = false
  loadCate()
  resetFields()
}


function loadCate() {
  apiJson.get({
    "BookmarkCate[]": {
      "count":0,
    }
  }).then(data => {
    console.log(data)
    treeData.value = toCateTree((data["BookmarkCate[]"].map(item => {
      return {
        ...item,
        label: item.title,
        value: item.cateId,
      }
    })))
    parentIdTreeData.value = toCateTree([].concat(data["BookmarkCate[]"].map(item => {
      return {
        ...item,
        label: item.title,
        value: item.cateId,
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
