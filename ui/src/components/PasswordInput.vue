<template>
  <v-text-field
    label="Password (optional)"
    :type="showPassword ? 'text' : 'password'"
    :model-value="modelValue"
    @update:modelValue="$emit('update:modelValue', $event)"
    variant="outlined"
    prepend-inner-icon="mdi-lock"
  >
    <template v-slot:append-inner>
      <v-tooltip text="Toggle Password Visibility">
        <template v-slot:activator="{ props }">
          <v-btn icon v-bind="props" @click="showPassword = !showPassword" size="small">
            <v-icon>{{ showPassword ? 'mdi-eye-off' : 'mdi-eye' }}</v-icon>
          </v-btn>
        </template>
      </v-tooltip>
      <v-tooltip text="Generate Random Password">
        <template v-slot:activator="{ props }">
          <v-btn icon color="primary" v-bind="props" @click="generateNewPassword" size="small" style="margin-left: 4px">
            <v-icon>mdi-refresh</v-icon>
          </v-btn>
        </template>
      </v-tooltip>
    </template>
  </v-text-field>
</template>

<script setup>
import { ref } from 'vue';
import { generatePassword } from '../utils/passwordGenerator.js';

const props = defineProps({
  modelValue: String
});
const emit = defineEmits(['update:modelValue']);

const showPassword = ref(false);

function generateNewPassword() {
  emit('update:modelValue', generatePassword());
}
</script>