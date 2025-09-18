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
            <span class="btn-icon">🔄</span>
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
      <div class="error-icon">❌</div>
      <h3 class="error-title">加载失败</h3>
      <p class="error-message">{{ error }}</p>
      <button class="btn btn-primary" @click="refreshData">
        <span class="btn-icon">🔄</span>
        重试
      </button>
    </div>

    <!-- 空状态 -->
    <div v-else-if="evaluationRequests.length === 0" class="empty-container">
      <div class="empty-icon">📚</div>
      <h3 class="empty-title">暂无求评价课程</h3>
      <p class="empty-description">目前还没有用户发起求评价请求</p>
      <p class="empty-hint">浏览课程列表，为感兴趣的课程发起求评价吧！</p>
      <RouterLink to="/" class="btn btn-primary">
        <span class="btn-icon">🏠</span>
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
            <span class="btn-icon">‹</span>
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
            <span class="btn-icon">›</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 返回顶部按钮 -->
    <div class="back-to-top" @click="scrollToTop">
      <span class="btn-icon">⬆️</span>
    </div>
  </div>
</template>

<style scoped>
.evaluation-requests-container {
  min-height: 100vh;
  background: linear-gradient(135deg, var(--background-secondary) 0%, var(--natural-background) 100%);
  padding: var(--spacing-lg) 0;
}

/* 页面头部 */
.page-header {
  max-width: 1200px;
  margin: 0 auto var(--spacing-xl);
  padding: 0 var(--spacing-md);
}

.header-content {
  text-align: center;
}

.breadcrumb {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-lg);
  font-size: var(--font-size-body2);
}

.breadcrumb-link {
  color: var(--primary-color);
  text-decoration: none;
  font-weight: var(--font-weight-medium);
  transition: var(--transition-standard);
}

.breadcrumb-link:hover {
  color: var(--primary-color-dark);
  text-decoration: underline;
}

.breadcrumb-separator {
  color: var(--text-tertiary);
}

.breadcrumb-current {
  color: var(--text-secondary);
  font-weight: var(--font-weight-medium);
}

.header-title-section {
  margin-bottom: var(--spacing-xl);
}

.page-title {
  font-size: clamp(32px, 4vw, 40px);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin-bottom: var(--spacing-sm);
  letter-spacing: -0.5px;
  position: relative;
  display: inline-block;
}

.page-title::after {
  content: '';
  position: absolute;
  bottom: -8px;
  left: 50%;
  transform: translateX(-50%);
  width: 60px;
  height: 3px;
  background: var(--gradient-primary);
  border-radius: 2px;
}

.page-subtitle {
  font-size: clamp(16px, 1.5vw, 18px);
  color: var(--text-secondary);
  margin: var(--spacing-lg) 0 0;
  font-weight: var(--font-weight-regular);
  letter-spacing: 0.5px;
  opacity: 0.9;
}

.header-stats {
  display: flex;
  justify-content: center;
  gap: var(--spacing-2xl);
  flex-wrap: wrap;
}

.stat-item {
  text-align: center;
  min-width: 100px;
}

.stat-number {
  font-size: 28px;
  font-weight: var(--font-weight-bold);
  color: var(--primary-color);
  margin-bottom: var(--spacing-xs);
  display: block;
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.stat-label {
  font-size: var(--font-size-caption);
  color: var(--text-secondary);
  font-weight: var(--font-weight-medium);
}

/* 工具栏 */
.toolbar-section {
  max-width: 1200px;
  margin: 0 auto var(--spacing-lg);
  padding: 0 var(--spacing-md);
}

.toolbar-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-lg);
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
}

.toolbar-left {
  flex-grow: 1;
}

.section-title {
  font-size: 24px;
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin-bottom: var(--spacing-xs);
}

.section-description {
  font-size: var(--font-size-body2);
  color: var(--text-secondary);
  margin: 0;
}

/* 加载状态 */
.loading-container {
  max-width: 1200px;
  margin: var(--spacing-3xl) auto;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-md);
  padding: var(--spacing-3xl);
}

.loader {
  width: 40px;
  height: 40px;
  border: 4px solid var(--background-secondary);
  border-top: 4px solid var(--primary-color);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.loading-text {
  font-size: var(--font-size-body);
  color: var(--text-secondary);
  font-weight: var(--font-weight-medium);
}

/* 错误状态 */
.error-container {
  max-width: 600px;
  margin: var(--spacing-3xl) auto;
  text-align: center;
  padding: var(--spacing-3xl);
  background: rgba(255, 255, 255, 0.95);
  border-radius: 20px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
}

.error-icon {
  font-size: 48px;
  margin-bottom: var(--spacing-lg);
}

.error-title {
  font-size: 24px;
  font-weight: var(--font-weight-bold);
  color: var(--danger-color);
  margin-bottom: var(--spacing-sm);
}

.error-message {
  font-size: var(--font-size-body);
  color: var(--text-secondary);
  margin-bottom: var(--spacing-xl);
}

/* 空状态 */
.empty-container {
  max-width: 600px;
  margin: var(--spacing-3xl) auto;
  text-align: center;
  padding: var(--spacing-3xl);
  background: rgba(255, 255, 255, 0.95);
  border-radius: 20px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
}

.empty-icon {
  font-size: 48px;
  margin-bottom: var(--spacing-lg);
}

.empty-title {
  font-size: 24px;
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin-bottom: var(--spacing-sm);
}

.empty-description {
  font-size: var(--font-size-body);
  color: var(--text-secondary);
  margin-bottom: var(--spacing-sm);
}

.empty-hint {
  font-size: var(--font-size-body2);
  color: var(--text-tertiary);
  margin-bottom: var(--spacing-xl);
}

/* 课程网格 */
.course-grid-section {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 var(--spacing-md);
}

.grid-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
  flex-wrap: wrap;
  gap: var(--spacing-md);
}

.grid-title {
  font-size: 20px;
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin: 0;
}

.highlight {
  color: var(--primary-color);
  font-weight: var(--font-weight-bold);
}

.view-info {
  font-size: var(--font-size-body2);
  color: var(--text-secondary);
}

.course-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: var(--spacing-xl);
  margin-bottom: var(--spacing-3xl);
}

/* 分页 */
.pagination-container {
  margin-top: var(--spacing-3xl);
  padding-top: var(--spacing-xl);
  border-top: 1px solid var(--separator-color);
}

.pagination-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-lg);
}

.pagination-btn {
  padding: var(--spacing-sm) var(--spacing-md);
  border: 1px solid var(--separator-color);
  border-radius: 8px;
  background: white;
  color: var(--text-secondary);
  font-size: var(--font-size-body2);
  cursor: pointer;
  transition: var(--transition-standard);
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.pagination-btn:hover:not(:disabled) {
  background: var(--background-secondary);
  color: var(--text-primary);
  border-color: var(--primary-color);
}

.pagination-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.pagination-pages {
  display: flex;
  gap: var(--spacing-xs);
}

.pagination-page {
  width: 36px;
  height: 36px;
  border: 1px solid var(--separator-color);
  border-radius: 8px;
  background: white;
  color: var(--text-secondary);
  font-size: var(--font-size-body2);
  cursor: pointer;
  transition: var(--transition-standard);
}

.pagination-page:hover {
  background: var(--background-secondary);
}

.pagination-page.active {
  background: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

/* 返回顶部 */
.back-to-top {
  position: fixed;
  bottom: 30px;
  right: 30px;
  width: 50px;
  height: 50px;
  background: var(--primary-color);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: var(--shadow-lg);
  transition: var(--transition-standard);
  z-index: 100;
}

.back-to-top:hover {
  background: var(--primary-color-dark);
  transform: translateY(-2px);
  box-shadow: var(--shadow-xl);
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