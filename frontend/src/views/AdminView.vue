<template>
  <div class="admin-container">
    <div class="container-fluid">
      <div class="row">
        <!-- 侧边栏 -->
        <nav class="col-md-3 col-lg-2 d-md-block sidebar collapse">
          <div class="position-sticky pt-3">
            <ul class="nav flex-column">
              <li class="nav-item">
                <RouterLink to="/admin/courses" class="nav-link" active-class="active">
                  <i class="bi bi-book me-2"></i>
                  课程管理
                </RouterLink>
              </li>
              <li class="nav-item">
                <RouterLink to="/admin/users" class="nav-link" active-class="active">
                  <i class="bi bi-people me-2"></i>
                  用户管理
                </RouterLink>
              </li>
              <li class="nav-item">
                <RouterLink to="/admin/stats" class="nav-link" active-class="active">
                  <i class="bi bi-graph-up me-2"></i>
                  数据统计
                </RouterLink>
              </li>
            </ul>
          </div>
        </nav>

        <!-- 主要内容区域 -->
        <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
          <div class="main-header mb-4">
            <div class="d-flex justify-content-between align-items-center">
              <h1 class="mb-0">课程管理系统</h1>
              <button type="button" class="btn btn-sm btn-outline-secondary">
                <i class="bi bi-download me-1"></i>导出
              </button>
            </div>
          </div>
          <router-view></router-view>
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterView, RouterLink } from 'vue-router'
</script>

<style scoped>

.admin-container {
  background-color: #ffffff;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  font-size: 14px;
  color: #333333;
}
.sidebar {
  height: 100vh;
  background: #f8f9fa;
  border-right: 1px solid #dee2e6;
  padding-top: 1rem;
}
.sidebar .nav-link {
  color: #495057;
  padding: 0.75rem 1.25rem;
  font-weight: 500;
  transition: all 0.2s;
  border-radius: 4px;
}
.sidebar .nav-link:hover,
.sidebar .nav-link.active {
  color: #007bff;
  background-color: #e3f2fd;
}
.sidebar .nav-link i {
  width: 20px;
  margin-right: 0.75rem;
}
.main-header {
  background: #ffffff;
  border-bottom: 1px solid #dee2e6;
  padding: 1rem 0;
}
.main-header h1 {
  font-size: 1.5rem;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
}
</style>