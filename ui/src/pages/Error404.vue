<template>
  <v-container class="error-page d-flex flex-column align-center justify-center fill-height">
    <v-card class="pa-6 text-center" elevation="4" max-width="500">
      <v-icon size="80" color="red">mdi-alert-circle-outline</v-icon>
      <h2 class="text-h4 mt-4">Oops! Page Not Found</h2>
      <p class="mt-2 mb-4 text-body-1">
        The page you are looking for does not exist or has been moved.
      </p>
      <v-btn to="/" color="primary" large>
        <v-icon left>mdi-home</v-icon>
        Go Home
      </v-btn>
    </v-card>
  </v-container>
</template>

<script setup>
/**
 * 404 Error Page Component.
 */
</script>

<style scoped>
.error-page {
  min-height: 100vh;
}

.v-card {
  background-color: #f5f5f5;
  border-radius: 12px;
}
</style>
