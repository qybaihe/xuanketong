import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import CourseDetailView from '../views/CourseDetailView.vue'
import ProfileView from '../views/ProfileView.vue'
import AuthView from '../views/AuthView.vue'
import AdminView from '../views/AdminView.vue'
import CourseManagement from '../views/CourseManagement.vue'
import CourseDetailsAdminView from '../views/CourseDetailsAdminView.vue'
import DashboardView from '../views/DashboardView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/courses/:id',
      name: 'course-detail',
      component: CourseDetailView
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfileView
    },
    {
      path: '/auth',
      name: 'auth',
      component: AuthView
    },
    {
      path: '/admin',
      component: AdminView,
      children: [
        {
          path: '',
          name: 'admin-dashboard',
          component: DashboardView
        },
        {
          path: 'dashboard',
          name: 'admin-dashboard-explicit',
          component: DashboardView
        },
        {
          path: 'courses',
          name: 'course-management',
          component: CourseManagement
        },
        {
          path: 'courses/:id',
          name: 'course-details-admin',
          component: CourseDetailsAdminView
        },
        {
          path: 'users',
          name: 'user-management',
          component: { template: '<div>User Management placeholder</div>' }
        },
        {
          path: 'comments',
          name: 'comment-management',
          component: { template: '<div>Comment Management placeholder</div>' }
        },
        {
          path: 'ratings',
          name: 'rating-management',
          component: { template: '<div>Rating Management placeholder</div>' }
        }
      ]
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue')
    }
  ]
})

export default router
