<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import axios from 'axios'

interface User {
  ID: number;
  Username: string;
  Nickname: string;
  Email: string;
  Avatar?: string;
  JoinDate?: string;
}

interface Rating {
  ID: number;
  CourseID: number;
  Score: number;
  CourseName?: string;
  CourseImage?: string;
  CreatedAt?: string;
}

interface Comment {
  ID: number;
  CourseID: number;
  Content: string;
  CourseName?: string;
  CourseImage?: string;
  CreatedAt?: string;
}

const route = useRoute()
const user = ref<User | null>(null)
const ratings = ref<Rating[]>([])
const comments = ref<Comment[]>([])
const loading = ref(true)

// 只使用natural主题
const themeClass = computed(() => 'theme-natural')

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

// 获取用户名首字母作为头像
const getUserInitial = (name: string) => {
  return name.charAt(0).toUpperCase()
}

// 格式化日期
const formatDate = (dateString?: string) => {
  if (!dateString) return '未知时间'
  return new Date(dateString).toLocaleDateString()
}

onMounted(async () => {
  try {
    loading.value = true
    // Assuming the user ID is 1 for now
    const userId = 1
    const response = await axios.get(`http://localhost:8080/api/v1/users/${userId}`)
    user.value = response.data.user
    
    // 为评分和评论添加模拟数据
    ratings.value = response.data.ratings.map((rating: Rating) => ({
      ...rating,
      CourseName: `课程${rating.CourseID}`,
      CourseImage: `https://picsum.photos/seed/course-${rating.CourseID}/100/100.jpg`,
      CreatedAt: new Date().toLocaleDateString()
    }))
    
    comments.value = response.data.comments.map((comment: Comment) => ({
      ...comment,
      CourseName: `课程${comment.CourseID}`,
      CourseImage: `https://picsum.photos/seed/course-${comment.CourseID}/100/100.jpg`,
      CreatedAt: new Date().toLocaleDateString()
    }))
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="profile-container" :class="themeClass">
    <!-- 返回按钮 -->
    <div class="back-button-container">
      <RouterLink to="/" class="btn btn-glass back-button">
        <span class="btn-icon">←</span>
        返回首页
      </RouterLink>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="loading-container">
      <div class="loader"></div>
      <p class="loading-text">正在加载个人资料...</p>
    </div>

    <!-- Profile Content -->
    <div v-else-if="user" class="profile-content">
      <!-- User Header -->
      <div class="user-header card-glass">
        <div class="user-avatar-container">
          <div class="user-avatar">
            {{ getUserInitial(user.Nickname || user.Username) }}
          </div>
        </div>
        
        <div class="user-info">
          <h1 class="user-name text-shine">{{ user.Nickname || user.Username }}</h1>
          <p class="user-username">@{{ user.Username }}</p>
          
          <div class="user-meta">
            <div class="meta-item">
              <span class="meta-icon">✉️</span>
              <span class="meta-text">{{ user.Email }}</span>
            </div>
            <div class="meta-item">
              <span class="meta-icon">📅</span>
              <span class="meta-text">加入时间: {{ formatDate(user.JoinDate) }}</span>
            </div>
          </div>
          
          <div class="user-stats">
            <div class="stat-item">
              <span class="stat-value">{{ ratings.length }}</span>
              <span class="stat-label">评分</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ comments.length }}</span>
              <span class="stat-label">评论</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Ratings Section -->
      <div class="ratings-section card-glass">
        <h2 class="section-title">我的评分</h2>
        
        <div v-if="ratings.length === 0" class="empty-state">
          <div class="empty-icon">⭐</div>
          <h3 class="empty-title">暂无评分</h3>
          <p class="empty-description">您还没有对任何课程进行评分</p>
        </div>
        
        <div v-else class="ratings-list">
          <div v-for="rating in ratings" :key="rating.ID" class="rating-item">
            <div class="rating-course">
              <img
                :src="rating.CourseImage || `https://picsum.photos/seed/course-${rating.CourseID}/100/100.jpg`"
                :alt="rating.CourseName"
                class="course-image"
              />
              <div class="course-info">
                <h3 class="course-name">{{ rating.CourseName }}</h3>
                <p class="course-date">{{ formatDate(rating.CreatedAt) }}</p>
              </div>
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

      <!-- Comments Section -->
      <div class="comments-section card-glass">
        <h2 class="section-title">我的评论</h2>
        
        <div v-if="comments.length === 0" class="empty-state">
          <div class="empty-icon">💬</div>
          <h3 class="empty-title">暂无评论</h3>
          <p class="empty-description">您还没有对任何课程进行评论</p>
        </div>
        
        <div v-else class="comments-list">
          <div v-for="comment in comments" :key="comment.ID" class="comment-item">
            <div class="comment-header">
              <div class="comment-course">
                <img
                  :src="comment.CourseImage || `https://picsum.photos/seed/course-${comment.CourseID}/100/100.jpg`"
                  :alt="comment.CourseName"
                  class="course-image"
                />
                <div class="course-info">
                  <h3 class="course-name">{{ comment.CourseName }}</h3>
                  <p class="course-date">{{ formatDate(comment.CreatedAt) }}</p>
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
</template>

<style scoped>
/* 主容器 */
.profile-container {
  min-height: 100vh;
  background: linear-gradient(135deg, var(--background-secondary) 0%, var(--natural-background) 100%);
  padding: var(--spacing-lg) 0;
}

/* 返回按钮容器 */
.back-button-container {
  max-width: 1200px;
  margin: 0 auto var(--spacing-lg);
  padding: 0 var(--spacing-md);
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

/* 个人资料内容 */
.profile-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 var(--spacing-md);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xl);
}

/* 用户头部 */
.user-header {
  padding: var(--spacing-xl);
  border-radius: 16px;
  overflow: hidden;
}

.user-header-content {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: var(--spacing-xl);
  align-items: center;
}

.user-avatar-container {
  position: relative;
}

.user-avatar {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  background: var(--gradient-primary);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48px;
  font-weight: var(--font-weight-bold);
  box-shadow: var(--shadow-primary);
}

.user-info {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.user-name {
  font-size: var(--font-size-title1);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin: 0;
  line-height: 1.2;
}

.user-username {
  font-size: var(--font-size-body);
  color: var(--text-secondary);
  margin: 0;
}

.user-meta {
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

.user-stats {
  display: flex;
  gap: var(--spacing-lg);
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-value {
  font-size: var(--font-size-title2);
  font-weight: var(--font-weight-semibold);
  color: var(--primary-color);
}

.stat-label {
  font-size: var(--font-size-caption);
  color: var(--text-tertiary);
}

/* 评分部分 */
.ratings-section {
  padding: var(--spacing-xl);
  border-radius: 16px;
}

.section-title {
  font-size: var(--font-size-title2);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin: 0 0 var(--spacing-md) 0;
}

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
  font-size: var(--font-size-title3);
  font-weight: var(--font-weight-semibold);
  margin: 0 0 var(--spacing-sm) 0;
}

.empty-description {
  font-size: var(--font-size-body);
}

.ratings-list {
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
  border-radius: 12px;
  transition: all 0.3s ease;
}

.rating-item:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.rating-course {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.course-image {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  object-fit: cover;
}

.course-info {
  display: flex;
  flex-direction: column;
}

.course-name {
  font-size: var(--font-size-body);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin: 0 0 var(--spacing-xs) 0;
}

.course-date {
  font-size: var(--font-size-caption);
  color: var(--text-tertiary);
  margin: 0;
}

.rating-score {
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

.score-value {
  font-size: var(--font-size-body2);
  font-weight: var(--font-weight-semibold);
  color: var(--primary-color);
}

/* 评论部分 */
.comments-section {
  padding: var(--spacing-xl);
  border-radius: 16px;
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.comment-item {
  padding: var(--spacing-lg);
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  transition: all 0.3s ease;
}

.comment-item:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.comment-header {
  margin-bottom: var(--spacing-md);
}

.comment-course {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
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
  .user-header-content {
    grid-template-columns: 1fr;
    text-align: center;
  }
  
  .user-avatar {
    margin: 0 auto;
  }
  
  .user-meta {
    align-items: center;
  }
  
  .user-stats {
    justify-content: center;
  }
  
  .rating-item {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-sm);
  }
  
  .profile-content {
    padding: 0 var(--spacing-sm);
  }
}

@media (min-width: 769px) and (max-width: 1024px) {
  .user-header-content {
    grid-template-columns: auto 1fr;
  }
}

/* 动画效果 */
.user-header,
.ratings-section,
.comments-section {
  transition: all 0.3s ease;
}

.user-header:hover,
.ratings-section:hover,
.comments-section:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-xl);
}

/* 暗色主题支持 */
@media (prefers-color-scheme: dark) {
  .profile-container {
    background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  }
  
  .user-name {
    background: linear-gradient(135deg, #4fc830 0%, #2fa914 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }
  
  .card-glass {
    background: rgba(0, 0, 0, 0.6);
    border: 1px solid rgba(255, 255, 255, 0.1);
  }
}
</style>