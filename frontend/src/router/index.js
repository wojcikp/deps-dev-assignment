import { createRouter, createWebHashHistory } from 'vue-router'
// import MainView from '../views/MainView.vue'
import DependenciesView from '../views/DependenciesView.vue'
import ChartView from '../views/ChartView.vue'

const routes = [
  {
    path: '/',
    name: 'dependencies-view',
    component: DependenciesView
  },
  {
    path: '/chart',
    name: 'chart-view',
    component: ChartView
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
