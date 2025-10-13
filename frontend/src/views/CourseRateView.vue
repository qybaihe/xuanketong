<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/services/api'
import { marked } from 'marked'
// 移除图标组件导入，使用内联SVG

const route = useRoute()
const router = useRouter()
const courseId = ref(route.params.id)
const loading = ref(true)
const submitting = ref(false)

// 编辑模式
const isEditMode = ref(false)
const editRatingId = ref<number | null>(null)

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

// 预览模式
const showPreview = ref(false)

// Markdown解析函数
const parseMarkdown = (content: string): string => {
  if (!content) return ''
  return marked(content, {
    breaks: true,
    gfm: true
  })
}

// 预设评价选项
const ratingOptions = [
  { value: 5, label: '非常好', emoji: '★★★★★' },
  { value: 4, label: '很好', emoji: '★★★★☆' },
  { value: 3, label: '一般', emoji: '★★★☆☆' },
  { value: 2, label: '不太好', emoji: '★★☆☆☆' },
  { value: 1, label: '很不好', emoji: '★☆☆☆☆' }
]

// 获取课程信息
const fetchCourseInfo = async () => {
  loading.value = true
  try {
    const response = await api.get(`/courses/${courseId.value}`)
    const courseData = response.data.data
    Object.assign(course, {
      id: courseData.id,
      name: courseData.name,
      teacher: courseData.teacher,
      credits: courseData.credits,
      grade: courseData.grade,
      semester: courseData.semester,
      subject: courseData.subject
    })
  } catch (error) {
    console.error('获取课程信息失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取评价数据（编辑模式）
const fetchRating = async (ratingId: number) => {
  try {
    const response = await api.get(`/ratings/${ratingId}`)
    const ratingData = response.data.data
    Object.assign(ratingForm, {
      overall: ratingData.score,
      difficulty: ratingData.difficulty,
      usefulness: ratingData.usefulness,
      teaching: ratingData.teaching,
      content: ratingData.content,
      isAnonymous: false
    })
  } catch (error) {
    console.error('获取评价数据失败:', error)
    alert('获取评价数据失败，请重试')
  }
}

// 处理星级点击
const handleStarClick = (ratingType: string, value: number) => {
  // 使用类型断言来允许动态属性访问
  (ratingForm as any)[ratingType] = value
}

// 提交评价
const submitRating = async () => {
  if (ratingForm.overall === 0) {
    alert('请给出总体评分')
    return
  }
  
  if (ratingForm.difficulty === 0) {
    alert('请给出课程难度评分')
    return
  }
  
  if (ratingForm.usefulness === 0) {
    alert('请给出实用性评分')
    return
  }
  
  if (ratingForm.teaching === 0) {
    alert('请给出教学质量评分')
    return
  }
  
  if (ratingForm.content.trim().length < 10) {
    alert('请输入至少10个字的评价内容')
    return
  }

  submitting.value = true
  try {
    let response
    if (isEditMode.value && editRatingId.value) {
      // 更新现有评价
      response = await api.put(`/ratings/${editRatingId.value}`, {
        score: ratingForm.overall,
        difficulty: ratingForm.difficulty,
        usefulness: ratingForm.usefulness,
        teaching: ratingForm.teaching,
        content: ratingForm.content,
        isAnonymous: ratingForm.isAnonymous
      })
    } else {
      // 创建新评价
      response = await api.post(`/courses/${courseId.value}/ratings`, {
        score: ratingForm.overall,
        difficulty: ratingForm.difficulty,
        usefulness: ratingForm.usefulness,
        teaching: ratingForm.teaching,
        content: ratingForm.content,
        isAnonymous: ratingForm.isAnonymous
      })
    }
    
    // 显示成功提示
    alert(isEditMode.value ? '评价更新成功！感谢您的反馈。' : '评价提交成功！感谢您的反馈。')
    
    // 返回课程详情页面
    router.push({ name: 'course-detail', params: { id: courseId.value } })
  } catch (error: any) {
    console.error('提交评价失败:', error)
    const errorMessage = error.response?.data?.error || '提交失败，请稍后重试'
    alert(errorMessage)
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
  
  // 检查是否是编辑模式
  const editParam = route.query.edit
  if (editParam && typeof editParam === 'string') {
    const ratingId = parseInt(editParam)
    if (!isNaN(ratingId)) {
      isEditMode.value = true
      editRatingId.value = ratingId
      fetchRating(ratingId)
    }
  }
})
</script>

<template>
  <div class="course-rate-container">
    <!-- 返回按钮 -->
    <div class="back-button-container">
      <button @click="cancelRating" class="back-button" aria-label="返回课程详情">
        <span class="btn-icon" aria-hidden="true">←</span>
        返回课程详情
      </button>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <div class="loader"></div>
      <p class="loading-text">正在加载课程信息...</p>
    </div>

    <!-- 评价表单 -->
    <div v-else class="rate-content">
      <!-- 课程信息卡片 -->
      <div class="card course-info-card">
        <h2 class="card-title">课程信息</h2>
        <div class="course-info">
          <h3 class="course-name">{{ course.name }}</h3>
          <div class="course-meta">
            <div class="meta-item">
              <svg class="meta-icon" width="16" height="16" viewBox="0 0 24 24" fill="none"
                xmlns="http://www.w3.org/2000/svg">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" stroke="var(--color-text)" stroke-width="2" fill="none" />
                <circle cx="12" cy="7" r="4" stroke="var(--color-text)" stroke-width="2" fill="none" />
              </svg>
              <span class="meta-text">{{ course.teacher }}</span>
            </div>
            <div class="meta-item">
              <svg class="meta-icon" width="16" height="16" viewBox="0 0 24 24" fill="none"
                xmlns="http://www.w3.org/2000/svg">
                <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20" stroke="var(--color-text)" stroke-width="2" fill="none" />
                <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z" stroke="var(--color-text)" stroke-width="2" fill="none" />
              </svg>
              <span class="meta-text">{{ course.credits }} 学分</span>
            </div>
            <div class="meta-item">
              <svg class="meta-icon" width="16" height="16" viewBox="0 0 24 24" fill="none"
                xmlns="http://www.w3.org/2000/svg">
                <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z" stroke="var(--color-text)" stroke-width="2" fill="none" />
                <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z" stroke="var(--color-text)" stroke-width="2" fill="none" />
              </svg>
              <span class="meta-text">{{ course.grade }} · {{ course.semester }} · {{ course.subject }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 评价表单卡片 -->
      <div class="card rating-form-card">
        <h2 class="card-title">{{ isEditMode ? '编辑评价' : '课程评价' }}</h2>
        <p class="form-subtitle">您的真实评价将帮助其他同学做出更好的选择</p>

        <!-- 总体评分 -->
        <div class="rating-section">
          <div class="section-header">
            <h3 class="section-title">总体评分 *</h3>
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
              <svg v-if="i <= ratingForm.overall" width="24" height="24" viewBox="0 0 24 24" fill="#F7D074">
                <path d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z"/>
              </svg>
              <svg v-else width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#888888" stroke-width="2">
                <path d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z"/>
              </svg>
            </span>
          </div>
          <div class="rating-help">
            {{ ratingOptions.find(opt => opt.value === ratingForm.overall)?.label || '请选择评分' }}
          </div>
        </div>

        <!-- 详细评分 -->
        <div class="detailed-ratings">
          <h3 class="detailed-title">详细评分</h3>
          
          <div class="rating-item">
            <div class="rating-item-header">
              <h4 class="rating-item-title">课程难度</h4>
              <p class="rating-item-desc">课程的难易程度</p>
            </div>
            <div class="rating-stars">
              <span
                v-for="i in 5"
                :key="i"
                class="star"
                :class="{ active: i <= ratingForm.difficulty }"
                @click="handleStarClick('difficulty', i)"
              >
                <svg v-if="i <= ratingForm.difficulty" width="20" height="20" viewBox="0 0 24 24" fill="#F7D074">
                  <path d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z"/>
                </svg>
                <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#888888" stroke-width="2">
                  <path d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z"/>
                </svg>
              </span>
            </div>
          </div>

          <div class="rating-item">
            <div class="rating-item-header">
              <h4 class="rating-item-title">实用性</h4>
              <p class="rating-item-desc">课程内容对您是否有帮助</p>
            </div>
            <div class="rating-stars">
              <span
                v-for="i in 5"
                :key="i"
                class="star"
                :class="{ active: i <= ratingForm.usefulness }"
                @click="handleStarClick('usefulness', i)"
              >
                <svg v-if="i <= ratingForm.usefulness" width="20" height="20" viewBox="0 0 24 24" fill="#F7D074">
                  <path d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z"/>
                </svg>
                <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#888888" stroke-width="2">
                  <path d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z"/>
                </svg>
              </span>
            </div>
          </div>

          <div class="rating-item">
            <div class="rating-item-header">
              <h4 class="rating-item-title">教学质量</h4>
              <p class="rating-item-desc">教师的教学水平和方法</p>
            </div>
            <div class="rating-stars">
              <span
                v-for="i in 5"
                :key="i"
                class="star"
                :class="{ active: i <= ratingForm.teaching }"
                @click="handleStarClick('teaching', i)"
              >
                <svg v-if="i <= ratingForm.teaching" width="20" height="20" viewBox="0 0 24 24" fill="#F7D074">
                  <path d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z"/>
                </svg>
                <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#888888" stroke-width="2">
                  <path d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z"/>
                </svg>
              </span>
            </div>
          </div>
        </div>

        <!-- 文字评价 -->
        <div class="text-rating">
          <div class="section-header">
            <h3 class="section-title">详细评价 *</h3>
            <p class="section-desc">分享您的学习体验和建议，支持Markdown格式</p>
          </div>
          
          <!-- 编辑/预览切换 -->
          <div class="editor-tabs">
            <button 
              type="button"
              class="tab-button" 
              :class="{ active: !showPreview }"
              @click="showPreview = false"
            >
              编辑
            </button>
            <button 
              type="button"
              class="tab-button" 
              :class="{ active: showPreview }"
              @click="showPreview = true"
            >
              预览
            </button>
          </div>

          <!-- 编辑模式 -->
          <div v-if="!showPreview" class="editor-content">
            <textarea
              v-model="ratingForm.content"
              class="rating-textarea"
              placeholder="请详细描述您对这门课程的看法，包括课程内容、教学方法、考核方式等方面。您的评价将对其他同学有重要的参考价值...

支持Markdown格式：
- **粗体文本**
- *斜体文本*
- # 标题
- - 列表项
- > 引用"
              rows="8"
            ></textarea>
          </div>

          <!-- 预览模式 -->
          <div v-else class="preview-content">
            <div class="markdown-preview" v-html="parseMarkdown(ratingForm.content)"></div>
          </div>

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
            <span v-else>{{ isEditMode ? '更新评价' : '提交评价' }}</span>
          </button>
        </div>
      </div>

      <!-- 评价指南卡片 -->
      <div class="card rating-guide-card">
        <h2 class="card-title">评价指南</h2>
        <div class="guide-content">
          <div class="guide-item">
            <div class="guide-icon">✓</div>
            <div class="guide-text">
              <strong>真实客观：</strong>基于您的真实学习体验进行评价
            </div>
          </div>
          <div class="guide-item">
            <div class="guide-icon">✓</div>
            <div class="guide-text">
              <strong>详细具体：</strong>提供具体的例子和详细描述
            </div>
          </div>
          <div class="guide-item">
            <div class="guide-icon">✓</div>
            <div class="guide-text">
              <strong>尊重他人：</strong>使用文明用语，尊重教师和其他同学
            </div>
          </div>
          <div class="guide-item">
            <div class="guide-icon">✓</div>
            <div class="guide-text">
              <strong>建设性：</strong>提出有价值的建议和改进意见
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 主容器 */
.course-rate-container {
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

/* 返回按钮样式 */
.back-button {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background-color: var(--color-surface);
  color: var(--color-text);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 12px 20px;
  font-size: 16px;
  font-weight: 500;
  text-decoration: none;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-left: var(--card-margin);
  box-shadow: none;
}

.back-button:hover,
.back-button:active,
.back-button:focus {
  background-color: var(--color-primary);
  transform: none;
  box-shadow: none;
  outline: none;
  text-decoration: none;
  color: var(--color-text);
}

.back-button .btn-icon {
  font-size: 18px;
  font-weight: bold;
  transition: none;
}

.back-button:hover .btn-icon {
  transform: none;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .back-button {
    padding: 10px 16px;
    font-size: 14px;
  }
  
  .back-button .btn-icon {
    font-size: 16px;
  }
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

/* 评价内容 */
.rate-content {
  max-width: 800px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* 卡片通用样式 */
.card {
  background-color: #FFFFFF;
  border-radius: 12px;
  border: 1px solid var(--color-border);
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

/* 课程信息卡片 */
.course-info {
  text-align: center;
}

.course-name {
  font-size: 24px;
  font-weight: bold;
  margin: 0 0 16px 0;
}

.course-meta {
  display: flex;
  justify-content: center;
  gap: 16px;
  flex-wrap: wrap;
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

.meta-text {
  font-size: 16px;
}

/* 表单副标题 */
.form-subtitle {
  font-size: 16px;
  color: #888888;
  text-align: center;
  margin: 0 0 24px 0;
}

/* 评分部分 */
.rating-section {
  margin-bottom: 24px;
}

.section-header {
  text-align: center;
  margin-bottom: 16px;
}

.section-title {
  font-size: 18px;
  font-weight: bold;
  margin: 0 0 8px 0;
}

.section-desc {
  font-size: 14px;
  color: #888888;
  margin: 0;
}

.rating-stars {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin-bottom: 8px;
}

.star {
  font-size: 24px;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.star:hover {
}

.star.active {
  color: #F7D074;
}

.rating-help {
  text-align: center;
  font-size: 14px;
  color: #1A1A1A;
  font-weight: bold;
}

/* 详细评分 */
.detailed-ratings {
  margin-bottom: 24px;
}

.detailed-title {
  font-size: 18px;
  font-weight: bold;
  margin: 0 0 16px 0;
  text-align: center;
}

.rating-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background-color: #FEF6F7;
  border-radius: 8px;
  border: 1px solid var(--color-border);
  margin-bottom: 16px;
}

.rating-item-header {
  flex: 1;
}

.rating-item-title {
  font-size: 16px;
  font-weight: bold;
  margin: 0 0 4px 0;
}

.rating-item-desc {
  font-size: 14px;
  color: #888888;
  margin: 0;
}

/* 文字评价 */
.text-rating {
  margin-bottom: 24px;
}

/* 编辑器标签页 */
.editor-tabs {
  display: flex;
  border-bottom: 1px solid var(--color-border);
  margin-bottom: 16px;
}

.tab-button {
  padding: 8px 16px;
  border: none;
  background: transparent;
  color: var(--color-textSecondary);
  cursor: pointer;
  border-bottom: 2px solid transparent;
  transition: all 0.2s ease;
  font-size: 14px;
}

.tab-button.active {
  color: var(--color-primary);
  border-bottom-color: var(--color-primary);
}

.tab-button:hover {
  color: var(--color-text);
}

/* 编辑器内容 */
.editor-content {
  margin-bottom: 12px;
}

.rating-textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  font-size: 16px;
  font-family: sans-serif;
  resize: vertical;
  min-height: 120px;
  background-color: #FFFFFF;
  transition: all 0.2s ease;
}

.rating-textarea:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(var(--color-primary-rgb), 0.1);
}

/* 预览内容 */
.preview-content {
  min-height: 120px;
  padding: 12px;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  background-color: #FFFFFF;
  margin-bottom: 12px;
}

/* Markdown预览样式 */
.markdown-preview {
  line-height: 1.6;
  color: var(--color-text);
}

.markdown-preview h1,
.markdown-preview h2,
.markdown-preview h3 {
  margin: 16px 0 8px 0;
  color: var(--color-text);
}

.markdown-preview h1:first-child,
.markdown-preview h2:first-child,
.markdown-preview h3:first-child {
  margin-top: 0;
}

.markdown-preview p {
  margin: 8px 0;
}

.markdown-preview p:first-child {
  margin-top: 0;
}

.markdown-preview p:last-child {
  margin-bottom: 0;
}

.markdown-preview code {
  background: var(--color-border);
  padding: 2px 4px;
  border-radius: 3px;
  font-family: 'Courier New', monospace;
  font-size: 13px;
}

.markdown-preview pre {
  background: var(--color-border);
  padding: 12px;
  border-radius: 6px;
  overflow-x: auto;
  margin: 12px 0;
}

.markdown-preview blockquote {
  border-left: 4px solid var(--color-primary);
  padding-left: 12px;
  margin: 12px 0;
  color: var(--color-textSecondary);
}

.markdown-preview ul,
.markdown-preview ol {
  margin: 8px 0;
  padding-left: 20px;
}

.markdown-preview li {
  margin: 4px 0;
}

.char-count {
  text-align: right;
  font-size: 14px;
  color: #888888;
  margin-top: 8px;
}

/* 匿名选项 */
.anonymous-option {
  margin-bottom: 24px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  user-select: none;
}

.checkbox {
  display: none;
}

.checkmark {
  width: 24px;
  height: 24px;
  border: 1px solid var(--color-border);
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #FFFFFF;
  transition: transform 0.2s ease;
}

.checkbox:checked + .checkmark {
  background-color: #76D7C4;
}

.checkbox:checked + .checkmark::after {
  content: '✓';
  color: #1A1A1A;
  font-size: 16px;
  font-weight: bold;
}

.label-text {
  font-size: 16px;
  font-weight: bold;
  color: #1A1A1A;
}

.anonymous-help {
  font-size: 14px;
  color: #888888;
  margin-top: 8px;
  margin-left: 36px;
}

/* 按钮通用样式 */
.btn {
  background-color: #F7D074;
  border-radius: 8px;
  border: 1px solid var(--color-border);
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
  transform: none;
}

.btn:active {
  transform: none;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background-color: #76D7C4;
}

.btn-icon {
  font-size: 16px;
}

/* 表单操作按钮 */
.form-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
}

/* 评价指南 */
.guide-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.guide-item {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.guide-icon {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background-color: #76D7C4;
  color: #FFFFFF;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  border: 1px solid var(--color-border);
  flex-shrink: 0;
}

.guide-text {
  font-size: 16px;
  line-height: 1.5;
  color: #1A1A1A;
}

.guide-text strong {
  font-weight: bold;
}

/* 加载动画 */
.loader {
  width: 40px;
  height: 40px;
  border: 3px solid #FEF6F7;
  border-top: 3px solid #F7D074;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .course-rate-container {
    padding: 16px;
  }
  
  .rating-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .btn {
    width: 100%;
  }
}
</style>