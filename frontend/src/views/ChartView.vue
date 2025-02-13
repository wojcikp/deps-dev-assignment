<template>
  <Bar :data="chartData"/>
</template>

<script>
import { mapGetters } from 'vuex'

import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale } from 'chart.js'
import { Bar } from 'vue-chartjs'

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend)

export default {
  name: 'App',
  components: {
    Bar
  },
  computed: {
    ...mapGetters(['getAllDependencies']),
    chartData () {
      return {
        labels: this.getAllDependencies.map(item => item.projectKey?.id || 'Unknown'),
        datasets: [
          {
            label: 'Dependencies Score Chart',
            backgroundColor: '#4CAF50',
            data: this.getAllDependencies.map(item => item.scorecard?.overallScore || 0)
          }
        ]
      }
    }
  }
}
</script>
