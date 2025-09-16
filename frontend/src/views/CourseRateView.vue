<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const courseId = ref(route.params.id)
const loading = ref(true)
const submitting = ref(false)

// 课程信息
const course = reactive({
  id: courseId.value,
  name: '',
  teacher: '',
  credits: 0,
  grade: '',
  semester: '',
  subject: ''
})

// 评价表单
const ratingForm = reactive({
  overall: 0,
  difficulty: 0,
  usefulness: 0,
  teaching: 0,
  content: '',
  isAnonymous: false
})

// 预设评价选项
const ratingOptions = [
  { value: 5, label: '非常好', emoji: '⭐⭐⭐⭐⭐' },
  { value: 4, label: '很好', emoji: '⭐⭐⭐⭐' },
  { value: 3, label: '一般', emoji: '⭐⭐⭐' },
  { value: 2, label: '不太好', emoji: '⭐⭐' },
  { value: 1, label: '很不好', emoji: '⭐' }
]

// 获取课程信息
const fetchCourseInfo = async () => {
  loading.value = true
  try {
    // 模拟API调用
    const mockCourse = {
      id: courseId.value,
      name: '高等数学',
      teacher: '张教授',
      credits: 4,
      grade: '大一',
      semester: '2024春季',
      subject: '数学'
    }
    Object.assign(course, mockCourse)
  } catch (error) {
    console.error('获取课程信息失败:', error)
  } finally {
    loading.value = false
  }
}

// 处理星级点击
const handleStarClick = (ratingType: keyof typeof ratingForm, value: number) => {
  ratingForm[ratingType] = value
}

// 提交评价
const submitRating = async () => {
  if (ratingForm.overall === 0) {
    alert('请给出总体评分')
    return
  }
  
  if (ratingForm.content.trim().length < 10) {
    alert('请输入至少10个字的评价内容')
    return
  }

  submitting.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 显示成功提示
    alert('评价提交成功！感谢您的反馈。')
    
    // 返回课程详情页面
    router.push({ name: 'course-detail', params: { id: courseId.value } })
  } catch (error) {
    console.error('提交评价失败:', error)
    alert('提交失败，请稍后重试')
  } finally {
    submitting.value = false
  }
}

// 取消评价
const cancelRating = () => {
  if (confirm('确定要取消评价吗？已填写的内容将丢失。')) {
    router.push({ name: 'course-detail', params: { id: courseId.value } })
  }
}

onMounted(() => {
  fetchCourseInfo()
})
</script>

<template>
  <div class="course-rate-container">
    <!-- 顶部导航 -->
    <div class="rate-header">
      <div class="header-content">
        <button @click="cancelRating" class="btn-back">
          <span class="back-icon">←</span>
          返回
        </button>
        <h1 class="page-title">课程评价</h1>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <div class="loader"></div>
      <p class="loading-text">正在加载课程信息...</p>
    </div>

    <!-- 评价表单 -->
    <div v-else class="rate-content">
      <!-- 课程信息卡片 -->
      <div class="course-info-card">
        <div class="course-info">
          <h2 class="course-name">{{ course.name }}</h2>
          <div class="course-meta">
            <span class="meta-item">
              <span class="meta-icon">👨‍🏫</span>
              {{ course.teacher }}
            </span>
            <span class="meta-item">
              <span class="meta-icon">📚</span>
              {{ course.credits }} 学分
            </span>
            <span class="meta-item">
              <span class="meta-icon">🏷️</span>
              {{ course.grade }} · {{ course.semester }} · {{ course.subject }}
            </span>
          </div>
        </div>
      </div>

      <!-- 评价表单 -->
      <div class="rating-form">
        <h3 class="form-title">为这门课程打分</h3>
        <p class="form-subtitle">您的真实评价将帮助其他同学做出更好的选择</p>

        <!-- 总体评分 -->
        <div class="rating-section">
          <div class="section-header">
            <h4 class="section-title">总体评分 *</h4>
            <p class="section-desc">您对这门课程的总体满意度</p>
          </div>
          <div class="rating-stars">
            <span 
              v-for="i in 5" 
              :key="i"
              class="star"
              :class="{ active: i <= ratingForm.overall }"
              @click="handleStarClick('overall', i)"
            >
              {{ i <= ratingForm.overall ? '⭐' : '☆' }}
            </span>
          </div>
          <div class="rating-help">
            {{ ratingOptions.find(opt => opt.value === ratingForm.overall)?.label || '请选择评分' }}
          </div>
        </div>

        <!-- 详细评分 -->
        <div class="detailed-ratings">
          <div class="rating-section">
            <div class="section-header">
              <h4 class="section-title">课程难度</h4>
              <p class="section-desc">课程的难易程度</p>
            </div>
            <div class="rating-stars">
              <span 
                v-for="i in 5" 
                :key="i"
                class="star small"
                :class="{ active: i <= ratingForm.difficulty }"
                @click="handleStarClick('difficulty', i)"
              >
                {{ i <= ratingForm.difficulty ? '⭐' : '☆' }}
              </span>
            </div>
          </div>

          <div class="rating-section">
            <div class="section-header">
              <h4 class="section-title">实用性</h4>
              <p class="section-desc">课程内容对您是否有帮助</p>
            </div>
            <div class="rating-stars">
              <span 
                v-for="i in 5" 
                :key="i"
                class="star small"
                :class="{ active: i <= ratingForm.usefulness }"
                @click="handleStarClick('usefulness', i)"
              >
                {{ i <= ratingForm.usefulness ? '⭐' : '☆' }}
              </span>
            </div>
          </div>

          <div class="rating-section">
            <div class="section-header">
              <h4 class="section-title">教学质量</h4>
              <p class="section-desc">教师的教学水平和方法</p>
            </div>
            <div class="rating-stars">
              <span 
                v-for="i in 5" 
                :key="i"
                class="star small"
                :class="{ active: i <= ratingForm.teaching }"
                @click="handleStarClick('teaching', i)"
              >
                {{ i <= ratingForm.teaching ? '⭐' : '☆' }}
              </span>
            </div>
          </div>
        </div>

        <!-- 文字评价 -->
        <div class="text-rating">
          <div class="section-header">
            <h4 class="section-title">详细评价 *</h4>
            <p class="section-desc">分享您的学习体验和建议</p>
          </div>
          <textarea
            v-model="ratingForm.content"
            class="rating-textarea"
            placeholder="请详细描述您对这门课程的看法，包括课程内容、教学方法、考核方式等方面。您的评价将对其他同学有重要的参考价值..."
            rows="6"
          ></textarea>
          <div class="char-count">
            {{ ratingForm.content.length }}/500 字
          </div>
        </div>

        <!-- 匿名选项 -->
        <div class="anonymous-option">
          <label class="checkbox-label">
            <input
              type="checkbox"
              v-model="ratingForm.isAnonymous"
              class="checkbox"
            />
            <span class="checkmark"></span>
            <span class="label-text">匿名评价</span>
          </label>
          <p class="anonymous-help">选择匿名后，其他用户将看不到您的身份信息</p>
        </div>

        <!-- 提交按钮 -->
        <div class="form-actions">
          <button @click="cancelRating" class="btn btn-secondary">
            取消
          </button>
          <button 
            @click="submitRating" 
            class="btn btn-primary"
            :disabled="submitting"
          >
            <span v-if="submitting" class="btn-icon">⏳</span>
            <span v-else>提交评价</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 评价指南 -->
    <div class="rating-guide">
      <h3 class="guide-title">评价指南</h3>
      <div class="guide-content">
        <div class="guide-item">
          <span class="guide-icon">✓</span>
          <div class="guide-text">
            <strong>真实客观：</strong>基于您的真实学习体验进行评价
          </div>
        </div>
        <div class="guide-item">
          <span class="guide-icon">✓</span>
          <div class="guide-text">
            <strong>详细具体：</strong>提供具体的例子和详细描述
          </div>
        </div>
        <div class="guide-item">
          <span class="guide-icon">✓</span>
          <div class="guide-text">
            <strong>尊重他人：</strong>使用文明用语，尊重教师和其他同学
          </div>
        </div>
        <div class="guide-item">
          <span class="guide-icon">✓</span>
          <div class="guide-text">
            <strong>建设性：</strong>提出有价值的建议和改进意见
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* ===== 课程评价页面样式 ===== */
.course-rate-container {
  min-height: 100vh;
  background: var(--background-primary);
  padding-bottom: var(--spacing-3xl);
}

.rate-header {
  background: var(--background-secondary);
  border-bottom: 1px solid var(--separator-color);
  padding: var(--spacing-lg) 0;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 var(--spacing-lg);
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.btn-back {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  background: transparent;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  transition: var(--transition-standard);
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: 8px;
}

.btn-back:hover {
  background: var(--background-tertiary);
  color: var(--text-primary);
}

.back-icon {
  font-size: 18px;
  font-weight: var(--font-weight-bold);
}

.page-title {
  font-size: var(--font-size-title2);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin: 0;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-3xl);
  gap: var(--spacing-md);
}

.loader {
  width: 40px;
  height: 40px;
  border: 3px solid var(--background-secondary);
  border-top: 3px solid var(--primary-color);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.loading-text {
  font-size: var(--font-size-body);
  color: var(--text-secondary);
}

.rate-content {
  max-width: 800px;
  margin: 0 auto;
  padding: var(--spacing-lg);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xl);
}

.course-info-card {
  background: white;
  border-radius: 16px;
  padding: var(--spacing-xl);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--separator-color);
}

.course-info {
  text-align: center;
}

.course-name {
  font-size: var(--font-size-title3);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin: 0 0 var(--spacing-md) 0;
}

.course-meta {
  display: flex;
  justify-content: center;
  gap: var(--spacing-lg);
  flex-wrap: wrap;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  font-size: var(--font-size-body2);
  color: var(--text-secondary);
}

.meta-icon {
  font-size: 16px;
}

.rating-form {
  background: white;
  border-radius: 16px;
  padding: var(--spacing-xl);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--separator-color);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xl);
}

.form-title {
  font-size: var(--font-size-title2);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin: 0;
  text-align: center;
}

.form-subtitle {
  font-size: var(--font-size-body);
  color: var(--text-secondary);
  margin: 0 0 var(--spacing-xl) 0;
  text-align: center;
}

.rating-section {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.section-header {
  text-align: center;
}

.section-title {
  font-size: var(--font-size-body);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin: 0;
}

.section-desc {
  font-size: var(--font-size-body2);
  color: var(--text-secondary);
  margin: var(--spacing-xs) 0 0;
}

.rating-stars {
  display: flex;
  justify-content: center;
  gap: var(--spacing-sm);
}

.star {
  font-size: 32px;
  cursor: pointer;
  transition: var(--transition-standard);
  color: var(--text-tertiary);
}

.star:hover {
  transform: scale(1.1);
  color: var(--warning-color);
}

.star.active {
  color: var(--warning-color);
  animation: starPulse 0.3s ease;
}

.star.small {
  font-size: 24px;
}

@keyframes starPulse {
  0% { transform: scale(1); }
  50% { transform: scale(1.2); }
  100% { transform: scale(1); }
}

.rating-help {
  text-align: center;
  font-size: var(--font-size-body2);
  color: var(--primary-color);
  font-weight: var(--font-weight-medium);
  margin-top: var(--spacing-xs);
}

.detailed-ratings {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: var(--spacing-xl);
  padding: var(--spacing-xl) 0;
  border-top: 1px solid var(--separator-color);
  border-bottom: 1px solid var(--separator-color);
}

.text-rating {
  padding: var(--spacing-xl) 0;
}

.rating-textarea {
  width: 100%;
  padding: var(--spacing-md);
  border: 2px solid var(--separator-color);
  border-radius: 12px;
  font-size: var(--font-size-body);
  font-family: inherit;
  resize: vertical;
  min-height: 120px;
  transition: var(--transition-standard);
  background: var(--background-primary);
}

.rating-textarea:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 4px rgba(47, 169, 20, 0.1);
}

.char-count {
  text-align: right;
  font-size: var(--font-size-caption);
  color: var(--text-tertiary);
  margin-top: var(--spacing-xs);
}

.anonymous-option {
  padding: var(--spacing-lg) 0;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  cursor: pointer;
  user-select: none;
}

.checkbox {
  display: none;
}

.checkmark {
  width: 20px;
  height: 20px;
  border: 2px solid var(--separator-color);
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: var(--transition-standard);
  position: relative;
}

.checkbox:checked + .checkmark {
  background: var(--primary-color);
  border-color: var(--primary-color);
}

.checkbox:checked + .checkmark::after {
  content: '✓';
  color: white;
  font-size: 12px;
  font-weight: var(--font-weight-bold);
}

.label-text {
  font-size: var(--font-size-body);
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
}

.anonymous-help {
  font-size: var(--font-size-body2);
  color: var(--text-tertiary);
  margin: var(--spacing-xs) 0 0 0;
}

.form-actions {
  display: flex;
  justify-content: center;
  gap: var(--spacing-md);
  padding-top: var(--spacing-xl);
}

.btn {
  padding: var(--spacing-md) var(--spacing-xl);
  border-radius: 12px;
  font-size: var(--font-size-body);
  font-weight: var(--font-weight-semibold);
  cursor: pointer;
  transition: var(--transition-standard);
  border: none;
  min-width: 120px;
}

.btn-secondary {
  background: var(--background-secondary);
  color: var(--text-primary);
}

.btn-secondary:hover {
  background: var(--background-tertiary);
  transform: translateY(-1px);
}

.btn-primary {
  background: var(--gradient-primary);
  color: white;
}

.btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(47, 169, 20, 0.3);
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.btn-icon {
  font-size: 16px;
}

.rating-guide {
  max-width: 800px;
  margin: 0 auto;
  padding: var(--spacing-lg);
  background: rgba(47, 169, 20, 0.05);
  border-radius: 16px;
  border: 1px solid rgba(47, 169, 20, 0.2);
}

.guide-title {
  font-size: var(--font-size-title3);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin: 0 0 var(--spacing-lg) 0;
  text-align: center;
}

.guide-content {
  display: grid;
  gap: var(--spacing-md);
}

.guide-item {
  display: flex;
  gap: var(--spacing-sm);
  align-items: flex-start;
}

.guide-icon {
  font-size: 18px;
  color: var(--success-base);
  margin-top: 2px;
}

.guide-text {
  flex: 1;
  font-size: var(--font-size-body);
  color: var(--text-secondary);
  line-height: 1.5;
}

.guide-text strong {
  color: var(--text-primary);
  font-weight: var(--font-weight-semibold);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .header-content {
    padding: 0 var(--spacing-md);
  }
  
  .rate-content {
    padding: var(--spacing-md);
  }
  
  .course-info-card,
  .rating-form {
    padding: var(--spacing-lg);
  }
  
  .detailed-ratings {
    grid-template-columns: 1fr;
    gap: var(--spacing-lg);
    padding: var(--spacing-lg) 0;
  }
  
  .star {
    font-size: 28px;
  }
  
  .star.small {
    font-size: 20px;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .btn {
    width: 100%;
  }
}
</style>