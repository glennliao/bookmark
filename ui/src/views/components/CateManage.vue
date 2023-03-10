<template>
  <a-modal
    title="目录管理"
    v-model:visible="visible"
    width="70%"
    @ok="handleOk"
  >
    <p>
      <div>
        <a-button type="primary" @click="add">新增</a-button>
      </div>
      <div class="mt-4">
        <a-tree
          :tree-data="treeData"
        />
      </div>
    </p>
  </a-modal>

  <a-modal
    title="编辑"
    v-model:visible="addVisible"
    @ok="handleAdd"
    width="60%"
  >
    <p>
      <a-form autocomplete="off"
              :rules="rules">
        <a-form-item required label="上级目录" v-bind="validateInfos.parentId">
          <a-tree-select v-model:value="info.parentId" :tree-data="parentIdTreeData"/>
        </a-form-item>
        <a-form-item required label="目录名称" v-bind="validateInfos.title">
          <a-input v-model:value="info.title"/>
        </a-form-item>
      </a-form>
    </p>
  </a-modal>

</template>

<script setup>
import {reactive, ref, toRaw} from "vue"
import {apiJson} from "../../api";
import {toCateTree} from "@/utils/tree";
import {useBookmark} from "@/views/hook/bookmark";
import {useForm} from "ant-design-vue/es/form/index.js";

const visible = ref(false)
const addVisible = ref(false)
const treeData = ref([])
const parentIdTreeData = ref([])
const bookmark = useBookmark()
const info = ref({
  parentId: '',
  title: '',
  groupId: '',
})

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
  }),
);

function handleAdd() {
  console.log("asd")
  console.log(info.value)

  validate().then(data => {
    let _info = toRaw(info.value)
    _info.groupId = bookmark.curGroupId.value

    apiJson.post({
      BookmarkCate: _info,
      tag: 'BookmarkCate',
    }).then(data => {
      console.log(data)

      addVisible.value = false
      loadCate()
      bookmark.loadCate()
    })
  })


}


function open() {
  visible.value = true
  loadCate()
}


function loadCate() {
  apiJson.get({
    "BookmarkCate[]": {
      count: 0,
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
    parentIdTreeData.value = toCateTree([{
      value: 'root',
      label: 'root',
      parentId: 'root'
    }].concat(data["BookmarkCate[]"].map(item => {
      return {
        ...item,
        label: item.title,
        value: item.cateId,
      }
    })))
  })
}

function add() {
  resetFields()
  addVisible.value = true
}

const rules = {}


defineExpose({
  open
})

</script>

<style scoped>

</style>
