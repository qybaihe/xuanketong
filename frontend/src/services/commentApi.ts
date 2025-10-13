import api from './api'

export interface Comment {
  id: number
  userId: number
  courseId: number
  content: string
  createdAt: string
  updatedAt: string
  user: {
    id: number
    nickname: string
    avatar?: string
  }
  course?: {
    id: number
    name: string
    teacher: string
  }
  courseName?: string
  courseTeacher?: string
}

export interface CommentWithCourse extends Comment {
  courseName: string
  courseTeacher: string
}

export interface CreateCommentRequest {
  courseId: number
  content: string
}

export interface CommentListResponse {
  data: Comment[]
}

export interface LatestCommentsResponse {
  data: CommentWithCourse[]
}

const commentApi = {
  // 获取最新评论（用于首页展示）
  async getLatestComments(limit: number = 8): Promise<LatestCommentsResponse> {
    const response = await api.get(`/comments/latest?limit=${limit}`)
    return response.data
  },

  // 获取指定课程的评论
  async getCourseComments(courseId: number): Promise<CommentListResponse> {
    const response = await api.get(`/courses/${courseId}/comments`)
    return response.data
  },

  // 创建评论
  async createComment(data: CreateCommentRequest): Promise<void> {
    await api.post('/comments', data)
  },

  // 获取用户的所有评论
  async getUserComments(userId: number): Promise<CommentListResponse> {
    const response = await api.get(`/users/${userId}/comments`)
    return response.data
  }
}

export default commentApi
