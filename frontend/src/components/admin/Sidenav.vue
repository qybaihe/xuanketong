<template>
  <aside class="app-sidebar bg-body-secondary shadow" data-bs-theme="dark">
    <div class="sidebar-brand">
      <a href="#" class="brand-link">
        <img src="@/assets/img/AdminLTELogo.png" alt="AdminLTE Logo" class="brand-image opacity-75 shadow">
        <span class="brand-text fw-light">AdminLTE 4</span>
      </a>
    </div>
    <div class="sidebar-wrapper">
      <nav class="mt-2">
        <ul class="nav sidebar-menu flex-column" data-lte-toggle="treeview" role="menu" data-accordion="false">
          <li class="nav-item">
            <RouterLink to="/admin/dashboard" class="nav-link" active-class="active">
              <i class="nav-icon bi bi-speedometer"></i>
              <p>仪表盘</p>
            </RouterLink>
          </li>
          <li class="nav-item">
            <RouterLink to="/admin/courses" class="nav-link" active-class="active">
              <i class="nav-icon bi bi-book"></i>
              <p>课程管理</p>
            </RouterLink>
          </li>
          <li class="nav-item">
            <RouterLink to="/admin/users" class="nav-link" active-class="active">
              <i class="nav-icon bi bi-people"></i>
              <p>用户管理</p>
            </RouterLink>
          </li>
          <li class="nav-item">
            <RouterLink to="/admin/comments" class="nav-link" active-class="active">
              <i class="nav-icon bi bi-chat-dots"></i>
              <p>评论管理</p>
            </RouterLink>
          </li>
          <li class="nav-item">
            <RouterLink to="/admin/ratings" class="nav-link" active-class="active">
              <i class="nav-icon bi bi-star"></i>
              <p>评分管理</p>
            </RouterLink>
          </li>
        </ul>
      </nav>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>