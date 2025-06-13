<template>
  <v-app>
    <v-main class="d-flex flex-column">
      <v-container fluid class="pa-0 flex-grow-1">
        <v-app-bar color="transparent" flat class="px-4">
          <v-toolbar-title>
            <router-link to="/" class="logo-link animate__animated animate__pulse" @click="requestFormReset">
              <img src="./assets/Images/logo.png" alt="Logo" height="40" />
            </router-link>
          </v-toolbar-title>
          <v-spacer></v-spacer>
          <v-btn to="/" text class="animate__animated animate__fadeIn" @click="requestFormReset" rounded>
            <v-icon left>mdi-plus</v-icon> Create New
          </v-btn>
          <v-tooltip :text="isDarkMode ? 'Switch to Light Mode' : 'Switch to Dark Mode'">
            <template v-slot:activator="{ props }">
              <v-btn v-bind="props" icon @click="toggleTheme" class="ml-2">
                <v-icon>{{ isDarkMode ? 'mdi-weather-sunny' : 'mdi-weather-night' }}</v-icon>
              </v-btn>
            </template>
          </v-tooltip>
        </v-app-bar>

        <v-container class="mt-4">
          <router-view />
        </v-container>
      </v-container>

      <v-footer color="transparent" class="justify-center mt-8 pa-4">
        <span>
          © {{ new Date().getFullYear() }} GopherDrop |
          <a
            href="https://github.com/kek-Sec/gopherdrop"
            target="_blank"
            rel="noopener noreferrer"
          >
            GitHub Repository
          </a>
        </span>
      </v-footer>
    </v-main>
  </v-app>
</template>

<script setup>
/**
 * The root component of the application.
 * Provides a header, navigation, and footer.
 */
import { ref, onMounted, computed } from 'vue';
import { formStore } from './stores/formStore.js';
import { useTheme } from 'vuetify';

const themeInstance = useTheme();
const isDarkMode = computed(() => themeInstance.global.current.value.dark);

// On component mount, check localStorage for a saved theme
onMounted(() => {
  const savedTheme = localStorage.getItem('theme');
  if (savedTheme) {
    themeInstance.global.name.value = savedTheme;
  }
});

function toggleTheme() {
  const newTheme = isDarkMode.value ? 'customLightTheme' : 'customDarkTheme';
  themeInstance.global.name.value = newTheme;
  // Save the new theme preference to localStorage
  localStorage.setItem('theme', newTheme);
}

/**
 * Triggers the form reset via the store.
 * This will be picked up by a watcher in Create.vue.
 */
function requestFormReset() {
  formStore.triggerReset();
}
</script>

<style scoped>
.v-footer a {
  font-weight: 500;
}

.logo-link {
  display: inline-block;
  text-decoration: none;
  vertical-align: middle;
}

.v-main {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}
</style>