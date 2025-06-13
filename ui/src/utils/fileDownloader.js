/**
 * @module fileDownloader
 * @description Provides a function for downloading files in the browser.
 */

/**
 * Triggers a browser download for a file blob.
 * @param {Blob} blob - The file blob to download.
 * @param {string} filename - The desired name for the downloaded file.
 */
export function downloadFile(blob, filename) {
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename || 'download';
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}