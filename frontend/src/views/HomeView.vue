<script setup lang="ts">
import { ref, onMounted, reactive, watch } from 'vue'
import axios from 'axios'
import { RouterLink } from 'vue-router'

interface Course {
  ID: number;
  Name: string;
  Description: string;
  Teacher: string;
  Credits: number;
  ImageURL: string;
  Grade: string;
  Semester: string;
  Subject: string;
}

const courses = ref<Course[]>([])
const filters = reactive({
  grade: '',
  semester: '',
  subject: ''
})

const fetchCourses = async () => {
  try {
    const response = await axios.get('http://localhost:8080/api/v1/courses', { params: filters })
    courses.value = response.data.data
  } catch (error) {
    console.error(error)
  }
}

watch(filters, fetchCourses)

onMounted(() => {
  fetchCourses()
})
</script>

<template>
  <main class="card">
    <h1>Courses</h1>
    <form @submit.prevent="fetchCourses" class="filter-form">
      <div class="form-group">
        <label for="grade">Grade</label>
        <input type="text" id="grade" v-model="filters.grade" />
      </div>
      <div class="form-group">
        <label for="semester">Semester</label>
        <input type="text" id="semester" v-model="filters.semester" />
      </div>
      <div class="form-group">
        <label for="subject">Subject</label>
        <input type="text" id="subject" v-model="filters.subject" />
      </div>
    </form>
    <div class="course-grid">
      <router-link v-for="course in courses" :key="course.ID" :to="`/courses/${course.ID}`" class="course-card">
        <img :src="course.ImageURL || 'https://via.placeholder.com/150'" alt="Course Image">
        <div class="course-info">
          <h2>{{ course.Name }}</h2>
          <p>{{ course.Description }}</p>
          <p><strong>Teacher:</strong> {{ course.Teacher }}</p>
          <p><strong>Credits:</strong> {{ course.Credits }}</p>
        </div>
      </router-link>
    </div>
  </main>
</template>

<style scoped>
.filter-form {
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
}
.form-group {
  display: flex;
  flex-direction: column;
}
.course-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1rem;
}

.course-card {
  border: 1px solid #ccc;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.course-card img {
  width: 100%;
  height: 150px;
  object-fit: cover;
}

.course-info {
  padding: 1rem;
}
</style>
