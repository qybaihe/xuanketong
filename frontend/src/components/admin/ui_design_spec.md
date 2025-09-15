# 选课通系统 UI 设计规范

## 1. 设计原则

### 1.1 iOS Human Interface Guidelines (HIG) 核心原则
- **清晰性 (Clarity)**: 确保界面元素清晰易读，文字大小和对比度符合iOS标准
- **遵从性 (Deference)**: UI应该流畅且不干扰内容，让用户专注于任务
- **深度感 (Depth)**: 通过视觉层次和过渡动画传达界面的深度和交互状态

### 1.2 高技感设计原则
- **玻璃态效果 (Glassmorphism)**: 使用半透明背景和模糊效果营造现代感
- **动态阴影系统**: 根据光线角度和距离计算阴影，增强立体感
- **微交互动画**: 细腻的过渡动画和反馈效果
- **渐变色彩系统**: 使用现代渐变色彩增强视觉吸引力

## 2. 视觉系统

### 2.1 色彩系统

#### 主色调配色方案
基于您提供的专业色彩配置，我们采用以下配色方案：

```css
:root {
  /* 主色调 - 基于AdminLTE品牌色 #2fa914 */
  --primary-color: #2fa914;
  --primary-color-light: #4fc830;
  --primary-color-dark: #1e7a0f;
  
  /* Discreet Palette 配色方案 */
  --discreet-primary: #2fa914;
  --discreet-secondary: #7d9871;
  --discreet-background: #f1fdeb;
  --discreet-accent: #9999cc;
  
  /* Matching Palette 配色方案 */
  --match-primary: #2fa914;
  --match-secondary: #404a3b;
  --match-tertiary: #a3ae9e;
  --match-accent: #00a8ff;
  --match-accent-dark: #0074e2;
  
  /* Natural Palette 配色方案 */
  --natural-primary: #2fa914;
  --natural-secondary: #7d9871;
  --natural-background: #f5fbf2;
  --natural-accent: #f6edd9;
  
  /* 辅助色 */
  --secondary-color: #007aff;
  --accent-color: #ff3b30;
  --success-color: #34c759;
  --warning-color: #ff9500;
  --info-color: #5ac8fa;
  
  /* 中性色 */
  --text-primary: #1d1d1f;
  --text-secondary: #6e6e73;
  --text-tertiary: #8e8e93;
  --text-quaternary: #c7c7cc;
  
  /* 背景色 */
  --background-primary: #ffffff;
  --background-secondary: #f2f2f7;
  --background-tertiary: #e5e5ea;
  --background-blur: rgba(255, 255, 255, 0.72);
  
  /* 边框色 */
  --separator-color: #c6c6c8;
  --separator-opaque: rgba(60, 60, 67, 0.36);
}
```

#### 渐变色彩系统
```css
/* 主色调渐变 */
--gradient-primary: linear-gradient(135deg, #2fa914 0%, #00a350 100%);
--gradient-primary-alt: linear-gradient(135deg, #2fa914 0%, #009a52 100%);

/* Matching 渐变 */
--gradient-match: linear-gradient(135deg, #2fa914 0%, #00a350 0%, #00997e 33%, #008da5 66%, #007fbf 100%);
--gradient-match-simplified: linear-gradient(135deg, #2fa914 0%, #006ec8 100%);

/* Skip 渐变 */
--gradient-skip: linear-gradient(135deg, #2fa914 0%, #ffd688 25%, #dd9f54 50%, #a26b22 100%);
--gradient-skip-shade: linear-gradient(135deg, #2fa914 0%, #9bb823 33%, #e2c650 66%, #ffd688 100%);

/* Natural 渐变 */
--gradient-natural: linear-gradient(135deg, #2fa914 0%, #7d9871 50%, #f5fbf2 100%);
--gradient-natural-warm: linear-gradient(135deg, #2fa914 0%, #f5fbf2 50%, #f6edd9 100%);

/* 玻璃态渐变 */
--glass-gradient: linear-gradient(135deg,
  rgba(255, 255, 255, 0.25) 0%,
  rgba(255, 255, 255, 0.1) 100%);
--glass-gradient-natural: linear-gradient(135deg,
  rgba(245, 251, 242, 0.25) 0%,
  rgba(246, 237, 217, 0.1) 100%);
```

#### 主题色彩方案
```css
/* Discreet 主题 */
.theme-discreet {
  --primary: var(--discreet-primary);
  --secondary: var(--discreet-secondary);
  --background: var(--discreet-background);
  --accent: var(--discreet-accent);
}

/* Matching 主题 */
.theme-matching {
  --primary: var(--match-primary);
  --secondary: var(--match-secondary);
  --tertiary: var(--match-tertiary);
  --accent: var(--match-accent);
  --accent-dark: var(--match-accent-dark);
}

/* Natural 主题 */
.theme-natural {
  --primary: var(--natural-primary);
  --secondary: var(--natural-secondary);
  --background: var(--natural-background);
  --accent: var(--natural-accent);
}
```

#### 状态色彩
```css
/* 成功状态 - 基于主色调 */
--success-light: #e0f8d5;
--success-base: #2fa914;
--success-dark: #1e7a0f;

/* 信息状态 */
--info-light: #cdf6ff;
--info-base: #004d72;
--info-dark: #003660;

/* 警告状态 */
--warning-light: #ffd688;
--warning-base: #dd9f54;
--warning-dark: #a26b22;

/* 错误状态 */
--error-light: #ffe0e0;
--error-base: #ff3b30;
--error-dark: #d70015;
```

#### 渐变色彩
```css
/* 主渐变 */
--gradient-primary: linear-gradient(135deg, #2fa914 0%, #4fc830 100%);
--gradient-secondary: linear-gradient(135deg, #007aff 0%, #5ac8fa 100%);
--gradient-accent: linear-gradient(135deg, #ff3b30 0%, #ff9500 100%);

/* 玻璃态渐变 */
--glass-gradient: linear-gradient(135deg, 
  rgba(255, 255, 255, 0.25) 0%, 
  rgba(255, 255, 255, 0.1) 100%);
```

### 2.2 字体系统

#### 字体栈
```css
font-family: 
  -apple-system, 
  BlinkMacSystemFont, 
  'SF Pro Display', 
  'SF Pro Text', 
  'Helvetica Neue', 
  Helvetica, 
  Arial, 
  sans-serif;
```

#### 字体大小和行高
```css
/* 标题字体 */
--font-size-title1: 28px;    /* line-height: 34px */
--font-size-title2: 22px;    /* line-height: 28px */
--font-size-title3: 20px;    /* line-height: 25px */

/* 正文字体 */
--font-size-body: 17px;      /* line-height: 22px */
--font-size-body2: 15px;     /* line-height: 20px */
--font-size-caption: 12px;  /* line-height: 16px */
--font-size-caption2: 11px; /* line-height: 13px */
```

#### 字重
```css
--font-weight-ultralight: 100;
--font-weight-thin: 200;
--font-weight-light: 300;
--font-weight-regular: 400;
--font-weight-medium: 500;
--font-weight-semibold: 600;
--font-weight-bold: 700;
--font-weight-heavy: 800;
--font-weight-black: 900;
```

### 2.3 间距系统

#### 8pt 网格系统
```css
/* 基础间距单位 */
--spacing-xs: 4px;    /* 0.5 * 8 */
--spacing-sm: 8px;    /* 1 * 8 */
--spacing-md: 16px;   /* 2 * 8 */
--spacing-lg: 24px;   /* 3 * 8 */
--spacing-xl: 32px;   /* 4 * 8 */
--spacing-2xl: 48px;  /* 6 * 8 */
--spacing-3xl: 64px;  /* 8 * 8 */
```

#### 组件间距
```css
/* 卡片间距 */
--card-padding: 16px;
--card-margin: 8px;
--card-radius: 12px;

/* 列表项间距 */
--list-item-height: 44px;
--list-item-padding: 12px 16px;
```

## 3. 布局系统

### 3.1 响应式断点
```css
/* 移动设备 */
--mobile-breakpoint: 320px;
--mobile-max: 768px;

/* 平板设备 */
--tablet-breakpoint: 768px;
--tablet-max: 1024px;

/* 桌面设备 */
--desktop-breakpoint: 1024px;
```

### 3.2 网格系统
```css
/* 容器 */
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 var(--spacing-md);
}

/* 网格布局 */
.grid {
  display: grid;
  gap: var(--spacing-md);
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
}

/* 响应式网格 */
@media (min-width: var(--tablet-breakpoint)) {
  .grid-2 {
    grid-template-columns: repeat(2, 1fr);
  }
  .grid-3 {
    grid-template-columns: repeat(3, 1fr);
  }
  .grid-4 {
    grid-template-columns: repeat(4, 1fr);
  }
}
```

### 3.3 安全区域
```css
/* 处理刘海屏和底部指示条 */
.safe-area-top {
  padding-top: env(safe-area-inset-top);
}

.safe-area-bottom {
  padding-bottom: env(safe-area-inset-bottom);
}

.safe-area-left {
  padding-left: env(safe-area-inset-left);
}

.safe-area-right {
  padding-right: env(safe-area-inset-right);
}
```

## 4. 组件设计

### 4.1 卡片组件

#### 基础卡片
```css
.card {
  background: var(--background-primary);
  border-radius: var(--card-radius);
  padding: var(--card-padding);
  margin: var(--card-margin);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: all 0.3s cubic-bezier(0.25, 0.46, 0.45, 0.94);
}

.card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
```

#### 玻璃态卡片
```css
.card-glass {
  background: var(--background-blur);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.18);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.card-glass:hover {
  background: rgba(255, 255, 255, 0.85);
  box-shadow: 0 12px 48px rgba(0, 0, 0, 0.15);
}
```

### 4.2 按钮组件

#### 基础按钮
```css
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 12px 24px;
  border-radius: 8px;
  font-size: var(--font-size-body);
  font-weight: var(--font-weight-medium);
  border: none;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.25, 0.46, 0.45, 0.94);
  min-height: 44px;
  position: relative;
  overflow: hidden;
}

/* 主色调按钮 */
.btn-primary {
  background: var(--gradient-primary);
  color: white;
  box-shadow: 0 2px 8px rgba(47, 169, 20, 0.2);
}

.btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(47, 169, 20, 0.3);
  background: var(--gradient-primary-alt);
}

.btn-primary:active {
  transform: translateY(0);
  box-shadow: 0 1px 4px rgba(47, 169, 20, 0.2);
}

/* Matching 主题按钮 */
.btn-matching {
  background: var(--gradient-match);
  color: white;
  box-shadow: 0 2px 8px rgba(0, 168, 255, 0.2);
}

.btn-matching:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 168, 255, 0.3);
}

/* Natural 主题按钮 */
.btn-natural {
  background: var(--gradient-natural);
  color: var(--text-primary);
  box-shadow: 0 2px 8px rgba(125, 152, 113, 0.2);
}

.btn-natural:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(125, 152, 113, 0.3);
}

/* Skip 渐变按钮 */
.btn-skip {
  background: var(--gradient-skip);
  color: var(--text-primary);
  box-shadow: 0 2px 8px rgba(221, 159, 84, 0.2);
}

.btn-skip:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(221, 159, 84, 0.3);
}

/* 状态按钮 */
.btn-success {
  background: var(--success-base);
  color: white;
  box-shadow: 0 2px 8px rgba(47, 169, 20, 0.2);
}

.btn-info {
  background: var(--info-base);
  color: white;
  box-shadow: 0 2px 8px rgba(0, 77, 114, 0.2);
}

.btn-warning {
  background: var(--warning-base);
  color: white;
  box-shadow: 0 2px 8px rgba(221, 159, 84, 0.2);
}

.btn-error {
  background: var(--error-base);
  color: white;
  box-shadow: 0 2px 8px rgba(255, 59, 48, 0.2);
}
```

#### 玻璃态按钮
```css
.btn-glass {
  background: var(--glass-gradient);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.18);
  color: var(--text-primary);
}

.btn-glass:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
}
```

### 4.3 输入框组件

#### 基础输入框
```css
.input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid var(--separator-color);
  border-radius: 8px;
  font-size: var(--font-size-body);
  background: var(--background-primary);
  transition: all 0.2s ease;
}

.input:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(47, 169, 20, 0.1);
}

.input-glass {
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.18);
}
```

### 4.4 课程卡片组件

#### 课程卡片布局
```css
.course-card {
  background: var(--background-primary);
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s cubic-bezier(0.25, 0.46, 0.45, 0.94);
  display: flex;
  flex-direction: column;
  height: 100%;
  position: relative;
}

.course-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

/* 课程卡片主题变体 */
.course-card-primary {
  border: 1px solid rgba(47, 169, 20, 0.1);
  box-shadow: 0 2px 8px rgba(47, 169, 20, 0.1);
}

.course-card-primary:hover {
  box-shadow: 0 8px 24px rgba(47, 169, 20, 0.2);
}

.course-card-matching {
  border: 1px solid rgba(0, 168, 255, 0.1);
  box-shadow: 0 2px 8px rgba(0, 168, 255, 0.1);
}

.course-card-matching:hover {
  box-shadow: 0 8px 24px rgba(0, 168, 255, 0.2);
}

.course-card-natural {
  border: 1px solid rgba(125, 152, 113, 0.1);
  box-shadow: 0 2px 8px rgba(125, 152, 113, 0.1);
}

.course-card-natural:hover {
  box-shadow: 0 8px 24px rgba(125, 152, 113, 0.2);
}

.course-card-image {
  width: 100%;
  height: 200px;
  object-fit: cover;
  background: linear-gradient(135deg, #f5f5f5 0%, #e0e0e0 100%);
  position: relative;
}

/* 课程卡片图片主题变体 */
.course-card-primary .course-card-image {
  background: linear-gradient(135deg, #e0f8d5 0%, #2fa914 100%);
}

.course-card-matching .course-card-image {
  background: linear-gradient(135deg, #cdf6ff 0%, #004d72 100%);
}

.course-card-natural .course-card-image {
  background: linear-gradient(135deg, #f5fbf2 0%, #7d9871 100%);
}

.course-card-content {
  padding: var(--spacing-md);
  flex: 1;
  display: flex;
  flex-direction: column;
}

.course-card-title {
  font-size: var(--font-size-title3);
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--spacing-sm);
  color: var(--text-primary);
}

.course-card-description {
  font-size: var(--font-size-body2);
  color: var(--text-secondary);
  margin-bottom: var(--spacing-md);
  flex: 1;
  line-height: 1.5;
}

.course-card-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: var(--font-size-caption);
  color: var(--text-tertiary);
}

.course-card-credits {
  background: var(--gradient-primary);
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: var(--font-size-caption);
  font-weight: var(--font-weight-medium);
}

/* 课程卡片学分主题变体 */
.course-card-primary .course-card-credits {
  background: var(--gradient-primary);
}

.course-card-matching .course-card-credits {
  background: var(--gradient-match-simplified);
}

.course-card-natural .course-card-credits {
  background: var(--gradient-natural);
}

/* 课程卡片标签 */
.course-card-tags {
  display: flex;
  gap: var(--spacing-xs);
  margin-bottom: var(--spacing-sm);
  flex-wrap: wrap;
}

.course-card-tag {
  background: var(--background-secondary);
  color: var(--text-tertiary);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: var(--font-size-caption2);
  font-weight: var(--font-weight-medium);
}

.course-card-tag.primary {
  background: var(--success-light);
  color: var(--success-base);
}

.course-card-tag.secondary {
  background: var(--info-light);
  color: var(--info-base);
}

.course-card-tag.accent {
  background: var(--warning-light);
  color: var(--warning-base);
}
```

#### 玻璃态课程卡片
```css
.course-card-glass {
  background: var(--background-blur);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.18);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.course-card-glass:hover {
  background: rgba(255, 255, 255, 0.85);
  box-shadow: 0 12px 48px rgba(0, 0, 0, 0.15);
}
```

## 5. 动画系统

### 5.1 过渡动画

#### 基础过渡
```css
/* 快速过渡 */
.transition-fast {
  transition: all 0.15s cubic-bezier(0.25, 0.46, 0.45, 0.94);
}

/* 标准过渡 */
.transition-standard {
  transition: all 0.3s cubic-bezier(0.25, 0.46, 0.45, 0.94);
}

/* 慢速过渡 */
.transition-slow {
  transition: all 0.5s cubic-bezier(0.25, 0.46, 0.45, 0.94);
}
```

#### Spring 动画参数
```css
/* 轻量级弹簧 */
.spring-light {
  animation: spring-light 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

/* 标准弹簧 */
.spring-standard {
  animation: spring-standard 0.6s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

/* 重量级弹簧 */
.spring-heavy {
  animation: spring-heavy 0.8s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}
```

### 5.2 微交互动画

#### 按钮点击效果
```css
.btn:active {
  transform: scale(0.95);
}

.btn::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 0;
  height: 0;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.5);
  transform: translate(-50%, -50%);
  transition: width 0.6s, height 0.6s;
}

.btn:active::after {
  width: 300px;
  height: 300px;
}
```

#### 卡片悬停效果
```css
.course-card {
  position: relative;
  overflow: hidden;
}

.course-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, 
    transparent 0%, 
    rgba(255, 255, 255, 0.1) 50%, 
    transparent 100%);
  transform: translateX(-100%);
  transition: transform 0.6s;
}

.course-card:hover::before {
  transform: translateX(100%);
}
```

### 5.3 页面转场动画

#### 淡入淡出
```css
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
```

#### 滑动转场
```css
.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s cubic-bezier(0.25, 0.46, 0.45, 0.94);
}

.slide-enter-from {
  transform: translateX(100%);
  opacity: 0;
}

.slide-leave-to {
  transform: translateX(-100%);
  opacity: 0;
}
```

## 6. 阴影系统

### 6.1 动态阴影
```css
/* 轻量级阴影 */
.shadow-sm {
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

/* 标准阴影 */
.shadow-md {
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

/* 重量级阴影 */
.shadow-lg {
  box-shadow: 0 10px 15px rgba(0, 0, 0, 0.1);
}

/* 超重量级阴影 */
.shadow-xl {
  box-shadow: 0 20px 25px rgba(0, 0, 0, 0.1);
}

/* 玻璃态阴影 */
.shadow-glass {
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}
```

### 6.2 彩色阴影
```css
/* 主色调阴影 */
.shadow-primary {
  box-shadow: 0 4px 12px rgba(47, 169, 20, 0.2);
}

/* 辅助色阴影 */
.shadow-secondary {
  box-shadow: 0 4px 12px rgba(0, 122, 255, 0.2);
}

/* 强调色阴影 */
.shadow-accent {
  box-shadow: 0 4px 12px rgba(255, 59, 48, 0.2);
}
```

## 7. 响应式设计

### 7.1 移动端适配
```css
/* 移动端基础样式 */
@media (max-width: 768px) {
  .container {
    padding: 0 var(--spacing-sm);
  }
  
  .course-grid {
    grid-template-columns: 1fr;
  }
  
  .card {
    margin: var(--spacing-sm) 0;
  }
  
  .btn {
    width: 100%;
  }
}
```

### 7.2 平板端适配
```css
/* 平板端样式 */
@media (min-width: 769px) and (max-width: 1024px) {
  .course-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .container {
    max-width: 900px;
  }
}
```

### 7.3 桌面端适配
```css
/* 桌面端样式 */
@media (min-width: 1025px) {
  .course-grid {
    grid-template-columns: repeat(3, 1fr);
  }
  
  .container {
    max-width: 1200px;
  }
}
```

## 8. 兼容性改造方案

### 8.1 现有组件改造

#### HomeView.vue 改造
```vue
<!-- 现有结构保持不变，主要更新样式 -->
<template>
  <main class="home-container">
    <div class="hero-section">
      <h1 class="hero-title">课程浏览</h1>
      <p class="hero-subtitle">发现优质课程，提升学习体验</p>
    </div>
    
    <div class="filter-section card-glass">
      <form @submit.prevent="fetchCourses" class="filter-form">
        <!-- 现有表单结构 -->
      </form>
    </div>
    
    <div class="course-grid">
      <router-link 
        v-for="course in courses" 
        :key="course.ID" 
        :to="`/courses/${course.ID}`" 
        class="course-card course-card-glass"
      >
        <!-- 现有卡片结构 -->
      </router-link>
    </div>
  </main>
</template>

<style scoped>
.home-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f5f5 0%, #e8f5e8 100%);
  padding: var(--spacing-lg) 0;
}

.hero-section {
  text-align: center;
  margin-bottom: var(--spacing-2xl);
  padding: var(--spacing-lg);
}

.hero-title {
  font-size: var(--font-size-title1);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin-bottom: var(--spacing-sm);
}

.hero-subtitle {
  font-size: var(--font-size-body);
  color: var(--text-secondary);
}

.filter-section {
  margin-bottom: var(--spacing-xl);
  padding: var(--spacing-lg);
}

/* 其他样式更新 */
</style>
```

#### CourseDetailView.vue 改造
```vue
<template>
  <div class="course-detail-container">
    <div class="course-header card-glass">
      <!-- 现有头部结构 -->
    </div>
    
    <div class="course-content">
      <div class="course-info card-glass">
        <!-- 现有信息结构 -->
      </div>
      
      <div class="rating-section card-glass">
        <!-- 现有评分结构 -->
      </div>
      
      <div class="comment-section card-glass">
        <!-- 现有评论结构 -->
      </div>
    </div>
  </div>
</template>

<style scoped>
.course-detail-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f5f5 0%, #e8f5e8 100%);
  padding: var(--spacing-lg) 0;
}

.course-header {
  margin-bottom: var(--spacing-xl);
  padding: var(--spacing-xl);
}

.course-content {
  display: grid;
  gap: var(--spacing-lg);
}

/* 其他样式更新 */
</style>
```

### 8.2 CSS 变量覆盖

#### 创建新的样式文件
```css
/* frontend/src/assets/modern-design.css */
/* 覆盖现有 AdminLTE 样式 */

:root {
  /* 覆盖 AdminLTE 颜色 */
  --primary: var(--primary-color);
  --secondary: var(--secondary-color);
  --success: var(--success-color);
  --info: var(--info-color);
  --warning: var(--warning-color);
  --danger: var(--accent-color);
  
  /* 覆盖 AdminLTE 间距 */
  --padding-base-vertical: 12px;
  --padding-base-horizontal: 16px;
  
  /* 覆盖 AdminLTE 圆角 */
  --border-radius-base: 8px;
  --border-radius-large: 12px;
}

/* 覆盖 AdminLTE 卡片样式 */
.card {
  border: none;
  border-radius: var(--border-radius-large);
  box-shadow: var(--shadow-md);
  background: var(--background-primary);
}

.card:hover {
  box-shadow: var(--shadow-lg);
}

/* 覆盖 AdminLTE 按钮样式 */
.btn {
  border-radius: var(--border-radius-base);
  font-weight: var(--font-weight-medium);
  border: none;
  transition: var(--transition-standard);
}

.btn-primary {
  background: var(--gradient-primary);
}

/* 覆盖 AdminLTE 输入框样式 */
.form-control {
  border-radius: var(--border-radius-base);
  border: 1px solid var(--separator-color);
  padding: var(--padding-base-vertical) var(--padding-base-horizontal);
  transition: var(--transition-standard);
}

.form-control:focus {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(47, 169, 20, 0.1);
}
```

### 8.3 渐进式增强

#### 基础支持检测
```javascript
// 检测浏览器特性支持
const supportsBackdropFilter = CSS.supports('backdrop-filter', 'blur(10px)');
const supportsGrid = CSS.supports('display', 'grid');
const supportsFlexbox = CSS.supports('display', 'flex');

// 根据支持情况添加类名
if (supportsBackdropFilter) {
  document.documentElement.classList.add('supports-backdrop-filter');
}

if (supportsGrid) {
  document.documentElement.classList.add('supports-grid');
}

if (supportsFlexbox) {
  document.documentElement.classList.add('supports-flexbox');
}
```

#### 降级方案
```css
/* 不支持 backdrop-filter 的降级方案 */
.no-backdrop-filter .card-glass {
  background: rgba(255, 255, 255, 0.95);
  border: 1px solid var(--separator-color);
}

/* 不支持 grid 的降级方案 */
.no-grid .course-grid {
  display: flex;
  flex-wrap: wrap;
}

.no-grid .course-card {
  flex: 1 1 300px;
  margin: var(--spacing-sm);
}
```

## 9. 实施建议

### 9.1 实施顺序

1. **第一阶段：基础样式更新**
   - 更新 CSS 变量系统
   - 实现新的色彩方案
   - 更新字体和间距系统

2. **第二阶段：组件改造**
   - 改造卡片组件
   - 更新按钮和输入框样式
   - 实现玻璃态效果

3. **第三阶段：动画和交互**
   - 添加微交互动画
   - 实现页面转场效果
   - 优化用户交互体验

4. **第四阶段：响应式优化**
   - 完善移动端适配
   - 优化平板和桌面端体验
   - 实现渐进式增强

### 9.2 性能优化建议

1. **CSS 优化**
   - 使用 CSS 变量便于主题切换
   - 合理使用 will-change 属性
   - 优化动画性能

2. **图片优化**
   - 使用现代图片格式 (WebP)
   - 实现懒加载
   - 使用响应式图片

3. **JavaScript 优化**
   - 使用 Intersection Observer 实现懒加载
   - 优化事件监听器
   - 使用 requestAnimationFrame 优化动画

### 9.3 测试建议

1. **跨浏览器测试**
   - Chrome、Firefox、Safari、Edge
   - 移动端浏览器
   - 旧版本浏览器兼容性

2. **设备测试**
   - 不同屏幕尺寸
   - 不同分辨率
   - 触摸设备 vs 鼠标设备

3. **性能测试**
   - 页面加载速度
   - 动画流畅度
   - 内存使用情况

## 10. 维护和扩展

### 10.1 设计系统维护

1. **版本控制**
   - 使用语义化版本控制
   - 维护变更日志
   - 定期更新设计规范

2. **组件文档**
   - 创建组件使用文档
   - 提供代码示例
   - 维护设计原则

### 10.2 扩展建议

1. **暗色主题**
   - 设计暗色主题色彩方案
   - 实现主题切换功能
   - 确保可访问性

2. **国际化支持**
   - 支持多语言
   - 考虑 RTL 布局
   - 本地化设计元素

3. **可访问性优化**
   - 键盘导航支持
   - 屏幕阅读器兼容
   - 高对比度模式

---

这个设计规范文档提供了完整的 UI 设计指导，包括视觉系统、组件设计、动画效果、响应式布局等方面。它基于 iOS HIG 原则，结合了现代高技感设计元素，同时考虑了与现有 AdminLTE 样式的兼容性。通过渐进式改造，可以在保持现有功能的基础上，逐步提升用户体验和视觉效果。