import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

import vueJsx from '@vitejs/plugin-vue-jsx'
import vueSetupExtend from 'vite-plugin-vue-setup-extend'
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'
import Components from 'unplugin-vue-components/vite'

import { execSync } from 'child_process'

let tag = 'dev'
try {
  tag = execSync('git describe --abbrev=10 --tags').toString()
} catch (e) {
  console.log(e.toString())
}
// @ts-ignore
process.env.VITE_app_version = tag

// https://vitejs.dev/config/
export default defineConfig({
  base: './',
  plugins: [
    vue(),
    vueSetupExtend(),
    vueJsx(),
    // eslint(),
    createSvgIconsPlugin({
      // 要缓存的图标文件夹
      iconDirs: [path.resolve(__dirname, 'src/svg')],
      // 执行 icon name 的格式
      symbolId: 'icon-[name]'
    }),
    Components({
      resolvers: [
        // AntDesignVueResolver(),
        // NaiveUiResolver()
      ]
    })
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src')
    }
  },
  server: {
    port: 3000,
    cors: true,

    proxy: {
      '/api': {
        target: 'http://127.0.0.1:8082',
        changeOrigin: true

        // rewrite: (path) => path.replace(/^\/api/, '')
      }

    }
  },
  envDir: path.resolve(__dirname, './env')
})
