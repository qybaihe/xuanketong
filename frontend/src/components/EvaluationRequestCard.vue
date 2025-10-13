<script setup lang="ts">
import { useRouter } from 'vue-router'
import { RouterLink } from 'vue-router'
import type { EvaluationRequest } from '@/services/api'

const router = useRouter()

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
  router.push(`/courses/${props.request.course.id}`)
}

const goToEvaluate = (event: Event) => {
  // 阻止事件冒泡，避免触发卡片的点击事件
  event.stopPropagation()
  // 点击"去评价"按钮跳转到课程评价页
  router.push(`/courses/${props.request.course.id}/rate`)
}
</script>

<template>
  <div class="evaluation-request-card" @click="handleCardClick">
    <!-- 装饰性边框 -->
    <div class="card-border"></div>
    
    <!-- 热门标记 -->
    <!--
    <div v-if="request.status === 'pending'" class="status-badge status-pending">
      求评价中
    </div>
    -->
    
    <!-- 卡片内容 -->
    <div class="card-content">
      <!-- 课程信息 -->
      <div class="course-info">
        <h3 class="course-name">{{ request.course.name }}</h3>
        <div class="course-teacher">
          <svg class="teacher-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" stroke="var(--color-text)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <circle cx="12" cy="7" r="4" stroke="var(--color-text)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
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
        <button class="btn btn-primary btn-evaluate" @click="goToEvaluate($event)">
          <!--
          <svg class="btn-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 20h9" stroke="var(--color-text)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z" stroke="var(--color-text)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          -->
          去评价
        </button>
        <RouterLink :to="`/courses/${request.course.id}`" class="btn btn-secondary" @click.stop>
          <!--
          <svg class="btn-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="var(--color-text)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <circle cx="12" cy="12" r="3" stroke="var(--color-text)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          -->
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
  background-color: var(--color-surface);
  border-radius: 12px;
  overflow: hidden;
  transition: transform 0.2s ease;
  min-height: 280px;
  display: flex;
  flex-direction: column;
  border: 1px solid var(--color-border);
  cursor: pointer;
}

.evaluation-request-card:hover {
}

.card-border {
  display: none;
}

.status-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  padding: 6px 10px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: bold;
  display: flex;
  align-items: center;
  gap: 6px;
  z-index: 2;
  border: 1px solid var(--color-border);
}

.status-pending {
  background-color: var(--color-secondary);
  color: var(--color-text);
}


.badge-icon {
  font-size: 14px;
}


.card-content {
  padding: 20px;
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.course-info {
  text-align: center;
}

.course-name {
  font-size: 18px;
  font-weight: bold;
  color: var(--color-text);
  margin: 0 0 8px 0;
  line-height: 1.4;
}

.course-teacher {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 14px;
  color: var(--color-text);
}

.teacher-icon {
  width: 16px;
  height: 16px;
}

.teacher-name {
  font-weight: bold;
}

.requester-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background-color: var(--color-primary);
  border-radius: 8px;
  border: 1px solid var(--color-border);
}

.requester-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: var(--color-accent);
  color: var(--color-text);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 16px;
  border: 1px solid var(--color-border);
}

.requester-details {
  flex-grow: 1;
}

.requester-name {
  font-size: 14px;
  font-weight: bold;
  color: var(--color-text);
  margin-bottom: 4px;
}

.request-time {
  font-size: 12px;
  color: var(--color-text);
  font-weight: 500;
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
  font-weight: bold;
  background-color: var(--color-primary);
  color: var(--color-text);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  box-shadow: none;
  transition: transform 0.2s ease;
}

.btn-evaluate:hover,
.btn-evaluate:active,
.btn-evaluate:focus {
  box-shadow: none;
  outline: none;
  transform: none;
}

.btn-secondary {
  padding: 12px 16px;
  background-color: var(--color-accent);
  color: var(--color-text);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  font-size: 14px;
  font-weight: bold;
  text-decoration: none;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: transform 0.2s ease;
  box-shadow: none;
}

.btn-secondary:hover,
.btn-secondary:active,
.btn-secondary:focus {
  box-shadow: none;
  outline: none;
  transform: none;
}

.card-footer {
  display: none;
}

.btn-icon {
  width: 16px;
  height: 16px;
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
    background: var(--color-text);
    border-color: rgba(255, 255, 255, 0.1);
  }
  
  .course-name {
    color: white;
  }
  
  .status-pending {
    background: linear-gradient(135deg, var(--color-secondary), var(--color-secondary));
  }
  
  .requester-info {
    background: linear-gradient(135deg, rgba(47, 169, 20, 0.15) 0%, rgba(47, 169, 20, 0.08) 100%);
    border-color: rgba(47, 169, 20, 0.25);
  }
  
  .btn-secondary {
    background: rgba(47, 169, 20, 0.15);
    border-color: rgba(47, 169, 20, 0.3);
    color: var(--color-success);
  }
  
  .footer-decoration {
    background: linear-gradient(90deg, transparent, var(--color-success), transparent);
  }
}
</style>