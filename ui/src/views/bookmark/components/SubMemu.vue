<template>
  <a-sub-menu :key="menuInfo.cateId" @titleClick="click([])">
    <template #title >
      <span>{{ menuInfo.title }}  <span v-if="menuInfo.count > 0" style="font-size: 12px;margin-left: 3px;color:#8b8b8b">({{menuInfo.count}})</span></span>

<!--      <span >{{ menuInfo.title }}</span>-->
    </template>
    <template v-for="item in menuInfo.children" :key="item.cateId">
      <template v-if="!item.children">
        <a-menu-item :key="item.cateId">
          <div class="flex items-end" @click="click([item.cateId])">
            <span>{{ item.title }}</span>
            <span v-if="item.count > 0" style="font-size: 12px;margin-left: 3px;color:#8b8b8b">({{item.count}})</span>
          </div>
        </a-menu-item>
      </template>
      <template v-else>
        <SubMemu :menu-info="item" @click="click"  :key="item.cateId" />
      </template>
    </template>
  </a-sub-menu>
</template>

<script setup>

const props = defineProps({
  menuInfo: {
    type: Object,
    default: () => ({})
  }
})


const emit = defineEmits(['click'])

function click (keyList) {
  keyList = [props.menuInfo.cateId,...keyList]
  console.log("keyList", keyList)
  emit('click', keyList)
}
</script>

<style scoped>

</style>
