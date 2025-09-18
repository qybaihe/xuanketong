// OAuth2相关类型定义

// SSE集市用户信息接口
export interface SSEUserInfo {
  user_id: number
  name: string
  email: string
  avatar_url: string
  intro: string
  identity: 'student' | 'teacher' | 'organization'
}

// SSE集市Token响应接口
export interface SSETokenResponse {
  code: number
  msg: string
  data: {
    access_token: string
    token_type: string
    expires_in: number
    scope: string
  }
}

// SSE集市用户信息响应接口
export interface SSEUserInfoResponse {
  code: number
  msg: string
  data: SSEUserInfo
}

// OAuth2配置接口
export interface OAuth2Config {
  BASE_URL: string
  APP_ID: string
  REDIRECT_URI: string
  SCOPE: string
}

// OAuth2 State响应接口
export interface OAuth2StateResponse {
  state: string
}

// OAuth2回调响应接口
export interface OAuth2CallbackResponse {
  token: string
  user: {
    id: number
    username: string
    email: string
    nickname: string
    avatar: string
    role: string
    createdAt: string
    updatedAt: string
  }
  message: string
}

// OAuth2错误类型
export interface OAuth2Error {
  error: string
  description?: string
  code?: number
}

