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
          <el-button type="success" style="margin-left: 8px" @click="openBatchImportDialog">
            批量导入
          </el-button>
          <el-button type="danger" style="margin-left: 8px" :disabled="!selectedCourseIds.length" @click="confirmBatchDelete">
            批量删除
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
        ref="coursesTableRef"
        v-loading="loading"
        :data="courses"
        stripe
        style="width: 100%"
        empty-text="暂无课程数据"
        @selection-change="onSelectionChange"
      >
        <el-table-column type="selection" width="48" />
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
        <el-table-column prop="courseUnitName" label="开课单位" min-width="140">
          <template #default="{ row }">
            {{ row.courseUnitName || '-' }}
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
        <el-row :gutter="20">
          <el-col :span="24">
            <el-form-item label="开课单位">
              <el-input v-model="courseForm.courseUnitName" placeholder="例如：软件工程学院/体育部" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="24">
            <el-form-item label="课程分类">
              <el-select v-model="(courseForm as any).courseCategory" placeholder="选择分类" style="width: 100%">
                <el-option label="公必" value="公必" />
                <el-option label="专必" value="专必" />
                <el-option label="公选" value="公选" />
                <el-option label="专选" value="专选" />
                <el-option label="其他" value="其他" />
              </el-select>
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

    <!-- 批量导入对话框 -->
    <el-dialog v-model="batchDialogVisible" title="批量导入课程" fullscreen>
      <el-form label-width="100px">
        <el-form-item label="统一年级">
          <el-input v-model="batchGrade" placeholder="例如：大二/2022级" />
        </el-form-item>
        <el-form-item label="JSON 粘贴区">
          <el-input
            v-model="batchRawJson"
            type="textarea"
            :rows="10"
            placeholder="粘贴提供的JSON（包含 total 与 rows 数组）"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="previewBatch">解析预览</el-button>
        </el-form-item>
      </el-form>

      <div v-if="batchPreview.length">
        <p>共解析到 {{ batchPreview.length }} 条，已选 {{ batchSelected.length }} 条：</p>
        <el-table ref="batchTableRef" :data="batchPreview" style="width: 100%" :height="batchTableHeight" @selection-change="onBatchSelectChange">
          <el-table-column type="selection" width="48" />
          <el-table-column prop="Name" label="课程名称" min-width="160" />
          <!--<el-table-column prop="Description" label="课程描述" min-width="180" />-->
          <el-table-column prop="Semester" label="学期" min-width="120" />
          <el-table-column prop="Teacher" label="教师" min-width="120" />
          <el-table-column prop="CourseUnitName" label="开课单位" min-width="160" />
          <el-table-column prop="Grade" label="年级" min-width="120" />
          <el-table-column prop="Subject" label="科目" min-width="140">
            <template #default="{ row }">
              <el-select v-model="row.Subject" placeholder="选择科目">
                <el-option label="思政" value="思政" />
                <el-option label="英语" value="英语" />
                <el-option label="数学" value="数学" />
                <el-option label="物理" value="物理" />
                <el-option label="计算机" value="计算机" />
                <el-option label="其他" value="其他" />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column prop="Credits" label="学分" min-width="160">
            <template #default="{ row }">
              <el-input-number v-model="row.Credits" :min="0" :step="0.5" />
            </template>
          </el-table-column>
          <el-table-column prop="CourseCategory" label="课程分类" min-width="160">
            <template #default="{ row }">
              <el-select v-model="row.CourseCategory" placeholder="选择分类">
                <el-option label="公必" value="公必" />
                <el-option label="专必" value="专必" />
                <el-option label="公选" value="公选" />
                <el-option label="专选" value="专选" />
                <el-option label="其他" value="其他" />
              </el-select>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <template #footer>
        <el-button @click="batchDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="batchSubmitting" @click="submitBatch" :disabled="!batchSelected.length">确认提交 ({{ batchSelected.length }})</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api, { type Course } from '@/services/api'

const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const courses = ref<Course[]>([])
const coursesTableRef = ref()
const selectedCourseIds = ref<number[]>([])
const courseDialogVisible = ref(false)
const deleteDialogVisible = ref(false)
const batchDialogVisible = ref(false)
const isEdit = ref(false)
const currentDeleteId = ref<number | null>(null)
const deleteCourseName = ref('')
const batchRawJson = ref('')
type BatchCoursePayload = {
  Name: string
  Description: string
  Grade: string
  Semester: string
  Subject: string
  Teacher: string
  Credits: number
  ImageURL: string
  CourseUnitName: string
}
const batchPreview = ref<BatchCoursePayload[]>([])
const batchSelected = ref<BatchCoursePayload[]>([])
const batchTableRef = ref<{ toggleAllSelection?: () => void; clearSelection?: () => void } | null>(null)
const batchGrade = ref('')
const batchSubmitting = ref(false)
const batchTableHeight = ref<number | string>('60vh')

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
  imageURL: '',
  courseUnitName: ''
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
  } catch (error: unknown) {
    console.error('加载课程失败:', error)
    ElMessage.error('加载课程数据失败，请检查网络连接')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  loadCourses()
}

const onSelectionChange = (selection: Course[]) => {
  selectedCourseIds.value = selection.map((c) => c.id)
}

const confirmBatchDelete = async () => {
  if (!selectedCourseIds.value.length) return
  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedCourseIds.value.length} 门课程吗？此操作不可撤销。`, '确认批量删除', { type: 'warning' })
  } catch {
    return
  }
  try {
    deleting.value = true
    const ids = [...selectedCourseIds.value]
    const results = await Promise.allSettled(ids.map((id) => api.delete(`/admin/courses/${id}`)))
    const success = results.filter(r => r.status === 'fulfilled').length
    const failed = results.length - success
    if (success) ElMessage.success(`成功删除 ${success} 门课程${failed ? `，失败 ${failed}` : ''}`)
    if (failed) ElMessage.error(`有 ${failed} 门删除失败`)
    const table = coursesTableRef.value as { clearSelection?: () => void } | null
    table?.clearSelection?.()
    selectedCourseIds.value = []
    loadCourses()
  } finally {
    deleting.value = false
  }
}

const showAddCourseModal = () => {
  isEdit.value = false
  resetForm()
  courseDialogVisible.value = true
}

const openBatchImportDialog = () => {
  batchDialogVisible.value = true
}

type IncomingRow = {
  courseName?: string
  courseNameEng?: string
  yearTerm?: string
  classMessage?: string
  courseUnitName?: string
  courseUnitNum?: string | null
  courseUnit?: string | null
}

const cnIndex = (i: string) => ({
  '1': '第一',
  '2': '第二',
  '3': '第三',
  '4': '第四'
}[i] || i)

const formatSemester = (term?: string): string => {
  if (!term) return ''
  const m = term.match(/^(\d{4})-(\d)$/)
  if (m) {
    return `${cnIndex(m[2])}学期`
  }
  return term
}

const previewBatch = () => {
  try {
    const parsed = JSON.parse(batchRawJson.value)
    const rows: IncomingRow[] = Array.isArray(parsed)
      ? parsed
      : Array.isArray(parsed?.data?.rows)
        ? parsed.data.rows
        : []

    if (!rows.length) {
      ElMessage.warning('未解析到 rows 数据')
      batchPreview.value = []
      return
    }

    const toCreate = rows.map((r) => {
      const teacher = (r.classMessage || '').split(';')[0] || ''
      const rawCredit = (r as unknown as { credit?: number | string })?.credit
      const creditNum = typeof rawCredit === 'string' ? parseFloat(rawCredit) : (rawCredit ?? 0)
      const subjectGuess = (() => {
        const name = (r.courseName || '').toLowerCase()
        if (name.includes('思政')) return '思政'
        if (name.includes('英语') || name.includes('english')) return '英语'
        if (name.includes('数学') || name.includes('math')) return '数学'
        if (name.includes('物理') || name.includes('physics')) return '物理'
        if (name.includes('程序') || name.includes('计算机') || name.includes('computer')) return '计算机'
        return '其他'
      })()
      return {
        Name: r.courseName || '',
        Description: r.courseNameEng || '',
        Semester: formatSemester(r.yearTerm || ''),
        Teacher: teacher,
        Subject: subjectGuess,
        Credits: Number.isFinite(creditNum as number) ? (creditNum as number) : 0,
        Grade: batchGrade.value || '',
        ImageURL: '',
        CourseUnitName: r.courseUnitName || r.courseUnit || '',
        CourseCategory: (r as unknown as { courseCategory?: string })?.courseCategory || ''
      }
    })

    batchPreview.value = toCreate
    if (!toCreate.length) {
      ElMessage.warning('没有可预览的数据')
    } else {
      ElMessage.success('解析完成，可检查并确认提交')
      nextTick(() => {
        batchTableRef.value?.clearSelection?.()
        batchTableRef.value?.toggleAllSelection?.()
      })
    }
  } catch (e: unknown) {
    console.error(e)
    ElMessage.error('JSON 解析失败，请检查格式')
    batchPreview.value = []
  }
}

const onBatchSelectChange = (selection: BatchCoursePayload[]) => {
  batchSelected.value = selection
}

const submitBatch = async () => {
  if (!batchSelected.value.length) {
    ElMessage.warning('请先选择需要导入的课程')
    return
  }
  try {
    batchSubmitting.value = true
    const resp = await api.post('/admin/courses/batch', batchSelected.value)
    if (resp.status === 200) {
      ElMessage.success(`批量创建成功，共 ${resp.data?.count ?? batchSelected.value.length} 条`)
      batchDialogVisible.value = false
      batchPreview.value = []
      batchSelected.value = []
      batchRawJson.value = ''
      batchGrade.value = ''
      loadCourses()
    }
  } catch (error: unknown) {
    console.error('批量创建失败:', error)
    const msg = (error as { response?: { data?: { error?: string }}}).response?.data?.error || '未知错误'
    ElMessage.error('批量创建失败: ' + msg)
  } finally {
    batchSubmitting.value = false
  }
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
  } catch (error: unknown) {
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
  } catch (error: unknown) {
    console.error('保存课程失败:', error)
    const errorMessage = (error as { response?: { data?: { error?: string }}}).response?.data?.error || '未知错误'
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
  } catch (error: unknown) {
    console.error('删除课程失败:', error)
    const errorMessage = (error as { response?: { data?: { error?: string }}}).response?.data?.error || '未知错误'
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
  const calc = () => {
    // 顶部表单和边距等占位，给表格留足高度
    const reserve = 260
    const h = Math.max(window.innerHeight - reserve, 300)
    batchTableHeight.value = h
  }
  calc()
  window.addEventListener('resize', calc)
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
  /* 过滤栏列纵向堆叠 */
  .filter-card .el-col {
    margin-bottom: 12px;
    flex: 0 0 100%;
    max-width: 100%;
  }
  /* 表格区域横向滚动 */
  .table-card .el-table__body-wrapper {
    overflow-x: auto;
  }
  :deep(.el-table th),
  :deep(.el-table td) {
    white-space: nowrap;
  }
}

/* 手机竖屏全屏对话框 */
@media (max-width: 576px) {
  :deep(.el-dialog) {
    width: 100% !important;
    margin: 0 !important;
  }
  .card-title { font-size: 1rem; }
}
</style>