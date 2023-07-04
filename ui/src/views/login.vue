<template>
  <div>
    <div class="login">
      <div>
        <a-form
          :form="form"
          layout="vertical"
        >
          <a-form-item label="email">
            <a-input v-model:value="form.email" />
          </a-form-item>
          <a-form-item label="password">
            <a-input-password @pressEnter="handleOk" v-model:value="form.password" />
          </a-form-item>
          <a-form-item>
            <div class="mx-auto">
              <a-button type="primary" block @click="handleOk">登录</a-button>
            </div>
          </a-form-item>
        </a-form>
      </div>
    </div>
  </div>
</template>

<script setup>
import {ref} from "vue";
import {auth} from '../api'
import {useRouter} from "vue-router";
import {useUser} from "@/views/hook/user";

const form = ref({})

const router = useRouter()
const user = useUser()

function handleOk(){
  auth(form.value).then(data=>{
    user.token.value = data.token
    localStorage.setItem('token',data.token)
    router.push("/")
  })
}
</script>

<style scoped>
.login{
  padding: 100px 0;
  width: 320px;
  margin: 0 auto;
}
</style>
