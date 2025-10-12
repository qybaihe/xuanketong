<script setup lang="ts">
import { ref, onMounted, reactive, watch, nextTick } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { courseService, type CourseWithRating } from '@/services/api'

const router = useRouter()

interface CourseWithDisplay extends CourseWithRating {
  Students?: number;
  isPopular?: boolean;
  isNew?: boolean;
  RatingDistribution?: Record<number, number>; // 1-5星评分分布
  MaterialsCount?: number; // 资料数量
  HotComments?: string[]; // 热门评论
  HotMaterials?: string[]; // 热门资料
  Difficulty?: 'easy' | 'medium' | 'hard'; // 难度等级
  Workload?: 'light' | 'medium' | 'heavy'; // 工作量
}

const courses = ref<CourseWithDisplay[]>([])
const loading = ref(true)
const courseSection = ref<HTMLElement>()
const filters = reactive({
  subject: '',
  teacher: ''
})

// 首页搜索相关状态
const searchQuery = ref('')
const searchType = ref<'course' | 'teacher'>('course')
const showSearchDropdown = ref(false)


// 清空筛选
const clearFilters = () => {
  filters.subject = ''
  filters.teacher = ''
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
    router.push(`/courses/${courseId}/rate`)
  }
}

// 首页搜索相关方法
const toggleSearchDropdown = () => {
  showSearchDropdown.value = !showSearchDropdown.value
}

const selectSearchType = (type: 'course' | 'teacher') => {
  searchType.value = type
  showSearchDropdown.value = false
}

const handleHomeSearch = () => {
  if (!searchQuery.value.trim()) return
  
  router.push({
    path: '/search',
    query: {
      q: searchQuery.value,
      type: searchType.value
    }
  })
}

// 跳转到上传资料页面
const goToUploadMaterials = (courseId: number | undefined) => {
  if (courseId) {
    // 这里可以跳转到上传资料页面，暂时先显示提示
    console.log('跳转到上传资料页面:', courseId)
    // router.push(`/courses/${courseId}/upload`)
  }
}


// 数据字段映射函数 - 将后端小写字段映射到前端大写字段
const mapCourseData = (course: CourseWithRating): CourseWithDisplay => {
  return {
    // 原始后端字段
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
    // 兼容字段映射（大写）
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


// 将AverageRating映射到Rating字段（兼容现有模板）
const getDisplayRating = (course: CourseWithDisplay): number => {
  return course.AverageRating || 0
}

// 生成模拟数据
const generateMockData = (course: CourseWithDisplay): CourseWithDisplay => {
  const hotComments = [
    "老师讲课很清晰，作业量适中",
    "课程内容实用，推荐！",
    "考试难度合理，给分不错",
    "老师人很好，有问题会耐心解答",
    "课程设计很棒，学到了很多",
    "作业有点多，但很有收获",
    "老师讲课生动有趣",
    "课程难度适中，适合初学者"
  ]
  
  const hotMaterials = [
    "课程PPT",
    "作业答案",
    "考试重点",
    "参考书籍",
    "实验报告",
    "课程笔记",
    "历年真题",
    "学习资料"
  ]
  
  const difficulties: ('easy' | 'medium' | 'hard')[] = ['easy', 'medium', 'hard']
  const workloads: ('light' | 'medium' | 'heavy')[] = ['light', 'medium', 'heavy']
  
  return {
    ...course,
    Students: Math.floor(Math.random() * 200) + 10,
    isPopular: (course.AverageRating || 0) > 4.0,
    isNew: course.CreatedAt ? new Date().getTime() - new Date(course.CreatedAt).getTime() < 30 * 24 * 60 * 60 * 1000 : false,
    MaterialsCount: Math.floor(Math.random() * 15) + 5,
    HotComments: hotComments.slice(0, Math.floor(Math.random() * 3) + 1),
    HotMaterials: hotMaterials.slice(0, Math.floor(Math.random() * 4) + 2),
    Difficulty: difficulties[Math.floor(Math.random() * 3)],
    Workload: workloads[Math.floor(Math.random() * 3)],
    RatingDistribution: {
      5: Math.floor(Math.random() * 20) + 10,
      4: Math.floor(Math.random() * 15) + 5,
      3: Math.floor(Math.random() * 10) + 2,
      2: Math.floor(Math.random() * 5) + 1,
      1: Math.floor(Math.random() * 3)
    }
  }
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
    courses.value = coursesData.map((course: CourseWithRating) => {
      const mappedCourse = mapCourseData(course)
      return generateMockData(mappedCourse)
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

// 添加磁性效果 - 已禁用
const addMagneticEffect = (element: HTMLElement) => {
  // 磁性效果已禁用以提升美观性
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


watch(filters, debouncedFetchCourses, { deep: true })

onMounted(() => {
  fetchCourses()

  // 初始化性能监控
  initPerformanceMonitoring()

  // 为按钮添加磁性效果
  nextTick(() => {
    document.querySelectorAll('.btn').forEach(btn => {
      addMagneticEffect(btn as HTMLElement)
    })
  })
})
</script>

<template>
  <main class="home-container">
    <!-- 顶部布局：左侧内容 + 右侧筛选 -->
    <div class="top-layout">
      <!-- 左侧内容区域 -->
      <div class="left-content">
        <!-- Header Section -->
        <header class="home-header">
          <div class="header-content">
            <h1 class="page-title">选课通</h1>
            <h2 class="page-subtitle">为软工人打造的评价选课平台</h2>
            
            <!-- 首页搜索框 -->
            <div class="home-search-section">
              <div class="search-container">
                <div class="search-type-dropdown">
                  <button @click="toggleSearchDropdown" class="search-type-button">
                    <span>{{ searchType === 'course' ? '课程' : '老师' }}</span>
                    <svg class="dropdown-arrow" :class="{ 'rotated': showSearchDropdown }" width="12" height="12" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                      <path d="M6 9l6 6 6-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    </svg>
                  </button>
                  <div v-if="showSearchDropdown" class="search-type-options">
                    <button @click="selectSearchType('course')" class="search-type-option" :class="{ active: searchType === 'course' }">
                      课程
                    </button>
                    <button @click="selectSearchType('teacher')" class="search-type-option" :class="{ active: searchType === 'teacher' }">
                      老师
                    </button>
                  </div>
                </div>
                <input
                  v-model="searchQuery"
                  @keyup.enter="handleHomeSearch"
                  type="text"
                  class="search-input"
                  :placeholder="searchType === 'course' ? '搜索课程...' : '搜索老师...'"
                />
                <button @click="handleHomeSearch" class="search-button">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <circle cx="11" cy="11" r="8" stroke="currentColor" stroke-width="2"/>
                    <path d="m21 21-4.35-4.35" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                  搜索
                </button>
              </div>
            </div>
          </div>
        </header>

        <!-- 求评价中心入口 -->
        <section class="evaluation-request-entry-section">
          <div class="entry-content">
            <div class="entry-icon">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z" fill="var(--color-primary)"/>
              </svg>
            </div>
            <div class="entry-text">
              <h3 class="entry-title">发现需要评价的课程</h3>
              <p class="entry-description">帮助同学们找到优质课程，分享你的学习体验</p>
            </div>
            <RouterLink to="/evaluation-requests" class="entry-button">
              求评价中心
            </RouterLink>
          </div>
        </section>
      </div>

      <!-- 右侧筛选区域 -->
      <div class="right-content">
        <!-- Filter Section -->
        <section class="filter-section">
          <h2 class="filter-title">筛选课程</h2>
          <p class="filter-subtitle">按科目、老师快速找到您感兴趣的课程</p>
          <form @submit.prevent="fetchCourses">
            <div class="filter-inputs">
              <div class="form-group">
                <label for="subject" class="form-label">
                  科目
                </label>
                <input type="text" id="subject" v-model="filters.subject" class="input-glass" placeholder="输入科目..." />
              </div>
              <div class="form-group">
                <label for="teacher" class="form-label">
                  老师
                </label>
                <input type="text" id="teacher" v-model="filters.teacher" class="input-glass" placeholder="输入老师..." />
              </div>
              <div class="form-group form-group--button">
                <button type="submit" class="btn btn-primary">
                  搜索课程
                </button>
              </div>
            </div>
          </form>
        </section>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="loading-container">
      <div class="loader"></div>
      <p class="loading-text">正在加载课程...</p>
    </div>

    <!-- Course Grid -->
    <section ref="courseSection" v-else class="course-section">
      <div class="course-header">
        <!--
        <div class="course-header-content">
          <h2 class="course-title">热门课程评价</h2>
          <p class="course-description">查看真实学生评价，找到最适合您的课程</p>
        </div>
        -->
        <div class="course-filters">
          <button class="filter-chip" :class="{ active: !filters.subject && !filters.teacher }"
            @click="clearFilters">
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
        <router-link v-for="course in courses" :key="course.ID" :to="`/courses/${course.ID}`"
          class="course-card">
          <!-- 课程内容 -->
          <div class="course-card-content">
            <!-- 课程标签 -->
            <div class="course-card-tags">
              <span v-for="(tag, index) in getTags(course)" :key="index" :class="['course-card-tag', tag.type]">
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
                    <div class="teacher-credits">{{ course.Credits }} 学分</div>
                  </div>
                </div>
              </div>
            </div>

            <!-- 评价统计区域 - 压缩版 -->
            <div class="course-rating-section">
              <div class="rating-overview-compact">
                <div class="rating-main">
                  <span class="rating-number">{{ getDisplayRating(course).toFixed(1) }}</span>
                  <div class="rating-stars">
                    <span v-for="i in 5" :key="i" :class="['star',
                      i <= Math.floor(getDisplayRating(course)) ? 'star-filled' :
                        (i - 0.5 <= getDisplayRating(course) ? 'star-half' : 'star-empty')]"></span>
                  </div>
                  <span class="rating-count">{{ (course.TotalRatings || 0) }}人</span>
                </div>
                
                <!-- 课程属性标签 -->
                <div class="course-attributes">
                  <span class="attribute-tag difficulty" :class="course.Difficulty">
                    {{ course.Difficulty === 'easy' ? '简单' : course.Difficulty === 'medium' ? '中等' : '困难' }}
                  </span>
                  <span class="attribute-tag workload" :class="course.Workload">
                    {{ course.Workload === 'light' ? '轻松' : course.Workload === 'medium' ? '适中' : '繁重' }}
                  </span>
                </div>
              </div>
            </div>

            <!-- 热门评论区域 -->
            <div v-if="course.HotComments && course.HotComments.length > 0" class="hot-comments-section">
              <div class="hot-comments-header">
                <svg class="comments-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" stroke="var(--color-text)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
                <span class="comments-title">热门评论</span>
              </div>
              <div class="hot-comments-list">
                <div v-for="(comment, index) in course.HotComments.slice(0, 2)" :key="index" class="hot-comment">
                  {{ comment }}
                </div>
              </div>
            </div>

            <!-- 资料信息区域 -->
            <div class="materials-section">
              <div class="materials-info">
                <span class="materials-count">
                  <svg class="materials-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20" stroke="var(--color-text)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z" stroke="var(--color-text)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                  {{ course.MaterialsCount }} 份资料
                </span>
                <div v-if="course.HotMaterials && course.HotMaterials.length > 0" class="hot-materials">
                  <span class="hot-materials-label">热门：</span>
                  <span v-for="(material, index) in course.HotMaterials.slice(0, 2)" :key="index" class="hot-material-tag">
                    {{ material }}
                  </span>
                </div>
              </div>
            </div>

            <!-- 课程操作区域 -->
            <div class="course-card-actions">
              <button class="btn-rate-course" @click.prevent="goToRateCourse(course.ID)">
                评价课程
              </button>
              <button class="btn-upload-materials" @click.prevent="goToUploadMaterials(course.ID)">
                上传资料
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
              <a href="#" class="social-link">集市</a>
              <a href="#" class="social-link">邮箱</a>
            </div>
          </div>
        </div>
      </div>
      <div class="footer-bottom">
        <p class="footer-text">© 2025 选课通 - 让学习更简单</p>
      </div>
    </footer>
  </main>
</template>

<style scoped>
/* ===== 俏皮的新拟物风格 (Playful Neobrutalism) ===== */

/* 全局样式 */
.home-container {
  font-family: sans-serif;
  color: var(--color-text);
  padding: 20px;
  max-width: 420px;
  margin: 0 auto;
  background-color: var(--color-background);
  min-height: 100vh;
}

/* 顶部布局样式 */
.top-layout {
  display: flex;
  flex-direction: column;
  gap: 24px;
  margin-bottom: 24px;
}

.left-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.right-content {
  display: flex;
  flex-direction: column;
}

.filter-section {
  background-color: var(--color-surface);
  border-radius: 12px;
  border: 1px solid var(--color-border);
  padding: 32px 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  min-height: 280px;
}

.filter-header {
  text-align: center;
  margin-bottom: 0;
  padding: 16px 0;
}

.filter-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
  flex: 1;
}

/* 增大间距 */
.course-grid {
  display: grid;
  gap: 16px;
  margin-bottom: 32px;
}

.filter-form {
  background-color: var(--color-surface);
  border-radius: 12px;
  border: 1px solid var(--color-border);
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
  background-color: var(--color-surface);
  border-radius: 12px;
  border: 1px solid var(--color-border);
  padding: 32px 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  transition: none;
  min-height: 120px;
}

/* 移动端纵向排布 */
@media (max-width: 767px) {
  .entry-content {
    flex-direction: column;
    text-align: center;
    gap: 20px;
  }

  .entry-text {
    order: 2;
  }

  .entry-button {
    order: 3;
    width: 100%;
    max-width: 200px;
  }

  .entry-icon {
    order: 1;
    font-size: 40px;
  }
}

/* 平板端优化 */
@media (min-width: 768px) and (max-width: 1023px) {
  .entry-content {
    gap: 20px;
  }

  .entry-icon {
    font-size: 36px;
  }

  .entry-title {
    font-size: 22px;
  }

  .entry-description {
    font-size: 16px;
  }

  .entry-button {
    padding: 14px 28px;
    font-size: 16px;
  }
}

.entry-content:hover {
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
  color: var(--color-text);
  margin-bottom: 4px;
}

.entry-description {
  font-size: 16px;
  color: var(--color-textSecondary);
  margin: 0;
}

.entry-button {
  background-color: var(--color-primary);
  border-radius: 8px;
  border: 1px solid var(--color-border);
  padding: 12px 24px;
  font-size: 16px;
  font-weight: bold;
  color: var(--color-text);
  cursor: pointer;
  transition: none;
  flex-shrink: 0;
}

.entry-button:hover {
}

/* Header Section 样式 */
.home-header {
  margin-bottom: 24px;
  text-align: center;
}

.header-content {
  background-color: var(--color-surface);
  border-radius: 12px;
  border: 1px solid var(--color-border);
  padding: 32px 24px;
  min-height: 120px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

h1.page-title {
  font-size: 28px;
  font-weight: bold;
  text-align: center;
  color: var(--color-text);
  margin-bottom: 8px;
}

h2.page-subtitle {
  font-size: 18px;
  font-weight: normal;
  text-align: center;
  color: var(--color-textSecondary);
  margin: 0;
}

/* 首页搜索框样式 */
.home-search-section {
  width: 100%;
  max-width: 600px;
  margin-top: 16px;
}

.search-container {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #FFFFFF;
  border: 1px solid var(--color-border);
  border-radius: 12px;
  padding: 8px;
  box-shadow: none;
  transition: none;
}

.search-container:focus-within {
  box-shadow: none;
}

.search-type-dropdown {
  position: relative;
}

.search-type-button {
  display: flex;
  align-items: center;
  gap: 4px;
  background: #F7D074;
  border: 1px solid var(--color-border);
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

.dropdown-arrow {
  transition: none;
}

.dropdown-arrow.rotated {
  transform: rotate(180deg);
}

.search-type-options {
  position: absolute;
  top: 100%;
  left: 0;
  z-index: 1000;
  background: #FFFFFF;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  margin-top: 4px;
  min-width: 100%;
  overflow: hidden;
  box-shadow: none;
}

.search-type-option {
  display: block;
  width: 100%;
  background: #FFFFFF;
  border: none;
  padding: 10px 12px;
  font-size: 14px;
  font-weight: bold;
  color: #1A1A1A;
  cursor: pointer;
  transition: background-color 0.2s ease;
  text-align: left;
}

.search-type-option:hover {
  background: #F7D074;
}

.search-type-option.active {
  background: #F7D074;
}

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

.search-button {
  display: flex;
  align-items: center;
  gap: 6px;
  background: #76D7C4;
  color: #1A1A1A;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 8px 16px;
  font-size: 14px;
  font-weight: bold;
  cursor: pointer;
  transition: background-color 0.2s ease;
  white-space: nowrap;
}

.search-button:hover {
  background: #6BC4B1;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .home-search-section {
    max-width: 100%;
  }
  
  .search-container {
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

.page-subtitle {
  font-size: 16px;
  color: var(--color-textSecondary);
  margin: 0 0 16px 0;
}

.page-badge {
  display: inline-block;
  background-color: var(--color-primary);
  color: var(--color-text);
  padding: 6px 12px;
  border-radius: 8px;
  border: 1px solid var(--color-border);
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
  color: var(--color-primary);
  display: block;
}

.header-stats .stat-label {
  font-size: 14px;
  color: var(--color-textSecondary);
}

/* Filter Section 样式 */
.filter-section {
  margin-bottom: 24px;
}

.filter-header {
  background-color: var(--color-surface);
  border-radius: 12px;
  border: 1px solid var(--color-border);
  padding: 24px;
  margin-bottom: 16px;
  text-align: center;
}

.filter-title {
  font-size: 28px;
  font-weight: bold;
  color: var(--color-text);
  margin-bottom: 8px;
}

.filter-subtitle {
  font-size: 18px;
  color: var(--color-textSecondary);
  margin: 0;
}

.filter-form {
  background-color: var(--color-surface);
  border-radius: 12px;
  border: 1px solid var(--color-border);
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
  align-items: end;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-size: 14px;
  font-weight: bold;
  color: var(--color-text);
}

.input-glass {
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 12px;
  font-size: 16px;
  font-family: sans-serif;
  transition: none;
}

.input-glass:hover,
.input-glass:active,
.input-glass:focus {
  outline: none !important;
  box-shadow: none;
  transform: none;
}

.form-group--button {
  margin-top: 8px;
}

.btn {
  background-color: var(--color-primary);
  border-radius: 8px;
  border: 1px solid var(--color-border);
  padding: 12px 24px;
  font-size: 16px;
  font-weight: bold;
  color: var(--color-text);
  cursor: pointer;
  transition: none;
}

.btn:hover,
.btn:active,
.btn:focus {
  box-shadow: none;
  outline: none;
  transform: none;
}

.btn-primary {
  background-color: var(--color-primary);
}

/* Course Section 样式 */
.course-section {
  margin-bottom: 24px;
}

.course-header {
  margin-bottom: 16px;
}

.course-header-content {
  background-color: var(--color-surface);
  border-radius: 12px;
  border: 1px solid var(--color-border);
  padding: 24px;
  margin-bottom: 16px;
  text-align: center;
}

.course-title {
  font-size: 20px;
  font-weight: bold;
  color: var(--color-text);
  margin-bottom: 8px;
}

.course-description {
  font-size: 14px;
  color: var(--color-textSecondary);
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
  background-color: var(--color-surface);
  border-radius: 8px;
  border: 1px solid var(--color-border);
  padding: 8px 16px;
  font-size: 14px;
  font-weight: bold;
  color: var(--color-text);
  cursor: pointer;
  transition: none;
}

.filter-chip:hover {
}

.filter-chip.active {
  background-color: var(--color-primary);
}

.course-grid {
  display: grid;
  gap: 16px;
  margin-bottom: 24px;
}

/* ===== 新卡片设计 ===== */
.course-card {
  background-color: var(--color-surface);
  border-radius: 12px;
  border: 1px solid var(--color-border);
  overflow: hidden;
  text-decoration: none;
  color: var(--color-text);
  display: flex;
  flex-direction: column;
  transition: none;
}

.course-card:hover {
}

.course-card-content {
  padding: 20px;
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  gap: 16px;
}

.course-card-tags {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.course-card-tag {
  padding: 4px 8px;
  border-radius: 6px;
  border: 1px solid var(--color-border);
  font-size: 12px;
  font-weight: bold;
}

.course-card-tag.primary {
  background-color: var(--color-primary);
  color: var(--color-text);
}

.course-card-tag.secondary {
  background-color: var(--color-accent);
  color: var(--color-text);
}

.course-card-tag.accent {
  background-color: var(--color-primary);
  color: var(--color-text);
}

.course-card-title {
  font-size: 20px;
  font-weight: bold;
  color: var(--color-text);
  margin-bottom: 12px;
  line-height: 1.3;
}

.course-card-description {
  color: var(--color-textSecondary);
  margin-bottom: 16px;
}

.rating-count {
  color: var(--color-textSecondary);
}

/* Footer 样式 */
.home-footer {
  margin-top: 24px;
}

.footer-content {
  background-color: var(--color-surface);
  border-radius: 12px;
  border: 1px solid var(--color-border);
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
  color: var(--color-primary);
}

.footer-tagline {
  font-size: 14px;
  color: var(--color-textSecondary);
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
  color: var(--color-text);
  margin-bottom: 8px;
}

.footer-link {
  display: block;
  color: var(--color-textSecondary);
  text-decoration: none;
  font-size: 14px;
  margin-bottom: 8px;
  transition: color 0.2s ease;
}

.footer-link:hover {
  color: var(--color-primary);
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
  padding: 8px;
  min-width: 36px;
  height: 36px;
  background-color: var(--color-primary);
  border-radius: 8px;
  border: 1px solid var(--color-border);
  color: var(--color-text);
  text-decoration: none;
  font-size: 14px;
  font-weight: bold;
  transition: none;
}

.social-link:hover {
}

.footer-bottom {
  text-align: center;
  padding-top: 8px;
}

.footer-text {
  font-size: 14px;
  color: var(--color-textSecondary);
  margin: 0;
}

/* 加载状态样式 */
.loading-container {
  text-align: center;
  padding: 40px 0;
}

.loader {
  border: 4px solid var(--color-border);
  border-top: 4px solid var(--color-primary);
  border-radius: 50%;
  width: 40px;
  height: 40px;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }

  100% {
    transform: rotate(360deg);
  }
}

.loading-text {
  font-size: 16px;
  color: var(--color-text);
}

/* 空状态样式 */
.empty-state {
  text-align: center;
  padding: 40px 0;
  background-color: var(--color-surface);
  border-radius: 12px;
  border: 1px solid var(--color-border);
}

.empty-title {
  font-size: 20px;
  font-weight: bold;
  color: var(--color-text);
  margin-bottom: 8px;
}

.empty-description {
  font-size: 14px;
  color: var(--color-textSecondary);
  margin: 0 0 16px 0;
}

.mt-lg {
  margin-top: 16px;
}

/* ===== 课程评价功能样式 ===== */
.course-card-header {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 2px solid var(--color-border);
}

/* 授课老师部分 */
.course-teacher-section {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background-color: var(--color-surface);
  border-radius: 8px;
  border: 1px solid var(--color-border);
  transition: none;
}

.course-teacher-section:hover {
}

.teacher-label {
  font-size: 12px;
  color: var(--color-text);
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
  gap: 12px;
}

.teacher-name {
  font-size: 16px;
  font-weight: bold;
  color: var(--color-text);
  flex-shrink: 0;
}

.teacher-credits {
  font-size: 14px;
  color: var(--color-text);
  font-weight: bold;
  background-color: var(--color-primary);
  padding: 4px 8px;
  border-radius: 6px;
  border: 1px solid var(--color-border);
  white-space: nowrap;
  flex-shrink: 0;
}

.course-rating-section {
  padding: 16px;
  background-color: var(--color-surface);
  border-radius: 8px;
  border: 1px solid var(--color-border);
  margin-bottom: 16px;
}

.rating-overview {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-direction: column;
  gap: 16px;
}

/* 压缩版评分样式 */
.rating-overview-compact {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.rating-main {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.rating-number {
  font-size: 20px;
  font-weight: bold;
  color: var(--color-primary);
}

.rating-stars {
  display: flex;
  gap: 2px;
}

.rating-count {
  font-size: 14px;
  color: var(--color-text);
  margin-left: auto;
  font-weight: 500;
}

/* 课程属性标签 */
.course-attributes {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.attribute-tag {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: bold;
  text-transform: uppercase;
  border: 1px solid var(--color-border);
}

.attribute-tag.difficulty.easy {
  background: var(--color-accent);
  color: var(--color-text);
}

.attribute-tag.difficulty.medium {
  background: var(--color-primary);
  color: var(--color-text);
}

.attribute-tag.difficulty.hard {
  background: var(--color-secondary);
  color: #FFFFFF;
}

.attribute-tag.workload.light {
  background: var(--color-accent);
  color: var(--color-text);
}

.attribute-tag.workload.medium {
  background: var(--color-primary);
  color: var(--color-text);
}

.attribute-tag.workload.heavy {
  background: var(--color-secondary);
  color: #FFFFFF;
}

/* 热门评论区域 */
.hot-comments-section {
  margin-bottom: 12px;
  padding: 12px;
  background-color: var(--color-surface);
  border-radius: 8px;
  border: 1px solid var(--color-border);
}

.hot-comments-header {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 8px;
}

.comments-icon {
  font-size: 12px;
}

.comments-title {
  font-size: 14px;
  font-weight: bold;
  color: var(--color-text);
}

.hot-comments-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.hot-comment {
  font-size: 13px;
  color: var(--color-text);
  line-height: 1.4;
  padding: 4px 0;
  border-left: 3px solid var(--color-primary);
  padding-left: 8px;
  font-weight: 500;
}

/* 资料信息区域 */
.materials-section {
  margin-bottom: 12px;
  padding: 12px;
  background-color: var(--color-surface);
  border-radius: 8px;
  border: 1px solid var(--color-border);
}

.materials-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.materials-count {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  font-weight: bold;
  color: var(--color-text);
}

.materials-icon {
  font-size: 12px;
}

.hot-materials {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.hot-materials-label {
  font-size: 13px;
  color: var(--color-text);
  font-weight: bold;
}

.hot-material-tag {
  padding: 3px 8px;
  background-color: var(--color-primary);
  color: var(--color-text);
  border-radius: 4px;
  font-size: 12px;
  border: 1px solid var(--color-border);
  font-weight: bold;
}

.rating-average {
  text-align: center;
}

.rating-number {
  font-size: 36px;
  font-weight: bold;
  color: var(--color-primary);
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
  color: var(--color-textSecondary);
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
  color: var(--color-textSecondary);
  min-width: 30px;
  text-align: right;
}

.bar-container {
  flex: 1;
  height: 8px;
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 4px;
  overflow: hidden;
  position: relative;
}

.bar-fill {
  height: 100%;
  background-color: var(--color-primary);
  transition: width 0.3s ease;
}

.bar-count {
  font-size: 14px;
  color: var(--color-textSecondary);
  text-align: left;
}

.course-card-actions {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: auto;
  padding-top: 16px;
  gap: 12px;
}

.btn-rate-course {
  background-color: var(--color-accent);
  border-radius: 8px;
  border: 1px solid var(--color-border);
  padding: 12px 20px;
  font-size: 14px;
  font-weight: bold;
  color: var(--color-text);
  cursor: pointer;
  transition: none;
  flex: 1;
}

.btn-rate-course:hover,
.btn-rate-course:active,
.btn-rate-course:focus {
  box-shadow: none;
  outline: none;
  transform: none;
}

.btn-upload-materials {
  background-color: var(--color-primary);
  border-radius: 8px;
  border: 1px solid var(--color-border);
  color: var(--color-text);
  font-weight: bold;
  font-size: 14px;
  padding: 12px 20px;
  cursor: pointer;
  transition: none;
  flex: 1;
}

.btn-upload-materials:hover,
.btn-upload-materials:active,
.btn-upload-materials:focus {
  box-shadow: none;
  outline: none;
  transform: none;
}


.btn-rate-course .btn-icon {
  font-size: 16px;
}

/* 评分星级样式 */
.star {
  font-size: 18px;
  transition: none;
  position: relative;
  display: inline-block;
  width: 18px;
  height: 18px;
}

.star-filled::before {
  content: "★";
  color: var(--color-primary);
}

.star-half::before {
  content: "★";
  color: var(--color-primary);
  background: linear-gradient(90deg, var(--color-primary) 50%, var(--color-textSecondary) 50%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.star-empty::before {
  content: "★";
  color: var(--color-textSecondary);
}

.course-card:hover .star {
  transform: none;
}

/* 响应式设计 */
@media (max-width: 767px) {
  /* 移动端课程网格优化 */
  .course-grid {
    grid-template-columns: 1fr;
    gap: 12px;
    margin-bottom: 16px;
  }

  /* 移动端课程卡片紧凑布局 */
  .course-card {
    border: 1px solid var(--color-border);
    border-radius: 12px;
    box-shadow: none;
  }

  .course-card-content {
    padding: 14px;
    gap: 10px;
  }

  .course-card-tags {
    margin-bottom: 10px;
    gap: 6px;
  }

  .course-card-tag {
    padding: 3px 8px;
    font-size: 11px;
    border-radius: 4px;
  }

  .course-card-title {
    font-size: 18px;
    margin-bottom: 8px;
    line-height: 1.3;
  }

  .course-card-header {
    gap: 10px;
    margin-bottom: 10px;
    padding-bottom: 10px;
  }

  .course-teacher-section {
    gap: 6px;
  }

  .teacher-label {
    font-size: 11px;
    padding: 3px 6px;
  }

  .teacher-name {
    font-size: 14px;
  }

  .course-rating-section {
    gap: 10px;
    margin-bottom: 10px;
  }

  .rating-display {
    gap: 6px;
  }

  .rating-stars {
    gap: 2px;
  }

  .star {
    font-size: 14px;
  }

  .rating-text {
    font-size: 13px;
  }

  .rating-count {
    font-size: 12px;
  }

  .course-stats {
    gap: 10px;
    margin-bottom: 10px;
  }

  .stat-item {
    gap: 4px;
  }

  .stat-value {
    font-size: 14px;
  }

  .stat-label {
    font-size: 11px;
  }

  .course-card-actions {
    padding-top: 10px;
    gap: 10px;
  }

  .btn-rate-course,
  .btn-upload-materials {
    padding: 8px 16px;
    font-size: 13px;
    border-radius: 6px;
  }

  /* 移动端课程卡片平衡排版 */
  .course-card-content {
    padding: 24px;
    gap: 8px;
  }

  .course-card-tags {
    margin-bottom: 6px;
    gap: 4px;
  }

  .course-card-tag {
    padding: 4px 8px;
    font-size: 11px;
    font-weight: bold;
  }

  .course-card-title {
    font-size: 18px;
    margin-bottom: 8px;
    line-height: 1.3;
    font-weight: bold;
  }

  /* 紧凑教师信息 */
  .course-teacher-section {
    border: none;
    background: transparent;
    margin-bottom: 8px;
    align-items: center;
    gap: 6px;
    padding: 0;
    border-bottom: 1px solid var(--color-border);
    padding-bottom: 8px;
    border-radius: 0;
  }

  .teacher-label {
    border: none;
    padding: 4px 8px;
    font-size: 12px;
    background-color: var(--color-background);
    border-radius: 6px;
  }

  .teacher-name {
    border: none;
    font-size: 16px;
    font-weight: 500;
  }

  /* 紧凑评分区域 */
  .course-rating-section {
    border: none;
    background: transparent;
    margin-bottom: 8px;
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 0;
    border-bottom: 1px solid var(--color-border);
    padding-bottom: 8px;
    border-radius: 0;
  }

  .rating-display {
    border: none;
    background: transparent;
    gap: 3px;
    align-items: center;
  }

  .rating-stars {
    border: none;
    background: transparent;
    gap: 1px;
  }

  .star {
    font-size: 16px;
  }

  .rating-text {
    font-size: 16px;
    font-weight: bold;
  }

  .rating-count {
    font-size: 14px;
    color: var(--color-text-secondary);
  }

  /* 紧凑统计信息 */
  .course-stats {
    border: none;
    background: transparent;
    margin-bottom: 8px;
    display: flex;
    gap: 12px;
    padding: 0;
    border-bottom: 1px solid var(--color-border);
    padding-bottom: 8px;
    border-radius: 0;
  }

  .stat-item {
    border: none;
    background: transparent;
    gap: 1px;
    text-align: left;
  }

  .stat-value {
    font-size: 16px;
    font-weight: bold;
  }

  .stat-label {
    font-size: 13px;
    color: var(--color-text-secondary);
  }

  /* 紧凑操作按钮 */
  .course-card-actions {
    border: none;
    background: transparent;
    gap: 6px;
    justify-content: center;
    padding: 0;
  }

  .btn-rate-course,
  .btn-upload-materials {
    padding: 8px 12px;
    font-size: 15px;
    border-radius: 6px;
    flex: 1;
    max-width: 120px;
  }

  /* 紧凑其他section元素 */
  .course-card-header {
    border: none;
    background: transparent;
    margin-bottom: 6px;
    gap: 6px;
    padding: 0;
  }

  .hot-comments-section {
    border: none;
    background: transparent;
    margin-bottom: 6px;
    padding: 0;
  }

  .materials-section {
    border: none;
    background: transparent;
    margin-bottom: 6px;
    padding: 0;
  }

  .evaluation-request-entry-section {
    border: none;
    background: transparent;
    margin-bottom: 6px;
    padding: 0;
  }

  .footer-links {
    grid-template-columns: 1fr;
    gap: 24px;
  }

  .footer-section {
    text-align: center;
  }

  .footer-section-title {
    text-align: center;
    margin-bottom: 12px;
  }

  .footer-link {
    display: inline-block;
    margin: 0 8px 8px 8px;
  }
}

/* 超小屏幕设备优化 (iPhone SE等) */
@media (max-width: 480px) {
  .course-card-content {
    padding: 10px;
    gap: 8px;
  }

  .course-card-title {
    font-size: 16px;
    margin-bottom: 6px;
  }

  .course-card-tags {
    margin-bottom: 6px;
  }

  .course-card-tag {
    padding: 2px 6px;
    font-size: 10px;
  }

  .course-card-header {
    gap: 8px;
    margin-bottom: 8px;
    padding-bottom: 8px;
  }

  .teacher-label {
    font-size: 10px;
    padding: 2px 4px;
  }

  .teacher-name {
    font-size: 13px;
  }

  .course-rating-section {
    gap: 8px;
    margin-bottom: 8px;
  }

  .star {
    font-size: 12px;
  }

  .rating-text {
    font-size: 12px;
  }

  .rating-count {
    font-size: 11px;
  }

  .course-stats {
    gap: 8px;
    margin-bottom: 8px;
  }

  .stat-value {
    font-size: 13px;
  }

  .stat-label {
    font-size: 10px;
  }

  .course-card-actions {
    padding-top: 8px;
    gap: 8px;
  }

  .btn-rate-course,
  .btn-upload-materials {
    padding: 6px 12px;
    font-size: 12px;
  }

  /* 超小屏幕平衡排版 */
  .course-card-content {
    padding: 20px;
    gap: 6px;
  }

  .course-card-title {
    font-size: 18px;
    margin-bottom: 6px;
    line-height: 1.2;
  }

  .course-card-tags {
    margin-bottom: 5px;
    gap: 3px;
  }

  .course-card-tag {
    padding: 3px 6px;
    font-size: 12px;
    font-weight: bold;
  }

  .course-teacher-section {
    margin-bottom: 5px;
    gap: 4px;
    border-bottom: 1px solid var(--color-border);
    padding-bottom: 6px;
    border-radius: 0;
  }

  .teacher-name {
    font-size: 15px;
  }

  .course-rating-section {
    margin-bottom: 5px;
    gap: 4px;
    border-bottom: 1px solid var(--color-border);
    padding-bottom: 6px;
    border-radius: 0;
  }

  .star {
    font-size: 15px;
  }

  .rating-text {
    font-size: 15px;
  }

  .rating-count {
    font-size: 13px;
  }

  .course-stats {
    margin-bottom: 6px;
    gap: 8px;
    border-bottom: 1px solid var(--color-border);
    padding-bottom: 6px;
    border-radius: 0;
  }

  .stat-value {
    font-size: 15px;
  }

  .stat-label {
    font-size: 12px;
  }

  .course-card-actions {
    gap: 4px;
  }

  .btn-rate-course,
  .btn-upload-materials {
    padding: 6px 10px;
    font-size: 14px;
    max-width: 100px;
  }
}

@media (min-width: 768px) {
  .home-container {
    max-width: 768px;
    background-color: var(--color-background);
  }

  .top-layout {
    flex-direction: row;
    align-items: flex-start;
    gap: 32px;
  }

  .left-content {
    flex: 1;
    gap: 24px;
    display: flex;
    flex-direction: column;
  }

  .right-content {
    flex: 0 0 300px;
    display: flex;
    flex-direction: column;
  }

  .header-content {
    padding: 32px 24px;
    display: flex;
    flex-direction: column;
    justify-content: center;
  }

  .entry-content {
    padding: 32px 24px;
    display: flex;
    align-items: center;
  }

  .filter-section {
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 32px 24px;
  }

  .filter-title {
    text-align: center;
  }

  .filter-subtitle {
    text-align: center;
  }

  form {
    display: flex;
    flex-direction: column;
    gap: 20px;
    margin-top: 16px;
  }

  .filter-inputs {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
    align-items: end;
  }

  .form-group {
    margin-bottom: 0;
  }

  .form-group--button {
    grid-column: 1 / -1;
    display: flex;
    justify-content: center;
    margin-top: 12px;
  }

  .course-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 1024px) {
  .home-container {
    max-width: 1400px;
    padding: 40px 40px 16px 40px;
    background-color: var(--color-background);
  }

  .top-layout {
    gap: 40px;
    align-items: flex-start;
  }

  .right-content {
    flex: 0 0 350px;
    display: flex;
    flex-direction: column;
  }

  .header-content {
    padding: 36px 24px;
    display: flex;
    flex-direction: column;
    justify-content: center;
  }

  .entry-content {
    padding: 36px 24px;
    display: flex;
    align-items: center;
  }

  .filter-section {
    display: flex;
    flex-direction: column;
    gap: 6px;
    padding: 36px 24px;
  }

  .filter-title {
    text-align: center;
  }

  .filter-subtitle {
    text-align: center;
  }

  form {
    display: flex;
    flex-direction: column;
    gap: 24px;
    margin-top: 20px;
  }

  .filter-inputs {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 20px;
    align-items: end;
  }

  .form-group {
    margin-bottom: 0;
  }

  .form-group--button {
    grid-column: 1 / -1;
    display: flex;
    justify-content: center;
    margin-top: 16px;
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
    justify-content: center;
  }
}

@media (min-width: 1440px) {
  .home-container {
    max-width: 1600px;
    background-color: var(--color-background);
  }

  .top-layout {
    align-items: flex-start;
  }

  .right-content {
    flex: 0 0 400px;
    display: flex;
    flex-direction: column;
  }

  .header-content {
    padding: 40px 24px;
    display: flex;
    flex-direction: column;
    justify-content: center;
  }

  .entry-content {
    padding: 40px 24px;
    display: flex;
    align-items: center;
  }

  .filter-section {
    display: flex;
    flex-direction: column;
    gap: 6px;
    padding: 40px 24px;
  }

  .filter-title {
    margin: 0 0 8px 0;
    text-align: center;
  }

  .filter-subtitle {
    margin: 0 0 16px 0;
    text-align: center;
  }

  form {
    display: flex;
    flex-direction: column;
    gap: 28px;
    margin-top: 24px;
  }

  .filter-inputs {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 24px;
    align-items: end;
  }

  .form-group {
    margin-bottom: 0;
  }

  .form-group--button {
    grid-column: 1 / -1;
    display: flex;
    justify-content: center;
    margin-top: 20px;
  }

  .course-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}

@media (min-width: 1600px) {
  .home-container {
    max-width: 1800px;
    background-color: var(--color-background);
  }

  .top-layout {
    align-items: flex-start;
  }

  .right-content {
    flex: 0 0 450px;
    display: flex;
    flex-direction: column;
  }

  .header-content {
    padding: 36px 24px;
    display: flex;
    flex-direction: column;
    justify-content: center;
  }

  .entry-content {
    padding: 36px 24px;
    display: flex;
    align-items: center;
  }

  .filter-section {
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 44px 24px;
  }

  .filter-title {
    text-align: center;
  }

  .filter-subtitle {
    text-align: center;
  }

  form {
    display: flex;
    flex-direction: column;
    gap: 32px;
    margin-top: 28px;
  }

  .filter-inputs {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 28px;
    align-items: end;
  }

  .form-group {
    margin-bottom: 0;
  }

  .form-group--button {
    grid-column: 1 / -1;
    display: flex;
    justify-content: center;
    margin-top: 24px;
  }

  .course-grid {
    grid-template-columns: repeat(5, 1fr);
  }
}

@media (min-width: 1920px) {
  .home-container {
    max-width: 2200px;
    background-color: var(--color-background);
  }

  .top-layout {
    align-items: flex-start;
  }

  .right-content {
    flex: 0 0 500px;
    display: flex;
    flex-direction: column;
  }

  .header-content {
    padding: 40px 24px;
    display: flex;
    flex-direction: column;
    justify-content: center;
  }

  .entry-content {
    padding: 40px 24px;
    display: flex;
    align-items: center;
  }

  .filter-section {
    display: flex;
    flex-direction: column;
    gap: 6px;
    padding: 48px 24px;
  }

  .filter-title {
    margin: 0 0 8px 0;
    text-align: center;
  }

  .filter-subtitle {
    margin: 0 0 16px 0;
    text-align: center;
  }

  form {
    display: flex;
    flex-direction: column;
    gap: 36px;
    margin-top: 32px;
  }

  .filter-inputs {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 32px;
    align-items: end;
  }

  .form-group {
    margin-bottom: 0;
  }

  .form-group--button {
    grid-column: 1 / -1;
    display: flex;
    justify-content: center;
    margin-top: 28px;
  }

  .course-grid {
    grid-template-columns: repeat(6, 1fr);
  }
}
</style>
