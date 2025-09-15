<template>
  <div class="content-wrapper">
    <div class="content-header">
      <div class="container-fluid">
        <div class="row mb-2">
          <div class="col-sm-6">
            <h1 class="m-0">仪表盘</h1>
          </div>
        </div>
      </div>
    </div>

    <section class="content">
      <div class="container-fluid">
        <div class="row">
          <div class="col-lg-3 col-6">
            <div class="small-box bg-info">
              <div class="inner">
                <h3>{{ stats.total_courses }}</h3>
                <p>总课程数</p>
              </div>
              <div class="icon">
                <i class="ion ion-bag"></i>
              </div>
            </div>
          </div>
          <div class="col-lg-3 col-6">
            <div class="small-box bg-success">
              <div class="inner">
                <h3>{{ stats.average_rating.toFixed(2) }}<sup style="font-size: 20px">/5</sup></h3>
                <p>平均评分</p>
              </div>
              <div class="icon">
                <i class="ion ion-stats-bars"></i>
              </div>
            </div>
          </div>
          <div class="col-lg-3 col-6">
            <div class="small-box bg-warning">
              <div class="inner">
                <h3>{{ stats.total_users }}</h3>
                <p>注册用户</p>
              </div>
              <div class="icon">
                <i class="ion ion-person-add"></i>
              </div>
            </div>
          </div>
          <div class="col-lg-3 col-6">
            <div class="small-box bg-danger">
              <div class="inner">
                <h3>{{ stats.total_comments }}</h3>
                <p>总评论数</p>
              </div>
              <div class="icon">
                <i class="ion ion-pie-graph"></i>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import axios from 'axios';

const stats = ref({
  total_users: 0,
  total_courses: 0,
  total_comments: 0,
  average_rating: 0,
});

const fetchStats = async () => {
  try {
    const response = await axios.get('/api/v1/admin/stats');
    if (response.data) {
      stats.value = response.data;
    }
  } catch (error) {
    console.error('获取统计数据失败:', error);
  }
};

onMounted(() => {
  fetchStats();
});
</script>

<style scoped>
/* 如果需要，可以在这里添加特定于此组件的样式 */
.small-box .icon {
    color: rgba(0,0,0,.15);
    z-index: 0;
}

.small-box .icon > i {
    font-size: 90px;
    position: absolute;
    right: 15px;
    top: 15px;
    transition: all .3s linear;
}
</style>