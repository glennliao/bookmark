<template>
  <div class="header ">
    <div class="nav flex justify-around">
      <router-link class="nav-item" to="/bookmark">Bookmark</router-link>
      <router-link class="nav-item" to="/note">Note</router-link>
      <div class="nav-item" style="width:40px">
        <a-dropdown>
          <a class="ant-dropdown-link" @click.prevent>
            <a-avatar  :size="32">
              <template #icon><UserOutlined /></template>
            </a-avatar>
          </a>
          <template #overlay>
            <a-menu>
              <a-menu-item @click="logout">
                <a href="javascript:;">logout</a>
              </a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>
      </div>
    </div>
  </div>
</template>

<script>
export default {
}
</script>

<style scoped lang="scss">
.nav{
  width:240px;
  margin: 0 auto;
  border-radius: 0 0 25px 25px ;
  height: 40px;
  position: relative;
  background: rgba(255, 255, 255, .7);
}

.nav-item {
  color:darkseagreen;
  text-decoration: none;
  padding: 4px;
  text-align: center;
  //background: rgba(56, 56, 56, 0.66);
  line-height: 30px;
  cursor: pointer;
}
</style>
