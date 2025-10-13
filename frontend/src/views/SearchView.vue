<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { courseService, type CourseWithRating } from '@/services/api'
import CourseCard from '@/components/CourseCard.vue'

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



// 搜索课程
const searchCourses = async () => {
  if (!searchQuery.value.trim()) return
  
  loading.value = true
  try {
    // 根据搜索类型构建正确的查询参数
    const filters: any = {}
    if (searchType.value === 'course') {
      filters.name = searchQuery.value.trim()
    } else {
      filters.teacher = searchQuery.value.trim()
    }
    
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
    courses.value = []
  } finally {
    loading.value = false
  }
}


// 清空搜索
const clearSearch = () => {
  searchQuery.value = ''
  courses.value = []
}

// 切换搜索类型下拉菜单
const toggleSearchType = () => {
  showSearchDropdown.value = !showSearchDropdown.value
}

// 选择搜索类型
const selectSearchType = (type: 'course' | 'teacher') => {
  searchType.value = type
  showSearchDropdown.value = false
}

const showSearchDropdown = ref(false)

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
          <div class="search-input-container">
            <!-- 搜索类型选择器 -->
            <div class="search-type-dropdown">
              <button 
                type="button" 
                class="search-type-button"
                @click="toggleSearchType"
              >
                <span class="search-type-text">{{ searchType === 'course' ? '课程' : '老师' }}</span>
                <svg class="search-type-arrow" width="12" height="12" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M6 9l6 6 6-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              </button>
              <div v-if="showSearchDropdown" class="search-type-options">
                <button 
                  type="button"
                  class="search-type-option"
                  :class="{ active: searchType === 'course' }"
                  @click="selectSearchType('course')"
                >
                  课程
                </button>
                <button 
                  type="button"
                  class="search-type-option"
                  :class="{ active: searchType === 'teacher' }"
                  @click="selectSearchType('teacher')"
                >
                  老师
                </button>
              </div>
            </div>
            
            <!-- 搜索输入框 -->
            <input
              type="text"
              v-model="searchQuery"
              class="search-input"
              :placeholder="searchType === 'course' ? '输入课程名称...' : '输入老师姓名...'"
              required
            />
            
            <!-- 搜索按钮 -->
            <button type="submit" class="search-button" :disabled="loading">
              <span v-if="loading" class="loading-spinner"></span>
              <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <circle cx="11" cy="11" r="8" stroke="currentColor" stroke-width="2"/>
                <path d="m21 21-4.35-4.35" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              {{ loading ? '搜索中...' : '搜索' }}
            </button>
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
        <CourseCard 
          v-for="course in courses"
          :key="course.ID"
          :course="course"
          :show-upload-button="true"
          :show-comments="true"
          :show-materials="true"
        />
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
  max-width: 1200px;
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
  border: 1px solid #000000;
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
  border: 1px solid #000000;
  padding: 24px;
}

.search-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* 搜索输入容器样式 - 参考HomeView设计 */
.search-input-container {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #FFFFFF;
  border: 1px solid #000000;
  border-radius: 12px;
  padding: 8px;
  box-shadow: none;
  transition: none;
}

.search-input-container:focus-within {
  box-shadow: none;
}

/* 搜索类型下拉菜单 */
.search-type-dropdown {
  position: relative;
}

.search-type-button {
  display: flex;
  align-items: center;
  gap: 4px;
  background: #F7D074;
  border: 1px solid #000000;
  border-radius: 8px;
  padding: 8px 12px;
  font-size: 14px;
  font-weight: bold;
  color: #1A1A1A;
  cursor: pointer;
  transition: background-color 0.2s ease;
  white-space: nowrap;
}

.search-type-button:hover {
  background: #F0C74A;
}

.search-type-arrow {
  transition: transform 0.2s ease;
}

.search-type-dropdown.open .search-type-arrow {
  transform: rotate(180deg);
}

.search-type-options {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: #FFFFFF;
  border: 1px solid #000000;
  border-radius: 8px;
  z-index: 1000;
  box-shadow: none;
}

.search-type-option {
  display: block;
  width: 100%;
  background: #FFFFFF;
  border: none;
  border-bottom: 1px solid #000000;
  padding: 8px 12px;
  font-size: 14px;
  font-weight: bold;
  color: #1A1A1A;
  cursor: pointer;
  transition: background-color 0.2s ease;
  text-align: left;
}

.search-type-option:last-child {
  border-bottom: none;
}

.search-type-option:hover {
  background: #F7D074;
}

.search-type-option.active {
  background: #F7D074;
}

/* 搜索输入框 */
.search-input {
  border: none;
  outline: none !important;
  padding: 8px 12px;
  font-size: 14px;
  font-weight: bold;
  min-width: 200px;
  background: transparent;
  color: #1A1A1A;
  flex: 1;
}

.search-input::placeholder {
  color: #999;
  font-weight: normal;
}

/* 搜索按钮 */
.search-button {
  display: flex;
  align-items: center;
  gap: 6px;
  background: #76D7C4;
  color: #1A1A1A;
  border: 1px solid #000000;
  border-radius: 8px;
  padding: 8px 16px;
  font-size: 14px;
  font-weight: bold;
  cursor: pointer;
  transition: background-color 0.2s ease;
  white-space: nowrap;
}

.search-button:hover:not(:disabled) {
  background: #6BC4B1;
}

.search-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Results Section */
.results-section {
  margin-bottom: 24px;
}

.results-header {
  background-color: #FFFFFF;
  border-radius: 12px;
  border: 1px solid #000000;
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
  border: 1px solid #000000;
  padding: 8px 16px;
  font-size: 14px;
  font-weight: bold;
  color: #1A1A1A;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.clear-button:hover {
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
  border: 1px solid #000000;
  overflow: hidden;
  text-decoration: none;
  color: #1A1A1A;
  display: flex;
  flex-direction: column;
  transition: transform 0.2s ease;
}

.course-card:hover {
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
  border: 1px solid #000000;
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
  border: 1px solid #000000;
  transition: transform 0.2s ease;
}

.course-teacher-section:hover {
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
  border: 1px solid #000000;
  white-space: nowrap;
  flex-shrink: 0;
}

.course-rating-section {
  padding: 24px;
  background-color: #FFFFFF;
  border-radius: 12px;
  border: 1px solid #000000;
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
  border: 1px solid #000000;
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
  border: 1px solid #000000;
  padding: 12px 24px;
  font-size: 15px;
  font-weight: bold;
  color: #1A1A1A;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.btn-rate-course:hover {
}

.btn-rate-course:active {
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
  transform: none;
}

/* Loading State */
.loading-container {
  text-align: center;
  padding: 40px 0;
}

.loader {
  border: 1px solid var(--color-border);
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
  border: 1px solid transparent;
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
  border: 1px solid #000000;
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
  border: 1px solid #000000;
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
  border: 1px solid #000000;
  padding: 12px 24px;
  font-size: 16px;
  font-weight: bold;
  color: #1A1A1A;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.btn:hover {
}

.btn-primary {
  background-color: #F7D074;
}

.mt-lg {
  margin-top: 16px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .search-input-container {
    flex-direction: column;
    gap: 8px;
    padding: 12px;
  }
  
  .search-type-button {
    width: 100%;
    justify-content: center;
  }
  
  .search-input {
    min-width: auto;
    width: 100%;
  }
  
  .search-button {
    width: 100%;
    justify-content: center;
  }
}

@media (min-width: 768px) {
  .search-container {
    max-width: 768px;
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

