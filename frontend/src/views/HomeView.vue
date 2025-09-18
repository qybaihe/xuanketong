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

// 只使用natural主题
const themeClass = computed(() => 'theme-natural')

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
  <main class="home-container" :class="themeClass">
    <!-- Header Section -->
    <header class="home-header">
      <div class="header-content">
        <div class="page-badge">课程评价平台</div>
        <h1 class="page-title">发现精品课程</h1>
        <p class="page-subtitle">查看真实评价，分享学习体验，帮助同学做出明智选择</p>
        <div class="header-stats">
          <div class="stat-item">
            <div class="stat-number">5000+</div>
            <div class="stat-label">课程评价</div>
          </div>
          <div class="stat-item">
            <div class="stat-number">98%</div>
            <div class="stat-label">学生满意度</div>
          </div>
          <div class="stat-item">
            <div class="stat-number">1000+</div>
            <div class="stat-label">优质课程</div>
          </div>
        </div>
      </div>
    </header>

    <!-- Filter Section -->
    <section class="filter-section card-glass">
      <div class="filter-header">
        <h2 class="filter-title">筛选课程</h2>
        <p class="filter-subtitle">按年级、学期、科目快速找到您感兴趣的课程</p>
      </div>
      <form @submit.prevent="fetchCourses" class="filter-form">
        <div class="form-group">
          <label for="grade" class="form-label">
            年级
          </label>
          <input
            type="text"
            id="grade"
            v-model="filters.grade"
            class="input input-glass"
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
            class="input input-glass"
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
            class="input input-glass"
            placeholder="输入科目..."
          />
        </div>
        <div class="form-group form-group--button">
          <button type="submit" class="btn btn-primary">
            搜索课程
          </button>
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
          :class="getCardClass(course, index)"
          class="scroll-reveal"
        >
          <!-- 课程图片 -->
          <!--
          <div class="course-card-image hardware-accelerated">
            <img
              :src="course.ImageURL || `https://picsum.photos/seed/course-${course.ID}/400/300.jpg`"
              :alt="course.Name || course.name || '课程'"
              loading="lazy"
              class="hardware-accelerated"
              @load="handleImageLoad"
              @error="handleImageError"
            />
            <div class="course-card-rating">
              <div class="rating-stars">
                <span v-for="i in 5" :key="i" class="star">
                  {{ i <= Math.floor(getDisplayRating(course)) ? '⭐' : (i - 0.5 <= getDisplayRating(course) ? '🌟' : '☆') }}
                </span>
              </div>
              <span class="rating-value">{{ getDisplayRating(course).toFixed(1) }}</span>
            </div>
          </div>
          -->

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
              <h3 class="course-card-title text-shine">{{ course.Name }}</h3>
              <div class="course-teacher-section">
                <div class="teacher-label">授课老师</div>
                <div class="teacher-name">{{ course.Teacher }}</div>
              </div>
              <div class="course-meta-info">
                <div class="meta-item">
                  <span class="meta-text">{{ course.Credits }} 学分</span>
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
              <div class="engagement-stats">
                <div class="meta-item">
                  <span class="meta-text">{{ course.Credits }} 学分</span>
                </div>
                <div class="meta-item">
                  <span class="meta-text">{{ (course.TotalRatings || 0) }} 人评价</span>
                </div>
              </div>
              <button class="btn btn-rate-course" @click.prevent="goToRateCourse(course.ID)">
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
/* ===== 简洁优雅样式 ===== */

/* Header Section 样式 */
.home-header {
  background: var(--background-primary);
  border-bottom: 1px solid var(--separator-color);
  padding: var(--spacing-2xl) 0;
  position: relative;
  overflow: hidden;
}

.home-header::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, 
    transparent 0%, 
    var(--primary-color) 50%, 
    transparent 100%);
  animation: headerShine 3s ease-in-out infinite;
}

@keyframes headerShine {
  0%, 100% { opacity: 0; transform: translateX(-100%); }
  50% { opacity: 1; transform: translateX(100%); }
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 var(--spacing-lg);
  text-align: center;
}

.page-title {
  font-size: clamp(32px, 4vw, 40px);
  font-weight: var(--font-weight-light);
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
  width: 40px;
  height: 2px;
  background: var(--gradient-primary);
  border-radius: 1px;
  transition: width 0.3s ease;
}

.page-title:hover::after {
  width: 60px;
}

.page-subtitle {
  font-size: clamp(16px, 1.5vw, 18px);
  color: var(--text-secondary);
  margin: 0;
  font-weight: var(--font-weight-regular);
  letter-spacing: 0.5px;
  opacity: 0.9;
}

.page-badge {
  display: inline-block;
  background: var(--gradient-primary);
  color: white;
  padding: var(--spacing-xs) var(--spacing-md);
  border-radius: 20px;
  font-size: var(--font-size-caption);
  font-weight: var(--font-weight-medium);
  margin-bottom: var(--spacing-lg);
  animation: badgeFloat 3s ease-in-out infinite;
}

@keyframes badgeFloat {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-3px); }
}

.header-stats {
  display: flex;
  justify-content: center;
  gap: var(--spacing-2xl);
  margin-top: var(--spacing-xl);
  flex-wrap: wrap;
}

.header-stats .stat-item {
  text-align: center;
  min-width: 80px;
}

.header-stats .stat-number {
  font-size: 32px;
  font-weight: var(--font-weight-bold);
  color: var(--primary-color);
  margin-bottom: var(--spacing-xs);
  display: block;
}

.header-stats .stat-label {
  font-size: var(--font-size-body2);
  color: var(--text-secondary);
  font-weight: var(--font-weight-medium);
}

/* Filter Section 增强样式 */
.filter-section {
  max-width: 1200px;
  margin: var(--spacing-3xl) auto;
  padding: var(--spacing-xl);
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 
    0 20px 40px rgba(0, 0, 0, 0.1),
    0 0 0 1px rgba(255, 255, 255, 0.05);
}

.filter-header {
  text-align: center;
  margin-bottom: var(--spacing-xl);
}

.filter-title {
  font-size: 32px;
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin-bottom: var(--spacing-sm);
}

.filter-subtitle {
  font-size: var(--font-size-body);
  color: var(--text-secondary);
  margin: 0;
}

.filter-form {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: var(--spacing-lg);
  align-items: end;
}

.form-label {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  font-size: var(--font-size-body2);
  font-weight: var(--font-weight-medium);
  color: var(--text-secondary);
  margin-bottom: var(--spacing-sm);
}

.label-icon {
  font-size: 16px;
}

.input-glass {
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  border: 2px solid rgba(255, 255, 255, 0.3);
  padding: 14px 16px;
  font-size: var(--font-size-body);
  border-radius: 12px;
  transition: var(--transition-standard);
}

.input-glass:focus {
  background: rgba(255, 255, 255, 0.95);
  border-color: var(--primary-color);
  box-shadow: 0 0 0 4px rgba(47, 169, 20, 0.1);
  transform: translateY(-2px);
}

/* Course Section 增强样式 */
.course-section {
  max-width: 1400px;
  margin: 0 auto;
  padding: var(--spacing-xl) var(--spacing-md);
}

.course-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-2xl);
  flex-wrap: wrap;
  gap: var(--spacing-lg);
}

.course-header-content {
  flex: 1;
  min-width: 300px;
}

.course-title {
  font-size: 36px;
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin-bottom: var(--spacing-sm);
}

.course-description {
  font-size: var(--font-size-body);
  color: var(--text-secondary);
  margin: 0;
}

.course-filters {
  display: flex;
  gap: var(--spacing-sm);
  flex-wrap: wrap;
}

.filter-chip {
  padding: 8px 16px;
  border-radius: 20px;
  background: var(--background-secondary);
  color: var(--text-secondary);
  border: 2px solid transparent;
  font-size: var(--font-size-body2);
  font-weight: var(--font-weight-medium);
  cursor: pointer;
  transition: var(--transition-standard);
}

.filter-chip:hover {
  background: var(--background-tertiary);
  transform: translateY(-2px);
}

.filter-chip.active {
  background: var(--gradient-primary);
  color: white;
  border-color: var(--primary-color);
}

.course-grid {
  display: grid;
  gap: var(--spacing-2xl);
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  margin-bottom: var(--spacing-3xl);
}

/* ===== 新卡片设计 ===== */
.course-card {
  background: #fff;
  border-radius: 20px;
  overflow: hidden;
  text-decoration: none;
  color: var(--text-primary);
  display: flex;
  flex-direction: column;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  transition: all 0.3s cubic-bezier(0.25, 0.46, 0.45, 0.94);
  min-height: 500px;
  position: relative;
  border: 2px solid transparent;
  background-clip: padding-box;
}

.course-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  border-radius: 20px;
  padding: 2px;
  background: var(--gradient-primary);
  -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
  mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
  mask-composite: exclude;
  -webkit-mask-composite: exclude;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.course-card:hover {
  transform: translateY(-8px) scale(1.02);
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.18);
  border-color: var(--primary-color);
}

.course-card:hover::before {
  opacity: 1;
}

.course-card-content {
  padding: 28px;
  display: flex;
  flex-direction: column;
  flex-grow: 1;
}

.course-card-tags {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
}

.course-card-tag {
  padding: 6px 12px;
  border-radius: 14px;
  font-size: 13px;
  font-weight: 500;
}

.course-card-tag.primary {
  background-color: #eef4ff;
  color: #4d8dff;
}

.course-card-tag.secondary {
  background-color: #f0f9f4;
  color: #28a745;
}

.course-card-tag.accent {
  background-color: #fff8e1;
  color: #f59e0b;
}

.course-card-title {
  color: var(--text-primary);
  font-size: 22px;
  font-weight: 700;
  margin-bottom: 12px;
  line-height: 1.3;
  position: relative;
  transition: all 0.3s ease;
}

.course-card-title::after {
  content: '';
  position: absolute;
  bottom: -4px;
  left: 0;
  width: 0;
  height: 3px;
  background: var(--gradient-primary);
  border-radius: 2px;
  transition: width 0.3s ease;
}

.course-card:hover .course-card-title {
  color: var(--primary-color);
}

.course-card:hover .course-card-title::after {
  width: 40px;
}

.course-card-description,
.rating-count {
  color: var(--text-secondary); /* 将描述和评价数量文字颜色设置为次要文字颜色 */
}

/* Footer 增强样式 */
.home-footer {
  background: var(--background-primary);
  border-top: 1px solid var(--separator-color);
  padding: var(--spacing-3xl) 0;
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 var(--spacing-md);
  display: grid;
  grid-template-columns: 2fr 3fr;
  gap: var(--spacing-2xl);
  margin-bottom: var(--spacing-2xl);
}

.footer-brand {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.footer-logo {
  font-size: 32px;
  font-weight: var(--font-weight-bold);
  color: var(--primary-color);
}

.footer-tagline {
  font-size: var(--font-size-body2);
  color: var(--text-secondary);
  margin: 0;
}

.footer-links {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--spacing-xl);
}

.footer-section {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.footer-section-title {
  font-size: var(--font-size-body);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin: 0;
}

.footer-link {
  color: var(--text-secondary);
  text-decoration: none;
  font-size: var(--font-size-body2);
  transition: var(--transition-standard);
}

.footer-link:hover {
  color: var(--primary-color);
  transform: translateX(4px);
}

.social-links {
  display: flex;
  gap: var(--spacing-sm);
}

.social-link {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  background: var(--background-secondary);
  border-radius: 50%;
  color: var(--text-secondary);
  text-decoration: none;
  transition: var(--transition-standard);
}

.social-link:hover {
  background: var(--primary-color);
  color: white;
  transform: translateY(-2px);
}

.footer-bottom {
  text-align: center;
  padding-top: var(--spacing-lg);
  border-top: 1px solid var(--separator-color);
  max-width: 1200px;
  margin: 0 auto;
  padding: var(--spacing-lg) var(--spacing-md);
}

.footer-text {
  font-size: var(--font-size-body2);
  color: var(--text-tertiary);
  margin: 0;
}

/* Course Card Enhancements */
/*
.course-card-image {
  position: relative;
  aspect-ratio: 16 / 10;
  width: 100%;
}

.course-card-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.4s ease;
}

.course-card:hover .course-card-image img {
  transform: scale(1.05);
}

.course-card-rating {
  position: absolute;
  top: 12px;
  right: 12px;
  background: rgba(0, 0, 0, 0.6);
  color: white;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 6px;
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.rating-stars {
  font-size: 12px;
}

.meta-info {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.meta-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  font-size: var(--font-size-caption);
  color: var(--text-tertiary);
}

.meta-icon {
  font-size: 14px;
}
*/

/* Footer */
.home-footer {
  text-align: center;
  padding: var(--spacing-2xl) var(--spacing-lg);
  margin-top: var(--spacing-3xl);
}

.footer-text {
  font-size: var(--font-size-body2);
  color: var(--text-tertiary);
}

/* ===== 响应式设计增强 ===== */

/* ===== 响应式设计增强 ===== */

/* 超大屏幕适配 (2K+ 分辨率) */
@media (min-width: 1600px) {
  .home-header {
    padding: var(--spacing-3xl) 0;
  }
  
  .page-title {
    font-size: 44px;
  }
  
  .page-subtitle {
    font-size: 20px;
  }
  
  .filter-section {
    max-width: 1400px;
  }
  
  .filter-form {
    grid-template-columns: repeat(4, 1fr);
  }
  
  .course-section {
    max-width: 1800px;
  }
  
  .course-grid {
    grid-template-columns: repeat(3, 1fr);
    gap: var(--spacing-3xl);
  }
  
  .course-header {
    align-items: flex-start;
  }
  
  .course-filters {
    margin-top: var(--spacing-sm);
  }
}

/* 大屏幕适配 (PC 端主要尺寸) */
@media (min-width: 1200px) {
  .home-header {
    padding: var(--spacing-2xl) 0;
  }
  
  .page-title {
    font-size: 40px;
  }
  
  .page-subtitle {
    font-size: 18px;
  }
  
  .filter-section {
    max-width: 1200px;
  }
  
  .filter-form {
    grid-template-columns: repeat(4, 1fr);
    gap: var(--spacing-xl);
  }
  
  .course-section {
    max-width: 1400px;
  }
  
  .course-grid {
    grid-template-columns: repeat(3, 1fr);
    gap: var(--spacing-xl);
  }
  
  .course-card {
    max-width: none;
  }
  
  .footer-content {
    grid-template-columns: 3fr 4fr;
    gap: var(--spacing-3xl);
  }
  
  .footer-links {
    grid-template-columns: repeat(3, 1fr);
    gap: var(--spacing-2xl);
  }
}

/* 中等屏幕适配 (平板和小型桌面) */
@media (min-width: 769px) and (max-width: 1199px) {
  .home-header {
    padding: var(--spacing-xl) 0;
  }
  
  .page-title {
    font-size: 36px;
  }
  
  .page-subtitle {
    font-size: 17px;
  }
  
  .filter-form {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .course-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: var(--spacing-lg);
  }
  
  .course-header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .course-filters {
    margin-top: var(--spacing-md);
  }
  
  .footer-content {
    grid-template-columns: 1fr;
    gap: var(--spacing-2xl);
  }
  
  .footer-links {
    grid-template-columns: repeat(3, 1fr);
  }
}

/* 移动端适配 */
@media (max-width: 768px) {
  .home-header {
    padding: var(--spacing-lg) 0;
  }
  
  .page-title {
    font-size: 32px;
    margin-bottom: var(--spacing-xs);
  }
  
  .page-subtitle {
    font-size: 16px;
  }
  
  .filter-section {
    margin: var(--spacing-xl) var(--spacing-md);
    padding: var(--spacing-lg);
  }
  
  .filter-header {
    margin-bottom: var(--spacing-lg);
  }
  
  .filter-title {
    font-size: 24px;
  }
  
  .filter-form {
    grid-template-columns: 1fr;
    gap: var(--spacing-md);
  }
  
  .course-section {
    padding: var(--spacing-lg) var(--spacing-sm);
  }
  
  .course-header {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-md);
  }
  
  .course-title {
    font-size: 28px;
  }
  
  .course-filters {
    width: 100%;
    overflow-x: auto;
    padding-bottom: var(--spacing-sm);
  }
  
  .filter-chip {
    flex-shrink: 0;
  }
  
  .course-grid {
    grid-template-columns: 1fr;
    gap: var(--spacing-lg);
  }
  
  .footer-content {
    grid-template-columns: 1fr;
    gap: var(--spacing-xl);
  }
  
  .footer-links {
    grid-template-columns: repeat(2, 1fr);
    gap: var(--spacing-lg);
  }
  
  .footer-brand {
    text-align: center;
  }
  
  .social-links {
    justify-content: center;
  }
}

/* 小型移动端适配 */
@media (max-width: 480px) {
  .home-header {
    padding: var(--spacing-md) 0;
  }
  
  .page-title {
    font-size: 28px;
    margin-bottom: var(--spacing-xs);
  }
  
  .page-subtitle {
    font-size: 15px;
  }
  
  .page-title::after {
    width: 30px;
  }
  
  .page-title:hover::after {
    width: 40px;
  }
  
  .filter-section {
    margin: var(--spacing-lg) var(--spacing-sm);
    padding: var(--spacing-md);
  }
  
  .filter-title {
    font-size: 20px;
  }
  
  .course-header {
    gap: var(--spacing-sm);
  }
  
  .course-title {
    font-size: 24px;
  }
  
  .course-filters {
    gap: var(--spacing-xs);
  }
  
  .filter-chip {
    padding: 6px 12px;
    font-size: 12px;
  }
  
  .footer-content {
    padding: 0 var(--spacing-sm);
  }
  
  .footer-links {
    grid-template-columns: 1fr;
    gap: var(--spacing-lg);
  }
}

/* 横屏模式优化 */
@media (orientation: landscape) and (max-height: 600px) {
  .home-header {
    padding: var(--spacing-lg) 0;
  }
}

/* 高对比度模式支持 */
@media (prefers-contrast: high) {
  .home-header {
    border-bottom: 2px solid var(--text-primary);
  }
  
  .page-title::after {
    background: var(--text-primary);
    height: 3px;
  }
  
  .filter-section {
    background: var(--background-primary);
    border: 2px solid var(--text-primary);
  }
  
  .course-card {
    border: 2px solid var(--text-primary);
  }
  
  .footer {
    border-top: 2px solid var(--text-primary);
  }
}

/* 课程徽章样式 */
.course-badge {
  position: absolute;
  top: var(--spacing-sm);
  left: var(--spacing-sm);
  background: rgba(0, 0, 0, 0.8);
  color: white;
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: 16px;
  font-size: var(--font-size-caption);
  font-weight: var(--font-weight-medium);
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  z-index: 2;
}

.course-badge-popular {
  background: linear-gradient(135deg, #ff6b6b, #ff8e8e);
  box-shadow: 0 2px 8px rgba(255, 107, 107, 0.3);
}

.course-badge-new {
  background: linear-gradient(135deg, #4ecdc4, #44a08d);
  box-shadow: 0 2px 8px rgba(78, 205, 196, 0.3);
}

.badge-icon {
  font-size: 12px;
}

.badge-text {
  font-weight: var(--font-weight-semibold);
}

/* 评分星级样式 */
.rating-stars {
  display: flex;
  gap: 2px;
}

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
  color: #ffc107;
}

.star-half::before {
  content: "★";
  color: #ffc107;
  background: linear-gradient(90deg, #ffc107 50%, #ccc 50%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.star-empty::before {
  content: "★";
  color: #ccc;
}

.course-card:hover .star {
  transform: scale(1.1);
}

/* 动画优化 */
.course-card {
  animation-delay: calc(var(--index, 0) * 0.1s);
  animation: cardEntry 0.6s cubic-bezier(0.25, 0.46, 0.45, 0.94) both;
}

@keyframes cardEntry {
  0% {
    opacity: 0;
    transform: translateY(30px) scale(0.95);
  }
  50% {
    opacity: 0.8;
    transform: translateY(-5px) scale(1.02);
  }
  100% {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.course-card:nth-child(1) { --index: 1; }
.course-card:nth-child(2) { --index: 2; }
.course-card:nth-child(3) { --index: 3; }
.course-card:nth-child(4) { --index: 4; }
.course-card:nth-child(5) { --index: 5; }
.course-card:nth-child(6) { --index: 6; }
.course-card:nth-child(7) { --index: 7; }
.course-card:nth-child(8) { --index: 8; }
.course-card:nth-child(9) { --index: 9; }
.course-card:nth-child(10) { --index: 10; }
.course-card:nth-child(11) { --index: 11; }
.course-card:nth-child(12) { --index: 12; }


/* 加载动画优化 */
.loader {
  border: 3px solid var(--background-secondary);
  border-top: 3px solid var(--primary-color);
}

.loading-text {
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  font-weight: var(--font-weight-semibold);
}

/* 空状态动画 */
.empty-state {
  animation: fadeInUp 0.6s ease-out;
}

.empty-icon {
  animation: float 3s ease-in-out infinite;
}

/* 触摸优化 */
@media (hover: none) {
  .course-card:hover {
    transform: none;
  }
  
  .course-card:active {
    transform: scale(0.98);
  }
  
  .btn:hover {
    transform: none;
  }
  
  .magnetic {
    transform: none !important;
  }
}

/* 高对比度模式支持 */
@media (prefers-contrast: high) {
  .course-card-glass {
    background: var(--background-primary);
    border: 2px solid var(--text-primary);
  }
  
  .glass-enhanced {
    background: var(--background-primary);
    border: 2px solid var(--text-primary);
  }
}

/* 减少动画支持 */
@media (prefers-reduced-motion: reduce) {
  .course-card,
  .hero-title,
  .hero-subtitle,
  .empty-state {
    animation: none;
    transition: none;
  }
  
  .float,
  .breathe,
  .sway {
    animation: none;
  }
}

/* 打印样式 */
@media print {
  .home-container {
    background: white;
  }
  
  .hero-section,
  .filter-section,
  .home-footer {
    display: none;
  }
  
  .course-card {
    break-inside: avoid;
    box-shadow: none;
    border: 1px solid #ccc;
  }
  
  .course-badge {
    background: black !important;
    color: white !important;
  }
}

/* 性能优化 */
.course-card {
  /* 移除contain属性，因为它可能影响动画渲染 */
  /* contain: layout style paint; */
}

.course-card-image img {
  content-visibility: auto;
}

/* 自定义滚动条 */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: var(--background-secondary);
}

::-webkit-scrollbar-thumb {
  background: var(--primary-color);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: var(--primary-color-dark);
}

/* 焦点样式优化 */
.btn:focus,
.input:focus {
  outline: 2px solid var(--primary-color);
  outline-offset: 2px;
}

/* 键盘导航支持 */
.btn:focus-visible,
.input:focus-visible {
  outline: 2px solid var(--primary-color);
  outline-offset: 2px;
}

/* 暗色主题支持 */
@media (prefers-color-scheme: dark) {
  .home-container {
    background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  }
  
  .hero-title {
    background: linear-gradient(135deg, #4fc830 0%, #2fa914 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }
  
  .course-card-glass {
    background: rgba(0, 0, 0, 0.6);
    border: 1px solid rgba(255, 255, 255, 0.1);
  }
  
  .course-card-content {
    color: white;
  }
  
  .course-card-title {
    color: white;
  }
  
  .course-card-description {
    color: rgba(255, 255, 255, 0.8);
  }
}

/* ===== 课程评价功能样式 ===== */
.course-card-header {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}

/* 授课老师部分 */
.course-teacher-section {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: linear-gradient(135deg, rgba(5, 170, 105, 0.08) 0%, rgba(5, 170, 105, 0.04) 100%);
  border-radius: 12px;
  border: 1px solid rgba(5, 170, 105, 0.15);
  transition: all 0.3s ease;
}

.course-teacher-section:hover {
  background: linear-gradient(135deg, rgba(5, 170, 105, 0.12) 0%, rgba(5, 170, 105, 0.06) 100%);
  border-color: rgba(5, 170, 105, 0.25);
  transform: translateY(-1px);
}

.teacher-label {
  font-size: 12px;
  color: var(--primary-color);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  opacity: 0.8;
}

.teacher-name {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: 8px;
}


/* 课程元信息 */
.course-meta-info {
  display: flex;
  gap: var(--spacing-lg);
  flex-wrap: wrap;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: var(--text-secondary);
  padding: 6px 12px;
  background: var(--background-secondary);
  border-radius: 8px;
  transition: all 0.2s ease;
}

.meta-item:hover {
  background: var(--background-tertiary);
  transform: translateY(-1px);
}

.meta-text {
  font-weight: var(--font-weight-medium);
}

.course-rating-section {
  padding: 20px;
  background: #f7f8fa;
  border-radius: 16px;
  border: 1px solid #e9eaee;
  margin-bottom: 20px;
}

.rating-overview {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.rating-average {
  text-align: center;
}

.rating-number {
  font-size: 36px;
  font-weight: var(--font-weight-bold);
  color: var(--primary-color);
  line-height: 1;
  display: block;
}

.rating-stars {
  display: flex;
  gap: 4px;
  margin: var(--spacing-sm) 0;
}

.rating-count {
  font-size: var(--font-size-caption);
  color: var(--text-tertiary);
  white-space: nowrap;
}

.rating-breakdown {
  width: 100%;
}

.rating-bar {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-sm);
}

.bar-label {
  font-size: var(--font-size-caption);
  color: var(--text-secondary);
  min-width: 30px;
  text-align: right;
}

.bar-container {
  flex: 1;
  height: 8px;
  background: var(--separator-color);
  border-radius: 4px;
  overflow: hidden;
  position: relative;
}

.bar-fill {
  height: 100%;
  background: var(--gradient-primary);
  border-radius: 3px;
  transition: width 0.3s ease;
}

.bar-count {
  font-size: var(--font-size-caption);
  color: var(--text-tertiary);
  text-align: left;
}

.course-card-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: auto; /* Push to bottom */
  padding-top: 16px;
}


.btn-rate-course {
  padding: 14px 24px;
  background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
  color: white;
  border: none;
  border-radius: 22px;
  font-size: 15px;
  font-weight: var(--font-weight-semibold);
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.25, 0.46, 0.45, 0.94);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-sm);
  box-shadow: 0 6px 16px rgba(40, 167, 69, 0.25);
  position: relative;
  overflow: hidden;
}

.btn-rate-course::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.6s ease;
}

.btn-rate-course:hover::before {
  left: 100%;
}

.btn-rate-course:hover {
  transform: translateY(-3px) scale(1.05);
  background: linear-gradient(135deg, #218838 0%, #1ebd8d 100%);
  box-shadow: 0 8px 24px rgba(40, 167, 69, 0.35);
}

.btn-rate-course:active {
  transform: translateY(-1px) scale(1.02);
  box-shadow: 0 4px 12px rgba(40, 167, 69, 0.3);
}

.btn-rate-course .btn-icon {
  font-size: 16px;
}
</style>
