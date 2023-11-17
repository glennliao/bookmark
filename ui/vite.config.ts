import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

import vueJsx from '@vitejs/plugin-vue-jsx'
import vueSetupExtend from 'vite-plugin-vue-setup-extend'
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'
import Components from 'unplugin-vue-components/vite'
import AutoImport from 'unplugin-auto-import/vite'
import { AntDesignVueResolver } from 'unplugin-vue-components/resolvers'
import { execSync } from 'child_process'
import UnoCSS from 'unocss/vite'

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
    UnoCSS(),
    // eslint(),
    createSvgIconsPlugin({
      // 要缓存的图标文件夹
      iconDirs: [path.resolve(__dirname, 'src/svg')],
      // 执行 icon name 的格式
      symbolId: 'icon-[name]'
    }),
    AutoImport({
      // targets to transform
      include: [
        /\.[tj]sx?$/,
        /\.vue$/,
        /\.vue\?vue/,
        /\.md$/
      ],

      // global imports to register
      imports: [
        'vue',
        'vue-router',
        'pinia'
      ],

      // Generate corresponding .eslintrc-auto-import.json file.
      // eslint globals Docs - https://eslint.org/docs/user-guide/configuring/language-options#specifying-globals
      eslintrc: {
        enabled: false,
        filepath: './.eslintrc-auto-import.json', // Default `./.eslintrc-auto-import.json`
        globalsPropValue: true // Default `true`, (true | false | 'readonly' | 'readable' | 'writable' | 'writeable')
      },
      dts: './auto-imports.d.ts'
    }),
    Components({
      resolvers: [
        AntDesignVueResolver({
          importStyle: false // css in js
        })
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
    https: false,
    host: '0.0.0.0',

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
