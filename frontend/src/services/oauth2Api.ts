import axios from 'axios'
import type { 
  OAuth2Config, 
  OAuth2StateResponse, 
  OAuth2CallbackResponse
} from '@/types/oauth2'
import { OAuth2Utils } from '@/utils/oauth2'

const BACKEND_BASE_URL = import.meta.env.VITE_BACKEND_BASE_URL

// OAuth2配置
export const OAUTH2_CONFIG: OAuth2Config = {
  BASE_URL: import.meta.env.VITE_OAUTH2_MARKET_BASE_URL,
  APP_ID: import.meta.env.VITE_OAUTH2_MARKET_APP_ID,
  REDIRECT_URI: import.meta.env.VITE_OAUTH2_MARKET_REDIRECT_URI,
  SCOPE: import.meta.env.VITE_OAUTH2_MARKET_SCOPE,
}

// OAuth2服务类
export class OAuth2Service {
  /**
   * 获取OAuth2 state参数
   */
  static async getOAuth2State(): Promise<string> {
    try {
      const response = await axios.get<OAuth2StateResponse>(
        `${BACKEND_BASE_URL}/auth/oauth2/state`
      )
      
      const state = response.data.state
      if (!OAuth2Utils.validateState(state)) {
        throw new Error('无效的state参数')
      }
      
      return state
    } catch (error) {
      console.error('获取OAuth2 state失败:', error)
      const oauth2Error = OAuth2Utils.parseError(error)
      throw new Error(OAuth2Utils.formatErrorMessage(oauth2Error))
    }
  }

  /**
   * 构建OAuth2授权URL
   */
  static buildOAuth2URL(state: string): string {
    const params = new URLSearchParams({
      appid: OAUTH2_CONFIG.APP_ID,
      redirect_uri: OAUTH2_CONFIG.REDIRECT_URI,
      response_type: 'code',
      scope: OAUTH2_CONFIG.SCOPE,
      state: state
    })
    
    return `${OAUTH2_CONFIG.BASE_URL}/new/connect/?${params.toString()}`
  }

  /**
   * 处理OAuth2回调
   */
  static async handleOAuth2Callback(code: string, state: string): Promise<OAuth2CallbackResponse> {
    try {
      // 验证参数
      if (!code || !state) {
        throw new Error('缺少必要的OAuth2参数')
      }
      
      if (!OAuth2Utils.validateState(state)) {
        throw new Error('无效的state参数')
      }
      
      const response = await axios.get<OAuth2CallbackResponse>(
        `${BACKEND_BASE_URL}/auth/oauth2/callback`,
        {
          params: { code, state }
        }
      )
      return response.data
    } catch (error) {
      console.error('OAuth2回调处理失败:', error)
      const oauth2Error = OAuth2Utils.parseError(error)
      throw new Error(OAuth2Utils.formatErrorMessage(oauth2Error))
    }
  }

  /**
   * 启动OAuth2登录流程
   */
  static async startOAuth2Login(): Promise<void> {
    try {
      // 1. 获取state
      const state = await this.getOAuth2State()
      
      // 2. 构建授权URL
      const authURL = this.buildOAuth2URL(state)
      
      // 3. 跳转到授权页面
      window.location.href = authURL
    } catch (error) {
      console.error('OAuth2登录启动失败:', error)
      throw error
    }
  }
}

// 导出默认实例
export default OAuth2Service
