import Vue from 'vue'
import App from './App.vue'
import router from './router'
import ElementUI from 'element-ui'
import BaiduMap from 'vue-baidu-map'
import * as echarts from 'echarts'
import moment from 'moment'
import 'github-markdown-css/github-markdown.css'

Vue.config.productionTip = false
Vue.use(BaiduMap, {
  ak: ''
})
Vue.use(ElementUI)
Vue.prototype.$echarts = echarts
Vue.prototype.$moment = moment

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
