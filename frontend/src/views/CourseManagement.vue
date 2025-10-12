<template>
  <div class="course-management">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">课程管理</h1>
      <p class="page-subtitle">管理系统中的所有课程信息</p>
      </div>

    <!-- 操作栏 -->
    <el-card class="filter-card">
      <el-row :gutter="20" align="center">
        <el-col :xs="24" :sm="8" :md="6">
          <el-form-item label="年级">
            <el-input
              v-model="filters.grade"
              placeholder="输入年级..."
              clearable
              @input="handleSearch"
            />
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="8" :md="6">
          <el-form-item label="学期">
            <el-input
              v-model="filters.semester"
              placeholder="输入学期..."
              clearable
              @input="handleSearch"
            />
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="8" :md="6">
          <el-form-item label="科目">
            <el-input
              v-model="filters.subject"
              placeholder="输入科目..."
              clearable
              @input="handleSearch"
            />
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="24" :md="6">
          <el-button type="primary" @click="showAddCourseModal">
            <el-icon><Plus /></el-icon>
            添加课程
          </el-button>
        </el-col>
      </el-row>
    </el-card>

      <!-- 课程列表 -->
    <el-card class="table-card">
      <template #header>
        <h3 class="card-title">课程列表</h3>
      </template>
      
      <el-table
        v-loading="loading"
        :data="courses"
        stripe
        style="width: 100%"
        empty-text="暂无课程数据"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="课程名称" min-width="150">
          <template #default="{ row }">
            <strong>{{ row.name }}</strong>
          </template>
        </el-table-column>
        <el-table-column prop="teacher" label="教师" width="120">
          <template #default="{ row }">
            {{ row.teacher || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="grade" label="年级" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.grade" type="primary" size="small">
              {{ row.grade }}
            </el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="semester" label="学期" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.semester" type="info" size="small">
              {{ row.semester }}
            </el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="subject" label="科目" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.subject" type="success" size="small">
              {{ row.subject }}
            </el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="credits" label="学分" width="100">
          <template #default="{ row }">
            <el-tag type="warning" size="small">
              {{ row.credits }} 学分
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              @click="editCourse(row.id)"
            >
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="deleteCourse(row.id, row.name)"
            >
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑课程对话框 -->
    <el-dialog
      v-model="courseDialogVisible"
      :title="isEdit ? '编辑课程' : '添加课程'"
      width="600px"
    >
      <el-form
        ref="courseFormRef"
        :model="courseForm"
        :rules="courseRules"
        label-width="100px"
      >
        <el-form-item label="课程名称" prop="name">
          <el-input v-model="courseForm.name" placeholder="请输入课程名称" />
        </el-form-item>
        <el-form-item label="课程描述">
          <el-input
            v-model="courseForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入课程描述"
          />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="教师">
              <el-input v-model="courseForm.teacher" placeholder="请输入教师姓名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="学分" prop="credits">
              <el-input-number
                v-model="courseForm.credits"
                :min="0"
                :step="0.5"
                placeholder="请输入学分"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="年级">
              <el-input v-model="courseForm.grade" placeholder="请输入年级" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="学期">
              <el-input v-model="courseForm.semester" placeholder="请输入学期" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="科目">
              <el-input v-model="courseForm.subject" placeholder="请输入科目" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      
      <template #footer>
        <el-button @click="courseDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveCourse" :loading="saving">
          {{ isEdit ? '更新' : '添加' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 删除确认对话框 -->
    <el-dialog
      v-model="deleteDialogVisible"
      title="确认删除"
      width="400px"
    >
      <p>确定要删除课程 "<strong>{{ deleteCourseName }}</strong>" 吗？此操作无法撤销。</p>
      
      <template #footer>
        <el-button @click="deleteDialogVisible = false">取消</el-button>
        <el-button type="danger" @click="confirmDelete" :loading="deleting">
          删除
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/services/api'

const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const courses = ref<any[]>([])
const courseDialogVisible = ref(false)
const deleteDialogVisible = ref(false)
const isEdit = ref(false)
const currentDeleteId = ref<number | null>(null)
const deleteCourseName = ref('')

const filters = reactive({
  grade: '',
  semester: '',
  subject: ''
})

const courseForm = reactive({
  id: null,
  name: '',
  description: '',
  teacher: '',
  credits: 0,
  grade: '',
  semester: '',
  subject: '',
  imageURL: ''
})

const courseRules = {
  name: [
    { required: true, message: '请输入课程名称', trigger: 'blur' }
  ],
  credits: [
    { required: true, message: '请输入学分', trigger: 'blur' }
  ]
}

const courseFormRef = ref()

const loadCourses = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (filters.grade) params.append('grade', filters.grade)
    if (filters.semester) params.append('semester', filters.semester)
    if (filters.subject) params.append('subject', filters.subject)
    
    const response = await api.get(`/admin/courses?${params}`)
    
    if (response.data.data) {
      courses.value = response.data.data
    } else {
      courses.value = []
    }
  } catch (error: any) {
    console.error('加载课程失败:', error)
    ElMessage.error('加载课程数据失败，请检查网络连接')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  loadCourses()
}

const showAddCourseModal = () => {
  isEdit.value = false
  resetForm()
  courseDialogVisible.value = true
}

const editCourse = async (id: number) => {
  try {
    const response = await api.get(`/admin/courses/${id}`)
    
    if (response.data.data) {
      const course = response.data.data
      Object.assign(courseForm, course)
      isEdit.value = true
      courseDialogVisible.value = true
    }
  } catch (error: any) {
    console.error('加载课程详情失败:', error)
    ElMessage.error('加载课程详情失败')
  }
}

const saveCourse = async () => {
  if (!courseFormRef.value) return
  
  try {
    await courseFormRef.value.validate()
    saving.value = true
    
    const courseData = { ...courseForm }
    let response
    
    if (isEdit.value) {
      response = await api.put(`/admin/courses/${courseData.id}`, courseData)
    } else {
      response = await api.post(`/admin/courses`, courseData)
    }

    if (response.status === 200 || response.status === 201) {
      courseDialogVisible.value = false
      loadCourses()
      ElMessage.success(isEdit.value ? '课程更新成功' : '课程添加成功')
    }
  } catch (error: any) {
    console.error('保存课程失败:', error)
    const errorMessage = error.response?.data?.error || '未知错误'
    ElMessage.error('操作失败: ' + errorMessage)
  } finally {
    saving.value = false
  }
}

const deleteCourse = (id: number, name: string) => {
  currentDeleteId.value = id
  deleteCourseName.value = name
  deleteDialogVisible.value = true
}

const confirmDelete = async () => {
  if (!currentDeleteId.value) return

  try {
    deleting.value = true
    const response = await api.delete(`/admin/courses/${currentDeleteId.value}`)

    if (response.status === 200) {
      deleteDialogVisible.value = false
      loadCourses()
      ElMessage.success('课程删除成功')
    }
  } catch (error: any) {
    console.error('删除课程失败:', error)
    const errorMessage = error.response?.data?.error || '未知错误'
    ElMessage.error('删除失败: ' + errorMessage)
  } finally {
    deleting.value = false
  }
}

const resetForm = () => {
  Object.assign(courseForm, {
    id: null,
    name: '',
    description: '',
    teacher: '',
    credits: 0,
    grade: '',
    semester: '',
    subject: '',
    imageURL: ''
  })
  courseFormRef.value?.clearValidate()
}

onMounted(() => {
  loadCourses()
})
</script>

<style scoped>
.course-management {
  padding: 20px;
  background-color: var(--color-background);
  min-height: 100vh;
}

.page-header {
  margin-bottom: 30px;
  text-align: center;
}

.page-title {
  font-size: 2rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0 0 10px 0;
}

.page-subtitle {
  font-size: 1rem;
  color: var(--color-text-secondary);
  margin: 0;
}

.filter-card {
  margin-bottom: 20px;
}

.filter-card .el-form-item {
  margin-bottom: 0;
}

.table-card {
  margin-bottom: 20px;
}

.card-title {
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .course-management {
    padding: 15px;
  }
  
  .page-title {
    font-size: 1.5rem;
  }
}
</style>