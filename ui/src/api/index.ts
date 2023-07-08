import http, { Param } from './http'

export const apiJson = {
  get: (data: Param) => http.post('/data/get', data),
  post: (data: Param) => http.post('/data/post', data),
  put: (data: Param) => http.post('/data/put', data),
  delete: (data: Param) => http.post('/data/delete', data)
}

export function auth (param: { email: string, password: string }) {
  return http.post('/auth', param)
}

export function register (param: { email: string, password: string, code: string }) {
  return http.post('/register', param)
}

export function registerLayout (param: { email: string, password: string, code: string }) {
  return http.get('/register', param)
}

export function upload (param: any) {
  return http.post('/upload', param)
}
