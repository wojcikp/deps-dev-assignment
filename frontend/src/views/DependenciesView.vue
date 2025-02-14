<template>
  <v-container>
    <v-card>
      <v-card-title class="text-center my-4">Dependency Scorecard</v-card-title>

      <v-card-text>
        <v-row>
          <v-col cols="6">
            <v-text-field
              v-model="searchQuery"
              label="Search by Dependency"
              outlined
              dense
              clearable
            ></v-text-field>
          </v-col>

          <v-col cols="6">
            <v-slider
              v-model="minScore"
              label="Minimum Score"
              min="0"
              max="10"
              step="1"
              thumb-label
              class="pt-4"
              color="green"
            ></v-slider>
          </v-col>
        </v-row>
        <v-row justify="center">
          <v-col cols="3"><v-btn variant="outlined" @click="this.updateDependencies()">Update dependencies</v-btn></v-col>
        </v-row>
        <v-row v-if="this.showUpdatedDependencies">
          <v-col class="ml-6">
            Updated dependencies:
            <span v-for="(dependency, i) in this.getUpdatedDependencies" :key="i">{{ dependency }}, </span>
            <br><v-btn class="mt-4" variant="outlined dense" @click="this.dismissUpdateInfo()">Dismiss</v-btn>
          </v-col>
        </v-row>
      </v-card-text>

      <v-card-text>
        <v-table>
          <thead>
            <tr>
              <th class="font-weight-bold">Dependency</th>
              <th class="font-weight-bold">Overall Score</th>
            </tr>
          </thead>
          <tbody>
            <!-- eslint-disable -->
            <template
              v-for="(dependency, index) in filteredDependencies"
              :key="index"
            >
            <!-- eslint-enable -->
              <tr @click="toggleExpand(index)" class="clickable-row">
                <td>
                  <v-icon class="pr-4" :icon="toggleIcon(index)" />
                  {{ dependency.projectKey.id }}
                </td>
                <td>
                  <v-chip
                    class="ml-6"
                    :color="getScoreColor(dependency.scorecard.overallScore)"
                    dark
                  >
                    {{ dependency.scorecard.overallScore }}
                  </v-chip>
                </td>
              </tr>

              <tr v-if="expandedDependencies.includes(index)">
                <td colspan="2">
                  <v-expand-transition>
                    <v-card class="pa-3" outlined>
                      <v-list dense>
                        <v-list-item-group>
                          <v-list-item
                            v-for="(check, i) in dependency.scorecard.checks"
                            :key="i"
                          >
                            <v-list-item-content>
                              <v-list-item-title class="font-weight-bold">
                                {{ check.name }}:
                                <span :class="getScoreColorClass(check.score)">
                                  {{ check.score }}
                                </span>
                              </v-list-item-title>
                              <v-list-item-subtitle>
                                {{ check.reason }}
                              </v-list-item-subtitle>
                            </v-list-item-content>
                          </v-list-item>
                        </v-list-item-group>
                      </v-list>
                    </v-card>
                  </v-expand-transition>
                </td>
              </tr>
            </template>
          </tbody>
        </v-table>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script>
import { mapActions, mapGetters, mapMutations } from 'vuex'

export default {
  name: 'DependenciesList',

  data () {
    return {
      searchQuery: '',
      minScore: 0,
      expandedDependencies: [],
      showUpdatedDependencies: false
    }
  },

  computed: {
    ...mapGetters(['getAllDependencies', 'getUpdatedDependencies']),

    dependencies () {
      return this.getAllDependencies || []
    },

    filteredDependencies () {
      return this.dependencies.filter(dependency => {
        const matchesSearch = !this.searchQuery || dependency.projectKey.id.toLowerCase().includes(this.searchQuery.toLowerCase())
        const matchesScore = dependency.scorecard.overallScore >= this.minScore
        return matchesSearch && matchesScore
      })
    }
  },

  methods: {
    ...mapActions(['getAllDependenciesAction', 'updateDependenciesAction', 'testDeleteBackend']),
    ...mapMutations(['setUpdatedDependencies']),

    getScoreColor (score) {
      if (score >= 8) return 'green'
      if (score >= 5) return 'orange'
      return 'red'
    },

    getScoreColorClass (score) {
      if (score >= 8) return 'text-green'
      if (score >= 5) return 'text-orange'
      return 'text-red'
    },

    toggleExpand (index) {
      if (this.expandedDependencies.includes(index)) {
        this.expandedDependencies = this.expandedDependencies.filter((i) => i !== index)
      } else {
        this.expandedDependencies.push(index)
      }
    },

    toggleIcon (index) {
      if (this.expandedDependencies.includes(index)) {
        return 'mdi-chevron-down'
      } else {
        return 'mdi-chevron-right'
      }
    },

    async updateDependencies () {
      await this.updateDependenciesAction()
      this.getAllDependenciesAction()
      this.showUpdatedDependencies = true
    },

    dismissUpdateInfo () {
      this.showUpdatedDependencies = false
      this.setUpdatedDependencies([])
    }
  }
}
</script>

<style scoped>
.v-chip {
  font-weight: bold;
}
</style>
