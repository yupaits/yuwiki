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
    pageId: undefined,
    part: {},
    keyword: '',
    menuVisible: true,
    sortPart: {
      list: [],
      fromIndex: 0,
      toIndex: 0
    },
    bookStar: false,
    partStar: false,
    pageStar: false
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
    },
    part(state) {
      return state.part;
    },
    keyword(state) {
      return state.keyword;
    },
    menuVisible(state) {
      return state.menuVisible;
    },
    sortPart(state) {
      return state.sortPart;
    },
    bookStar(state) {
      return state.bookStar;
    },
    partStar(state) {
      return state.partStar;
    },
    pageStar(state) {
      return state.pageStar;
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
    setPart(state, part) {
      state.part = part;
    },
    setKeyword(state, keyword) {
      state.keyword = keyword;
    },
    setMenuVisible(state, visible) {
      state.menuVisible = visible;
    },
    setSortPart(state, sortPart) {
      state.sortPart = sortPart;
    },
    setBookStar(state, star) {
      state.bookStar = star;
    },
    setPartStar(state, star) {
      state.partStar = star;
    },
    setPageStar(state, star) {
      state.pageStar = star;
    }
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
    setPart({commit}, part) {
      commit('setPart', part);
    },
    setKeyword({commit}, keyword) {
      commit('setKeyword', keyword);
    },
    setMenuVisible({commit}, visible) {
      commit('setMenuVisible', visible);
    },
    setSortPart({commit}, sortPart) {
      commit('setSortPart', sortPart);
    },
    setBookStar({commit}, star) {
      commit('setBookStar', star);
    },
    setPartStar({commit}, star) {
      commit('setPartStar', star);
    },
    setPageStar({commit}, star) {
      commit('setPageStar', star);
    }
  }
})