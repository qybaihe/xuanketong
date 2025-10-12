<script setup lang="ts">
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router'
import { computed, onMounted, ref, nextTick } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const themeStore = useThemeStore()
const isAdminRoute = computed(() => route.path.startsWith('/admin'))
const isHomeRoute = computed(() => route.path === '/')

const currentUserName = computed(() => authStore.user?.nickname || authStore.user?.username || '游客')

// 搜索相关状态
const searchQuery = ref('')
const searchType = ref<'course' | 'teacher'>('course')
const showSearchInput = ref(false)
const showDropdown = ref(false)

onMounted(async () => {
  // Try to get current user if token exists
  await authStore.getCurrentUser()
  
  // 初始化主题
  themeStore.initTheme()
})

const handleLogout = async () => {
  await authStore.logout()
}

// 搜索功能
const handleSearch = () => {
  if (!searchQuery.value.trim()) return
  
  router.push({
    path: '/search',
    query: {
      q: searchQuery.value,
      type: searchType.value
    }
  })
  
  // 清空搜索框
  searchQuery.value = ''
  showSearchInput.value = false
}

const toggleSearchInput = () => {
  showSearchInput.value = !showSearchInput.value
  if (showSearchInput.value) {
    // 移动端搜索展开时添加class
    if (window.innerWidth <= 768) {
      document.querySelector('.header-content')?.classList.add('search-expanded')
    }
    nextTick(() => {
      const input = document.querySelector('.search-input') as HTMLInputElement
      if (input) {
        input.focus()
      }
    })
  } else {
    // 移动端搜索收起时移除class
    if (window.innerWidth <= 768) {
      document.querySelector('.header-content')?.classList.remove('search-expanded')
    }
    searchQuery.value = ''
    showDropdown.value = false
  }
}

const toggleDropdown = () => {
  showDropdown.value = !showDropdown.value
}

const selectSearchType = (type: 'course' | 'teacher') => {
  searchType.value = type
  showDropdown.value = false
}
</script>

<template>
  <header v-if="!isAdminRoute" class="app-header">
    <div class="header-content">
      <nav class="main-nav">
        <RouterLink to="/" class="nav-link">首页</RouterLink>
        <!--
        <RouterLink to="/stats" class="nav-link">数据统计</RouterLink>
        <RouterLink to="/about" class="nav-link">关于</RouterLink>
        -->
        <RouterLink to="/evaluation-requests" class="nav-link">求评价</RouterLink>
        <RouterLink to="/profile" v-if="authStore.isAuthenticated" class="nav-link">个人中心</RouterLink>
        <RouterLink to="/admin" v-if="authStore.isAdmin" class="nav-link">管理后台</RouterLink>
      </nav>
      
      <div class="header-right">
        <!-- 搜索功能 - 只在非首页显示 -->
        <div v-if="!isHomeRoute" class="search-section">
          <!-- 桌面端搜索框 -->
          <div v-if="!showSearchInput" class="search-toggle" @click="toggleSearchInput">
            <svg class="search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <circle cx="11" cy="11" r="8" stroke="currentColor" stroke-width="2"/>
              <path d="m21 21-4.35-4.35" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            <span class="search-text">搜索课程</span>
          </div>
          
          <!-- 搜索输入框 -->
          <div v-if="showSearchInput" class="search-input-container">
            <div class="search-type-dropdown">
              <button @click="toggleDropdown" class="search-type-button">
                <span>{{ searchType === 'course' ? '课程' : '老师' }}</span>
                <svg class="dropdown-arrow" :class="{ 'rotated': showDropdown }" width="12" height="12" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M6 9l6 6 6-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              </button>
              <div v-if="showDropdown" class="search-type-options">
                <button @click="selectSearchType('course')" class="search-type-option" :class="{ active: searchType === 'course' }">
                  课程
                </button>
                <button @click="selectSearchType('teacher')" class="search-type-option" :class="{ active: searchType === 'teacher' }">
                  老师
                </button>
              </div>
            </div>
            <input
              v-model="searchQuery"
              @keyup.enter="handleSearch"
              type="text"
              class="search-input"
              :placeholder="searchType === 'course' ? '搜索课程...' : '搜索老师...'"
              ref="searchInput"
            />
            <button @click="handleSearch" class="search-button">搜索</button>
            <button @click="toggleSearchInput" class="search-cancel">取消</button>
          </div>
        </div>
        
        <!-- 移动端搜索按钮 - 只在非首页显示 -->
        <div v-if="!isHomeRoute" class="mobile-search">
          <button @click="toggleSearchInput" class="mobile-search-btn">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <circle cx="11" cy="11" r="8" stroke="currentColor" stroke-width="2"/>
              <path d="m21 21-4.35-4.35" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </button>
          
          <!-- 移动端搜索输入框 -->
          <div v-if="showSearchInput" class="search-input-container">
            <div class="search-type-dropdown">
              <button @click="toggleDropdown" class="search-type-button">
                {{ searchType === 'course' ? '课程' : '老师' }}
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" :class="{ 'rotate-180': showDropdown }">
                  <path d="m6 9 6 6 6-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              </button>
              <div v-if="showDropdown" class="search-type-options">
                <button @click="selectSearchType('course')" class="search-type-option" :class="{ active: searchType === 'course' }">
                  课程
                </button>
                <button @click="selectSearchType('teacher')" class="search-type-option" :class="{ active: searchType === 'teacher' }">
                  老师
                </button>
              </div>
            </div>
            
            <input
              v-model="searchQuery"
              @keyup.enter="handleSearch"
              type="text"
              :placeholder="searchType === 'course' ? '搜索课程...' : '搜索老师...'"
              class="search-input"
            />
            
            <button @click="handleSearch" class="search-button">
              搜索
            </button>
            
            <button @click="toggleSearchInput" class="search-cancel">
              取消
            </button>
          </div>
        </div>
        
        <div class="auth-section">
          <template v-if="authStore.isAuthenticated">
            <span class="welcome-text desktop-only">欢迎，{{ currentUserName }}</span>
            <button @click="handleLogout" class="logout-btn btn-secondary desktop-only">
              退出登录
            </button>
          </template>
          <template v-else>
            <RouterLink to="/auth" class="login-btn btn-primary">
              登录
            </RouterLink>
          </template>
        </div>
      </div>
    </div>
  </header>

  <main class="main-content">
    <RouterView />
  </main>
</template>

<style scoped>
.app-header {
  background: var(--color-surface);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid var(--color-border);
  padding: 16px 24px;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.main-nav {
  display: flex;
  gap: 24px;
}

.nav-link {
  padding: 8px 16px;
  border-radius: 8px;
  text-decoration: none;
  color: var(--color-text-secondary);
  font-weight: 500;
  transition: all 0.2s ease;
}

.nav-link:hover {
  background: var(--color-background);
  color: var(--color-text);
}

.nav-link.router-link-exact-active {
  background: var(--color-primary);
  color: white;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.search-section {
  position: relative;
}

.search-toggle {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.search-toggle:hover {
}

.search-icon {
  font-size: 16px;
}

.search-text {
  color: var(--color-text);
  font-size: 14px;
  font-weight: bold;
}

.search-input-container {
  position: absolute;
  top: 50%;
  left: 50%;
  z-index: 1000;
  display: flex;
  align-items: center;
  gap: 8px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 8px;
  height: 48px;
  min-width: 400px;
  animation: expandWidth 0.3s ease-out;
}

@keyframes expandWidth {
  from {
    width: 0;
    opacity: 0;
  }
  to {
    width: auto;
    opacity: 1;
  }
}

.search-type-dropdown {
  position: relative;
}

.search-type-button {
  display: flex;
  align-items: center;
  gap: 4px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 6px;
  padding: 6px 10px;
  font-size: 14px;
  font-weight: bold;
  color: var(--color-text);
  cursor: pointer;
  transition: background-color 0.2s ease;
  height: 32px;
  white-space: nowrap;
}

.search-type-button:hover {
  background: var(--color-primary);
}

.dropdown-arrow {
  transition: transform 0.2s ease;
}

.dropdown-arrow.rotated {
}

.search-type-options {
  position: absolute;
  top: 100%;
  left: 0;
  z-index: 1000;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  margin-top: 4px;
  min-width: 100%;
  overflow: hidden;
}

.search-type-option {
  display: block;
  width: 100%;
  background: var(--color-surface);
  border: none;
  padding: 8px 12px;
  font-size: 14px;
  font-weight: bold;
  color: var(--color-text);
  cursor: pointer;
  transition: background-color 0.2s ease;
  text-align: left;
}

.search-type-option:hover {
  background: var(--color-primary);
}

.search-type-option.active {
  background: var(--color-primary);
}

.search-input {
  border: none;
  outline: none !important;
  padding: 8px;
  font-size: 14px;
  font-weight: bold;
  min-width: 200px;
  background: transparent;
  color: var(--color-text);
  height: 100%;
  flex: 1;
  box-shadow: none;
}

.search-input::placeholder {
  color: var(--color-textSecondary);
  font-weight: normal;
}

.search-button {
  background: var(--color-primary);
  color: var(--color-text);
  border: 1px solid var(--color-border);
  border-radius: 6px;
  padding: 6px 12px;
  font-size: 14px;
  font-weight: bold;
  cursor: pointer;
  transition: background-color 0.2s ease;
  height: 32px;
  white-space: nowrap;
}

.search-button:hover {
  background: var(--color-primary);
}

.search-cancel {
  background: var(--color-surface);
  color: var(--color-text);
  border: 1px solid var(--color-border);
  border-radius: 6px;
  padding: 6px 12px;
  font-size: 14px;
  font-weight: bold;
  cursor: pointer;
  transition: background-color 0.2s ease;
  height: 32px;
  white-space: nowrap;
}

.search-cancel:hover {
  background: var(--color-background);
}

.mobile-search {
  display: none;
}

.mobile-search-btn {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  cursor: pointer;
  padding: 8px;
  transition: transform 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.mobile-search-btn:hover {
}

.auth-section {
  display: flex;
  align-items: center;
  gap: 16px;
}

.welcome-text {
  color: var(--color-text-secondary);
  font-size: 14px;
}

.logout-btn,
.login-btn {
  padding: 8px 16px;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 500;
  transition: none;
  cursor: pointer;
}

.logout-btn {
  background: var(--color-background);
  color: var(--color-text);
  border: 1px solid var(--color-border);
}

.logout-btn:hover,
.logout-btn:active,
.logout-btn:focus {
  background: var(--color-background);
  box-shadow: none;
  outline: none;
  transform: none;
}

.login-btn {
  background: var(--color-primary);
  color: var(--color-text);
  border: 1px solid var(--color-border);
}

.login-btn:hover,
.login-btn:active,
.login-btn:focus {
  background: var(--color-primary);
  box-shadow: none;
  outline: none;
  transform: none;
}

.main-content {
  min-height: calc(100vh - 80px); /* Account for header height */
}

/* Desktop-only elements */
.desktop-only {
  display: block;
}

/* Responsive Design */
@media (max-width: 768px) {
  .header-content {
    gap: var(--spacing-md);
  }
  
  .main-nav {
    flex-wrap: wrap;
    justify-content: center;
  }
  
  .nav-link {
    padding: var(--spacing-xs) var(--spacing-sm);
    font-size: var(--font-size-body2);
  }
  
  .search-section {
    display: none;
  }
  
  .mobile-search {
    display: block;
  }
  
  .desktop-only {
    display: none;
  }
  
  .mobile-search .search-input-container {
    position: absolute !important;
    top: 0 !important;
    left: 0 !important;
    right: 0 !important;
    bottom: 0 !important;
    z-index: 1000;
    background: var(--color-surface);
    width: 100% !important;
    height: 100% !important;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    padding: 0 16px !important;
    margin: 0 !important;
    border: none !important;
    box-shadow: none !important;
    border-radius: 0 !important;
    transform: none !important;
    min-width: auto !important;
    animation: none !important;
  }
  
  .mobile-search .search-input {
    min-width: 120px;
    flex: 1;
  }
  
  .mobile-search .search-type-button {
    min-width: 60px;
    padding: 4px 8px;
    font-size: 12px;
  }
  
  .mobile-search .search-button,
  .mobile-search .search-cancel {
    padding: 4px 8px;
    font-size: 12px;
    min-width: 50px;
  }
  
  /* 移动端搜索展开时隐藏其他元素 */
  .header-content.search-expanded .main-nav,
  .header-content.search-expanded .auth-section {
    opacity: 0;
    pointer-events: none;
  }
}

@keyframes expandFullWidth {
  from {
    width: 0;
    opacity: 0;
    transform: translateX(100%);
  }
  to {
    width: 100%;
    opacity: 1;
    transform: translateX(0);
  }
}
</style>
