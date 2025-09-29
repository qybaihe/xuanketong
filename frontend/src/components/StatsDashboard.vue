<template>
  <div class="stats-dashboard">
    <!-- 总体统计卡片区域 -->
    <div class="stats-overview-section mb-5">
      <div class="section-header">
        <h2 class="section-title">平台数据统计</h2>
        <p class="section-subtitle">实时了解平台运行状况和用户活跃度</p>
      </div>
      
      <div class="stats-cards-grid">
        <div class="stat-card">
          <div class="stat-card-header">
            <div class="stat-icon-wrapper bg-primary">
              <i class="bi bi-book"></i>
            </div>
            <div class="stat-trend positive">
              <i class="bi bi-arrow-up"></i>
              <span>+12%</span>
            </div>
          </div>
          <div class="stat-card-body">
            <h3 class="stat-number">{{ formatNumber(statsData.overview_data.total_courses) }}</h3>
            <p class="stat-label">总课程数</p>
            <div class="stat-comparison">
              <span class="comparison-text">较上月新增 {{ calculateMonthlyIncrease('courses') }} 门课程</span>
            </div>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-card-header">
            <div class="stat-icon-wrapper bg-success">
              <i class="bi bi-people"></i>
            </div>
            <div class="stat-trend positive">
              <i class="bi bi-arrow-up"></i>
              <span>+8%</span>
            </div>
          </div>
          <div class="stat-card-body">
            <h3 class="stat-number">{{ formatNumber(statsData.overview_data.total_users) }}</h3>
            <p class="stat-label">总用户数</p>
            <div class="stat-comparison">
              <span class="comparison-text">本周活跃用户 {{ statsData.overview_data.active_users_this_week }} 人</span>
            </div>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-card-header">
            <div class="stat-icon-wrapper bg-warning">
              <i class="bi bi-star"></i>
            </div>
            <div class="stat-trend positive">
              <i class="bi bi-arrow-up"></i>
              <span>+5%</span>
            </div>
          </div>
          <div class="stat-card-body">
            <h3 class="stat-number">{{ statsData.overview_data.average_rating.toFixed(1) }}</h3>
            <p class="stat-label">平均评分</p>
            <div class="stat-comparison">
              <span class="comparison-text">基于 {{ formatNumber(statsData.overview_data.total_ratings) }} 条评价</span>
            </div>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-card-header">
            <div class="stat-icon-wrapper bg-info">
              <i class="bi bi-chat-text"></i>
            </div>
            <div class="stat-trend positive">
              <i class="bi bi-arrow-up"></i>
              <span>+15%</span>
            </div>
          </div>
          <div class="stat-card-body">
            <h3 class="stat-number">{{ formatNumber(statsData.overview_data.total_comments) }}</h3>
            <p class="stat-label">总评论数</p>
            <div class="stat-comparison">
              <span class="comparison-text">用户参与度持续提升</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 用户活动统计 -->
    <div class="user-activity-section mb-5">
      <div class="section-header">
        <h3 class="section-subtitle">用户活动分析</h3>
      </div>
      
      <div class="user-activity-grid">
        <div class="activity-card">
          <div class="activity-header">
            <h4 class="activity-title">今日活跃</h4>
            <i class="bi bi-calendar-day"></i>
          </div>
          <div class="activity-number">{{ statsData.user_activity_stats.active_users_today }}</div>
          <div class="activity-change positive">
            <i class="bi bi-arrow-up"></i>
            {{ calculateGrowth(
              statsData.user_activity_stats.active_users_today,
              statsData.user_activity_stats.active_users_today
            ) }}%
          </div>
        </div>

        <div class="activity-card">
          <div class="activity-header">
            <h4 class="activity-title">本周活跃</h4>
            <i class="bi bi-calendar-week"></i>
          </div>
          <div class="activity-number">{{ statsData.user_activity_stats.active_users_this_week }}</div>
          <div class="activity-change positive">
            <i class="bi bi-arrow-up"></i>
            +{{ statsData.user_activity_stats.new_users_this_week }} 新用户
          </div>
        </div>

        <div class="activity-card">
          <div class="activity-header">
            <h4 class="activity-title">本月活跃</h4>
            <i class="bi bi-calendar-month"></i>
          </div>
          <div class="activity-number">{{ statsData.user_activity_stats.active_users_this_month }}</div>
          <div class="activity-change positive">
            <i class="bi bi-arrow-up"></i>
            +{{ statsData.user_activity_stats.new_users_this_month }} 新用户
          </div>
        </div>
      </div>
    </div>

    <!-- 课程分类统计 -->
    <div class="course-distribution-section mb-5">
      <div class="section-header">
        <h3 class="section-subtitle">课程分布分析</h3>
      </div>
      
      <div class="distribution-cards">
        <div class="distribution-card">
          <h4 class="distribution-title">按科目分析</h4>
          <div class="distribution-items">
            <div v-for="(count, subject) in sortedSubjects" :key="subject" class="distribution-item">
              <span class="distribution-label">{{ subject }}</span>
              <div class="distribution-bar">
                <div class="distribution-fill" :style="{ width: getDistributionPercentage(count, 'subject') + '%' }"></div>
              </div>
              <span class="distribution-count">{{ count }}</span>
            </div>
          </div>
        </div>

        <div class="distribution-card">
          <h4 class="distribution-title">按年级分析</h4>
          <div class="distribution-items">
            <div v-for="(count, grade) in sortedGrades" :key="grade" class="distribution-item">
              <span class="distribution-label">{{ grade }}</span>
              <div class="distribution-bar">
                <div class="distribution-fill" :style="{ width: getDistributionPercentage(count, 'grade') + '%' }"></div>
              </div>
              <span class="distribution-count">{{ count }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 月度趋势图表 -->
    <div class="monthly-trends-section mb-5">
      <div class="section-header">
        <h3 class="section-subtitle">月度增长趋势</h3>
      </div>
      
      <div class="trends-chart-card">
        <div class="chart-header">
          <div class="chart-tabs">
            <button 
              v-for="tab in chartTabs" 
              :key="tab.key"
              class="chart-tab"
              :class="{ active: activeChartTab === tab.key }"
              @click="activeChartTab = tab.key"
            >
              {{ tab.label }}
            </button>
          </div>
          <div class="chart-period">
            <span>最近6个月</span>
          </div>
        </div>
        
        <div class="chart-container">
          <!-- 简化的图表展示，实际项目中可以使用 Chart.js 或 ECharts -->
          <div class="chart-data">
            <div class="chart-bars">
              <div 
                v-for="(stat, index) in statsData.monthly_stats" 
                :key="index"
                class="chart-bar-group"
              >
                <div class="chart-bar-label">{{ getMonthLabel(stat.month) }}</div>
                <div class="chart-bars-container">
                  <div
                    class="chart-bar bar-courses"
                    :style="{ height: getChartBarHeight(stat.courses, activeChartTab) + '%' }"
                    :title="getChartBarTitle(stat, activeChartTab)"
                  ></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 热门课程展示 -->
    <div class="popular-courses-section">
      <div class="section-header">
        <h3 class="section-subtitle">热门推荐课程</h3>
        <button class="section-action-btn" @click="loadMoreCourses">
          <i class="bi bi-arrow-clockwise"></i>
          刷新
        </button>
      </div>
      
      <div class="courses-grid">
        <div v-for="course in statsData.popular_courses" :key="course.id" class="stats-course-card">
          <div class="course-image">
            <img :src="course.image_url" :alt="course.name" />
            <div class="course-overlay">
              <div class="course-rating-overlay">
                <svg class="rating-stars" width="16" height="16" viewBox="0 0 24 24" fill="#F7D074">
                  <path d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z"/>
                </svg>
                <span class="rating-value">{{ course.average_rating.toFixed(1) }}</span>
              </div>
              <div class="course-engagement">
                <i class="bi bi-people"></i>
                <span>{{ course.student_engagement }}</span>
              </div>
            </div>
          </div>
          <div class="course-info">
            <h4 class="course-name">{{ course.name }}</h4>
            <p class="course-teacher">{{ course.teacher }}</p>
            <div class="course-meta">
              <span class="course-subject">{{ course.subject }}</span>
              <span class="course-grade">{{ course.grade }}</span>
            </div>
            <div class="course-stats">
              <span class="ratings-count">{{ course.total_ratings }} 人评价</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 数据更新时间 -->
    <div class="data-update-info">
      <i class="bi bi-clock"></i>
      <span>数据更新时间：{{ lastUpdateTime }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { statsApi, calculateStats } from '@/services/statsApi'
import type { EnhancedHomeStats } from '@/services/statsApi'

// 响应式数据
const statsData = ref<EnhancedHomeStats>({
  overview_data: {
    total_courses: 0,
    total_users: 0,
    total_ratings: 0,
    total_comments: 0,
    average_rating: 0,
    active_users_this_week: 0
  },
  top_rated_courses: [],
  recent_courses: [],
  popular_courses: [],
  user_activity_stats: {
    active_users_today: 0,
    active_users_this_week: 0,
    active_users_this_month: 0,
    new_users_today: 0,
    new_users_this_week: 0,
    new_users_this_month: 0
  },
  course_distribution: {
    by_subject: {},
    by_grade: {},
    by_semester: {}
  },
  monthly_stats: []
})

const loading = ref(false)
const lastUpdateTime = ref<string>('')
const activeChartTab = ref('courses')
const autoRefreshInterval = ref<number | null>(null)

// 图表选项卡
const chartTabs = [
  { key: 'courses', label: '新增课程' },
  { key: 'users', label: '新增用户' },
  { key: 'ratings', label: '新增评价' },
  { key: 'comments', label: '新增评论' }
]

// 计算属性
const sortedSubjects = computed(() => {
  const subjects = statsData.value.course_distribution.by_subject
  return Object.entries(subjects)
    .sort((a, b) => b[1] - a[1])
    .reduce((acc, [key, value]) => ({ ...acc, [key]: value }), {})
})

const sortedGrades = computed(() => {
  const grades = statsData.value.course_distribution.by_grade
  return Object.entries(grades)
    .sort((a, b) => b[1] - a[1])
    .reduce((acc, [key, value]) => ({ ...acc, [key]: value }), {})
})

// 方法
const loadStatsData = async () => {
  loading.value = true
  try {
    const data = await statsApi.getEnhancedHomeStats()
    statsData.value = data
    lastUpdateTime.value = new Date().toLocaleString('zh-CN')
  } catch (error) {
    console.error('加载统计数据失败:', error)
  } finally {
    loading.value = false
  }
}

const formatNumber = (num: number): string => {
  return calculateStats.formatNumber(num)
}

const getDistributionPercentage = (count: number, type: 'subject' | 'grade'): number => {
  const total = type === 'subject' 
    ? Object.values(statsData.value.course_distribution.by_subject).reduce((sum, c) => sum + c, 0)
    : Object.values(statsData.value.course_distribution.by_grade).reduce((sum, c) => sum + c, 0)
  return total > 0 ? (count / total) * 100 : 0
}

const getMonthLabel = (monthStr: string): string => {
  const [year, month] = monthStr.split('-')
  return `${month}月`
}

const getChartBarHeight = (value: number, chartType: string): number => {
    const chartData = statsData.value.monthly_stats.map(stat => stat[chartType as keyof typeof stat] as number)
    const maxValue = Math.max(...chartData)
  return maxValue > 0 ? (value / maxValue) * 100 : 0
}

const getChartBarTitle = (stat: any, chartType: string): string => {
  return `${getMonthLabel(stat.month)}: ${stat[chartType]} ${getChartLabel(chartType)}`
}

const getChartLabel = (chartType: string): string => {
  const labels = {
    courses: '门课程',
    users: '个用户',
    ratings: '条评价',
    comments: '条评论'
  }
  return labels[chartType as keyof typeof labels] || ''
}

const calculateGrowth = (current: number, previous: number): number => {
  return calculateStats.growthRate(current, previous)
}

const calculateMonthlyIncrease = (type: 'courses' | 'users' | 'ratings' | 'comments'): number => {
  const monthlyData = statsData.value.monthly_stats
  if (monthlyData.length < 2) return 0
  
  const currentMonth = monthlyData[monthlyData.length - 1]
  const previousMonth = monthlyData[monthlyData.length - 2]
  
  return (currentMonth[type] || 0) - (previousMonth[type] || 0)
}

const loadMoreCourses = async () => {
  // 模拟刷新数据
  await loadStatsData()
}

const startAutoRefresh = () => {
  // 每30秒自动刷新一次数据
  autoRefreshInterval.value = window.setInterval(async () => {
    await loadStatsData()
  }, 30000)
}

const stopAutoRefresh = () => {
  if (autoRefreshInterval.value) {
    clearInterval(autoRefreshInterval.value)
    autoRefreshInterval.value = null
  }
}

// 生命周期
onMounted(async () => {
  await loadStatsData()
  startAutoRefresh()
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>

<style scoped>
.stats-dashboard {
  padding: 2rem;
  background: linear-gradient(135deg, #f8fafb 0%, #e3f2fd 100%);
  min-height: 100vh;
}

.section-header {
  text-align: center;
  margin-bottom: 2rem;
}

.section-title {
  font-size: 2.5rem;
  font-weight: 700;
  color: #2c3e50;
  margin-bottom: 0.5rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.section-subtitle {
  font-size: 1.1rem;
  color: #6c757d;
  margin-bottom: 1rem;
}

/* 统计卡片样式 */
.stats-cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: white;
  border-radius: 16px;
  padding: 1.5rem;
  box-shadow: 0 8px 32px rgba(0,0,0,0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  border: 1px solid rgba(255,255,255,0.2);
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 16px 48px rgba(0,0,0,0.15);
}

.stat-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.stat-icon-wrapper {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 1.5rem;
}

.stat-icon-wrapper.bg-primary { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
.stat-icon-wrapper.bg-success { background: linear-gradient(135deg, #28a745 0%, #20c997 100%); }
.stat-icon-wrapper.bg-warning { background: linear-gradient(135deg, #ffc107 0%, #fd7e14 100%); }
.stat-icon-wrapper.bg-info { background: linear-gradient(135deg, #17a2b8 0%, #138496 100%); }

/* 用户活动统计样式 */
.user-activity-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.activity-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 4px 16px rgba(0,0,0,0.08);
  text-align: center;
  transition: all 0.3s ease;
}

.activity-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.12);
}

.activity-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1rem;
}

.activity-title {
  font-size: 1rem;
  font-weight: 600;
  color: #495057;
  margin: 0;
}

.activity-number {
  font-size: 2rem;
  font-weight: 700;
  color: #2c3e50;
  margin-bottom: 0.5rem;
}

.activity-change {
  font-size: 0.85rem;
  font-weight: 500;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.25rem;
}

.activity-change.positive {
  color: #28a745;
}

.activity-change.negative {
  color: #dc3545;
}

/* 分布统计样式 */
.distribution-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 2rem;
  margin-bottom: 2rem;
}

.distribution-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 4px 16px rgba(0,0,0,0.08);
}

.distribution-title {
  font-size: 1.2rem;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 1rem;
}

.distribution-items {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.distribution-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.distribution-label {
  min-width: 60px;
  font-size: 0.9rem;
  color: #6c757d;
}

.distribution-bar {
  flex: 1;
  height: 8px;
  background: #e9ecef;
  border-radius: 4px;
  overflow: hidden;
}

.distribution-fill {
  height: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  transition: width 0.3s ease;
}

.distribution-count {
  min-width: 30px;
  font-size: 0.85rem;
  font-weight: 600;
  color: #2c3e50;
  text-align: right;
}

/* 图表样式 */
.trends-chart-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 4px 16px rgba(0,0,0,0.08);
  margin-bottom: 2rem;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.chart-tabs {
  display: flex;
  gap: 0.5rem;
}

.chart-tab {
  padding: 0.5rem 1rem;
  border: none;
  background: #f8f9fa;
  color: #6c757d;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 0.9rem;
}

.chart-tab:hover {
  background: #e9ecef;
}

.chart-tab.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.chart-period {
  font-size: 0.85rem;
  color: #6c757d;
}

.chart-data {
  padding: 1rem 0;
}

.chart-bars {
  display: flex;
  align-items: end;
  gap: 1rem;
  height: 200px;
}

.chart-bar-group {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
  min-width: 0;
}

.chart-bar-label {
  font-size: 0.8rem;
  color: #6c757d;
  margin-bottom: 0.5rem;
  text-align: center;
}

.chart-bars-container {
  display: flex;
  align-items: end;
  width: 100%;
  height: 160px;
}

.chart-bar {
  width: 100%;
  max-width: 30px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 4px 4px 0 0;
  transition: all 0.3s ease;
  cursor: pointer;
}

.chart-bar:hover {
  opacity: 0.8;
}

/* 热门课程卡片样式 */
.popular-courses-section {
  margin-bottom: 2rem;
}

.courses-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stats-course-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 8px 24px rgba(0,0,0,0.08);
  transition: all 0.3s ease;
}

.stats-course-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 16px 40px rgba(0,0,0,0.12);
}

.course-image {
  position: relative;
  aspect-ratio: 16/9;
  overflow: hidden;
}

.course-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.stats-course-card:hover .course-image img {
  transform: scale(1.05);
}

.course-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(transparent 60%, rgba(0,0,0,0.8));
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  padding: 1rem;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.stats-course-card:hover .course-overlay {
  opacity: 1;
}

.course-rating-overlay {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  color: white;
}

.rating-stars {
  font-size: 0.8rem;
}

.rating-value {
  font-size: 0.9rem;
  font-weight: 600;
}

.course-engagement {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  color: white;
  font-size: 0.8rem;
}

.course-info {
  padding: 1rem;
}

.course-name {
  font-size: 1.1rem;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 0.5rem;
  line-height: 1.3;
}

.course-teacher {
  font-size: 0.9rem;
  color: #6c757d;
  margin-bottom: 0.75rem;
}

.course-meta {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
}

.course-subject, .course-grade {
  padding: 0.25rem 0.5rem;
  background: #f8f9fa;
  border-radius: 4px;
  font-size: 0.75rem;
  color: #495057;
}

.course-stats {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 0.75rem;
  border-top: 1px solid #e9ecef;
}

.ratings-count {
  font-size: 0.8rem;
  color: #6c757d;
}

/* 数据更新信息 */
.data-update-info {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 1rem;
  color: #6c757d;
  font-size: 0.9rem;
}

.data-update-info i {
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% { opacity: 1; }
  50% { opacity: 0.5; }
  100% { opacity: 1; }
}

.section-action-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.section-action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(102, 126, 234, 0.4);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .stats-dashboard {
    padding: 1rem;
  }
  
  .section-title {
    font-size: 2rem;
  }
  
  .stats-cards-grid {
    grid-template-columns: 1fr;
  }
  
  .courses-grid {
    grid-template-columns: 1fr;
  }
  
  .distribution-cards {
    grid-template-columns: 1fr;
  }
  
  .chart-bars {
    height: 150px;
  }
}
</style>