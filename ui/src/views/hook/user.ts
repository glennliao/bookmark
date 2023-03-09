import {ref} from "vue";

const token = ref(localStorage.getItem("token") || "")

export function useUser() {
  return {
    token
  }
}
