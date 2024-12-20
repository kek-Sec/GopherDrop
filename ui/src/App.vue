<template>
  <v-app>
    <v-main>
      <v-container fluid class="pa-0">
        <!-- Header -->
        <v-toolbar color="primary" dark flat>
          <v-toolbar-title>
            <router-link to="/" class="logo-link">
              <img src="./assets/Images/logo.png" alt="Logo" height="50" />
            </router-link>
          </v-toolbar-title>
          <v-spacer></v-spacer>
          <v-btn to="/" text>
            <v-icon left>mdi-plus</v-icon> Create
          </v-btn>
        </v-toolbar>

        <!-- Main Content -->
        <v-container class="mt-4">
          <router-view />
        </v-container>

        <!-- Footer -->
        <v-footer color="primary" dark class="justify-center mt-4">
          <span>
            © {{ new Date().getFullYear() }} GopherDrop | 
            <a
              href="https://github.com/kek-Sec/gopherdrop"
              target="_blank"
              rel="noopener noreferrer"
              class="text-white"
            >
              GitHub Repository
            </a>
          </span>
        </v-footer>
      </v-container>
    </v-main>
  </v-app>
</template>

<script setup>
/**
 * The root component of the application.
 * Provides a header, navigation, and footer.
 */
</script>

<style scoped>
.v-toolbar {
  width: 100%;
  margin: 0;
}

.v-footer {
  width: 100%;
  padding: 16px;
  font-size: 14px;
  text-align: center;
}

.v-footer a {
  text-decoration: none;
  color: inherit;
}

.v-footer a:hover {
  text-decoration: underline;
}

.logo-link {
  display: inline-block;
  text-decoration: none;
}

</style>
