import axios from "axios";
import router from "../router";
import { MessageBox, Message } from 'element-ui'

const request = axios.create({
  baseURL: "http://10.143.50.69:8017",
});

// 请求拦截器
request.interceptors.request.use(
  config => {
    config.signal = null
    // 在请求发送前可以添加 token
    if (config.needToken !== false) {
      const token = localStorage.getItem('token')
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

let isShowingLoginExpired = false

// 响应拦截器
request.interceptors.response.use(
  response => {
    if (response.request?.responseType === 'arraybuffer') {
      return response;
    }

    const res = response.data
    if (res.code !== 200) {

      if (res.code == 201 || res.code == 203) {
        Message({
          message: res.msg || '发生错误',
          type: 'error',
          duration: 5 * 1000
        })
        router.push('/login');
      }

      if (res.code === 202 && !isShowingLoginExpired) {

        isShowingLoginExpired = true;

        MessageBox.confirm('您的登录已经失效，您可以选择取消以停留在本页面或者重新登录', '登录失效', {
          confirmButtonText: '重新登录',
          cancelButtonText: '取消',
          type: 'warning'
        })
        .then(() => {
          localStorage.setItem('token', '')
          router.push('/login');
          isShowingLoginExpired = false; 
        })
        .catch(() => {
          isShowingLoginExpired = false;
        })

        return Promise.reject(new Error(res.msg || '登录失效'));
      }

      return Promise.reject(new Error(res.msg || '发生错误'))
    }else{
      return response.data
    }
  }
)


export default request;
