<template>
  <a-sub-menu :key="menuInfo.cateId">
<!--    <template #icon><MailOutlined /></template>-->
    <template #title ><span>{{ menuInfo.title }}</span></template>
    <template v-for="item in menuInfo.children" :key="item.cateId">
      <template v-if="!item.children">
        <a-menu-item :key="menuInfo.cateId+'/'+item.cateId">
<!--          <template #icon>-->
<!--            <PieChartOutlined />-->
<!--          </template>-->
          <div >
            {{ item.title }}
          </div>
        </a-menu-item>
      </template>
      <template v-else>
        <SubMemu :menu-info="item" :key="item.cateId" />
      </template>
    </template>
  </a-sub-menu>
</template>

<script>
export default {
  name: "SubMemu",
  props: {
    menuInfo: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(){
    return {
    }
  }
}
</script>

<style scoped>

</style>
