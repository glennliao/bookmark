<script setup lang="ts">
import { ref, type Ref, type UnwrapRef } from 'vue'
import { PlusOutlined, DeleteOutlined, EditOutlined } from '@ant-design/icons-vue'
import EditModal from './components/EditModal.vue'
import Window from './components/Window/index.vue'
import TitleIcon from './components/TitleIcon.vue'
import { useDesktop } from './useDesktop'
import { useAppList, AppItem } from './useAppList'

const appList = useAppList()
const desktop = useDesktop()

appList.loadList()

const open = (item: AppItem) => {
  if ((item.openType || '') === '_blank') {
    window.open(item.url, '_blank')
    return
  }

  const openParam = JSON.parse(item.openParam || JSON.stringify({
    width: 500,
    height: 400
  }))

  desktop.fork({
    url: item.url,
    title: item.title,
    icon: item.icon,
    appId: item.appId,
    windowStyle: {
      ...openParam
    }
  })

}

const editModal: Ref<UnwrapRef<null>> = ref(null)
const add = () => {
  editModal.value && editModal.value.open()
}

const edit = (item) => {
  editModal.value.open(item)
}


</script>

<template>
  <div class="desktop">
    <div class="app-list">
      <a-dropdown v-for="item in appList.list.value" :key="item.url" :trigger="['contextmenu']">
        <div class="app-list-item" @click="open(item)">
          <div>
            <img :src="item.icon" alt=""/>
          </div>
          <div class="title">
            {{ item.title }}
          </div>
        </div>

        <template #overlay>
          <a-menu>
            <a-menu-item key="edit" @click="edit(item)">
              <EditOutlined/>
              edit
            </a-menu-item>
            <a-menu-item key="del" @click="appList.del(item.appId)">
              <DeleteOutlined/>
              delete
            </a-menu-item>

          </a-menu>
        </template>
      </a-dropdown>
      <div class="app-list-item" style="border:1px dashed #b9b9b9;line-height: 48px" @click="add">
        <div>
          <PlusOutlined/>
        </div>
        <div class="title">

        </div>
      </div>
    </div>

    <div class="task-bar">
      <div v-for="process in desktop.processList.value" :key="process.pid" class="task-bar-item"
           :class="{active:process.focus}">
        <TitleIcon :title="process.title" :icon="process.icon"
                   @click="desktop.action(process.pid,'toggle')"></TitleIcon>
      </div>
    </div>

    <Window v-for="process in desktop.processList.value"
            :key="process.pid" :title="process.title"
            v-show="!process.hidden"
            @mousedown="desktop.action(process.pid,'focus')"
            @action="desktop.action(process.pid,$event.action)"
            :windowStyle="process.windowStyle" style="position: fixed">
      <iframe :src="process.url" class="w-full h-full"></iframe>
    </Window>

    <EditModal ref="editModal" @refresh="appList.loadList()"/>


  </div>
</template>

<style scoped lang="less">
.desktop {
  height: 100%;
  width: 100%;
  overflow: hidden;
  position: relative;
}

.app-list {
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: flex-start;
  align-content: flex-start;
  flex-wrap: wrap;
  padding: 10px;
  border-radius: 6px;
  transition: all .35s ease-in-out;
  outline: 1px solid transparent;
  height: 100vh;

  .app-list-item {
    width: 64px;
    height: 64px;
    //background: #fff;
    //box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    margin-bottom: 26px;
    margin-right: 26px;
    cursor: pointer;
    user-select: none;
    border-radius: 8px;
    text-align: center;
    text-overflow: ellipsis;
    word-break: break-all;
    font-size: 13px;
    color: white;
    padding: 8px;

    img {
      width: 32px;
      height: 32px;
      margin-bottom: 8px;
    }
  }
}

.task-bar {
  position: fixed;
  display: flex;
  justify-content: center;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 10002;

  .task-bar-item {
    width: 44px;
    height: 34px;
    background: rgba(22, 119, 255, 0.16);
    padding: 2px 6px 0 6px;
    margin-right: 2px;
    cursor: pointer;

    &.active {
      background-color: #009688;
    }
  }

}

</style>
