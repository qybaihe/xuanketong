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
  background-color: #FEF6F7;
  font-family: sans-serif;
  color: #1A1A1A;
}

.profile-header {
  text-align: center;
  margin-bottom: 20px;
}

.profile-header h1 {
  font-size: 32px;
  font-weight: bold;
  color: #1A1A1A;
}

.profile-content {
  max-width: 600px;
  margin: 0 auto;
  padding: 24px;
  background-color: #FFFFFF;
  border-radius: 12px;
  border: 3px solid #000000;
  box-shadow: 5px 5px 0px 0px #000000;
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
  border: 3px solid #000000;
  box-shadow: 3px 3px 0px 0px #000000;
}

.avatar-placeholder {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background-color: #F7D074;
  color: #1A1A1A;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 32px;
  border: 3px solid #000000;
  box-shadow: 3px 3px 0px 0px #000000;
}

.user-details h2 {
  font-size: 24px;
  font-weight: bold;
  color: #1A1A1A;
  margin-bottom: 6px;
}

.username {
  color: #1A1A1A;
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 4px;
}

.email {
  color: #1A1A1A;
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 4px;
}

.role {
  color: #1A1A1A;
  font-size: 14px;
  font-weight: bold;
  margin-top: 6px;
}

.created-at {
  color: #666666;
  font-size: 12px;
  margin-top: 6px;
}

.profile-actions {
  display: flex;
  justify-content: center;
}

.logout-btn {
  padding: 12px 24px;
  background-color: #FF6B6B;
  color: #1A1A1A;
  border: 3px solid #000000;
  border-radius: 8px;
  font-weight: bold;
  font-size: 14px;
  cursor: pointer;
  transition: transform 0.2s ease;
  box-shadow: 4px 4px 0px 0px #000000;
}

.logout-btn:hover {
  transform: translate(-2px, -2px);
  box-shadow: 6px 6px 0px 0px #000000;
}

.loading-container,
.error-container {
  text-align: center;
  padding: 40px;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #F7D074;
  border-top: 4px solid #76D7C4;
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
  background-color: #F7D074;
  color: #1A1A1A;
  border: 3px solid #000000;
  border-radius: 8px;
  font-weight: bold;
  font-size: 14px;
  cursor: pointer;
  transition: transform 0.2s ease;
  box-shadow: 4px 4px 0px 0px #000000;
}

.retry-btn:hover {
  transform: translate(-2px, -2px);
  box-shadow: 6px 6px 0px 0px #000000;
}

/* Responsive Design */
@media (max-width: 768px) {
  .profile-container {
    padding: 16px;
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