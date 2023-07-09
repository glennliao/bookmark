<template>
  <div id="bg"></div>
  <div  class="body z-10 w-full">
    <router-view/>
  </div>
  <Footer/>

</template>
<script setup>
import Footer from '@/layout/Footer.vue'

</script>

<style>
.body{
  min-height: calc(100vh - 26px);
  font-size: 14px;
  position: relative;
}
</style>
