<script setup lang="ts">
import { ref, onMounted, reactive, watch, computed, nextTick } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { courseService, type Course, type CourseWithRating } from '@/services/api'

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

// 滚动到课程区域
const scrollToCourses = () => {
  if (courseSection.value) {
    courseSection.value.scrollIntoView({ behavior: 'smooth' })
  }
}

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

// 跳转到上传资料页面
const goToUploadMaterials = (courseId: number | undefined) => {
  if (courseId) {
    // 这里可以跳转到上传资料页面，暂时先显示提示
    console.log('跳转到上传资料页面:', courseId)
    // router.push(`/courses/${courseId}/upload`)
  }
}

// 获取评分百分比（用于评分条）
const getRatingPercentage = (course: CourseWithDisplay, starLevel: number) => {
  // 使用真实的评分分布数据
  const ratingDistribution = course.RatingDistribution || {}
  const countForStar = ratingDistribution[starLevel] || 0
  const totalRatings = course.TotalRatings || 1

  // 计算百分比
  const percentage = (countForStar / totalRatings) * 100
  return Math.max(5, Math.min(95, percentage)) // 限制在5-95之间，确保可见
}

// 获取评分数量（用于评分条）
const getRatingCount = (course: CourseWithDisplay, starLevel: number) => {
  // 使用真实的评分分布数据
  const ratingDistribution = course.RatingDistribution || {}
  return ratingDistribution[starLevel] || 0
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

// 生成模拟数据
const generateMockData = (course: any): CourseWithDisplay => {
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
    courses.value = coursesData.map((course: any, index) => {
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
    <!-- Header Section -->
    <header class="home-header">
      <div class="header-content">
        <h1 class="page-title">选课通</h1>
        <h2 class="page-title">为软工人打造的评价选课平台</h2>
      </div>
    </header>

    <!-- 求评价中心入口 -->
    <section class="evaluation-request-entry-section">
      <div class="entry-content">
        <div class="entry-icon">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z" fill="#F7D074"/>
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

    <!-- Filter Section -->
    <section class="filter-section">
      <form @submit.prevent="fetchCourses" class="filter-form">
        <div class="filter-header">
          <h2 class="filter-title">筛选课程</h2>
          <p class="filter-subtitle">按科目、老师快速找到您感兴趣的课程</p>
        </div>
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
        <router-link v-for="(course, index) in courses" :key="course.ID" :to="`/courses/${course.ID}`"
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
                  <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" stroke="#1A1A1A" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
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
                    <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20" stroke="#1A1A1A" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z" stroke="#1A1A1A" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
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
  gap: 16px;
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

h1.page-title {
  font-size: 28px;
  font-weight: bold;
  text-align: center;
  color: #1A1A1A;
  margin-bottom: 8px;
}

h2.page-title{
  font-size: 20px;
  font-weight: bold;
  text-align: center;
  color: #1A1A1A;
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
  margin-bottom: 12px;
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
  padding: 8px;
  min-width: 36px;
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
  padding-top: 8px;
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
  0% {
    transform: rotate(0deg);
  }

  100% {
    transform: rotate(360deg);
  }
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
  gap: 12px;
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 2px solid #000000;
}

/* 授课老师部分 */
.course-teacher-section {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background-color: #FFFFFF;
  border-radius: 8px;
  border: 2px solid #000000;
  box-shadow: 3px 3px 0px 0px #000000;
  transition: transform 0.2s ease;
}

.course-teacher-section:hover {
  transform: translate(-2px, -2px);
  box-shadow: 5px 5px 0px 0px #000000;
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
  gap: 12px;
}

.teacher-name {
  font-size: 16px;
  font-weight: bold;
  color: #1A1A1A;
  flex-shrink: 0;
}

.teacher-credits {
  font-size: 14px;
  color: #1A1A1A;
  font-weight: bold;
  background-color: #F7D074;
  padding: 4px 8px;
  border-radius: 6px;
  border: 2px solid #000000;
  white-space: nowrap;
  flex-shrink: 0;
}

.course-rating-section {
  padding: 16px;
  background-color: #FFFFFF;
  border-radius: 8px;
  border: 2px solid #000000;
  box-shadow: 3px 3px 0px 0px #000000;
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
  color: #F7D074;
}

.rating-stars {
  display: flex;
  gap: 2px;
}

.rating-count {
  font-size: 14px;
  color: #1A1A1A;
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
  border: 2px solid #000000;
  box-shadow: 2px 2px 0px 0px #000000;
}

.attribute-tag.difficulty.easy {
  background: #76D7C4;
  color: #1A1A1A;
}

.attribute-tag.difficulty.medium {
  background: #F7D074;
  color: #1A1A1A;
}

.attribute-tag.difficulty.hard {
  background: #FF6B6B;
  color: #FFFFFF;
}

.attribute-tag.workload.light {
  background: #76D7C4;
  color: #1A1A1A;
}

.attribute-tag.workload.medium {
  background: #F7D074;
  color: #1A1A1A;
}

.attribute-tag.workload.heavy {
  background: #FF6B6B;
  color: #FFFFFF;
}

/* 热门评论区域 */
.hot-comments-section {
  margin-bottom: 12px;
  padding: 12px;
  background-color: #FFFFFF;
  border-radius: 8px;
  border: 2px solid #000000;
  box-shadow: 3px 3px 0px 0px #000000;
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
  color: #1A1A1A;
}

.hot-comments-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.hot-comment {
  font-size: 13px;
  color: #1A1A1A;
  line-height: 1.4;
  padding: 4px 0;
  border-left: 3px solid #F7D074;
  padding-left: 8px;
  font-weight: 500;
}

/* 资料信息区域 */
.materials-section {
  margin-bottom: 12px;
  padding: 12px;
  background-color: #FFFFFF;
  border-radius: 8px;
  border: 2px solid #000000;
  box-shadow: 3px 3px 0px 0px #000000;
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
  color: #1A1A1A;
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
  color: #1A1A1A;
  font-weight: bold;
}

.hot-material-tag {
  padding: 3px 8px;
  background-color: #F7D074;
  color: #1A1A1A;
  border-radius: 4px;
  font-size: 12px;
  border: 2px solid #000000;
  font-weight: bold;
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
  gap: 12px;
}

.btn-rate-course {
  background-color: #76D7C4;
  border-radius: 8px;
  border: 3px solid #000000;
  box-shadow: 4px 4px 0px 0px #000000;
  padding: 12px 20px;
  font-size: 14px;
  font-weight: bold;
  color: #1A1A1A;
  cursor: pointer;
  transition: transform 0.2s ease;
  flex: 1;
}

.btn-rate-course:hover {
  transform: translate(-2px, -2px);
  box-shadow: 6px 6px 0px 0px #000000;
}

.btn-upload-materials {
  background-color: #F7D074;
  border-radius: 8px;
  border: 3px solid #000000;
  box-shadow: 4px 4px 0px 0px #000000;
  color: #1A1A1A;
  font-weight: bold;
  font-size: 14px;
  padding: 12px 20px;
  cursor: pointer;
  transition: transform 0.2s ease;
  flex: 1;
}

.btn-upload-materials:hover {
  transform: translate(-2px, -2px);
  box-shadow: 6px 6px 0px 0px #000000;
}

.btn-rate-course:active {
  transform: translate(-1px, -1px);
  box-shadow: 3px 3px 0px 0px #000000;
}

.btn-upload-materials:active {
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
@media (max-width: 767px) {
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
    padding: 40px 40px 16px 40px;
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
    justify-content: center;
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
