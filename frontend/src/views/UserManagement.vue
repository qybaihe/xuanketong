<template>
  <div class="user-management">
    <div class="header-section mb-4">
      <div class="row align-items-center">
        <div class="col-md-6">
          <h2 class="mb-0">用户管理</h2>
          <p class="text-muted mb-0">管理系统中的所有用户账户</p>
        </div>
        <div class="col-md-6 text-md-end">
          <button class="btn btn-primary" @click="showAddUserModal = true">
            <i class="bi bi-plus-circle me-1"></i>添加用户
          </button>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="row mb-4">
      <div class="col-md-3">
        <div class="card stats-card">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-center">
              <div>
                <h5 class="card-title text-muted mb-1">总用户数</h5>
                <h3 class="mb-0">{{ userStats.total_users || 0 }}</h3>
              </div>
              <div class="stats-icon">
                <i class="bi bi-people-fill text-primary"></i>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card stats-card">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-center">
              <div>
                <h5 class="card-title text-muted mb-1">管理员</h5>
                <h3 class="mb-0">{{ userStats.total_admins || 0 }}</h3>
              </div>
              <div class="stats-icon">
                <i class="bi bi-person-check-fill text-success"></i>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card stats-card">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-center">
              <div>
                <h5 class="card-title text-muted mb-1">普通用户</h5>
                <h3 class="mb-0">{{ userStats.total_regular_users || 0 }}</h3>
              </div>
              <div class="stats-icon">
                <i class="bi bi-person-fill text-info"></i>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card stats-card">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-center">
              <div>
                <h5 class="card-title text-muted mb-1">活跃用户</h5>
                <h3 class="mb-0">{{ userStats.users_with_comment || 0 }}</h3>
              </div>
              <div class="stats-icon">
                <i class="bi bi-chat-dots-fill text-warning"></i>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <div class="card mb-4">
      <div class="card-body">
        <div class="row align-items-center">
          <div class="col-md-6">
            <div class="input-group">
              <span class="input-group-text"><i class="bi bi-search"></i></span>
              <input 
                type="text" 
                class="form-control" 
                placeholder="搜索用户名、昵称或邮箱..."
                v-model="searchKeyword"
                @input="handleSearch"
              >
            </div>
          </div>
          <div class="col-md-3">
            <select class="form-select" v-model="selectedRole" @change="handleSearch">
              <option value="">所有角色</option>
              <option value="user">普通用户</option>
              <option value="admin">管理员</option>
            </select>
          </div>
          <div class="col-md-3">
            <button class="btn btn-outline-secondary w-100" @click="resetSearch">
              <i class="bi bi-arrow-clockwise me-1"></i>重置
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 用户列表 -->
    <div class="card">
      <div class="card-header d-flex justify-content-between align-items-center">
        <h5 class="mb-0">用户列表</h5>
        <span class="text-muted">共 {{ totalUsers }} 条记录</span>
      </div>
      <div class="card-body p-0">
        <div class="table-responsive">
          <table class="table table-hover mb-0">
            <thead class="table-light">
              <tr>
                <th width="50">ID</th>
                <th>用户信息</th>
                <th>角色</th>
                <th>状态</th>
                <th>注册时间</th>
                <th width="150">操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="user in users" :key="user.id">
                <td>{{ user.id }}</td>
                <td>
                  <div class="d-flex align-items-center">
                    <img 
                      :src="user.avatar || '/default-avatar.png'" 
                      class="rounded-circle me-3" 
                      width="40" 
                      height="40"
                      :alt="user.nickname"
                    >
                    <div>
                      <h6 class="mb-0">{{ user.nickname || user.username }}</h6>
                      <small class="text-muted">{{ user.username }}</small>
                      <div class="text-muted small">{{ user.email }}</div>
                    </div>
                  </div>
                </td>
                <td>
                  <span :class="['badge', user.role === 'admin' ? 'bg-success' : 'bg-info']">
                    {{ user.role === 'admin' ? '管理员' : '普通用户' }}
                  </span>
                </td>
                <td>
                  <span class="badge bg-success">活跃</span>
                </td>
                <td>
                  <small class="text-muted">{{ formatDate(user.createdAt) }}</small>
                </td>
                <td>
                  <div class="btn-group btn-group-sm">
                    <button class="btn btn-outline-primary" @click="editUser(user)">
                      <i class="bi bi-pencil"></i>
                    </button>
                    <button class="btn btn-outline-danger" @click="deleteUser(user)" :disabled="user.role === 'admin'">
                      <i class="bi bi-trash"></i>
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="users.length === 0">
                <td colspan="6" class="text-center py-4 text-muted">
                  <i class="bi bi-inbox display-4"></i>
                  <p class="mt-2">暂无用户数据</p>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- 分页 -->
    <nav class="mt-4" v-if="totalPages > 1">
      <ul class="pagination justify-content-center">
        <li class="page-item" :class="{ disabled: currentPage === 1 }">
          <a class="page-link" href="#" @click.prevent="changePage(currentPage - 1)">上一页</a>
        </li>
        <li 
          v-for="page in visiblePages" 
          :key="page" 
          class="page-item" 
          :class="{ active: page === currentPage }"
        >
          <a class="page-link" href="#" @click.prevent="changePage(page)">{{ page }}</a>
        </li>
        <li class="page-item" :class="{ disabled: currentPage === totalPages }">
          <a class="page-link" href="#" @click.prevent="changePage(currentPage + 1)">下一页</a>
        </li>
      </ul>
    </nav>

    <!-- 编辑用户模态框 -->
    <div class="modal fade" :class="{ show: showEditModal }" tabindex="-1" v-if="showEditModal">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">编辑用户</h5>
            <button type="button" class="btn-close" @click="closeEditModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="updateUser">
              <div class="mb-3">
                <label for="nickname" class="form-label">昵称</label>
                <input 
                  type="text" 
                  class="form-control" 
                  id="nickname"
                  v-model="editingUser.nickname"
                  required
                >
              </div>
              <div class="mb-3">
                <label for="email" class="form-label">邮箱</label>
                <input 
                  type="email" 
                  class="form-control" 
                  id="email"
                  v-model="editingUser.email"
                  required
                >
              </div>
              <div class="mb-3">
                <label for="role" class="form-label">角色</label>
                <select class="form-select" id="role" v-model="editingUser.role">
                  <option value="user">普通用户</option>
                  <option value="admin">管理员</option>
                </select>
              </div>
              <div class="mb-3">
                <label for="avatar" class="form-label">头像URL</label>
                <input 
                  type="url" 
                  class="form-control" 
                  id="avatar"
                  v-model="editingUser.avatar"
                  placeholder="https://example.com/avatar.jpg"
                >
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="closeEditModal">取消</button>
            <button type="button" class="btn btn-primary" @click="updateUser">保存</button>
          </div>
        </div>
      </div>
    </div>

    <!-- 删除确认模态框 -->
    <div class="modal fade" :class="{ show: showDeleteModal }" tabindex="-1" v-if="showDeleteModal">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">确认删除</h5>
            <button type="button" class="btn-close" @click="closeDeleteModal"></button>
          </div>
          <div class="modal-body">
            <p>确定要删除用户 <strong>{{ deletingUser?.nickname || deletingUser?.username }}</strong> 吗？</p>
            <p class="text-danger">此操作将同时删除该用户的所有评分和评论，且无法恢复。</p>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="closeDeleteModal">取消</button>
            <button type="button" class="btn btn-danger" @click="confirmDelete">删除</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import type { User } from '@/stores/auth'

interface UserStats {
  total_users: number
  total_admins: number
  total_regular_users: number
  average_rating: number
  users_with_comment: number
}

const authStore = useAuthStore()

// 数据状态
const users = ref<User[]>([])
const userStats = ref<UserStats>({} as UserStats)
const totalUsers = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchKeyword = ref('')
const selectedRole = ref('')
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const showAddUserModal = ref(false)
const editingUser = ref<User>({} as User)
const deletingUser = ref<User | null>(null)

// 计算属性
const totalPages = computed(() => Math.ceil(totalUsers.value / pageSize.value))

const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, currentPage.value - 2)
  const end = Math.min(totalPages.value, currentPage.value + 2)
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

// 方法
const fetchUsers = async () => {
  try {
    const params = new URLSearchParams({
      page: currentPage.value.toString(),
      pageSize: pageSize.value.toString()
    })
    
    if (searchKeyword.value) {
      params.append('keyword', searchKeyword.value)
    }
    
    if (selectedRole.value) {
      params.append('role', selectedRole.value)
    }
    
    const response = await fetch(`/api/v1/admin/users?${params}`, {
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    })
    
    const data = await response.json()
    users.value = data.data
    totalUsers.value = data.total
  } catch (error) {
    console.error('获取用户列表失败:', error)
  }
}

const fetchUserStats = async () => {
  try {
    const response = await fetch('/api/v1/admin/user-stats', {
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    })
    
    userStats.value = await response.json()
  } catch (error) {
    console.error('获取用户统计失败:', error)
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchUsers()
}

const resetSearch = () => {
  searchKeyword.value = ''
  selectedRole.value = ''
  currentPage.value = 1
  fetchUsers()
}

const changePage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    fetchUsers()
  }
}

const editUser = (user: User) => {
  editingUser.value = { ...user }
  showEditModal.value = true
}

const closeEditModal = () => {
  showEditModal.value = false
  editingUser.value = {} as User
}

const updateUser = async () => {
  try {
    const response = await fetch(`/api/v1/admin/users/${editingUser.value.id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authStore.token}`
      },
      body: JSON.stringify({
        nickname: editingUser.value.nickname,
        email: editingUser.value.email,
        role: editingUser.value.role,
        avatar: editingUser.value.avatar
      })
    })
    
    if (response.ok) {
      closeEditModal()
      fetchUsers()
      fetchUserStats()
    } else {
      const error = await response.json()
      alert(error.error || '更新用户失败')
    }
  } catch (error) {
    console.error('更新用户失败:', error)
    alert('更新用户失败')
  }
}

const deleteUser = (user: User) => {
  deletingUser.value = user
  showDeleteModal.value = true
}

const closeDeleteModal = () => {
  showDeleteModal.value = false
  deletingUser.value = null
}

const confirmDelete = async () => {
  if (!deletingUser.value) return
  
  try {
    const response = await fetch(`/api/v1/admin/users/${deletingUser.value.id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    })
    
    if (response.ok) {
      closeDeleteModal()
      fetchUsers()
      fetchUserStats()
    } else {
      const error = await response.json()
      alert(error.error || '删除用户失败')
    }
  } catch (error) {
    console.error('删除用户失败:', error)
    alert('删除用户失败')
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('zh-CN')
}

// 生命周期
onMounted(() => {
  fetchUsers()
  fetchUserStats()
})
</script>

<style scoped>
.user-management {
  padding: 2rem;
}

.header-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 2rem;
  border-radius: 12px;
}

.stats-card {
  border: none;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s ease;
}

.stats-card:hover {
  transform: translateY(-2px);
}

.stats-icon {
  font-size: 2rem;
}

.modal {
  display: block;
  background-color: rgba(0, 0, 0, 0.5);
}

.modal.show {
  opacity: 1;
}

.table-hover tbody tr:hover {
  background-color: #f8f9fa;
}

.btn-group-sm .btn {
  padding: 0.25rem 0.5rem;
  font-size: 0.875rem;
}

.badge {
  padding: 0.5em 0.75em;
  font-weight: 500;
}

@media (max-width: 768px) {
  .user-management {
    padding: 1rem;
  }
  
  .header-section {
    padding: 1.5rem;
  }
  
  .stats-card {
    margin-bottom: 1rem;
  }
}
</style>