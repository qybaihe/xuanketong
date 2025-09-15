<template>
  <div>
    <div class="app-content-header">
      <div class="container-fluid">
        <div class="row">
          <div class="col-sm-6">
            <h3 class="mb-0">课程管理</h3>
          </div>
          <div class="col-sm-6">
            <ol class="breadcrumb float-sm-end">
              <li class="breadcrumb-item"><a href="#">主页</a></li>
              <li class="breadcrumb-item active" aria-current="page">
                课程管理
              </li>
            </ol>
          </div>
        </div>
      </div>
    </div>
    <div class="app-content">
      <div class="container-fluid">
        <div class="card">
          <div class="card-header">
            <h3 class="card-title">课程</h3>
            <div class="card-tools d-flex align-items-center">
              <input type="text" class="form-control me-2" placeholder="年级" v-model="filter.grade" @input="fetchCourses" style="width: 150px;">
              <input type="text" class="form-control me-2" placeholder="学期" v-model="filter.semester" @input="fetchCourses" style="width: 150px;">
              <input type="text" class="form-control me-2" placeholder="科目" v-model="filter.subject" @input="fetchCourses" style="width: 150px;">
              <button class="btn btn-primary" @click="showAddModal = true">添加课程</button>
            </div>
          </div>
          <div class="card-body">
            <table class="table table-striped table-hover">
              <thead>
                <tr>
                  <th style="width: 10px">#</th>
                  <th>名称</th>
                  <th>教师</th>
                  <th>年级</th>
                  <th>学期</th>
                  <th>科目</th>
                  <th style="width: 180px">操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="course in courses" :key="course.id">
                  <td>{{ course.id }}</td>
                  <td>{{ course.name }}</td>
                  <td>{{ course.teacher }}</td>
                  <td>{{ course.grade }}</td>
                  <td>{{ course.semester }}</td>
                  <td>{{ course.subject }}</td>
                  <td>
                    <button class="btn btn-sm btn-info me-2" @click="viewCourse(course.id)">查看</button>
                    <button class="btn btn-sm btn-warning me-2" @click="editCourse(course)">编辑</button>
                    <button class="btn btn-sm btn-danger" @click="confirmDelete(course.id)">删除</button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <!-- Add/Edit Modal -->
    <div class="modal" :class="{ 'd-block': showAddModal }" style="background-color: rgba(0,0,0,0.5);">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ editingCourse ? '编辑课程' : '添加课程' }}</h5>
            <button type="button" class="btn-close" @click="closeModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveCourse">
              <div class="mb-3">
                <label for="name" class="form-label">名称</label>
                <input type="text" class="form-control" id="name" v-model="courseForm.name" required>
              </div>
              <div class="mb-3">
                <label for="description" class="form-label">描述</label>
                <textarea class="form-control" id="description" v-model="courseForm.description"></textarea>
              </div>
              <div class="mb-3">
                <label for="teacher" class="form-label">教师</label>
                <input type="text" class="form-control" id="teacher" v-model="courseForm.teacher">
              </div>
              <div class="mb-3">
                <label for="grade" class="form-label">年级</label>
                <input type="text" class="form-control" id="grade" v-model="courseForm.grade">
              </div>
              <div class="mb-3">
                <label for="semester" class="form-label">学期</label>
                <input type="text" class="form-control" id="semester" v-model="courseForm.semester">
              </div>
              <div class="mb-3">
                <label for="subject" class="form-label">科目</label>
                <input type="text" class="form-control" id="subject" v-model="courseForm.subject">
              </div>
               <div class="mb-3">
                <label for="credits" class="form-label">学分</label>
                <input type="number" class="form-control" id="credits" v-model.number="courseForm.credits">
              </div>
              <button type="submit" class="btn btn-primary">保存</button>
            </form>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

interface Course {
  id: number;
  name: string;
  description: string;
  teacher: string;
  grade: string;
  semester: string;
  subject: string;
  credits: number;
}

const courses = ref<Course[]>([])
const showAddModal = ref(false)
const editingCourse = ref<Course | null>(null)
const router = useRouter()

const filter = reactive({
  grade: '',
  semester: '',
  subject: ''
})

const courseForm = reactive({
  id: 0,
  name: '',
  description: '',
  teacher: '',
  grade: '',
  semester: '',
  subject: '',
  credits: 0
})

const fetchCourses = async () => {
  try {
    const params = new URLSearchParams()
    if (filter.grade) params.append('grade', filter.grade)
    if (filter.semester) params.append('semester', filter.semester)
    if (filter.subject) params.append('subject', filter.subject)
    const response = await axios.get('http://localhost:8080/api/v1/courses', { params })
    courses.value = response.data.data
  } catch (error) {
    console.error('获取课程失败:', error)
  }
}

const closeModal = () => {
  showAddModal.value = false
  editingCourse.value = null
  Object.assign(courseForm, { id:0, name: '', description: '', teacher: '', grade: '', semester: '', subject: '', credits: 0 })
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
    if (editingCourse.value) {
      await axios.put(`http://localhost:8080/api/v1/admin/courses/${courseForm.id}`, courseForm)
    } else {
      await axios.post('http://localhost:8080/api/v1/admin/courses', courseForm)
    }
    closeModal()
    fetchCourses()
  } catch (error) {
    console.error('保存课程失败:', error)
  }
}

const confirmDelete = (id: number) => {
  if (window.confirm('您确定要删除此课程吗？')) {
    deleteCourse(id)
  }
}

const deleteCourse = async (id: number) => {
  try {
    await axios.delete(`http://localhost:8080/api/v1/admin/courses/${id}`)
    fetchCourses()
  } catch (error) {
    console.error('删除课程失败:', error)
  }
}

onMounted(() => {
  fetchCourses()
})
</script>

<style>
.modal.d-block {
  display: block;
}
.card-tools {
  gap: 0.5rem;
}
</style>