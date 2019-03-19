import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    user: {
      roles: []
    },
    record: {},
  },
  getters: {
    user(state) {
      return state.user;
    },
    record(state) {
      return state.record;
    },
  },
  mutations: {
    setUser(state, user) {
      state.user = user;
    },
    setRecord(state, record) {
      state.record = record;
    },
  },
  actions: {
    setUser({commit}, user) {
      commit('setUser', user);
    },
    setRecord({commit}, record) {
      commit('setRecord', record);
    },
  }
})