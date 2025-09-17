<template>
  <div>
    <!-- 课程管理区域 -->
    <div id="courses-section">
      <!-- 操作栏 -->
      <div class="card mb-3">
        <div class="card-body">
          <div class="row g-3 align-items-end">
            <div class="col-md-3">
              <label class="form-label">年级</label>
              <input type="text" class="form-control" id="filter-grade" placeholder="输入年级...">
            </div>
            <div class="col-md-3">
              <label class="form-label">学期</label>
              <input type="text" class="form-control" id="filter-semester" placeholder="输入学期...">
            </div>
            <div class="col-md-3">
              <label class="form-label">科目</label>
              <input type="text" class="form-control" id="filter-subject" placeholder="输入科目...">
            </div>
            <div class="col-md-3 d-flex justify-content-end">
              <button class="btn btn-primary" @click="showAddCourseModal">
                <i class="bi bi-plus-circle me-1"></i>添加课程
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 加载状态 -->
      <div class="loading" id="loading">
        <div class="spinner-border text-primary" role="status">
          <span class="visually-hidden">加载中...</span>
        </div>
        <p class="mt-2">正在加载课程数据...</p>
      </div>

      <!-- 课程列表 -->
      <div class="card">
        <div class="card-header">
          <h5 class="mb-0">课程列表</h5>
        </div>
        <div class="card-body p-0">
          <div class="table-responsive">
            <table class="table table-hover mb-0">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>课程名称</th>
                  <th>教师</th>
                  <th>年级</th>
                  <th>学期</th>
                  <th>科目</th>
                  <th>学分</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody id="courses-tbody">
                <tr v-for="course in courses" :key="course.id">
                  <td>{{ course.id }}</td>
                  <td><strong>{{ course.name }}</strong></td>
                  <td>{{ course.teacher || '-' }}</td>
                  <td><span class="badge bg-primary badge-status">{{ course.grade || '-' }}</span></td>
                  <td><span class="badge bg-secondary badge-status">{{ course.semester || '-' }}</span></td>
                  <td><span class="badge bg-info badge-status">{{ course.subject || '-' }}</span></td>
                  <td><span class="badge bg-warning text-dark badge-status">{{ course.credits }} 学分</span></td>
                  <td>
                    <button class="btn btn-sm btn-outline-primary btn-action" @click="editCourse(course.id)" title="编辑">
                      <i class="bi bi-pencil"></i>
                    </button>
                    <button class="btn btn-sm btn-outline-danger btn-action" @click="deleteCourse(course.id, course.name)" title="删除">
                      <i class="bi bi-trash"></i>
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          
          <div class="empty-state" id="empty-state" style="display: none;">
            <i class="bi bi-book fs-2 text-muted"></i>
            <h5 class="mt-3">暂无课程数据</h5>
            <p class="text-muted">请调整筛选条件或添加新课程</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 添加/编辑课程模态框 -->
    <div class="modal fade" id="courseModal" tabindex="-1">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="courseModalTitle">添加课程</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
          </div>
          <div class="modal-body">
            <form id="courseForm">
              <input type="hidden" id="course-id">
              <div class="mb-3">
                <label class="form-label">课程名称 *</label>
                <input type="text" class="form-control" id="course-name" required>
              </div>
              <div class="mb-3">
                <label class="form-label">课程描述</label>
                <textarea class="form-control" id="course-description" rows="3"></textarea>
              </div>
              <div class="row">
                <div class="col-md-6 mb-3">
                  <label class="form-label">教师</label>
                  <input type="text" class="form-control" id="course-teacher">
                </div>
                <div class="col-md-6 mb-3">
                  <label class="form-label">学分</label>
                  <input type="number" class="form-control" id="course-credits" min="0" step="0.5">
                </div>
              </div>
              <div class="row">
                <div class="col-md-4 mb-3">
                  <label class="form-label">年级</label>
                  <input type="text" class="form-control" id="course-grade">
                </div>
                <div class="col-md-4 mb-3">
                  <label class="form-label">学期</label>
                  <input type="text" class="form-control" id="course-semester">
                </div>
                <div class="col-md-4 mb-3">
                  <label class="form-label">科目</label>
                  <input type="text" class="form-control" id="course-subject">
                </div>
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">取消</button>
            <button type="button" class="btn btn-primary" @click="saveCourse">保存</button>
          </div>
        </div>
      </div>
    </div>

    <!-- 删除确认模态框 -->
    <div class="modal fade" id="deleteModal" tabindex="-1">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">确认删除</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
          </div>
          <div class="modal-body">
            <p>确定要删除课程 "<span id="delete-course-name"></span>" 吗？此操作无法撤销。</p>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">取消</button>
            <button type="button" class="btn btn-danger" @click="confirmDelete">删除</button>
          </div>
        </div>
      </div>
    </div>

    <!-- 成功提示 -->
    <div class="position-fixed bottom-0 end-0 p-3" style="z-index: 11">
      <div id="successToast" class="toast" role="alert" aria-live="assertive" aria-atomic="true">
        <div class="toast-header">
          <i class="bi bi-check-circle-fill text-success me-2"></i>
          <strong class="me-auto">成功</strong>
          <button type="button" class="btn-close" data-bs-dismiss="toast"></button>
        </div>
        <div class="toast-body" id="toast-message">
          操作成功！
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import * as bootstrap from 'bootstrap';
import api from '@/services/api';
import { useAuthStore } from '@/stores/auth';

const API_BASE = 'http://localhost:8080/api/v1';
const courses = ref<any[]>([]);
const currentDeleteId = ref<number | null>(null);
const authStore = useAuthStore();
let courseModal: bootstrap.Modal | null = null;
let deleteModal: bootstrap.Modal | null = null;
let successToast: bootstrap.Toast | null = null;

onMounted(() => {
  courseModal = new bootstrap.Modal(document.getElementById('courseModal')!);
  deleteModal = new bootstrap.Modal(document.getElementById('deleteModal')!);
  successToast = new bootstrap.Toast(document.getElementById('successToast')!);
  
  loadCourses();
  bindFilterEvents();
});

function bindFilterEvents() {
  const filters = ['grade', 'semester', 'subject'];
  filters.forEach(filter => {
    const element = document.getElementById(`filter-${filter}`);
    if (element) {
      element.addEventListener('input', debounce(loadCourses, 300));
    }
  });
}

function debounce(func: Function, wait: number) {
  let timeout: number | undefined;
  return function executedFunction(...args: any[]) {
    const later = () => {
      clearTimeout(timeout);
      func(...args);
    };
    clearTimeout(timeout);
    timeout = window.setTimeout(later, wait);
  };
}

async function loadCourses() {
  try {
    const loadingEl = document.getElementById('loading');
    if(loadingEl) loadingEl.style.display = 'block';
    
    const grade = (document.getElementById('filter-grade') as HTMLInputElement).value;
    const semester = (document.getElementById('filter-semester') as HTMLInputElement).value;
    const subject = (document.getElementById('filter-subject') as HTMLInputElement).value;
    
    const params = new URLSearchParams();
    if (grade) params.append('grade', grade);
    if (semester) params.append('semester', semester);
    if (subject) params.append('subject', subject);
    
    const response = await api.get(`/admin/courses?${params}`);
    
    if (response.data.data) {
      courses.value = response.data.data;
    } else {
      courses.value = [];
    }
  } catch (error: any) {
    console.error('加载课程失败:', error);
    alert('加载课程数据失败，请检查网络连接');
  } finally {
    const loadingEl = document.getElementById('loading');
    if(loadingEl) loadingEl.style.display = 'none';
  }
}

function showAddCourseModal() {
  (document.getElementById('courseModalTitle') as HTMLElement).textContent = '添加课程';
  (document.getElementById('courseForm') as HTMLFormElement).reset();
  (document.getElementById('course-id') as HTMLInputElement).value = '';
  courseModal?.show();
}

async function editCourse(id: number) {
  try {
    const response = await api.get(`/admin/courses/${id}`);
    
    if (response.data.data) {
      const course = response.data.data;
      (document.getElementById('courseModalTitle') as HTMLElement).textContent = '编辑课程';
      (document.getElementById('course-id') as HTMLInputElement).value = course.id;
      (document.getElementById('course-name') as HTMLInputElement).value = course.name;
      (document.getElementById('course-description') as HTMLTextAreaElement).value = course.description || '';
      (document.getElementById('course-teacher') as HTMLInputElement).value = course.teacher || '';
      (document.getElementById('course-credits') as HTMLInputElement).value = course.credits;
      (document.getElementById('course-grade') as HTMLInputElement).value = course.grade || '';
      (document.getElementById('course-semester') as HTMLInputElement).value = course.semester || '';
      (document.getElementById('course-subject') as HTMLInputElement).value = course.subject || '';
      
      courseModal?.show();
    }
  } catch (error: any) {
    console.error('加载课程详情失败:', error);
    alert('加载课程详情失败');
  }
}

async function saveCourse() {
  const courseId = (document.getElementById('course-id') as HTMLInputElement).value;
  const courseData = {
    name: (document.getElementById('course-name') as HTMLInputElement).value,
    description: (document.getElementById('course-description') as HTMLTextAreaElement).value,
    teacher: (document.getElementById('course-teacher') as HTMLInputElement).value,
    credits: parseInt((document.getElementById('course-credits') as HTMLInputElement).value) || 0,
    grade: (document.getElementById('course-grade') as HTMLInputElement).value,
    semester: (document.getElementById('course-semester') as HTMLInputElement).value,
    subject: (document.getElementById('course-subject') as HTMLInputElement).value,
    imageURL: ''
  };

  try {
    let response;
    if (courseId) {
      response = await api.put(`/admin/courses/${courseId}`, courseData);
    } else {
      response = await api.post(`/admin/courses`, courseData);
    }

    if (response.status === 200 || response.status === 201) {
      courseModal?.hide();
      loadCourses();
      showToast(courseId ? '课程更新成功' : '课程添加成功');
    }
  } catch (error: any) {
    console.error('保存课程失败:', error);
    const errorMessage = error.response?.data?.error || '未知错误';
    alert('操作失败: ' + errorMessage);
  }
}

function deleteCourse(id: number, name: string) {
  currentDeleteId.value = id;
  (document.getElementById('delete-course-name') as HTMLElement).textContent = name;
  deleteModal?.show();
}

async function confirmDelete() {
  if (!currentDeleteId.value) return;

  try {
    const response = await api.delete(`/admin/courses/${currentDeleteId.value}`);

    if (response.status === 200) {
      deleteModal?.hide();
      loadCourses();
      showToast('课程删除成功');
    }
  } catch (error: any) {
    console.error('删除课程失败:', error);
    const errorMessage = error.response?.data?.error || '未知错误';
    alert('删除失败: ' + errorMessage);
  }
}

function showToast(message: string) {
  (document.getElementById('toast-message') as HTMLElement).textContent = message;
  successToast?.show();
}

</script>

<style scoped>

body {
  background-color: #ffffff;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  font-size: 14px;
  color: #333333;
}
.card {
  border: 1px solid #dee2e6;
  border-radius: 6px;
  box-shadow: none;
  transition: border-color 0.2s;
}
.card:hover {
  border-color: #999999;
}
.card-header {
  background: #f8f9fa;
  border-bottom: 1px solid #dee2e6;
  padding: 1rem 1.25rem;
}
.card-header h5 {
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
  color: #1a1a1a;
}
.card-body {
  padding: 1.25rem;
}
.btn {
  padding: 0.375rem 0.75rem;
  font-size: 0.875rem;
  border-radius: 4px;
  font-weight: 500;
}
.btn-primary {
  background-color: #007bff;
  border-color: #007bff;
}
.btn-primary:hover {
  background-color: #0056b3;
  border-color: #0056b3;
}
.btn-action {
  padding: 0.25rem 0.5rem;
  margin: 0 0.125rem;
  font-size: 0.8125rem;
}
.btn-action i {
  font-size: 0.875rem;
}
.form-control {
  font-size: 0.875rem;
  padding: 0.375rem 0.75rem;
  border: 1px solid #ced4da;
  border-radius: 4px;
}
.form-control:focus {
  border-color: #80bdff;
  box-shadow: 0 0 0 0.2rem rgba(0,123,255,.25);
}
.form-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #495057;
  margin-bottom: 0.5rem;
}
.table {
  font-size: 0.875rem;
  margin-bottom: 0;
}
.table th {
  background-color: #f8f9fa;
  border-bottom: 2px solid #dee2e6;
  font-weight: 600;
  color: #495057;
  padding: 0.75rem;
  white-space: nowrap;
}
.table td {
  padding: 0.75rem;
  vertical-align: middle;
  border-top: 1px solid #dee2e6;
}
.table tbody tr:hover {
  background-color: #f5f5f5;
}
.badge-status {
  font-size: 0.75rem;
  font-weight: 500;
}
.loading {
  display: none;
  text-align: center;
  padding: 2rem;
  background: #ffffff;
}
.empty-state {
  text-align: center;
  padding: 3rem;
  color: #6c757d;
  background: #ffffff;
}
.empty-state h5 {
  color: #495057;
  font-weight: 600;
}
.table-responsive {
  border: 1px solid #dee2e6;
  border-radius: 6px;
}
.modal-header {
  background: #f8f9fa;
  border-bottom: 1px solid #dee2e6;
}
.modal-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: #1a1a1a;
}
.modal-footer {
  background: #f8f9fa;
  border-top: 1px solid #dee2e6;
}
</style>