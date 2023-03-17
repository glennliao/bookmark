<template>
  <div id="bg"></div>
  <div class="z-10 absolute">
    <router-view/>
  </div>

  <div style="color:#a1a1a1;position: absolute;bottom:4px;text-align: center;width: 100vw">
    {{version}}
  </div>
</template>
<script setup>
const version = import.meta.env.VITE_app_version.trim()
</script>
