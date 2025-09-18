<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { OAuth2Service } from '@/services/oauth2Api'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

onMounted(async () => {
  const code = route.query.code as string
  const state = route.query.state as string
  
  if (!code || !state) {
    console.error('缺少OAuth2回调参数')
    router.push('/auth?error=missing_params')
    return
  }
  
  try {
    // 调用OAuth2服务处理回调
    const response = await OAuth2Service.handleOAuth2Callback(code, state)
    
    const { token, user } = response
    
    // 使用OAuth2登录
    await authStore.oauth2Login(token, user)
    
    // 跳转到首页
    router.push('/')
    
  } catch (error) {
    console.error('OAuth2登录失败:', error)
    router.push('/auth?error=oauth2_failed')
  }
})
</script>

<template>
  <div class="oauth2-callback-container">
    <div class="loading-card card-glass">
      <div class="loading-content">
        <div class="loading-spinner-large"></div>
        <h2>正在处理登录...</h2>
        <p>请稍候，我们正在验证您的身份</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.oauth2-callback-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f5f5f5 0%, #e8f5e8 100%);
}

.loading-card {
  padding: 48px 32px;
  text-align: center;
  max-width: 400px;
  width: 90%;
}

.loading-content h2 {
  margin: 24px 0 8px 0;
  color: var(--text-primary);
  font-size: 24px;
  font-weight: 600;
}

.loading-content p {
  color: var(--text-secondary);
  font-size: 16px;
  margin: 0;
}

.loading-spinner-large {
  width: 48px;
  height: 48px;
  border: 4px solid rgba(47, 169, 20, 0.2);
  border-top: 4px solid var(--primary-color);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
