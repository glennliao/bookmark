<template>
  <a-sub-menu :key="menuInfo.cateId" @titleClick="click([])">
    <template #title ><span >{{ menuInfo.title }}</span></template>
    <template v-for="item in menuInfo.children" :key="item.cateId">
      <template v-if="!item.children">
        <a-menu-item :key="item.cateId">
          <div @click="click([item.cateId])">
            {{ item.title }}
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
