<template>
  <v-container class="d-flex justify-center align-center" style="height: 100vh;">
    <v-btn color="primary" @click="handleClick">
      Click Me
    </v-btn>
  </v-container>
</template>

<script setup>
import { useStore } from 'vuex'

const store = useStore()

function handleClick () {
  store.dispatch('getTestBackend')
}
</script>

<style scoped>
.v-container {
  background-color: #f5f5f5;
}
</style>
