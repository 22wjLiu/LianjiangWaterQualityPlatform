import axios from 'axios'

const request = axios.create({
  baseURL: 'http://www.shantouliu.site:3000/api'
})

export default request
