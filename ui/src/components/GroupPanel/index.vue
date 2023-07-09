<template>
  <div class="toggle-bar" @click="groupsVisible = !groupsVisible">
    <div class="n-layout-toggle-bar">
      <div class="n-layout-toggle-bar__top"></div>
      <div class="n-l  ayout-toggle-bar__bottom"></div>
    </div>
  </div>

  <a-drawer
    title="Groups"
    placement="left"
    :closable="true"
    v-model:open="groupsVisible"
    width="180px"
  >
    <div @click="changeGroup(item)" v-for="item in groups" style="cursor: pointer"
         class="shadow p-2 rounded bg-orange-500 text-white" :key="item.groupId">
      {{ item.title }}
    </div>
  </a-drawer>
</template>


<script setup>

import { ref } from 'vue'
import { useBookmark } from '@/views/bookmark/hook/bookmark.js'

const {

  curGroupId,
  groups
} = useBookmark()
const groupsVisible = ref(false)

function changeGroup (item) {
  curGroupId.value = item.groupId
  // loadCate()
  groupsVisible.value = false
}

</script>

<style>

.n-layout-toggle-bar:hover .n-layout-toggle-bar__top {
  transform: rotate(-12deg) scale(1.15) translateY(-2px);
}

.n-layout-toggle-bar:hover .n-layout-toggle-bar__bottom {
  transform: rotate(12deg) scale(1.15) translateY(2px);
}

.n-layout-toggle-bar {

//left: -28px;
  transform: rotate(180deg);

}

.n-layout-toggle-bar:hover .n-layout-toggle-bar__top {
  transform: rotate(12deg) scale(1.15) translateY(-2px);
}

.n-layout-toggle-bar:hover .n-layout-toggle-bar__bottom {
  transform: rotate(-12deg) scale(1.15) translateY(2px);
}

.n-layout-toggle-bar:hover .n-layout-toggle-bar__top {
  transform: rotate(-12deg) scale(1.15) translateY(-2px);
}

.n-layout-toggle-bar:hover .n-layout-toggle-bar__bottom {
  transform: rotate(12deg) scale(1.15) translateY(2px);
}

.toggle-bar {
  --n-bezier: cubic-bezier(0.4, 0, 0.2, 1);
  --n-toggle-button-color: #FFF;
  --n-toggle-button-border: 1px solid rgb(239, 239, 245);
  --n-toggle-bar-color: rgba(191, 191, 191, 1);
  --n-toggle-bar-color-hover: rgba(153, 153, 153, 1);
  --n-color: #fff;
  --n-text-color: rgb(51, 54, 57);
  --n-border-color: rgb(239, 239, 245);
  --n-toggle-button-icon-color: rgb(51, 54, 57);
  position: absolute;
  left: -4px;
  top: 0;
  width: 12px;
  bottom: 0;
}

.n-layout-toggle-bar {

  cursor: pointer;
  height: 72px;
  width: 32px;
  position: absolute;
  top: calc(50% - 36px);
//right: -28px;

}

.n-layout-toggle-bar .n-layout-toggle-bar__top, .n-layout-toggle-bar .n-layout-toggle-bar__bottom {

  position: absolute;
  width: 4px;
  border-radius: 2px;
  height: 38px;
  left: 14px;
  transition: background-color .3s var(--n-bezier),
  transform .3s var(--n-bezier);

}

.n-layout-toggle-bar .n-layout-toggle-bar__bottom {

  position: absolute;
  top: 34px;

}

.n-layout-toggle-bar:hover .n-layout-toggle-bar__top {
  transform: rotate(12deg) scale(1.15) translateY(-2px);
}

.n-layout-toggle-bar:hover .n-layout-toggle-bar__bottom {
  transform: rotate(-12deg) scale(1.15) translateY(2px);
}

.n-layout-toggle-bar .n-layout-toggle-bar__top, .n-layout-toggle-bar .n-layout-toggle-bar__bottom {
  background-color: var(--n-toggle-bar-color);
}

.n-layout-toggle-bar:hover .n-layout-toggle-bar__top, .n-layout-toggle-bar:hover .n-layout-toggle-bar__bottom {
  background-color: var(--n-toggle-bar-color-hover);
}
</style>
