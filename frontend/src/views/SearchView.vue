<script setup lang="ts">
import { ref, onMounted, reactive, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { courseService, type Course, type CourseWithRating } from '@/services/api'

const router = useRouter()
const route = useRoute()

interface CourseWithDisplay extends CourseWithRating {
  Students?: number;
  isPopular?: boolean;
  isNew?: boolean;
  RatingDistribution?: Record<number, number>;
}

const courses = ref<CourseWithDisplay[]>([])
const loading = ref(true)
const searchQuery = ref('')
const searchType = ref<'course' | 'teacher'>('course')

// 从路由参数获取搜索条件
const searchParams = reactive({
  query: '',
  type: 'course' as 'course' | 'teacher'
})

// 数据字段映射函数
const mapCourseData = (course: any): CourseWithDisplay => {
  return {
    id: course.id,
    name: course.name,
    description: course.description,
    teacher: course.teacher,
    credits: course.credits,
    grade: course.grade,
    semester: course.semester,
    subject: course.subject,
    createdAt: course.createdAt,
    updatedAt: course.updatedAt,
    averageRating: course.averageRating,
    totalRatings: course.totalRatings,
    ratingDistribution: course.ratingDistribution,
    ID: course.id,
    Name: course.name,
    Description: course.description,
    Teacher: course.teacher,
    Credits: course.credits,
    Grade: course.grade,
    Semester: course.semester,
    Subject: course.subject,
    CreatedAt: course.createdAt,
    UpdatedAt: course.updatedAt,
    AverageRating: course.averageRating,
    TotalRatings: course.totalRatings,
    RatingDistribution: course.ratingDistribution
  }
}

// 获取标签
const getTags = (course: CourseWithDisplay) => {
  const tags = []
  if (course.Grade) tags.push({ text: course.Grade, type: 'primary' })
  if (course.Semester) tags.push({ text: course.Semester, type: 'secondary' })
  if (course.Subject) tags.push({ text: course.Subject, type: 'accent' })
  return tags
}

// 获取评分百分比
const getRatingPercentage = (course: CourseWithDisplay, starLevel: number) => {
  const ratingDistribution = course.RatingDistribution || {}
  const countForStar = ratingDistribution[starLevel] || 0
  const totalRatings = course.TotalRatings || 1
  const percentage = (countForStar / totalRatings) * 100
  return Math.max(5, Math.min(95, percentage))
}

// 获取评分数量
const getRatingCount = (course: CourseWithDisplay, starLevel: number) => {
  const ratingDistribution = course.RatingDistribution || {}
  return ratingDistribution[starLevel] || 0
}

// 获取显示评分
const getDisplayRating = (course: CourseWithDisplay): number => {
  return course.AverageRating || 0
}

// 搜索课程
const searchCourses = async () => {
  if (!searchQuery.value.trim()) return
  
  loading.value = true
  try {
    const filters = searchType.value === 'course' 
      ? { name: searchQuery.value }
      : { teacher: searchQuery.value }
    
    const coursesData = await courseService.getCourses(filters)
    
    courses.value = coursesData.map((course: any) => {
      const mappedCourse = mapCourseData(course)
      return {
        ...mappedCourse,
        Students: Math.floor(Math.random() * 200) + 10,
        isPopular: (mappedCourse.AverageRating || 0) > 4.0,
        isNew: mappedCourse.CreatedAt ? new Date().getTime() - new Date(mappedCourse.CreatedAt).getTime() < 30 * 24 * 60 * 60 * 1000 : false
      }
    })
  } catch (error) {
    console.error('搜索课程失败', error)
  } finally {
    loading.value = false
  }
}

// 跳转到课程评价页面
const goToRateCourse = (courseId: number | undefined) => {
  if (courseId) {
    router.push(`/courses/${courseId}/rate`)
  }
}

// 清空搜索
const clearSearch = () => {
  searchQuery.value = ''
  courses.value = []
}

// 计算搜索结果统计
const searchStats = computed(() => {
  const total = courses.value.length
  const popular = courses.value.filter(c => c.isPopular).length
  const newCourses = courses.value.filter(c => c.isNew).length
  return { total, popular, newCourses }
})

// 初始化
onMounted(() => {
  // 从路由参数获取搜索条件
  if (route.query.q) {
    searchQuery.value = route.query.q as string
  }
  if (route.query.type) {
    searchType.value = route.query.type as 'course' | 'teacher'
  }
  
  // 如果有搜索条件，执行搜索
  if (searchQuery.value) {
    searchCourses()
  } else {
    loading.value = false
  }
})
</script>

<template>
  <main class="search-container">
    <!-- Header Section -->
    <header class="search-header">
      <div class="header-content">
        <h1 class="page-title">课程搜索</h1>
        <p class="page-subtitle">找到最适合您的课程</p>
      </div>
    </header>

    <!-- Search Section -->
    <section class="search-section">
      <div class="search-form-container">
        <form @submit.prevent="searchCourses" class="search-form">
          <div class="search-inputs">
            <div class="search-type-selector">
              <label class="search-type-label">搜索类型</label>
              <div class="search-type-options">
                <label class="search-type-option">
                  <input 
                    type="radio" 
                    v-model="searchType" 
                    value="course" 
                    class="search-type-radio"
                  />
                  <span class="search-type-text">按课程名</span>
                </label>
                <label class="search-type-option">
                  <input 
                    type="radio" 
                    v-model="searchType" 
                    value="teacher" 
                    class="search-type-radio"
                  />
                  <span class="search-type-text">按老师名</span>
                </label>
              </div>
            </div>
            
            <div class="search-input-group">
              <input
                type="text"
                v-model="searchQuery"
                class="search-input"
                :placeholder="searchType === 'course' ? '输入课程名称...' : '输入老师姓名...'"
                required
              />
              <button type="submit" class="search-button" :disabled="loading">
                <span v-if="loading" class="loading-spinner"></span>
                {{ loading ? '搜索中...' : '搜索' }}
              </button>
            </div>
          </div>
        </form>
      </div>
    </section>

    <!-- Search Results -->
    <section v-if="searchQuery" class="results-section">
      <div class="results-header">
        <div class="results-stats">
          <h2 class="results-title">
            搜索结果
            <span class="results-count">({{ searchStats.total }} 个结果)</span>
          </h2>
          <div class="results-summary">
            <span class="stat-item">热门课程: {{ searchStats.popular }}</span>
            <span class="stat-item">新课程: {{ searchStats.newCourses }}</span>
          </div>
        </div>
        <button @click="clearSearch" class="clear-button">
          清空搜索
        </button>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="loading-container">
        <div class="loader"></div>
        <p class="loading-text">正在搜索课程...</p>
      </div>

      <!-- Results Grid -->
      <div v-else-if="courses.length > 0" class="course-grid">
        <router-link
          v-for="course in courses"
          :key="course.ID"
          :to="`/courses/${course.ID}`"
          class="course-card"
        >
          <div class="course-card-content">
            <!-- 课程标签 -->
            <div class="course-card-tags">
              <span
                v-for="(tag, index) in getTags(course)"
                :key="index"
                :class="['course-card-tag', tag.type]"
              >
                {{ tag.text }}
              </span>
            </div>

            <!-- 课程标题和基本信息 -->
            <div class="course-card-header">
              <h3 class="course-card-title">{{ course.Name }}</h3>
              <div class="course-teacher-section">
                <div class="teacher-label">授课老师</div>
                <div class="teacher-info">
                  <div class="teacher-details">
                    <div class="teacher-name">{{ course.Teacher }}</div>
                    <div class="teacher-credits">{{ course.Credits }} credits</div>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- 评价统计区域 -->
            <div class="course-rating-section">
              <div class="rating-overview">
                <div class="rating-average">
                  <span class="rating-number">{{ getDisplayRating(course).toFixed(1) }}</span>
                  <div class="rating-stars">
                    <span
                      v-for="i in 5"
                      :key="i"
                      :class="['star',
                        i <= Math.floor(getDisplayRating(course)) ? 'star-filled' :
                        (i - 0.5 <= getDisplayRating(course) ? 'star-half' : 'star-empty')]"
                    ></span>
                  </div>
                  <span class="rating-count">{{ (course.TotalRatings || 0) }} 人评价</span>
                </div>
                <div class="rating-breakdown">
                  <div class="rating-bar" v-for="i in 5" :key="i">
                    <span class="bar-label">{{ 6-i }}星</span>
                    <div class="bar-container">
                      <div class="bar-fill" :style="{ width: getRatingPercentage(course, 6-i) + '%' }"></div>
                    </div>
                    <span class="bar-count">{{ getRatingCount(course, 6-i) }}</span>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- 课程操作区域 -->
            <div class="course-card-actions">
              <button class="btn-rate-course" @click.prevent="goToRateCourse(course.ID)">
                评价课程
              </button>
            </div>
          </div>
        </router-link>
      </div>

      <!-- Empty State -->
      <div v-else class="empty-state">
        <h3 class="empty-title">未找到相关课程</h3>
        <p class="empty-description">请尝试其他搜索关键词</p>
        <button @click="clearSearch" class="btn btn-primary mt-lg">
          重新搜索
        </button>
      </div>
    </section>

    <!-- Initial State -->
    <section v-else class="initial-state">
      <div class="initial-content">
        <h2 class="initial-title">开始搜索课程</h2>
        <p class="initial-description">输入课程名称或老师姓名来查找相关课程</p>
      </div>
    </section>
  </main>
</template>

<style scoped>
/* ===== 俏皮的新拟物风格 (Playful Neobrutalism) ===== */

/* 全局样式 */
.search-container {
  font-family: sans-serif;
  color: #1A1A1A;
  padding: 20px;
  max-width: 420px;
  margin: 0 auto;
  background-color: #FEF6F7;
  min-height: 100vh;
}

/* Header Section */
.search-header {
  margin-bottom: 24px;
  text-align: center;
}

.header-content {
  background-color: #FFFFFF;
  border-radius: 12px;
  border: 3px solid #000000;
  box-shadow: 5px 5px 0px 0px #000000;
  padding: 24px;
}

.page-title {
  font-size: 28px;
  font-weight: bold;
  text-align: center;
  color: #1A1A1A;
  margin-bottom: 8px;
}

.page-subtitle {
  font-size: 16px;
  color: #888888;
  margin: 0;
}

/* Search Section */
.search-section {
  margin-bottom: 24px;
}

.search-form-container {
  background-color: #FFFFFF;
  border-radius: 12px;
  border: 3px solid #000000;
  box-shadow: 5px 5px 0px 0px #000000;
  padding: 24px;
}

.search-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.search-inputs {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.search-type-selector {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.search-type-label {
  font-size: 14px;
  font-weight: bold;
  color: #1A1A1A;
}

.search-type-options {
  display: flex;
  gap: 16px;
}

.search-type-option {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.search-type-radio {
  width: 16px;
  height: 16px;
  accent-color: #F7D074;
}

.search-type-text {
  font-size: 14px;
  color: #1A1A1A;
}

.search-input-group {
  display: flex;
  gap: 12px;
}

.search-input {
  flex: 1;
  background-color: #FFFFFF;
  border: 3px solid #000000;
  border-radius: 8px;
  padding: 12px;
  font-size: 16px;
  font-family: sans-serif;
  transition: transform 0.2s ease;
}

.search-input:focus {
  transform: translate(-2px, -2px);
  box-shadow: 4px 4px 0px 0px #000000;
  outline: none;
}

.search-button {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 24px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: bold;
  background-color: #F7D074;
  border: 3px solid #000000;
  box-shadow: 4px 4px 0px 0px #000000;
  color: #1A1A1A;
  cursor: pointer;
  transition: transform 0.2s ease;
  min-width: 100px;
}

.search-button:hover:not(:disabled) {
  transform: translate(-2px, -2px);
  box-shadow: 6px 6px 0px 0px #000000;
}

.search-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
  box-shadow: 4px 4px 0px 0px #000000;
}

/* Results Section */
.results-section {
  margin-bottom: 24px;
}

.results-header {
  background-color: #FFFFFF;
  border-radius: 12px;
  border: 3px solid #000000;
  box-shadow: 5px 5px 0px 0px #000000;
  padding: 24px;
  margin-bottom: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
}

.results-stats {
  flex: 1;
}

.results-title {
  font-size: 20px;
  font-weight: bold;
  color: #1A1A1A;
  margin-bottom: 8px;
}

.results-count {
  font-size: 16px;
  color: #888888;
  font-weight: normal;
}

.results-summary {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.stat-item {
  font-size: 14px;
  color: #888888;
}

.clear-button {
  background-color: #FFFFFF;
  border-radius: 8px;
  border: 3px solid #000000;
  box-shadow: 4px 4px 0px 0px #000000;
  padding: 8px 16px;
  font-size: 14px;
  font-weight: bold;
  color: #1A1A1A;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.clear-button:hover {
  transform: translate(-2px, -2px);
  box-shadow: 6px 6px 0px 0px #000000;
}

/* Course Grid */
.course-grid {
  display: grid;
  gap: 16px;
  margin-bottom: 24px;
}

/* Course Card Styles (复用HomeView的样式) */
.course-card {
  background-color: #FFFFFF;
  border-radius: 12px;
  border: 3px solid #000000;
  box-shadow: 5px 5px 0px 0px #000000;
  overflow: hidden;
  text-decoration: none;
  color: #1A1A1A;
  display: flex;
  flex-direction: column;
  transition: transform 0.2s ease;
}

.course-card:hover {
  transform: translate(-2px, -2px);
  box-shadow: 7px 7px 0px 0px #000000;
}

.course-card-content {
  padding: 24px;
  display: flex;
  flex-direction: column;
  flex-grow: 1;
}

.course-card-tags {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.course-card-tag {
  padding: 4px 8px;
  border-radius: 8px;
  border: 2px solid #000000;
  font-size: 12px;
  font-weight: bold;
}

.course-card-tag.primary {
  background-color: #F7D074;
  color: #1A1A1A;
}

.course-card-tag.secondary {
  background-color: #76D7C4;
  color: #1A1A1A;
}

.course-card-tag.accent {
  background-color: #F7D074;
  color: #1A1A1A;
}

.course-card-title {
  font-size: 20px;
  font-weight: bold;
  color: #1A1A1A;
  margin-bottom: 16px;
  line-height: 1.3;
}

.course-card-header {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 3px solid #000000;
}

.course-teacher-section {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background-color: #FFFFFF;
  border-radius: 8px;
  border: 3px solid #000000;
  box-shadow: 4px 4px 0px 0px #000000;
  transition: transform 0.2s ease;
}

.course-teacher-section:hover {
  transform: translate(-2px, -2px);
  box-shadow: 6px 6px 0px 0px #000000;
}

.teacher-label {
  font-size: 12px;
  color: #1A1A1A;
  font-weight: bold;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.teacher-info {
  display: flex;
  align-items: center;
  flex-grow: 1;
}

.teacher-details {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  gap: 16px;
}

.teacher-name {
  font-size: 16px;
  font-weight: bold;
  color: #1A1A1A;
  flex-shrink: 0;
}

.teacher-credits {
  font-size: 14px;
  color: #888888;
  font-weight: bold;
  background-color: #FFFFFF;
  padding: 4px 10px;
  border-radius: 8px;
  border: 2px solid #000000;
  white-space: nowrap;
  flex-shrink: 0;
}

.course-rating-section {
  padding: 24px;
  background-color: #FFFFFF;
  border-radius: 12px;
  border: 3px solid #000000;
  box-shadow: 5px 5px 0px 0px #000000;
  margin-bottom: 24px;
}

.rating-overview {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-direction: column;
  gap: 16px;
}

.rating-average {
  text-align: center;
}

.rating-number {
  font-size: 36px;
  font-weight: bold;
  color: #F7D074;
  line-height: 1;
  display: block;
}

.rating-stars {
  display: flex;
  gap: 4px;
  margin: 8px 0;
}

.rating-count {
  font-size: 14px;
  color: #888888;
  white-space: nowrap;
}

.rating-breakdown {
  width: 100%;
}

.rating-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.bar-label {
  font-size: 14px;
  color: #888888;
  min-width: 30px;
  text-align: right;
}

.bar-container {
  flex: 1;
  height: 8px;
  background-color: #FFFFFF;
  border: 2px solid #000000;
  border-radius: 4px;
  overflow: hidden;
  position: relative;
}

.bar-fill {
  height: 100%;
  background-color: #F7D074;
  transition: width 0.3s ease;
}

.bar-count {
  font-size: 14px;
  color: #888888;
  text-align: left;
}

.course-card-actions {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: auto;
  padding-top: 16px;
}

.btn-rate-course {
  background-color: #76D7C4;
  border-radius: 8px;
  border: 3px solid #000000;
  box-shadow: 4px 4px 0px 0px #000000;
  padding: 12px 24px;
  font-size: 15px;
  font-weight: bold;
  color: #1A1A1A;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.btn-rate-course:hover {
  transform: translate(-2px, -2px);
  box-shadow: 6px 6px 0px 0px #000000;
}

.btn-rate-course:active {
  transform: translate(-1px, -1px);
  box-shadow: 3px 3px 0px 0px #000000;
}

/* 评分星级样式 */
.star {
  font-size: 18px;
  transition: transform 0.2s ease;
  position: relative;
  display: inline-block;
  width: 18px;
  height: 18px;
}

.star-filled::before {
  content: "★";
  color: #F7D074;
}

.star-half::before {
  content: "★";
  color: #F7D074;
  background: linear-gradient(90deg, #F7D074 50%, #888888 50%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.star-empty::before {
  content: "★";
  color: #888888;
}

.course-card:hover .star {
  transform: scale(1.1);
}

/* Loading State */
.loading-container {
  text-align: center;
  padding: 40px 0;
}

.loader {
  border: 4px solid #000000;
  border-top: 4px solid #F7D074;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.loading-text {
  font-size: 16px;
  color: #1A1A1A;
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid transparent;
  border-top: 2px solid #1A1A1A;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 40px 0;
  background-color: #FFFFFF;
  border-radius: 12px;
  border: 3px solid #000000;
  box-shadow: 5px 5px 0px 0px #000000;
}

.empty-title {
  font-size: 20px;
  font-weight: bold;
  color: #1A1A1A;
  margin-bottom: 8px;
}

.empty-description {
  font-size: 14px;
  color: #888888;
  margin: 0 0 16px 0;
}

/* Initial State */
.initial-state {
  text-align: center;
  padding: 40px 0;
}

.initial-content {
  background-color: #FFFFFF;
  border-radius: 12px;
  border: 3px solid #000000;
  box-shadow: 5px 5px 0px 0px #000000;
  padding: 40px 24px;
}

.initial-title {
  font-size: 24px;
  font-weight: bold;
  color: #1A1A1A;
  margin-bottom: 12px;
}

.initial-description {
  font-size: 16px;
  color: #888888;
  margin: 0;
}

/* Button Styles */
.btn {
  background-color: #F7D074;
  border-radius: 8px;
  border: 3px solid #000000;
  box-shadow: 4px 4px 0px 0px #000000;
  padding: 12px 24px;
  font-size: 16px;
  font-weight: bold;
  color: #1A1A1A;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.btn:hover {
  transform: translate(-2px, -2px);
  box-shadow: 6px 6px 0px 0px #000000;
}

.btn-primary {
  background-color: #F7D074;
}

.mt-lg {
  margin-top: 16px;
}

/* 响应式设计 */
@media (min-width: 768px) {
  .search-container {
    max-width: 768px;
  }
  
  .search-input-group {
    flex-direction: row;
  }
  
  .course-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .results-header {
    flex-direction: row;
    align-items: center;
  }
}

@media (min-width: 1024px) {
  .search-container {
    max-width: 1400px;
    padding: 40px;
  }
  
  .course-grid {
    grid-template-columns: repeat(3, 1fr);
    gap: 32px;
  }
  
  .search-form {
    gap: 32px;
  }
  
  .search-inputs {
    gap: 24px;
  }
}

@media (min-width: 1440px) {
  .search-container {
    max-width: 1600px;
  }
  
  .course-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}

@media (min-width: 1600px) {
  .search-container {
    max-width: 1800px;
  }
  
  .course-grid {
    grid-template-columns: repeat(5, 1fr);
  }
}

@media (min-width: 1920px) {
  .search-container {
    max-width: 2200px;
  }
  
  .course-grid {
    grid-template-columns: repeat(6, 1fr);
  }
}
</style>

