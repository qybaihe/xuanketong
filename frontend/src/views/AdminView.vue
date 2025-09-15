<template>
  <div class="admin-container theme-natural">
    <!-- 顶部导航栏 -->
    <nav class="admin-topbar card-glass">
      <div class="topbar-content">
        <div class="topbar-left">
          <button class="menu-toggle btn btn-glass" @click="toggleSidebar">
            <span class="menu-icon">☰</span>
          </button>
          <div class="logo">
            <span class="logo-text text-shine">选课通管理</span>
          </div>
        </div>
        
        <div class="topbar-right">
          <!-- 搜索框 -->
          <div class="search-container">
            <input type="text" class="search-input input-glass" placeholder="搜索...">
            <button class="search-btn btn btn-glass">
              <span class="search-icon">🔍</span>
            </button>
          </div>
          
          <!-- 通知图标 -->
          <div class="notification-icons">
            <button class="notification-btn btn btn-glass">
              <span class="notification-icon">💬</span>
              <span class="notification-badge">3</span>
            </button>
            <button class="notification-btn btn btn-glass">
              <span class="notification-icon">🔔</span>
              <span class="notification-badge">15</span>
            </button>
          </div>
          
          <!-- 用户菜单 -->
          <div class="user-menu">
            <div class="user-avatar">
              <span class="user-initial">管</span>
            </div>
            <span class="user-name">管理员</span>
          </div>
        </div>
      </div>
    </nav>
    
    <div class="admin-layout">
      <!-- 侧边导航栏 -->
      <aside class="admin-sidebar card-glass" :class="{ 'sidebar-collapsed': sidebarCollapsed }">
        <div class="sidebar-header">
          <h2 class="sidebar-title">管理菜单</h2>
        </div>
        
        <nav class="sidebar-nav">
          <ul class="nav-list">
            <li class="nav-item">
              <RouterLink to="/admin" class="nav-link" active-class="active">
                <span class="nav-icon">📊</span>
                <span class="nav-text">仪表盘</span>
              </RouterLink>
            </li>
            <li class="nav-item">
              <RouterLink to="/admin/courses" class="nav-link" active-class="active">
                <span class="nav-icon">📚</span>
                <span class="nav-text">课程管理</span>
              </RouterLink>
            </li>
            <li class="nav-item">
              <RouterLink to="/admin/users" class="nav-link" active-class="active">
                <span class="nav-icon">👥</span>
                <span class="nav-text">用户管理</span>
              </RouterLink>
            </li>
            <li class="nav-item">
              <RouterLink to="/admin/comments" class="nav-link" active-class="active">
                <span class="nav-icon">💬</span>
                <span class="nav-text">评论管理</span>
              </RouterLink>
            </li>
            <li class="nav-item">
              <RouterLink to="/admin/ratings" class="nav-link" active-class="active">
                <span class="nav-icon">⭐</span>
                <span class="nav-text">评分管理</span>
              </RouterLink>
            </li>
          </ul>
        </nav>
      </aside>
      
      <!-- 主内容区域 -->
      <main class="admin-main">
        <div class="main-content">
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>
      </main>
    </div>
    
    <!-- 底部 -->
    <footer class="admin-footer card-glass">
      <div class="footer-content">
        <span class="footer-text">© 2024 选课通管理系统 - 保留所有权利</span>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterView, RouterLink } from 'vue-router'

const sidebarCollapsed = ref(false)

// 切换侧边栏
const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

// 检查屏幕尺寸，自动折叠侧边栏
const checkScreenSize = () => {
  if (window.innerWidth < 768) {
    sidebarCollapsed.value = true
  } else {
    sidebarCollapsed.value = false
  }
}

onMounted(() => {
  checkScreenSize()
  window.addEventListener('resize', checkScreenSize)
})
</script>

<style scoped>
/* 主容器 */
.admin-container {
  min-height: 100vh;
  background: linear-gradient(135deg, var(--background-secondary) 0%, var(--natural-background) 100%);
  display: flex;
  flex-direction: column;
}

/* 顶部导航栏 */
.admin-topbar {
  position: sticky;
  top: 0;
  z-index: 100;
  padding: var(--spacing-md) var(--spacing-lg);
  border-radius: 0;
  margin-bottom: var(--spacing-md);
}

.topbar-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1400px;
  margin: 0 auto;
}

.topbar-left {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.menu-toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  padding: 0;
  border-radius: 8px;
}

.menu-icon {
  font-size: 18px;
}

.logo-text {
  font-size: var(--font-size-title2);
  font-weight: var(--font-weight-bold);
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.topbar-right {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.search-container {
  display: flex;
  align-items: center;
  position: relative;
}

.search-input {
  width: 200px;
  padding-right: 40px;
}

.search-btn {
  position: absolute;
  right: 4px;
  width: 32px;
  height: 32px;
  padding: 0;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.search-icon {
  font-size: 16px;
}

.notification-icons {
  display: flex;
  gap: var(--spacing-xs);
}

.notification-btn {
  position: relative;
  width: 36px;
  height: 36px;
  padding: 0;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.notification-icon {
  font-size: 16px;
}

.notification-badge {
  position: absolute;
  top: 0;
  right: 0;
  width: 16px;
  height: 16px;
  background: var(--error-base);
  color: white;
  font-size: 10px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.user-menu {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--gradient-primary);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: var(--font-weight-semibold);
}

.user-name {
  font-size: var(--font-size-body);
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
}

/* 管理布局 */
.admin-layout {
  display: flex;
  flex: 1;
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
  padding: 0 var(--spacing-md);
}

/* 侧边导航栏 */
.admin-sidebar {
  width: 240px;
  flex-shrink: 0;
  padding: var(--spacing-md);
  border-radius: 16px;
  margin-right: var(--spacing-md);
  height: fit-content;
  transition: all 0.3s ease;
}

.sidebar-collapsed {
  width: 60px;
}

.sidebar-header {
  margin-bottom: var(--spacing-lg);
}

.sidebar-title {
  font-size: var(--font-size-title3);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin: 0;
}

.nav-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.nav-item {
  margin-bottom: var(--spacing-xs);
}

.nav-link {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-sm) var(--spacing-md);
  border-radius: 8px;
  color: var(--text-secondary);
  text-decoration: none;
  transition: all 0.2s ease;
}

.nav-link:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--text-primary);
}

.nav-link.active {
  background: var(--gradient-primary);
  color: white;
}

.nav-icon {
  font-size: 18px;
  width: 24px;
  text-align: center;
}

.nav-text {
  font-size: var(--font-size-body);
  font-weight: var(--font-weight-medium);
}

.sidebar-collapsed .nav-text,
.sidebar-collapsed .sidebar-title {
  display: none;
}

/* 主内容区域 */
.admin-main {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  padding: var(--spacing-md);
  background: rgba(255, 255, 255, 0.05);
  border-radius: 16px;
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
}

/* 底部 */
.admin-footer {
  padding: var(--spacing-md) var(--spacing-lg);
  border-radius: 0;
  margin-top: var(--spacing-md);
}

.footer-content {
  max-width: 1400px;
  margin: 0 auto;
  text-align: center;
}

.footer-text {
  font-size: var(--font-size-body2);
  color: var(--text-tertiary);
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .admin-layout {
    flex-direction: column;
  }
  
  .admin-sidebar {
    width: 100%;
    margin-right: 0;
    margin-bottom: var(--spacing-md);
    position: relative;
    height: auto;
  }
  
  .sidebar-collapsed {
    height: 60px;
    overflow: hidden;
  }
  
  .search-input {
    width: 120px;
  }
  
  .user-name {
    display: none;
  }
  
  .topbar-right {
    gap: var(--spacing-xs);
  }
}

@media (min-width: 769px) and (max-width: 1024px) {
  .admin-sidebar {
    width: 200px;
  }
  
  .sidebar-collapsed {
    width: 60px;
  }
}

/* 暗色主题支持 */
@media (prefers-color-scheme: dark) {
  .admin-container {
    background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  }
  
  .logo-text {
    background: linear-gradient(135deg, #4fc830 0%, #2fa914 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }
  
  .card-glass {
    background: rgba(0, 0, 0, 0.6);
    border: 1px solid rgba(255, 255, 255, 0.1);
  }
  
  .main-content {
    background: rgba(0, 0, 0, 0.4);
  }
}
</style>