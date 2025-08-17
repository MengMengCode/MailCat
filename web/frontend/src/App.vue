<template>
  <div id="app">
    <router-view />
    <Toast />
    <ConfirmDialog />
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
/* Vercel 风格的全局样式 */
:root {
  /* 主色调 */
  --primary-color: #000;
  --secondary-color: #fff;
  --accent-color: #0070f3;
  
  /* 边框和分割线 */
  --border-color: #eaeaea;
  --border-color-light: #f0f0f0;
  --border-color-dark: #d4d4d4;
  
  /* 背景色 */
  --background: #fff;
  --background-secondary: #fafafa;
  --surface: #fff;
  --surface-hover: #f9f9f9;
  --surface-active: #f0f0f0;
  
  /* 文字颜色 */
  --text-primary: #000;
  --text-secondary: #666;
  --text-tertiary: #999;
  --text-inverse: #fff;
  
  /* 状态颜色 */
  --success-color: #0070f3;
  --warning-color: #f5a623;
  --error-color: #e00;
  --info-color: #0070f3;
  
  /* 阴影 */
  --shadow-small: 0 1px 3px rgba(0, 0, 0, 0.1);
  --shadow-medium: 0 4px 12px rgba(0, 0, 0, 0.08);
  --shadow-large: 0 8px 30px rgba(0, 0, 0, 0.12);
  
  /* 圆角 */
  --radius-small: 4px;
  --radius-medium: 6px;
  --radius-large: 8px;
  --radius-xl: 12px;
  
  /* 间距 */
  --spacing-xs: 4px;
  --spacing-sm: 8px;
  --spacing-md: 16px;
  --spacing-lg: 24px;
  --spacing-xl: 32px;
  --spacing-2xl: 48px;
  
  /* 字体 */
  --font-mono: 'SF Mono', 'Monaco', 'Inconsolata', 'Roboto Mono', 'Consolas', monospace;
  --font-sans: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: var(--font-sans);
  background: var(--background);
  color: var(--text-primary);
  line-height: 1.6;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  min-height: 100vh;
  background: var(--background);
}

/* 覆盖 PrimeVue 默认样式以匹配 Vercel 风格 */
.p-component {
  font-family: var(--font-sans);
}

.p-button {
  background: var(--primary-color);
  color: var(--text-inverse);
  border: 1px solid var(--primary-color);
  border-radius: var(--radius-medium);
  font-weight: 500;
  font-size: 0.875rem;
  padding: 0.75rem 1rem;
  transition: all 0.15s ease;
  cursor: pointer;
}

.p-button:hover:not(:disabled) {
  background: #333;
  border-color: #333;
  transform: translateY(-1px);
  box-shadow: var(--shadow-medium);
}

.p-button:active {
  transform: translateY(0);
}

.p-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.p-button.p-button-outlined {
  background: transparent;
  color: var(--text-primary);
  border: 1px solid var(--border-color);
}

.p-button.p-button-outlined:hover:not(:disabled) {
  background: var(--surface-hover);
  border-color: var(--border-color-dark);
  color: var(--text-primary);
}

.p-button.p-button-text {
  background: transparent;
  color: var(--text-secondary);
  border: none;
  padding: 0.5rem 0.75rem;
}

.p-button.p-button-text:hover:not(:disabled) {
  background: var(--surface-hover);
  color: var(--text-primary);
}

.p-inputtext {
  background: var(--surface);
  border: 1px solid var(--border-color);
  color: var(--text-primary);
  border-radius: var(--radius-medium);
  padding: 0.75rem;
  font-size: 0.875rem;
  transition: all 0.15s ease;
  font-family: var(--font-sans);
}

.p-inputtext:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px rgba(0, 0, 0, 0.1);
}

.p-inputtext::placeholder {
  color: var(--text-tertiary);
}

.p-card {
  background: var(--surface);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-large);
  box-shadow: var(--shadow-small);
  overflow: hidden;
}

.p-card .p-card-header {
  background: var(--surface);
  border-bottom: 1px solid var(--border-color);
  padding: var(--spacing-lg);
}

.p-card .p-card-content {
  padding: var(--spacing-lg);
}

.p-datatable {
  background: var(--surface);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-large);
  overflow: hidden;
}

.p-datatable .p-datatable-thead > tr > th {
  background: var(--background-secondary);
  border-bottom: 1px solid var(--border-color);
  color: var(--text-primary);
  font-weight: 600;
  font-size: 0.875rem;
  padding: var(--spacing-md);
}

.p-datatable .p-datatable-tbody > tr {
  background: var(--surface);
  border-bottom: 1px solid var(--border-color-light);
  transition: background-color 0.15s ease;
}

.p-datatable .p-datatable-tbody > tr:hover {
  background: var(--surface-hover);
}

.p-datatable .p-datatable-tbody > tr > td {
  color: var(--text-primary);
  border: none;
  padding: var(--spacing-md);
  font-size: 0.875rem;
}

.p-dialog {
  background: var(--surface);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-large);
  overflow: hidden;
}

.p-dialog .p-dialog-header {
  background: var(--surface);
  border-bottom: 1px solid var(--border-color);
  color: var(--text-primary);
  padding: var(--spacing-lg);
}

.p-dialog .p-dialog-content {
  background: var(--surface);
  color: var(--text-primary);
  padding: var(--spacing-lg);
}

.p-tag {
  background: var(--background-secondary);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-small);
  font-size: 0.75rem;
  font-weight: 500;
  padding: 0.25rem 0.5rem;
}

.p-tag.p-tag-info {
  background: rgba(0, 112, 243, 0.1);
  color: var(--info-color);
  border-color: rgba(0, 112, 243, 0.2);
}

.p-tag.p-tag-success {
  background: rgba(0, 112, 243, 0.1);
  color: var(--success-color);
  border-color: rgba(0, 112, 243, 0.2);
}

/* Toast 样式 */
.p-toast .p-toast-message {
  background: var(--surface);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-large);
  box-shadow: var(--shadow-medium);
}

/* 自定义滚动条 */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: var(--background-secondary);
}

::-webkit-scrollbar-thumb {
  background: var(--border-color-dark);
  border-radius: var(--radius-small);
}

::-webkit-scrollbar-thumb:hover {
  background: var(--text-tertiary);
}

/* 通用工具类 */
.text-primary { color: var(--text-primary); }
.text-secondary { color: var(--text-secondary); }
.text-tertiary { color: var(--text-tertiary); }
.bg-surface { background: var(--surface); }
.bg-secondary { background: var(--background-secondary); }
.border-color { border-color: var(--border-color); }
.shadow-small { box-shadow: var(--shadow-small); }
.shadow-medium { box-shadow: var(--shadow-medium); }
.radius-medium { border-radius: var(--radius-medium); }
.radius-large { border-radius: var(--radius-large); }

/* 全局动画和交互效果 */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slideInLeft {
  from {
    opacity: 0;
    transform: translateX(-30px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@keyframes slideInRight {
  from {
    opacity: 0;
    transform: translateX(30px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@keyframes scaleIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

@keyframes shimmer {
  0% {
    background-position: -200px 0;
  }
  100% {
    background-position: calc(200px + 100%) 0;
  }
}

/* 页面过渡动画 */
.page-enter-active,
.page-leave-active {
  transition: all 0.3s ease;
}

.page-enter-from {
  opacity: 0;
  transform: translateX(20px);
}

.page-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}

/* 路由视图动画 */
.router-view {
  animation: fadeInUp 0.4s ease-out;
}

/* 增强的卡片悬停效果 */
.p-card {
  transition: all 0.2s ease;
}

.p-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-medium);
}

/* 增强的表格行悬停效果 */
.p-datatable .p-datatable-tbody > tr {
  transition: all 0.15s ease;
}

.p-datatable .p-datatable-tbody > tr:hover {
  transform: translateX(2px);
  box-shadow: inset 3px 0 0 var(--accent-color);
}

/* 标签动画 */
.p-tag {
  transition: all 0.15s ease;
}

.p-tag:hover {
  transform: scale(1.05);
}

/* 对话框进入动画 */
.p-dialog {
  animation: scaleIn 0.3s ease-out;
}

/* 吐司通知动画 */
.p-toast .p-toast-message {
  animation: slideInRight 0.3s ease-out;
}

/* 按钮点击反馈增强 */
.p-button:active {
  transform: translateY(1px) scale(0.98);
}

/* 输入框焦点动画增强 */
.p-inputtext {
  position: relative;
  overflow: hidden;
}

.p-inputtext::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  width: 0;
  height: 2px;
  background: var(--accent-color);
  transition: all 0.3s ease;
  transform: translateX(-50%);
}

.p-inputtext:focus::after {
  width: 100%;
}

/* 加载状态动画 */
.loading-skeleton {
  background: linear-gradient(90deg, var(--surface) 25%, var(--surface-hover) 50%, var(--surface) 75%);
  background-size: 200px 100%;
  animation: shimmer 1.5s infinite;
}

/* 悬停提升效果 */
.hover-lift {
  transition: all 0.15s ease;
}

.hover-lift:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-medium);
}

/* 菜单项悬停效果 */
.p-menu .p-menuitem-link {
  transition: all 0.15s ease;
}

.p-menu .p-menuitem-link:hover {
  transform: translateX(4px);
}

/* 进度条动画 */
.p-progressbar .p-progressbar-value {
  transition: width 0.6s ease;
}

/* 下拉菜单动画 */
.p-dropdown-panel {
  animation: fadeInDown 0.2s ease-out;
}

/* 工具提示动画 */
.p-tooltip {
  animation: fadeInUp 0.2s ease-out;
}

/* 确认对话框动画 */
.p-confirm-dialog {
  animation: scaleIn 0.25s ease-out;
}

/* 选择状态动画 */
.p-checkbox .p-checkbox-box,
.p-radiobutton .p-radiobutton-box {
  transition: all 0.2s ease;
}

/* 滚动条动画 */
::-webkit-scrollbar-thumb {
  transition: background-color 0.2s ease;
}

/* 响应式动画优化 */
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}

/* 触摸设备优化 */
@media (hover: none) and (pointer: coarse) {
  .hover-lift:hover,
  .p-card:hover,
  .p-datatable .p-datatable-tbody > tr:hover {
    transform: none;
  }
  
  .p-tag:hover {
    transform: none;
  }
}

/* 焦点可见性增强 */
:focus-visible {
  outline: 2px solid var(--accent-color);
  outline-offset: 2px;
  border-radius: var(--radius-small);
}

/* 选择文本样式 */
::selection {
  background: rgba(0, 112, 243, 0.2);
  color: var(--text-primary);
}

/* 平滑滚动 */
html {
  scroll-behavior: smooth;
}

/* 防止布局偏移 */
img, video {
  max-width: 100%;
  height: auto;
}

/* 改善触摸体验 */
button, [role="button"] {
  touch-action: manipulation;
}
</style>