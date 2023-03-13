<template>
  <div class="bookmark cursor-pointer  rounded shadow p-1 items-baseline" :class="{'simple':simple}" @click="toURL(item)">
    <div class="flex">
      <div>
        <img  v-lazy="item.icon" v-if="item.icon && item.icon !== ''" class="logo" />
        <div v-else class="logo text-center text-white" :style="{'background': 'orange'}">{{ item.title[0] }}</div>
      </div>
      <div class="ml-1 title truncate ...">
        {{ item.title }}
      </div>
    </div>
    <div v-if="!simple" class="remark ..." style="height: 58px">
      {{item.remark || item.description}}
    </div>
    <div v-if="!simple" class="flex justify-end" style="height: 12px;line-height: 12px;color:#777" @click.stop="void">
      <a-popover :title="item.title">
        <template #content>
          <div style="width:300px">
            <p>{{item.url}}</p>
            <p>{{item.description}}</p>
          </div>
        </template>
        <div >
          ...
        </div>
      </a-popover>
    </div>

  </div>
</template>

<script setup>
import {addUse, apiJson} from "@/api/index";

const props = defineProps({
  item:{},
  simple:{
    type:Boolean,
    default:false
  }
})
function toURL(item){
  apiJson.put({
    tag:"BookmarkUseAdd",
    "BookmarkUse":{
      "bmId":item.bmId
    }
  })
  // addUse({bmId:item.bmId})
  window.open(item.url,"_blank")
}
</script>

<style scoped lang="scss">
.bookmark{
  width:200px;
  height: 100px;
  background: white;
  margin: 4px 4px;



  .logo{
    width: 24px;
    min-width: 24px;
    height: 24px;
    line-height: 24px;
    border-radius: 100%;
    text-transform: uppercase;
    font-weight: bold;
  }

  .title{
    font-size: 14px;
    font-weight: 600;
    color: #333;
  }

  &.simple{
    height: 60px;
    width: 240px;
    line-height: 60px;

    .logo{
      width: 30px;
      min-width: 30px;
      height: 30px;
      line-height: 30px;
      font-size: 16px;
      text-transform: uppercase;
      font-weight: bold;
      margin-left: 4px;
      margin-top: 10px;
      border-radius: 100%;
    }

    .title{
      line-height: 50px;
      margin-left: 8px;
      font-size: 15px;
      font-weight: 600;
      color: #333;

    }
  }

  .remark{
    font-size: 12px;
    color: #acacac;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: break-spaces;
    line-height: 1.2;
  }

  &:hover{
    /*filter: drop-shadow(0 0 4px #42b883aa);*/
    filter: drop-shadow(0 0 4px #646cffaa);
  }
}

</style>
