import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'
import eslint from 'vite-plugin-eslint'
import vueJsx from '@vitejs/plugin-vue-jsx'
import vueSetupExtend from 'vite-plugin-vue-setup-extend'
import {createSvgIconsPlugin} from 'vite-plugin-svg-icons'
import Components from 'unplugin-vue-components/vite'
import {
  AntDesignVueResolver,NaiveUiResolver
} from 'unplugin-vue-components/resolvers'

// https://vitejs.dev/config/
export default defineConfig({
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
        AntDesignVueResolver(),NaiveUiResolver()
      ],
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
        rewrite: (path) => path.replace(/^\/api/, '')
      }
    }
  },
  envDir: path.resolve(__dirname, './env')
})
