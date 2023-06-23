<template>
  <div>
    <div>
      <div>
        <a-button type="primary" @click="add">新增</a-button>
      </div>
      <div class="mt-4">
        <a-tree
          :tree-data="treeData"
          draggable
          show-line
          :show-icon="false"
          block-node
          @drop="drop"
        >
          <template #title="{key, title, parentId}">
            <div style="height: 36px" class="flex justify-between items-center">
              <div>
                {{ title }}
              </div>
              <div>
                <a-button @click="edit(key,title,parentId)" type="text">编辑</a-button>
                <a-button @click="del(key,title)" type="text">删除</a-button>
              </div>
            </div>
          </template>
        </a-tree>
      </div>
    </div>
    <a-modal
      title="编辑"
      v-model:open="addVisible"
      @ok="handleAdd"
      width="460px"
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
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, toRaw } from 'vue'
import { apiJson } from '@/api'
import { toCateTree } from '@/utils/tree'
import { useBookmark } from '@/views/hook/bookmark'
import { useForm } from 'ant-design-vue/es/form/index.js'
import { message, Modal, TreeDataItem, DropEvent } from 'ant-design-vue'

const addVisible = ref(false)
const treeData = ref([] as any[])
const parentIdTreeData = ref([] as any[])
const bookmark = useBookmark()
const info = ref({
  parentId: '',
  title: '',
  groupId: '',
  cateId: ''
})

const {
  resetFields,
  validate,
  validateInfos
} = useForm(
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
    ]
  })
)

function handleAdd () {
  validate().then(data => {
    const _info = toRaw(info.value)
    _info.groupId = bookmark.curGroupId.value

    if (_info.cateId) {
      apiJson.put({
        BookmarkCate: _info,
        tag: 'BookmarkCate'
      }).then(data => {
        console.log(data)
        addVisible.value = false
        loadCate()
        bookmark.loadCate()
      })
      return
    }
    apiJson.post({
      BookmarkCate: _info,
      tag: 'BookmarkCate'
    }).then(data => {
      console.log(data)

      addVisible.value = false
      loadCate()
      bookmark.loadCate()
    })
  })
}

onMounted(() => {
  loadCate()
})

function loadCate () {
  apiJson.get({
    'BookmarkCate[]': {
      count: 0,
      '@alias': 'list'
    }
  }).then(data => {
    treeData.value = toCateTree((data.list.map(item => {
      return {
        ...item,
        label: item.title,
        key: item.cateId,
        value: item.cateId
      }
    })))
    parentIdTreeData.value = toCateTree([{
      value: 'root',
      label: 'root',
      parentId: 'root'
    }].concat(data.list.map(item => {
      return {
        ...item,
        label: item.title,
        key: item.cateId,
        value: item.cateId
      }
    })))
  })
}

function add () {
  resetFields()
  info.value = {
    parentId: '',
    title: '',
    cateId: ''
  }
  addVisible.value = true
}

const rules = {}

function del (key: any, title: any) {
  Modal.confirm({
    title: '删除',
    content: `确认删除 ${title} 吗？`,
    onOk: () => {
      apiJson.delete({
        BookmarkCate: {
          cateId: key
        },
        tag: 'BookmarkCate'
      }).then(data => {
        message.success('删除成功')
        loadCate()
        bookmark.loadCate()
      })
    }
  })
}

function edit (key: any, title: any, parentId: any) {
  info.value = {
    parentId,
    title,
    cateId: key
  }
  addVisible.value = true
}

function drop (info) {
  console.log('dragend', info)
  // 目标节点
  const dropKey = info.node.eventKey
  // 拖动节点
  const dragKey = info.dragNode.eventKey
  const dropPos = info.node.pos.split('-')
  // -1 上
  // 0 内
  // 1 下
  // 拖动的位置
  const dropPosition =
    info.dropPosition - Number(dropPos[dropPos.length - 1])

  const loop = (data: TreeDataItem[], key: string, callback: any) => {
    data.forEach((item, index, arr) => {
      if (item.key === key) {
        return callback(item, index, arr)
      }
      if (item.children) {
        return loop(item.children, key, callback)
      }
    })
  }
  const data = [...treeData.value]

  // 拖动节点的对象
  let dragObj: TreeDataItem = {}

  if (!info.dropToGap) {
    return false
  } else if (
    (info.node.children || []).length > 0 &&
    info.node.expanded &&
    dropPosition === 1
  ) {
    return false
  } else {
    let a: TreeDataItem[] = []
    let ii = 0
    loop(
      data,
      dragKey,
      (item: TreeDataItem, index: number, arr: TreeDataItem[]) => {
        a = arr
        ii = index
        dragObj = item
      }
    )
    // 只允许当前节点下的子节点拖动排序
    if (a.some((item) => item.key === dropKey)) {
      a.splice(ii, 1)
      let ar: TreeDataItem[] = []
      let i = 0
      loop(
        data,
        dropKey,
        (item: TreeDataItem, index: number, arr: TreeDataItem[]) => {
          ar = arr
          i = index
        }
      )
      if (dropPosition === -1) {
        ar.splice(i, 0, dragObj)
      } else {
        ar.splice(i + 1, 0, dragObj)
      }
    }
  }

  console.log(data)
}

</script>

<style scoped>

</style>
