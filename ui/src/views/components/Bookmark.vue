<template>
  <div>
    <div class="bookmark cursor-pointer  rounded shadow p-1 items-baseline" :class="{'simple':simple}"
         @click="toURL(item)" v-contextmenu:contextmenu @contextmenu.stop="">
      <div class="flex">
        <div>
          <img v-lazy="item.icon" v-if="item.icon && item.icon !== ''" class="logo"/>
          <div v-else class="logo text-center text-white" :style="{'background': colorByURL(item.url)}">{{
              item.title[0]
            }}
          </div>
        </div>
        <div class="ml-1 title truncate ...">
          {{ item.title }}
        </div>
      </div>
      <div v-if="!simple" class="remark ..." style="height: 58px; margin-top: 4px">
        {{ item.remark || item.description }}
      </div>
      <div v-if="!simple" class="flex justify-end" style="height: 12px;line-height: 12px;color:#777" @click.stop="void">
        <a-popover :title="item.title">
          <template #content>
            <div style="width:300px">
              <p>{{ item.url }}</p>
              <p>{{ item.description }}</p>
            </div>
          </template>
          <div>
            ...
          </div>
        </a-popover>
      </div>
    </div>

    <v-contextmenu ref="contextmenu">
      <v-contextmenu-item @click="edit">🖊 编辑</v-contextmenu-item>
      <v-contextmenu-item @click="drop">❌ 删除</v-contextmenu-item>
    </v-contextmenu>
  </div>

</template>

<script setup>
import { apiJson } from '@/api/index'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { createVNode } from 'vue'
import { Modal } from 'ant-design-vue'
import dayjs from 'dayjs'
import { useBookmark } from '@/views/hook/bookmark'
import { colorByURL } from '../../utils/str-utils'

const bookmark = useBookmark()

const props = defineProps({
  item: {},
  simple: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['edit'])

function toURL (item) {
  apiJson.put({
    tag: 'BookmarkUse',
    BookmarkUse: {
      bmId: item.bmId
    }
  })
  window.open(item.url, '_blank')
}

function edit () {
  emit('edit')
}

function drop () {
  Modal.confirm({
    title: 'Do you Want to delete these items?',
    icon: createVNode(ExclamationCircleOutlined),
    content: createVNode('div', { style: 'color:red;' }, ''),
    onOk () {
      let cateId = bookmark.curSubCateId.value
      if (!cateId) {
        cateId = bookmark.curCate.value[bookmark.curCate.value.length - 1]
      }
      apiJson.put({
        tag: 'GroupBookmark',
        GroupBookmark: {
          bmId: props.item.bmId,
          groupId: bookmark.curGroupId.value,
          cateId,
          dropAt: dayjs().format('YYYY-MM-DD HH:mm:ss')
        }
      }).then(data => {
        bookmark.loadBookmarkList()
        bookmark.loadSubCateBookmark()
      })
    },
    onCancel () {
    }
  })
}


</script>

<style scoped lang="scss">
.bookmark {
  width: 200px;
  height: 106px;
  background: white;
  margin: 4px 4px;

  .logo {
    width: 24px;
    min-width: 24px;
    height: 24px;
    line-height: 24px;
    border-radius: 100%;
    text-transform: uppercase;
    font-weight: bold;
  }

  .title {
    font-size: 14px;
    font-weight: 600;
    color: #333;
    line-height: 24px;
  }

  &.simple {
    height: 60px;
    width: 240px;
    //line-height: 60px;

    .logo {
      width: 30px;
      min-width: 30px;
      height: 30px;
      line-height: 30px;
      font-size: 16px;
      text-transform: uppercase;
      font-weight: bold;
      margin-left: 4px;
      margin-top: 7px;
      border-radius: 100%;
    }

    .title {
      line-height: 50px;
      margin-left: 8px;
      font-size: 14px;
      font-weight: 600;
      color: #333;

    }
  }


  .remark {
    font-size: 12px;
    color: #acacac;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: break-spaces;
    line-height: 1.2;
  }

  &:hover {
    /*filter: drop-shadow(0 0 4px #42b883aa);*/
    filter: drop-shadow(0 0 4px #646cffaa);
  }
}

@media only screen and (max-width: 415px) {
  .bookmark {
    width: calc(50vw - 1rem);

    &.simple {
      width: calc(50vw - 1rem);

      .logo {
        width: 22px;
        min-width: 22px;
        height: 22px;
        line-height: 22px;
        font-size: 14px;
        text-transform: uppercase;
        font-weight: bold;
        margin-left: 4px;
        margin-top: 14px;
        border-radius: 100%;
      }

      .title {
        line-height: 50px;
        margin-left: 8px;
        font-size: 14px;
        font-weight: 600;
        color: #333;

      }
    }
  }
}

</style>
