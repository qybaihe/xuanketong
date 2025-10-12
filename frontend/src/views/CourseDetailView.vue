<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { evaluationRequestService, type EvaluationRequest } from '@/services/api'
import axios from 'axios'
import CourseMaterialApi from '@/services/courseMaterialApi'
import CourseNoteApi from '@/services/courseNoteApi'
import { marked } from 'marked'

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
const mapCourseData = (course: Record<string, unknown>): Course => {
  return {
    // 原始后端字段
    id: course.id as number,
    name: course.name as string,
    description: course.description as string,
    grade: course.grade as string,
    semester: course.semester as string,
    subject: course.subject as string,
    teacher: course.teacher as string,
    credits: course.credits as number,
    createdAt: course.createdAt as string,
    updatedAt: course.updatedAt as string,
    // 兼容字段映射
    ID: course.id as number,
    Name: course.name as string,
    Description: course.description as string,
    Grade: course.grade as string,
    Semester: course.semester as string,
    Subject: course.subject as string,
    Teacher: course.teacher as string,
    Credits: course.credits as number,
    CreatedAt: course.createdAt as string,
    UpdatedAt: course.updatedAt as string
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
  Nickname?: string;
}

interface Comment {
  ID: number;
  UserID: number;
  CourseID: number;
  Content: string;
  Username?: string;
  Nickname?: string;
  CreatedAt?: string;
}

const route = useRoute()
const authStore = useAuthStore()
const course = ref<Course | null>(null)
const ratings = ref<Rating[]>([])
const comments = ref<Comment[]>([])
const loading = ref(true)
const evaluationRequestLoading = ref(false)
const hasEvaluationRequested = ref(false)

// 笔记相关状态
const noteContent = ref('')
const noteLoading = ref(false)
const noteSaving = ref(false)
const showNoteEditor = ref(false)
const activeTab = ref('edit')
interface Note {
  id: number
  content: string
  authorName: string
  createdAt: string
  isOwner: boolean
}

interface User {
  username?: string
  nickname?: string
}

interface ApiResponse {
  user?: User
  UserID: number
  score?: number
  Score?: number
  createdAt?: string
}

const notes = ref<Note[]>([])
const editingNote = ref<Note | null>(null)

// 资料相关状态
interface CourseMaterial {
  id: number
  courseId: number
  materialName: string
  type: 'cloud' | 'note'
  service: string
  link: string
  accessPassword?: string
  description?: string
  uploaderId: number
  uploaderName: string
  uploadTime: string
  downloadCount: number
}

const materials = ref<CourseMaterial[]>([])
const materialsLoading = ref(false)
const uploadLoading = ref(false)
const showUploadModal = ref(false)
const materialName = ref('')
const materialType = ref<'cloud' | 'note'>('cloud')
const cloudService = ref('百度网盘')
const noteService = ref('幕布')
const materialLink = ref('')
const accessPassword = ref('')
const materialDescription = ref('')
const showCloudTypeDropdown = ref(false)
const showNoteTypeDropdown = ref(false)

// 网盘类型选项
const cloudServices = [
  '百度网盘',
  '夸克网盘',
  '阿里云盘',
  '坚果云',
  'Dropbox',
  'OneDrive',
  '123Pan',
  '蓝奏云',
  '天翼云盘',
  '其他'
]

const noteServices = [
  '幕布',
  'OneNote',
  '坚果云',
  '知悉',
  '飞书文档',
  '语雀',
  'Notion',
  '印象笔记',
  '其他'
]

const activeMaterialTab = ref<'all' | 'cloud' | 'note'>('all')

const filteredMaterials = computed(() => {
  if (activeMaterialTab.value === 'all') return materials.value
  return materials.value.filter((item) => item.type === activeMaterialTab.value)
})


// 检查用户是否登录
const canSubmit = computed(() => authStore.isAuthenticated)

// 获取当前用户信息
const currentUser = computed(() => authStore.user)


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
  } catch (error: unknown) {
    console.error('求评价请求失败:', error)
    const errorMessage = error && typeof error === 'object' && 'response' in error
      ? (error as { response?: { data?: { error?: string } } }).response?.data?.error
      : '求评价请求失败，请重试！'
    alert(errorMessage || '求评价请求失败，请重试！')
  } finally {
    evaluationRequestLoading.value = false
  }
}

// 笔记相关方法
const fetchNotes = async (courseId: number) => {
  try {
    noteLoading.value = true
    const response = await CourseNoteApi.getCourseNotes(courseId)
    if (response.success) {
      notes.value = response.data.notes || []
    }
  } catch (error) {
    console.error('获取笔记失败:', error)
    notes.value = []
  } finally {
    noteLoading.value = false
  }
}

const saveNote = async () => {
  if (!course.value?.ID) return

  try {
    noteSaving.value = true
    const response = await CourseNoteApi.saveCourseNote(course.value.ID, noteContent.value, false)
    if (response.success) {
      await fetchNotes(course.value.ID)
      closeNoteEditor()
      console.log('笔记保存成功')
    }
  } catch (error) {
    console.error('保存笔记失败:', error)
  } finally {
    noteSaving.value = false
  }
}

const startNewNote = () => {
  editingNote.value = null
  noteContent.value = ''
  showNoteEditor.value = true
  activeTab.value = 'edit'
}

const editNote = (note: Note) => {
  editingNote.value = note
  noteContent.value = note.content
  showNoteEditor.value = true
  activeTab.value = 'edit'
}

const closeNoteEditor = () => {
  showNoteEditor.value = false
  editingNote.value = null
  noteContent.value = ''
  activeTab.value = 'edit'
}

// 网盘类型下拉框方法
const toggleServiceDropdown = (type: 'cloud' | 'note') => {
  if (type === 'cloud') {
    showCloudTypeDropdown.value = !showCloudTypeDropdown.value
    if (showCloudTypeDropdown.value) showNoteTypeDropdown.value = false
  } else {
    showNoteTypeDropdown.value = !showNoteTypeDropdown.value
    if (showNoteTypeDropdown.value) showCloudTypeDropdown.value = false
  }
}

const selectService = (type: 'cloud' | 'note', service: string) => {
  if (type === 'cloud') {
    cloudService.value = service
    showCloudTypeDropdown.value = false
  } else {
    noteService.value = service
    showNoteTypeDropdown.value = false
  }
}

// 点击外部关闭下拉框
const handleClickOutside = (event: Event) => {
  const target = event.target as HTMLElement
  if (!target.closest('.cloud-type-dropdown')) {
    showCloudTypeDropdown.value = false
  }
  if (!target.closest('.note-type-dropdown')) {
    showNoteTypeDropdown.value = false
  }
}

const openMaterialModal = () => {
  showUploadModal.value = true
}

// 资料相关方法
const fetchMaterials = async (courseId: number) => {
  try {
    materialsLoading.value = true
    const response = await CourseMaterialApi.getCourseMaterials(courseId)
    if (response.success) {
      materials.value = response.data.materials || []
    }
  } catch (error) {
    console.error('获取资料失败:', error)
    materials.value = []
  } finally {
    materialsLoading.value = false
  }
}

const uploadMaterialFile = async () => {
  if (!materialName.value || !materialLink.value || !course.value?.ID) {
    alert('请填写资料名称和资源链接！')
    return
  }

  try {
    uploadLoading.value = true
    const materialData = {
      materialName: materialName.value,
      type: materialType.value,
      service: materialType.value === 'cloud' ? cloudService.value : noteService.value,
      link: materialLink.value,
      accessPassword: accessPassword.value,
      description: materialDescription.value
    }

    const response = await CourseMaterialApi.uploadMaterial(course.value.ID, materialData)
    if (response.success) {
      // 重新获取资料列表
      await fetchMaterials(course.value.ID)
      // 关闭上传模态框并清空表单
      showUploadModal.value = false
      materialName.value = ''
      materialType.value = 'cloud'
      cloudService.value = '百度网盘'
      noteService.value = '幕布'
      materialLink.value = ''
      accessPassword.value = ''
      materialDescription.value = ''
      // 可以添加成功提示
      console.log('资料上传成功')
    }
  } catch (error) {
    console.error('上传资料失败:', error)
  } finally {
    uploadLoading.value = false
  }
}

// 简单的markdown解析器
marked.setOptions({
  breaks: true,
  gfm: true
})

const parseMarkdown = (text: string) => {
  const rawHtml = marked.parse(text || '')
  return rawHtml
}


// 格式化时间
const formatTime = (timeString: string) => {
  return new Date(timeString).toLocaleString('zh-CN')
}

onMounted(async () => {
  try {
    loading.value = true
    const courseId = parseInt(route.params.id as string)
    const [courseResponse, ratingsResponse, commentsResponse] = await Promise.all([
      axios.get(`${import.meta.env.VITE_BACKEND_BASE_URL}/courses/${courseId}`),
      axios.get(`${import.meta.env.VITE_BACKEND_BASE_URL}/courses/${courseId}/ratings`),
      axios.get(`${import.meta.env.VITE_BACKEND_BASE_URL}/courses/${courseId}/comments`)
    ])
    course.value = mapCourseData(courseResponse.data.data)
    // 使用后端返回的username，如果没有则使用nickname，如果都没有才使用默认
    ratings.value = ratingsResponse.data.data.map((rating: ApiResponse) => ({
      ...rating,
      Username: rating.user?.username || `用户${rating.UserID}`,
      Nickname: rating.user?.nickname || `用户${rating.UserID}`,
      Score: !isNaN(rating.score as number) && isFinite(rating.score as number) ? rating.score : (!isNaN(rating.Score as number) && isFinite(rating.Score as number) ? rating.Score : 0)
    }))
    // 使用后端返回的username，如果没有则使用nickname，如果都没有才使用默认
    comments.value = (commentsResponse.data.data || []).map((comment: ApiResponse) => ({
      ...comment,
      Username: comment.user?.username || `用户${comment.UserID}`,
      Nickname: comment.user?.nickname || `用户${comment.UserID}`,
      CreatedAt: comment.createdAt ? new Date(comment.createdAt).toLocaleDateString() : new Date().toLocaleDateString()
    }))

    // 加载笔记和资料
    await Promise.all([
      fetchNotes(courseId),
      fetchMaterials(courseId)
    ])

    // 检查是否已经发起过求评价请求
    if (authStore.isAuthenticated) {
      try {
        const evaluationResponse = await evaluationRequestService.getEvaluationRequests()
        hasEvaluationRequested.value = evaluationResponse.items.some((req: EvaluationRequest) =>
          req.courseId === courseId && req.status === 'pending'
        )
      } catch (error) {
        console.error('检查求评价请求失败:', error)
      }
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
})

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div class="course-detail-container">
    <!-- 返回按钮 -->
    <div class="back-button-container">
      <RouterLink to="/" class="back-button" aria-label="返回课程列表">
        <span class="btn-icon" aria-hidden="true">←</span>
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
      <!-- Main Content Grid -->
      <div class="content-grid">
        <!-- Left Column: Course Info and Evaluations -->
        <div class="left-column">
          <!-- Course Header Card -->
          <div class="card course-header">
            <h1 class="course-title">{{ course.Name }}</h1>

            <div class="course-meta">
              <div class="meta-item">
                <svg class="meta-icon" width="16" height="16" viewBox="0 0 24 24" fill="none"
                  xmlns="http://www.w3.org/2000/svg">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" stroke="var(--color-text)" stroke-width="2" fill="none" />
                  <circle cx="12" cy="7" r="4" stroke="var(--color-text)" stroke-width="2" fill="none" />
                </svg>
                <span class="meta-text">{{ course.Teacher }}</span>
              </div>
              <div class="meta-item">
                <svg class="meta-icon" width="16" height="16" viewBox="0 0 24 24" fill="none"
                  xmlns="http://www.w3.org/2000/svg">
                  <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20" stroke="var(--color-text)" stroke-width="2" fill="none" />
                  <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z" stroke="var(--color-text)" stroke-width="2"
                    fill="none" />
                </svg>
                <span class="meta-text">{{ course.Grade }} · {{ course.Semester }}</span>
              </div>
              <div class="meta-item">
                <svg class="meta-icon" width="16" height="16" viewBox="0 0 24 24" fill="none"
                  xmlns="http://www.w3.org/2000/svg">
                  <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z" stroke="var(--color-text)" stroke-width="2" fill="none" />
                  <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z" stroke="var(--color-text)" stroke-width="2" fill="none" />
                </svg>
                <span class="meta-text">{{ course.Subject }}</span>
              </div>
              <div class="meta-item">
                <svg class="meta-icon" width="16" height="16" viewBox="0 0 24 24" fill="none"
                  xmlns="http://www.w3.org/2000/svg">
                  <path d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z" stroke="var(--color-text)" stroke-width="2" fill="none" />
                </svg>
                <span class="meta-text">{{ course.Credits }} 学分</span>
              </div>
            </div>，尝试

            <div class="course-rating-summary">
              <div class="rating-stars">
                <span v-for="i in 5" :key="i" class="star">
                  <svg v-if="i <= Math.floor(averageRating)" width="16" height="16" viewBox="0 0 24 24" fill="var(--color-primary)">
                    <path d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z" />
                  </svg>
                  <svg v-else-if="i - 0.5 <= averageRating" width="16" height="16" viewBox="0 0 24 24" fill="var(--color-primary)">
                    <path d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z" />
                  </svg>
                  <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="var(--color-textSecondary)" stroke-width="1.5">
                    <path d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z" />
                  </svg>
                </span>
              </div>
              <span class="rating-value">{{ averageRating }}</span>
              <span class="rating-count">({{ ratings.length }} 人评分)</span>
            </div>

            <!-- Action Buttons -->
            <div class="action-buttons">
              <RouterLink :to="`/courses/${course.ID}/rate`" class="btn btn-primary">
                <svg class="btn-icon" width="16" height="16" viewBox="0 0 24 24" fill="var(--color-text)">
                  <path d="M12 2L15.09 8.26L22 9L17 14.74L18.18 22L12 18.77L5.82 22L7 14.74L2 9L8.91 8.26L12 2Z" />
                </svg>
                去评价课程
              </RouterLink>

              <button v-if="canSubmit" @click="submitEvaluationRequest" class="btn btn-secondary"
                :disabled="evaluationRequestLoading || hasEvaluationRequested">
                <svg class="btn-icon" width="16" height="16" viewBox="0 0 24 24" fill="none"
                  xmlns="http://www.w3.org/2000/svg">
                  <path d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z"
                    fill="var(--color-text)" />
                </svg>
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
                        <svg v-if="i <= Math.floor(averageRating)" width="16" height="16" viewBox="0 0 24 24"
                          fill="var(--color-primary)">
                          <path
                            d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z" />
                        </svg>
                        <svg v-else-if="i - 0.5 <= averageRating" width="16" height="16" viewBox="0 0 24 24"
                          fill="var(--color-primary)">
                          <path
                            d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z" />
                        </svg>
                        <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                          stroke-width="2">
                          <path
                            d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z" />
                        </svg>
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
                        <svg v-if="i <= Math.floor(averageDifficulty)" width="16" height="16" viewBox="0 0 24 24"
                          fill="var(--color-primary)">
                          <path
                            d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z" />
                        </svg>
                        <svg v-else-if="i - 0.5 <= averageDifficulty" width="16" height="16" viewBox="0 0 24 24"
                          fill="var(--color-primary)">
                          <path
                            d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z" />
                        </svg>
                        <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                          stroke-width="2">
                          <path
                            d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z" />
                        </svg>
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
                        <svg v-if="i <= Math.floor(averageUsefulness)" width="16" height="16" viewBox="0 0 24 24"
                          fill="var(--color-primary)">
                          <path
                            d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z" />
                        </svg>
                        <svg v-else-if="i - 0.5 <= averageUsefulness" width="16" height="16" viewBox="0 0 24 24"
                          fill="var(--color-primary)">
                          <path
                            d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z" />
                        </svg>
                        <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                          stroke-width="2">
                          <path
                            d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z" />
                        </svg>
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
                        <svg v-if="i <= Math.floor(averageTeaching)" width="16" height="16" viewBox="0 0 24 24"
                          fill="var(--color-primary)">
                          <path
                            d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z" />
                        </svg>
                        <svg v-else-if="i - 0.5 <= averageTeaching" width="16" height="16" viewBox="0 0 24 24"
                          fill="var(--color-primary)">
                          <path
                            d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z" />
                        </svg>
                        <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                          stroke-width="2">
                          <path
                            d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z" />
                        </svg>
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
                    <span class="user-avatar">{{ (rating.Nickname || rating.Username)?.charAt(0) || 'U' }}</span>
                    <span class="user-name">{{ rating.Nickname || rating.Username }}</span>
                  </div>
                  <div class="rating-score">
                    <div class="rating-stars">
                      <span v-for="i in 5" :key="i" class="star">
                        <svg v-if="i <= Math.floor(rating.Score)" width="16" height="16" viewBox="0 0 24 24"
                          fill="var(--color-primary)">
                          <path
                            d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z" />
                        </svg>
                        <svg v-else-if="i - 0.5 <= rating.Score" width="16" height="16" viewBox="0 0 24 24"
                          fill="var(--color-primary)">
                          <path
                            d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z" />
                        </svg>
                        <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                          stroke-width="2">
                          <path
                            d="M12 2l3.09 6.26L22 9l-5 4.87 1.18 6.88L12 17.77l-6.18 2.98L7 13.87 2 9l6.91-1.74L12 2z" />
                        </svg>
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
                    <span class="user-avatar">{{ (comment.Nickname || comment.Username)?.charAt(0) || 'U' }}</span>
                    <div class="user-info">
                      <span class="user-name">{{ comment.Nickname || comment.Username }}</span>
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

        <!-- Right Column: Materials and Notes -->
        <div class="right-column">
          <!-- Materials Section -->
          <div class="materials-section">
            <div class="card materials-card">
              <div class="materials-header">
                <h2 class="card-title">课程资料</h2>
                <button @click="openMaterialModal" class="btn btn-primary">
                  <svg class="btn-icon" width="16" height="16" viewBox="0 0 24 24" fill="none"
                    xmlns="http://www.w3.org/2000/svg">
                    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" stroke="currentColor"
                      stroke-width="2" fill="none" />
                    <polyline points="14,2 14,8 20,8" stroke="currentColor" stroke-width="2" fill="none" />
                    <line x1="16" y1="13" x2="8" y2="13" stroke="currentColor" stroke-width="2" />
                    <line x1="16" y1="17" x2="8" y2="17" stroke="currentColor" stroke-width="2" />
                    <polyline points="10,9 9,9 8,9" stroke="currentColor" stroke-width="2" fill="none" />
                  </svg>
                  分享资料
                </button>
              </div>

              <div v-if="materialsLoading" class="loading-container">
                <p>加载资料中...</p>
              </div>

              <div v-else-if="materials.length === 0" class="empty-materials">
                <svg class="empty-icon" width="48" height="48" viewBox="0 0 24 24" fill="none"
                  xmlns="http://www.w3.org/2000/svg">
                  <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" stroke="currentColor" stroke-width="2"
                    fill="none" />
                  <polyline points="14,2 14,8 20,8" stroke="currentColor" stroke-width="2" fill="none" />
                  <line x1="16" y1="13" x2="8" y2="13" stroke="currentColor" stroke-width="2" />
                  <line x1="16" y1="17" x2="8" y2="17" stroke="currentColor" stroke-width="2" />
                  <polyline points="10,9 9,9 8,9" stroke="currentColor" stroke-width="2" fill="none" />
                </svg>
                <p>还没有资料，点击"分享资料"上传第一个吧！</p>
              </div>

              <div v-else class="materials-list">
                <div v-if="filteredMaterials.length === 0" class="empty-materials">
                  <svg class="empty-icon" width="48" height="48" viewBox="0 0 24 24" fill="none"
                    xmlns="http://www.w3.org/2000/svg">
                    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" stroke="currentColor"
                      stroke-width="2" fill="none" />
                    <polyline points="14,2 14,8 20,8" stroke="currentColor" stroke-width="2" fill="none" />
                    <line x1="16" y1="13" x2="8" y2="13" stroke="currentColor" stroke-width="2" />
                    <line x1="16" y1="17" x2="8" y2="17" stroke="currentColor" stroke-width="2" />
                    <polyline points="10,9 9,9 8,9" stroke="currentColor" stroke-width="2" fill="none" />
                  </svg>
                  <p>暂时还没有此类型的资料，点击"分享资料"发布吧！</p>
                </div>
                <div v-for="material in filteredMaterials" :key="material.id" class="material-item">
                  <div class="material-info">
                    <div class="material-header">
                      <div class="material-title">
                        <h3 class="material-name">{{ material.materialName }}</h3>
                      </div>
                      <span class="cloud-type">{{ material.service }}（{{ material.type === 'cloud' ? '网盘' : '笔记' }}）</span>
                    </div>
                    <p class="material-description">{{ material.description }}</p>
                    <div class="material-meta">
                      <span class="uploader">上传者：{{ material.uploaderName }}</span>
                      <span class="upload-time">{{ formatTime(material.uploadTime) }}</span>
                      <span class="download-count">浏览：{{ material.downloadCount }}次</span>
                    </div>
                  </div>
                  <div class="material-actions">
                    <a :href="material.link" target="_blank" class="btn btn-primary">
                      <svg class="btn-icon" width="16" height="16" viewBox="0 0 24 24" fill="none"
                        xmlns="http://www.w3.org/2000/svg">
                        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" stroke="currentColor" stroke-width="2"
                          fill="none" />
                        <polyline points="7,10 12,15 17,10" stroke="currentColor" stroke-width="2" fill="none" />
                        <line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" />
                      </svg>
                      访问链接
                    </a>
                    <div v-if="material.accessPassword" class="password-info">
                      <span class="password-label">密码：</span>
                      <span class="password-value">{{ material.accessPassword }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Notes Section -->
          <div class="notes-section">
            <div class="card notes-card" :class="{ editing: showNoteEditor }">
              <div class="notes-header">
                <h2 class="card-title">课程笔记</h2>
                <div class="notes-actions">
                  <button v-if="!showNoteEditor" @click="startNewNote" class="btn btn-primary">
                    <svg class="btn-icon" width="16" height="16" viewBox="0 0 24 24" fill="none"
                      xmlns="http://www.w3.org/2000/svg">
                      <path
                        d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.39-.39-1.02-.39-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z"
                        fill="var(--color-text)" />
                    </svg>
                    写笔记
                  </button>
                  <button v-if="showNoteEditor" @click="closeNoteEditor" class="btn btn-secondary">
                    取消编辑
                  </button>
                </div>
              </div>

              <div v-if="noteLoading" class="loading-container">
                <p>加载笔记中...</p>
              </div>

              <div v-else-if="showNoteEditor" class="note-editor">
                <div class="editor-container">
                  <div class="editor-tabs">
                    <button class="tab-button active" @click="activeTab = 'edit'">
                      编辑
                    </button>
                    <button class="tab-button" @click="activeTab = 'preview'">
                      预览
                    </button>
                  </div>

                  <div v-if="activeTab === 'edit'" class="editor-content">
                    <textarea v-model="noteContent" class="note-textarea"
                      placeholder="在这里输入你的笔记内容，支持Markdown格式..."></textarea>
                  </div>

                  <div v-else class="preview-content">
                    <div class="markdown-preview" v-html="parseMarkdown(noteContent)"></div>
                  </div>
                </div>

                <div class="editor-actions">
                  <button @click="saveNote" class="btn btn-primary" :disabled="noteSaving">
                    <svg v-if="noteSaving" class="btn-icon" width="16" height="16" viewBox="0 0 24 24" fill="none"
                      xmlns="http://www.w3.org/2000/svg">
                      <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2" fill="none" />
                      <path d="M12 6v6l4 2" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
                    </svg>
                    {{ noteSaving ? '保存中...' : '保存笔记' }}
                  </button>
                </div>
              </div>

              <div v-else class="notes-list">
                <div v-if="notes.length === 0" class="empty-note">
                  <svg class="empty-icon" width="48" height="48" viewBox="0 0 24 24" fill="none"
                    xmlns="http://www.w3.org/2000/svg">
                    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" stroke="currentColor"
                      stroke-width="2" fill="none" />
                    <polyline points="14,2 14,8 20,8" stroke="currentColor" stroke-width="2" fill="none" />
                    <line x1="16" y1="13" x2="8" y2="13" stroke="currentColor" stroke-width="2" />
                    <line x1="16" y1="17" x2="8" y2="17" stroke="currentColor" stroke-width="2" />
                    <polyline points="10,9 9,9 8,9" stroke="currentColor" stroke-width="2" fill="none" />
                  </svg>
                  <p>还没有笔记，点击"写笔记"开始记录吧！</p>
                </div>
                <div v-else>
                  <div v-for="note in notes" :key="note.id" class="note-item">
                    <div class="note-header">
                      <div class="note-author">
                        <div class="author-avatar">{{ note.authorName?.charAt(0) || 'U' }}</div>
                        <div class="author-info">
                          <div class="author-name">{{ note.authorName }}</div>
                          <div class="note-date">{{ formatTime(note.createdAt) }}</div>
                        </div>
                      </div>
                      <div v-if="note.isOwner" class="note-actions">
                        <button @click="editNote(note)" class="btn-icon-small">
                          <svg width="14" height="14" viewBox="0 0 24 24" fill="none"
                            xmlns="http://www.w3.org/2000/svg">
                            <path
                              d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.39-.39-1.02-.39-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z"
                              fill="var(--color-text)" />
                          </svg>
                        </button>
                      </div>
                    </div>
                    <div class="note-content">
                      <div class="markdown-preview" v-html="parseMarkdown(note.content)"></div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Upload Modal -->
      <div v-if="showUploadModal" class="modal-overlay" @click="showUploadModal = false">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <h3>分享资料</h3>
            <button @click="showUploadModal = false" class="modal-close">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <line x1="18" y1="6" x2="6" y2="18" stroke="currentColor" stroke-width="2" />
                <line x1="6" y1="6" x2="18" y2="18" stroke="currentColor" stroke-width="2" />
              </svg>
            </button>
          </div>

          <form @submit.prevent="uploadMaterialFile">
            <div class="modal-body">
              <div class="form-group">
                <label for="materialName">资料名称 *</label>
                <input id="materialName" v-model="materialName" type="text" placeholder="请输入资料名称" class="form-input">
              </div>

            <div class="form-group">
              <label>资料类型 *</label>
              <div class="material-type-selector">
                <button class="type-toggle" :class="{ active: materialType === 'cloud' }"
                  @click="materialType = 'cloud'">
                  <svg class="type-icon" width="16" height="16" viewBox="0 0 24 24" fill="none"
                    xmlns="http://www.w3.org/2000/svg">
                    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" stroke="currentColor"
                      stroke-width="2" fill="none" />
                    <polyline points="14,2 14,8 20,8" stroke="currentColor" stroke-width="2" fill="none" />
                    <line x1="16" y1="13" x2="8" y2="13" stroke="currentColor" stroke-width="2" />
                    <line x1="16" y1="17" x2="8" y2="17" stroke="currentColor" stroke-width="2" />
                  </svg>
                  网盘资源
                </button>
                <button class="type-toggle" :class="{ active: materialType === 'note' }" @click="materialType = 'note'">
                  <svg class="type-icon" width="16" height="16" viewBox="0 0 24 24" fill="none"
                    xmlns="http://www.w3.org/2000/svg">
                    <path
                      d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.39-.39-1.02-.39-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z"
                      fill="currentColor" />
                  </svg>
                  课程笔记
                </button>
              </div>
            </div>

            <div class="form-group">
              <label>{{ materialType === 'cloud' ? '网盘类型 *' : '笔记平台 *' }}</label>
              <div
                :class="[materialType === 'cloud' ? 'cloud-type-dropdown' : 'note-type-dropdown', 'dropdown-container']">
                <button class="cloud-type-button" @click="toggleServiceDropdown(materialType)">
                  <span>{{ materialType === 'cloud' ? cloudService : noteService }}</span>
                  <svg class="dropdown-arrow"
                    :class="{ 'rotated': materialType === 'cloud' ? showCloudTypeDropdown : showNoteTypeDropdown }"
                    width="12" height="12" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M6 9l6 6 6-6" stroke="currentColor" stroke-width="2" stroke-linecap="round"
                      stroke-linejoin="round" />
                  </svg>
                </button>
                <div v-if="materialType === 'cloud' ? showCloudTypeDropdown : showNoteTypeDropdown"
                  class="cloud-type-options">
                  <button v-for="service in (materialType === 'cloud' ? cloudServices : noteServices)" :key="service"
                    @click="selectService(materialType, service)" class="cloud-type-option"
                    :class="{ active: (materialType === 'cloud' ? cloudService : noteService) === service }">
                    {{ service }}
                  </button>
                </div>
              </div>
            </div>

            <div class="form-group">
              <label for="materialLink">{{ materialType === 'cloud' ? '网盘链接 *' : '笔记链接 *' }}</label>
              <input id="materialLink" v-model="materialLink" type="url"
                :placeholder="materialType === 'cloud' ? '请输入网盘分享链接' : '请输入笔记访问链接'" class="form-input">
            </div>

            <div class="form-group">
              <label for="accessPassword">访问密码</label>
              <input id="accessPassword" v-model="accessPassword" type="text" placeholder="如果有访问密码请填写（可选）"
                class="form-input">
            </div>

              <div class="form-group">
                <label for="materialDescription">资料描述</label>
                <textarea id="materialDescription" v-model="materialDescription" placeholder="请简单描述一下这个资料的内容..."
                  class="form-textarea" rows="3"></textarea>
              </div>
            </div>

            <div class="modal-footer">
              <button type="button" @click="showUploadModal = false" class="btn btn-secondary">
                取消
              </button>
              <button type="submit" class="btn btn-primary" :disabled="uploadLoading">
                {{ uploadLoading ? '上传中...' : '确认分享' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 主容器 */
.course-detail-container {
  min-height: 100vh;
  background-color: var(--color-background);
  font-family: sans-serif;
  color: var(--color-text);
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
  font-weight: bold;
  text-decoration: none;
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
  color: var(--color-text-secondary);
  font-weight: bold;
}

/* 课程内容 */
.course-content {
  max-width: 1400px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.content-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 32px;
  align-items: start;
}

.left-column {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.right-column {
  display: flex;
  flex-direction: column;
  gap: 24px;
  position: sticky;
  top: 20px;
  max-height: calc(100vh - 40px);
  overflow-y: auto;
}

.right-column::-webkit-scrollbar {
  width: 8px;
}

.right-column::-webkit-scrollbar-track {
  background: var(--color-background);
  border-radius: 4px;
}

.right-column::-webkit-scrollbar-thumb {
  background: var(--color-text-secondary);
  border-radius: 4px;
}

.right-column::-webkit-scrollbar-thumb:hover {
  background: var(--color-text);
}

.notes-section {
  flex: 1;
}

.materials-section {
  flex: 1;
}

/* 卡片通用样式 */
.card {
  background-color: var(--color-surface);
  border-radius: 12px;
  border: 1px solid var(--color-border);
  padding: 20px;
  margin-bottom: 16px;
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
  color: var(--color-text);
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
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: bold;
  border: 1px solid var(--color-border);
}

.course-tag.primary {
  background-color: var(--color-accent);
  color: var(--color-text);
}

.course-tag.secondary {
  background-color: var(--color-primary);
  color: var(--color-text);
}

.course-tag.accent {
  background-color: var(--color-primary);
  color: var(--color-text);
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
  color: var(--color-text);
}

.rating-count {
  font-size: 14px;
  color: var(--color-text-secondary);
}

/* 按钮通用样式 */
.btn {
  background-color: var(--color-primary);
  border-radius: 8px;
  border: 1px solid var(--color-border);
  padding: 12px 20px;
  font-size: 14px;
  font-weight: bold;
  color: var(--color-text);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  text-decoration: none;
  transition: none;
}

.btn:hover,
.btn:active,
.btn:focus {
  box-shadow: none;
  outline: none;
  transform: none;
}


.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.btn-secondary {
  background-color: var(--color-accent);
}

/* 操作按钮区域 */
.action-buttons {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}

.evaluation-request-success {
  font-size: 14px;
  color: var(--color-accent);
  font-weight: bold;
  text-align: center;
  margin-top: 8px;
}

/* 课程描述 */
.description-text {
  font-size: 16px;
  line-height: 1.6;
  color: var(--color-text);
  margin: 0;
}

/* 评分和评论部分 */
.empty-state {
  text-align: center;
  padding: 24px;
  color: var(--color-text-secondary);
  font-size: 16px;
}

/* 平均评分展示 */
.average-ratings {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 24px;
  padding: 16px;
  background-color: var(--color-background);
  border-radius: 8px;
  border: 1px solid var(--color-border);
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
  color: var(--color-text);
}

.rating-value-stars {
  display: flex;
  align-items: center;
  gap: 8px;
}

.rating-number {
  font-size: 16px;
  font-weight: bold;
  color: var(--color-text);
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
  padding: 12px;
  background-color: var(--color-surface);
  border-radius: 8px;
  border: 1px solid var(--color-border);
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
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background-color: var(--color-primary);
  color: var(--color-text);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 14px;
  border: 1px solid var(--color-border);
}

.user-name {
  font-size: 14px;
  font-weight: bold;
  color: var(--color-text);
}

.rating-score {
  display: flex;
  align-items: center;
  gap: 8px;
}

.score-value {
  font-size: 14px;
  font-weight: bold;
  color: var(--color-text);
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
  color: var(--color-text-secondary);
}

.comment-content {
  font-size: 14px;
  line-height: 1.5;
  color: var(--color-text);
  width: 100%;
}

.comment-content p {
  margin: 0;
}

/* 笔记部分 */
.notes-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.notes-actions {
  display: flex;
  gap: 12px;
}

.notes-card.editing {
  border-color: var(--color-error);
  transition: border-color 0.2s ease;
}

.notes-card.editing .card-title {
  color: var(--color-error);
}

.notes-card.editing .tab-button.active {
  background-color: var(--color-error-light);
}

.note-editor {
  margin-top: 20px;
}

.editor-container {
  border: 1px solid var(--color-border);
  border-radius: 8px;
  background-color: var(--color-surface);
  overflow: hidden;
}

.editor-tabs {
  display: flex;
  background-color: var(--color-primary);
  border-bottom: 2px solid var(--color-border);
}

.tab-button {
  flex: 1;
  padding: 12px 16px;
  background: none;
  border: none;
  font-weight: bold;
  color: var(--color-text);
  cursor: pointer;
  transition: background-color 0.2s;
}

.tab-button:hover,
.tab-button:active,
.tab-button:focus {
  background-color: var(--color-primary);
  box-shadow: none;
  outline: none;
  transform: none;
}

.tab-button.active {
  background-color: var(--color-surface);
  border-bottom: 2px solid var(--color-border);
}

.editor-content,
.preview-content {
  min-height: 300px;
}

.note-textarea {
  width: 100%;
  height: 300px;
  padding: 16px;
  border: none;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.5;
  resize: vertical;
  outline: none !important;
}

.preview-content {
  padding: 16px;
}

.markdown-preview {
  line-height: 1.6;
  font-family: 'Inter', sans-serif;
  color: var(--color-text);
}

.markdown-preview h1,
.markdown-preview h2,
.markdown-preview h3,
.markdown-preview h4,
.markdown-preview h5 {
  font-weight: 700;
  margin: 20px 0 12px 0;
  color: var(--color-text);
}

.markdown-preview h1 {
  font-size: 22px;
  padding-bottom: 8px;
  border-bottom: 3px solid var(--color-border);
}

.markdown-preview h2 {
  font-size: 18px;
  position: relative;
  padding-left: 14px;
}

.markdown-preview h2::before {
  content: '';
  position: absolute;
  left: 0;
  top: 4px;
  bottom: 4px;
  width: 6px;
  background-color: var(--color-primary);
  border: 1px solid var(--color-border);
  border-radius: 3px;
}

.markdown-preview h3 {
  font-size: 16px;
  padding: 4px 8px;
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.markdown-preview strong {
  font-weight: 700;
  color: var(--color-text);
}

.markdown-preview em {
  font-style: italic;
}

.markdown-preview p {
  margin: 10px 0;
  line-height: 1.7;
}

.markdown-preview ul,
.markdown-preview ol {
  margin: 10px 0;
  padding-left: 24px;
}

.markdown-preview li {
  margin: 6px 0;
  padding-left: 6px;
}

.markdown-preview li::marker {
  font-weight: bold;
  color: var(--color-error);
}

.markdown-preview blockquote {
  border-left: 5px solid var(--color-error);
  padding: 12px 16px;
  margin: 16px 0;
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
}

.markdown-preview code {
  background-color: var(--color-background);
  padding: 2px 6px;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  font-family: 'Fira Code', 'Courier New', monospace;
  font-size: 13px;
}

.markdown-preview pre {
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 16px;
  overflow-x: auto;
  margin: 16px 0;
}

.markdown-preview pre code {
  background: none;
  border: none;
  padding: 0;
  font-size: 13px;
}

.markdown-preview a {
  color: var(--color-text);
  font-weight: bold;
  text-decoration: underline;
}

.markdown-preview a:hover {
  color: var(--color-error);
}

.markdown-preview table {
  width: 100%;
  border-collapse: collapse;
  margin: 16px 0;
  border: 1px solid var(--color-border);
}

.markdown-preview th,
.markdown-preview td {
  border: 1px solid var(--color-border);
  padding: 8px 12px;
  text-align: left;
  background-color: var(--color-surface);
}

.markdown-preview th {
  background-color: var(--color-primary);
}

.markdown-preview hr {
  margin: 16px 0;
  border: none;
  border-top: 3px solid var(--color-border);
}

.markdown-preview img {
  max-width: 100%;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  margin: 12px 0;
}

.editor-actions {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.note-display {
  margin-top: 20px;
}

.empty-note {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
}

.empty-note p {
  font-size: 16px;
  color: var(--color-text-secondary);
  margin-top: 16px;
}

.empty-icon {
  opacity: 0.5;
}

/* 笔记列表样式 */
.notes-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.note-item {
  padding: 16px;
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  transition: none;
}

.note-item:hover,
.note-item:active,
.note-item:focus {
  box-shadow: none;
  outline: none;
  transform: none;
}

.note-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.note-author {
  display: flex;
  align-items: center;
  gap: 8px;
}

.author-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background-color: var(--color-accent);
  color: var(--color-text);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 14px;
  border: 1px solid var(--color-border);
}

.author-info {
  display: flex;
  flex-direction: column;
}

.author-name {
  font-weight: bold;
  color: var(--color-text);
  font-size: 14px;
}

.note-date {
  font-size: 12px;
  color: var(--color-text-secondary);
}

.note-actions {
  display: flex;
  gap: 8px;
}

.btn-icon-small {
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.btn-icon-small:hover,
.btn-icon-small:active,
.btn-icon-small:focus {
  background-color: var(--color-background);
  box-shadow: none;
  outline: none;
  transform: none;
}

.note-content {
  margin-top: 8px;
}

.note-content .markdown-preview {
  font-size: 14px;
  line-height: 1.5;
}

/* 资料部分 */
.materials-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.materials-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.material-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 16px;
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  transition: none;
}

.material-item:hover,
.material-item:active,
.material-item:focus {
  box-shadow: none;
  outline: none;
  transform: none;
}

.material-info {
  flex: 1;
  margin-right: 16px;
}

.material-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.material-title {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.material-name {
  font-size: 16px;
  font-weight: bold;
  color: var(--color-text);
  margin: 0;
}

.material-type-tag {
  padding: 4px 8px;
  color: var(--color-text);
  border: 1px solid var(--color-border);
  border-radius: 4px;
  font-size: 12px;
  font-weight: bold;
  margin-right: 8px;
}

.material-type-tag.cloud {
  background-color: var(--color-accent);
}

.material-type-tag.note {
  background-color: var(--color-primary);
}

.cloud-type {
  padding: 4px 8px;
  background-color: var(--color-accent);
  color: var(--color-text);
  border: 1px solid var(--color-border);
  border-radius: 4px;
  font-size: 12px;
  font-weight: bold;
  text-align: center;
}

.material-description {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin: 8px 0;
  line-height: 1.4;
}

.material-meta {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: var(--color-text-secondary);
}

.material-actions {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 8px;
}

.password-info {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--color-text-secondary);
}

.password-label {
  font-weight: bold;
}

.password-value {
  background-color: var(--color-background);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
}

.empty-materials {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
}

.empty-materials p {
  font-size: 16px;
  color: var(--color-text-secondary);
  margin-top: 16px;
}

/* 模态框样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal-content {
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 12px;
  width: 100%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 2px solid var(--color-border);
  background-color: var(--color-primary);
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: bold;
  color: var(--color-text);
}

.modal-close {
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.modal-close:hover,
.modal-close:active,
.modal-close:focus {
  background-color: var(--color-primary);
  box-shadow: none;
  outline: none;
  transform: none;
}

.modal-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-weight: bold;
  color: var(--color-text);
  margin-bottom: 6px;
  font-size: 14px;
}

.form-input,
.form-textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--color-border);
  border-radius: 6px;
  font-size: 14px;
  background-color: var(--color-surface);
  color: var(--color-text);
  box-sizing: border-box;
  box-shadow: none;
}

.form-input:hover,
.form-input:active,
.form-input:focus,
.form-textarea:hover,
.form-textarea:active,
.form-textarea:focus {
  outline: none !important;
  border-color: var(--color-primary);
  box-shadow: none;
  transform: none;
}

/* 网盘类型下拉框样式 */
.cloud-type-dropdown,
.note-type-dropdown {
  position: relative;
  width: 100%;
}

.cloud-type-button {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 6px;
  padding: 10px 12px;
  font-size: 14px;
  font-weight: bold;
  color: var(--color-text);
  cursor: pointer;
  transition: background-color 0.2s ease;
  box-sizing: border-box;
}

.cloud-type-button:hover,
.cloud-type-button:active,
.cloud-type-button:focus {
  background: var(--color-primary);
  box-shadow: none;
  outline: none;
  transform: none;
}

.dropdown-arrow {
  transition: none;
  color: var(--color-text);
}

.dropdown-arrow.rotated {
  transform: rotate(180deg);
}

.cloud-type-options {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  right: 0;
  z-index: 1001;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  overflow: hidden;
  max-height: 200px;
  overflow-y: auto;
}

.cloud-type-option {
  display: block;
  width: 100%;
  background: var(--color-surface);
  border: none;
  padding: 10px 12px;
  font-size: 14px;
  font-weight: bold;
  color: var(--color-text);
  cursor: pointer;
  transition: background-color 0.2s ease;
  text-align: left;
  border-bottom: 1px solid var(--color-border);
}

.cloud-type-option:last-child {
  border-bottom: none;
}

.cloud-type-option:hover,
.cloud-type-option:active,
.cloud-type-option:focus {
  background: var(--color-primary);
  box-shadow: none;
  outline: none;
  transform: none;
}

.cloud-type-option.active {
  background: var(--color-primary);
}

.form-textarea {
  resize: vertical;
  min-height: 80px;
}

/* 资料类型选择器样式 - Segmented Button */
.material-type-selector {
  display: flex;
  background-color: var(--color-background);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  overflow: hidden;
  position: relative;
}

.type-toggle {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 16px;
  background-color: var(--color-surface);
  border: none;
  border-right: 2px solid var(--color-border);
  font-size: 14px;
  font-weight: bold;
  color: var(--color-text);
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
}

.type-toggle:last-child {
  border-right: none;
}

.type-toggle:hover,
.type-toggle:active,
.type-toggle:focus {
  background-color: var(--color-background);
  box-shadow: none;
  outline: none;
  transform: none;
}

.type-toggle.active {
  background-color: var(--color-primary);
  color: var(--color-text);
}

.type-icon {
  flex-shrink: 0;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 20px;
  border-top: 2px solid var(--color-border);
  background-color: var(--color-background);
}

/* 大屏幕优化 */
@media (min-width: 1200px) {
  .course-detail-container {
    padding: 32px;
  }

  .course-content {
    max-width: 1600px;
    gap: 32px;
  }

  .content-grid {
    grid-template-columns: 1.1fr 0.9fr;
    gap: 40px;
  }

  .left-column,
  .right-column {
    gap: 32px;
  }

  .card {
    padding: 24px;
  }

  .course-title {
    font-size: 32px;
  }
}

/* 平板设备优化 */
@media (min-width: 769px) and (max-width: 1024px) {
  .course-detail-container {
    padding: 24px;
  }

  .content-grid {
    grid-template-columns: 1fr 1fr;
    gap: 24px;
  }

  .left-column,
  .right-column {
    gap: 20px;
  }

  .card {
    padding: 18px;
  }
}

/* 中等屏幕优化 */
@media (min-width: 1025px) and (max-width: 1199px) {
  .content-grid {
    grid-template-columns: 1fr 1fr;
    gap: 28px;
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .course-detail-container {
    padding: 16px;
  }

  .content-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .left-column,
  .right-column {
    gap: 16px;
    position: static;
    max-height: none;
    overflow-y: visible;
  }

  .right-column::-webkit-scrollbar {
    display: none;
  }

  .action-buttons {
    flex-direction: column;
  }

  .course-content {
    gap: 16px;
  }

  /* 移动端课程卡片优化排版 */
  .card {
    padding: 16px;
    margin-bottom: 16px;
  }

  .card-title {
    font-size: 18px;
    margin-bottom: 12px;
    font-weight: bold;
  }

  .course-content {
    gap: 12px;
  }

  .course-rating-section {
    border: none;
    background: transparent;
    margin-bottom: 12px;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .course-stats {
    border: none;
    background: transparent;
    margin-bottom: 12px;
    display: flex;
    gap: 16px;
  }

  .course-teacher-section {
    border: none;
    background: transparent;
    margin-bottom: 12px;
    align-items: center;
    gap: 8px;
  }

  .course-card-actions {
    border: none;
    background: transparent;
    gap: 8px;
    justify-content: center;
  }

  .rating-display {
    border: none;
    background: transparent;
    gap: 4px;
    align-items: center;
  }

  .rating-stars {
    border: none;
    background: transparent;
    gap: 1px;
  }

  .star {
    font-size: 14px;
  }

  .rating-text {
    font-size: 14px;
    font-weight: bold;
  }

  .rating-count {
    font-size: 12px;
    color: var(--color-text-secondary);
  }

  .stat-item {
    border: none;
    background: transparent;
    gap: 2px;
    text-align: left;
  }

  .stat-value {
    font-size: 14px;
    font-weight: bold;
  }

  .stat-label {
    font-size: 11px;
    color: var(--color-text-secondary);
  }

  .teacher-label {
    border: none;
    padding: 4px 8px;
    font-size: 11px;
    background-color: var(--color-background);
    border-radius: 4px;
  }

  .teacher-name {
    border: none;
    font-size: 14px;
    font-weight: 500;
  }

  /* 优化其他section元素 */
  .notes-section {
    border: none;
    background: transparent;
    margin-bottom: 12px;
  }

  .materials-section {
    border: none;
    background: transparent;
    margin-bottom: 12px;
  }

  .course-card-header {
    border: none;
    background: transparent;
    margin-bottom: 12px;
    gap: 8px;
  }
}

/* 超小屏幕设备进一步优化 */
@media (max-width: 480px) {
  .card {
    padding: 12px;
    margin-bottom: 12px;
  }

  .course-detail-container {
    padding: 12px;
  }

  .content-grid {
    gap: 12px;
  }

  .left-column,
  .right-column {
    gap: 12px;
  }

  /* 超小屏幕简化排版 */
  .card {
    padding: 14px;
    margin-bottom: 14px;
  }

  .card-title {
    font-size: 17px;
    margin-bottom: 10px;
  }

  .course-content {
    gap: 12px;
  }

  .course-rating-section,
  .course-stats,
  .course-teacher-section,
  .course-card-actions {
    border: none;
    background: transparent;
    margin-bottom: 10px;
  }

  .course-rating-section {
    gap: 6px;
  }

  .course-stats {
    gap: 8px;
  }

  .course-teacher-section {
    gap: 6px;
  }

  .rating-display,
  .rating-stars,
  .stat-item {
    border: none;
    background: transparent;
    gap: 3px;
  }

  .teacher-name {
    font-size: 13px;
  }

  .star {
    font-size: 13px;
  }

  .rating-text {
    font-size: 12px;
  }

  .rating-count {
    font-size: 11px;
  }

  .stat-value {
    font-size: 13px;
  }

  .stat-label {
    font-size: 10px;
  }

  .notes-section,
  .materials-section {
    border: none;
    background: transparent;
    margin-bottom: 10px;
  }

  .rating-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  /* 移动端弹窗样式 */
  .modal-overlay {
    padding: 0;
    align-items: stretch;
  }

  .modal-content {
    width: 100%;
    height: 100vh;
    max-width: none;
    max-height: none;
    border-radius: 0;
    border: none;
    box-shadow: none !important;
  }

  .modal-header {
    border-radius: 0;
    border-left: none;
    border-right: none;
    border-top: none;
  }

  .modal-footer {
    border-radius: 0;
    border-left: none;
    border-right: none;
    border-bottom: none;
    margin-top: auto;
  }

  .modal-body {
    flex: 1;
    overflow-y: auto;
  }

  /* 移动端资料类型选择器 */
  .material-type-selector {
    flex-direction: row;
  }

  .type-toggle {
    padding: 14px 12px;
    font-size: 13px;
  }

  .type-toggle .type-icon {
    width: 14px;
    height: 14px;
  }
}
</style>