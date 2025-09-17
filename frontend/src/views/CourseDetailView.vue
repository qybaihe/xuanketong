<script setup lang="ts">
import { ref, onMounted, computed, nextTick } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { commentService } from '@/services/api'
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
  imageURL: string;
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
  ImageURL?: string;
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
    imageURL: course.imageURL || `https://picsum.photos/seed/course-${course.id}/800/400.jpg`,
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
    ImageURL: course.imageURL || `https://picsum.photos/seed/course-${course.id}/800/400.jpg`,
    CreatedAt: course.createdAt,
    UpdatedAt: course.updatedAt
  }
}

interface Rating {
  ID: number;
  UserID: number;
  CourseID: number;
  Score: number;
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
const newScore = ref(5)
const newComment = ref('')
const loading = ref(true)
const submitLoading = ref(false)

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
      Username: rating.username || rating.nickname || `用户${rating.UserID}`,
      Score: rating.score || rating.Score || 0
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
      Username: comment.username || comment.nickname || `用户${comment.UserID}`,
      CreatedAt: comment.createdAt ? new Date(comment.createdAt).toLocaleDateString() : new Date().toLocaleDateString()
    }))
  } catch (error) {
    console.error('获取评论失败:', error)
  }
}

// 计算平均评分
const averageRating = computed(() => {
  if (ratings.value.length === 0) return 0
  const sum = ratings.value.reduce((acc, rating) => acc + rating.Score, 0)
  return Number((sum / ratings.value.length).toFixed(1))
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

const submitRating = async () => {
  if (!authStore.isAuthenticated) {
    alert('请先登录后再评分！')
    return
  }
  
  if (!currentUser.value) {
    alert('无法获取用户信息，请重新登录！')
    return
  }
  
  submitLoading.value = true
  try {
    const courseId = Number(route.params.id)
    
    // 直接调用API创建评分
    await api.post('/ratings', {
      UserID: currentUser.value.id,
      CourseID: courseId,
      Score: newScore.value
    })
    
    newScore.value = 5
    alert('评分提交成功！')
    
    // 刷新评分列表
    await fetchRatings(courseId)
  } catch (error: any) {
    console.error('评分提交失败:', error)
    const errorMessage = error.response?.data?.error || '评分提交失败，请重试！'
    alert(errorMessage)
  } finally {
    submitLoading.value = false
  }
}

const submitComment = async () => {
  if (!authStore.isAuthenticated) {
    alert('请先登录后再评论！')
    return
  }
  
  if (!currentUser.value) {
    alert('无法获取用户信息，请重新登录！')
    return
  }
  
  if (!newComment.value.trim()) {
    alert('请输入评论内容！')
    return
  }
  
  submitLoading.value = true
  try {
    const courseId = Number(route.params.id)
    
    // 使用评论服务创建评论
    await commentService.createComment(
      currentUser.value.id,
      courseId,
      newComment.value.trim()
    )
    
    newComment.value = ''
    alert('评论提交成功！')
    
    // 刷新评论列表
    await fetchComments(courseId)
  } catch (error: any) {
    console.error('评论提交失败:', error)
    const errorMessage = error.response?.data?.error || '评论提交失败，请重试！'
    alert(errorMessage)
  } finally {
    submitLoading.value = false
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
      Username: (rating as any).username || (rating as any).nickname || `用户${rating.UserID}`
    }))
    // 使用后端返回的username，如果没有则使用nickname，如果都没有才使用默认
    comments.value = (commentsResponse.data.data || []).map((comment: Comment) => ({
      ...comment,
      Username: (comment as any).username || (comment as any).nickname || `用户${comment.UserID}`,
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
  <div class="course-detail-container" :class="themeClass">
    <!-- 返回按钮 -->
    <div class="back-button-container">
      <RouterLink to="/" class="btn btn-glass back-button">
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
      <!-- Course Header -->
      <div class="course-header card-glass">
        <div class="course-header-content">
          <div class="course-image-container">
            <img
              :src="course.ImageURL || course.imageURL || `https://picsum.photos/seed/course-${course?.ID || course?.id || 'default'}/800/400.jpg`"
              :alt="course.Name || course.name || '课程图片'"
              class="course-image"
              @error="(e) => {
                const img = e.target as HTMLImageElement;
                const courseId = course?.ID || course?.id || 'default';
                img.src = `https://picsum.photos/seed/course-${courseId}/800/400.jpg`;
              }"
            />
            <div class="course-credits-badge">
              {{ course.Credits }} 学分
            </div>
          </div>
          
          <div class="course-info">
            <div class="course-tags">
              <span
                v-for="(tag, index) in getTags(course)"
                :key="index"
                :class="['course-tag', tag.type]"
              >
                {{ tag.text }}
              </span>
            </div>
            
            <h1 class="course-title text-shine">{{ course.Name }}</h1>
            
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
            
            <div class="course-rating-summary">
              <div class="rating-stars">
                <span v-for="i in 5" :key="i" class="star">
                  {{ i <= Math.floor(averageRating) ? '⭐' : (i - 0.5 <= averageRating ? '🌟' : '☆') }}
                </span>
              </div>
              <span class="rating-value">{{ averageRating }}</span>
              <span class="rating-count">({{ ratings.length }} 人评分)</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Course Description -->
      <div class="course-description card-glass">
        <h2 class="section-title">课程介绍</h2>
        <p class="description-text">{{ course.Description }}</p>
      </div>

      <!-- Rating Section -->
      <div class="rating-section card-glass">
        <h2 class="section-title">课程评分</h2>
        
        <div class="rating-form">
          <h3>为这门课程评分</h3>
          <div class="score-input-container">
            <input
              type="range"
              v-model.number="newScore"
              min="1"
              max="5"
              step="0.5"
              class="score-slider"
            />
            <div class="score-display">{{ newScore }}</div>
          </div>
          <button
            @click="submitRating"
            class="btn btn-primary"
            :disabled="submitLoading || !canSubmit"
          >
            <span class="btn-icon">⭐</span>
            {{ submitLoading ? '提交中...' : '提交评分' }}
          </button>
          <div v-if="!canSubmit" class="login-prompt">
            <span class="prompt-icon">🔒</span>
            <span>请先登录后进行评分</span>
            <RouterLink to="/auth" class="login-link">立即登录</RouterLink>
          </div>
        </div>
        
        <div class="ratings-list">
          <h3>用户评分</h3>
          <div v-if="ratings.length === 0" class="empty-state">
            <p>暂无评分</p>
          </div>
          <div v-else class="rating-items">
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

      <!-- Comment Section -->
      <div class="comment-section card-glass">
        <h2 class="section-title">课程评论</h2>
        
        <div class="comment-form">
          <h3>发表评论</h3>
          <textarea
            v-model="newComment"
            placeholder="分享你对这门课程的看法..."
            class="comment-input"
            rows="4"
          ></textarea>
          <button
            @click="submitComment"
            class="btn btn-primary"
            :disabled="submitLoading || !canSubmit"
          >
            <span class="btn-icon">💬</span>
            {{ submitLoading ? '提交中...' : '提交评论' }}
          </button>
          <div v-if="!canSubmit" class="login-prompt">
            <span class="prompt-icon">🔒</span>
            <span>请先登录后发表评论</span>
            <RouterLink to="/auth" class="login-link">立即登录</RouterLink>
          </div>
        </div>
        
        <div class="comments-list">
          <h3>用户评论</h3>
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
    </div>
  </div>
</template>

<style scoped>
/* 主容器 */
.course-detail-container {
  min-height: 100vh;
  background: linear-gradient(135deg, var(--background-secondary) 0%, var(--natural-background) 100%);
  padding: var(--spacing-lg) 0;
}

/* 返回按钮容器 */
.back-button-container {
  max-width: 1200px;
  margin: 0 auto var(--spacing-lg);
  padding: 0 var(--spacing-md);
  display: flex;
  justify-content: flex-start;
  align-items: center;
}

.back-button {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

/* 加载状态 */
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
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  font-weight: var(--font-weight-semibold);
}

/* 课程内容 */
.course-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 var(--spacing-md);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xl);
}

/* 课程头部 */
.course-header {
  padding: var(--spacing-xl);
  border-radius: 16px;
  overflow: hidden;
}

.course-header-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--spacing-xl);
  align-items: center;
}

.course-image-container {
  position: relative;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: var(--shadow-lg);
}

.course-image {
  width: 100%;
  height: 300px;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.course-image-container:hover .course-image {
  transform: scale(1.05);
}

.course-credits-badge {
  position: absolute;
  top: var(--spacing-md);
  right: var(--spacing-md);
  background: var(--gradient-primary);
  color: white;
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: 20px;
  font-size: var(--font-size-caption);
  font-weight: var(--font-weight-medium);
  box-shadow: var(--shadow-primary);
}

.course-info {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.course-tags {
  display: flex;
  gap: var(--spacing-xs);
  flex-wrap: wrap;
}

.course-tag {
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: 4px;
  font-size: var(--font-size-caption);
  font-weight: var(--font-weight-medium);
}

.course-tag.primary {
  background: var(--success-light);
  color: var(--success-base);
}

.course-tag.secondary {
  background: var(--info-light);
  color: var(--info-base);
}

.course-tag.accent {
  background: var(--warning-light);
  color: var(--warning-base);
}

.course-title {
  font-size: var(--font-size-title1);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin: 0;
  line-height: 1.2;
}

.course-meta {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.meta-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  font-size: var(--font-size-body);
  color: var(--text-secondary);
}

.meta-icon {
  font-size: 18px;
}

.course-rating-summary {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.rating-stars {
  display: flex;
  gap: 2px;
}

.star {
  font-size: 16px;
}

.rating-value {
  font-size: var(--font-size-title3);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
}

.rating-count {
  font-size: var(--font-size-body2);
  color: var(--text-tertiary);
}

/* 课程描述 */
.course-description {
  padding: var(--spacing-xl);
  border-radius: 16px;
}

.section-title {
  font-size: var(--font-size-title2);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin: 0 0 var(--spacing-md) 0;
}

.description-text {
  font-size: var(--font-size-body);
  line-height: 1.6;
  color: var(--text-secondary);
  margin: 0;
}

/* 评分部分 */
.rating-section {
  padding: var(--spacing-xl);
  border-radius: 16px;
}

.rating-form {
  margin-bottom: var(--spacing-xl);
  padding: var(--spacing-lg);
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
}

.rating-form h3 {
  font-size: var(--font-size-title3);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin: 0 0 var(--spacing-md) 0;
}

.score-input-container {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-md);
}

.score-slider {
  flex: 1;
  height: 8px;
  border-radius: 4px;
  background: var(--background-tertiary);
  outline: none;
  -webkit-appearance: none;
}

.score-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--gradient-primary);
  cursor: pointer;
}

.score-slider::-moz-range-thumb {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--gradient-primary);
  cursor: pointer;
  border: none;
}

.score-display {
  font-size: var(--font-size-title2);
  font-weight: var(--font-weight-bold);
  color: var(--primary-color);
  min-width: 40px;
  text-align: center;
}

.ratings-list h3 {
  font-size: var(--font-size-title3);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin: 0 0 var(--spacing-md) 0;
}

.empty-state {
  text-align: center;
  padding: var(--spacing-lg);
  color: var(--text-tertiary);
}

.rating-items {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.rating-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-md);
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
}

.rating-user {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--gradient-primary);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: var(--font-weight-semibold);
}

.user-name {
  font-size: var(--font-size-body);
  color: var(--text-primary);
}

.rating-score {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.score-value {
  font-size: var(--font-size-body2);
  font-weight: var(--font-weight-semibold);
  color: var(--primary-color);
}

/* 评论部分 */
.comment-section {
  padding: var(--spacing-xl);
  border-radius: 16px;
}

.comment-form {
  margin-bottom: var(--spacing-xl);
  padding: var(--spacing-lg);
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
}

.comment-form h3 {
  font-size: var(--font-size-title3);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin: 0 0 var(--spacing-md) 0;
}

.comment-input {
  width: 100%;
  padding: var(--spacing-md);
  border: 1px solid var(--separator-color);
  border-radius: 8px;
  font-size: var(--font-size-body);
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  margin-bottom: var(--spacing-md);
  resize: vertical;
  min-height: 100px;
}

.comment-input:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(47, 169, 20, 0.1);
}

.comments-list h3 {
  font-size: var(--font-size-title3);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin: 0 0 var(--spacing-md) 0;
}

.comment-items {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.comment-item {
  padding: var(--spacing-lg);
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
}

.comment-header {
  margin-bottom: var(--spacing-md);
}

.comment-user {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.user-info {
  display: flex;
  flex-direction: column;
}

.comment-date {
  font-size: var(--font-size-caption);
  color: var(--text-tertiary);
}

.comment-content {
  font-size: var(--font-size-body);
  line-height: 1.6;
  color: var(--text-secondary);
}

.comment-content p {
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .course-header-content {
    grid-template-columns: 1fr;
  }
  
  .course-image {
    height: 200px;
  }
  
  .back-button-container {
    flex-direction: column;
    gap: var(--spacing-md);
  }
  
  .course-content {
    padding: 0 var(--spacing-sm);
  }
  
  .rating-item {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-sm);
  }
}

@media (min-width: 769px) and (max-width: 1024px) {
  .course-header-content {
    grid-template-columns: 1fr;
  }
  
  .course-image {
    height: 250px;
  }
}

/* 动画效果 */
.course-header,
.course-description,
.rating-section,
.comment-section {
  transition: all 0.3s ease;
}

.course-header:hover,
.course-description:hover,
.rating-section:hover,
.comment-section:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-xl);
}

/* 主题切换动画 */
.elastic-in {
  animation: elasticIn 0.6s cubic-bezier(0.68, -0.55, 0.265, 1.55);
}

@keyframes elasticIn {
  0% {
    transform: scale(0.8);
    opacity: 0;
  }
  50% {
    transform: scale(1.05);
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

/* 登录提示样式 */
.login-prompt {
  margin-top: var(--spacing-md);
  padding: var(--spacing-sm) var(--spacing-md);
  background: var(--background-secondary);
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  font-size: var(--font-size-caption);
  color: var(--text-tertiary);
}

.prompt-icon {
  font-size: 14px;
}

.login-link {
  color: var(--primary-color);
  text-decoration: none;
  font-weight: var(--font-weight-medium);
  margin-left: auto;
  transition: var(--transition-standard);
}

.login-link:hover {
  color: var(--primary-color-dark);
  text-decoration: underline;
}

/* 按钮禁用状态 */
.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none !important;
}

/* 暗色主题支持 */
@media (prefers-color-scheme: dark) {
  .course-detail-container {
    background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  }
  
  .course-title {
    background: linear-gradient(135deg, #4fc830 0%, #2fa914 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }
  
  .card-glass {
    background: rgba(0, 0, 0, 0.6);
    border: 1px solid rgba(255, 255, 255, 0.1);
  }
  
  .comment-input {
    background: rgba(0, 0, 0, 0.6);
    border: 1px solid rgba(255, 255, 255, 0.1);
    color: white;
  }

  .login-prompt {
    background: rgba(255, 255, 255, 0.1);
    color: rgba(255, 255, 255, 0.8);
  }

  .login-link {
    color: #4fc830;
  }

  .login-link:hover {
    color: #2fa914;
  }
}
</style>