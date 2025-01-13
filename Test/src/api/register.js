import request from '@/util/request'
export const submit = (data, config) => {
  return request.post('/regist', data, config)
}
export const getVerify = id => {
  return request.get(`/verify/${id}`)
}
