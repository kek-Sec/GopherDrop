<template>
  <v-container class="d-flex justify-center align-center fill-height">
    <v-card class="pa-4" max-width="600" outlined>
      <v-card-title class="text-h5 text-center">View Secret</v-card-title>
      <v-card-text>
        <v-alert v-if="errorMessage" type="error" class="mb-4">
          {{ errorMessage }}
        </v-alert>

        <v-alert v-if="notFound" type="error" class="mb-4">
          Secret not found or has expired.
        </v-alert>

        <v-form @submit.prevent="loadSecret" v-if="!secretLoaded">
          <v-text-field
            v-if="requiresPassword"
            label="Password"
            v-model="password"
            type="password"
            prepend-icon="mdi-lock"
            required
          ></v-text-field>

          <v-btn type="submit" color="primary" block>Load Secret</v-btn>
        </v-form>

        <div v-if="secretLoaded">
          <!-- Display text secret -->
          <v-alert v-if="secretContent" type="info" class="mb-4">
            {{ secretContent }}
          </v-alert>

          <!-- Download button for file secret -->
          <div v-if="fileBlob">
            <v-btn color="primary" @click="downloadFile" block>
              <v-icon left>mdi-download</v-icon> Download File
            </v-btn>
          </div>

          <!-- Copy button for text secret -->
          <v-btn v-if="secretContent" color="primary" @click="copyContent" block class="mt-4">
            <v-icon left>mdi-content-copy</v-icon> Copy Secret
          </v-btn>

          <v-snackbar v-model="snackbar" timeout="2000">
            Link copied to clipboard!
          </v-snackbar>
        </div>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getSend, checkPasswordProtection } from '../api.js'

const route = useRoute()
const hash = route.params.hash
const password = ref('')
const errorMessage = ref('')
const notFound = ref(false)
const requiresPassword = ref(false)
const secretContent = ref('')
const fileBlob = ref(null)
const filename = ref('')
const secretLoaded = ref(false)
const snackbar = ref(false)

async function loadSecret() {
  errorMessage.value = ''
  notFound.value = false
  secretContent.value = ''
  fileBlob.value = null

  try {
    const result = await getSend(hash, password.value)
    if (result.notFound) {
      notFound.value = true
      return
    }

    if (result.file) {
      fileBlob.value = result.file
      filename.value = result.filename
    } else {
      secretContent.value = result.text
    }

    secretLoaded.value = true
  } catch (err) {
    errorMessage.value = 'Failed to load secret. Incorrect password or secret has expired.'
  }
}

function copyContent() {
  navigator.clipboard.writeText(secretContent.value)
  snackbar.value = true
}

function downloadFile() {
  const url = URL.createObjectURL(fileBlob.value)
  const a = document.createElement('a')
  a.href = url
  a.download = filename.value || 'download'
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

onMounted(async () => {
  try {
    const result = await checkPasswordProtection(hash)
    requiresPassword.value = result.requiresPassword || false
  } catch {
    notFound.value = true
  }
})
</script>

<style scoped>
.v-container {
  min-height: 100vh;
}

.v-card {
  width: 100%;
}
</style>
