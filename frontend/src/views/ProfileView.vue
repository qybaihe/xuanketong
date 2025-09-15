<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'

interface User {
  ID: number;
  Username: string;
  Nickname: string;
  Email: string;
}

interface Rating {
  ID: number;
  CourseID: number;
  Score: number;
}

interface Comment {
  ID: number;
  CourseID: number;
  Content: string;
}

const route = useRoute()
const user = ref<User | null>(null)
const ratings = ref<Rating[]>([])
const comments = ref<Comment[]>([])

onMounted(async () => {
  try {
    // Assuming the user ID is 1 for now
    const userId = 1
    const response = await axios.get(`http://localhost:8080/api/v1/users/${userId}`)
    user.value = response.data.user
    ratings.value = response.data.ratings
    comments.value = response.data.comments
  } catch (error) {
    console.error(error)
  }
})
</script>

<template>
  <div v-if="user" class="card">
    <h1>{{ user.Nickname }}的个人资料</h1>
    <p><strong>用户名:</strong> {{ user.Username }}</p>
    <p><strong>电子邮箱:</strong> {{ user.Email }}</p>

    <hr>

    <div class="card">
      <h2>我的评分</h2>
      <ul>
        <li v-for="rating in ratings" :key="rating.ID">
          <p><strong>课程 {{ rating.CourseID }}:</strong> {{ rating.Score }}</p>
        </li>
      </ul>
    </div>

    <hr>

    <div class="card">
      <h2>我的评论</h2>
      <ul>
        <li v-for="comment in comments" :key="comment.ID">
          <p><strong>课程 {{ comment.CourseID }}:</strong></p>
          <p>{{ comment.Content }}</p>
        </li>
      </ul>
    </div>
  </div>
</template>