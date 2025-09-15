<template>
  <div class="app-wrapper">
    <Topbar />
    <Sidenav />
    <main class="app-main">
      <router-view></router-view>
    </main>
    <Footer />
  </div>
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router'
import Sidenav from '@/components/admin/Sidenav.vue'
import Topbar from '@/components/admin/Topbar.vue'
import Footer from '@/components/admin/Footer.vue'
</script>

<style>
@import '@/assets/adminlte.css';

body {
  --bs-primary: #2fa914;
  --bs-secondary: #7d9871;
  --bs-light: #f1fdeb;
  --bs-dark: #9999cc;
}
</style>