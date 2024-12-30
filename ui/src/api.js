/**
 * Provides functions to interact with the backend API.
 * Uses fetch for HTTP requests.
 */
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
console.log('API_URL:', API_URL);

export async function createSend(formData) {
  const res = await fetch(`${API_URL}/send`, {
    method: 'POST',
    body: formData
  })
  if (!res.ok) {
    if (res.status === 413) {
      throw new Error('File too large')
    } else if (res.status === 422) {
      throw new Error('Invalid form data')
    } else if (res.status === 429) {
      throw new Error('Too many requests – please try again later')
    } else {
      throw new Error('Failed to create send')
    }
  }
  return res.json()
}

export async function getSend(hash, password = '') {
  const API_BASE_URL = `${window.location.origin}/api`
  const url = new URL(`${API_BASE_URL}/send/${hash}`)

  if (password) url.searchParams.set('password', password)

  const res = await fetch(url)


  if (res.status === 404) {
    console.log('Secret not found.')
    return { notFound: true }
  }

  if (!res.ok) {
    console.error('Failed to retrieve send:', res.statusText)
    throw new Error('Failed to retrieve send')
  }

  const contentDisposition = res.headers.get('Content-Disposition')
  let filename = `download-${hash}`

  if (contentDisposition) {
    const matches = contentDisposition.match(/filename="(.+?)"/)
    if (matches && matches[1]) {
      filename = matches[1]
    }
  }

  const contentType = res.headers.get('content-type')

  if (contentType.includes('application/octet-stream')) {
    const blob = await res.blob()
    return { file: blob, filename }
  }

  const text = await res.text()
  return { text }
}

export async function checkPasswordProtection(hash) {
  const res = await fetch(`${API_URL}/send/${hash}/check`)

  if (res.status === 404) {
    return { notFound: true }
  }

  if (!res.ok) {
    throw new Error('Failed to check password protection')
  }

  return await res.json()
}