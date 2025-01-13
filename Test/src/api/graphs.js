import request from '@/util/request'
import router from '../router'
import Vue from 'vue'
import bus from '../util/eventBus'

export const getLineData = (name, system, params, config) => {
  // 拦截器
  let hasWarning = false
  request.interceptors.response.use((res) => {
    if (res.data.code === 201 && !hasWarning) {
      Vue.prototype.$message.warning('登录已过期，请重新登录')
      router.push('/login')
      hasWarning = true
      return
    }

    if (res.data.code === 400 && !hasWarning) {
      Vue.prototype.$message.warning('数据缺失')
      bus.$emit('hideLoading')
      hasWarning = true
      return
    }
    return res
  })
  return request.get(`/data/${name}/${system}?${params}`, config)
}

export const getPieData = (name, system, field, config) => {
  return request.get(`/data/rowall/${name}/${system}/${field}`, config)
}

export const getNameList = (id, config) => {
  return request.get(`/map/${id}`, config)
}
