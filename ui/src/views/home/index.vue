<script setup lang="ts">
import { ref, type Ref, type UnwrapRef } from 'vue'
import { CaretRightOutlined, PlusOutlined, DeleteOutlined } from '@ant-design/icons-vue'
import { AppItem } from '@/views/home/types'
import { apiJson } from '@/api'
import EditModal from '@/views/home/components/EditModal.vue'
import Win from "@/views/home/win.vue";

const appList: Ref<AppItem[]> = ref([])

const loadList = () => {
  apiJson.get({
    'AppList[]': {
      '@alias': 'list'
    }
  }).then(data => {
    appList.value = data.list
  })
}

loadList()

const to = (item: AppItem) => {
  window.open(item.url, '_blank')

  return

  const taskItem = winTasks.value.find(taskItem=>taskItem.appId === item.appId)

  if(taskItem){
    taskItem.zIndex = 101
    let preTaskItem = winTasks.value.find(taskItem=>taskItem.zIndex === 101)
    if(preTaskItem){
      preTaskItem.zIndex = 100
    }
  }else{
    winTasks.value.push({
      appId: item.appId,
      url: item.url,
      zIndex: 100,
    })
  }
}

const editModal: Ref<UnwrapRef<null>> = ref(null)
const add = () => {
  editModal.value && editModal.value.open()
}

const del = (appId: string) => {
  apiJson.delete({
    AppList: {
      appId
    },
    tag: 'AppList'
  }).then(data => {
    loadList()
  })
}

const winTasks = ref([])
</script>

<template>
  <div class="home">
    <div class="apps">
      <a-dropdown v-for="item in appList" :key="item.url" :trigger="['contextmenu']">
        <div class="item" @click="to(item)">

          <div>
            <img :src="item.icon" alt=""/>
          </div>
          <div class="title">
            {{ item.title }}
          </div>
          <div class="to-icon">
            <CaretRightOutlined/>
          </div>
        </div>

        <template #overlay>
          <a-menu>
            <a-menu-item key="del" @click="del(item.appId)">
              <DeleteOutlined/>
              delete
            </a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>

    </div>

    <win v-for="(task,index) in winTasks" :key="index" >
      <iframe :src="task.url" style="width:100%;height: 100%"/>
    </win>

    <EditModal ref="editModal" @refresh="loadList()"/>

    <a-float-button-group shape="square" :style="{ right: '24px' }">
      <a-float-button type="primary" @click="add">
        <template #icon>
          <PlusOutlined/>
        </template>
      </a-float-button>
    </a-float-button-group>

  </div>
</template>

<style lang="scss" scoped>

.apps {
  display: flex;
  flex-wrap: wrap;
  padding: 200px;

  .item {
    width: 280px;
    height: 90px;
    flex: 0 0 280px;
    border-radius: 6px;
    padding: 15px 55px 15px 15px;
    color: #000;
    overflow: hidden;
    position: relative;
    transition: all .35s ease-in-out;
    outline: 1px solid transparent;
    display: flex;
    align-items: center;
    border: 1px solid rgba(76, 76, 76, 0.4);
    background-clip: padding-box;
    background: #fff linear-gradient(90deg, rgba(255, 255, 255, 0), rgba(255, 255, 255, 0.25));
    margin: 12px 0 20px 12px;
    cursor: pointer;

    img {
      width: 60px;
      height: 60px;
    }

    .title {
      margin-left: 12px;
      font-size: 16px;
      font-weight: bold;
    }

    .to-icon {
      position: absolute;
      right: 6px;
      font-size: 24px;
      top: calc(50% - 12px);
    }

    &:after {
      content: "";
      position: absolute;
      width: 90px;
      height: 90px;
      border-radius: 50%;
      right: -48px;
      top: 0px;
      background: rgba(255, 255, 255, 0.1);
      box-shadow: 0 0 40px 0px rgba(0, 0, 0, 0.2);
    }

    &:hover {
      box-shadow: 0 0 20px 0 rgba(255, 255, 255, 0.6);
    }
  }
}
</style>
