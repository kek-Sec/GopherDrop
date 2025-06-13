/**
 * @module api
 * @description Provides functions to interact with the backend API.
 */

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';
console.log('API_URL:', API_URL);

/**
 * Maps HTTP status codes to error messages for the createSend function.
 * @type {Object<number, string>}
 */
const createSendErrorMessages = {
  413: 'File too large',
  422: 'Invalid form data',
  429: 'Too many requests – please try again later'
};

/**
 * Creates a new "send" with the provided form data.
 * @param {FormData} formData - The data for creating the send.
 * @returns {Promise<object>} The server's JSON response.
 * @throws {Error} If the request fails.
 */
export async function createSend(formData) {
  const res = await fetch(`${API_URL}/send`, {
    method: 'POST',
    body: formData
  });

  if (!res.ok) {
    const errorMessage = createSendErrorMessages[res.status] || 'Failed to create send';
    throw new Error(errorMessage);
  }
  return res.json();
}

/**
 * Retrieves a "send" by its hash.
 * @param {string} hash - The hash of the send to retrieve.
 * @param {string} [password=''] - The password for the send.
 * @returns {Promise<object>} An object containing the send data or a notFound flag.
 * @throws {Error} If the retrieval fails.
 */
export async function getSend(hash, password = '') {
  const url = new URL(`${API_URL}/send/${hash}`);
  if (password) {
    url.searchParams.set('password', password);
  }

  const res = await fetch(url);

  if (res.status === 404) {
    console.log('Secret not found.');
    return { notFound: true };
  }

  if (!res.ok) {
    console.error('Failed to retrieve send:', res.statusText);
    throw new Error('Failed to retrieve send');
  }

  const contentDisposition = res.headers.get('Content-Disposition');
  let filename = `download-${hash}`;
  if (contentDisposition) {
    const matches = contentDisposition.match(/filename="(.+?)"/);
    if (matches && matches[1]) {
      filename = matches[1];
    }
  }

  const contentType = res.headers.get('content-type');
  if (contentType.includes('application/octet-stream')) {
    const blob = await res.blob();
    return { file: blob, filename };
  }

  const text = await res.text();
  return { text };
}

/**
 * Checks if a "send" is password protected.
 * @param {string} hash - The hash of the send to check.
 * @returns {Promise<object>} An object indicating if a password is required or a notFound flag.
 * @throws {Error} If the check fails.
 */
export async function checkPasswordProtection(hash) {
  const res = await fetch(`${API_URL}/send/${hash}/check`);

  if (res.status === 404) {
    return { notFound: true };
  }

  if (!res.ok) {
    throw new Error('Failed to check password protection');
  }

  return res.json();
}