import Vue from 'vue'
import App from './App.vue'
import router from './router'
import Antd from 'ant-design-vue'
import MavonEditor from 'mavon-editor'

import 'ant-design-vue/dist/antd.min.css'
import 'mavon-editor/dist/css/index.css'
import './styles/index.css'

Vue.config.productionTip = false

Vue.use(Antd)
Vue.use(MavonEditor)

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
