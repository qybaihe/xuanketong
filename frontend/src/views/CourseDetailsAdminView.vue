<template>
  <div>
    <div class="app-content-header">
      <div class="container-fluid">
        <div class="row">
          <div class="col-sm-6">
            <h3 class="mb-0">课程详情</h3>
          </div>
          <div class="col-sm-6">
            <ol class="breadcrumb float-sm-end">
              <li class="breadcrumb-item"><a href="#">主页</a></li>
              <li class="breadcrumb-item">
                <RouterLink to="/admin/courses">课程管理</RouterLink>
              </li>
              <li class="breadcrumb-item active" aria-current="page">
                课程详情
              </li>
            </ol>
          </div>
        </div>
      </div>
    </div>
    <div class="app-content">
      <div class="container-fluid">
        <div class="card mb-4" v-if="course">
          <div class="card-header">
            <h3 class="card-title">{{ course.name }}</h3>
          </div>
          <div class="card-body">
            <p><strong>教师:</strong> {{ course.teacher }}</p>
            <p><strong>描述:</strong> {{ course.description }}</p>
            <p><strong>年级:</strong> {{ course.grade }} | <strong>学期:</strong> {{ course.semester }} | <strong>科目:</strong> {{ course.subject }} | <strong>学分:</strong> {{ course.credits }}</p>
          </div>
        </div>
        <div class="row">
          <div class="col-md-6">
            <div class="card">
              <div class="card-header">
                <h3 class="card-title">评分</h3>
              </div>
              <div class="card-body">
                <table class="table table-bordered">
                  <thead>
                    <tr>
                      <th>用户ID</th>
                      <th>分数</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="rating in ratings" :key="rating.id">
                      <td>{{ rating.userID }}</td>
                      <td>{{ rating.score }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
          <div class="col-md-6">
            <div class="card">
              <div class="card-header">
                <h3 class="card-title">评论</h3>
              </div>
              <div class="card-body">
                <table class="table table-bordered">
                  <thead>
                    <tr>
                      <th>用户ID</th>
                      <th>内容</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="comment in comments" :key="comment.id">
                      <td>{{ comment.userID }}</td>
                      <td>{{ comment.content }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import axios from 'axios'

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

interface Rating {
  id: number;
  userID: number;
  score: number;
}

interface Comment {
  id: number;
  userID: number;
  content: string;
}

const route = useRoute()
const course = ref<Course | null>(null)
const ratings = ref<Rating[]>([])
const comments = ref<Comment[]>([])

const fetchDetails = async () => {
  const courseId = route.params.id
  try {
    const [courseRes, ratingsRes, commentsRes] = await Promise.all([
      axios.get(`${import.meta.env.VITE_BACKEND_BASE_URL}/courses/${courseId}`),
      axios.get(`${import.meta.env.VITE_BACKEND_BASE_URL}/courses/${courseId}/ratings`),
      axios.get(`${import.meta.env.VITE_BACKEND_BASE_URL}/courses/${courseId}/comments`)
    ])
    course.value = courseRes.data.data
    ratings.value = ratingsRes.data.data || []
    comments.value = commentsRes.data.data || []
  } catch (error) {
    console.error('获取课程详情失败:', error)
  }
}

onMounted(() => {
  fetchDetails()
})
</script>