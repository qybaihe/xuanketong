<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'

const isRegister = ref(false)
const username = ref('')
const password = ref('')
const email = ref('')
const nickname = ref('')

const handleSubmit = async () => {
  try {
    if (isRegister.value) {
      await axios.post('http://localhost:8080/api/v1/auth/register', {
        username: username.value,
        password: password.value,
        email: email.value,
        nickname: nickname.value
      })
      alert('注册成功！')
    } else {
      // Login logic will be implemented later
      alert('登录功能尚未实现。')
    }
  } catch (error) {
    console.error(error)
    alert('发生错误。')
  }
}
</script>

<template>
  <div>
    <h1>{{ isRegister ? '注册' : '登录' }}</h1>
    <form @submit.prevent="handleSubmit">
      <div>
        <label for="username">用户名:</label>
        <input type="text" id="username" v-model="username" required>
      </div>
      <div>
        <label for="password">密码:</label>
        <input type="password" id="password" v-model="password" required>
      </div>
      <div v-if="isRegister">
        <label for="email">电子邮箱:</label>
        <input type="email" id="email" v-model="email" required>
      </div>
      <div v-if="isRegister">
        <label for="nickname">昵称:</label>
        <input type="text" id="nickname" v-model="nickname" required>
      </div>
      <button type="submit">{{ isRegister ? '注册' : '登录' }}</button>
    </form>
    <button @click="isRegister = !isRegister">
      {{ isRegister ? '切换到登录' : '切换到注册' }}
    </button>
  </div>
</template>