import Vue from 'vue'
import App from './App.vue'
import api from '../../api'
import Antd from 'ant-design-vue'

import 'ant-design-vue/dist/antd.min.css'

Vue.config.productionTip = false

Vue.use(Antd)

Vue.prototype.$api = api

new Vue({
  render: h => h(App)
}).$mount('#app')
