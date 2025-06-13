<template>
  <v-container class="d-flex flex-column align-center justify-center fill-height text-center">
    <v-card class="pa-6 pa-md-10" elevation="6" max-width="500" rounded="lg">
      <v-icon size="80" color="warning">mdi-alert-circle-outline</v-icon>
      <h1 class="text-h4 mt-4 font-weight-bold">404 - Page Not Found</h1>
      <p class="mt-2 mb-6 text-body-1">
        Sorry, the page you are looking for does not exist or has been moved.
      </p>
      <v-btn to="/" color="primary" large rounded>
        <v-icon left>mdi-home</v-icon>
        Go Home
      </v-btn>
    </v-card>
  </v-container>
</template>

<script setup>
/**
 * A 404 "Page Not Found" error component.
 * It is displayed when a user navigates to a route that does not exist.
 */
</script>

<style scoped>
.fill-height {
  min-height: 80vh;
}
</style>