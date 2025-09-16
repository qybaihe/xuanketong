<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const loading = ref(false)

onMounted(async () => {
  if (!authStore.user) {
    loading.value = true
    await authStore.getCurrentUser()
    loading.value = false
  }
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
    
    <div v-else-if="authStore.user" class="profile-content card-glass">
      <div class="user-info">
        <div class="user-avatar">
          <img 
            :src="authStore.user.avatar || '/default-avatar.png'" 
            :alt="authStore.user.nickname"
            class="avatar-img"
          />
        </div>
        <div class="user-details">
          <h2>{{ authStore.user.nickname }}</h2>
          <p class="username">@{{ authStore.user.username }}</p>
          <p class="email">{{ authStore.user.email }}</p>
          <p class="role">角色: {{ authStore.user.role === 'admin' ? '管理员' : '普通用户' }}</p>
          <p class="created-at">注册时间: {{ new Date(authStore.user.createdAt).toLocaleString('zh-CN') }}</p>
        </div>
      </div>
      
      <div class="profile-actions">
        <button @click="authStore.logout()" class="logout-btn btn-primary">
          退出登录
        </button>
      </div>
    </div>
    
    <div v-else class="error-container">
      <p>无法加载用户信息</p>
      <button @click="authStore.getCurrentUser()" class="retry-btn btn-secondary">
        重试
      </button>
    </div>
  </div>
</template>

<style scoped>
.profile-container {
  min-height: 100vh;
  padding: var(--spacing-lg);
  background: linear-gradient(135deg, #f5f5f5 0%, #e8f5e8 100%);
}

.profile-header {
  text-align: center;
  margin-bottom: var(--spacing-2xl);
}

.profile-header h1 {
  font-size: var(--font-size-title2);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
}

.profile-content {
  max-width: 600px;
  margin: 0 auto;
  padding: var(--spacing-xl);
}

.user-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
  margin-bottom: var(--spacing-xl);
}

.user-avatar {
  flex-shrink: 0;
}

.avatar-img {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid var(--primary-color);
}

.user-details h2 {
  font-size: var(--font-size-title3);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin-bottom: var(--spacing-xs);
}

.username {
  color: var(--text-secondary);
  font-size: var(--font-size-body2);
}

.email {
  color: var(--text-tertiary);
  font-size: var(--font-size-caption);
}

.role {
  color: var(--text-secondary);
  font-size: var(--font-size-caption);
  font-weight: var(--font-weight-medium);
  margin-top: var(--spacing-xs);
}

.created-at {
  color: var(--text-quaternary);
  font-size: var(--font-size-caption2);
  margin-top: var(--spacing-xs);
}

.profile-actions {
  display: flex;
  justify-content: center;
}

.logout-btn {
  padding: var(--spacing-sm) var(--spacing-lg);
  background: var(--gradient-primary);
  color: white;
  border: none;
  border-radius: var(--border-radius-base);
  font-weight: var(--font-weight-medium);
  cursor: pointer;
  transition: all 0.2s ease;
}

.logout-btn:hover {
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

.loading-container,
.error-container {
  text-align: center;
  padding: var(--spacing-2xl);
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--text-quaternary);
  border-top: 3px solid var(--primary-color);
  border-radius: 50%;
  margin: 0 auto var(--spacing-md);
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.retry-btn {
  padding: var(--spacing-sm) var(--spacing-lg);
  background: var(--background-secondary);
  color: var(--text-primary);
  border: 1px solid var(--separator-color);
  border-radius: var(--border-radius-base);
  font-weight: var(--font-weight-medium);
  cursor: pointer;
  transition: all 0.2s ease;
}

.retry-btn:hover {
  background: var(--background-tertiary);
}

/* Responsive Design */
@media (max-width: 768px) {
  .profile-container {
    padding: var(--spacing-md);
  }
  
  .user-info {
    flex-direction: column;
    text-align: center;
  }
  
  .profile-content {
    padding: var(--spacing-lg);
  }
}
</style>