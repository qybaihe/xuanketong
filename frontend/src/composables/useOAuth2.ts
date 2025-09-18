import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { OAuth2Service } from '@/services/oauth2Api'
import { OAuth2Utils } from '@/utils/oauth2'
import type { OAuth2Error } from '@/types/oauth2'

/**
 * OAuth2组合式函数
 */
export function useOAuth2() {
  const router = useRouter()
  const loading = ref(false)
  const error = ref<string | null>(null)

  // 计算属性
  const isLoading = computed(() => loading.value)
  const hasError = computed(() => !!error.value)

  /**
   * 启动OAuth2登录
   */
  const startLogin = async (): Promise<void> => {
    loading.value = true
    error.value = null

    try {
      await OAuth2Service.startOAuth2Login()
    } catch (err) {
      const oauth2Error = OAuth2Utils.parseError(err)
      error.value = OAuth2Utils.formatErrorMessage(oauth2Error)
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * 处理OAuth2回调
   */
  const handleCallback = async (code: string, state: string): Promise<void> => {
    loading.value = true
    error.value = null

    try {
      const response = await OAuth2Service.handleOAuth2Callback(code, state)
      
      // 这里可以触发登录成功的事件或调用其他服务
      console.log('OAuth2登录成功:', response)
      
      // 跳转到首页
      router.push('/')
    } catch (err) {
      const oauth2Error = OAuth2Utils.parseError(err)
      error.value = OAuth2Utils.formatErrorMessage(oauth2Error)
      
      // 跳转到登录页面并显示错误
      router.push('/auth?error=oauth2_failed')
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * 清除错误
   */
  const clearError = (): void => {
    error.value = null
  }

  /**
   * 检查URL是否包含OAuth2参数
   */
  const checkOAuth2Params = (url: string): boolean => {
    return OAuth2Utils.hasOAuth2Params(url)
  }

  /**
   * 从URL提取OAuth2参数
   */
  const extractOAuth2Params = (url: string) => {
    return OAuth2Utils.extractOAuth2Params(url)
  }

  return {
    // 状态
    loading: isLoading,
    error: computed(() => error.value),
    hasError,
    
    // 方法
    startLogin,
    handleCallback,
    clearError,
    checkOAuth2Params,
    extractOAuth2Params
  }
}

export default useOAuth2

