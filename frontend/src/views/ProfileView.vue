<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore, type DarkModePreference } from '@/stores/theme'

const authStore = useAuthStore()
const themeStore = useThemeStore()

const loading = ref(false)

// 主题设置
const availableThemes = computed(() => themeStore.themes)
const currentTheme = computed(() => themeStore.currentTheme)
const darkModePreference = computed(() => themeStore.darkModePreference)

// 夜间模式下拉选择状态
const showDarkModeDropdown = ref(false)

const setTheme = (themeName: string) => {
  themeStore.setTheme(themeName)
}

const setDarkMode = (preference: DarkModePreference) => {
  themeStore.setDarkMode(preference)
}

// 夜间模式下拉选择方法
const toggleDarkModeDropdown = () => {
  showDarkModeDropdown.value = !showDarkModeDropdown.value
}

const selectDarkModeOption = (preference: DarkModePreference) => {
  setDarkMode(preference)
  showDarkModeDropdown.value = false
}

// 获取夜间模式显示文本
const getDarkModeDisplayText = (preference: DarkModePreference) => {
  switch (preference) {
    case 'system':
      return '跟随系统'
    case 'enabled':
      return '一直启用'
    case 'disabled':
      return '一直禁用'
    default:
      return '跟随系统'
  }
}

onMounted(async () => {
  if (!authStore.user) {
    loading.value = true
    await authStore.getCurrentUser()
    loading.value = false
  }
  
  // 初始化主题
  themeStore.initTheme()
})
</script>

<template>
  <div class="profile-container">
    <div class="profile-header">
      <h1>个人中心</h1>
    </div>
    
    <div v-if="loading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>
    
    <div v-else-if="authStore.user" class="profile-content">
      <div class="user-info">
        <div class="user-avatar">
          <img 
            v-if="authStore.user.avatar" 
            :src="authStore.user.avatar" 
            :alt="authStore.user.nickname"
            class="avatar-img"
          />
          <div v-else class="avatar-placeholder">
            {{ authStore.user.nickname.charAt(0).toUpperCase() || 'U' }}
          </div>
        </div>
        <div class="user-details">
          <h2>{{ authStore.user.nickname }}</h2>
          <p class="username">@{{ authStore.user.username }}</p>
          <p class="email">{{ authStore.user.email }}</p>
          <p v-if="authStore.user.role === 'admin'" class="role">管理员</p>
          <p class="created-at">注册时间: {{ new Date(authStore.user.createdAt).toLocaleString('zh-CN') }}</p>
        </div>
      </div>
      
      <!-- 主题设置 -->
      <div class="theme-settings">
        <h3>主题设置</h3>
        
        <!-- 主题配色选择 -->
        <div class="setting-group">
          <label>主题配色</label>
          <div class="theme-options">
            <div 
              v-for="theme in availableThemes" 
              :key="theme.name"
              class="theme-option"
              :class="{ active: currentTheme === theme.name }"
              @click="setTheme(theme.name)"
            >
              <div class="theme-preview">
                <div class="preview-primary" :style="{ backgroundColor: theme.colors.primary }"></div>
                <div class="preview-secondary" :style="{ backgroundColor: theme.colors.secondary }"></div>
                <div class="preview-accent" :style="{ backgroundColor: theme.colors.accent }"></div>
              </div>
              <span class="theme-name">{{ theme.name === 'default' ? '默认' : '幽蓝' }}</span>
            </div>
          </div>
        </div>
        
        <!-- 夜间模式设置 -->
        <div class="setting-group">
          <label>夜间模式</label>
          <div class="dark-mode-dropdown-container">
            <div class="dark-mode-dropdown">
              <button @click="toggleDarkModeDropdown" class="dark-mode-button">
                <span>{{ getDarkModeDisplayText(darkModePreference) }}</span>
                <svg class="dropdown-arrow" :class="{ 'rotated': showDarkModeDropdown }" width="12" height="12" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M6 9l6 6 6-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              </button>
              <div v-if="showDarkModeDropdown" class="dark-mode-options">
                <button @click="selectDarkModeOption('system')" class="dark-mode-option" :class="{ active: darkModePreference === 'system' }">
                  跟随系统
                </button>
                <button @click="selectDarkModeOption('enabled')" class="dark-mode-option" :class="{ active: darkModePreference === 'enabled' }">
                  一直启用
                </button>
                <button @click="selectDarkModeOption('disabled')" class="dark-mode-option" :class="{ active: darkModePreference === 'disabled' }">
                  一直禁用
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <div class="profile-actions">
        <button @click="authStore.logout()" class="logout-btn">
          退出登录
        </button>
      </div>
    </div>
    
    <div v-else class="error-container">
      <p>无法加载用户信息</p>
      <button @click="authStore.getCurrentUser()" class="retry-btn">
        重试
      </button>
    </div>
  </div>
</template>

<style scoped>
.profile-container {
  min-height: 100vh;
  padding: 20px;
  background-color: var(--color-background);
  font-family: sans-serif;
  color: var(--color-text);
}

.profile-header {
  text-align: center;
  margin-bottom: 20px;
}

.profile-header h1 {
  font-size: 32px;
  font-weight: bold;
  color: var(--color-text);
}

.profile-content {
  max-width: 600px;
  margin: 0 auto;
  padding: 24px;
  background-color: var(--color-surface);
  border-radius: 12px;
  border: 1px solid var(--color-border);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 24px;
}

.user-avatar {
  flex-shrink: 0;
}

.avatar-img {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid var(--color-border);
}

.avatar-placeholder {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background-color: var(--color-primary);
  color: var(--color-text);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 32px;
  border: 1px solid var(--color-border);
}

.user-details h2 {
  font-size: 24px;
  font-weight: bold;
  color: var(--color-text);
  margin-bottom: 6px;
}

.username {
  color: var(--color-text);
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 4px;
}

.email {
  color: var(--color-text);
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 4px;
}

.role {
  color: var(--color-text);
  font-size: 14px;
  font-weight: bold;
  margin-top: 6px;
}

.created-at {
  color: var(--color-textSecondary);
  font-size: 12px;
  margin-top: 6px;
}

.profile-actions {
  display: flex;
  justify-content: center;
}

.logout-btn {
  padding: 12px 24px;
  background-color: var(--color-secondary);
  color: var(--color-text);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  font-weight: bold;
  font-size: 14px;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.logout-btn:hover {
}

.loading-container,
.error-container {
  text-align: center;
  padding: 40px;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 2px solid var(--color-primary);
  border-top: 2px solid var(--color-accent);
  border-radius: 50%;
  margin: 0 auto 16px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.retry-btn {
  padding: 12px 24px;
  background-color: var(--color-primary);
  color: var(--color-text);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  font-weight: bold;
  font-size: 14px;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.retry-btn:hover {
}

/* 主题设置样式 */
.theme-settings {
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 20px;
  margin: 20px 0;
}

.theme-settings h3 {
  margin: 0 0 20px 0;
  color: var(--color-text);
  font-size: 18px;
}

.setting-group {
  margin-bottom: 24px;
}

.setting-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: bold;
  color: var(--color-text);
}

.theme-options {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.theme-option {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  background-color: var(--color-surface);
}

.theme-option:hover {
}

.theme-option.active {
  border-color: var(--color-primary);
  box-shadow: none;
}

.theme-preview {
  display: flex;
  gap: 4px;
  margin-bottom: 8px;
}

.preview-primary,
.preview-secondary,
.preview-accent {
  width: 16px;
  height: 16px;
  border-radius: 4px;
  border: 1px solid var(--color-border);
}

.theme-name {
  font-size: 14px;
  color: var(--color-text);
}

/* 夜间模式下拉选择样式 */
.dark-mode-dropdown-container {
  width: 100%;
  max-width: 200px;
}

.dark-mode-dropdown {
  position: relative;
}

.dark-mode-button {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  background: #FFFFFF;
  color: #1A1A1A;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 10px 12px;
  font-size: 14px;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.2s ease;
  width: 100%;
  text-align: left;
  box-shadow: none;
}

.dark-mode-button:hover {
  background: #F7D074;
  border-color: #F0C74A;
}

.dark-mode-button:focus {
  outline: none;
}

.dropdown-arrow {
  transition: transform 0.2s ease;
  flex-shrink: 0;
}

.dropdown-arrow.rotated {
  transform: rotate(180deg);
}

.dark-mode-options {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  z-index: 1000;
  background: #FFFFFF;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  margin-top: 4px;
  overflow: hidden;
  box-shadow: none;
}

.dark-mode-option {
  display: block;
  width: 100%;
  background: #FFFFFF;
  border: none;
  padding: 10px 12px;
  font-size: 14px;
  font-weight: bold;
  color: #1A1A1A;
  cursor: pointer;
  transition: background-color 0.2s ease;
  text-align: left;
  border-bottom: 1px solid #E0E0E0;
}

.dark-mode-option:last-child {
  border-bottom: none;
}

.dark-mode-option:hover {
  background: #F7D074;
}

.dark-mode-option.active {
  background: #F7D074;
  color: #1A1A1A;
}

/* Responsive Design */
@media (max-width: 768px) {
  .profile-container {
    padding: 16px;
  }
  
  .dark-mode-dropdown-container {
    max-width: 100%;
  }
  
  .dark-mode-button {
    padding: 12px 16px;
    font-size: 16px;
  }
  
  .dark-mode-option {
    padding: 12px 16px;
    font-size: 16px;
  }
  
  .user-info {
    flex-direction: column;
    text-align: center;
    gap: 16px;
  }
  
  .profile-content {
    padding: 20px;
  }
}
</style>