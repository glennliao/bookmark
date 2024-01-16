<script setup lang="ts">
defineProps({
  title: {
    type: String,
    default: ''
  },
  icon: {
    type: String,
    default: ''
  }
})
const titleIconText = (title: string) => {
  return title[0].toUpperCase()
}
</script>

<template>
  <div class="w-full h-full text-center" >
    <img v-if="icon" :src="icon" style="height: 32px;width: 32px">
    <span v-else style="font-size: 20px;color: white">
      {{ titleIconText(title) }}
    </span>
  </div>
</template>

<style scoped>

</style>
