<template>
  <div id="bg"></div>
  <div class="z-10 absolute w-full">
    <router-view/>
  </div>
  <Footer/>

</template>
<script setup>
import Footer from '@/views/layout/Footer.vue'

</script>

<style>

</style>
