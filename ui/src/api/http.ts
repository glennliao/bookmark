import router from '../router'
import { useUser } from '@/views/bookmark/hook/user'
import { message } from 'ant-design-vue'
const user = useUser()
export function request (method: string, url: string, query: Record<string, any> | undefined, data: Record<string, any> | undefined, headers: Record<string, any> = {}): Promise<Record<string, any>> {
  if (query) {
    const paramsArray: string[] = []
    Object.keys(query).forEach(key => paramsArray.push(key + '=' + query[key]))
    if (url.search(/\?/) === -1) {
      url += '?' + paramsArray.join('&')
    } else {
      url += '&' + paramsArray.join('&')
    }
  }

  headers.Authorization = user.token.value
  url = './api' + url

  return fetch(url, {
    method,
    headers,
    body: data ? JSON.stringify(data) : data
  }).then(res => res.json()).then(data => {
    if (data.code === 200) {
      return data.data
    }
    if (data.code === 401 && router.currentRoute.value.path !== '/login') {
      router.push('/login')
    }
    message.warning(data.msg)
    throw data
  })
}

export type Param = Record<string, any>

export default {
  get: (url: string, query: Param) => request('get', url, query, undefined),
  post: (url: string, data: Param) => request('post', url, undefined, data),
  put: (url: string, data: Param) => request('put', url, undefined, data),
  delete: (url: string, query: Param) => request('delete', url, query, undefined)
}
