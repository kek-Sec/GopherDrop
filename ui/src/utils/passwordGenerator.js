/**
 * @module passwordGenerator
 * @description Provides a function for generating random passwords.
 */

/**
 * Generates a random password.
 * @param {number} length - The length of the password.
 * @param {string} charset - The character set to use for the password.
 * @returns {string} The generated password.
 */
export function generatePassword(length = 12, charset = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()_+') {
  let generatedPassword = '';
  const randomValues = new Uint32Array(length);
  window.crypto.getRandomValues(randomValues);
  for (let i = 0; i < length; i++) {
    const randomIndex = randomValues[i] % charset.length;
    generatedPassword += charset[randomIndex];
  }
  return generatedPassword;
}