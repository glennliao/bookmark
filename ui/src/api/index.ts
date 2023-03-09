import http, { Param } from './http'

export const apiJson = {
  get: (data: Param) => http.post('/data/get', data),
  post: (data: Param) => http.post('/data/post', data),
  put: (data: Param) => http.post('/data/put', data),
  delete: (data: Param) => http.post('/data/get', data)
}

export function auth(param:{email:string,password:string}){
  return http.post('/auth', param)
}

export function addUse(param:{bmId:string}){
  return http.post('/addUse', param)
}
