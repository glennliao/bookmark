import { apiJson } from '@/api'
import type { Ref } from 'vue'
import { ref } from 'vue'

export interface AppItem {
  url: string
  title: string
  icon: string
}

export const useAppList = () => {

  const list: Ref<AppItem[]> = ref([])

  const loadList = () => {
    apiJson.get({
      'AppList[]': {
        '@alias': 'list'
      }
    }).then(data => {
      list.value = data.list.sort((a,b)=>a.sorts-b.sorts)
    })
  }

  const del = (appId: string) => {
    apiJson.delete({
      AppList: {
        appId
      },
      tag: 'AppList'
    }).then(data => {
      loadList()
    })
  }

  return {
    loadList,
    del,
    list
  }
}
