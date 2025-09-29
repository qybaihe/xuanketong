import axios from 'axios'

const API_BASE_URL = `${import.meta.env.VITE_BACKEND_BASE_URL}`

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      // 处理未授权情况
      localStorage.removeItem('token')
      window.location.href = '/auth'
    }
    return Promise.reject(error)
  }
)

export interface Course {
  id: number
  name: string
  description: string
  teacher: string
  credits: number
  grade: string
  semester: string
  subject: string
  createdAt: string
  updatedAt: string
  averageRating?: number
  totalRatings?: number
}

export interface Rating {
  ID: number
  UserID: number
  CourseID: number
  Score: number
  CreatedAt: string
  UpdatedAt: string
}

export interface Comment {
  ID: number
  UserID: number
  CourseID: number
  Content: string
  CreatedAt: string
  UpdatedAt: string
}

export interface EvaluationRequest {
  id: number
  userId: number
  courseId: number
  status: string
  createdAt: string
  updatedAt: string
  user: {
    id: number
    nickname: string
  }
  course: {
    id: number
    name: string
    teacher: string
  }
}

export interface EvaluationRequestListResponse {
  total: number
  page: number
  pageSize: number
  items: EvaluationRequest[]
}

export interface CreateEvaluationRequestRequest {
  courseId: number
}

export interface CourseWithRating extends Course {
  AverageRating?: number
  TotalRatings?: number
  ratingDistribution?: Record<number, number> // 1-5星评分分布
  // 兼容字段映射
  ID?: number
  Name?: string
  Description?: string
  Teacher?: string
  Credits?: number
  Grade?: string
  Semester?: string
  Subject?: string
  CreatedAt?: string
  UpdatedAt?: string
  RatingDistribution?: Record<number, number> // 1-5星评分分布（大写形式）
}

export interface CourseFilters {
  grade?: string
  semester?: string
  subject?: string
  name?: string
  teacher?: string
}

const courseService = {
  // 获取课程列表（包含评分信息）
  async getCourses(filters?: CourseFilters): Promise<CourseWithRating[]> {
    const response = await api.get('/courses', { params: filters })
    return response.data.data
  },

  // 获取单个课程
  async getCourse(id: number): Promise<CourseWithRating> {
    const response = await api.get(`/courses/${id}`)
    return response.data.data
  },

  // 获取课程的评分
  async getCourseRatings(courseId: number): Promise<Rating[]> {
    const response = await api.get(`/courses/${courseId}/ratings`)
    return response.data.data
  },

  // 获取课程的平均评分（现在后端API已经包含评分信息，这个方法保留用于兼容）
  async getCourseAverageRating(courseId: number): Promise<{ average: number; total: number }> {
    try {
      const course = await this.getCourse(courseId)
      return {
        average: course.AverageRating || 0,
        total: course.TotalRatings || 0
      }
    } catch (error) {
      console.error(`获取课程 ${courseId} 评分失败:`, error)
      return { average: 0, total: 0 }
    }
  },

  // 批量获取课程评分信息（现在后端API已经包含评分信息，这个方法保留用于兼容）
  async getCoursesWithRatings(courses: Course[]): Promise<CourseWithRating[]> {
    // 由于后端API已经返回了包含评分信息的课程数据，直接返回即可
    return courses as CourseWithRating[]
  }
}

const ratingService = {
  // 创建评分
  async createRating(userId: number, courseId: number, score: number): Promise<void> {
    await api.post('/ratings', {
      UserID: userId,
      CourseID: courseId,
      Score: score
    })
  },

  // 获取用户评分
  async getUserRating(userId: number, courseId: number): Promise<Rating | null> {
    try {
      const response = await api.get(`/courses/${courseId}/ratings`)
      const ratings = response.data.data as Rating[]
      const userRating = ratings.find(rating => rating.UserID === userId)
      return userRating || null
    } catch (error) {
      return null
    }
  }
}

const commentService = {
  // 创建评论
  async createComment(userId: number, courseId: number, content: string): Promise<void> {
    await api.post('/comments', {
      UserID: userId,
      CourseID: courseId,
      Content: content
    })
  },

  // 获取课程的评论
  async getCourseComments(courseId: number): Promise<Comment[]> {
    const response = await api.get(`/courses/${courseId}/comments`)
    return response.data.data
  },

  // 获取用户的评论
  async getUserComments(userId: number): Promise<Comment[]> {
    const response = await api.get(`/users/${userId}/comments`)
    return response.data.data
  }
}

const evaluationRequestService = {
  // 获取求评价列表
  async getEvaluationRequests(page: number = 1, pageSize: number = 10): Promise<EvaluationRequestListResponse> {
    const response = await api.get('/evaluation-requests', {
      params: { page, pageSize }
    })
    return response.data
  },

  // 创建新的求评价请求
  async createEvaluationRequest(courseId: number): Promise<EvaluationRequest> {
    const response = await api.post('/evaluation-requests', { courseId })
    return response.data
  }
}

export { courseService, ratingService, commentService, evaluationRequestService }
export default api