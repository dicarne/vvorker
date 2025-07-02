import axios, { AxiosInstance } from 'axios'
import { setupAxiosEncryption } from '@/lib/encrypt'
import { $user } from '@/store/userState'

const instance: AxiosInstance = axios.create({})

// Setup request interceptor for authentication
instance.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
}, (error) => {
  return Promise.reject(error)
})

// Setup response interceptor for token refresh and error handling
instance.interceptors.response.use((response) => {
  // Update token if new one is provided
  if (response.headers?.['x-authorization-token']) {
    const newToken = response.headers['x-authorization-token']
    localStorage.setItem('token', newToken)
  }

  // Handle error responses
  if (response.data?.code && response.data.code !== 0) {
    throw new Error(response.data.msg || 'Request failed')
  }

  return response
}, (error) => {
  // Handle 401 Unauthorized
  if (error.response?.status === 401) {
    // Clear auth data and redirect to login
    localStorage.removeItem('token')
    window.location.href = '/login'
  }
  return Promise.reject(error)
})

// Setup encryption if user is authenticated
const setupEncryption = () => {
  if (!$user.value || !$user.value.userName) {
    return false
  }
  // Get the encryption key from localStorage or user data
  const user = $user.value
  const encryptionKey = user?.vk

  if (encryptionKey) {
    try {
      setupAxiosEncryption(instance, encryptionKey)
    } catch (error) {
      console.error('Failed to setup encryption:', error)
    }
  }
  return true
}

// Initialize encryption when the module loads
if (typeof window !== 'undefined') {
  // Only run in browser environment
  let int = setInterval(() => {
    if (setupEncryption()) {
      clearInterval(int)
    }
  }, 100)
}

export default instance
