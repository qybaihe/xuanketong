<script setup lang="ts">
import { RouterLink } from 'vue-router'
import type { EvaluationRequest } from '@/services/api'

interface Props {
  request: EvaluationRequest
}

const props = defineProps<Props>()

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const handleCardClick = () => {
  // 点击卡片跳转到课程详情页
  window.location.href = `#/courses/${props.request.course.id}`
}

const goToEvaluate = () => {
  // 点击"去评价"按钮跳转到课程评价页
  window.location.href = `#/courses/${props.request.course.id}/rate`
}
</script>

<template>
  <div class="evaluation-request-card" @click="handleCardClick">
    <!-- 装饰性边框 -->
    <div class="card-border"></div>
    
    <!-- 热门标记 -->
    <div v-if="request.status === 'pending'" class="status-badge status-pending">
      <span class="badge-icon">🔥</span>
      求评价中
    </div>
    
    <!-- 课程图片 -->
    <div class="course-image-container">
      <img
        :src="request.course.imageURL || `https://picsum.photos/seed/course-${request.course.id}/400/200.jpg`"
        :alt="request.course.name"
        class="course-image"
        @error="(e) => {
          const img = e.target as HTMLImageElement
          img.src = `https://picsum.photos/seed/course-${request.course.id}/400/200.jpg`
        }"
      />
      <div class="image-overlay">
        <div class="course-rating">
          <span class="rating-icon">⭐</span>
          <span class="rating-text">去评价</span>
        </div>
      </div>
    </div>
    
    <!-- 卡片内容 -->
    <div class="card-content">
      <!-- 课程信息 -->
      <div class="course-info">
        <h3 class="course-name">{{ request.course.name }}</h3>
        <div class="course-teacher">
          <span class="teacher-icon">👨‍🏫</span>
          <span class="teacher-name">{{ request.course.teacher }}</span>
        </div>
      </div>
      
      <!-- 发起人信息 -->
      <div class="requester-info">
        <div class="requester-avatar">
          {{ request.user.nickname.charAt(0) || 'U' }}
        </div>
        <div class="requester-details">
          <div class="requester-name">{{ request.user.nickname }}</div>
          <div class="request-time">发起于 {{ formatDate(request.createdAt) }}</div>
        </div>
      </div>
      
      <!-- 操作按钮 -->
      <div class="card-actions">
        <button class="btn btn-primary btn-evaluate" @click.stop="goToEvaluate">
          <span class="btn-icon">✍️</span>
          去评价
        </button>
        <RouterLink :to="`/courses/${request.course.id}`" class="btn btn-secondary" @click.stop>
          <span class="btn-icon">👀</span>
          查看详情
        </RouterLink>
      </div>
    </div>
    
    <!-- 卡片底部装饰 -->
    <div class="card-footer">
      <div class="footer-decoration"></div>
    </div>
  </div>
</template>

<style scoped>
.evaluation-request-card {
  position: relative;
  background: #ffffff;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  transition: all 0.3s cubic-bezier(0.25, 0.46, 0.45, 0.94);
  cursor: pointer;
  min-height: 380px;
  display: flex;
  flex-direction: column;
  border: 2px solid transparent;
}

.evaluation-request-card:hover {
  transform: translateY(-8px) scale(1.02);
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.18);
  border-color: var(--primary-color);
}

.card-border {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  border-radius: 20px;
  padding: 2px;
  background: linear-gradient(135deg, var(--primary-color) 0%, #1ebd8d 100%);
  -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
  mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
  mask-composite: exclude;
  -webkit-mask-composite: exclude;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.evaluation-request-card:hover .card-border {
  opacity: 1;
}

.status-badge {
  position: absolute;
  top: 16px;
  right: 16px;
  padding: 8px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 6px;
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  z-index: 2;
}

.status-pending {
  background: linear-gradient(135deg, #ff6b6b, #ff8e8e);
  color: white;
  box-shadow: 0 4px 12px rgba(255, 107, 107, 0.3);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% {
    transform: scale(1);
    box-shadow: 0 4px 12px rgba(255, 107, 107, 0.3);
  }
  50% {
    transform: scale(1.05);
    box-shadow: 0 6px 20px rgba(255, 107, 107, 0.5);
  }
  100% {
    transform: scale(1);
    box-shadow: 0 4px 12px rgba(255, 107, 107, 0.3);
  }
}

.badge-icon {
  font-size: 14px;
}

.course-image-container {
  position: relative;
  height: 160px;
  overflow: hidden;
  border-radius: 16px 16px 0 0;
}

.course-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.4s ease;
}

.evaluation-request-card:hover .course-image {
  transform: scale(1.05);
}

.image-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(0, 0, 0, 0.3) 0%, rgba(0, 0, 0, 0.5) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.evaluation-request-card:hover .image-overlay {
  opacity: 1;
}

.course-rating {
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  padding: 8px 16px;
  border-radius: 20px;
  color: var(--primary-color);
  font-weight: 600;
}

.rating-icon {
  font-size: 16px;
}

.rating-text {
  font-size: 14px;
}

.card-content {
  padding: 24px;
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.course-info {
  text-align: center;
}

.course-name {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 8px 0;
  line-height: 1.3;
}

.course-teacher {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 14px;
  color: var(--text-secondary);
}

.teacher-icon {
  font-size: 16px;
}

.teacher-name {
  font-weight: 500;
}

.requester-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: linear-gradient(135deg, rgba(47, 169, 20, 0.08) 0%, rgba(47, 169, 20, 0.04) 100%);
  border-radius: 12px;
  border: 1px solid rgba(47, 169, 20, 0.15);
}

.requester-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary-color), #1ebd8d);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 18px;
}

.requester-details {
  flex-grow: 1;
}

.requester-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.request-time {
  font-size: 12px;
  color: var(--text-tertiary);
}

.card-actions {
  display: flex;
  gap: 12px;
  margin-top: auto;
}

.btn-evaluate {
  flex-grow: 1;
  padding: 12px 16px;
  font-size: 14px;
  font-weight: 600;
}

.btn-secondary {
  padding: 12px 16px;
  background: rgba(47, 169, 20, 0.1);
  color: var(--primary-color);
  border: 1px solid rgba(47, 169, 20, 0.2);
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  text-decoration: none;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  transition: all 0.3s ease;
}

.btn-secondary:hover {
  background: rgba(47, 169, 20, 0.2);
  transform: translateY(-1px);
}

.card-footer {
  padding: 0 24px 24px;
}

.footer-decoration {
  height: 3px;
  background: linear-gradient(90deg, transparent, var(--primary-color), transparent);
  border-radius: 2px;
  opacity: 0.5;
}

.btn-icon {
  font-size: 16px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .card-actions {
    flex-direction: column;
  }
  
  .btn-evaluate,
  .btn-secondary {
    width: 100%;
  }
}

/* 暗色主题支持 */
@media (prefers-color-scheme: dark) {
  .evaluation-request-card {
    background: #1a1a1a;
    border-color: rgba(255, 255, 255, 0.1);
  }
  
  .course-name {
    color: white;
  }
  
  .status-pending {
    background: linear-gradient(135deg, #ff4757, #ff6b6b);
  }
  
  .requester-info {
    background: linear-gradient(135deg, rgba(47, 169, 20, 0.15) 0%, rgba(47, 169, 20, 0.08) 100%);
    border-color: rgba(47, 169, 20, 0.25);
  }
  
  .btn-secondary {
    background: rgba(47, 169, 20, 0.15);
    border-color: rgba(47, 169, 20, 0.3);
    color: #4fc830;
  }
  
  .footer-decoration {
    background: linear-gradient(90deg, transparent, #4fc830, transparent);
  }
}
</style>