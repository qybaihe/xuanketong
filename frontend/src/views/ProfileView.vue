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
    <h1>{{ user.Nickname }}'s Profile</h1>
    <p><strong>Username:</strong> {{ user.Username }}</p>
    <p><strong>Email:</strong> {{ user.Email }}</p>

    <hr>

    <div class="card">
      <h2>My Ratings</h2>
      <ul>
        <li v-for="rating in ratings" :key="rating.ID">
          <p><strong>Course {{ rating.CourseID }}:</strong> {{ rating.Score }}</p>
        </li>
      </ul>
    </div>

    <hr>

    <div class="card">
      <h2>My Comments</h2>
      <ul>
        <li v-for="comment in comments" :key="comment.ID">
          <p><strong>Course {{ comment.CourseID }}:</strong></p>
          <p>{{ comment.Content }}</p>
        </li>
      </ul>
    </div>
  </div>
</template>