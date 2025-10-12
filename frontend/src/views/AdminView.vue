<template>
  <el-container class="admin-container">
    <!-- 侧边栏 -->
    <el-aside width="250px" class="admin-sidebar">
      <el-menu
        :default-active="activeMenu"
        class="admin-menu"
        router
        background-color="var(--color-background)"
        text-color="var(--color-text-secondary)"
        active-text-color="var(--color-primary)"
      >
        <el-menu-item index="/admin">
          <el-icon><Odometer /></el-icon>
          <span>仪表盘</span>
        </el-menu-item>
        <el-menu-item index="/admin/courses">
          <el-icon><Reading /></el-icon>
          <span>课程管理</span>
        </el-menu-item>
        <el-menu-item index="/admin/users">
          <el-icon><User /></el-icon>
          <span>用户管理</span>
        </el-menu-item>
        <el-menu-item index="/admin/stats">
          <el-icon><TrendCharts /></el-icon>
          <span>数据统计</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <!-- 主要内容区域 -->
    <el-container>
      <el-header class="admin-header">
        <div class="header-content">
          <h1 class="page-title">课程管理系统</h1>
        </div>
      </el-header>
      
      <el-main class="admin-main">
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const activeMenu = computed(() => {
  return route.path
})
</script>

<style scoped>
.admin-container {
  height: 100vh;
  background-color: var(--color-background);
}

.admin-sidebar {
  background-color: var(--color-background);
  border-right: 1px solid var(--color-border);
}

.admin-menu {
  border-right: none;
  height: 100%;
}

.admin-header {
  background-color: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
  padding: 0 20px;
  display: flex;
  align-items: center;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.page-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0;
}

.admin-main {
  background-color: var(--color-background);
  padding: 20px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .admin-sidebar {
    width: 200px !important;
  }
}

@media (max-width: 576px) {
  .admin-container {
    flex-direction: column;
  }
  
  .admin-sidebar {
    width: 100% !important;
    height: auto;
  }
  
  .admin-main {
    padding: 10px;
  }
}
</style>