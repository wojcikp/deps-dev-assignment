import { createStore } from 'vuex'
import axios from 'axios'

export default createStore({
  state: {
    testBackend: null
  },
  getters: {
    testBackend: state => { return state.testBackend }
  },
  mutations: {
    setTestBackend (state, payload) {
      state.testBackend = payload
    }
  },
  actions: {
    getTestBackend ({ commit, state }) {
      return axios.get('http://localhost:3000/api/hello')
        .then(response => response.data)
        .then(data => {
          console.log(data)
          commit('setTestBackend', data)
        })
        .catch(err => {
          console.error(err)
        })
    }
  }
})
