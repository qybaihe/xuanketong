<script setup lang="ts">
import { RouterLink, RouterView, useRoute } from 'vue-router'
import { computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const authStore = useAuthStore()
const isAdminRoute = computed(() => route.path.startsWith('/admin'))

const currentUserName = computed(() => authStore.user?.nickname || authStore.user?.username || '游客')

onMounted(async () => {
  // Try to get current user if token exists
  await authStore.getCurrentUser()
})

const handleLogout = async () => {
  await authStore.logout()
}
</script>

<template>
  <header v-if="!isAdminRoute" class="app-header">
    <div class="header-content">
      <nav class="main-nav">
        <RouterLink to="/" class="nav-link">首页</RouterLink>
        <RouterLink to="/stats" class="nav-link">数据统计</RouterLink>
        <RouterLink to="/about" class="nav-link">关于</RouterLink>
        <RouterLink to="/profile" class="nav-link">个人中心</RouterLink>
        <RouterLink to="/admin" v-if="authStore.isAdmin" class="nav-link">管理后台</RouterLink>
      </nav>
      
      <div class="auth-section">
        <template v-if="authStore.isAuthenticated">
          <span class="welcome-text">欢迎，{{ currentUserName }}</span>
          <button @click="handleLogout" class="logout-btn btn-secondary">
            退出登录
          </button>
        </template>
        <template v-else>
          <RouterLink to="/auth" class="login-btn btn-primary">
            登录 / 注册
          </RouterLink>
        </template>
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
  color: white;
  box-shadow: 0 2px 8px rgba(47, 169, 20, 0.2);
}

.login-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(47, 169, 20, 0.3);
}

.main-content {
  min-height: calc(100vh - 80px); /* Account for header height */
}

/* Responsive Design */
@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
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
}
</style>
