<template>
  <div class="user-management">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">用户管理</h1>
      <p class="page-subtitle">管理系统中的所有用户信息</p>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="8" :md="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon size="32" color="var(--color-primary)">
                <User />
              </el-icon>
              </div>
            <div class="stat-info">
              <h3 class="stat-value">{{ stats.total_users }}</h3>
              <p class="stat-label">总用户数</p>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="8" :md="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon size="32" color="var(--color-success)">
                <UserFilled />
              </el-icon>
        </div>
            <div class="stat-info">
              <h3 class="stat-value">{{ stats.active_users }}</h3>
              <p class="stat-label">活跃用户</p>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="8" :md="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon size="32" color="var(--color-warning)">
                <Clock />
              </el-icon>
        </div>
            <div class="stat-info">
              <h3 class="stat-value">{{ stats.new_users_today }}</h3>
              <p class="stat-label">今日新增</p>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="8" :md="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon size="32" color="var(--color-error)">
                <Lock />
              </el-icon>
            </div>
            <div class="stat-info">
              <h3 class="stat-value">{{ stats.banned_users }}</h3>
              <p class="stat-label">被封禁用户</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 操作栏 -->
    <el-card class="filter-card">
      <el-row :gutter="20" align="center">
        <el-col :xs="24" :sm="8" :md="6">
          <el-form-item label="用户名">
            <el-input
              v-model="filters.username"
              placeholder="输入用户名..."
              clearable
              @input="handleSearch"
            />
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="8" :md="6">
          <el-form-item label="邮箱">
            <el-input
              v-model="filters.email"
              placeholder="输入邮箱..."
              clearable
              @input="handleSearch"
            />
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="8" :md="6">
          <el-form-item label="状态">
            <el-select
              v-model="filters.status"
              placeholder="选择状态"
              clearable
              @change="handleSearch"
            >
              <el-option label="全部" value="" />
              <el-option label="正常" value="active" />
              <el-option label="封禁" value="banned" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="24" :md="6">
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 用户列表 -->
    <el-card class="table-card">
      <template #header>
        <h3 class="card-title">用户列表</h3>
      </template>
      
      <el-table
        v-loading="loading"
        :data="users"
        stripe
        style="width: 100%"
        empty-text="暂无用户数据"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" min-width="120">
          <template #default="{ row }">
            <strong>{{ row.username }}</strong>
          </template>
        </el-table-column>
        <el-table-column prop="email" label="邮箱" min-width="200">
          <template #default="{ row }">
            {{ row.email || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="role" label="角色" width="100">
          <template #default="{ row }">
            <el-tag
              :type="row.role === 'admin' ? 'danger' : 'primary'"
              size="small"
            >
              {{ row.role === 'admin' ? '管理员' : '普通用户' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag
              :type="row.status === 'active' ? 'success' : 'danger'"
              size="small"
            >
              {{ row.status === 'active' ? '正常' : '封禁' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="last_login" label="最后登录" width="180">
          <template #default="{ row }">
            {{ row.last_login ? formatDate(row.last_login) : '从未登录' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              @click="viewUserDetails(row.id)"
            >
              <el-icon><View /></el-icon>
              查看
            </el-button>
            <el-button
              v-if="row.status === 'active'"
              type="warning"
              size="small"
              @click="banUser(row.id, row.username)"
            >
              <el-icon><Lock /></el-icon>
              封禁
            </el-button>
            <el-button
              v-else
              type="success"
              size="small"
              @click="unbanUser(row.id, row.username)"
            >
              <el-icon><Unlock /></el-icon>
              解封
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="deleteUser(row.id, row.username)"
            >
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 用户详情对话框 -->
    <el-dialog
      v-model="userDetailVisible"
      title="用户详情"
      width="600px"
    >
      <div v-if="selectedUser" class="user-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="用户ID">
            {{ selectedUser.id }}
          </el-descriptions-item>
          <el-descriptions-item label="用户名">
            {{ selectedUser.username }}
          </el-descriptions-item>
          <el-descriptions-item label="邮箱">
            {{ selectedUser.email || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="角色">
            <el-tag
              :type="selectedUser.role === 'admin' ? 'danger' : 'primary'"
              size="small"
            >
              {{ selectedUser.role === 'admin' ? '管理员' : '普通用户' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag
              :type="selectedUser.status === 'active' ? 'success' : 'danger'"
              size="small"
            >
              {{ selectedUser.status === 'active' ? '正常' : '封禁' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="注册时间">
            {{ formatDate(selectedUser.created_at) }}
          </el-descriptions-item>
          <el-descriptions-item label="最后登录">
            {{ selectedUser.last_login ? formatDate(selectedUser.last_login) : '从未登录' }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>

    <!-- 删除确认对话框 -->
    <el-dialog
      v-model="deleteDialogVisible"
      title="确认删除"
      width="400px"
    >
      <p>确定要删除用户 "<strong>{{ deleteUserName }}</strong>" 吗？此操作无法撤销。</p>
      
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
const deleting = ref(false)
interface User {
  id: number
  username: string
  email?: string
  role: string
  status: string
  created_at: string
  last_login?: string
}

const users = ref<User[]>([])
const userDetailVisible = ref(false)
const deleteDialogVisible = ref(false)
const selectedUser = ref<User | null>(null)
const currentDeleteId = ref<number | null>(null)
const deleteUserName = ref('')

const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const stats = ref({
  total_users: 0,
  active_users: 0,
  new_users_today: 0,
  banned_users: 0
})

const filters = reactive({
  username: '',
  email: '',
  status: ''
})

const loadUsers = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (filters.username) params.append('username', filters.username)
    if (filters.email) params.append('email', filters.email)
    if (filters.status) params.append('status', filters.status)
    params.append('page', currentPage.value.toString())
    params.append('limit', pageSize.value.toString())
    
    const response = await api.get(`/admin/users?${params}`)
    
    if (response.data.data) {
      users.value = response.data.data.users || []
      total.value = response.data.data.total || 0
    } else {
      users.value = []
      total.value = 0
    }
  } catch (error: unknown) {
    console.error('加载用户失败:', error)
    ElMessage.error('加载用户数据失败，请检查网络连接')
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  try {
    const response = await api.get('/admin/stats')
    if (response.data.success) {
      stats.value = response.data.data
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

const handleSearch = () => {
  currentPage.value = 1
  loadUsers()
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  currentPage.value = 1
  loadUsers()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  loadUsers()
}

const viewUserDetails = async (id: number) => {
  try {
    const response = await api.get(`/admin/users/${id}`)
    if (response.data.data) {
      selectedUser.value = response.data.data
      userDetailVisible.value = true
    }
  } catch (error: unknown) {
    console.error('获取用户详情失败:', error)
    ElMessage.error('获取用户详情失败')
  }
}

const banUser = async (id: number, username: string) => {
  try {
    await ElMessageBox.confirm(
      `确定要封禁用户 "${username}" 吗？`,
      '确认封禁',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    const response = await api.put(`/admin/users/${id}/ban`)
    if (response.status === 200) {
      ElMessage.success('用户封禁成功')
      loadUsers()
      loadStats()
    }
  } catch (error: unknown) {
    if (error !== 'cancel') {
      console.error('封禁用户失败:', error)
      ElMessage.error('封禁用户失败')
    }
  }
}

const unbanUser = async (id: number, username: string) => {
  try {
    await ElMessageBox.confirm(
      `确定要解封用户 "${username}" 吗？`,
      '确认解封',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'success',
      }
    )
    
    const response = await api.put(`/admin/users/${id}/unban`)
    if (response.status === 200) {
      ElMessage.success('用户解封成功')
      loadUsers()
      loadStats()
    }
  } catch (error: unknown) {
    if (error !== 'cancel') {
      console.error('解封用户失败:', error)
      ElMessage.error('解封用户失败')
    }
  }
}

const deleteUser = (id: number, username: string) => {
  currentDeleteId.value = id
  deleteUserName.value = username
  deleteDialogVisible.value = true
}

const confirmDelete = async () => {
  if (!currentDeleteId.value) return

  try {
    deleting.value = true
    const response = await api.delete(`/admin/users/${currentDeleteId.value}`)

    if (response.status === 200) {
      deleteDialogVisible.value = false
      loadUsers()
      loadStats()
      ElMessage.success('用户删除成功')
    }
  } catch (error: unknown) {
    console.error('删除用户失败:', error)
    const errorMessage = (error as { response?: { data?: { error?: string } } })?.response?.data?.error || '未知错误'
    ElMessage.error('删除失败: ' + errorMessage)
  } finally {
    deleting.value = false
  }
}

const formatDate = (dateString: string) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('zh-CN')
}

onMounted(() => {
  loadUsers()
  loadStats()
})
</script>

<style scoped>
.user-management {
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

.stats-row {
  margin-bottom: 30px;
}

.stat-card {
  height: 100%;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 15px;
}

.stat-icon {
  width: 48px;
  height: 48px;
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
  font-size: 1.8rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0 0 5px 0;
}

.stat-label {
  font-size: 0.9rem;
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

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.user-detail {
  padding: 10px 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .user-management {
    padding: 15px;
  }
  
  .page-title {
    font-size: 1.5rem;
  }
  
  .stat-content {
    flex-direction: column;
    text-align: center;
    gap: 10px;
  }
  
  .stat-icon {
    width: 40px;
    height: 40px;
  }
  
  .stat-value {
    font-size: 1.5rem;
  }
}
</style>