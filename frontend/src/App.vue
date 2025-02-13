<template>
  <v-app>
    <v-navigation-drawer app>
      <v-list>
        <v-list-item-title class="text-h6 ml-8 py-6">deps.dev App</v-list-item-title>

        <v-divider></v-divider><div class="pt-6"></div>

        <v-list-item
          v-for="(item, index) in menuItems"
          :key="index"
          :to="item.route"
          link
        >
          <v-list-item-content>
            <v-list-item-title>{{ item.title }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-main>
      <v-container>
        <router-view></router-view>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
import { mapActions } from 'vuex'

export default {
  data () {
    return {
      menuItems: [
        { title: 'Dependencies', route: '/dependencies' },
        { title: 'Chart', route: '/chart' }
      ]
    }
  },

  mounted () {
    this.getAllDependenciesAction()
  },

  methods: {
    ...mapActions(['getAllDependenciesAction'])
  }
}
</script>

<style scoped>
.v-navigation-drawer {
  width: 250px;
}
</style>
