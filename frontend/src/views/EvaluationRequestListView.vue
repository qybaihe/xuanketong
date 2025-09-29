<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { evaluationRequestService } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import type { EvaluationRequest, EvaluationRequestListResponse } from '@/services/api'
import EvaluationRequestCard from '@/components/EvaluationRequestCard.vue'

const authStore = useAuthStore()
const evaluationRequests = ref<EvaluationRequest[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const currentPage = ref(1)
const pageSize = ref(12)
const totalItems = ref(0)

// 计算总页数
const totalPages = computed(() => Math.ceil(totalItems.value / pageSize.value))

// 获取求评价列表
const fetchEvaluationRequests = async () => {
  loading.value = true
  error.value = null
  
  try {
    const response: EvaluationRequestListResponse = await evaluationRequestService.getEvaluationRequests(
      currentPage.value,
      pageSize.value
    )
    
    evaluationRequests.value = response.items
    totalItems.value = response.total
  } catch (err: any) {
    console.error('获取求评价列表失败:', err)
    error.value = err.response?.data?.error || '获取求评价列表失败，请稍后重试'
  } finally {
    loading.value = false
  }
}

// 分页处理
const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchEvaluationRequests()
}

// 返回顶部函数
const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

// 刷新数据
const refreshData = () => {
  fetchEvaluationRequests()
}

onMounted(() => {
  fetchEvaluationRequests()
})

// 格式化日期
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}
</script>

<template>
  <div class="evaluation-requests-container">
    <!-- 页面头部 -->
    <header class="page-header">
      <div class="header-content">
        <div class="breadcrumb">
          <RouterLink to="/" class="breadcrumb-link">首页</RouterLink>
          <span class="breadcrumb-separator">›</span>
          <span class="breadcrumb-current">求评价中心</span>
        </div>
        
        <div class="header-title-section">
          <h1 class="page-title">求评价中心</h1>
          <p class="page-subtitle">
            发现需要帮助的课程，为同学们提供宝贵的评价和建议
          </p>
        </div>
        
        <div class="header-stats">
          <div class="stat-item">
            <div class="stat-number">{{ totalItems }}</div>
            <div class="stat-label">正在求评价</div>
          </div>
          <div class="stat-item">
            <div class="stat-number">{{ evaluationRequests.length }}</div>
            <div class="stat-label">当前显示</div>
          </div>
        </div>
      </div>
    </header>

    <!-- 操作栏 -->
    <section class="toolbar-section">
      <div class="toolbar-content">
        <div class="toolbar-left">
          <h2 class="section-title">求评价课程列表</h2>
          <p class="section-description">点击下方课程卡片查看详情并进行评价</p>
        </div>
        
        <div class="toolbar-right">
          <button class="btn btn-secondary" @click="refreshData" :disabled="loading">
            <svg class="btn-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8" stroke="#1A1A1A" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M21 3v5h-5" stroke="#1A1A1A" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16" stroke="#1A1A1A" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M3 21v-5h5" stroke="#1A1A1A" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            {{ loading ? '刷新中...' : '刷新' }}
          </button>
        </div>
      </div>
    </section>

    <!-- 加载状态 -->
    <div v-if="loading && evaluationRequests.length === 0" class="loading-container">
      <div class="loader"></div>
      <p class="loading-text">正在加载求评价列表...</p>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="error-container">
      <svg class="error-icon" width="48" height="48" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <circle cx="12" cy="12" r="10" stroke="#FF6B6B" stroke-width="2"/>
        <line x1="15" y1="9" x2="9" y2="15" stroke="#FF6B6B" stroke-width="2" stroke-linecap="round"/>
        <line x1="9" y1="9" x2="15" y2="15" stroke="#FF6B6B" stroke-width="2" stroke-linecap="round"/>
      </svg>
      <h3 class="error-title">加载失败</h3>
      <p class="error-message">{{ error }}</p>
      <button class="btn btn-primary" @click="refreshData">
        <svg class="btn-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8" stroke="#1A1A1A" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M21 3v5h-5" stroke="#1A1A1A" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16" stroke="#1A1A1A" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M3 21v-5h5" stroke="#1A1A1A" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        重试
      </button>
    </div>

    <!-- 空状态 -->
    <div v-else-if="evaluationRequests.length === 0" class="empty-container">
      <svg class="empty-icon" width="48" height="48" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20" stroke="#76D7C4" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z" stroke="#76D7C4" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      <h3 class="empty-title">暂无求评价课程</h3>
      <p class="empty-description">目前还没有用户发起求评价请求</p>
      <p class="empty-hint">浏览课程列表，为感兴趣的课程发起求评价吧！</p>
      <RouterLink to="/" class="btn btn-primary">
        <svg class="btn-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" stroke="#1A1A1A" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          <polyline points="9,22 9,12 15,12 15,22" stroke="#1A1A1A" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        返回首页
      </RouterLink>
    </div>

    <!-- 课程列表 -->
    <div v-else class="course-grid-section">
      <!-- 列表头部 -->
      <div class="grid-header">
        <h3 class="grid-title">
          共找到 <span class="highlight">{{ totalItems }}</span> 个求评价请求
        </h3>
        <div class="view-info">
          第 {{ currentPage }} 页，共 {{ totalPages }} 页
        </div>
      </div>

      <!-- 课程卡片网格 -->
      <div class="course-grid">
        <EvaluationRequestCard
          v-for="request in evaluationRequests"
          :key="request.id"
          :request="request"
        />
      </div>

      <!-- 分页控件 -->
      <div v-if="totalPages > 1" class="pagination-container">
        <div class="pagination-content">
          <button
            class="pagination-btn"
            :disabled="currentPage === 1"
            @click="handlePageChange(currentPage - 1)"
          >
            <svg class="btn-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <polyline points="15,18 9,12 15,6" stroke="#1A1A1A" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            上一页
          </button>

          <div class="pagination-pages">
            <button
              v-for="page in Math.min(5, totalPages)"
              :key="page"
              class="pagination-page"
              :class="{ active: page === currentPage }"
              @click="handlePageChange(page)"
            >
              {{ page }}
            </button>
          </div>

          <button
            class="pagination-btn"
            :disabled="currentPage === totalPages"
            @click="handlePageChange(currentPage + 1)"
          >
            下一页
            <svg class="btn-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <polyline points="9,18 15,12 9,6" stroke="#1A1A1A" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- 返回顶部按钮 -->
    <div class="back-to-top" @click="scrollToTop">
      <svg class="btn-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <polyline points="18,15 12,9 6,15" stroke="#1A1A1A" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
    </div>
  </div>
</template>

<style scoped>
.evaluation-requests-container {
  min-height: 100vh;
  background-color: #FEF6F7;
  padding: 20px;
  font-family: sans-serif;
  color: #1A1A1A;
}

/* 页面头部 */
.page-header {
  max-width: 1200px;
  margin: 0 auto 32px;
  padding: 0 16px;
}

.header-content {
  text-align: center;
}

.breadcrumb {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-bottom: 24px;
  font-size: 14px;
}

.breadcrumb-link {
  color: #76D7C4;
  text-decoration: none;
  font-weight: bold;
  transition: color 0.2s ease;
}

.breadcrumb-link:hover {
  color: #1A1A1A;
  text-decoration: underline;
}

.breadcrumb-separator {
  color: #666666;
}

.breadcrumb-current {
  color: #1A1A1A;
  font-weight: bold;
}

.header-title-section {
  margin-bottom: 32px;
}

.page-title {
  font-size: 32px;
  font-weight: bold;
  color: #1A1A1A;
  margin-bottom: 8px;
  text-align: center;
}

.page-subtitle {
  font-size: 16px;
  color: #1A1A1A;
  margin: 16px 0 0;
  font-weight: 500;
}

.header-stats {
  display: flex;
  justify-content: center;
  gap: 32px;
  flex-wrap: wrap;
}

.stat-item {
  text-align: center;
  min-width: 100px;
  background-color: #FFFFFF;
  padding: 16px;
  border-radius: 8px;
  border: 2px solid #000000;
  box-shadow: 3px 3px 0px 0px #000000;
}

.stat-number {
  font-size: 24px;
  font-weight: bold;
  color: #76D7C4;
  margin-bottom: 4px;
  display: block;
}

.stat-label {
  font-size: 12px;
  color: #1A1A1A;
  font-weight: bold;
}

/* 按钮基础样式 */
.btn {
  background-color: #F7D074;
  border-radius: 8px;
  border: 3px solid #000000;
  box-shadow: 4px 4px 0px 0px #000000;
  padding: 12px 20px;
  font-size: 14px;
  font-weight: bold;
  color: #1A1A1A;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  text-decoration: none;
  transition: transform 0.2s ease;
}

.btn:hover {
  transform: translate(-2px, -2px);
  box-shadow: 6px 6px 0px 0px #000000;
}

.btn:active {
  transform: translate(-1px, -1px);
  box-shadow: 3px 3px 0px 0px #000000;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
  box-shadow: 4px 4px 0px 0px #000000;
}

.btn-secondary {
  background-color: #76D7C4;
}

/* 工具栏 */
.toolbar-section {
  max-width: 1200px;
  margin: 0 auto 24px;
  padding: 0 16px;
}

.toolbar-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  background-color: #FFFFFF;
  border-radius: 12px;
  border: 3px solid #000000;
  box-shadow: 5px 5px 0px 0px #000000;
}

.toolbar-left {
  flex-grow: 1;
}

.section-title {
  font-size: 20px;
  font-weight: bold;
  color: #1A1A1A;
  margin-bottom: 4px;
}

.section-description {
  font-size: 14px;
  color: #1A1A1A;
  margin: 0;
}

/* 加载状态 */
.loading-container {
  max-width: 1200px;
  margin: 60px auto;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 60px;
}

.loader {
  width: 40px;
  height: 40px;
  border: 4px solid #F7D074;
  border-top: 4px solid #76D7C4;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.loading-text {
  font-size: 16px;
  color: #1A1A1A;
  font-weight: bold;
}

/* 错误状态 */
.error-container {
  max-width: 600px;
  margin: 60px auto;
  text-align: center;
  padding: 40px;
  background-color: #FFFFFF;
  border-radius: 12px;
  border: 3px solid #000000;
  box-shadow: 5px 5px 0px 0px #000000;
}

.error-icon {
  margin-bottom: 24px;
}

.error-title {
  font-size: 20px;
  font-weight: bold;
  color: #FF6B6B;
  margin-bottom: 8px;
}

.error-message {
  font-size: 16px;
  color: #1A1A1A;
  margin-bottom: 32px;
}

/* 空状态 */
.empty-container {
  max-width: 600px;
  margin: 60px auto;
  text-align: center;
  padding: 40px;
  background-color: #FFFFFF;
  border-radius: 12px;
  border: 3px solid #000000;
  box-shadow: 5px 5px 0px 0px #000000;
}

.empty-icon {
  margin-bottom: 24px;
}

.empty-title {
  font-size: 20px;
  font-weight: bold;
  color: #1A1A1A;
  margin-bottom: 8px;
}

.empty-description {
  font-size: 16px;
  color: #1A1A1A;
  margin-bottom: 8px;
}

.empty-hint {
  font-size: 14px;
  color: #666666;
  margin-bottom: 32px;
}

/* 课程网格 */
.course-grid-section {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 16px;
}

.grid-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.grid-title {
  font-size: 18px;
  font-weight: bold;
  color: #1A1A1A;
  margin: 0;
}

.highlight {
  color: #76D7C4;
  font-weight: bold;
}

.view-info {
  font-size: 14px;
  color: #1A1A1A;
  font-weight: bold;
}

.course-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 24px;
  margin-bottom: 60px;
}

/* 分页 */
.pagination-container {
  margin-top: 60px;
  padding-top: 24px;
  border-top: 2px solid #000000;
}

.pagination-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 24px;
}

.pagination-btn {
  padding: 8px 16px;
  border: 2px solid #000000;
  border-radius: 8px;
  background-color: #FFFFFF;
  color: #1A1A1A;
  font-size: 14px;
  font-weight: bold;
  cursor: pointer;
  transition: transform 0.2s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  box-shadow: 3px 3px 0px 0px #000000;
}

.pagination-btn:hover:not(:disabled) {
  transform: translate(-2px, -2px);
  box-shadow: 5px 5px 0px 0px #000000;
}

.pagination-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
  box-shadow: 3px 3px 0px 0px #000000;
}

.pagination-pages {
  display: flex;
  gap: 8px;
}

.pagination-page {
  width: 36px;
  height: 36px;
  border: 2px solid #000000;
  border-radius: 8px;
  background-color: #FFFFFF;
  color: #1A1A1A;
  font-size: 14px;
  font-weight: bold;
  cursor: pointer;
  transition: transform 0.2s ease;
  box-shadow: 2px 2px 0px 0px #000000;
}

.pagination-page:hover {
  transform: translate(-1px, -1px);
  box-shadow: 3px 3px 0px 0px #000000;
}

.pagination-page.active {
  background-color: #76D7C4;
  color: #1A1A1A;
  border-color: #000000;
}

/* 返回顶部 */
.back-to-top {
  position: fixed;
  bottom: 30px;
  right: 30px;
  width: 50px;
  height: 50px;
  background-color: #F7D074;
  color: #1A1A1A;
  border-radius: 50%;
  border: 3px solid #000000;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 4px 4px 0px 0px #000000;
  transition: transform 0.2s ease;
  z-index: 100;
}

.back-to-top:hover {
  transform: translate(-2px, -2px);
  box-shadow: 6px 6px 0px 0px #000000;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .toolbar-content {
    flex-direction: column;
    text-align: center;
  }
  
  .course-grid {
    grid-template-columns: 1fr;
  }
  
  .grid-header {
    flex-direction: column;
    text-align: center;
  }
  
  .header-stats {
    gap: var(--spacing-lg);
  }
  
  .pagination-content {
    flex-direction: column;
    gap: var(--spacing-md);
  }
}

/* 暗色主题支持 */
@media (prefers-color-scheme: dark) {
  .toolbar-content {
    background: rgba(0, 0, 0, 0.6);
    border-color: rgba(255, 255, 255, 0.1);
  }
  
  .evaluation-requests-container {
    background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  }
  
  .page-title {
    background: linear-gradient(135deg, #4fc830 0%, #2fa914 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }
}
</style>