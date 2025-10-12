<template>
  <div class="stats-view">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">数据统计</h1>
      <p class="page-subtitle">查看系统详细的数据统计和分析</p>
    </div>

    <!-- 统计概览 -->
    <el-row :gutter="20" class="stats-overview">
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
                <User />
              </el-icon>
            </div>
            <div class="stat-info">
              <h3 class="stat-value">{{ stats.total_users }}</h3>
              <p class="stat-label">总用户数</p>
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
    </el-row>

    <!-- 操作栏 -->
    <el-card class="control-card">
      <el-row :gutter="20" align="center">
        <el-col :xs="24" :sm="8" :md="6">
          <el-form-item label="时间范围">
            <el-select
              v-model="timeRange"
              placeholder="选择时间范围"
              @change="handleTimeRangeChange"
            >
              <el-option label="最近7天" value="7d" />
              <el-option label="最近30天" value="30d" />
              <el-option label="最近90天" value="90d" />
              <el-option label="最近一年" value="1y" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="8" :md="6">
          <el-form-item label="数据类型">
            <el-select
              v-model="dataType"
              placeholder="选择数据类型"
              @change="handleDataTypeChange"
            >
              <el-option label="用户活跃度" value="users" />
              <el-option label="课程访问量" value="courses" />
              <el-option label="评论数量" value="comments" />
              <el-option label="评分分布" value="ratings" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="8" :md="6">
          <el-button type="primary" @click="refreshData" :loading="refreshing">
            <el-icon><Refresh /></el-icon>
            刷新数据
          </el-button>
        </el-col>
        <el-col :xs="24" :sm="24" :md="6">
          <el-button type="success" @click="exportData">
            <el-icon><Download /></el-icon>
            导出报告
          </el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 图表区域 -->
    <el-row :gutter="20" class="charts-section">
      <el-col :xs="24" :lg="12">
        <el-card class="chart-card">
          <template #header>
            <div class="chart-header">
              <h3 class="chart-title">用户增长趋势</h3>
              <el-tag type="info" size="small">{{ timeRangeText }}</el-tag>
            </div>
          </template>
          <div class="chart-container">
            <div class="chart-placeholder">
              <el-icon size="48" color="var(--color-text-secondary)">
                <TrendCharts />
              </el-icon>
              <p class="chart-text">用户增长趋势图表将显示在这里</p>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :lg="12">
        <el-card class="chart-card">
          <template #header>
            <div class="chart-header">
              <h3 class="chart-title">课程评分分布</h3>
              <el-tag type="success" size="small">实时数据</el-tag>
            </div>
          </template>
          <div class="chart-container">
            <div class="chart-placeholder">
              <el-icon size="48" color="var(--color-text-secondary)">
                <PieChart />
              </el-icon>
              <p class="chart-text">课程评分分布图表将显示在这里</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 详细统计 -->
    <el-card class="detailed-stats-section">
      <template #header>
        <h3 class="section-title">详细统计信息</h3>
      </template>
      
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="用户统计" name="users">
          <el-table :data="userStats" stripe>
            <el-table-column prop="metric" label="指标" />
            <el-table-column prop="value" label="数值" />
            <el-table-column prop="change" label="变化" />
            <el-table-column prop="trend" label="趋势">
              <template #default="{ row }">
                <el-tag :type="row.trend === 'up' ? 'success' : 'danger'" size="small">
                  {{ row.trend === 'up' ? '↗' : '↘' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        
        <el-tab-pane label="课程统计" name="courses">
          <el-table :data="courseStats" stripe>
            <el-table-column prop="metric" label="指标" />
            <el-table-column prop="value" label="数值" />
            <el-table-column prop="change" label="变化" />
            <el-table-column prop="trend" label="趋势">
              <template #default="{ row }">
                <el-tag :type="row.trend === 'up' ? 'success' : 'danger'" size="small">
                  {{ row.trend === 'up' ? '↗' : '↘' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        
        <el-tab-pane label="评论统计" name="comments">
          <el-table :data="commentStats" stripe>
            <el-table-column prop="metric" label="指标" />
            <el-table-column prop="value" label="数值" />
            <el-table-column prop="change" label="变化" />
            <el-table-column prop="trend" label="趋势">
              <template #default="{ row }">
                <el-tag :type="row.trend === 'up' ? 'success' : 'danger'" size="small">
                  {{ row.trend === 'up' ? '↗' : '↘' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import api from '@/services/api'

const refreshing = ref(false)
const timeRange = ref('30d')
const dataType = ref('users')
const activeTab = ref('users')

const stats = ref({
  total_users: 0,
  total_courses: 0,
  total_comments: 0,
  average_rating: 0,
})

const userStats = ref([
  { metric: '新注册用户', value: '156', change: '+12%', trend: 'up' },
  { metric: '活跃用户', value: '1,234', change: '+8%', trend: 'up' },
  { metric: '用户留存率', value: '78.5%', change: '+2.3%', trend: 'up' },
  { metric: '平均在线时长', value: '45分钟', change: '-5分钟', trend: 'down' }
])

const courseStats = ref([
  { metric: '新增课程', value: '23', change: '+15%', trend: 'up' },
  { metric: '课程访问量', value: '8,945', change: '+22%', trend: 'up' },
  { metric: '平均评分', value: '4.2', change: '+0.3', trend: 'up' },
  { metric: '课程完成率', value: '65.8%', change: '+3.2%', trend: 'up' }
])

const commentStats = ref([
  { metric: '新增评论', value: '456', change: '+18%', trend: 'up' },
  { metric: '评论回复率', value: '89.2%', change: '+5.1%', trend: 'up' },
  { metric: '平均评论长度', value: '128字', change: '+12字', trend: 'up' },
  { metric: '评论质量评分', value: '4.1', change: '+0.2', trend: 'up' }
])

const timeRangeText = computed(() => {
  const rangeMap: Record<string, string> = {
    '7d': '最近7天',
    '30d': '最近30天',
    '90d': '最近90天',
    '1y': '最近一年'
  }
  return rangeMap[timeRange.value] || '最近30天'
})

const loadStats = async () => {
  try {
    const response = await api.get('/admin/stats')
    if (response.data.success) {
      stats.value = response.data.data
    }
  } catch {
    console.error('获取统计数据失败')
  }
}

const refreshData = async () => {
  refreshing.value = true
  try {
    await loadStats()
    ElMessage.success('数据刷新成功')
  } catch {
    ElMessage.error('数据刷新失败')
  } finally {
    refreshing.value = false
  }
}

const exportData = () => {
  ElMessage.info('导出功能开发中...')
}

const handleTimeRangeChange = (value: string) => {
  console.log('时间范围变更:', value)
  // 这里可以根据时间范围重新加载数据
}

const handleDataTypeChange = (value: string) => {
  console.log('数据类型变更:', value)
  // 这里可以根据数据类型重新加载数据
}

const handleTabChange = (tabName: string) => {
  console.log('切换标签页:', tabName)
  // 这里可以根据标签页重新加载数据
}

onMounted(() => {
  loadStats()
})
</script>

<style scoped>
.stats-view {
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

.control-card {
  margin-bottom: 30px;
}

.control-card .el-form-item {
  margin-bottom: 0;
}

.charts-section {
  margin-bottom: 30px;
}

.chart-card {
  height: 400px;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-title {
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0;
}

.chart-container {
  height: 300px;
}

.chart-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  background: var(--color-background);
  border: 1px solid var(--color-border);
  border-radius: 8px;
}

.chart-text {
  font-size: 0.9rem;
  color: var(--color-text-secondary);
  margin: 15px 0 0 0;
}

.detailed-stats-section {
  margin-bottom: 30px;
}

.section-title {
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .stats-view {
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
  
  .chart-card {
    height: 300px;
  }
  
  .chart-container {
    height: 200px;
  }
}
</style>