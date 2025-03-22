import axios from 'axios'

const request = axios.create({
  baseURL: 'http://192.168.10.100:1017'
})

export default request
