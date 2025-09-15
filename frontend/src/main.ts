import '../../AdminLTE/src/html/public/css/adminlte.min.css'
import './assets/main.css'
import './assets/modern-design.css'
import './assets/animations.css'
import './assets/responsive.css'
import './assets/performance.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

// 浏览器特性检测
const supportsBackdropFilter = CSS.supports('backdrop-filter', 'blur(10px)')
const supportsGrid = CSS.supports('display', 'grid')
const supportsFlexbox = CSS.supports('display', 'flex')

// 根据支持情况添加类名
if (supportsBackdropFilter) {
  document.documentElement.classList.add('supports-backdrop-filter')
} else {
  document.documentElement.classList.add('no-backdrop-filter')
}

if (supportsGrid) {
  document.documentElement.classList.add('supports-grid')
} else {
  document.documentElement.classList.add('no-grid')
}

if (supportsFlexbox) {
  document.documentElement.classList.add('supports-flexbox')
} else {
  document.documentElement.classList.add('no-flexbox')
}

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')
