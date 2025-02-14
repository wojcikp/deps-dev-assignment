import { createStore } from 'vuex'
import axios from 'axios'

axios.defaults.baseURL = process.env.VUE_APP_API_URL

export default createStore({
  state: {
    allDependencies: [],
    updatedDependencies: []
  },
  getters: {
    getAllDependencies: state => { return state.allDependencies },
    getUpdatedDependencies: state => { return state.updatedDependencies }
  },
  mutations: {
    setAllDependencies (state, payload) {
      state.allDependencies = payload
    },
    setUpdatedDependencies (state, payload) {
      state.updatedDependencies = payload
    }
  },
  actions: {
    getAllDependenciesAction ({ commit, state }) {
      return axios.get('/dependency/all')
        .then(response => response.data)
        .then(data => {
          commit('setAllDependencies', data)
        })
        .catch(err => {
          console.error(err)
        })
    },
    updateDependenciesAction ({ commit, state }) {
      return axios.get('/dependency/update')
        .then(response => response.data)
        .then(data => {
          commit('setUpdatedDependencies', data)
        })
        .catch(err => {
          console.error(err)
        })
    }
  }
})
