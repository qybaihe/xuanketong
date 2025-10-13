<template>
  <el-container class="admin-container">
    <!-- 顶部标题栏 -->
    <el-header class="admin-header">
      <el-button
        class="mobile-menu-btn"
        type="text"
        @click="drawerVisible = true"
        circle
        plain
      >
        <el-icon><Menu /></el-icon>
      </el-button>
      <div class="header-content">
        <h1 class="page-title">管理后台</h1>
      </div>
    </el-header>

    <!-- 主要内容区域 -->
    <el-container>
      <!-- 侧边栏 -->
      <el-aside width="250px" class="admin-sidebar desktop-only">
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
      <!-- 移动端抽屉侧栏 -->
      <el-drawer
        v-model="drawerVisible"
        size="200px"
        :with-header="false"
        class="mobile-only"
      >
        <el-menu
          :default-active="activeMenu"
          class="admin-menu"
          router
          background-color="var(--color-background)"
          text-color="var(--color-text-secondary)"
          active-text-color="var(--color-primary)"
          @close="drawerVisible = false"
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
      </el-drawer>

      <el-main class="admin-main">
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onBeforeUnmount } from 'vue'
import { useRoute } from 'vue-router'
import { watch } from 'vue'

const drawerVisible = ref(false)

// 监听路由变换自动收起抽屉
const route = useRoute()

const activeMenu = computed(() => {
  return route.path
})

// 当路由切换时关闭移动抽屉
watch(() => route.path, () => {
  drawerVisible.value = false
})
</script>

<style scoped>
.admin-container {
  height: 100vh;
  background-color: var(--color-background);
  flex-direction: column;
}

.admin-header {
  background-color: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
  padding: 0 20px;
  display: flex;
  align-items: center;
  height: 60px;
  flex-shrink: 0;
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

.admin-sidebar {
  background-color: var(--color-background);
  border-right: 1px solid var(--color-border);
  height: calc(100vh - 60px);
}

.admin-menu {
  border-right: none;
  height: 100%;
}

.admin-main {
  background-color: var(--color-background);
  padding: 20px;
  height: calc(100vh - 60px);
  overflow-y: auto;
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
    height: auto;
  }
}

/* 隐藏 PC 端不需要的菜单按钮 */
.mobile-menu-btn {
  display: none;
  margin-right: 10px;
}

/* 工具类 */
.desktop-only { display: block; }
.mobile-only { display: none; }

@media (max-width: 576px) {
  .mobile-menu-btn { display: inline-flex; }
  .desktop-only { display: none !important; }
  .mobile-only { display: block !important; }
}
</style>