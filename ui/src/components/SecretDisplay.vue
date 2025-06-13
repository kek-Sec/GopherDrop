<template>
  <div>
    <v-alert v-if="textContent" type="info" class="mb-4 text-left" variant="tonal">
      {{ textContent }}
    </v-alert>

    <div v-if="file">
      <v-btn color="primary" @click="download" block x-large height="50" rounded>
        <v-icon left>mdi-download</v-icon> Download File
      </v-btn>
    </div>

    <v-btn v-if="textContent" color="primary" @click="copy" block x-large height="50" rounded class="mt-4">
      <v-icon left>mdi-content-copy</v-icon> Copy Secret
    </v-btn>

    <v-snackbar v-model="snackbar" timeout="2000">
      Text copied to clipboard!
    </v-snackbar>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { downloadFile } from '../utils/fileDownloader.js';

const props = defineProps({
  textContent: String,
  file: Blob,
  filename: String
});

const snackbar = ref(false);

function download() {
  if (props.file) {
    downloadFile(props.file, props.filename);
  }
}

function copy() {
  if (props.textContent) {
    navigator.clipboard.writeText(props.textContent);
    snackbar.value = true;
  }
}
</script>