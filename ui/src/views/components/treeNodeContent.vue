<template>
  <span v-contextmenu:contextmenu >
    <slot/>
  </span>
  <v-contextmenu ref="contextmenu">
    <v-contextmenu-item @click="onContextmenu('add')">管理</v-contextmenu-item>
<!--    <v-contextmenu-item @click="onContextmenu('update')">编辑</v-contextmenu-item>-->
<!--    <v-contextmenu-item @click="onContextmenu('delete')">删除</v-contextmenu-item>-->
  </v-contextmenu>
</template>

<script setup>
const emit = defineEmits(['contextmenu','click'])

function onContextmenu(action){
  emit("contextmenu",action)
}
</script>

<style scoped>

</style>
