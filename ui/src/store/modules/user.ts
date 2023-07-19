import { defineStore } from 'pinia'
import { ref } from 'vue'

const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('bm_token') || '')

  const setToken = (t: string) => {
    token.value = t
  }

  const logout = () => {
    token.value = ''
    window.location.reload()
  }

  return {
    token,
    logout,
    setToken
  }
}, {
  persist: {
    key: 'bm_token',
    storage: localStorage
  }
})

export default useUserStore
