<template>
  <div class="stats-view">
    <!-- 页面头部 -->
    <div class="stats-view-header">
      <div class="header-content">
        <div class="header-left">
          <h1 class="page-title">数据统计中心</h1>
          <p class="page-description">全面了解平台运行情况和用户行为数据</p>
        </div>
        <div class="header-right">
          <div class="refresh-controls">
            <span class="last-update">最后更新：{{ lastUpdateTime }}</span>
            <button class="refresh-btn" @click="refreshData">
              <i class="bi bi-arrow-clockwise"></i>
              刷新数据
            </button>
          </div>
          <div class="view-controls">
            <button 
              class="view-btn" 
              :class="{ active: currentView === 'home' }"
              @click="currentView = 'home'"
            >
              <i class="bi bi-house-door"></i>
              首页概览
            </button>
            <button 
              class="view-btn" 
              :class="{ active: currentView === 'detailed' }"
              @click="currentView === 'detailed'"
            >
              <i class="bi bi-bar-chart-line"></i>
              详细分析
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <main class="stats-content">
      <!-- 首页概览视图 -->
      <div v-if="currentView === 'home'" class="home-stats-section">
        <StatsDashboard />
      </div>

      <!-- 详细分析视图 -->
      <div v-if="currentView === 'detailed'" class="detailed-stats-section">
        <div class="detailed-stats-placeholder">
          <div class="placeholder-icon">
            <i class="bi bi-graph-up"></i>
          </div>
          <h3>详细分析功能开发中</h3>
          <p>更多详细的数据分析功能即将推出，包括：用户行为分析、课程趋势预测、数据导出等高级功能。</p>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import StatsDashboard from '@/components/StatsDashboard.vue'

// 响应式数据
const currentView = ref<'home' | 'detailed'>('home')
const lastUpdateTime = ref<string>('')

// 方法
const refreshData = async () => {
  // 这里可以添加刷新数据的逻辑
  // 例如重新加载子组件的数据
  lastUpdateTime.value = new Date().toLocaleString('zh-CN')
}

// 初始化
onMounted(() => {
  lastUpdateTime.value = new Date().toLocaleString('zh-CN')
})
</script>

<style scoped>
.stats-view {
  min-height: 100vh;
  background: linear-gradient(135deg, #f8fafb 0%, #e3f2fd 100%);
}

.stats-view-header {
  background: white;
  border-bottom: 1px solid #e9ecef;
  padding: 2rem 0;
  margin-bottom: 2rem;
}

.header-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 2rem;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 2rem;
}

.header-left {
  flex: 1;
}

.page-title {
  font-size: 2.5rem;
  font-weight: 700;
  color: #2c3e50;
  margin-bottom: 0.5rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.page-description {
  font-size: 1.1rem;
  color: #6c757d;
  margin: 0;
}

.header-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 1rem;
}

.refresh-controls {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.last-update {
  font-size: 0.9rem;
  color: #6c757d;
}

.refresh-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.refresh-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(40, 167, 69, 0.4);
}

.view-controls {
  display: flex;
  gap: 0.5rem;
}

.view-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: #f8f9fa;
  color: #6c757d;
  border: 1px solid #e9ecef;
  border-radius: 6px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.view-btn:hover {
  background: #e9ecef;
  color: #495057;
}

.view-btn.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-color: transparent;
}

.stats-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 2rem 2rem;
}

.home-stats-section {
  margin-bottom: 2rem;
}

.detailed-stats-section {
  background: white;
  border-radius: 16px;
  padding: 3rem;
  box-shadow: 0 8px 32px rgba(0,0,0,0.08);
  text-align: center;
}

.detailed-stats-placeholder {
  max-width: 600px;
  margin: 0 auto;
}

.placeholder-icon {
  font-size: 4rem;
  color: #667eea;
  margin-bottom: 1.5rem;
}

.placeholder-icon i {
  animation: pulse 2s infinite;
}

.detailed-stats-placeholder h3 {
  font-size: 1.5rem;
  color: #2c3e50;
  margin-bottom: 1rem;
}

.detailed-stats-placeholder p {
  font-size: 1rem;
  color: #6c757d;
  line-height: 1.6;
}

@keyframes pulse {
  0% { opacity: 1; }
  50% { opacity: 0.5; }
  100% { opacity: 1; }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .stats-view-header {
    padding: 1.5rem 0;
  }
  
  .header-content {
    flex-direction: column;
    text-align: center;
    padding: 0 1rem;
  }
  
  .header-right {
    align-items: center;
  }
  
  .refresh-controls {
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .stats-content {
    padding: 0 1rem 1rem;
  }
}
</style>