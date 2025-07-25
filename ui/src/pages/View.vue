<template>
  <v-container class="d-flex justify-center align-center fill-height">
    <v-card class="pa-4 pa-md-8 animate__animated animate__fadeIn" max-width="600" elevation="6" rounded="lg">
      <v-card-title class="text-h5 text-md-h4 font-weight-bold text-center mb-4">View Secret 🔒</v-card-title>
      <v-card-text>
        <v-alert v-if="errorMessage" type="error" class="mb-4" variant="tonal">
          {{ errorMessage }}
        </v-alert>

        <v-alert v-if="notFound" type="error" class="mb-4" variant="tonal">
          Secret not found or has expired.
        </v-alert>

        <v-form @submit.prevent="loadSecret" v-if="!secretLoaded">
          <v-text-field
            v-if="requiresPassword"
            label="Password"
            v-model="password"
            type="password"
            prepend-inner-icon="mdi-lock"
            required
            variant="outlined"
            class="animate__animated animate__fadeIn"
          ></v-text-field>

          <v-btn type="submit" color="primary" class="animate__animated animate__fadeIn mt-4" block x-large height="50" rounded>Load Secret</v-btn>
        </v-form>

        <SecretDisplay
          v-if="secretLoaded"
          :text-content="secretContent"
          :file="fileBlob"
          :filename="filename"
        />
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { getSend, checkPasswordProtection } from '../services/api.js';
import SecretDisplay from '../components/SecretDisplay.vue';

const route = useRoute();
const hash = route.params.hash;

const password = ref('');
const errorMessage = ref('');
const notFound = ref(false);
const requiresPassword = ref(false);
const secretContent = ref('');
const fileBlob = ref(null);
const filename = ref('');
const secretLoaded = ref(false);

async function loadSecret() {
  errorMessage.value = '';
  try {
    const result = await getSend(hash, password.value);
    if (result.notFound) {
      notFound.value = true;
      return;
    }

    if (result.file) {
      fileBlob.value = result.file;
      filename.value = result.filename;
    } else {
      secretContent.value = result.text;
    }
    secretLoaded.value = true;
  } catch (err) {
    errorMessage.value = err.message || 'Failed to load secret. Incorrect password or secret has expired.';
  }
}

onMounted(async () => {
  try {
    const result = await checkPasswordProtection(hash);
    if (result.notFound) {
      notFound.value = true;
    } else {
      requiresPassword.value = result.requiresPassword;
      if (!requiresPassword.value) {
        loadSecret();
      }
    }
  } catch (err) {
    notFound.value = true;
  }
});
</script>

<style scoped>
.v-container {
  min-height: 85vh;
}

.v-card {
  width: 100%;
}
</style>