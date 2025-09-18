import type { OAuth2Error } from '@/types/oauth2'

/**
 * OAuth2工具函数
 */
export class OAuth2Utils {
  /**
   * 生成随机state参数
   */
  static generateState(): string {
    const array = new Uint8Array(16)
    crypto.getRandomValues(array)
    return Array.from(array, byte => byte.toString(16).padStart(2, '0')).join('')
  }

  /**
   * 验证state参数格式
   */
  static validateState(state: string): boolean {
    return /^[a-f0-9]{32}$/.test(state)
  }

  /**
   * 解析OAuth2错误
   */
  static parseError(error: any): OAuth2Error {
    if (error.response?.data) {
      return {
        error: error.response.data.error || '未知错误',
        description: error.response.data.message,
        code: error.response.status
      }
    }
    
    return {
      error: error.message || '网络错误',
      description: '请检查网络连接'
    }
  }

  /**
   * 格式化OAuth2错误信息
   */
  static formatErrorMessage(error: OAuth2Error): string {
    let message = error.error
    
    if (error.description) {
      message += `: ${error.description}`
    }
    
    if (error.code) {
      message += ` (错误代码: ${error.code})`
    }
    
    return message
  }

  /**
   * 检查URL参数是否包含OAuth2回调参数
   */
  static hasOAuth2Params(url: string): boolean {
    const urlObj = new URL(url)
    return !!(urlObj.searchParams.get('code') && urlObj.searchParams.get('state'))
  }

  /**
   * 从URL中提取OAuth2参数
   */
  static extractOAuth2Params(url: string): { code: string; state: string } | null {
    const urlObj = new URL(url)
    const code = urlObj.searchParams.get('code')
    const state = urlObj.searchParams.get('state')
    
    if (code && state) {
      return { code, state }
    }
    
    return null
  }

  /**
   * 清理URL中的OAuth2参数
   */
  static cleanOAuth2Params(url: string): string {
    const urlObj = new URL(url)
    urlObj.searchParams.delete('code')
    urlObj.searchParams.delete('state')
    urlObj.searchParams.delete('error')
    urlObj.searchParams.delete('error_description')
    return urlObj.toString()
  }
}

export default OAuth2Utils

