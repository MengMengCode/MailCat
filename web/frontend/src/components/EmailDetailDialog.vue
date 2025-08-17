<template>
  <Dialog
    :visible="visible"
    :modal="true"
    :closable="true"
    :draggable="false"
    class="email-dialog"
    :style="{ width: '90vw', maxWidth: '1200px' }"
    @update:visible="$emit('hide')"
    @hide="$emit('hide')"
  >
    <template #header>
      <div class="dialog-header">
        <h3 class="dialog-title">邮件详情</h3>
        <div class="email-id">ID: {{ email?.id }}</div>
      </div>
    </template>

    <div v-if="loading" class="loading-container">
      <ProgressSpinner />
      <p>加载中...</p>
    </div>

    <div v-else-if="email" class="email-content">
      <!-- 邮件基本信息 -->
      <div class="email-info">
        <div class="info-grid">
          <div class="info-item">
            <label>主题</label>
            <div class="info-value">{{ email.subject || '(无主题)' }}</div>
          </div>
          <div class="info-item">
            <label>发件人</label>
            <div class="info-value">
              <Tag severity="info">{{ email.from }}</Tag>
            </div>
          </div>
          <div class="info-item">
            <label>收件人</label>
            <div class="info-value">
              <Tag severity="success">{{ email.to }}</Tag>
            </div>
          </div>
          <div class="info-item">
            <label>接收时间</label>
            <div class="info-value">{{ formatTime(email.received_at) }}</div>
          </div>
        </div>
      </div>

      <Divider />

      <!-- 邮件内容标签页 -->
      <div class="email-tabs">
        <div class="tab-buttons">
          <button
            v-for="tab in tabs"
            :key="tab.key"
            :class="['tab-button', { active: activeTab === tab.key }]"
            @click="activeTab = tab.key"
          >
            <i :class="tab.icon"></i>
            {{ tab.label }}
          </button>
        </div>

        <div class="tab-content">
          <!-- HTML 渲染 -->
          <div v-if="activeTab === 'html'" class="tab-pane">
            <div v-if="email.html_body && email.html_body.trim()" class="email-html-content" v-html="email.html_body"></div>
            <div v-else class="empty-content">
              <i class="pi pi-info-circle"></i>
              <p>此邮件没有HTML内容</p>
            </div>
          </div>

          <!-- HTML 源码 -->
          <div v-if="activeTab === 'source'" class="tab-pane">
            <pre v-if="email.html_body && email.html_body.trim()" class="email-source-content">{{ email.html_body }}</pre>
            <div v-else class="empty-content">
              <i class="pi pi-info-circle"></i>
              <p>此邮件没有HTML源码</p>
            </div>
          </div>
        </div>
      </div>
    </div>

  </Dialog>
</template>

<script>
import { ref, watch, computed } from 'vue'
import { emailAPI } from '../services/api'

export default {
  name: 'EmailDetailDialog',
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    emailId: {
      type: [Number, String],
      default: null
    }
  },
  emits: ['hide'],
  setup(props) {
    const email = ref(null)
    const loading = ref(false)
    const activeTab = ref('html')

    const tabs = [
      { key: 'html', label: 'HTML渲染', icon: 'pi pi-eye' },
      { key: 'source', label: 'HTML源码', icon: 'pi pi-code' }
    ]

    const formatTime = (timeStr) => {
      if (!timeStr) return '-'
      const date = new Date(timeStr)
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      })
    }

    const loadEmailDetail = async (id) => {
      if (!id) return

      loading.value = true
      try {
        const response = await emailAPI.getEmailById(id)
        email.value = response.data
        console.log('邮件详情数据:', email.value)
        console.log('HTML内容:', email.value?.html_body)
        console.log('纯文本内容:', email.value?.body)
        console.log('头部信息:', email.value?.headers)
      } catch (error) {
        console.error('加载邮件详情失败:', error)
      } finally {
        loading.value = false
      }
    }

    // 监听 emailId 变化
    watch(() => props.emailId, (newId) => {
      if (newId && props.visible) {
        loadEmailDetail(newId)
      }
    })

    // 监听对话框显示状态
    watch(() => props.visible, (visible) => {
      if (visible && props.emailId) {
        loadEmailDetail(props.emailId)
      } else if (!visible) {
        email.value = null
        activeTab.value = 'html'
      }
    })

    return {
      email,
      loading,
      activeTab,
      tabs,
      formatTime
    }
  }
}
</script>

<style scoped>
.email-dialog {
  --dialog-bg: var(--surface);
  --border-color: var(--border-color);
  --text-color: var(--text-primary);
  --text-muted: var(--text-secondary);
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.dialog-title {
  color: var(--text-primary);
  font-size: 1.375rem;
  font-weight: 700;
  margin: 0;
  letter-spacing: -0.025em;
}

.email-id {
  color: var(--text-secondary);
  font-size: 0.8125rem;
  font-family: var(--font-mono);
  font-weight: 600;
  background: var(--background-secondary);
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--radius-small);
  border: 1px solid var(--border-color);
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-2xl);
  color: var(--text-secondary);
  gap: var(--spacing-md);
}

.email-content {
  max-height: 75vh;
  overflow-y: auto;
}

.email-info {
  margin-bottom: var(--spacing-lg);
  background: var(--background-secondary);
  border: 1px solid var(--border-color-light);
  border-radius: var(--radius-large);
  padding: var(--spacing-lg);
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: var(--spacing-lg);
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.info-item label {
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.1em;
}

.info-value {
  color: var(--text-primary);
  font-size: 0.9375rem;
  font-weight: 500;
}

.email-tabs {
  margin-top: var(--spacing-lg);
}

.tab-buttons {
  display: flex;
  gap: 0;
  margin-bottom: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
  background: var(--background-secondary);
  border-radius: var(--radius-medium) var(--radius-medium) 0 0;
  padding: var(--spacing-xs);
}

.tab-button {
  padding: var(--spacing-sm) var(--spacing-md);
  background: transparent;
  border: none;
  color: var(--text-secondary);
  font-size: 0.8125rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
  border-radius: var(--radius-small);
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  position: relative;
}

.tab-button:hover {
  color: var(--text-primary);
  background: var(--surface-hover);
}

.tab-button.active {
  color: var(--text-primary);
  background: var(--surface);
  box-shadow: var(--shadow-small);
}

.tab-content {
  min-height: 350px;
  background: var(--surface);
  border: 1px solid var(--border-color);
  border-radius: 0 0 var(--radius-large) var(--radius-large);
}

.tab-pane {
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.email-html-content {
  background: var(--surface);
  border: none;
  border-radius: 0 0 var(--radius-large) var(--radius-large);
  padding: var(--spacing-xl);
  min-height: 350px;
  overflow: auto;
  line-height: 1.6;
}

.email-text-content,
.email-source-content,
.email-headers-content {
  background: var(--background-secondary);
  border: none;
  border-radius: 0 0 var(--radius-large) var(--radius-large);
  padding: var(--spacing-lg);
  color: var(--text-primary);
  font-family: var(--font-mono);
  font-size: 0.8125rem;
  line-height: 1.6;
  white-space: pre-wrap;
  word-wrap: break-word;
  max-height: 450px;
  overflow: auto;
  margin: 0;
}

.email-source-content,
.email-headers-content {
  background: #1a1a1a;
  color: #e5e5e5;
  border: 1px solid #333;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-md);
  padding-top: var(--spacing-md);
  border-top: 1px solid var(--border-color-light);
}

/* 滚动条样式 */
.email-content::-webkit-scrollbar,
.email-html-content::-webkit-scrollbar,
.email-text-content::-webkit-scrollbar,
.email-source-content::-webkit-scrollbar,
.email-headers-content::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

.email-content::-webkit-scrollbar-track,
.email-html-content::-webkit-scrollbar-track,
.email-text-content::-webkit-scrollbar-track,
.email-source-content::-webkit-scrollbar-track,
.email-headers-content::-webkit-scrollbar-track {
  background: var(--background-secondary);
}

.email-content::-webkit-scrollbar-thumb,
.email-html-content::-webkit-scrollbar-thumb,
.email-text-content::-webkit-scrollbar-thumb,
.email-source-content::-webkit-scrollbar-thumb,
.email-headers-content::-webkit-scrollbar-thumb {
  background: var(--border-color-dark);
  border-radius: var(--radius-small);
}

.email-content::-webkit-scrollbar-thumb:hover,
.email-html-content::-webkit-scrollbar-thumb:hover,
.email-text-content::-webkit-scrollbar-thumb:hover,
.email-source-content::-webkit-scrollbar-thumb:hover,
.email-headers-content::-webkit-scrollbar-thumb:hover {
  background: var(--text-tertiary);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .info-grid {
    grid-template-columns: 1fr;
    gap: var(--spacing-md);
  }
  
  .tab-buttons {
    flex-wrap: wrap;
    gap: var(--spacing-xs);
  }
  
  .tab-button {
    flex: 1;
    min-width: 100px;
    justify-content: center;
  }
  
  .email-info {
    padding: var(--spacing-md);
  }
  
  .email-html-content {
    padding: var(--spacing-md);
  }
  
  .email-text-content,
  .email-source-content,
  .email-headers-content {
    padding: var(--spacing-md);
    font-size: 0.75rem;
  }
}

@media (max-width: 480px) {
  .dialog-title {
    font-size: 1.25rem;
  }
  
  .email-id {
    font-size: 0.75rem;
  }
  
  .tab-button {
    padding: var(--spacing-xs) var(--spacing-sm);
    font-size: 0.75rem;
  }
  
  .info-grid {
    gap: var(--spacing-sm);
  }
  
  .email-info {
    padding: var(--spacing-sm);
  }
}

/* 深色模式代码块样式 */
.email-source-content {
  background: #0d1117;
  color: #c9d1d9;
  border-color: #30363d;
}

.email-headers-content {
  background: #161b22;
  color: #8b949e;
  border-color: #21262d;
}

/* 空内容样式 */
.empty-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-2xl);
  color: var(--text-secondary);
  background: var(--background-secondary);
  border-radius: 0 0 var(--radius-large) var(--radius-large);
  min-height: 350px;
  gap: var(--spacing-md);
}

.empty-content i {
  font-size: 2rem;
  opacity: 0.5;
}

.empty-content p {
  margin: 0;
  font-size: 0.9375rem;
  font-weight: 500;
}
</style>