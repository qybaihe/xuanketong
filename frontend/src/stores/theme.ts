import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface ThemeColors {
  primary: string
  secondary: string
  background: string
  surface: string
  text: string
  textSecondary: string
  border: string
  accent: string
  success: string
  warning: string
  error: string
}

export interface Theme {
  name: string
  colors: ThemeColors
}

export type DarkModePreference = 'system' | 'enabled' | 'disabled'

export const useThemeStore = defineStore('theme', () => {
  // 主题配置
  const themes: Theme[] = [
    {
      name: 'default',
      colors: {
        primary: '#F7D074',
        secondary: '#FF6B6B',
        background: '#FEF6F7',
        surface: '#FFFFFF',
        text: '#1A1A1A',
        textSecondary: '#888888',
        border: '#000000',
        accent: '#4ECDC4',
        success: '#4CAF50',
        warning: '#FF9800',
        error: '#F44336'
      }
    },
    {
      name: 'modern',
      colors: {
        primary: '#93C5FD',
        secondary: '#8E8E93',
        background: '#F2F2F7',
        surface: '#FFFFFF',
        text: '#000000',
        textSecondary: '#8E8E93',
        border: '#C6C6C8',
        accent: '#60A5FA',
        success: '#34C759',
        warning: '#FF9500',
        error: '#FF3B30'
      }
    }
  ]

  // 当前主题
  const currentTheme = ref<string>('default')
  
  // 夜间模式设置
  const darkModePreference = ref<DarkModePreference>('system')
  
  // 是否启用夜间模式
  const isDarkMode = computed(() => {
    if (darkModePreference.value === 'enabled') return true
    if (darkModePreference.value === 'disabled') return false
    
    // 跟随系统
    return window.matchMedia('(prefers-color-scheme: dark)').matches
  })

  // 当前主题配置
  const currentThemeConfig = computed(() => {
    const theme = themes.find(t => t.name === currentTheme.value)
    if (!theme) return themes[0]
    
    if (isDarkMode.value) {
      // 夜间模式颜色调整
      return {
        ...theme,
        colors: {
          ...theme.colors,
          background: '#0F172A',
          surface: '#1E293B',
          text: '#F8FAFC',
          textSecondary: '#94A3B8',
          border: '#334155'
        }
      }
    }
    
    return theme
  })

  // CSS 变量
  const cssVariables = computed(() => {
    const colors = currentThemeConfig.value.colors
    return {
      '--color-primary': colors.primary,
      '--color-secondary': colors.secondary,
      '--color-background': colors.background,
      '--color-surface': colors.surface,
      '--color-text': colors.text,
      '--color-textSecondary': colors.textSecondary,
      '--color-border': colors.border,
      '--color-accent': colors.accent,
      '--color-success': colors.success,
      '--color-warning': colors.warning,
      '--color-error': colors.error
    }
  })

  // 应用主题到DOM
  const applyTheme = () => {
    const root = document.documentElement
    const body = document.body
    
    // 设置CSS变量
    Object.entries(cssVariables.value).forEach(([key, value]) => {
      root.style.setProperty(key, value)
    })
    
    // 设置夜间模式类
    if (isDarkMode.value) {
      body.classList.add('dark-mode')
    } else {
      body.classList.remove('dark-mode')
    }
  }

  // 设置主题
  const setTheme = (themeName: string) => {
    currentTheme.value = themeName
    localStorage.setItem('theme', themeName)
    applyTheme()
  }

  // 设置夜间模式
  const setDarkMode = (preference: DarkModePreference) => {
    darkModePreference.value = preference
    localStorage.setItem('darkMode', preference)
    applyTheme()
  }

  // 初始化主题
  const initTheme = () => {
    // 从localStorage加载设置
    const savedTheme = localStorage.getItem('theme')
    const savedDarkMode = localStorage.getItem('darkMode') as DarkModePreference
    
    if (savedTheme) currentTheme.value = savedTheme
    if (savedDarkMode) darkModePreference.value = savedDarkMode
    
    // 监听系统主题变化
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    mediaQuery.addEventListener('change', applyTheme)
    
    applyTheme()
  }

  return {
    themes,
    currentTheme,
    darkModePreference,
    isDarkMode,
    currentThemeConfig,
    cssVariables,
    setTheme,
    setDarkMode,
    initTheme,
    applyTheme
  }
})
