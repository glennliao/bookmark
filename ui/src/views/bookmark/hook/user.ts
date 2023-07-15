import {ref} from "vue";

const token = ref(localStorage.getItem("bm_token") || "")

export function useUser() {
  return {
    token
  }
}
