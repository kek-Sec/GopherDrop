import { reactive } from 'vue';

/**
 * A simple reactive store to handle cross-component communication
 * for resetting the create form.
 */
export const formStore = reactive({
  // A counter that components can watch. Incrementing it signals a reset.
  resetCounter: 0,

  // Function to trigger the reset from anywhere in the app.
  triggerReset() {
    this.resetCounter++;
  }
});