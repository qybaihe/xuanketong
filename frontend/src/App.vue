<script setup lang="ts">
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router'
import { computed, onMounted, ref, nextTick } from 'vue'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
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
  background: var(--background-blur);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid var(--separator-color);
  padding: var(--spacing-md) var(--spacing-lg);
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
  gap: var(--spacing-lg);
}

.nav-link {
  padding: var(--spacing-sm) var(--spacing-md);
  border-radius: var(--border-radius-base);
  text-decoration: none;
  color: var(--text-secondary);
  font-weight: var(--font-weight-medium);
  transition: all 0.2s ease;
}

.nav-link:hover {
  background: var(--background-secondary);
  color: var(--text-primary);
}

.nav-link.router-link-exact-active {
  background: var(--primary-color);
  color: white;
}

.header-right {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.search-section {
  position: relative;
}

.search-toggle {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: #FFFFFF;
  border: 3px solid #000000;
  border-radius: 8px;
  box-shadow: 4px 4px 0px 0px #000000;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.search-toggle:hover {
  transform: translate(-2px, -2px);
  box-shadow: 6px 6px 0px 0px #000000;
}

.search-icon {
  font-size: 16px;
}

.search-text {
  color: #1A1A1A;
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
  background: #FFFFFF;
  border: 3px solid #000000;
  border-radius: 8px;
  padding: 8px;
  box-shadow: 4px 4px 0px 0px #000000;
  height: 48px;
  min-width: 400px;
  transform: translate(-50%, -50%);
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
  background: #FFFFFF;
  border: 2px solid #000000;
  border-radius: 6px;
  padding: 6px 10px;
  font-size: 14px;
  font-weight: bold;
  color: #1A1A1A;
  cursor: pointer;
  transition: background-color 0.2s ease;
  height: 32px;
  white-space: nowrap;
}

.search-type-button:hover {
  background: #F7D074;
}

.dropdown-arrow {
  transition: transform 0.2s ease;
}

.dropdown-arrow.rotated {
  transform: rotate(180deg);
}

.search-type-options {
  position: absolute;
  top: 100%;
  left: 0;
  z-index: 1000;
  background: #FFFFFF;
  border: 3px solid #000000;
  border-radius: 8px;
  box-shadow: 4px 4px 0px 0px #000000;
  margin-top: 4px;
  min-width: 100%;
  overflow: hidden;
}

.search-type-option {
  display: block;
  width: 100%;
  background: #FFFFFF;
  border: none;
  padding: 8px 12px;
  font-size: 14px;
  font-weight: bold;
  color: #1A1A1A;
  cursor: pointer;
  transition: background-color 0.2s ease;
  text-align: left;
}

.search-type-option:hover {
  background: #F7D074;
}

.search-type-option.active {
  background: #F7D074;
}

.search-input {
  border: none;
  outline: none;
  padding: 8px;
  font-size: 14px;
  font-weight: bold;
  min-width: 200px;
  background: transparent;
  color: #1A1A1A;
  height: 100%;
  flex: 1;
}

.search-input::placeholder {
  color: #888888;
  font-weight: normal;
}

.search-button {
  background: #F7D074;
  color: #1A1A1A;
  border: 2px solid #000000;
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
  background: #E6C066;
}

.search-cancel {
  background: #FFFFFF;
  color: #1A1A1A;
  border: 2px solid #000000;
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
  background: #F0F0F0;
}

.mobile-search {
  display: none;
}

.mobile-search-btn {
  background: #FFFFFF;
  border: 3px solid #000000;
  border-radius: 8px;
  box-shadow: 4px 4px 0px 0px #000000;
  cursor: pointer;
  padding: 8px;
  transition: transform 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.mobile-search-btn:hover {
  transform: translate(-2px, -2px);
  box-shadow: 6px 6px 0px 0px #000000;
}

.auth-section {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.welcome-text {
  color: var(--text-secondary);
  font-size: var(--font-size-body2);
}

.logout-btn,
.login-btn {
  padding: var(--spacing-sm) var(--spacing-md);
  border-radius: var(--border-radius-base);
  text-decoration: none;
  font-weight: var(--font-weight-medium);
  transition: all 0.2s ease;
  cursor: pointer;
}

.logout-btn {
  background: var(--background-secondary);
  color: var(--text-primary);
  border: 1px solid var(--separator-color);
}

.logout-btn:hover {
  background: var(--background-tertiary);
}

.login-btn {
  background: var(--gradient-primary);
  color: #1A1A1A;
  border: 2px solid #000000;
  /*box-shadow: 0 2px 8px rgba(47, 169, 20, 0.2);*/
}

.login-btn:hover {
  transform: translateY(-1px);
  /*box-shadow: 0 4px 12px rgba(47, 169, 20, 0.3);*/
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
    background: #FFFFFF;
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
