<script lang="ts" setup>
import {vDrag,vResize} from './directive'

</script>
<template>
  <div ref="el"   style="position: fixed" class="win" v-resize="true">
      <div v-drag="true" style="background: #2299dd">
          title
      </div>
    <slot>

    </slot>

  </div>
</template>

<style scope>
.win{
  border: 1px solid #2299dd;
  background: #fff;
  border-radius: 4px;
  width: 600px;
  height: 600px;

}
</style>
