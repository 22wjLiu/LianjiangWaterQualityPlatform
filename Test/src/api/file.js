import request from '@/util/request'

export const download = (path, file, config) => {
  return request.get(`/download?${path}${file}`, config)
}
export const upload = (system, data, config) => {
  return request.post(`/upload/${system}`, data, config)
}
export const getFileList = (path, config) => {
  return request.get(`/files?${path}`, config)
}
