<template>
  <div class="dashboard-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">数据仪表盘</h1>
      <p class="page-subtitle">查看系统整体运行情况和关键指标</p>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon size="32" color="var(--color-primary)">
                <Reading />
              </el-icon>
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
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon size="32" color="var(--color-primary)">
                <Star />
              </el-icon>
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
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon size="32" color="var(--color-primary)">
                <User />
              </el-icon>
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
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon size="32" color="var(--color-primary)">
                <ChatDotRound />
              </el-icon>
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
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表区域 -->
    <el-row :gutter="20" class="charts-row">
      <el-col :xs="24" :lg="12">
        <el-card class="chart-card">
          <template #header>
            <h3 class="chart-title">课程评分分布</h3>
          </template>
          <div class="chart-placeholder">
            <el-icon size="48" color="var(--color-text-secondary)">
              <TrendCharts />
            </el-icon>
            <p class="chart-text">评分分布图表将显示在这里</p>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :lg="12">
        <el-card class="chart-card">
          <template #header>
            <h3 class="chart-title">用户活跃度</h3>
          </template>
          <div class="chart-placeholder">
            <el-icon size="48" color="var(--color-text-secondary)">
              <TrendCharts />
            </el-icon>
            <p class="chart-text">用户活跃度图表将显示在这里</p>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 最近活动 -->
    <el-card class="activity-card">
      <template #header>
        <h3 class="section-title">最近活动</h3>
      </template>
      <el-timeline>
        <el-timeline-item
          v-for="(activity, index) in activities"
          :key="index"
          :icon="activity.icon"
          :color="activity.color"
        >
          <div class="activity-content">
            <p class="activity-text">{{ activity.text }}</p>
            <p class="activity-time">{{ activity.time }}</p>
          </div>
        </el-timeline-item>
      </el-timeline>
    </el-card>
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

const activities = ref([
  {
    icon: 'ChatDotRound',
    color: 'var(--color-primary)',
    text: '用户 张三 评论了课程 高等数学',
    time: '5分钟前'
  },
  {
    icon: 'Star',
    color: 'var(--color-primary)',
    text: '用户 李四 评分了课程 线性代数',
    time: '15分钟前'
  },
  {
    icon: 'User',
    color: 'var(--color-primary)',
    text: '新用户 王五 注册了账号',
    time: '30分钟前'
  }
]);

const fetchStats = async () => {
  try {
    const response = await axios.get(`${import.meta.env.VITE_BACKEND_BASE_URL}/admin/stats`);
    if (response.data.success) {
      stats.value = response.data.data;
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
.dashboard-container {
  padding: 20px;
  background-color: var(--color-background);
  min-height: 100vh;
}

.page-header {
  margin-bottom: 30px;
  text-align: center;
}

.page-title {
  font-size: 2rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0 0 10px 0;
}

.page-subtitle {
  font-size: 1rem;
  color: var(--color-text-secondary);
  margin: 0;
}

.stats-row {
  margin-bottom: 30px;
}

.stat-card {
  height: 100%;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 15px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  background: var(--color-background);
  border: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 1.8rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0 0 5px 0;
}

.stat-label {
  font-size: 0.9rem;
  color: var(--color-text-secondary);
  margin: 0;
}

.stat-footer {
  display: flex;
  align-items: center;
  gap: 10px;
}

.stat-trend {
  font-size: 0.8rem;
  font-weight: 500;
}

.stat-trend.positive {
  color: var(--color-success);
}

.stat-period {
  font-size: 0.8rem;
  color: var(--color-text-secondary);
}

.charts-row {
  margin-bottom: 30px;
}

.chart-card {
  height: 300px;
}

.chart-title {
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0;
}

.chart-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  background: var(--color-background);
  border: 1px solid var(--color-border);
  border-radius: 8px;
}

.chart-text {
  font-size: 0.9rem;
  color: var(--color-text-secondary);
  margin: 15px 0 0 0;
}

.activity-card {
  margin-bottom: 30px;
}

.section-title {
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0;
}

.activity-content {
  margin-left: 10px;
}

.activity-text {
  font-size: 0.9rem;
  color: var(--color-text);
  margin: 0 0 5px 0;
}

.activity-time {
  font-size: 0.8rem;
  color: var(--color-text-secondary);
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .dashboard-container {
    padding: 15px;
  }
  
  .page-title {
    font-size: 1.5rem;
  }
  
  .stat-content {
    flex-direction: column;
    text-align: center;
    gap: 10px;
  }
  
  .stat-icon {
    width: 40px;
    height: 40px;
  }
  
  .stat-value {
    font-size: 1.5rem;
  }
}
</style>