<script setup lang="ts">
import { ref, onMounted, computed, nextTick } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { commentService, evaluationRequestService } from '@/services/api'
import api from '@/services/api'
import axios from 'axios'

interface Course {
  // 后端原始字段（小写）
  id: number;
  name: string;
  description: string;
  grade: string;
  semester: string;
  subject: string;
  teacher: string;
  credits: number;
  createdAt: string;
  updatedAt: string;
  // 兼容字段（大写）
  ID?: number;
  Name?: string;
  Description?: string;
  Grade?: string;
  Semester?: string;
  Subject?: string;
  Teacher?: string;
  Credits?: number;
  CreatedAt?: string;
  UpdatedAt?: string;
}

// 数据映射函数
const mapCourseData = (course: any): Course => {
  return {
    // 原始后端字段
    id: course.id,
    name: course.name,
    description: course.description,
    grade: course.grade,
    semester: course.semester,
    subject: course.subject,
    teacher: course.teacher,
    credits: course.credits,
    createdAt: course.createdAt,
    updatedAt: course.updatedAt,
    // 兼容字段映射
    ID: course.id,
    Name: course.name,
    Description: course.description,
    Grade: course.grade,
    Semester: course.semester,
    Subject: course.subject,
    Teacher: course.teacher,
    Credits: course.credits,
    CreatedAt: course.createdAt,
    UpdatedAt: course.updatedAt
  }
}

interface Rating {
  ID: number;
  UserID: number;
  CourseID: number;
  Score: number;
  Difficulty?: number;
  Usefulness?: number;
  Teaching?: number;
  Username?: string;
  // 后端API返回的字段
  username?: string;
  nickname?: string;
}

interface Comment {
  ID: number;
  UserID: number;
  CourseID: number;
  Content: string;
  Username?: string;
  CreatedAt?: string;
  // 后端API返回的字段
  createdAt?: string;
}

const route = useRoute()
const authStore = useAuthStore()
const course = ref<Course | null>(null)
const ratings = ref<Rating[]>([])
const comments = ref<Comment[]>([])
const loading = ref(true)
const evaluationRequestLoading = ref(false)
const hasEvaluationRequested = ref(false)

// 只使用natural主题
const themeClass = computed(() => 'theme-natural')

// 检查用户是否登录
const canSubmit = computed(() => authStore.isAuthenticated)

// 获取当前用户信息
const currentUser = computed(() => authStore.user)

// 添加数据刷新函数
const fetchRatings = async (courseId: number) => {
  try {
    const response = await api.get(`/courses/${courseId}/ratings`)
    // 使用后端返回的username，如果没有则使用nickname，如果都没有才使用默认
    ratings.value = response.data.data.map((rating: any) => ({
      ...rating,
      Username: rating.user?.username || rating.user?.nickname || `用户${rating.UserID}`,
      Score: !isNaN(rating.score) && isFinite(rating.score) ? rating.score : (!isNaN(rating.Score) && isFinite(rating.Score) ? rating.Score : 0)
    }))
  } catch (error) {
    console.error('获取评分失败:', error)
  }
}

const fetchComments = async (courseId: number) => {
  try {
    const response = await api.get(`/courses/${courseId}/comments`)
    comments.value = (response.data.data || []).map((comment: any) => ({
      ...comment,
      Username: comment.user?.username || comment.user?.nickname || `用户${comment.UserID}`,
      CreatedAt: comment.createdAt ? new Date(comment.createdAt).toLocaleDateString() : new Date().toLocaleDateString()
    }))
  } catch (error) {
    console.error('获取评论失败:', error)
  }
}

// 计算平均评分
const averageRating = computed(() => {
  if (ratings.value.length === 0) return 0
  const validRatings = ratings.value.filter(rating => !isNaN(rating.Score) && isFinite(rating.Score))
  if (validRatings.length === 0) return 0
  const sum = validRatings.reduce((acc, rating) => acc + rating.Score, 0)
  return Number((sum / validRatings.length).toFixed(1))
})

// 计算平均难度评分
const averageDifficulty = computed(() => {
  if (ratings.value.length === 0) return 0
  const validRatings = ratings.value.filter(rating =>
    rating.Difficulty !== undefined && !isNaN(rating.Difficulty) && isFinite(rating.Difficulty)
  )
  if (validRatings.length === 0) return 0
  const sum = validRatings.reduce((acc, rating) => acc + (rating.Difficulty as number), 0)
  return Number((sum / validRatings.length).toFixed(1))
})

// 计算平均实用性评分
const averageUsefulness = computed(() => {
  if (ratings.value.length === 0) return 0
  const validRatings = ratings.value.filter(rating =>
    rating.Usefulness !== undefined && !isNaN(rating.Usefulness) && isFinite(rating.Usefulness)
  )
  if (validRatings.length === 0) return 0
  const sum = validRatings.reduce((acc, rating) => acc + (rating.Usefulness as number), 0)
  return Number((sum / validRatings.length).toFixed(1))
})

// 计算平均教学质量评分
const averageTeaching = computed(() => {
  if (ratings.value.length === 0) return 0
  const validRatings = ratings.value.filter(rating =>
    rating.Teaching !== undefined && !isNaN(rating.Teaching) && isFinite(rating.Teaching)
  )
  if (validRatings.length === 0) return 0
  const sum = validRatings.reduce((acc, rating) => acc + (rating.Teaching as number), 0)
  return Number((sum / validRatings.length).toFixed(1))
})

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

// 获取标签
const getTags = (course: Course) => {
  const tags = []
  if (course.Grade) tags.push({ text: course.Grade, type: 'primary' })
  if (course.Semester) tags.push({ text: course.Semester, type: 'secondary' })
  if (course.Subject) tags.push({ text: course.Subject, type: 'accent' })
  return tags
}


// 发起求评价
const submitEvaluationRequest = async () => {
  if (!authStore.isAuthenticated) {
    alert('请先登录后再发起求评价！')
    return
  }
  
  if (!currentUser.value) {
    alert('无法获取用户信息，请重新登录！')
    return
  }
  
  const courseId = Number(route.params.id)
  if (!courseId) {
    alert('课程ID无效')
    return
  }
  
  evaluationRequestLoading.value = true
  try {
    await evaluationRequestService.createEvaluationRequest(courseId)
    hasEvaluationRequested.value = true
    alert('求评价请求发送成功！')
  } catch (error: any) {
    console.error('求评价请求失败:', error)
    const errorMessage = error.response?.data?.error || '求评价请求失败，请重试！'
    alert(errorMessage)
  } finally {
    evaluationRequestLoading.value = false
  }
}

onMounted(async () => {
  try {
    loading.value = true
    const courseId = route.params.id
    const [courseResponse, ratingsResponse, commentsResponse] = await Promise.all([
      axios.get(`http://localhost:8080/api/v1/courses/${courseId}`),
      axios.get(`http://localhost:8080/api/v1/courses/${courseId}/ratings`),
      axios.get(`http://localhost:8080/api/v1/courses/${courseId}/comments`)
    ])
    course.value = mapCourseData(courseResponse.data.data)
    // 使用后端返回的username，如果没有则使用nickname，如果都没有才使用默认
    ratings.value = ratingsResponse.data.data.map((rating: Rating) => ({
      ...rating,
      Username: (rating as any).user?.username || (rating as any).user?.nickname || `用户${rating.UserID}`,
      Score: !isNaN((rating as any).score) && isFinite((rating as any).score) ? (rating as any).score : (!isNaN(rating.Score) && isFinite(rating.Score) ? rating.Score : 0)
    }))
    // 使用后端返回的username，如果没有则使用nickname，如果都没有才使用默认
    comments.value = (commentsResponse.data.data || []).map((comment: Comment) => ({
      ...comment,
      Username: (comment as any).user?.username || (comment as any).user?.nickname || `用户${comment.UserID}`,
      CreatedAt: (comment as any).createdAt ? new Date((comment as any).createdAt).toLocaleDateString() : new Date().toLocaleDateString()
    }))
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="course-detail-container">
    <!-- 返回按钮 -->
    <div class="back-button-container">
      <RouterLink to="/" class="btn back-button">
        <span class="btn-icon">←</span>
        返回课程列表
      </RouterLink>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="loading-container">
      <div class="loader"></div>
      <p class="loading-text">正在加载课程详情...</p>
    </div>

    <!-- Course Content -->
    <div v-else-if="course" class="course-content">
      <div class="main-content">
        <!-- Course Header Card -->
        <div class="card course-header">
          <div class="course-credits-badge">
            {{ course.Credits }} 学分
          </div>
          
          <h1 class="course-title">{{ course.Name }}</h1>
          
          <div class="course-meta">
            <div class="meta-item">
              <span class="meta-icon">👨‍🏫</span>
              <span class="meta-text">{{ course.Teacher }}</span>
            </div>
            <div class="meta-item">
              <span class="meta-icon">📚</span>
              <span class="meta-text">{{ course.Grade }} · {{ course.Semester }}</span>
            </div>
            <div class="meta-item">
              <span class="meta-icon">📖</span>
              <span class="meta-text">{{ course.Subject }}</span>
            </div>
          </div>
          
          <div class="course-tags">
            <span
              v-for="(tag, index) in getTags(course)"
              :key="index"
              :class="['course-tag', tag.type]"
            >
              {{ tag.text }}
            </span>
          </div>
          
          <div class="course-rating-summary">
            <div class="rating-stars">
              <span v-for="i in 5" :key="i" class="star">
                {{ i <= Math.floor(averageRating) ? '⭐' : (i - 0.5 <= averageRating ? '🌟' : '☆') }}
              </span>
            </div>
            <span class="rating-value">{{ averageRating }}</span>
            <span class="rating-count">({{ ratings.length }} 人评分)</span>
          </div>
          
          <!-- Action Buttons -->
          <div class="action-buttons">
            <RouterLink
              :to="`/courses/${course.ID}/rate`"
              class="btn btn-primary"
            >
              <span class="btn-icon">⭐</span>
              去评价课程
            </RouterLink>
            
            <button
              v-if="canSubmit"
              @click="submitEvaluationRequest"
              class="btn btn-secondary"
              :disabled="evaluationRequestLoading || hasEvaluationRequested"
            >
              <span class="btn-icon">📢</span>
              {{ evaluationRequestLoading ? '发送中...' : (hasEvaluationRequested ? '已求评价' : '求评价') }}
            </button>
          </div>
          
          <p v-if="hasEvaluationRequested" class="evaluation-request-success">
            已成功发起求评价请求，请耐心等待其他同学的评价！
          </p>
        </div>

        <!-- Course Description Card -->
        <div class="card course-description">
          <h2 class="card-title">课程介绍</h2>
          <p class="description-text">{{ course.Description }}</p>
        </div>

        <!-- Rating Section Card -->
        <div class="card rating-section">
          <h2 class="card-title">课程评分</h2>
          
          <div v-if="ratings.length === 0" class="empty-state">
            <p>暂无评分</p>
          </div>
          <div v-else>
            <!-- 平均评分展示 -->
            <div class="average-ratings">
              <div class="average-rating-item">
                <div class="rating-label">总体评分</div>
                <div class="rating-value-stars">
                  <div class="rating-stars">
                    <span v-for="i in 5" :key="i" class="star">
                      {{ i <= Math.floor(averageRating) ? '⭐' : (i - 0.5 <= averageRating ? '🌟' : '☆') }}
                    </span>
                  </div>
                  <span class="rating-number">{{ averageRating }}</span>
                </div>
              </div>
              
              <div class="average-rating-item">
                <div class="rating-label">课程难度</div>
                <div class="rating-value-stars">
                  <div class="rating-stars">
                    <span v-for="i in 5" :key="i" class="star">
                      {{ i <= Math.floor(averageDifficulty) ? '⭐' : (i - 0.5 <= averageDifficulty ? '🌟' : '☆') }}
                    </span>
                  </div>
                  <span class="rating-number">{{ averageDifficulty }}</span>
                </div>
              </div>
              
              <div class="average-rating-item">
                <div class="rating-label">实用性</div>
                <div class="rating-value-stars">
                  <div class="rating-stars">
                    <span v-for="i in 5" :key="i" class="star">
                      {{ i <= Math.floor(averageUsefulness) ? '⭐' : (i - 0.5 <= averageUsefulness ? '🌟' : '☆') }}
                    </span>
                  </div>
                  <span class="rating-number">{{ averageUsefulness }}</span>
                </div>
              </div>
              
              <div class="average-rating-item">
                <div class="rating-label">教学质量</div>
                <div class="rating-value-stars">
                  <div class="rating-stars">
                    <span v-for="i in 5" :key="i" class="star">
                      {{ i <= Math.floor(averageTeaching) ? '⭐' : (i - 0.5 <= averageTeaching ? '🌟' : '☆') }}
                    </span>
                  </div>
                  <span class="rating-number">{{ averageTeaching }}</span>
                </div>
              </div>
            </div>
            
            <!-- 用户评分列表 -->
            <h3 class="user-ratings-title">用户评分</h3>
            <div class="rating-items">
              <div v-for="rating in ratings" :key="rating.ID" class="rating-item">
                <div class="rating-user">
                  <span class="user-avatar">{{ rating.Username?.charAt(0) || 'U' }}</span>
                  <span class="user-name">{{ rating.Username }}</span>
                </div>
                <div class="rating-score">
                  <div class="rating-stars">
                    <span v-for="i in 5" :key="i" class="star">
                      {{ i <= Math.floor(rating.Score) ? '⭐' : (i - 0.5 <= rating.Score ? '🌟' : '☆') }}
                    </span>
                  </div>
                  <span class="score-value">{{ rating.Score }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Comment Section Card -->
        <div class="card comment-section">
          <h2 class="card-title">课程评论</h2>
          
          <div v-if="comments.length === 0" class="empty-state">
            <p>暂无评论</p>
          </div>
          <div v-else class="comment-items">
            <div v-for="comment in comments" :key="comment.ID" class="comment-item">
              <div class="comment-header">
                <div class="comment-user">
                  <span class="user-avatar">{{ comment.Username?.charAt(0) || 'U' }}</span>
                  <div class="user-info">
                    <span class="user-name">{{ comment.Username }}</span>
                    <span class="comment-date">{{ comment.CreatedAt }}</span>
                  </div>
                </div>
              </div>
              <div class="comment-content">
                <p>{{ comment.Content }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Notes Section (Placeholder) -->
      <div class="notes-section">
        <div class="card notes-card">
          <h2 class="card-title">课程笔记</h2>
          <div class="notes-placeholder">
            <p>笔记功能即将上线...</p>
            <div class="notes-icon">📝</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 主容器 */
.course-detail-container {
  min-height: 100vh;
  background-color: #FEF6F7;
  font-family: sans-serif;
  color: #1A1A1A;
  padding: 20px;
}

/* 返回按钮容器 */
.back-button-container {
  margin-bottom: 24px;
}

/* 加载状态 */
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px;
  gap: 16px;
}

.loading-text {
  font-size: 16px;
  color: #888888;
  font-weight: bold;
}

/* 课程内容 */
.course-content {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  gap: 24px;
}

.main-content {
  flex: 2;
}

.notes-section {
  flex: 1;
}

/* 卡片通用样式 */
.card {
  background-color: #FFFFFF;
  border-radius: 12px;
  border: 3px solid #000000;
  box-shadow: 5px 5px 0px 0px #000000;
  padding: 24px;
  margin-bottom: 24px;
}

/* 卡片标题 */
.card-title {
  font-size: 20px;
  font-weight: bold;
  margin: 0 0 16px 0;
  text-align: center;
}

/* 课程头部 */
.course-header {
  position: relative;
}

.course-credits-badge {
  background-color: #F7D074;
  color: #1A1A1A;
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: bold;
  display: inline-block;
  margin-bottom: 16px;
  border: 2px solid #000000;
}

.course-title {
  font-size: 28px;
  font-weight: bold;
  margin: 0 0 16px 0;
  text-align: center;
}

.course-meta {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  color: #1A1A1A;
}

.meta-icon {
  font-size: 18px;
}

.course-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 16px;
}

.course-tag {
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: bold;
  border: 2px solid #000000;
}

.course-tag.primary {
  background-color: #76D7C4;
  color: #1A1A1A;
}

.course-tag.secondary {
  background-color: #F7D074;
  color: #1A1A1A;
}

.course-tag.accent {
  background-color: #F7D074;
  color: #1A1A1A;
}

.course-rating-summary {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-bottom: 24px;
}

.rating-stars {
  display: flex;
  gap: 2px;
}

.star {
  font-size: 18px;
}

.rating-value {
  font-size: 20px;
  font-weight: bold;
  color: #1A1A1A;
}

.rating-count {
  font-size: 14px;
  color: #888888;
}

/* 按钮通用样式 */
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
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  text-decoration: none;
  transition: transform 0.2s ease;
}

.btn:hover {
  transform: translateY(-2px);
}

.btn:active {
  transform: translateY(0);
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background-color: #76D7C4;
}

/* 操作按钮区域 */
.action-buttons {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}

.evaluation-request-success {
  font-size: 14px;
  color: #76D7C4;
  font-weight: bold;
  text-align: center;
  margin-top: 8px;
}

/* 课程描述 */
.description-text {
  font-size: 16px;
  line-height: 1.6;
  color: #1A1A1A;
  margin: 0;
}

/* 评分和评论部分 */
.empty-state {
  text-align: center;
  padding: 24px;
  color: #888888;
  font-size: 16px;
}

/* 平均评分展示 */
.average-ratings {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 24px;
  padding: 16px;
  background-color: #FEF6F7;
  border-radius: 8px;
  border: 2px solid #000000;
}

.average-rating-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.rating-label {
  font-size: 16px;
  font-weight: bold;
  color: #1A1A1A;
}

.rating-value-stars {
  display: flex;
  align-items: center;
  gap: 8px;
}

.rating-number {
  font-size: 16px;
  font-weight: bold;
  color: #1A1A1A;
}

.user-ratings-title {
  font-size: 18px;
  font-weight: bold;
  margin: 0 0 16px 0;
  text-align: center;
}

.rating-items,
.comment-items {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.rating-item,
.comment-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background-color: #FEF6F7;
  border-radius: 8px;
  border: 2px solid #000000;
}

.comment-item {
  flex-direction: column;
  align-items: flex-start;
}

.rating-user,
.comment-user {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: #F7D074;
  color: #1A1A1A;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 18px;
  border: 2px solid #000000;
}

.user-name {
  font-size: 16px;
  font-weight: bold;
  color: #1A1A1A;
}

.rating-score {
  display: flex;
  align-items: center;
  gap: 8px;
}

.score-value {
  font-size: 16px;
  font-weight: bold;
  color: #1A1A1A;
}

.comment-header {
  width: 100%;
  margin-bottom: 8px;
}

.user-info {
  display: flex;
  flex-direction: column;
}

.comment-date {
  font-size: 14px;
  color: #888888;
}

.comment-content {
  font-size: 16px;
  line-height: 1.6;
  color: #1A1A1A;
  width: 100%;
}

.comment-content p {
  margin: 0;
}

/* 笔记部分 */
.notes-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
}

.notes-placeholder p {
  font-size: 16px;
  color: #888888;
  margin-bottom: 16px;
}

.notes-icon {
  font-size: 48px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .course-content {
    flex-direction: column;
  }
  
  .action-buttons {
    flex-direction: column;
  }
  
  .rating-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}
</style>