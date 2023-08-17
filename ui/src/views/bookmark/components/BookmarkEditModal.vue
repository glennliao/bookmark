<template>
  <a-modal
    title="书签"
    v-model:open="visible"

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

          <div v-if="equalBm.bmId" class="mt-4">
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
          <a-form-item label="目录" v-bind="validateInfos.cateId">
            <a-tree-select block v-model:value="info.cateId" :tree-data="parentIdTreeData"/>
          </a-form-item>
          <a-alert show-icon v-if="info.title && info.title.length > 20"
                   style="width:60%;text-align: center;margin-left: 10%;margin-bottom: 4px">
            <template #message>标题过长, 不利于查看</template>
          </a-alert>
          <a-form-item label="名称" v-bind="validateInfos.title">
            <a-input v-model:value="info.title"/>
          </a-form-item>
          <a-form-item label="地址" v-bind="validateInfos.url">
            <a-input v-model:value="info.url"/>
          </a-form-item>
          <a-form-item label="图标" help="图标为空时将显示标题的首个字符" v-bind="validateInfos.icon">
            <a-input v-model:value="info.icon">
              <template #addonAfter>

                <div>
                  <a-upload-dragger
                    name="file"
                    :showUploadList="false"
                    :multiple="false"
                    action="./api/upload"
                    :headers="uploadHeader"
                    @change="onFileChange"
                  >
                    <img v-if="info.icon" :src="info.icon" style="width: 20px;height: 20px;max-width: 20px"/>
                    <div v-else class="logo text-center text-white" :style="{'background': colorByURL(info.url)}">
                      {{ info.title[0] }}
                    </div>
                  </a-upload-dragger>
                </div>
              </template>
            </a-input>
          </a-form-item>
          <a-form-item label="描述" v-bind="validateInfos.description">
            <a-textarea :rows="3" v-model:value="info.description" type="textarea"/>
          </a-form-item>
          <a-form-item label="备注" v-bind="validateInfos.remark">
            <a-textarea :rows="3" v-model:value="info.remark" type="textarea"/>
          </a-form-item>
          <a-form-item label="标签" v-bind="validateInfos.tags">
            <a-select
                v-model:value="info.tags"
                mode="tags"
                style="width: 100%"
                placeholder="Tags"
                :options="tagsOptions"
                @change="handleChange"
            ></a-select>
          </a-form-item>
        </div>
      </a-form>
    </div>

    <template #footer>
      <div class="flex justify-between">
        <div>
          <a-popover title="">
            <template #content>
              <div>
                1. 将右侧拖动到浏览器书签 (推荐放在第一位)<a class="ml-2" :href="saveBookJS"
                                                             title="保存书签">保存书签</a>
              </div>
              <div>
                2. 如遇到需要保存到此处的书签,点击第一步保存的书签即可
              </div>
              <div>
                <img src="../../../assets/add_bookmark.gif"/>
              </div>
            </template>
            从书签栏添加
          </a-popover>
        </div>
        <div>

          <a-button key="back" @click="visible = false" class="mr-2">Close</a-button>

          <a-dropdown-button type="primary" v-if="hasFetchURL" :loading="addLoading" @click="handleAdd(false)">
            Submit
            <template #overlay>
              <a-menu>
                <a-menu-item v-if="hasFetchURL && !info.bmId" key="1" @click="handleAdd(true)">Submit & Next
                </a-menu-item>
              </a-menu>
            </template>
            <template #icon>
              <DownOutlined/>
            </template>
          </a-dropdown-button>
        </div>
      </div>
    </template>
  </a-modal>

</template>

<script setup>
import { createVNode, reactive, ref, toRaw } from 'vue'
import { apiJson } from '../../../api/index'
import { toCateTree } from '@/utils/tree'
import { useBookmark } from '@/views/bookmark/hook/bookmark'
import { useForm } from 'ant-design-vue/es/form/index'
import { message, Modal } from 'ant-design-vue'
import { ExclamationCircleOutlined, DownOutlined } from '@ant-design/icons-vue'
import { colorByURL } from '@/utils/str-utils'
import useUserStore from '@/store/modules/user'


const visible = ref(false)
const addVisible = ref(false)
const treeData = ref([])
const parentIdTreeData = ref([])
const bookmark = useBookmark()
const info = ref({
  cateId: '',
  title: '',
  groupId: '',
  url: '',
  tags:[]
})

const hasFetchURL = ref(false)
const fetchLoading = ref(false)

const equalBm = ref({})

const {
  resetFields,
  validate,
  validateInfos
} = useForm(
  info,
  reactive({
    cateId: [
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

function onSearch () {
  const url = info.value.url.trim()

  if (!url) {
    message.warn('url 不能为空')
    return
  }
  fetchLoading.value = true

  const param = {
    url,
    'meta()': 'fetchURL(url)',
    equal: {
      Bookmark: {
        url // or suffix /
      },
      BookmarkCate: {
        'cateId@': '/Bookmark/cateId'
      }
    }
  }

  try {
    const host = (new URL(url)).host
    param.domain = {
      'Bookmark[]': {
        url$: '%' + host + '%'
      }
    }
  } catch (e) {

  }

  apiJson.get(param).then(data => {
    const equal = data.equal.Bookmark
    if (equal.bmId) {
      equal.cate = data.equal.BookmarkCate.title
      equalBm.value = equal
    } else {
      hasFetchURL.value = true
    }

    if (bookmark.curSubCateId.value) {
      info.value.cateId = bookmark.curSubCateId.value || ''
    } else {
      if (bookmark.curCate.value.length) {
        info.value.cateId = bookmark.curCate.value[bookmark.curCate.value.length - 1] || ''
      }
    }
    info.value.url = data.meta.url
    info.value.title = data.meta.title
    info.value.icon = data.meta.icon
    info.value.description = data.meta.description
    info.value.tags = []
  }).catch(() => {
    Modal.confirm({
      title: '操作确认',
      icon: createVNode(ExclamationCircleOutlined),
      content: '获取网址信息错误, 是否手动录入',
      okText: '确认',
      cancelText: '取消',
      onOk: () => {
        info.value.url = url
        info.value.title = ''
        hasFetchURL.value = true
      }
    })
  }).finally(() => {
    fetchLoading.value = false
  })
}

const addLoading = ref(false)

function handleAdd (next) {
  addLoading.value = true
  const _info = toRaw(info.value)

  if (!_info.cateId) {
    message.warn('请选择分类')
    addLoading.value = false
    return
  }

  validate().then(data => {
    let api = apiJson.post
    _info.groupId = bookmark.curGroupId.value

    if (!_info.groupId) {
      message.warn('groupId丢失，请刷新页面')
      addLoading.value = false
      return
    }

    if (_info.bmId) {
      api = apiJson.put
    }
    api({
      Bookmark: _info,
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
      if (bookmark.curSubCateId.value) {
        bookmark.loadSubCateBookmark()
      }
    })
  }).finally(() => {
    addLoading.value = false
  })
}

function open (_info = {}) {
  resetFields()
  equalBm.value = {}
  visible.value = true
  info.value = {
    ..._info
  }

  if (_info.url && !_info.bmId) {
    setTimeout(() => {
      onSearch()
    }, 512)
  }

  hasFetchURL.value = _info.bmId
  loadCate()
  loadTagList()
}

function loadCate () {
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

const user = useUserStore()
const uploadHeader = {
  Authorization: user.token
}

function onFileChange (e) {
  if (e.file.response) {
    const res = e.file.response
    if (res.code === 200) {
      info.value.icon = './api/icon?name=' + res.data.path
    } else {
      message.warn(res.msg)
    }
  }
  console.log(e)
}

function add () {
  addVisible.value = true
}

const saveBookJS = `javascript:window.open('${window.location.href.split('#')[0]}#/?url='+window.location,'_blank')`


const tagsOptions = ref([])

function loadTagList() {
  apiJson.get({
    'tags()': 'bmTags()'
  }).then(data => {
    tagsOptions.value = data.tags.map(item=>{
      return {
        value:item
      }
    })
  })
}

function handleChange(e){
  console.log(e)
}


defineExpose({
  open
})

</script>

<style scoped>
.logo {
  width: 22px;
  min-width: 22px;
  height: 22px;
  line-height: 22px;
  font-size: 14px;
  text-transform: uppercase;
  font-weight: bold;
  border-radius: 100%;
}

:deep(.ant-form-item) {
  margin-bottom: 16px;
}

:deep(.ant-upload-wrapper .ant-upload-drag) {
  border: none;
}

:deep(.ant-upload-wrapper .ant-upload-drag :hover) {
  border: 1px dashed #d9d9d9;
  border-radius: 8px;
}

:deep(.ant-upload-wrapper .ant-upload-drag .ant-upload) {
  padding: 0;
}

@media (max-width: 575px) {

  :deep(.ant-form-item .ant-form-item-label) {
    padding: 0;
  }
}
</style>
