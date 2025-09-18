<script setup lang="ts">
import { ref, onMounted, reactive, watch, computed, nextTick } from 'vue'
import { RouterLink } from 'vue-router'
import { courseService, type Course, type CourseWithRating } from '@/services/api'

interface CourseWithDisplay extends CourseWithRating {
  Students?: number;
  isPopular?: boolean;
  isNew?: boolean;
}

const courses = ref<CourseWithDisplay[]>([])
const loading = ref(true)
const courseSection = ref<HTMLElement>()
const filters = reactive({
  grade: '',
  semester: '',
  subject: ''
})

// 滚动到课程区域
const scrollToCourses = () => {
  if (courseSection.value) {
    courseSection.value.scrollIntoView({ behavior: 'smooth' })
  }
}

// 清空筛选
const clearFilters = () => {
  filters.grade = ''
  filters.semester = ''
  filters.subject = ''
  fetchCourses()
}

// 显示热门课程
const showPopularOnly = () => {
  // 这里可以添加筛选逻辑
  console.log('显示热门课程')
}

// 显示新课程
const showNewOnly = () => {
  // 这里可以添加筛选逻辑
  console.log('显示新课程')
}

// 跳转到课程评价页面
const goToRateCourse = (courseId: number | undefined) => {
  if (courseId) {
    // 使用Vue Router进行导航
    window.location.href = `#/courses/${courseId}/rate`
  }
}

// 获取评分百分比（用于评分条）
const getRatingPercentage = (rating: number, starLevel: number) => {
  // 根据评分级别计算百分比，模拟真实数据分布
  const baseDistribution = [0.05, 0.15, 0.35, 0.30, 0.15] // 1-5星的基准分布
  const ratingIndex = Math.max(0, Math.min(4, Math.floor(rating) - 1))
  const basePercent = baseDistribution[4 - (starLevel - 1)] * 100
  
  // 根据实际评分调整分布
  const adjustment = (rating % 1) * (starLevel <= rating ? 10 : -10)
  
  return Math.max(5, Math.min(95, basePercent + adjustment))
}

// 获取评分数量（用于评分条）
const getRatingCount = (rating: number, starLevel: number) => {
  const totalReviews = Math.floor((rating || 0) * 20) // 假设每0.1分对应2条评价
  const percentage = getRatingPercentage(rating, starLevel)
  return Math.floor(totalReviews * percentage / 100)
}

// 获取卡片类名
const getCardClass = (course: Course, index: number) => {
  // 简化类名以应用新设计
  return 'course-card'
}

// 数据字段映射函数 - 将后端小写字段映射到前端大写字段
const mapCourseData = (course: any): CourseWithDisplay => {
  return {
    // 原始后端字段
    id: course.id,
    name: course.name,
    description: course.description,
    teacher: course.teacher,
    credits: course.credits,
    imageURL: course.imageURL,
    grade: course.grade,
    semester: course.semester,
    subject: course.subject,
    createdAt: course.createdAt,
    updatedAt: course.updatedAt,
    averageRating: course.averageRating,
    totalRatings: course.totalRatings,
    // 兼容字段映射（大写）
    ID: course.id,
    Name: course.name,
    Description: course.description,
    Teacher: course.teacher,
    Credits: course.credits,
    ImageURL: course.imageURL || `https://picsum.photos/seed/course-${course.id}/400/200.jpg`,
    Grade: course.grade,
    Semester: course.semester,
    Subject: course.subject,
    CreatedAt: course.createdAt,
    UpdatedAt: course.updatedAt,
    AverageRating: course.averageRating,
    TotalRatings: course.totalRatings
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

// 获取评分星级
const getRatingStars = (rating: number) => {
  const fullStars = Math.floor(rating)
  const hasHalfStar = rating % 1 >= 0.5
  const emptyStars = 5 - fullStars - (hasHalfStar ? 1 : 0)
  
  return {
    full: fullStars,
    half: hasHalfStar,
    empty: emptyStars
  }
}

// 将AverageRating映射到Rating字段（兼容现有模板）
const getDisplayRating = (course: CourseWithDisplay): number => {
  return course.AverageRating || 0
}

const fetchCourses = async () => {
  console.log('fetchCourses: 开始获取课程数据，loading状态:', loading.value)
  loading.value = true
  try {
    // 获取课程列表
    const coursesData = await courseService.getCourses(filters)
    console.log('fetchCourses: 获取到课程数据，数量:', coursesData.length)
    console.log('fetchCourses: 第一个课程数据:', coursesData[0])
    
    // 添加显示用的额外属性
    courses.value = coursesData.map((course: any, index) => {
      const mappedCourse = mapCourseData(course)
      return {
        ...mappedCourse,
        Students: Math.floor(Math.random() * 200) + 10, // 暂时模拟学生数，后续可以从后端获取
        // 添加热门标记（评分大于4.0的课程为热门）
        isPopular: (mappedCourse.AverageRating || 0) > 4.0,
        // 添加新课程标记（基于创建时间，30天内为新课程）
        isNew: mappedCourse.CreatedAt ? new Date().getTime() - new Date(mappedCourse.CreatedAt).getTime() < 30 * 24 * 60 * 60 * 1000 : false
      }
    })
    
    // 添加滚动显示动画，确保DOM完全渲染
    nextTick(() => {
      console.log('fetchCourses: nextTick回调，课程卡片数量:', document.querySelectorAll('.course-card').length)
      // 额外延迟确保所有样式和动画都应用完成
      setTimeout(() => {
        observeElements()
        // 最后才设置loading为false，确保所有准备工作完成
        loading.value = false
        console.log('fetchCourses: 所有准备工作完成，loading状态设置为false')
      }, 200)
    })
  } catch (error) {
    console.error('fetchCourses: 获取课程数据失败', error)
    loading.value = false
  }
}

// 防抖处理
let debounceTimer: number
const debouncedFetchCourses = () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(fetchCourses, 300)
}

// 滚动显示观察器
const observeElements = () => {
  console.log('observeElements: 开始设置观察器')
  
  // 延迟设置观察器，确保元素完全渲染
  setTimeout(() => {
    const observer = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          console.log('observeElements: 元素进入视口，添加revealed类', entry.target)
          entry.target.classList.add('revealed')
        }
      })
    }, {
      threshold: 0.1,
      rootMargin: '0px 0px -50px 0px'
    })

    const cards = document.querySelectorAll('.course-card')
    console.log('observeElements: 找到课程卡片数量:', cards.length)
    cards.forEach(el => {
      console.log('observeElements: 观察卡片元素', el)
      observer.observe(el)
    })
  }, 100) // 延迟100ms确保元素完全渲染
}

// 添加磁性效果
const addMagneticEffect = (element: HTMLElement) => {
  element.addEventListener('mousemove', (e) => {
    const rect = element.getBoundingClientRect()
    const x = e.clientX - rect.left - rect.width / 2
    const y = e.clientY - rect.top - rect.height / 2
    
    element.style.transform = `translate(${x * 0.1}px, ${y * 0.1}px) scale(1.02)`
  })
  
  element.addEventListener('mouseleave', () => {
    element.style.transform = 'translate(0, 0) scale(1)'
  })
}

// 图片加载处理
const handleImageLoad = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.classList.add('loaded')
  console.log('图片加载成功:', img.src)
}

// 图片加载错误处理
const handleImageError = (event: Event) => {
  const img = event.target as HTMLImageElement
  console.log('图片加载失败，使用备用URL:', img.src)
  // 使用备用图片URL
  img.src = `https://picsum.photos/seed/course-${img.alt || 'default'}/400/200.jpg`
}

// 添加性能监控
const initPerformanceMonitoring = () => {
  // 页面加载时间
  window.addEventListener('load', () => {
    const loadTime = performance.now()
    console.log(`页面加载时间: ${loadTime.toFixed(2)}ms`)
  })
  
  // 添加加载指示器
  const loadIndicator = document.createElement('div')
  loadIndicator.className = 'load-indicator'
  loadIndicator.style.width = '0%'
  document.body.appendChild(loadIndicator)
  
  // 监听加载进度
  let loadProgress = 0
  const updateLoadProgress = () => {
    loadProgress += Math.random() * 30
    if (loadProgress > 90) loadProgress = 90
    loadIndicator.style.width = `${loadProgress}%`
    
    if (loadProgress < 90) {
      requestAnimationFrame(updateLoadProgress)
    } else {
      setTimeout(() => {
        loadProgress = 100
        loadIndicator.style.width = '100%'
        setTimeout(() => {
          loadIndicator.style.opacity = '0'
          setTimeout(() => {
            document.body.removeChild(loadIndicator)
          }, 300)
        }, 200)
      }, 500)
    }
  }
  
  requestAnimationFrame(updateLoadProgress)
}

// 添加离线支持
const initOfflineSupport = () => {
  const offlineIndicator = document.createElement('div')
  offlineIndicator.className = 'offline-indicator'
  offlineIndicator.textContent = '您当前处于离线状态'
  document.body.appendChild(offlineIndicator)
  
  window.addEventListener('online', () => {
    offlineIndicator.classList.remove('show')
  })
  
  window.addEventListener('offline', () => {
    offlineIndicator.classList.add('show')
  })
}

watch(filters, debouncedFetchCourses, { deep: true })

onMounted(() => {
  fetchCourses()
  
  // 初始化性能监控
  initPerformanceMonitoring()
  
  // 初始化离线支持
  initOfflineSupport()
  
  // 为按钮添加磁性效果
  nextTick(() => {
    document.querySelectorAll('.btn').forEach(btn => {
      addMagneticEffect(btn as HTMLElement)
    })
    
    // 为图片添加懒加载观察器
    const imageObserver = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          const img = entry.target as HTMLImageElement
          if (img.dataset.src) {
            img.src = img.dataset.src
            img.removeAttribute('data-src')
          }
          imageObserver.unobserve(img)
        }
      })
    })
    
    document.querySelectorAll('img[loading="lazy"]').forEach(img => {
      imageObserver.observe(img)
    })
  })
})
</script>

<template>
  <main class="home-container">
    <!-- Header Section -->
    <header class="home-header">
      <div class="header-content">
        <h1 class="page-title">选课通——为软工人打造的评价选课平台</h1>
      </div>
    </header>

    <!-- 求评价中心入口 -->
    <section class="evaluation-request-entry-section">
      <div class="entry-content">
        <div class="entry-icon">📢</div>
        <div class="entry-text">
          <h3 class="entry-title">发现需要评价的课程</h3>
          <p class="entry-description">帮助同学们找到优质课程，分享你的学习体验</p>
        </div>
        <RouterLink to="/evaluation-requests" class="entry-button">
          <span class="btn-icon">🔍</span>
          求评价中心
        </RouterLink>
      </div>
    </section>

    <!-- Filter Section -->
    <section class="filter-section">
      <form @submit.prevent="fetchCourses" class="filter-form">
        <div class="filter-header">
          <h2 class="filter-title">筛选课程</h2>
          <p class="filter-subtitle">按年级、学期、科目快速找到您感兴趣的课程</p>
        </div>
        <div class="filter-inputs">
          <div class="form-group">
            <label for="grade" class="form-label">
              年级
            </label>
            <input
              type="text"
              id="grade"
              v-model="filters.grade"
              class="input-glass"
              placeholder="输入年级..."
            />
          </div>
          <div class="form-group">
            <label for="semester" class="form-label">
              学期
            </label>
            <input
              type="text"
              id="semester"
              v-model="filters.semester"
              class="input-glass"
              placeholder="输入学期..."
            />
          </div>
          <div class="form-group">
            <label for="subject" class="form-label">
              科目
            </label>
            <input
              type="text"
              id="subject"
              v-model="filters.subject"
              class="input-glass"
              placeholder="输入科目..."
            />
          </div>
          <div class="form-group form-group--button">
            <button type="submit" class="btn btn-primary">
              搜索课程
            </button>
          </div>
        </div>
      </form>
    </section>

    <!-- Loading State -->
    <div v-if="loading" class="loading-container">
      <div class="loader"></div>
      <p class="loading-text">正在加载课程...</p>
    </div>

    <!-- Course Grid -->
    <section ref="courseSection" v-else class="course-section">
      <div class="course-header">
        <div class="course-header-content">
          <h2 class="course-title">热门课程评价</h2>
          <p class="course-description">查看真实学生评价，找到最适合您的课程</p>
        </div>
        <div class="course-filters">
          <button class="filter-chip" :class="{ active: !filters.grade && !filters.semester && !filters.subject }" @click="clearFilters">
            全部课程
          </button>
          <button class="filter-chip" @click="showPopularOnly">
            高分课程
          </button>
          <button class="filter-chip" @click="showNewOnly">
            最新评价
          </button>
        </div>
      </div>
      
      <div v-if="courses.length === 0" class="empty-state">
        <h3 class="empty-title">暂无课程</h3>
        <p class="empty-description">请调整筛选条件或稍后再试</p>
        <button @click="clearFilters" class="btn btn-primary mt-lg">
          重置筛选
        </button>
      </div>
      
      <div v-else class="course-grid">
        <router-link
          v-for="(course, index) in courses"
          :key="course.ID"
          :to="`/courses/${course.ID}`"
          class="course-card"
        >
          <!-- 课程内容 -->
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
                      <div class="bar-fill" :style="{ width: getRatingPercentage(getDisplayRating(course), 6-i) + '%' }"></div>
                    </div>
                    <span class="bar-count">{{ getRatingCount(getDisplayRating(course), 6-i) }}</span>
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
    </section>

    <!-- Footer -->
    <footer class="home-footer">
      <div class="footer-content">
        <div class="footer-brand">
          <div class="footer-logo">选课通</div>
          <p class="footer-tagline">让学习更简单，让选择更明智</p>
        </div>
        <div class="footer-links">
          <div class="footer-section">
            <h4 class="footer-section-title">产品</h4>
            <router-link to="/" class="footer-link">首页</router-link>
            <router-link to="/about" class="footer-link">关于我们</router-link>
            <router-link to="/profile" class="footer-link">个人中心</router-link>
          </div>
          <div class="footer-section">
            <h4 class="footer-section-title">支持</h4>
            <a href="#" class="footer-link">帮助中心</a>
            <a href="#" class="footer-link">联系我们</a>
            <a href="#" class="footer-link">意见反馈</a>
          </div>
          <div class="footer-section">
            <h4 class="footer-section-title">关注我们</h4>
            <div class="social-links">
              <a href="#" class="social-link">邮箱</a>
              <a href="#" class="social-link">微信</a>
              <a href="#" class="social-link">电话</a>
            </div>
          </div>
        </div>
      </div>
      <div class="footer-bottom">
        <p class="footer-text">© 2024 选课通 - 让学习更简单</p>
      </div>
    </footer>
  </main>
</template>

<style scoped>
/* ===== 俏皮的新拟物风格 (Playful Neobrutalism) ===== */

/* 全局样式 */
.home-container {
  font-family: sans-serif;
  color: #1A1A1A;
  padding: 20px;
  max-width: 420px;
  margin: 0 auto;
  background-color: #FEF6F7;
  min-height: 100vh;
}

/* 增大间距 */
.course-grid {
  display: grid;
  gap: 24px;
  margin-bottom: 32px;
}

.filter-form {
  background-color: #FFFFFF;
  border-radius: 12px;
  border: 3px solid #000000;
  box-shadow: 5px 5px 0px 0px #000000;
  padding: 24px;
  display: grid;
  grid-template-columns: 1fr;
  gap: 20px;
}

/* 求评价中心入口 */
.evaluation-request-entry-section {
  margin-bottom: 24px;
}

.entry-content {
  background-color: #FFFFFF;
  border-radius: 12px;
  border: 3px solid #000000;
  box-shadow: 5px 5px 0px 0px #000000;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  transition: transform 0.2s ease;
}

.entry-content:hover {
  transform: translate(-2px, -2px);
  box-shadow: 7px 7px 0px 0px #000000;
}

.entry-icon {
  font-size: 32px;
  flex-shrink: 0;
}

.entry-text {
  flex-grow: 1;
}

.entry-title {
  font-size: 20px;
  font-weight: bold;
  color: #1A1A1A;
  margin-bottom: 4px;
}

.entry-description {
  font-size: 16px;
  color: #888888;
  margin: 0;
}

.entry-button {
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
  flex-shrink: 0;
}

.entry-button:hover {
  transform: translate(-2px, -2px);
  box-shadow: 6px 6px 0px 0px #000000;
}

/* Header Section 样式 */
.home-header {
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
  margin: 0 0 16px 0;
}

.page-badge {
  display: inline-block;
  background-color: #F7D074;
  color: #1A1A1A;
  padding: 6px 12px;
  border-radius: 8px;
  border: 2px solid #000000;
  font-size: 14px;
  font-weight: bold;
  margin-bottom: 16px;
}

.header-stats {
  display: flex;
  justify-content: space-around;
  margin-top: 24px;
}

.header-stats .stat-item {
  text-align: center;
}

.header-stats .stat-number {
  font-size: 24px;
  font-weight: bold;
  color: #F7D074;
  display: block;
}

.header-stats .stat-label {
  font-size: 14px;
  color: #888888;
}

/* Filter Section 样式 */
.filter-section {
  margin-bottom: 24px;
}

.filter-header {
  background-color: #FFFFFF;
  border-radius: 12px;
  border: 3px solid #000000;
  box-shadow: 5px 5px 0px 0px #000000;
  padding: 24px;
  margin-bottom: 16px;
  text-align: center;
}

.filter-title {
  font-size: 20px;
  font-weight: bold;
  color: #1A1A1A;
  margin-bottom: 8px;
}

.filter-subtitle {
  font-size: 14px;
  color: #888888;
  margin: 0;
}

.filter-form {
  background-color: #FFFFFF;
  border-radius: 12px;
  border: 3px solid #000000;
  box-shadow: 5px 5px 0px 0px #000000;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.filter-header {
  text-align: center;
  margin-bottom: 8px;
}

.filter-inputs {
  display: grid;
  grid-template-columns: 1fr;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-size: 14px;
  font-weight: bold;
  color: #1A1A1A;
}

.input-glass {
  background-color: #FFFFFF;
  border: 3px solid #000000;
  border-radius: 8px;
  padding: 12px;
  font-size: 16px;
  font-family: sans-serif;
  transition: transform 0.2s ease;
}

.input-glass:focus {
  transform: translate(-2px, -2px);
  box-shadow: 4px 4px 0px 0px #000000;
  outline: none;
}

.form-group--button {
  margin-top: 8px;
}

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

/* Course Section 样式 */
.course-section {
  margin-bottom: 24px;
}

.course-header {
  margin-bottom: 16px;
}

.course-header-content {
  background-color: #FFFFFF;
  border-radius: 12px;
  border: 3px solid #000000;
  box-shadow: 5px 5px 0px 0px #000000;
  padding: 24px;
  margin-bottom: 16px;
  text-align: center;
}

.course-title {
  font-size: 20px;
  font-weight: bold;
  color: #1A1A1A;
  margin-bottom: 8px;
}

.course-description {
  font-size: 14px;
  color: #888888;
  margin: 0;
}

.course-filters {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  justify-content: center;
  margin-bottom: 16px;
}

.filter-chip {
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

.filter-chip:hover {
  transform: translate(-2px, -2px);
  box-shadow: 6px 6px 0px 0px #000000;
}

.filter-chip.active {
  background-color: #F7D074;
}

.course-grid {
  display: grid;
  gap: 16px;
  margin-bottom: 24px;
}

/* ===== 新卡片设计 ===== */
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

.course-card-description {
  color: #888888;
  margin-bottom: 16px;
}

.rating-count {
  color: #888888;
}

/* Footer 样式 */
.home-footer {
  margin-top: 24px;
}

.footer-content {
  background-color: #FFFFFF;
  border-radius: 12px;
  border: 3px solid #000000;
  box-shadow: 5px 5px 0px 0px #000000;
  padding: 24px;
  margin-bottom: 16px;
}

.footer-brand {
  text-align: center;
  margin-bottom: 16px;
}

.footer-logo {
  font-size: 24px;
  font-weight: bold;
  color: #F7D074;
}

.footer-tagline {
  font-size: 14px;
  color: #888888;
  margin: 8px 0 0 0;
}

.footer-links {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.footer-section {
  text-align: center;
}

.footer-section-title {
  font-size: 16px;
  font-weight: bold;
  color: #1A1A1A;
  margin-bottom: 8px;
}

.footer-link {
  display: block;
  color: #888888;
  text-decoration: none;
  font-size: 14px;
  margin-bottom: 8px;
  transition: color 0.2s ease;
}

.footer-link:hover {
  color: #F7D074;
}

.social-links {
  display: flex;
  justify-content: center;
  gap: 8px;
}

.social-link {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  background-color: #F7D074;
  border-radius: 8px;
  border: 2px solid #000000;
  color: #1A1A1A;
  text-decoration: none;
  font-size: 14px;
  font-weight: bold;
  transition: transform 0.2s ease;
}

.social-link:hover {
  transform: translate(-2px, -2px);
  box-shadow: 4px 4px 0px 0px #000000;
}

.footer-bottom {
  text-align: center;
  padding-top: 16px;
  border-top: 3px solid #000000;
}

.footer-text {
  font-size: 14px;
  color: #888888;
  margin: 0;
}

/* 加载状态样式 */
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

/* 空状态样式 */
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

.mt-lg {
  margin-top: 16px;
}

/* ===== 课程评价功能样式 ===== */
.course-card-header {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 3px solid #000000;
}

/* 授课老师部分 */
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

.btn-rate-course .btn-icon {
  font-size: 16px;
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

/* 响应式设计 */
@media (min-width: 768px) {
  .home-container {
    max-width: 768px;
    background-color: #FEF6F7;
  }
  
  .filter-inputs {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .course-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 1024px) {
  .home-container {
    max-width: 1400px;
    padding: 40px;
    background-color: #FEF6F7;
  }
  
  .filter-form {
    gap: 32px;
  }
  
  .filter-inputs {
    grid-template-columns: repeat(3, 1fr);
    gap: 24px;
  }
  
  .course-grid {
    grid-template-columns: repeat(3, 1fr);
    gap: 32px;
  }
  
  .header-stats {
    justify-content: center;
    gap: 40px;
  }
  
  .footer-links {
    grid-template-columns: repeat(3, 1fr);
    gap: 32px;
  }
  
  .course-filters {
    justify-content: flex-start;
  }
}

@media (min-width: 1440px) {
  .home-container {
    max-width: 1600px;
    background-color: #FEF6F7;
  }
  
  .course-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}

@media (min-width: 1600px) {
  .home-container {
    max-width: 1800px;
    background-color: #FEF6F7;
  }
  
  .course-grid {
    grid-template-columns: repeat(5, 1fr);
  }
}

@media (min-width: 1920px) {
  .home-container {
    max-width: 2200px;
    background-color: #FEF6F7;
  }
  
  .course-grid {
    grid-template-columns: repeat(6, 1fr);
  }
}


</style>
