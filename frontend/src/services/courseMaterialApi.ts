import api from './api'

// 课程资料相关类型定义
export interface MaterialData {
  materialName: string
  type: 'cloud' | 'note'
  service: string
  link: string
  accessPassword?: string
  description?: string
}

export interface CourseMaterial {
  id: number
  courseId: number
  courseName?: string
  materialName: string
  type: 'cloud' | 'note'
  service: string
  link: string
  accessPassword?: string
  description?: string
  uploaderId: number
  uploaderName: string
  uploadTime: string
  downloadCount: number
}

export interface MaterialsResponse {
  total: number
  materials: CourseMaterial[]
}

export interface ApiResponse<T> {
  success: boolean
  data: T
}

// 课程资料API服务
export class CourseMaterialApi {
  /**
   * 获取课程资料列表
   */
  static async getCourseMaterials(courseId: number): Promise<ApiResponse<MaterialsResponse>> {
    try {
      const response = await api.get(`/course-materials/course/${courseId}`)
      return response.data
    } catch (error) {
      console.error('获取课程资料失败:', error)
      throw new Error('获取课程资料失败')
    }
  }

  /**
   * 获取单个课程资料详情
   */
  static async getCourseMaterial(materialId: number): Promise<ApiResponse<CourseMaterial>> {
    try {
      const response = await api.get(`/materials/${materialId}`)
      return response.data
    } catch (error) {
      console.error('获取资料详情失败:', error)
      throw new Error('获取资料详情失败')
    }
  }

  /**
   * 上传课程资料
   */
  static async uploadMaterial(courseId: number, materialData: MaterialData): Promise<ApiResponse<CourseMaterial>> {
    try {
      const response = await api.post(`/course-materials/course/${courseId}`, materialData)
      return response.data
    } catch (error) {
      console.error('上传资料失败:', error)
      throw new Error('上传资料失败')
    }
  }

  /**
   * 更新课程资料
   */
  static async updateMaterial(materialId: number, materialData: MaterialData): Promise<ApiResponse<CourseMaterial>> {
    try {
      const response = await api.put(`/materials/${materialId}`, materialData)
      return response.data
    } catch (error) {
      console.error('更新资料失败:', error)
      throw new Error('更新资料失败')
    }
  }

  /**
   * 删除课程资料
   */
  static async deleteMaterial(materialId: number): Promise<ApiResponse<null>> {
    try {
      const response = await api.delete(`/materials/${materialId}`)
      return response.data
    } catch (error) {
      console.error('删除资料失败:', error)
      throw new Error('删除资料失败')
    }
  }

  /**
   * 增加下载次数
   */
  static async incrementDownloadCount(materialId: number): Promise<ApiResponse<{ downloadCount: number }>> {
    try {
      const response = await api.post(`/materials/${materialId}/download`)
      return response.data
    } catch (error) {
      console.error('更新下载次数失败:', error)
      throw new Error('更新下载次数失败')
    }
  }

  /**
   * 获取热门资料
   */
  static async getPopularMaterials(limit: number = 10): Promise<ApiResponse<CourseMaterial[]>> {
    try {
      const response = await api.get(`/course-materials/popular?limit=${limit}`)
      return response.data
    } catch (error) {
      console.error('获取热门资料失败:', error)
      throw new Error('获取热门资料失败')
    }
  }
}

// 导出默认实例
export default CourseMaterialApi