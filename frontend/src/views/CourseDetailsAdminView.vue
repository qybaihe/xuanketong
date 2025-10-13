<template>
  <div class="course-details-admin">
    <!-- 页面标题 -->
    <div class="page-header">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item>
          <router-link to="/admin">管理后台</router-link>
        </el-breadcrumb-item>
        <el-breadcrumb-item>
          <router-link to="/admin/courses">课程管理</router-link>
        </el-breadcrumb-item>
        <el-breadcrumb-item>课程详情</el-breadcrumb-item>
      </el-breadcrumb>
      <h1 class="page-title">课程详情</h1>
    </div>

    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="5" animated />
    </div>

    <div v-else-if="course" class="course-content">
      <!-- 基本信息 -->
      <el-card class="info-card">
        <template #header>
          <h3 class="card-title">基本信息</h3>
        </template>
        
        <el-descriptions :column="2" border>
          <el-descriptions-item label="课程ID">
            {{ course.id }}
          </el-descriptions-item>
          <el-descriptions-item label="课程名称">
            <strong>{{ course.name }}</strong>
          </el-descriptions-item>
          <el-descriptions-item label="课程描述">
            {{ course.description || '暂无描述' }}
          </el-descriptions-item>
          <el-descriptions-item label="教师">
            {{ course.teacher || '未指定' }}
          </el-descriptions-item>
          <el-descriptions-item label="学分">
            <el-tag type="warning">{{ course.credits }} 学分</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="年级">
            <el-tag v-if="course.grade" type="primary">{{ course.grade }}</el-tag>
            <span v-else>未指定</span>
          </el-descriptions-item>
          <el-descriptions-item label="学期">
            <el-tag v-if="course.semester" type="info">{{ course.semester }}</el-tag>
            <span v-else>未指定</span>
          </el-descriptions-item>
          <el-descriptions-item label="科目">
            <el-tag v-if="course.subject" type="success">{{ course.subject }}</el-tag>
            <span v-else>未指定</span>
          </el-descriptions-item>
          <el-descriptions-item label="开课单位">
            <el-tag v-if="course.courseUnitName" type="warning">{{ course.courseUnitName }}</el-tag>
            <span v-else>未指定</span>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">
            {{ formatDate(course.created_at) }}
          </el-descriptions-item>
          <el-descriptions-item label="更新时间">
            {{ formatDate(course.updated_at) }}
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 课程统计 -->
      <el-card class="stats-card">
        <template #header>
          <h3 class="card-title">课程统计</h3>
        </template>
        
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12" :md="6">
            <div class="stat-item">
              <div class="stat-icon">
                <el-icon size="24" color="var(--color-primary)">
                  <View />
                </el-icon>
              </div>
              <div class="stat-info">
                <h4 class="stat-value">{{ courseStats.views || 0 }}</h4>
                <p class="stat-label">浏览次数</p>
              </div>
            </div>
          </el-col>
          
          <el-col :xs="24" :sm="12" :md="6">
            <div class="stat-item">
              <div class="stat-icon">
                <el-icon size="24" color="var(--color-success)">
                  <Star />
                </el-icon>
              </div>
              <div class="stat-info">
                <h4 class="stat-value">{{ courseStats.ratings || 0 }}</h4>
                <p class="stat-label">评分次数</p>
              </div>
            </div>
          </el-col>
          
          <el-col :xs="24" :sm="12" :md="6">
            <div class="stat-item">
              <div class="stat-icon">
                <el-icon size="24" color="var(--color-warning)">
                  <ChatDotRound />
                </el-icon>
              </div>
              <div class="stat-info">
                <h4 class="stat-value">{{ courseStats.comments || 0 }}</h4>
                <p class="stat-label">评论数量</p>
              </div>
            </div>
          </el-col>
          
          <el-col :xs="24" :sm="12" :md="6">
            <div class="stat-item">
              <div class="stat-icon">
                <el-icon size="24" color="var(--color-error)">
                  <Document />
                </el-icon>
              </div>
              <div class="stat-info">
                <h4 class="stat-value">{{ courseStats.materials || 0 }}</h4>
                <p class="stat-label">资料数量</p>
              </div>
            </div>
          </el-col>
        </el-row>
      </el-card>

      <!-- 评分详情 -->
      <el-card class="ratings-card">
        <template #header>
          <h3 class="card-title">评分详情</h3>
        </template>
        
        <div v-if="ratings.length > 0">
          <el-table :data="ratings" stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="user.username" label="用户" width="120" />
            <el-table-column prop="rating" label="评分" width="100">
              <template #default="{ row }">
                <el-rate
                  v-model="row.rating"
                  disabled
                  show-score
                  text-color="#ff9900"
                />
              </template>
            </el-table-column>
            <el-table-column prop="comment" label="评论" min-width="200">
              <template #default="{ row }">
                {{ row.comment || '无评论' }}
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="评分时间" width="180">
              <template #default="{ row }">
                {{ formatDate(row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button
                  type="danger"
                  size="small"
                  @click="deleteRating(row.id)"
                >
                  <el-icon><Delete /></el-icon>
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
        
        <el-empty v-else description="暂无评分数据" />
      </el-card>

      <!-- 评论详情 -->
      <el-card class="comments-card">
        <template #header>
          <h3 class="card-title">评论详情</h3>
        </template>
        
        <div v-if="comments.length > 0">
          <el-table :data="comments" stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="user.username" label="用户" width="120" />
            <el-table-column prop="content" label="评论内容" min-width="300" />
            <el-table-column prop="created_at" label="评论时间" width="180">
              <template #default="{ row }">
                {{ formatDate(row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button
                  type="danger"
                  size="small"
                  @click="deleteComment(row.id)"
                >
                  <el-icon><Delete /></el-icon>
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
        
        <el-empty v-else description="暂无评论数据" />
      </el-card>

      <!-- 操作按钮 -->
      <div class="action-buttons">
        <el-button type="primary" @click="editCourse">
          <el-icon><Edit /></el-icon>
          编辑课程
        </el-button>
        <el-button type="danger" @click="deleteCourse">
          <el-icon><Delete /></el-icon>
          删除课程
        </el-button>
        <el-button @click="goBack">
          <el-icon><ArrowLeft /></el-icon>
          返回列表
        </el-button>
      </div>
    </div>

    <el-empty v-else description="课程不存在或已被删除" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/services/api'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const course = ref<any>(null)
const courseStats = ref({
  views: 0,
  ratings: 0,
  comments: 0,
  materials: 0
})
const ratings = ref<any[]>([])
const comments = ref<any[]>([])

const loadCourseDetails = async () => {
  const courseId = route.params.id
  if (!courseId) return

  loading.value = true
  try {
    const response = await api.get(`/admin/courses/${courseId}`)
    if (response.data.data) {
      course.value = response.data.data
    }
  } catch (error: any) {
    console.error('加载课程详情失败:', error)
    ElMessage.error('加载课程详情失败')
  } finally {
    loading.value = false
  }
}

const loadCourseStats = async () => {
  const courseId = route.params.id
  if (!courseId) return

  try {
    // 这里应该调用实际的统计API
    // const response = await api.get(`/admin/courses/${courseId}/stats`)
    // courseStats.value = response.data.data
    
    // 模拟数据
    courseStats.value = {
      views: Math.floor(Math.random() * 1000) + 100,
      ratings: Math.floor(Math.random() * 100) + 10,
      comments: Math.floor(Math.random() * 200) + 20,
      materials: Math.floor(Math.random() * 50) + 5
    }
  } catch (error) {
    console.error('加载课程统计失败:', error)
  }
}

const loadRatings = async () => {
  const courseId = route.params.id
  if (!courseId) return

  try {
    // 这里应该调用实际的评分API
    // const response = await api.get(`/admin/courses/${courseId}/ratings`)
    // ratings.value = response.data.data
    
    // 模拟数据
    ratings.value = [
      {
        id: 1,
        user: { username: '张三' },
        rating: 5,
        comment: '课程内容很丰富，老师讲解很详细',
        created_at: '2024-01-15T10:30:00Z'
      },
      {
        id: 2,
        user: { username: '李四' },
        rating: 4,
        comment: '整体不错，但有些地方可以更深入',
        created_at: '2024-01-14T15:20:00Z'
      }
    ]
  } catch (error) {
    console.error('加载评分数据失败:', error)
  }
}

const loadComments = async () => {
  const courseId = route.params.id
  if (!courseId) return

  try {
    // 这里应该调用实际的评论API
    // const response = await api.get(`/admin/courses/${courseId}/comments`)
    // comments.value = response.data.data
    
    // 模拟数据
    comments.value = [
      {
        id: 1,
        user: { username: '王五' },
        content: '这门课程对我的帮助很大，推荐大家学习',
        created_at: '2024-01-15T09:15:00Z'
      },
      {
        id: 2,
        user: { username: '赵六' },
        content: '老师很负责任，课程质量很高',
        created_at: '2024-01-14T14:30:00Z'
      }
    ]
  } catch (error) {
    console.error('加载评论数据失败:', error)
  }
}

const deleteRating = async (ratingId: number) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这个评分吗？',
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    // 这里应该调用实际的删除API
    // await api.delete(`/admin/ratings/${ratingId}`)
    
    ElMessage.success('评分删除成功')
    loadRatings()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('删除评分失败:', error)
      ElMessage.error('删除评分失败')
    }
  }
}

const deleteComment = async (commentId: number) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这个评论吗？',
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    // 这里应该调用实际的删除API
    // await api.delete(`/admin/comments/${commentId}`)
    
    ElMessage.success('评论删除成功')
    loadComments()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('删除评论失败:', error)
      ElMessage.error('删除评论失败')
    }
  }
}

const editCourse = () => {
  const courseId = route.params.id
  router.push(`/admin/courses/edit/${courseId}`)
}

const deleteCourse = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除课程 "${course.value?.name}" 吗？此操作无法撤销。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'error',
      }
    )
    
    const courseId = route.params.id
    await api.delete(`/admin/courses/${courseId}`)
    
    ElMessage.success('课程删除成功')
    router.push('/admin/courses')
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('删除课程失败:', error)
      ElMessage.error('删除课程失败')
    }
  }
}

const goBack = () => {
  router.push('/admin/courses')
}

const formatDate = (dateString: string) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('zh-CN')
}

onMounted(() => {
  loadCourseDetails()
  loadCourseStats()
  loadRatings()
  loadComments()
})
</script>

<style scoped>
.course-details-admin {
  padding: 20px;
  background-color: var(--color-background);
  min-height: 100vh;
}

.page-header {
  margin-bottom: 30px;
}

.page-title {
  font-size: 2rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 15px 0 0 0;
}

.loading-container {
  padding: 20px;
}

.course-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.info-card,
.stats-card,
.ratings-card,
.comments-card {
  margin-bottom: 20px;
}

.card-title {
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0;
}

.stats-card .el-row {
  margin-top: 20px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px;
  background: var(--color-background);
  border: 1px solid var(--color-border);
  border-radius: 8px;
}

.stat-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: var(--color-background);
  border: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0 0 5px 0;
}

.stat-label {
  font-size: 0.9rem;
  color: var(--color-text-secondary);
  margin: 0;
}

.action-buttons {
  display: flex;
  gap: 15px;
  justify-content: center;
  margin-top: 30px;
  padding: 20px;
  background: var(--color-background);
  border: 1px solid var(--color-border);
  border-radius: 8px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .course-details-admin {
    padding: 15px;
  }
  
  .page-title {
    font-size: 1.5rem;
  }
  
  .stat-item {
    flex-direction: column;
    text-align: center;
    gap: 10px;
  }
  
  .action-buttons {
    flex-direction: column;
    align-items: center;
  }

  /* 描述列表单列显示 */
  :deep(.el-descriptions__body) table {
    table-layout: auto;
  }
  :deep(.el-descriptions__label) {
    white-space: nowrap;
  }
  :deep(.el-descriptions__content) {
    white-space: nowrap;
  }
  /* 表格横向滚动 */
  .ratings-card .el-table__body-wrapper,
  .materials-card .el-table__body-wrapper {
    overflow-x: auto;
  }
  :deep(.el-table th),
  :deep(.el-table td) {
    white-space: nowrap;
  }
}
</style>