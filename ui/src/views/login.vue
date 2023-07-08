<template>
  <div>
    <div class="login">
      <div>
        <a-form
          :form="form"
          layout="vertical"
        >
          <a-form-item label="email">
            <a-input v-model:value="form.email"/>
          </a-form-item>
          <a-form-item label="password">
            <a-input-password @pressEnter="handleOk" v-model:value="form.password"/>
          </a-form-item>
          <a-form-item>
            <div class="mx-auto">
              <a-button type="primary" block @click="handleOk">登录</a-button>
            </div>
          </a-form-item>
        </a-form>

        <div v-if="layoutData.canReg" class="cursor-pointer text-center" style="color: #aaa;font-size: 14px;"
             @click="openRegister">
          注册账号

          <div v-if="layoutData.msg" style="font-size: 12px; margin-top: 12px;color: #b6b6b6">
            {{ layoutData.msg }}
          </div>
        </div>
      </div>
    </div>


    <a-modal
      title="register"
      :open="registerModalShow"
      @ok="registerSubmit"
      @cancel="cancel"
    >

      <a-form
        :form="form"
        layout="vertical"
      >
        <a-form-item label="email">
          <a-input v-model:value="form.email"/>
        </a-form-item>
        <a-form-item label="password">
          <a-input-password @pressEnter="registerSubmit" v-model:value="form.password"/>
        </a-form-item>
        <a-form-item label="code">
          <a-input-password @pressEnter="registerSubmit" v-model:value="form.code"/>
        </a-form-item>
      </a-form>

    </a-modal>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { auth, register, registerLayout } from '../api'
import { useRouter } from 'vue-router'
import { useUser } from '@/views/hook/user'
import { message } from 'ant-design-vue'

const form = ref({})

const router = useRouter()
const user = useUser()

function handleOk () {
  auth(form.value).then(data => {
    user.token.value = data.token
    localStorage.setItem('token', data.token)
    router.push('/')
  })
}

const registerModalShow = ref(false)

function openRegister () {
  registerModalShow.value = true
}

function registerSubmit () {
  register(form.value).then(data => {
    message.success('success')
    loadLayout()
    form.value = {}
    registerModalShow.value = false
  })
}

function cancel () {
  form.value = {}
  registerModalShow.value = false
}

const layoutData = reactive(({
  canReg: false,
  msg: ''
}))

function loadLayout() {
  registerLayout().then(data => {
    layoutData.msg = data.msg
    layoutData.canReg = data.canReg
  })
}

loadLayout()

</script>

<style scoped>
.login {
  padding: 100px 0;
  width: 320px;
  margin: 0 auto;
}
</style>
