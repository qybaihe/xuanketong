<template>
  <div class="admin-dashboard">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">管理后台仪表盘</h1>
      <p class="page-subtitle">实时监控系统运行状态和关键指标</p>
    </div>

    <!-- 统计概览卡片 -->
    <div class="stats-overview">
      <el-row :gutter="20">
        <el-col :xs="24" :sm="12" :md="8" :lg="6">
          <el-card class="stat-card primary">
            <div class="stat-content">
              <div class="stat-icon">
                <el-icon size="40"><Reading /></el-icon>
              </div>
              <div class="stat-info">
                <h3 class="stat-number">{{ stats.total_courses }}</h3>
                <p class="stat-label">总课程数</p>
                <div class="stat-trend">
                  <el-icon><ArrowUp /></el-icon>
                  <span>+12%</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :xs="24" :sm="12" :md="8" :lg="6">
          <el-card class="stat-card success">
            <div class="stat-content">
              <div class="stat-icon">
                <el-icon size="40"><User /></el-icon>
              </div>
              <div class="stat-info">
                <h3 class="stat-number">{{ stats.total_users }}</h3>
                <p class="stat-label">注册用户</p>
                <div class="stat-trend">
                  <el-icon><ArrowUp /></el-icon>
                  <span>+8%</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :xs="24" :sm="12" :md="8" :lg="6">
          <el-card class="stat-card warning">
            <div class="stat-content">
              <div class="stat-icon">
                <el-icon size="40"><Star /></el-icon>
              </div>
              <div class="stat-info">
                <h3 class="stat-number">{{ stats.average_rating.toFixed(1) }}</h3>
                <p class="stat-label">平均评分</p>
                <div class="stat-trend">
                  <el-icon><ArrowUp /></el-icon>
                  <span>+0.3</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :xs="24" :sm="12" :md="8" :lg="6">
          <el-card class="stat-card info">
            <div class="stat-content">
              <div class="stat-icon">
                <el-icon size="40"><ChatDotRound /></el-icon>
              </div>
              <div class="stat-info">
                <h3 class="stat-number">{{ stats.total_comments }}</h3>
                <p class="stat-label">总评论数</p>
                <div class="stat-trend">
                  <el-icon><ArrowUp /></el-icon>
                  <span>+15%</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 详细统计区域 -->
    <div class="detailed-stats">
      <el-row :gutter="20">
        <!-- 用户活动统计 -->
        <el-col :xs="24" :lg="12">
          <el-card class="chart-card">
            <template #header>
              <div class="card-header">
                <h3 class="card-title">用户活动统计</h3>
                <el-button type="primary" size="small" @click="refreshStats">
                  <el-icon><Refresh /></el-icon>
                  刷新
                </el-button>
              </div>
            </template>
            <div class="activity-stats">
              <div class="activity-item">
                <div class="activity-label">本周新增用户</div>
                <div class="activity-value">{{ stats.new_users_this_week }}</div>
              </div>
              <div class="activity-item">
                <div class="activity-label">总评分数</div>
                <div class="activity-value">{{ stats.total_ratings }}</div>
              </div>
              <div class="activity-item">
                <div class="activity-label">待审核评论</div>
                <div class="activity-value">{{ stats.pending_comments }}</div>
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- 系统状态 -->
        <el-col :xs="24" :lg="12">
          <el-card class="chart-card">
            <template #header>
              <h3 class="card-title">系统状态</h3>
            </template>
            <div class="system-status">
              <div class="status-item">
                <div class="status-indicator online"></div>
                <div class="status-info">
                  <div class="status-label">数据库连接</div>
                  <div class="status-value">正常</div>
                </div>
              </div>
              <div class="status-item">
                <div class="status-indicator online"></div>
                <div class="status-info">
                  <div class="status-label">API服务</div>
                  <div class="status-value">正常</div>
                </div>
              </div>
              <div class="status-item">
                <div class="status-indicator online"></div>
                <div class="status-info">
                  <div class="status-label">文件存储</div>
                  <div class="status-value">正常</div>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 快速操作 -->
    <div class="quick-actions">
      <el-card>
        <template #header>
          <h3 class="card-title">快速操作</h3>
        </template>
        <div class="action-buttons">
          <el-button type="primary" @click="goToCourseManagement">
            <el-icon><Reading /></el-icon>
            课程管理
          </el-button>
          <el-button type="success" @click="goToUserManagement">
            <el-icon><User /></el-icon>
            用户管理
          </el-button>
          <el-button type="warning" @click="goToStats">
            <el-icon><TrendCharts /></el-icon>
            数据统计
          </el-button>
          <el-button type="info" @click="goToSystemSettings">
            <el-icon><Setting /></el-icon>
            系统设置
          </el-button>
        </div>
      </el-card>
    </div>

    <!-- 最近活动 -->
    <div class="recent-activities">
      <el-card>
        <template #header>
          <h3 class="card-title">最近活动</h3>
        </template>
        <el-timeline>
          <el-timeline-item
            v-for="(activity, index) in activities"
            :key="index"
            :icon="activity.icon"
            :color="activity.color"
            :timestamp="activity.time"
          >
            <div class="activity-content">
              <p class="activity-text">{{ activity.text }}</p>
            </div>
          </el-timeline-item>
        </el-timeline>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus'

const router = useRouter()

const stats = ref({
  total_users: 0,
  total_courses: 0,
  total_comments: 0,
  total_ratings: 0,
  average_rating: 0,
  new_users_this_week: 0,
  pending_comments: 0,
})

const activities = ref([
  {
    icon: 'ChatDotRound',
    color: '#409EFF',
    text: '用户 张三 评论了课程 高等数学',
    time: '5分钟前'
  },
  {
    icon: 'Star',
    color: '#E6A23C',
    text: '用户 李四 评分了课程 线性代数',
    time: '15分钟前'
  },
  {
    icon: 'User',
    color: '#67C23A',
    text: '新用户 王五 注册了账号',
    time: '30分钟前'
  },
  {
    icon: 'Reading',
    color: '#F56C6C',
    text: '管理员添加了新课程 数据结构',
    time: '1小时前'
  }
])

const fetchStats = async () => {
  try {
    const response = await axios.get(`${import.meta.env.VITE_BACKEND_BASE_URL}/admin/stats`)
    if (response.data.success) {
      stats.value = response.data.data
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
    ElMessage.error('获取统计数据失败')
  }
}

const refreshStats = async () => {
  await fetchStats()
  ElMessage.success('数据刷新成功')
}

const goToCourseManagement = () => {
  router.push('/admin/courses')
}

const goToUserManagement = () => {
  router.push('/admin/users')
}

const goToStats = () => {
  router.push('/admin/stats')
}

const goToSystemSettings = () => {
  ElMessage.info('系统设置功能开发中...')
}

onMounted(() => {
  fetchStats()
})
</script>

<style scoped>
.admin-dashboard {
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

.stats-overview {
  margin-bottom: 30px;
}

.stat-card {
  height: 100%;
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 20px 0 rgba(0, 0, 0, 0.15);
}

.stat-card.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.stat-card.success {
  background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
  color: white;
}

.stat-card.warning {
  background: linear-gradient(135deg, #ffc107 0%, #fd7e14 100%);
  color: white;
}

.stat-card.info {
  background: linear-gradient(135deg, #17a2b8 0%, #138496 100%);
  color: white;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.stat-icon {
  flex-shrink: 0;
}

.stat-info {
  flex: 1;
}

.stat-number {
  font-size: 2.5rem;
  font-weight: 700;
  margin: 0 0 5px 0;
  line-height: 1;
}

.stat-label {
  font-size: 0.9rem;
  opacity: 0.9;
  margin: 0 0 10px 0;
}

.stat-trend {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 0.8rem;
  opacity: 0.8;
}

.detailed-stats {
  margin-bottom: 30px;
}

.chart-card {
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0;
}

.activity-stats {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.activity-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  background: var(--color-background);
  border-radius: 8px;
  border: 1px solid var(--color-border);
}

.activity-label {
  font-size: 0.9rem;
  color: var(--color-text-secondary);
}

.activity-value {
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--color-text);
}

.system-status {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px;
  background: var(--color-background);
  border-radius: 8px;
  border: 1px solid var(--color-border);
}

.status-indicator {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
}

.status-indicator.online {
  background: #67C23A;
}

.status-indicator.offline {
  background: #F56C6C;
}

.status-info {
  flex: 1;
}

.status-label {
  font-size: 0.9rem;
  color: var(--color-text-secondary);
  margin-bottom: 5px;
}

.status-value {
  font-size: 1rem;
  font-weight: 500;
  color: var(--color-text);
}

.quick-actions {
  margin-bottom: 30px;
}

.action-buttons {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
}

.action-buttons .el-button {
  flex: 1;
  min-width: 150px;
}

.recent-activities {
  margin-bottom: 30px;
}

.activity-content {
  margin-left: 10px;
}

.activity-text {
  font-size: 0.9rem;
  color: var(--color-text);
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .admin-dashboard {
    padding: 15px;
  }
  
  .page-title {
    font-size: 1.5rem;
  }
  
  .stat-content {
    flex-direction: column;
    text-align: center;
    gap: 15px;
  }
  
  .stat-number {
    font-size: 2rem;
  }
  
  .action-buttons {
    flex-direction: column;
  }
  
  .action-buttons .el-button {
    width: 100%;
  }
}
</style>
