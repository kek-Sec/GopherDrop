<template>
  <v-container class="d-flex justify-center align-center fill-height">
    <v-card class="pa-4" max-width="600" outlined>
      <v-card-title class="text-h5 text-center">Create a New Secret</v-card-title>
      <v-card-text>
        <v-form @submit.prevent="handleSubmit">
          <v-select
            label="Type"
            v-model="type"
            :items="['text', 'file']"
            required
          ></v-select>

          <v-textarea
            v-if="type === 'text'"
            label="Text Secret"
            v-model="textSecret"
            required
          ></v-textarea>

          <v-file-input
            v-if="type === 'file'"
            label="Select File"
            prepend-icon="mdi-upload"
            @update:modelValue="handleFile"
            show-size
            required
          ></v-file-input>

          <v-text-field
            label="Password (optional)"
            v-model="password"
            type="password"
          ></v-text-field>

          <v-select
            label="Expiration"
            v-model="expires"
            :items="expirationOptions"
            required
          ></v-select>

          <v-checkbox
            v-model="oneTime"
            label="One-Time Retrieval"
            class="mt-4"
          ></v-checkbox>

          <v-btn type="submit" color="primary" class="mt-4" block>Create</v-btn>

          <v-alert v-if="errorMessage" type="error" class="mt-4">
            {{ errorMessage }}
          </v-alert>
        </v-form>

        <v-alert v-if="resultHash" type="success" class="mt-4">
          Secret Created! Share this link:<br />
          <div class="d-flex align-center mt-2">
            <v-chip class="mr-2">{{ baseUrl }}/view/{{ resultHash }}</v-chip>
            <v-btn icon @click="copyLink" color="light-green">
              <v-icon>mdi-content-copy</v-icon>
            </v-btn>
          </div>
          <v-snackbar v-model="snackbar" timeout="2000">
            Link copied to clipboard!
          </v-snackbar>
        </v-alert>

        <!-- Loader Overlay -->
        <v-overlay v-model="loading" class="align-center justify-center">
          <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
        </v-overlay>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref } from 'vue'
import { createSend } from '../api.js'

const type = ref('text')
const textSecret = ref('')
const fileBlob = ref(null)
const password = ref('')
const oneTime = ref(false)
const expires = ref('24h')
const errorMessage = ref('')
const resultHash = ref('')
const baseUrl = window.location.origin
const snackbar = ref(false)
const loading = ref(false)

const expirationOptions = [
  { title: '1 Hour', value: '1h' },
  { title: '6 Hours', value: '6h' },
  { title: '12 Hours', value: '12h' },
  { title: '24 Hours', value: '24h' },
  { title: '3 Days', value: '72h' },
  { title: '1 Week', value: '168h' }
]

function handleFile(file) {
  if (file) {
    fileBlob.value = file
    console.log("Selected file:", fileBlob.value)
  } else {
    fileBlob.value = null
    console.error("No file selected")
  }
}

async function handleSubmit() {
  errorMessage.value = ''
  resultHash.value = ''
  loading.value = true

  const formData = new FormData()
  formData.append('type', type.value)

  if (type.value === 'text') {
    if (!textSecret.value.trim()) {
      errorMessage.value = 'Please enter some text'
      loading.value = false
      return
    }
    formData.append('data', textSecret.value)
  } else if (type.value === 'file') {
    if (!fileBlob.value) {
      errorMessage.value = 'Please select a file'
      loading.value = false
      return
    }
    formData.append('file', fileBlob.value)
  }

  if (password.value.trim()) formData.append('password', password.value)
  if (oneTime.value) formData.append('onetime', 'true')
  formData.append('expires', expires.value)

  try {
    const result = await createSend(formData)
    if (result.hash) {
      resultHash.value = result.hash
    } else {
      errorMessage.value = 'Unexpected error: no hash returned'
    }
  } catch (err) {
    console.error("Submit error:", err)
    if (err instanceof Error) {
      errorMessage.value = err.message
    } else {
      errorMessage.value = 'Failed to create secret'
    }
  } finally {
    loading.value = false
  }
}

function copyLink() {
  const link = `${baseUrl}/view/${resultHash.value}`
  navigator.clipboard.writeText(link)
  snackbar.value = true
}
</script>

<style scoped>
.v-container {
  min-height: 100vh;
}

.v-card {
  width: 100%;
}
</style>
