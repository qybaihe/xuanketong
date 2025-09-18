<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()

const isRegister = ref(false)
const loading = ref(false)
const formErrors = reactive({
  username: '',
  password: '',
  email: '',
  nickname: '',
  general: ''
})

const formData = reactive({
  username: '',
  password: '',
  email: '',
  nickname: '',
  avatar: ''
})

const validateForm = () => {
  let isValid = true
  
  // Reset errors
  Object.keys(formErrors).forEach(key => {
    formErrors[key as keyof typeof formErrors] = ''
  })
  
  if (!formData.username || formData.username.length < 3) {
    formErrors.username = '用户名至少需要3个字符'
    isValid = false
  }
  
  if (!formData.password || formData.password.length < 6) {
    formErrors.password = '密码至少需要6个字符'
    isValid = false
  }
  
  if (isRegister.value) {
    if (!formData.email || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      formErrors.email = '请输入有效的邮箱地址'
      isValid = false
    }
    
    if (!formData.nickname || formData.nickname.length < 2) {
      formErrors.nickname = '昵称至少需要2个字符'
      isValid = false
    }
  }
  
  return isValid
}

const handleSubmit = async () => {
  if (!validateForm()) return
  
  loading.value = true
  formErrors.general = ''
  
  try {
    let result
    
    if (isRegister.value) {
      result = await authStore.register(formData)
    } else {
      result = await authStore.login({
        username: formData.username,
        password: formData.password
      })
    }
    
    if (result.success) {
      // Redirect to home page on successful auth
      router.push('/')
    } else {
      formErrors.general = result.error || '操作失败，请重试'
    }
  } catch (error) {
    formErrors.general = '网络错误，请检查连接'
  } finally {
    loading.value = false
  }
}

const toggleMode = () => {
  isRegister.value = !isRegister.value
  // Reset form data when switching
  formData.username = ''
  formData.password = ''
  formData.email = ''
  formData.nickname = ''
  // Reset errors
  authStore.clearError()
  Object.keys(formErrors).forEach(key => {
    formErrors[key as keyof typeof formErrors] = ''
  })
}
</script>

<template>
  <div class="auth-container">
    <div class="auth-backdrop"></div>
    
    <div class="auth-card">
      <div class="auth-header">
        <h1 class="auth-title">{{ isRegister ? '创建账户' : '欢迎回来' }}</h1>
        <p class="auth-subtitle">
          {{ isRegister ? '加入我们，发现优质课程' : '登录后继续你的学习之旅' }}
        </p>
      </div>

      <form @submit.prevent="handleSubmit" class="auth-form">
        <div class="form-group">
          <label for="username" class="form-label">用户名</label>
          <input
            id="username"
            v-model="formData.username"
            type="text"
            class="input-glass"
            placeholder="请输入用户名"
            required
          />
          <span v-if="formErrors.username" class="error-text">{{ formErrors.username }}</span>
        </div>

        <div class="form-group">
          <label for="password" class="form-label">密码</label>
          <input
            id="password"
            v-model="formData.password"
            type="password"
            class="input-glass"
            placeholder="请输入密码"
            required
          />
          <span v-if="formErrors.password" class="error-text">{{ formErrors.password }}</span>
        </div>

        <div v-if="isRegister" class="register-fields">
          <div class="form-group">
            <label for="email" class="form-label">电子邮箱</label>
            <input
              id="email"
              v-model="formData.email"
              type="email"
              class="input-glass"
              placeholder="请输入邮箱地址"
              required
            />
            <span v-if="formErrors.email" class="error-text">{{ formErrors.email }}</span>
          </div>

          <div class="form-group">
            <label for="nickname" class="form-label">昵称</label>
            <input
              id="nickname"
              v-model="formData.nickname"
              type="text"
              class="input-glass"
              placeholder="请输入显示昵称"
              required
            />
            <span v-if="formErrors.nickname" class="error-text">{{ formErrors.nickname }}</span>
          </div>
        </div>

        <div v-if="formErrors.general" class="general-error">
          {{ formErrors.general }}
        </div>

        <button
          type="submit"
          class="submit-btn"
          :disabled="loading"
        >
          <span v-if="loading" class="loading-spinner"></span>
          {{ loading ? '处理中...' : (isRegister ? '注册' : '登录') }}
        </button>
      </form>

      <div class="auth-footer">
        <p class="toggle-text">
          {{ isRegister ? '已有账户？' : '还没有账户？' }}
          <button
            type="button"
            class="toggle-btn"
            @click="toggleMode"
          >
            {{ isRegister ? '立即登录' : '立即注册' }}
          </button>
        </p>
      </div>
    </div>

    <!-- Decorative elements -->
    <div class="floating-element element-1"></div>
    <div class="floating-element element-2"></div>
    <div class="floating-element element-3"></div>
  </div>
</template>

<style scoped>
/* ===== 俏皮的新拟物风格 (Playful Neobrutalism) ===== */

/* Auth Container */
.auth-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  background-color: #FEF6F7;
  font-family: sans-serif;
  color: #1A1A1A;
  padding: 20px;
  overflow: hidden;
}

.auth-backdrop {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #F7D074;
  opacity: 0.2;
}

/* Auth Card */
.auth-card {
  position: relative;
  z-index: 1;
  width: 90%;
  max-width: 420px;
  padding: 32px 24px;
  background-color: #FFFFFF;
  border-radius: 12px;
  border: 3px solid #000000;
  box-shadow: 5px 5px 0px 0px #000000;
}

/* Auth Header */
.auth-header {
  text-align: center;
  margin-bottom: 24px;
}

.auth-title {
  font-size: 28px;
  font-weight: bold;
  color: #1A1A1A;
  margin-bottom: 8px;
}

.auth-subtitle {
  font-size: 16px;
  color: #888888;
  margin: 0;
}

/* Form */
.auth-form {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.form-group {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-size: 14px;
  font-weight: bold;
  color: #1A1A1A;
}

.input-glass {
  background-color: #FFFFFF;
  border: 3px solid #000000;
  border-radius: 8px;
  padding: 12px;
  font-size: 16px;
  font-family: sans-serif;
  transition: transform 0.2s ease;
}

.input-glass:focus {
  transform: translate(-2px, -2px);
  box-shadow: 4px 4px 0px 0px #000000;
  outline: none;
}

.error-text {
  display: block;
  font-size: 12px;
  color: #FF3B30;
  margin-top: 4px;
}

.general-error {
  font-size: 14px;
  color: #FF3B30;
  text-align: center;
  padding: 8px;
  background-color: #FFFFFF;
  border-radius: 8px;
  border: 3px solid #FF3B30;
  box-shadow: 4px 4px 0px 0px rgba(255, 59, 48, 0.2);
}

/* Submit Button */
.submit-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 14px 24px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: bold;
  background-color: #F7D074;
  border: 3px solid #000000;
  box-shadow: 4px 4px 0px 0px #000000;
  color: #1A1A1A;
  cursor: pointer;
  transition: transform 0.2s ease;
  min-height: 48px;
}

.submit-btn:hover {
  transform: translate(-2px, -2px);
  box-shadow: 6px 6px 0px 0px #000000;
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
  box-shadow: 4px 4px 0px 0px #000000;
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid transparent;
  border-top: 2px solid #1A1A1A;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* Footer */
.auth-footer {
  text-align: center;
  margin-top: 24px;
}

.toggle-text {
  font-size: 14px;
  color: #888888;
}

.toggle-btn {
  background: none;
  border: none;
  color: #F7D074;
  font-weight: bold;
  cursor: pointer;
  transition: transform 0.2s ease;
  padding: 4px 8px;
  border-radius: 8px;
}

.toggle-btn:hover {
  transform: translate(-2px, -2px);
  box-shadow: 4px 4px 0px 0px #000000;
}

/* Register Fields */
.register-fields {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* Decorative Elements */
.floating-element {
  position: absolute;
  border-radius: 50%;
  background-color: #F7D074;
  opacity: 0.2;
  animation: float 6s ease-in-out infinite;
}

.element-1 {
  width: 80px;
  height: 80px;
  top: 20%;
  left: 10%;
  animation-delay: 0s;
}

.element-2 {
  width: 120px;
  height: 120px;
  bottom: 20%;
  right: 10%;
  animation-delay: -2s;
}

.element-3 {
  width: 60px;
  height: 60px;
  top: 60%;
  left: 5%;
  animation-delay: -4s;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  50% {
    transform: translateY(-20px) rotate(180deg);
  }
}

/* Responsive Design */
@media (min-width: 768px) {
  .auth-container {
    max-width: 768px;
    margin: 0 auto;
  }
}

@media (min-width: 1024px) {
  .auth-container {
    max-width: 1024px;
    margin: 0 auto;
  }
}

@media (max-width: 768px) {
  .auth-card {
    width: 95%;
    padding: 24px 16px;
  }
  
  .auth-title {
    font-size: 24px;
  }
  
  .input-glass {
    font-size: 16px; /* Prevent zoom on iOS */
  }
}

@media (max-width: 480px) {
  .auth-card {
    padding: 20px 16px;
  }
  
  .submit-btn {
    font-size: 16px; /* Prevent zoom on iOS */
  }
}
</style>