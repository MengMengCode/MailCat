import axios from 'axios'

// 创建 axios 实例
const api = axios.create({
  baseURL: '/',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('admin_session')
    if (token) {
      config.headers['X-Admin-Session'] = token
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('admin_session')
      window.location.href = '/admin/login'
    }
    return Promise.reject(error)
  }
)

// API 方法
export const authAPI = {
  // 登录
  login: (password) => {
    return api.post('/admin/login', { password })
  },
  
  // 登出
  logout: () => {
    return api.post('/admin/logout')
  }
}

export const emailAPI = {
  // 获取邮件列表
  getEmails: (page = 1, limit = 20) => {
    return api.get(`/admin/api/emails?page=${page}&limit=${limit}`)
  },
  
  // 获取邮件详情
  getEmailById: (id) => {
    return api.get(`/admin/api/emails/${id}`)
  },
  
  // 获取统计信息
  getStats: () => {
    return api.get('/admin/api/stats')
  }
}

export const configAPI = {
  // 获取配置
  getConfig: () => {
    return api.get('/admin/api/config')
  },
  
  // 保存配置
  saveConfig: (config) => {
    return api.post('/admin/api/config', config)
  }
}

export default api