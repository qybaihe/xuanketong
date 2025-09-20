import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'
import api from '@/services/api'

export interface User {
  id: number
  username: string
  email: string
  nickname: string
  avatar: string
  role: string
  createdAt: string
  updatedAt: string
}

export interface AuthResponse {
  token: string
  user: User
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(localStorage.getItem('token'))
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Computed
  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  // Set token
  const setToken = (newToken: string | null) => {
    token.value = newToken
    if (newToken) {
      localStorage.setItem('token', newToken)
      axios.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
      api.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
    } else {
      localStorage.removeItem('token')
      delete axios.defaults.headers.common['Authorization']
      delete api.defaults.headers.common['Authorization']
    }
  }

  // Set auth header
  if (token.value) {
    axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
    api.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
  }

  // Register
  const register = async (userData: {
    username: string
    password: string
    email: string
    nickname: string
  }) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await axios.post<AuthResponse>(`${import.meta.env.VITE_BACKEND_BASE_URL}/auth/register`, userData)
      
      const { token: newToken, user: userInfo } = response.data
      setToken(newToken)
      user.value = userInfo
      
      return { success: true, data: response.data }
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || 'Registration failed'
      error.value = errorMessage
      return { success: false, error: errorMessage }
    } finally {
      loading.value = false
    }
  }

  // Login
  const login = async (credentials: {
    username: string
    password: string
  }) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await axios.post<AuthResponse>(`${import.meta.env.VITE_BACKEND_BASE_URL}/auth/login`, credentials)
      
      const { token: newToken, user: userInfo } = response.data
      setToken(newToken)
      user.value = userInfo
      
      return { success: true, data: response.data }
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || 'Login failed'
      error.value = errorMessage
      return { success: false, error: errorMessage }
    } finally {
      loading.value = false
    }
  }

  // Logout
  const logout = () => {
    setToken(null)
    user.value = null
    error.value = null
  }

  // Get current user
  const getCurrentUser = async () => {
    if (!token.value) return
    
    loading.value = true
    error.value = null
    
    try {
      const response = await axios.get<User>(`${import.meta.env.VITE_BACKEND_BASE_URL}/auth/me`)
      user.value = response.data
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || 'Failed to get user info'
      error.value = errorMessage
      // If token is invalid, logout
      if (err.response?.status === 401) {
        logout()
      }
    } finally {
      loading.value = false
    }
  }

  // OAuth2 Login
  const oauth2Login = async (token: string, userInfo: User) => {
    setToken(token)
    user.value = userInfo
    error.value = null
  }

  // Clear error
  const clearError = () => {
    error.value = null
  }

  return {
    user,
    token,
    loading,
    error,
    isAuthenticated,
    isAdmin,
    register,
    login,
    logout,
    getCurrentUser,
    oauth2Login,
    clearError
  }
})