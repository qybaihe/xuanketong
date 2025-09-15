<script setup lang="ts">
import { ref, onMounted, reactive, watch, computed, nextTick } from 'vue'
import axios from 'axios'
import { RouterLink } from 'vue-router'

interface Course {
  ID: number;
  Name: string;
  Description: string;
  Teacher: string;
  Credits: number;
  ImageURL: string;
  Grade: string;
  Semester: string;
  Subject: string;
  Rating?: number;
  Students?: number;
  isPopular?: boolean;
  isNew?: boolean;
}

const courses = ref<Course[]>([])
const loading = ref(true)
const filters = reactive({
  grade: '',
  semester: '',
  subject: ''
})

// 主题状态
const currentTheme = ref('natural')
const themes = ['discreet', 'matching', 'natural'] as const
type Theme = typeof themes[number]

// 切换主题
const cycleTheme = () => {
  const currentIndex = themes.indexOf(currentTheme.value as Theme)
  currentTheme.value = themes[(currentIndex + 1) % themes.length]
  
  // 添加切换动画效果
  const container = document.querySelector('.home-container')
  if (container) {
    container.classList.add('elastic-in')
    setTimeout(() => {
      container.classList.remove('elastic-in')
    }, 600)
  }
}

// 获取主题类名
const themeClass = computed(() => `theme-${currentTheme.value}`)

// 获取卡片类名
const getCardClass = (course: Course, index: number) => {
  const baseClasses = 'course-card course-card-glass card-float shine-effect hardware-accelerated'
  const themeClass = `course-card-${currentTheme.value}`
  // 移除slide-in-up和slide-in-left动画，避免与scroll-reveal冲突
  // const animationClass = index % 2 === 0 ? 'slide-in-up' : 'slide-in-left'
  return `${baseClasses} ${themeClass}`
}

// 获取标签
const getTags = (course: Course) => {
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

const fetchCourses = async () => {
  console.log('fetchCourses: 开始获取课程数据，loading状态:', loading.value)
  loading.value = true
  try {
    const response = await axios.get('http://localhost:8080/api/v1/courses', { params: filters })
    console.log('fetchCourses: 获取到课程数据，数量:', response.data.data.length)
    courses.value = response.data.data.map((course: Course, index: number) => ({
      ...course,
      Rating: Math.random() * 5, // 模拟评分数据
      Students: Math.floor(Math.random() * 200) + 10, // 模拟学生数
      // 添加热门标记
      isPopular: Math.random() > 0.7,
      // 添加新课程标记
      isNew: Math.random() > 0.8
    }))
    
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
    <!-- Hero Section -->
    <section class="hero-section">
      <div class="hero-content">
        <h1 class="hero-title">课程浏览</h1>
        <p class="hero-subtitle">发现优质课程，提升学习体验</p>
        <div class="hero-actions">
          <button @click="cycleTheme" class="btn btn-glass">
            <span class="theme-icon">🎨</span>
            切换主题: {{ currentTheme }}
          </button>
        </div>
      </div>
    </section>

    <!-- Filter Section -->
    <section class="filter-section card-glass">
      <form @submit.prevent="fetchCourses" class="filter-form">
        <div class="form-group">
          <label for="grade" class="form-label">年级</label>
          <input
            type="text"
            id="grade"
            v-model="filters.grade"
            class="input input-glass"
            placeholder="输入年级..."
          />
        </div>
        <div class="form-group">
          <label for="semester" class="form-label">学期</label>
          <input
            type="text"
            id="semester"
            v-model="filters.semester"
            class="input input-glass"
            placeholder="输入学期..."
          />
        </div>
        <div class="form-group">
          <label for="subject" class="form-label">科目</label>
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
            <span class="btn-icon">🔍</span>
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
    <section v-else class="course-section">
      <div v-if="courses.length === 0" class="empty-state">
        <div class="empty-icon">📚</div>
        <h3 class="empty-title">暂无课程</h3>
        <p class="empty-description">请调整筛选条件或稍后再试</p>
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
          <div class="course-card-image hardware-accelerated">
            <img
              :src="course.ImageURL || `https://picsum.photos/seed/course-${course.ID}/400/200.jpg`"
              :alt="course.Name"
              loading="lazy"
              class="hardware-accelerated"
              @load="handleImageLoad"
            />
            <!-- 热门/新课程标记 -->
            <div v-if="course.isPopular" class="course-badge course-badge-popular">
              <span class="badge-icon">🔥</span>
              <span class="badge-text">热门</span>
            </div>
            <div v-if="course.isNew" class="course-badge course-badge-new">
              <span class="badge-icon">✨</span>
              <span class="badge-text">新课</span>
            </div>
            <!-- 评分显示 -->
            <div class="course-card-rating">
              <div class="rating-stars">
                <span v-for="i in 5" :key="i" class="star">
                  {{ i <= Math.floor(course.Rating || 0) ? '⭐' : (i - 0.5 <= (course.Rating || 0) ? '🌟' : '☆') }}
                </span>
              </div>
              <span class="rating-value">{{ course.Rating?.toFixed(1) || '0.0' }}</span>
            </div>
          </div>

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

            <!-- 课程标题 -->
            <h3 class="course-card-title text-shine">{{ course.Name }}</h3>
            
            <!-- 课程描述 -->
            <p class="course-card-description">{{ course.Description }}</p>
            
            <!-- 课程元信息 -->
            <div class="course-card-meta">
              <div class="meta-info">
                <div class="meta-item">
                  <span class="meta-icon">👨‍🏫</span>
                  <span class="meta-text">{{ course.Teacher }}</span>
                </div>
                <div class="meta-item">
                  <span class="meta-icon">👥</span>
                  <span class="meta-text">{{ course.Students }} 人</span>
                </div>
                <div v-if="course.Students && course.Students > 100" class="meta-item">
                  <span class="meta-icon">📈</span>
                  <span class="meta-text">热门</span>
                </div>
              </div>
              <div class="course-card-credits magnetic">
                {{ course.Credits }} 学分
              </div>
            </div>
          </div>
        </router-link>
      </div>
    </section>

    <!-- Footer -->
    <footer class="home-footer">
      <p class="footer-text">© 2024 选课通 - 让学习更简单</p>
    </footer>
  </main>
</template>

<style scoped>
/* 主容器 */
.home-container {
  min-height: 100vh;
  background: linear-gradient(135deg, var(--background-secondary) 0%, var(--natural-background) 100%);
  padding: var(--spacing-lg) 0;
}

/* Hero Section */
.hero-section {
  text-align: center;
  margin-bottom: var(--spacing-2xl);
  padding: var(--spacing-lg);
  position: relative;
}

.hero-content {
  max-width: 800px;
  margin: 0 auto;
}

.hero-title {
  font-size: var(--font-size-title1);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin-bottom: var(--spacing-sm);
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.hero-subtitle {
  font-size: var(--font-size-body);
  color: var(--text-secondary);
  margin-bottom: var(--spacing-lg);
}

.hero-actions {
  display: flex;
  justify-content: center;
  gap: var(--spacing-md);
}

.theme-icon {
  margin-right: var(--spacing-xs);
}

/* Filter Section */
.filter-section {
  max-width: 1000px;
  margin: 0 auto var(--spacing-xl);
  padding: var(--spacing-lg);
}

.filter-form {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: var(--spacing-md);
  align-items: end;
}

.form-group {
  display: flex;
  flex-direction: column;
}

.form-label {
  font-size: var(--font-size-body2);
  font-weight: var(--font-weight-medium);
  color: var(--text-secondary);
  margin-bottom: var(--spacing-xs);
}

.form-group--button {
  align-self: end;
}

.btn-icon {
  margin-right: var(--spacing-xs);
}

/* Loading State */
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-3xl);
  gap: var(--spacing-md);
}

.loading-text {
  font-size: var(--font-size-body);
  color: var(--text-secondary);
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: var(--spacing-3xl);
  color: var(--text-tertiary);
}

.empty-icon {
  font-size: 48px;
  margin-bottom: var(--spacing-md);
}

.empty-title {
  font-size: var(--font-size-title2);
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--spacing-sm);
}

.empty-description {
  font-size: var(--font-size-body);
}

/* Course Section */
.course-section {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 var(--spacing-md);
}

.course-grid {
  display: grid;
  gap: var(--spacing-lg);
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
}

/* Course Card Enhancements */
.course-card-image {
  position: relative;
  overflow: hidden;
}

.course-card-image img {
  transition: transform 0.3s ease;
}

.course-card:hover .course-card-image img {
  transform: scale(1.05);
}

.course-card-rating {
  position: absolute;
  top: var(--spacing-sm);
  right: var(--spacing-sm);
  background: rgba(0, 0, 0, 0.7);
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

/* 响应式设计 */
@media (max-width: 768px) {
  .hero-title {
    font-size: 24px;
  }
  
  .filter-form {
    grid-template-columns: 1fr;
  }
  
  .course-grid {
    grid-template-columns: 1fr;
    gap: var(--spacing-md);
  }
  
  .hero-actions {
    flex-direction: column;
    align-items: center;
  }
}

@media (min-width: 769px) and (max-width: 1024px) {
  .course-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 1025px) {
  .course-grid {
    grid-template-columns: repeat(3, 1fr);
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
  font-size: 12px;
  transition: transform 0.2s ease;
}

.course-card:hover .star {
  transform: scale(1.1);
}

/* 动画优化 */
.course-card {
  animation-delay: calc(var(--index, 0) * 0.1s);
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

/* 主题切换动画 */
.home-container {
  transition: background 0.5s ease;
}

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
</style>
