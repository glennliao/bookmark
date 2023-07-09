<script setup lang="ts">
import { computed, ref } from 'vue'
import { RightOutlined } from '@ant-design/icons-vue'
import { apiJson } from '../api'

const version = import.meta.env.VITE_app_version.trim()
const latestVersion = ref('')

function checkNewVersion () {
  apiJson.get({
    'version()': 'latestVersion()'
  }).then(data => {
    latestVersion.value = data.version
  })
}

const canUpdate = computed(() => {
  return latestVersion.value > version
})

checkNewVersion()

</script>

<template>
  <footer class="version">
    <a class="mr-2" target="_blank" href="https://github.com/glennliao/bookmark">
      <svg t="1686376535417" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"
           p-id="2370" width="16" height="16">
        <path
          d="M512 12.64c-282.752 0-512 229.216-512 512 0 226.208 146.72 418.144 350.144 485.824 25.6 4.736 35.008-11.104 35.008-24.64 0-12.192-0.48-52.544-0.704-95.328-142.464 30.976-172.512-60.416-172.512-60.416-23.296-59.168-56.832-74.912-56.832-74.912-46.464-31.776 3.52-31.136 3.52-31.136 51.392 3.616 78.464 52.768 78.464 52.768 45.664 78.272 119.776 55.648 148.992 42.56 4.576-33.088 17.856-55.68 32.512-68.48-113.728-12.928-233.28-56.864-233.28-253.024 0-55.904 20-101.568 52.768-137.44-5.312-12.896-22.848-64.96 4.96-135.488 0 0 43.008-13.76 140.832 52.48 40.832-11.36 84.64-17.024 128.16-17.248 43.488 0.192 87.328 5.888 128.256 17.248 97.728-66.24 140.64-52.48 140.64-52.48 27.872 70.528 10.336 122.592 5.024 135.488 32.832 35.84 52.704 81.536 52.704 137.44 0 196.64-119.776 239.936-233.792 252.64 18.368 15.904 34.72 47.04 34.72 94.816 0 68.512-0.608 123.648-0.608 140.512 0 13.632 9.216 29.6 35.168 24.576 203.328-67.776 349.856-259.616 349.856-485.76 0-282.784-229.248-512-512-512z"
          fill="#444444" p-id="2371"></path>
      </svg>
    </a>

    <span>bookmark {{ version }}  </span>
    <span v-if="canUpdate">
      <right-outlined class="ml-1" style="font-size: 12px"/>
      <a-badge dot>
    <a target="_blank" href="https://github.com/glennliao/bookmark" class="ml-1">{{ latestVersion }}</a>
  </a-badge>

    </span>
  </footer>
</template>

<style scoped>
.version {
  height: 26px;
  color: #a1a1a1;
  text-align: center;
  font-size: 14px;
  z-index: 10;
}
</style>
