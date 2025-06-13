<template>
  <v-container class="d-flex flex-column align-center justify-center fill-height text-center">
    <v-card class="pa-6 pa-md-10" elevation="6" max-width="500" rounded="lg">
      <v-icon size="80" color="error">mdi-close-octagon-outline</v-icon>
      <h1 class="text-h4 mt-4 font-weight-bold">Oops! Something Went Wrong</h1>
      <p class="mt-2 mb-6 text-body-1">
        An unexpected error occurred. Please try again later.
      </p>
      <v-btn to="/" color="primary" large rounded>
        <v-icon left>mdi-home</v-icon>
        Go Back Home
      </v-btn>
    </v-card>
  </v-container>
</template>

<script setup>
/**
 * A general error page for displaying unexpected application issues.
 */
</script>

<style scoped>
.fill-height {
  min-height: 80vh;
}
</style>