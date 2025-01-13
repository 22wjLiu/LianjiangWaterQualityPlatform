import request from '@/util/request'

export const getMapData = () => {
  return request.get('/map')
}
