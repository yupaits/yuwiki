import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home'
import PageTemplates from './views/PageTemplates'
import Profile from './views/Profile'
import PageEditor from './views/PageEditor'
import SiteSearch from './views/SiteSearch'
import ShareBook from './views/ShareBook'
import SharedBooks from './views/SharedBooks'

import store from './store'
import api from '../../api'

Vue.use(Router)

const router = new Router({
  routes: [
    {path: '/', component: Home},
    {path: '/templates', component: PageTemplates},
    {path: '/profile', component: Profile},
    {path: '/page/edit', component: PageEditor},
    {path: '/search', component: SiteSearch},
    {path: '/book/share', component: ShareBook},
    {path: '/books/shared', component: SharedBooks},
  ]
})

router.beforeEach((to, from, next) => {
  if (!store.getters.user.id) {
    api.getUserInfo().then(res => {
      store.dispatch('setUser', res.data)
    })
  }
  next()
})

export default router