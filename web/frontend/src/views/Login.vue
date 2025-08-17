<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <h2 class="login-subtitle">MailCat</h2>
        <p class="login-description">管理员登录</p>
      </div>
      
      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label for="password" class="form-label">密码</label>
          <InputText
            id="password"
            v-model="password"
            type="password"
            placeholder="请输入管理员密码"
            class="form-input"
            :class="{ 'error': error }"
            @keyup.enter="handleLogin"
          />
        </div>
        
        <div v-if="error" class="error-message">
          {{ error }}
        </div>
        
        <Button
          type="submit"
          :loading="loading"
          :disabled="!password || loading"
          class="login-button"
        >
          {{ loading ? '登录中...' : '登录' }}
        </Button>
      </form>
      
      <div class="login-footer">
        <p class="footer-text">Powered by MengMeng</p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { authAPI } from '../services/api'

export default {
  name: 'Login',
  setup() {
    const router = useRouter()
    const toast = useToast()
    
    const password = ref('')
    const loading = ref(false)
    const error = ref('')
    
    const handleLogin = async () => {
      if (!password.value.trim()) {
        error.value = '请输入密码'
        return
      }
      
      loading.value = true
      error.value = ''
      
      try {
        const response = await authAPI.login(password.value)
        
        if (response.data.session) {
          localStorage.setItem('admin_session', response.data.session)
          toast.add({
            severity: 'success',
            summary: '登录成功',
            detail: '欢迎回来！',
            life: 3000
          })
          router.push('/dashboard')
        }
      } catch (err) {
        error.value = err.response?.data?.error || '登录失败，请检查密码'
        toast.add({
          severity: 'error',
          summary: '登录失败',
          detail: error.value,
          life: 3000
        })
      } finally {
        loading.value = false
      }
    }
    
    return {
      password,
      loading,
      error,
      handleLogin
    }
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--background-secondary);
  padding: var(--spacing-lg);
  position: relative;
}

.login-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background:
    radial-gradient(circle at 25% 25%, rgba(0, 0, 0, 0.02) 0%, transparent 50%),
    radial-gradient(circle at 75% 75%, rgba(0, 0, 0, 0.02) 0%, transparent 50%);
  pointer-events: none;
}

.login-card {
  width: 100%;
  max-width: 420px;
  background: var(--surface);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-xl);
  padding: var(--spacing-2xl);
  box-shadow: var(--shadow-large);
  position: relative;
  z-index: 1;
  backdrop-filter: blur(10px);
}

.login-header {
  text-align: center;
  margin-bottom: var(--spacing-xl);
}

.login-title {
  font-size: 3.5rem;
  margin-bottom: var(--spacing-sm);
  background: linear-gradient(135deg, var(--primary-color) 0%, #333 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.login-subtitle {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: var(--spacing-sm);
  letter-spacing: -0.025em;
}

.login-description {
  color: var(--text-secondary);
  font-size: 0.9375rem;
  font-weight: 400;
  margin: 0;
}

.login-form {
  margin-bottom: var(--spacing-lg);
}

.form-group {
  margin-bottom: var(--spacing-lg);
}

.form-label {
  display: block;
  margin-bottom: var(--spacing-sm);
  color: var(--text-primary);
  font-size: 0.875rem;
  font-weight: 600;
  letter-spacing: 0.025em;
}

.form-input {
  width: 100%;
  padding: 0.875rem 1rem;
  background: var(--surface);
  border: 1.5px solid var(--border-color);
  border-radius: var(--radius-medium);
  color: var(--text-primary);
  font-size: 1rem;
  font-family: var(--font-sans);
  transition: all 0.15s ease;
  box-shadow: var(--shadow-small);
}

.form-input:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(0, 0, 0, 0.1);
  transform: translateY(-1px);
}

.form-input::placeholder {
  color: var(--text-tertiary);
}

.form-input.error {
  border-color: var(--error-color);
  box-shadow: 0 0 0 3px rgba(238, 0, 0, 0.1);
}

.error-message {
  color: var(--error-color);
  font-size: 0.875rem;
  font-weight: 500;
  margin-bottom: var(--spacing-md);
  text-align: center;
  padding: var(--spacing-sm);
  background: rgba(238, 0, 0, 0.05);
  border-radius: var(--radius-small);
  border: 1px solid rgba(238, 0, 0, 0.1);
}

.login-button {
  width: 100%;
  padding: 0.875rem 1rem;
  background: var(--primary-color);
  color: var(--text-inverse);
  border: none;
  border-radius: var(--radius-medium);
  font-size: 1rem;
  font-weight: 600;
  font-family: var(--font-sans);
  cursor: pointer;
  transition: all 0.15s ease;
  box-shadow: var(--shadow-small);
  position: relative;
  overflow: hidden;
}

.login-button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.1), transparent);
  transition: left 0.5s ease;
}

.login-button:hover:not(:disabled) {
  background: #333;
  transform: translateY(-2px);
  box-shadow: var(--shadow-medium);
}

.login-button:hover:not(:disabled)::before {
  left: 100%;
}

.login-button:active:not(:disabled) {
  transform: translateY(-1px);
}

.login-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
  box-shadow: var(--shadow-small);
}

.login-footer {
  text-align: center;
  padding-top: var(--spacing-lg);
  border-top: 1px solid var(--border-color-light);
  margin-top: var(--spacing-lg);
}

.footer-text {
  color: var(--text-tertiary);
  font-size: 0.8125rem;
  font-weight: 400;
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 480px) {
  .login-container {
    padding: var(--spacing-md);
  }
  
  .login-card {
    padding: var(--spacing-xl);
    margin: 0;
    border-radius: var(--radius-large);
  }
  
  .login-title {
    font-size: 2.75rem;
  }
  
  .login-subtitle {
    font-size: 1.5rem;
  }
  
  .form-input {
    padding: 0.75rem;
    font-size: 1rem;
  }
  
  .login-button {
    padding: 0.75rem;
    font-size: 1rem;
  }
}

@media (max-width: 360px) {
  .login-card {
    padding: var(--spacing-lg);
  }
  
  .login-title {
    font-size: 2.5rem;
  }
  
  .login-subtitle {
    font-size: 1.375rem;
  }
}

/* 动画效果 */
@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.login-card {
  animation: slideUp 0.4s ease-out;
}

/* 深色模式支持 */
@media (prefers-color-scheme: dark) {
  .login-container::before {
    background:
      radial-gradient(circle at 25% 25%, rgba(255, 255, 255, 0.02) 0%, transparent 50%),
      radial-gradient(circle at 75% 75%, rgba(255, 255, 255, 0.02) 0%, transparent 50%);
  }
}
</style>