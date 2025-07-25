<template>
  <v-container class="d-flex justify-center align-center fill-height">
    <v-card class="pa-4 pa-md-8 animate__animated animate__fadeIn" max-width="600" elevation="6" rounded="lg">
      <v-card-title class="text-h5 text-md-h4 font-weight-bold text-center mb-4">Create a New Secret 🔑</v-card-title>
      <v-card-text>
        <v-form @submit.prevent="handleSubmit">
          <v-btn-toggle
            v-model="type"
            mandatory
            class="mb-4 d-flex justify-center"
            color="primary"
            rounded
            group
          >
            <v-btn value="text" class="px-6" rounded>
              <v-icon left>mdi-text</v-icon> Text
            </v-btn>
            <v-btn value="file" class="px-6" rounded>
              <v-icon left>mdi-file</v-icon> File
            </v-btn>
          </v-btn-toggle>

          <v-textarea
            v-if="type === 'text'"
            label="Text Secret"
            v-model="textSecret"
            required
            variant="outlined"
            rows="4"
          ></v-textarea>

          <v-file-input
            v-if="type === 'file'"
            label="Select File"
            prepend-icon="mdi-upload"
            v-model="files"
            show-size
            required
            variant="outlined"
          ></v-file-input>

          <PasswordInput v-model="password" class="mt-2" />

          <v-select
            label="Expiration"
            v-model="expires"
            :items="expirationOptions"
            required
            variant="outlined"
            class="mt-2"
          ></v-select>

          <v-checkbox
            v-model="oneTime"
            label="One-Time Retrieval"
            color="primary"
            class="mt-2"
          ></v-checkbox>

          <v-btn type="submit" color="primary" class="mt-4" block large rounded x-large height="50">Create Secret</v-btn>

          <v-alert v-if="errorMessage" type="error" class="mt-4 animate__animated animate__bounceIn" variant="tonal">
            {{ errorMessage }}
          </v-alert>
        </v-form>

        <v-alert v-if="resultHash" type="success" class="mt-6 animate__animated animate__fadeIn" variant="tonal">
          <div class="text-h6 mb-2">Secret Created!</div>
          <p>Share this link to view the secret:</p>
          <div class="d-flex align-center mt-2 pa-2" style="background-color: rgba(var(--v-theme-on-surface), 0.05); border-radius: 8px;">
            <span class="mr-2 text-truncate">{{ baseUrl }}/view/{{ resultHash }}</span>
            <v-spacer></v-spacer>
            <v-tooltip text="Copy Link to Clipboard">
              <template v-slot:activator="{ props }">
                <v-btn icon v-bind="props" @click="copyLink">
                  <v-icon>mdi-content-copy</v-icon>
                </v-btn>
              </template>
            </v-tooltip>
          </div>
          <v-snackbar v-model="snackbar" timeout="2000" color="success">
            Link copied to clipboard!
          </v-snackbar>
        </v-alert>

        <v-overlay v-model="loading" class="align-center justify-center">
          <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
        </v-overlay>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref, watch } from 'vue';
import { createSend } from '../services/api.js';
import PasswordInput from '../components/PasswordInput.vue';
import { formStore } from '../stores/formStore.js';

const type = ref('text');
const textSecret = ref('');
// v-file-input uses an array for its model, so initialize it as such.
const files = ref([]);
const password = ref('');
const oneTime = ref(false);
const expires = ref('24h');
const errorMessage = ref('');
const resultHash = ref('');
const baseUrl = window.location.origin;
const snackbar = ref(false);
const loading = ref(false);

const expirationOptions = [
  { title: '1 Hour', value: '1h' },
  { title: '6 Hours', value: '6h' },
  { title: '12 Hours', value: '12h' },
  { title: '24 Hours', value: '24h' },
  { title: '3 Days', value: '72h' },
  { title: '1 Week', value: '168h' }
];

function resetForm() {
  type.value = 'text';
  textSecret.value = '';
  files.value = [];
  password.value = '';
  oneTime.value = false;
  expires.value = '24h';
  errorMessage.value = '';
  resultHash.value = '';
  loading.value = false;
}

watch(() => formStore.resetCounter, () => {
  resetForm();
});

async function handleSubmit() {
  errorMessage.value = '';
  resultHash.value = '';
  loading.value = true;

  const formData = new FormData();
  formData.append('type', type.value);

  if (type.value === 'text') {
    if (!textSecret.value.trim()) {
      errorMessage.value = 'Please enter some text';
      loading.value = false;
      return;
    }
    formData.append('data', textSecret.value);
  } else if (type.value === 'file') {
    // Debug log
    console.log('files.value:', files.value);
    // Ensure files.value is an array and has a File object
    const fileArr = Array.isArray(files.value) ? files.value : (files.value ? [files.value] : []);
    if (!fileArr.length || !(fileArr[0] instanceof File)) {
      errorMessage.value = 'Please select a file 😟';
      loading.value = false;
      return;
    }
    const fileToUpload = fileArr[0];
    formData.append('file', fileToUpload, fileToUpload.name);
  }

  if (password.value.trim()) {
    formData.append('password', password.value);
  }
  if (oneTime.value) {
    formData.append('onetime', 'true');
  }
  formData.append('expires', expires.value);

  try {
    const result = await createSend(formData);
    resultHash.value = result.hash;
    // Clear form inputs but keep the result hash visible
    type.value = 'text';
    textSecret.value = '';
    files.value = [];
    password.value = '';
    oneTime.value = false;
    expires.value = '24h';
  } catch (err) {
    errorMessage.value = err.message || 'Failed to create secret';
  } finally {
    loading.value = false;
  }
}

function copyLink() {
  const link = `${baseUrl}/view/${resultHash.value}`;
  navigator.clipboard.writeText(link);
  snackbar.value = true;
}
</script>

<style scoped>
.v-container {
  min-height: 85vh;
}

.v-card {
  width: 100%;
}
</style>