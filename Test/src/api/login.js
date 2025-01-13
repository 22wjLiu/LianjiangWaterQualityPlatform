import request from '@/util/request'

export const login = data => {
  return request.post('/login', data)
}
export const verifyEmail = id => {
  return request.get(`/verify/${id}`)
}
export const security = data => {
  return request.put('/security', data)
}
export const updatepass = (data, config) => {
  return request.put('/updatepass', data, config)
}
