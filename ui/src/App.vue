<template>
  <v-app>
    <v-main>
      <v-container fluid class="pa-0">
        <v-toolbar color="primary" dark flat>
          <v-toolbar-title>
            <img src="./assets/Images/logo.png" alt="Logo" height="50" />
          </v-toolbar-title>
          <v-spacer></v-spacer>
          <v-btn to="/" text>
            <v-icon left>mdi-plus</v-icon> Create </v-btn>
        </v-toolbar>
        <v-container class="mt-4">
          <router-view />
        </v-container>
      </v-container>
    </v-main>
  </v-app>
</template>

<script setup>
/**
 * The root component of the application.
 * Provides a header and navigation.
 */
</script>

<style scoped>
.v-toolbar {
  width: 100%;
  margin: 0;
}
</style>
