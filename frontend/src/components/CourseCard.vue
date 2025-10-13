<template>
  <router-link :to="`/courses/${course.ID}`" class="course-card">
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
          <div class="teacher-credits-badge">{{ course.Credits }} 学分</div>
          <div class="teacher-info">
            <div class="teacher-details">
              <div class="teacher-name">{{ course.Teacher }}</div>
              <div v-if="course.CourseUnitName" class="unit-name">{{ course.CourseUnitName }}</div>
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
      <div v-if="showComments" class="hot-comments-section">
        <div class="hot-comments-header">
          <svg class="comments-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" stroke="var(--color-text)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          <span class="comments-title">热门评论</span>
        </div>
        <div v-if="course.HotComments && course.HotComments.length > 0" class="hot-comments-list">
          <div v-for="(comment, index) in course.HotComments.slice(0, 2)" :key="index" class="hot-comment">
            {{ comment }}
          </div>
        </div>
        <div v-else class="no-comments">
          <span class="no-comments-text">暂无评论，快去分享你的学习体验吧！</span>
        </div>
      </div>

      <!-- 资料信息区域 -->
      <div v-if="showMaterials" class="materials-section">
        <div class="materials-info">
          <span class="materials-count">
            <svg class="materials-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20" stroke="var(--color-text)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z" stroke="var(--color-text)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            {{ course.MaterialsCount || 0 }} 份资料
          </span>
          <div v-if="course.HotMaterials && course.HotMaterials.length > 0" class="hot-materials">
            <span class="hot-materials-label">热门：</span>
            <span v-for="(material, index) in course.HotMaterials.slice(0, 2)" :key="index" class="hot-material-tag">
              {{ material }}
            </span>
          </div>
          <div v-else-if="(course.MaterialsCount || 0) === 0" class="no-materials">
            <span class="no-materials-text">暂无资料，快去分享学习资料吧！</span>
          </div>
        </div>
      </div>

      <!-- 课程操作区域 -->
      <div class="course-card-actions">
        <button class="btn-rate-course" @click.prevent="goToRateCourse(course.ID)">
          评价课程
        </button>
        <button v-if="showUploadButton" class="btn-upload-materials" @click.prevent="goToUploadMaterials(course.ID)">
          上传资料
        </button>
      </div>
    </div>
  </router-link>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

interface CourseWithDisplay {
  ID: number
  Name: string
  Description?: string
  Teacher: string
  Credits: number
  Grade?: string
  Semester?: string
  Subject?: string
  CourseUnitName?: string
  CreatedAt?: string
  UpdatedAt?: string
  AverageRating?: number
  TotalRatings?: number
  RatingDistribution?: Record<number, number>
  MaterialsCount?: number
  HotComments?: string[]
  HotMaterials?: string[]
  AverageDifficulty?: number
  Difficulty?: 'easy' | 'medium' | 'hard'
  Workload?: 'light' | 'medium' | 'heavy'
}

interface Props {
  course: CourseWithDisplay
  showUploadButton?: boolean
  showComments?: boolean
  showMaterials?: boolean
}

withDefaults(defineProps<Props>(), {
  showUploadButton: true,
  showComments: true,
  showMaterials: true
})

const router = useRouter()

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
</script>

<style scoped>
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
  gap: 12px;
}

.course-card-tags {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
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
  line-height: 1.2;
}

.course-card-description {
  color: var(--color-textSecondary);
  margin-bottom: 16px;
}

.rating-count {
  color: var(--color-textSecondary);
}

/* ===== 课程评价功能样式 ===== */
.course-card-header {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 2px solid var(--color-border);
}

/* 授课老师部分 */
.course-teacher-section {
  position: relative;
  padding: 12px 16px;
  background-color: var(--color-surface);
  border-radius: 8px;
  border: 1px solid var(--color-border);
  transition: none;
}

.course-teacher-section:hover {
}

.teacher-credits-badge {
  position: absolute;
  top: 8px;
  right: 12px;
  font-size: 12px;
  color: var(--color-text);
  font-weight: bold;
  background-color: var(--color-primary);
  padding: 4px 8px;
  border-radius: 6px;
  border: 1px solid var(--color-border);
  white-space: nowrap;
}

.teacher-info {
  display: flex;
  align-items: flex-start;
  flex-grow: 1;
  padding-right: 60px; /* 为学分徽章留出空间 */
}

.teacher-details {
  display: flex;
  flex-direction: column;
  gap: 4px;
  width: 100%;
}

.teacher-name {
  font-size: 16px;
  font-weight: bold;
  color: var(--color-text);
  line-height: 1.3;
}

.unit-name {
  font-size: 14px;
  color: var(--color-textSecondary);
  font-weight: 500;
  line-height: 1.3;
}


/* 评价统计区域 */
.course-rating-section {
  padding: 12px;
  background-color: var(--color-surface);
  border-radius: 8px;
  border: 1px solid var(--color-border);
  margin-bottom: 12px;
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
  gap: 8px;
}

.rating-main {
  display: flex;
  align-items: center;
  gap: 8px;
  justify-content: center;
}

.rating-average {
  text-align: center;
}

.rating-number {
  font-size: 28px;
  font-weight: bold;
  color: var(--color-primary);
  line-height: 1;
  display: block;
}

.rating-stars {
  display: flex;
  gap: 2px;
  margin: 4px 0;
}

.rating-count {
  font-size: 12px;
  color: var(--color-textSecondary);
  white-space: nowrap;
}

/* 课程属性标签 */
.course-attributes {
  display: flex;
  gap: 8px;
  justify-content: center;
  flex-wrap: wrap;
}

.attribute-tag {
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: bold;
  border: 1px solid var(--color-border);
}

.attribute-tag.difficulty.easy {
  background-color: #76D7C4;
  color: var(--color-text);
}

.attribute-tag.difficulty.medium {
  background-color: #F7D074;
  color: var(--color-text);
}

.attribute-tag.difficulty.hard {
  background-color: #FF6B6B;
  color: var(--color-text);
}

.attribute-tag.workload.light {
  background-color: #76D7C4;
  color: var(--color-text);
}

.attribute-tag.workload.medium {
  background-color: #F7D074;
  color: var(--color-text);
}

.attribute-tag.workload.heavy {
  background-color: #FF6B6B;
  color: var(--color-text);
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
  color: var(--color-text);
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
  padding-left: 4px;
  border-left: 1px solid var(--color-border);
}

.no-comments {
  text-align: center;
  padding: 8px 0;
}

.no-comments-text {
  font-size: 12px;
  color: var(--color-textSecondary);
  font-style: italic;
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
  color: var(--color-text);
}

.hot-materials {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.hot-materials-label {
  font-size: 12px;
  color: var(--color-textSecondary);
  font-weight: bold;
}

.hot-material-tag {
  font-size: 11px;
  color: var(--color-text);
  background-color: var(--color-primary);
  padding: 2px 6px;
  border-radius: 4px;
  border: 1px solid var(--color-border);
}

.no-materials {
  text-align: center;
  padding: 4px 0;
}

.no-materials-text {
  font-size: 12px;
  color: var(--color-textSecondary);
  font-style: italic;
}

.course-card-actions {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: auto;
  padding-top: 12px;
  gap: 10px;
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

.btn-rate-course:hover {
}

.btn-upload-materials {
  background-color: var(--color-primary);
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

.btn-upload-materials:hover {
}

/* 评分星级样式 */
.star {
  font-size: 16px;
  transition: transform 0.2s ease;
  position: relative;
  display: inline-block;
  width: 16px;
  height: 16px;
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
  /* 移动端课程卡片优化布局 */
  .course-card {
    border: 1px solid var(--color-border);
    border-radius: 12px;
    box-shadow: none;
  }

  .course-card-content {
    padding: 18px;
    gap: 14px;
  }

  .course-card-tags {
    margin-bottom: 12px;
    gap: 8px;
  }

  .course-card-tag {
    padding: 4px 8px;
    font-size: 12px;
    border-radius: 6px;
  }

  .course-card-title {
    font-size: 19px;
    margin-bottom: 12px;
    line-height: 1.2;
  }

  .course-card-header {
    gap: 12px;
    margin-bottom: 12px;
    padding-bottom: 12px;
  }

  .course-teacher-section {
    padding: 12px 14px;
  }

  .teacher-credits-badge {
    top: 6px;
    right: 10px;
    font-size: 11px;
    padding: 3px 6px;
  }

  .teacher-info {
    padding-right: 50px;
  }

  .teacher-name {
    font-size: 15px;
  }

  .unit-name {
    font-size: 13px;
  }

  .course-rating-section {
    gap: 8px;
    padding: 12px;
    margin-bottom: 12px;
  }

  .rating-number {
    font-size: 26px;
  }

  .rating-stars {
    gap: 2px;
  }

  .star {
    font-size: 15px;
    width: 15px;
    height: 15px;
  }

  .course-attributes {
    gap: 8px;
  }

  .attribute-tag {
    padding: 4px 8px;
    font-size: 11px;
  }

  .course-card-actions {
    padding-top: 12px;
    gap: 12px;
  }

  .btn-rate-course,
  .btn-upload-materials {
    padding: 12px 18px;
    font-size: 14px;
    border-radius: 8px;
  }
}

/* 移动端课程卡片平衡排版 */
@media (max-width: 600px) {
  .course-card-content {
    padding: 20px;
    gap: 16px;
  }

  .course-card-tags {
    margin-bottom: 8px;
    gap: 8px;
  }

  .course-card-tag {
    padding: 4px 8px;
    font-size: 12px;
    font-weight: bold;
  }

  .course-card-title {
    font-size: 18px;
    margin-bottom: 12px;
    line-height: 1.2;
    font-weight: bold;
  }

  /* 紧凑教师信息 */
  .course-teacher-section {
    border: none;
    background: transparent;
    margin-bottom: 12px;
    align-items: center;
    gap: 8px;
    padding: 8px 0;
    border-bottom: 1px solid var(--color-border);
  }

  .teacher-label {
    font-size: 10px;
    padding: 2px 4px;
    background-color: var(--color-primary);
    border-radius: 4px;
    color: var(--color-text);
  }

  .teacher-name {
    font-size: 14px;
    font-weight: bold;
  }

  .teacher-credits {
    font-size: 12px;
    padding: 2px 6px;
    background-color: var(--color-accent);
    border-radius: 4px;
  }

  /* 紧凑评分区域 */
  .course-rating-section {
    border: none;
    background: transparent;
    margin-bottom: 8px;
    padding: 0;
  }

  .rating-overview-compact {
    gap: 6px;
  }

  .rating-main {
    gap: 6px;
  }

  .rating-number {
    font-size: 20px;
  }

  .rating-stars {
    gap: 1px;
  }

  .star {
    font-size: 12px;
    width: 12px;
    height: 12px;
  }

  .rating-count {
    font-size: 11px;
  }

  .course-attributes {
    gap: 4px;
  }

  .attribute-tag {
    padding: 2px 6px;
    font-size: 9px;
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
    font-size: 12px;
    border-radius: 4px;
    flex: none;
    min-width: 80px;
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

  .hot-comments-header {
    margin-bottom: 4px;
  }

  .comments-title {
    font-size: 12px;
  }

  .hot-comment {
    font-size: 11px;
  }

  .materials-count {
    font-size: 12px;
  }

  .hot-material-tag {
    font-size: 10px;
    padding: 1px 4px;
  }
}

/* 超小屏幕设备优化 (iPhone SE等) */
@media (max-width: 480px) {
  .course-card-content {
    padding: 16px;
    gap: 12px;
  }

  .course-card-title {
    font-size: 16px;
    margin-bottom: 8px;
    line-height: 1.2;
  }

  .course-card-tags {
    margin-bottom: 8px;
    gap: 6px;
  }

  .course-card-tag {
    padding: 3px 6px;
    font-size: 11px;
  }

  .course-card-header {
    gap: 8px;
    margin-bottom: 8px;
    padding-bottom: 8px;
  }

  .course-teacher-section {
    padding: 10px 0;
    margin-bottom: 10px;
  }

  .teacher-credits-badge {
    top: 5px;
    right: 8px;
    font-size: 10px;
    padding: 2px 4px;
  }

  .teacher-info {
    padding-right: 40px;
  }

  .teacher-name {
    font-size: 13px;
  }

  .unit-name {
    font-size: 11px;
  }


  .course-rating-section {
    padding: 6px;
  }

  .rating-number {
    font-size: 18px;
  }

  .star {
    font-size: 11px;
    width: 11px;
    height: 11px;
  }

  .rating-count {
    font-size: 10px;
  }

  .attribute-tag {
    padding: 2px 4px;
    font-size: 8px;
  }

  .hot-comments-section {
    padding: 8px;
  }

  .comments-title {
    font-size: 11px;
  }

  .hot-comment {
    font-size: 10px;
  }

  .materials-section {
    padding: 8px;
  }

  .materials-count {
    font-size: 11px;
  }

  .hot-material-tag {
    font-size: 9px;
    padding: 1px 3px;
  }

  .course-card-actions {
    gap: 4px;
  }

  .btn-rate-course,
  .btn-upload-materials {
    padding: 6px 10px;
    font-size: 11px;
    min-width: 70px;
  }
}
</style>
