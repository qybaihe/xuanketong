import api from './api'

// 课程笔记相关类型定义
export interface CreateNoteRequest {
  content: string
  isPublic: boolean
}

export interface UpdateNoteRequest {
  content: string
  isPublic: boolean
}

export interface CourseNote {
  id: number
  courseId: number
  content: string
  authorId: number
  authorName: string
  isPublic: boolean
  isOwner: boolean
  createdAt: string
  updatedAt: string
}

export interface NotesResponse {
  notes: CourseNote[]
}

export interface SaveNoteResponse {
  id: number
  courseId: number
  content: string
  lastModified: string
}

export interface ApiResponse<T> {
  success: boolean
  data: T
}

// 课程笔记API服务
export class CourseNoteApi {
  /**
   * 获取课程笔记列表
   */
  static async getCourseNotes(courseId: number): Promise<ApiResponse<NotesResponse>> {
    try {
      const response = await api.get(`/course-notes/course/${courseId}`)
      return response.data
    } catch (error) {
      console.error('获取课程笔记失败:', error)
      throw new Error('获取课程笔记失败')
    }
  }

  /**
   * 保存课程笔记
   */
  static async saveCourseNote(courseId: number, content: string, isPublic: boolean = false): Promise<ApiResponse<CourseNote>> {
    try {
      const requestData: CreateNoteRequest = {
        content,
        isPublic
      }
      const response = await api.post(`/course-notes/course/${courseId}`, requestData)
      return response.data
    } catch (error) {
      console.error('保存笔记失败:', error)
      throw new Error('保存笔记失败')
    }
  }

  /**
   * 更新课程笔记
   */
  static async updateCourseNote(noteId: number, content: string, isPublic: boolean): Promise<ApiResponse<CourseNote>> {
    try {
      const requestData: UpdateNoteRequest = {
        content,
        isPublic
      }
      const response = await api.put(`/notes/${noteId}`, requestData)
      return response.data
    } catch (error) {
      console.error('更新笔记失败:', error)
      throw new Error('更新笔记失败')
    }
  }

  /**
   * 删除课程笔记
   */
  static async deleteCourseNote(noteId: number): Promise<ApiResponse<null>> {
    try {
      const response = await api.delete(`/notes/${noteId}`)
      return response.data
    } catch (error) {
      console.error('删除笔记失败:', error)
      throw new Error('删除笔记失败')
    }
  }

  /**
   * 获取用户的笔记列表
   */
  static async getUserNotes(): Promise<ApiResponse<NotesResponse>> {
    try {
      const response = await api.get('/users/notes')
      return response.data
    } catch (error) {
      console.error('获取用户笔记失败:', error)
      throw new Error('获取用户笔记失败')
    }
  }
}

// 导出默认实例
export default CourseNoteApi
