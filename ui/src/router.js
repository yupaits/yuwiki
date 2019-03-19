import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'
import SiteSearch from './views/SiteSearch'
import ShareBook from './views/ShareBook'
import SharedBooks from './views/SharedBooks'

Vue.use(Router)

export default new Router({
  routes: [
    {path: '/', component: Home},
    {path: '/search', component: SiteSearch},
    {path: '/book/share', component: ShareBook},
    {path: '/books/shared', component: SharedBooks},
  ]
})
