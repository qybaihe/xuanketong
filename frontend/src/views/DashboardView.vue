<template>
  <div class="dashboard-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title text-shine">数据仪表盘</h1>
      <p class="page-subtitle">查看系统整体运行情况和关键指标</p>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card card-glass">
        <div class="stat-content">
          <div class="stat-icon courses-icon">
            <span class="icon-emoji">📚</span>
          </div>
          <div class="stat-info">
            <h3 class="stat-value">{{ stats.total_courses }}</h3>
            <p class="stat-label">总课程数</p>
          </div>
        </div>
        <div class="stat-footer">
          <span class="stat-trend positive">↗ 12%</span>
          <span class="stat-period">较上月</span>
        </div>
      </div>

      <div class="stat-card card-glass">
        <div class="stat-content">
          <div class="stat-icon rating-icon">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="#F7D074">
              <path d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z"/>
            </svg>
          </div>
          <div class="stat-info">
            <h3 class="stat-value">{{ stats.average_rating.toFixed(1) }}</h3>
            <p class="stat-label">平均评分</p>
          </div>
        </div>
        <div class="stat-footer">
          <span class="stat-trend positive">↗ 0.3</span>
          <span class="stat-period">较上月</span>
        </div>
      </div>

      <div class="stat-card card-glass">
        <div class="stat-content">
          <div class="stat-icon users-icon">
            <span class="icon-emoji">👥</span>
          </div>
          <div class="stat-info">
            <h3 class="stat-value">{{ stats.total_users }}</h3>
            <p class="stat-label">注册用户</p>
          </div>
        </div>
        <div class="stat-footer">
          <span class="stat-trend positive">↗ 8%</span>
          <span class="stat-period">较上月</span>
        </div>
      </div>

      <div class="stat-card card-glass">
        <div class="stat-content">
          <div class="stat-icon comments-icon">
            <span class="icon-emoji">💬</span>
          </div>
          <div class="stat-info">
            <h3 class="stat-value">{{ stats.total_comments }}</h3>
            <p class="stat-label">总评论数</p>
          </div>
        </div>
        <div class="stat-footer">
          <span class="stat-trend positive">↗ 15%</span>
          <span class="stat-period">较上月</span>
        </div>
      </div>
    </div>

    <!-- 图表区域 -->
    <div class="charts-section">
      <div class="chart-container card-glass">
        <h2 class="chart-title">课程评分分布</h2>
        <div class="chart-placeholder">
          <div class="chart-icon">📊</div>
          <p class="chart-text">评分分布图表将显示在这里</p>
        </div>
      </div>

      <div class="chart-container card-glass">
        <h2 class="chart-title">用户活跃度</h2>
        <div class="chart-placeholder">
          <div class="chart-icon">📈</div>
          <p class="chart-text">用户活跃度图表将显示在这里</p>
        </div>
      </div>
    </div>

    <!-- 最近活动 -->
    <div class="recent-activity card-glass">
      <h2 class="section-title">最近活动</h2>
      <div class="activity-list">
        <div class="activity-item">
          <div class="activity-icon activity-comment">
            <span>💬</span>
          </div>
          <div class="activity-content">
            <p class="activity-text">用户 <span class="activity-user">张三</span> 评论了课程 <span class="activity-course">高等数学</span></p>
            <p class="activity-time">5分钟前</p>
          </div>
        </div>
        
        <div class="activity-item">
          <div class="activity-icon activity-rating">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="#F7D074">
              <path d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z"/>
            </svg>
          </div>
          <div class="activity-content">
            <p class="activity-text">用户 <span class="activity-user">李四</span> 评分了课程 <span class="activity-course">线性代数</span></p>
            <p class="activity-time">15分钟前</p>
          </div>
        </div>
        
        <div class="activity-item">
          <div class="activity-icon activity-user">
            <span>👤</span>
          </div>
          <div class="activity-content">
            <p class="activity-text">新用户 <span class="activity-user">王五</span> 注册了账号</p>
            <p class="activity-time">30分钟前</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import axios from 'axios';

const stats = ref({
  total_users: 0,
  total_courses: 0,
  total_comments: 0,
  average_rating: 0,
});

const fetchStats = async () => {
  try {
    const response = await axios.get('/api/v1/admin/stats');
    if (response.data) {
      stats.value = response.data;
    }
  } catch (error) {
    console.error('获取统计数据失败:', error);
  }
};

onMounted(() => {
  fetchStats();
});
</script>

<style scoped>
/* 仪表盘容器 */
.dashboard-container {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xl);
}

/* 页面标题 */
.page-header {
  text-align: center;
  margin-bottom: var(--spacing-lg);
}

.page-title {
  font-size: var(--font-size-title1);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin: 0 0 var(--spacing-sm) 0;
}

.page-subtitle {
  font-size: var(--font-size-body);
  color: var(--text-secondary);
  margin: 0;
}

/* 统计卡片网格 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: var(--spacing-lg);
}

.stat-card {
  padding: var(--spacing-lg);
  border-radius: 16px;
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-xl);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-md);
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.courses-icon {
  background: linear-gradient(135deg, #4fc830 0%, #2fa914 100%);
}

.rating-icon {
  background: linear-gradient(135deg, #ffd688 0%, #dd9f54 100%);
}

.users-icon {
  background: linear-gradient(135deg, #5ac8fa 0%, #007aff 100%);
}

.comments-icon {
  background: linear-gradient(135deg, #ff9500 0%, #ff3b30 100%);
}

.icon-emoji {
  font-size: 28px;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: var(--font-size-title2);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin: 0 0 var(--spacing-xs) 0;
}

.stat-label {
  font-size: var(--font-size-body);
  color: var(--text-secondary);
  margin: 0;
}

.stat-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: var(--spacing-md);
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.stat-trend {
  font-size: var(--font-size-caption);
  font-weight: var(--font-weight-semibold);
  padding: 2px 6px;
  border-radius: 4px;
}

.stat-trend.positive {
  color: var(--success-base);
  background: rgba(52, 199, 89, 0.1);
}

.stat-trend.negative {
  color: var(--error-base);
  background: rgba(255, 59, 48, 0.1);
}

.stat-period {
  font-size: var(--font-size-caption);
  color: var(--text-tertiary);
}

/* 图表区域 */
.charts-section {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: var(--spacing-lg);
}

.chart-container {
  padding: var(--spacing-lg);
  border-radius: 16px;
  transition: all 0.3s ease;
}

.chart-container:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-xl);
}

.chart-title {
  font-size: var(--font-size-title3);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin: 0 0 var(--spacing-lg) 0;
}

.chart-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
}

.chart-icon {
  font-size: 48px;
  margin-bottom: var(--spacing-md);
}

.chart-text {
  font-size: var(--font-size-body);
  color: var(--text-tertiary);
  margin: 0;
}

/* 最近活动 */
.recent-activity {
  padding: var(--spacing-lg);
  border-radius: 16px;
  transition: all 0.3s ease;
}

.recent-activity:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-xl);
}

.section-title {
  font-size: var(--font-size-title3);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin: 0 0 var(--spacing-lg) 0;
}

.activity-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.activity-item {
  display: flex;
  align-items: flex-start;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  transition: all 0.2s ease;
}

.activity-item:hover {
  background: rgba(255, 255, 255, 0.1);
}

.activity-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.activity-comment {
  background: linear-gradient(135deg, #ff9500 0%, #ff3b30 100%);
}

.activity-rating {
  background: linear-gradient(135deg, #ffd688 0%, #dd9f54 100%);
}

.activity-user {
  background: linear-gradient(135deg, #5ac8fa 0%, #007aff 100%);
}

.activity-content {
  flex: 1;
}

.activity-text {
  font-size: var(--font-size-body);
  color: var(--text-primary);
  margin: 0 0 var(--spacing-xs) 0;
}

.activity-user {
  font-weight: var(--font-weight-semibold);
  color: var(--primary-color);
}

.activity-course {
  font-weight: var(--font-weight-semibold);
  color: var(--info-color);
}

.activity-time {
  font-size: var(--font-size-caption);
  color: var(--text-tertiary);
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .charts-section {
    grid-template-columns: 1fr;
  }
  
  .stat-content {
    flex-direction: column;
    text-align: center;
  }
  
  .activity-item {
    flex-direction: column;
    text-align: center;
  }
}

@media (min-width: 769px) and (max-width: 1024px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>