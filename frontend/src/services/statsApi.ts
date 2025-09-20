import axios from 'axios'
import type { User } from '@/stores/auth'

// 首页统计数据接口定义
export interface OverviewStats {
  total_courses: number
  total_users: number
  total_ratings: number
  total_comments: number
  average_rating: number
  active_users_this_week: number
}

export interface TopRatedCourse {
  id: number
  name: string
  teacher: string
  average_rating: number
  total_ratings: number
  image_url: string
  subject: string
  grade: string
}

export interface RecentCourse {
  id: number
  name: string
  teacher: string
  description: string
  image_url: string
  subject: string
  grade: string
  created_at: string
}

export interface PopularCourse {
  id: number
  name: string
  teacher: string
  average_rating: number
  total_ratings: number
  image_url: string
  subject: string
  grade: string
  student_engagement: number
}

export interface CourseDistribution {
  by_subject: Record<string, number>
  by_grade: Record<string, number>
  by_semester: Record<string, number>
}

export interface MonthlyStat {
  month: string
  courses: number
  users: number
  ratings: number
  comments: number
  total_score: number
}

export interface UserActivityStats {
  active_users_today: number
  active_users_this_week: number
  active_users_this_month: number
  new_users_today: number
  new_users_this_week: number
  new_users_this_month: number
}

export interface EnhancedHomeStats {
  overview_data: OverviewStats
  top_rated_courses: TopRatedCourse[]
  recent_courses: RecentCourse[]
  popular_courses: PopularCourse[]
  user_activity_stats: UserActivityStats
  course_distribution: CourseDistribution
  monthly_stats: MonthlyStat[]
}

// API服务类
export class StatsApi {
  private baseURL = `${import.meta.env.VITE_BACKEND_BASE_URL}/`

  // 获取增强版首页统计数据
  async getEnhancedHomeStats(): Promise<EnhancedHomeStats> {
    try {
      const response = await axios.get(`${this.baseURL}/admin/stats/enhanced`)
      return response.data
    } catch (error) {
      console.error('获取首页统计数据失败:', error)
      throw new Error('获取首页统计数据失败')
    }
  }

  // 获取总体统计数据
  async getOverviewStats(): Promise<OverviewStats> {
    try {
      const response = await axios.get(`${this.baseURL}/admin/stats`)
      const stats = response.data as any
      // 适配新的统计结构
      return {
        total_courses: stats.total_courses || 0,
        total_users: stats.total_users || 0,
        total_ratings: stats.total_ratings || 0,
        total_comments: stats.total_comments || 0,
        average_rating: stats.average_rating || 0,
        active_users_this_week: 0 // 从basic stats中获取不到这个数据
      }
    } catch (error) {
      console.error('获取总体统计数据失败:', error)
      throw new Error('获取总体统计数据失败')
    }
  }

  // 获取用户活动统计
  async getUserActivityStats(): Promise<UserActivityStats> {
    try {
      const enhancedStats = await this.getEnhancedHomeStats()
      return enhancedStats.user_activity_stats
    } catch (error) {
      console.error('获取用户活动统计失败:', error)
      // 返回默认数据
      return {
        active_users_today: 0,
        active_users_this_week: 0,
        active_users_this_month: 0,
        new_users_today: 0,
        new_users_this_week: 0,
        new_users_this_month: 0
      }
    }
  }

  // 获取课程分布统计
  async getCourseDistribution(): Promise<CourseDistribution> {
    try {
      const enhancedStats = await this.getEnhancedHomeStats()
      return enhancedStats.course_distribution
    } catch (error) {
      console.error('获取课程分布统计失败:', error)
      throw new Error('获取课程分布统计失败')
    }
  }

  // 获取月度统计
  async getMonthlyStats(): Promise<MonthlyStat[]> {
    try {
      const enhancedStats = await this.getEnhancedHomeStats()
      return enhancedStats.monthly_stats
    } catch (error) {
      console.error('获取月度统计失败:', error)
      throw new Error('获取月度统计失败')
    }
  }

  // 获取TOP评分课程
  async getTopRatedCourses(limit: number = 6): Promise<TopRatedCourse[]> {
    try {
      const enhancedStats = await this.getEnhancedHomeStats()
      return enhancedStats.top_rated_courses.slice(0, limit)
    } catch (error) {
      console.error('获取TOP评分课程失败:', error)
      throw new Error('获取TOP评分课程失败')
    }
  }

  // 获取最新课程
  async getRecentCourses(limit: number = 6): Promise<RecentCourse[]> {
    try {
      const enhancedStats = await this.getEnhancedHomeStats()
      return enhancedStats.recent_courses.slice(0, limit)
    } catch (error) {
      console.error('获取最新课程失败:', error)
      throw new Error('获取最新课程失败')
    }
  }

  // 获取热门课程
  async getPopularCourses(limit: number = 6): Promise<PopularCourse[]> {
    try {
      const enhancedStats = await this.getEnhancedHomeStats()
      return enhancedStats.popular_courses.slice(0, limit)
    } catch (error) {
      console.error('获取热门课程失败:', error)
      throw new Error('获取热门课程失败')
    }
  }
}

// 创建统计API实例
export const statsApi = new StatsApi()

// 简单的统计计算方法
export const calculateStats = {
  // 计算增长率
  growthRate(current: number, previous: number): number {
    if (previous === 0) return current > 0 ? 100 : 0
    return ((current - previous) / previous) * 100
  },

  // 格式化数字
  formatNumber(num: number): string {
    if (num >= 1000000) {
      return (num / 1000000).toFixed(1) + 'M'
    } else if (num >= 1000) {
      return (num / 1000).toFixed(1) + 'K'
    }
    return num.toString()
  },

  // 格式化货币
  formatCurrency(amount: number): string {
    return new Intl.NumberFormat('zh-CN', {
      style: 'currency',
      currency: 'CNY'
    }).format(amount)
  },

  // 格式化百分比
  formatPercentage(value: number, decimals: number = 1): string {
    return value.toFixed(decimals) + '%'
  },

  // 格式化日期
  formatDate(dateStr: string): string {
    const date = new Date(dateStr)
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    })
  },

  // 计算平均评分
  calculateAverageRating(ratings: { score: number }[]): number {
    if (ratings.length === 0) return 0
    const sum = ratings.reduce((acc, rating) => acc + rating.score, 0)
    return sum / ratings.length
  }
}

// 导出类型定义
export default {
  statsApi,
  calculateStats
}