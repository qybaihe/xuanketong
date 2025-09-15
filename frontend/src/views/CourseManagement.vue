<template>
  <div class="course-management-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title text-shine">课程管理</h1>
      <p class="page-subtitle">管理系统中的所有课程信息</p>
    </div>

    <!-- 筛选和操作栏 -->
    <div class="filter-bar card-glass">
      <div class="filter-content">
        <div class="filter-group">
          <div class="filter-item">
            <label class="filter-label">年级</label>
            <input type="text" class="filter-input input-glass" placeholder="输入年级..." v-model="filter.grade" @input="debouncedFetchCourses">
          </div>
          <div class="filter-item">
            <label class="filter-label">学期</label>
            <input type="text" class="filter-input input-glass" placeholder="输入学期..." v-model="filter.semester" @input="debouncedFetchCourses">
          </div>
          <div class="filter-item">
            <label class="filter-label">科目</label>
            <input type="text" class="filter-input input-glass" placeholder="输入科目..." v-model="filter.subject" @input="debouncedFetchCourses">
          </div>
        </div>
        
        <div class="action-group">
          <button class="btn btn-primary" @click="showAddModal = true">
            <span class="btn-icon">➕</span>
            添加课程
          </button>
        </div>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <div class="loader"></div>
      <p class="loading-text">正在加载课程数据...</p>
    </div>

    <!-- 课程列表 -->
    <div v-else class="course-list card-glass">
      <div class="course-list-header">
        <h2 class="list-title">课程列表</h2>
        <p class="list-count">共 {{ courses.length }} 门课程</p>
      </div>
      
      <div class="course-table-container">
        <table class="course-table">
          <thead>
            <tr>
              <th class="table-header">ID</th>
              <th class="table-header">名称</th>
              <th class="table-header">教师</th>
              <th class="table-header">年级</th>
              <th class="table-header">学期</th>
              <th class="table-header">科目</th>
              <th class="table-header">学分</th>
              <th class="table-header">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="course in courses" :key="course.ID" class="course-row">
              <td class="table-cell">{{ course.ID }}</td>
              <td class="table-cell course-name">{{ course.Name }}</td>
              <td class="table-cell">{{ course.Teacher }}</td>
              <td class="table-cell">
                <span class="tag tag-primary">{{ course.Grade }}</span>
              </td>
              <td class="table-cell">
                <span class="tag tag-secondary">{{ course.Semester }}</span>
              </td>
              <td class="table-cell">
                <span class="tag tag-accent">{{ course.Subject }}</span>
              </td>
              <td class="table-cell">
                <span class="credit-badge">{{ course.Credits }} 学分</span>
              </td>
              <td class="table-cell actions-cell">
                <button class="btn-icon-btn btn-view" @click="viewCourse(course.ID)" title="查看">
                  <span class="icon">👁️</span>
                </button>
                <button class="btn-icon-btn btn-edit" @click="editCourse(course)" title="编辑">
                  <span class="icon">✏️</span>
                </button>
                <button class="btn-icon-btn btn-delete" @click="confirmDelete(course.ID)" title="删除">
                  <span class="icon">🗑️</span>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
        
        <div v-if="courses.length === 0" class="empty-state">
          <div class="empty-icon">📚</div>
          <h3 class="empty-title">暂无课程</h3>
          <p class="empty-description">请调整筛选条件或添加新课程</p>
        </div>
      </div>
    </div>

    <!-- 添加/编辑课程模态框 -->
    <div v-if="showAddModal" class="modal-overlay" @click="closeModalOnOverlay">
      <div class="modal-container card-glass" @click.stop>
        <div class="modal-header">
          <h3 class="modal-title">{{ editingCourse ? '编辑课程' : '添加课程' }}</h3>
          <button class="modal-close" @click="closeModal">
            <span class="close-icon">✕</span>
          </button>
        </div>
        
        <div class="modal-body">
          <form @submit.prevent="saveCourse" class="course-form">
            <div class="form-group">
              <label class="form-label">课程名称 *</label>
              <input type="text" class="form-input input-glass" id="name" v-model="courseForm.Name" required>
            </div>
            
            <div class="form-group">
              <label class="form-label">课程描述</label>
              <textarea class="form-textarea input-glass" id="description" v-model="courseForm.Description" rows="3"></textarea>
            </div>
            
            <div class="form-row">
              <div class="form-group">
                <label class="form-label">教师</label>
                <input type="text" class="form-input input-glass" id="teacher" v-model="courseForm.Teacher">
              </div>
              
              <div class="form-group">
                <label class="form-label">学分</label>
                <input type="number" class="form-input input-glass" id="credits" v-model.number="courseForm.Credits" min="0" step="0.5">
              </div>
            </div>
            
            <div class="form-row">
              <div class="form-group">
                <label class="form-label">年级</label>
                <input type="text" class="form-input input-glass" id="grade" v-model="courseForm.Grade">
              </div>
              
              <div class="form-group">
                <label class="form-label">学期</label>
                <input type="text" class="form-input input-glass" id="semester" v-model="courseForm.Semester">
              </div>
              
              <div class="form-group">
                <label class="form-label">科目</label>
                <input type="text" class="form-input input-glass" id="subject" v-model="courseForm.Subject">
              </div>
            </div>
            
            <div class="form-actions">
              <button type="button" class="btn btn-secondary" @click="closeModal">取消</button>
              <button type="submit" class="btn btn-primary" :disabled="saving">
                <span v-if="saving" class="btn-icon">⏳</span>
                <span v-else>保存</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- 删除确认模态框 -->
    <div v-if="showDeleteModal" class="modal-overlay" @click="closeDeleteModalOnOverlay">
      <div class="modal-container card-glass modal-sm" @click.stop>
        <div class="modal-header">
          <h3 class="modal-title">确认删除</h3>
          <button class="modal-close" @click="closeDeleteModal">
            <span class="close-icon">✕</span>
          </button>
        </div>
        
        <div class="modal-body">
          <p class="delete-confirm-text">您确定要删除课程 "{{ courseToDelete?.Name }}" 吗？此操作不可撤销。</p>
          
          <div class="form-actions">
            <button type="button" class="btn btn-secondary" @click="closeDeleteModal">取消</button>
            <button type="button" class="btn btn-danger" @click="deleteCourse" :disabled="deleting">
              <span v-if="deleting" class="btn-icon">⏳</span>
              <span v-else>删除</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 成功提示 -->
    <div v-if="showSuccessMessage" class="success-toast">
      <div class="toast-content">
        <span class="toast-icon">✓</span>
        <span class="toast-message">{{ successMessage }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

interface Course {
  ID: number;
  Name: string;
  Description: string;
  Teacher: string;
  Grade: string;
  Semester: string;
  Subject: string;
  Credits: number;
  ImageURL?: string;
}

const courses = ref<Course[]>([])
const showAddModal = ref(false)
const showDeleteModal = ref(false)
const showSuccessMessage = ref(false)
const editingCourse = ref<Course | null>(null)
const courseToDelete = ref<Course | null>(null)
const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const successMessage = ref('')
const router = useRouter()

const filter = reactive({
  grade: '',
  semester: '',
  subject: ''
})

const courseForm = reactive({
  ID: 0,
  Name: '',
  Description: '',
  Teacher: '',
  Grade: '',
  Semester: '',
  Subject: '',
  Credits: 0,
  ImageURL: ''
})

// 防抖处理
let debounceTimer: number
const debouncedFetchCourses = () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(fetchCourses, 300)
}

const fetchCourses = async () => {
  try {
    loading.value = true
    const params = new URLSearchParams()
    if (filter.grade) params.append('grade', filter.grade)
    if (filter.semester) params.append('semester', filter.semester)
    if (filter.subject) params.append('subject', filter.subject)
    const response = await axios.get('/api/v1/courses', { params })
    courses.value = response.data.data
  } catch (error) {
    console.error('获取课程失败:', error)
    showErrorMessage('获取课程列表失败')
  } finally {
    loading.value = false
  }
}

// 测试后端路由是否正常工作
const testBackendRoutes = async () => {
  try {
    // 测试ping路由
    const pingResponse = await axios.get('/ping')
    console.log('Ping路由响应:', pingResponse.data)
    
    // 测试管理员路由
    const adminResponse = await axios.get('/admin-routes')
    console.log('管理员路由响应:', adminResponse.data)
    
    // 测试管理员测试路由
    const testResponse = await axios.post('/api/v1/admin/courses/test', {})
    console.log('管理员测试路由响应:', testResponse.data)
  } catch (error) {
    console.error('测试路由失败:', error)
  }
}

const closeModal = () => {
  showAddModal.value = false
  editingCourse.value = null
  Object.assign(courseForm, {
    ID: 0,
    Name: '',
    Description: '',
    Teacher: '',
    Grade: '',
    Semester: '',
    Subject: '',
    Credits: 0,
    ImageURL: ''
  })
}

const closeModalOnOverlay = (event: MouseEvent) => {
  if (event.target === event.currentTarget) {
    closeModal()
  }
}

const closeDeleteModal = () => {
  showDeleteModal.value = false
  courseToDelete.value = null
}

const closeDeleteModalOnOverlay = (event: MouseEvent) => {
  if (event.target === event.currentTarget) {
    closeDeleteModal()
  }
}

const editCourse = (course: Course) => {
  editingCourse.value = course
  Object.assign(courseForm, course)
  showAddModal.value = true
}

const viewCourse = (id: number) => {
  router.push({ name: 'course-details-admin', params: { id } })
}

const saveCourse = async () => {
  try {
    saving.value = true
    
    // 确保发送的数据格式正确
    const courseData = {
      Name: courseForm.Name,
      Description: courseForm.Description,
      Teacher: courseForm.Teacher,
      Grade: courseForm.Grade,
      Semester: courseForm.Semester,
      Subject: courseForm.Subject,
      Credits: courseForm.Credits,
      ImageURL: courseForm.ImageURL || ''
    }
    
    if (editingCourse.value) {
      // 先测试测试路由
      try {
        const testResponse = await axios.post('/api/v1/admin/courses/test', {})
        console.log('测试路由响应:', testResponse.data)
      } catch (testError) {
        console.error('测试路由失败:', testError)
      }
      
      // 然后尝试实际的更新请求
      await axios.put(`/api/v1/admin/courses/${courseForm.ID}`, courseData)
      showSuccessMessageFunc('课程更新成功')
    } else {
      // 先测试测试路由
      try {
        const testResponse = await axios.post('/api/v1/admin/courses/test', {})
        console.log('测试路由响应:', testResponse.data)
      } catch (testError) {
        console.error('测试路由失败:', testError)
      }
      
      // 然后尝试实际的添加请求
      await axios.post('/api/v1/admin/courses', courseData)
      showSuccessMessageFunc('课程添加成功')
    }
    closeModal()
    fetchCourses()
  } catch (error) {
    console.error('保存课程失败:', error)
    const axiosError = error as any
    if (axiosError.response) {
      console.error('错误响应:', axiosError.response.data)
      console.error('错误状态:', axiosError.response.status)
    }
    showErrorMessage(editingCourse.value ? '更新课程失败' : '添加课程失败')
  } finally {
    saving.value = false
  }
}

const confirmDelete = (id: number) => {
  const course = courses.value.find(c => c.ID === id)
  if (course) {
    courseToDelete.value = course
    showDeleteModal.value = true
  }
}

const deleteCourse = async () => {
  if (!courseToDelete.value) return
  
  try {
    deleting.value = true
    await axios.delete(`/api/v1/admin/courses/${courseToDelete.value.ID}`)
    showSuccessMessageFunc('课程删除成功')
    closeDeleteModal()
    fetchCourses()
  } catch (error) {
    console.error('删除课程失败:', error)
    showErrorMessage('删除课程失败')
  } finally {
    deleting.value = false
  }
}

const showSuccessMessageFunc = (message: string) => {
  successMessage.value = message
  showSuccessMessage.value = true
  setTimeout(() => {
    showSuccessMessage.value = false
  }, 3000)
}

const showErrorMessage = (message: string) => {
  // 在实际应用中，这里可以使用更友好的错误提示方式
  alert(message)
}

onMounted(() => {
  fetchCourses()
  testBackendRoutes()
})
</script>

<style scoped>
/* 课程管理容器 */
.course-management-container {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xl);
}

/* 页面标题 */
.page-header {
  text-align: center;
}

.page-title {
  font-size: var(--font-size-title1);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin: 0 0 var(--spacing-sm) 0;
}

.page-subtitle {
  font-size: var(--font-size-body);
  color: var(--text-secondary);
  margin: 0;
}

/* 筛选栏 */
.filter-bar {
  padding: var(--spacing-lg);
  border-radius: 16px;
}

.filter-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  gap: var(--spacing-md);
}

.filter-group {
  display: flex;
  gap: var(--spacing-md);
  flex-wrap: wrap;
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.filter-label {
  font-size: var(--font-size-body2);
  font-weight: var(--font-weight-medium);
  color: var(--text-secondary);
}

.filter-input {
  width: 150px;
}

.action-group {
  display: flex;
  gap: var(--spacing-sm);
}

/* 课程列表 */
.course-list {
  padding: var(--spacing-lg);
  border-radius: 16px;
}

.course-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
}

.list-title {
  font-size: var(--font-size-title2);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin: 0;
}

.list-count {
  font-size: var(--font-size-body2);
  color: var(--text-tertiary);
  margin: 0;
}

.course-table-container {
  overflow-x: auto;
}

.course-table {
  width: 100%;
  border-collapse: collapse;
}

.table-header {
  text-align: left;
  padding: var(--spacing-sm) var(--spacing-md);
  font-size: var(--font-size-body2);
  font-weight: var(--font-weight-semibold);
  color: var(--text-secondary);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.course-row {
  transition: background-color 0.2s ease;
}

.course-row:hover {
  background-color: rgba(255, 255, 255, 0.05);
}

.table-cell {
  padding: var(--spacing-md);
  font-size: var(--font-size-body);
  color: var(--text-primary);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.course-name {
  font-weight: var(--font-weight-medium);
}

.tag {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: var(--font-size-caption);
  font-weight: var(--font-weight-medium);
}

.tag-primary {
  background: var(--success-light);
  color: var(--success-base);
}

.tag-secondary {
  background: var(--info-light);
  color: var(--info-base);
}

.tag-accent {
  background: var(--warning-light);
  color: var(--warning-base);
}

.credit-badge {
  background: var(--gradient-primary);
  color: white;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: var(--font-size-caption);
  font-weight: var(--font-weight-medium);
}

.actions-cell {
  width: 120px;
}

.btn-icon-btn {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  margin-right: var(--spacing-xs);
  transition: all 0.2s ease;
}

.btn-icon-btn:hover {
  background: rgba(255, 255, 255, 0.1);
}

.btn-view:hover {
  color: var(--info-color);
}

.btn-edit:hover {
  color: var(--warning-color);
}

.btn-delete:hover {
  color: var(--error-color);
}

.icon {
  font-size: 16px;
}

/* 空状态 */
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
  margin: 0;
}

/* 模态框 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(5px);
  -webkit-backdrop-filter: blur(5px);
}

.modal-container {
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  border-radius: 16px;
  animation: modalFadeIn 0.3s ease;
}

@keyframes modalFadeIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-lg);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.modal-title {
  font-size: var(--font-size-title2);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin: 0;
}

.modal-close {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: none;
  background: transparent;
  color: var(--text-tertiary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.modal-close:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--text-primary);
}

.close-icon {
  font-size: 16px;
}

.modal-body {
  padding: var(--spacing-lg);
}

.course-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.form-label {
  font-size: var(--font-size-body2);
  font-weight: var(--font-weight-medium);
  color: var(--text-secondary);
}

.form-input,
.form-textarea {
  width: 100%;
  padding: var(--spacing-sm);
  border: 1px solid var(--separator-color);
  border-radius: 8px;
  font-size: var(--font-size-body);
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  transition: all 0.2s ease;
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(47, 169, 20, 0.1);
}

.form-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: var(--spacing-md);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-md);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .filter-content {
    flex-direction: column;
    align-items: stretch;
  }
  
  .filter-group {
    justify-content: center;
  }
  
  .action-group {
    justify-content: center;
    margin-top: var(--spacing-md);
  }
  
  .course-table {
    font-size: var(--font-size-body2);
  }
  
  .table-cell {
    padding: var(--spacing-sm);
  }
  
  .actions-cell {
    width: 100px;
  }
  
  .btn-icon-btn {
    width: 28px;
    height: 28px;
  }
  
  .form-row {
    grid-template-columns: 1fr;
  }
  
  .modal-container {
    width: 95%;
    padding: var(--spacing-md);
  }
}

@media (min-width: 769px) and (max-width: 1024px) {
  .filter-group {
    flex-wrap: wrap;
  }
  
  .filter-input {
    width: 120px;
  }
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

/* 删除确认模态框 */
.modal-sm {
  max-width: 400px;
}

.delete-confirm-text {
  font-size: var(--font-size-body);
  color: var(--text-primary);
  margin: 0 0 var(--spacing-lg) 0;
  line-height: 1.5;
}

/* 成功提示 */
.success-toast {
  position: fixed;
  bottom: var(--spacing-xl);
  right: var(--spacing-xl);
  z-index: 2000;
  animation: slideInRight 0.3s ease, slideOutRight 0.3s ease 2.7s;
}

@keyframes slideInRight {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

@keyframes slideOutRight {
  from {
    transform: translateX(0);
    opacity: 1;
  }
  to {
    transform: translateX(100%);
    opacity: 0;
  }
}

.toast-content {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  background: var(--success-base);
  color: white;
  padding: var(--spacing-md) var(--spacing-lg);
  border-radius: 8px;
  box-shadow: var(--shadow-lg);
}

.toast-icon {
  font-size: 20px;
}

.toast-message {
  font-size: var(--font-size-body);
  font-weight: var(--font-weight-medium);
}

/* 按钮加载状态 */
.btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}
</style>