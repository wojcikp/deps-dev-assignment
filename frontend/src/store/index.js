import { createStore } from 'vuex'
import axios from 'axios'

axios.defaults.baseURL = process.env.VUE_APP_API_URL

export default createStore({
  state: {
    allDependencies: []
  },
  getters: {
    getAllDependencies: state => { return state.allDependencies }
  },
  mutations: {
    setAllDependencies (state, payload) {
      state.allDependencies = payload
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
          commit('setUpdatedDataFromDepsDev', data)
        })
        .catch(err => {
          console.error(err)
        })
    },
    getTestBackend ({ commit, state }, id) {
      return axios.get(`/dependency?id=${id}`)
        .then(response => response.data)
        .then(data => {
          commit('setFromDepsDev', data)
        })
        .catch(err => {
          console.error(err)
        })
    },
    testDeleteBackend ({ commit, state }, id) {
      return axios.delete(`/dependency?id=${id}`)
        .then(response => response.data)
        .then(data => {
          // commit('setTestBackend', data)
        })
        .catch(err => {
          console.error(err)
        })
    },
    testUpdateBackend ({ commit, state }, dependencyDetails) {
      return axios.put('/dependency', dependencyDetails)
        .then(response => response.data)
        .then(data => {
          // commit('setTestBackend', data)
        })
        .catch(err => {
          console.error(err)
        })
    },
    testAddBackend ({ commit, state }, dependencyDetails) {
      return axios.post('/dependency', dependencyDetails)
        .then(response => response.data)
        .then(data => {
          // commit('setTestBackend', data)
        })
        .catch(err => {
          console.error(err)
        })
    },
    testGetByScoreBackend ({ commit, state }, score) {
      return axios.get(`/dependency/score/${score}`)
        .then(response => response.data)
        .then(data => {
          commit('setTestBackend', data)
        })
        .catch(err => {
          console.error(err)
        })
    }
  }
})
