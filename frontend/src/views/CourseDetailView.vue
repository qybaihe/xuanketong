<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'

interface Course {
  ID: number;
  Name: string;
  Description: string;
  Grade: string;
  Semester: string;
  Subject: string;
  Teacher: string;
  Credits: number;
  ImageURL: string;
}

interface Rating {
  ID: number;
  UserID: number;
  CourseID: number;
  Score: number;
}

interface Comment {
  ID: number;
  UserID: number;
  CourseID: number;
  Content: string;
}

const route = useRoute()
const course = ref<Course | null>(null)
const ratings = ref<Rating[]>([])
const comments = ref<Comment[]>([])
const newScore = ref(5)
const newComment = ref('')

const submitRating = async () => {
  try {
    const courseId = route.params.id
    await axios.post('http://localhost:8080/api/v1/ratings', {
      courseId: Number(courseId),
      userId: 1, // Hardcoded for now
      score: newScore.value
    })
    newScore.value = 5
    // Refresh ratings
    const response = await axios.get(`http://localhost:8080/api/v1/courses/${courseId}/ratings`)
    ratings.value = response.data.data
  } catch (error) {
    console.error(error)
  }
}

const submitComment = async () => {
  try {
    const courseId = route.params.id
    await axios.post('http://localhost:8080/api/v1/comments', {
      courseId: Number(courseId),
      userId: 1, // Hardcoded for now
      content: newComment.value
    })
    newComment.value = ''
    // Refresh comments
    const response = await axios.get(`http://localhost:8080/api/v1/courses/${courseId}/comments`)
    comments.value = response.data.data
  } catch (error) {
    console.error(error)
  }
}

onMounted(async () => {
  try {
    const courseId = route.params.id
    const [courseResponse, ratingsResponse, commentsResponse] = await Promise.all([
      axios.get(`http://localhost:8080/api/v1/courses/${courseId}`),
      axios.get(`http://localhost:8080/api/v1/courses/${courseId}/ratings`),
      axios.get(`http://localhost:8080/api/v1/courses/${courseId}/comments`)
    ])
    course.value = courseResponse.data.data
    ratings.value = ratingsResponse.data.data
    comments.value = commentsResponse.data.data
  } catch (error) {
    console.error(error)
  }
})
</script>

<template>
  <div v-if="course" class="card">
    <h1>{{ course.Name }}</h1>
    <img :src="course.ImageURL" :alt="course.Name" />
    <p>{{ course.Description }}</p>
    <ul>
      <li><strong>年级:</strong> {{ course.Grade }}</li>
      <li><strong>学期:</strong> {{ course.Semester }}</li>
      <li><strong>科目:</strong> {{ course.Subject }}</li>
      <li><strong>教师:</strong> {{ course.Teacher }}</li>
      <li><strong>学分:</strong> {{ course.Credits }}</li>
    </ul>

    <hr>

    <div class="card">
      <h2>评分</h2>
      <ul>
        <li v-for="rating in ratings" :key="rating.ID">
          用户 {{ rating.UserID }}: {{ rating.Score }}
        </li>
      </ul>
      <form @submit.prevent="submitRating">
        <label for="score">分数:</label>
        <input type="number" id="score" v-model.number="newScore" min="1" max="5" required>
        <button type="submit">提交评分</button>
      </form>
    </div>

    <hr>

    <div class="card">
      <h2>评论</h2>
      <ul>
        <li v-for="comment in comments" :key="comment.ID">
          <p><strong>用户 {{ comment.UserID }}:</strong></p>
          <p>{{ comment.Content }}</p>
        </li>
      </ul>
      <form @submit.prevent="submitComment">
        <label for="comment">评论:</label>
        <textarea id="comment" v-model="newComment" required></textarea>
        <button type="submit">提交评论</button>
      </form>
    </div>
  </div>
  <div v-else>
    <p>正在加载课程详情...</p>
  </div>
</template>

<style scoped>
img {
  max-width: 100%;
  height: auto;
}
</style>