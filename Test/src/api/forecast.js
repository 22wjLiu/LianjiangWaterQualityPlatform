import request from '@/util/request'
export const forecast = (params, config) => {
  return request.get(`/forecast?${params}`, config)
}
