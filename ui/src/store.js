import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    user: {
      roles: []
    },
    record: {},
    bookId: undefined,
    partId: undefined,
    pageId: undefined
  },
  getters: {
    user(state) {
      return state.user;
    },
    record(state) {
      return state.record;
    },
    bookId(state) {
      return state.bookId;
    },
    partId(state) {
      return state.partId;
    },
    pageId(state) {
      return state.pageId;
    }
  },
  mutations: {
    setUser(state, user) {
      state.user = user;
    },
    setRecord(state, record) {
      state.record = record;
    },
    setBookId(state, bookId) {
      state.bookId = bookId;
    },
    setPartId(state, partId) {
      state.partId = partId;
    },
    setPageId(state, pageId) {
      state.pageId = pageId;
    },
  },
  actions: {
    setUser({commit}, user) {
      commit('setUser', user);
    },
    setRecord({commit}, record) {
      commit('setRecord', record);
    },
    setBookId({commit}, bookId) {
      commit('setBookId', bookId);
    },
    setPartId({commit}, partId) {
      commit('setPartId', partId);
    },
    setPageId({commit}, pageId) {
      commit('setPageId', pageId);
    },
  }
})