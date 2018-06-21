import Vue from "vue";
import Vuex from "vuex";
Vue.use(Vuex);

const store = new Vuex.Store({
  state: {
    authenticated: false
  },
  mutations: {
    setAuthenticated(state, payload) {
      state.authenticated = payload.isAuthenticated;
    }
  },
  actions: {
    updateAuthenticated({ commit }, isAuthenticated) {
      commit("setAuthenticated", { isAuthenticated });
    }
  },
  getters: {
    authenticated(state) {
      return state.authenticated;
    }
  }
});
export default store;
