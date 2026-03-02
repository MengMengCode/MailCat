<template>
  <div class="dashboard-container">
    <!-- 顶部导航栏 -->
    <header class="main-header">
      <div class="header-left">
        <h1 class="app-title">📧 MailCat</h1>
      </div>
      <div class="header-right">
        <Button
          icon="pi pi-refresh"
          @click="refreshData"
          :loading="refreshing"
          class="p-button-outlined"
          size="small"
        />
        <Button
          label="退出登录"
          icon="pi pi-sign-out"
          @click="handleLogout"
          class="p-button-outlined"
          size="small"
        />
      </div>
    </header>

    <!-- 主内容区域 -->
    <main class="content-area">
      <!-- 统计卡片 -->
      <div class="stats-grid">
        <Card class="stat-card">
          <template #content>
            <div class="stat-content">
              <div class="stat-icon">
                <i class="pi pi-envelope"></i>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ stats.total_emails || 0 }}</div>
                <div class="stat-label">总邮件数</div>
              </div>
            </div>
          </template>
        </Card>
        
        <Card class="stat-card">
          <template #content>
            <div class="stat-content">
              <div class="stat-icon">
                <i class="pi pi-calendar"></i>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ stats.today_emails || 0 }}</div>
                <div class="stat-label">今日邮件</div>
              </div>
            </div>
          </template>
        </Card>
      </div>

      <!-- 最近7天收件曲线图 -->
      <Card class="chart-card">
        <template #header>
          <h3>最近7天收件统计</h3>
        </template>
        <template #content>
          <div class="chart-container">
            <canvas ref="chartCanvas" width="800" height="300"></canvas>
          </div>
        </template>
      </Card>

      <!-- 设置区域 -->
      <div class="settings-grid">
        <Card class="settings-card">
          <template #header>
            <h3>API 邮件查询设置</h3>
          </template>
          <template #content>
            <div class="settings-form">
              <div class="form-group">
                <label>API 认证令牌</label>
                <div class="input-with-copy">
                  <InputText
                    v-model="config.api_token"
                    placeholder="当前令牌"
                    readonly
                    class="readonly-input long-input"
                  />
                  <Button
                    icon="pi pi-copy"
                    @click="copyToClipboard(config.api_token)"
                    class="p-button-outlined copy-btn"
                    size="small"
                    v-tooltip="'复制'"
                  />
                </div>
              </div>
              <div class="form-group">
                <label>查询端点</label>
                <div class="input-with-copy">
                  <InputText
                    :value="apiEndpoints.query"
                    readonly
                    class="readonly-input long-input"
                  />
                  <Button
                    icon="pi pi-copy"
                    @click="copyToClipboard(apiEndpoints.query)"
                    class="p-button-outlined copy-btn"
                    size="small"
                    v-tooltip="'复制'"
                  />
                </div>
              </div>
              
              <!-- 分页查询说明 -->
              <div class="form-group">
                <label>分页查询说明</label>
                <div class="query-examples">
                  <div class="example-section">
                    <h4>1. 按页查询：</h4>
                    <div class="example-item">
                      <code>{{ apiEndpoints.query }}</code>
                      <span class="example-desc">查询第1页，每页20条（默认）</span>
                    </div>
                    <div class="example-item">
                      <code>{{ apiEndpoints.query }}&page=2</code>
                      <span class="example-desc">查询第2页，每页20条</span>
                    </div>
                    <div class="example-item">
                      <code>{{ apiEndpoints.query }}&page=3&limit=50</code>
                      <span class="example-desc">查询第3页，每页50条</span>
                    </div>
                  </div>
                  
                  <div class="example-section">
                    <h4>2. 查询全部邮件：</h4>
                    <div class="method-section">
                      <strong>方法一：设置较大的limit值</strong>
                      <div class="example-item">
                        <code>{{ apiEndpoints.query }}&limit=100</code>
                        <span class="example-desc">设置limit为100（系统最大限制）</span>
                      </div>
                    </div>
                    
                    <div class="method-section">
                      <strong>方法二：分页遍历所有邮件</strong>
                      <div class="example-item">
                        <code>{{ apiEndpoints.query }}&page=1&limit=100</code>
                        <span class="example-desc">先查询第1页获取总数</span>
                      </div>
                      <div class="example-item">
                        <code>{{ apiEndpoints.query }}&page=2&limit=100</code>
                        <span class="example-desc">根据返回的total字段计算总页数，然后逐页查询</span>
                      </div>
                    </div>
                  </div>
                  
                  <div class="example-section">
                    <h4>参数说明：</h4>
                    <ul class="param-list">
                      <li><code>page</code>: 页码（默认：1，最小值：1）</li>
                      <li><code>limit</code>: 每页数量（默认：20，范围：1-100）</li>
                      <li><code>token</code>: 认证令牌（也可使用 Authorization: Bearer token 头部）</li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>
          </template>
        </Card>

      </div>

      <!-- 邮件列表 -->
      <Card class="emails-card">
        <template #header>
          <div class="card-header">
            <h3>所有邮件</h3>
            <Button
              icon="pi pi-refresh"
              @click="loadEmails"
              :loading="loadingEmails"
              class="p-button-outlined"
              size="small"
            />
          </div>
        </template>
        <template #content>
          <DataTable
            :value="emails"
            :loading="loadingEmails"
            :paginator="true"
            :rows="20"
            :totalRecords="totalEmails"
            :lazy="true"
            @page="onPageChange"
            @row-click="viewEmailDetail"
            class="emails-table"
          >
            <Column field="id" header="ID" style="width: 80px">
              <template #body="{ data }">
                <Tag>{{ data.id }}</Tag>
              </template>
            </Column>
            <Column field="subject" header="主题">
              <template #body="{ data }">
                <span class="email-subject">{{ data.subject || '(无主题)' }}</span>
              </template>
            </Column>
            <Column field="from" header="发件人" style="width: 200px">
              <template #body="{ data }">
                <Tag severity="info">{{ data.from }}</Tag>
              </template>
            </Column>
            <Column field="to" header="收件人" style="width: 200px">
              <template #body="{ data }">
                <Tag severity="success">{{ data.to }}</Tag>
              </template>
            </Column>
            <Column field="created_at" header="接收时间" style="width: 180px">
              <template #body="{ data }">
                <span class="email-time">{{ formatTime(data.created_at) }}</span>
              </template>
            </Column>
            <Column header="操作" style="width: 100px">
              <template #body="{ data }">
                <Button
                  icon="pi pi-eye"
                  @click.stop="viewEmailDetail(data)"
                  class="p-button-outlined"
                  size="small"
                />
              </template>
            </Column>
          </DataTable>
        </template>
      </Card>
    </main>

    <!-- 邮件详情弹窗 -->
    <EmailDetailDialog
      :visible="showEmailDialog"
      :email-id="selectedEmailId"
      @hide="showEmailDialog = false"
    />
  </div>
</template>

<script>
import { ref, reactive, onMounted, computed, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import { authAPI, emailAPI, configAPI } from '../services/api'
import EmailDetailDialog from '../components/EmailDetailDialog.vue'

export default {
  name: 'Dashboard',
  components: {
    EmailDetailDialog
  },
  setup() {
    const router = useRouter()
    const toast = useToast()
    const confirm = useConfirm()

    const refreshing = ref(false)
    const chartCanvas = ref(null)
    
    // 统计数据
    const stats = reactive({
      total_emails: 0,
      today_emails: 0,
      weekly_stats: []
    })

    // 邮件数据
    const emails = ref([])
    const totalEmails = ref(0)
    const currentPage = ref(1)
    const loadingEmails = ref(false)

    // 配置数据
    const config = reactive({
      api_token: ''
    })
    

    // 邮件详情弹窗
    const showEmailDialog = ref(false)
    const selectedEmailId = ref(null)

    // API端点信息
    const apiEndpoints = computed(() => ({
      query: `${window.location.origin}/api/v1/emails?token=${config.api_token}`,
      receive: `${window.location.origin}/api/v1/emails`
    }))


    const formatTime = (timeStr) => {
      if (!timeStr) return '未知时间'
      const date = new Date(timeStr)
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      })
    }

    // 绘制图表
    const drawChart = () => {
      if (!chartCanvas.value) {
        console.log('Canvas not available')
        return
      }

      const canvas = chartCanvas.value
      const ctx = canvas.getContext('2d')
      const { width, height } = canvas

      // 清空画布
      ctx.clearRect(0, 0, width, height)

      // 检查数据
      if (!stats.weekly_stats || stats.weekly_stats.length === 0) {
        // 显示无数据提示
        ctx.fillStyle = '#666'
        ctx.font = '16px -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto'
        ctx.textAlign = 'center'
        ctx.fillText('暂无数据', width / 2, height / 2)
        return
      }

      console.log('Drawing chart with data:', stats.weekly_stats)

      // 设置样式
      ctx.strokeStyle = '#0070f3'
      ctx.fillStyle = 'rgba(0, 112, 243, 0.1)'
      ctx.lineWidth = 3

      // 计算数据点
      const padding = 60
      const chartWidth = width - padding * 2
      const chartHeight = height - padding * 2

      const maxCount = Math.max(...stats.weekly_stats.map(item => item.count), 1)
      const points = stats.weekly_stats.map((item, index) => ({
        x: padding + (stats.weekly_stats.length > 1 ? (index / (stats.weekly_stats.length - 1)) * chartWidth : chartWidth / 2),
        y: padding + chartHeight - (item.count / maxCount) * chartHeight,
        count: item.count,
        date: item.date
      }))

      // 绘制网格线
      ctx.strokeStyle = '#f0f0f0'
      ctx.lineWidth = 1
      for (let i = 0; i <= 5; i++) {
        const y = padding + (i / 5) * chartHeight
        ctx.beginPath()
        ctx.moveTo(padding, y)
        ctx.lineTo(width - padding, y)
        ctx.stroke()
      }

      // 绘制填充区域
      if (points.length > 1) {
        ctx.fillStyle = 'rgba(0, 112, 243, 0.1)'
        ctx.beginPath()
        ctx.moveTo(points[0].x, height - padding)
        points.forEach(point => ctx.lineTo(point.x, point.y))
        ctx.lineTo(points[points.length - 1].x, height - padding)
        ctx.closePath()
        ctx.fill()
      }

      // 绘制曲线
      ctx.strokeStyle = '#0070f3'
      ctx.lineWidth = 3
      ctx.beginPath()
      if (points.length > 0) {
        ctx.moveTo(points[0].x, points[0].y)
        if (points.length > 1) {
          for (let i = 1; i < points.length; i++) {
            const prevPoint = points[i - 1]
            const currentPoint = points[i]
            const cpx = (prevPoint.x + currentPoint.x) / 2
            ctx.quadraticCurveTo(cpx, prevPoint.y, currentPoint.x, currentPoint.y)
          }
        }
      }
      ctx.stroke()

      // 绘制数据点
      points.forEach(point => {
        ctx.fillStyle = '#0070f3'
        ctx.beginPath()
        ctx.arc(point.x, point.y, 4, 0, Math.PI * 2)
        ctx.fill()
        
        ctx.fillStyle = '#fff'
        ctx.beginPath()
        ctx.arc(point.x, point.y, 2, 0, Math.PI * 2)
        ctx.fill()
      })

      // 绘制标签
      ctx.fillStyle = '#666'
      ctx.font = '12px -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto'
      ctx.textAlign = 'center'
      
      // Y轴标签 - 优化逻辑，避免重复标签
      ctx.textAlign = 'right'
      
      // 根据最大值动态确定合适的步数和标签值
      let yLabels = []
      if (maxCount <= 5) {
        // 当最大值较小时，使用整数步长
        for (let i = 0; i <= maxCount; i++) {
          yLabels.push(i)
        }
      } else {
        // 当最大值较大时，使用5个均匀分布的标签
        const step = maxCount / 5
        for (let i = 0; i <= 5; i++) {
          yLabels.push(Math.round(step * i))
        }
      }
      
      // 绘制Y轴标签 - 从上到下显示（最大值在顶部，0在底部）
      yLabels.forEach((value, index) => {
        const y = padding + ((yLabels.length - 1 - index) / (yLabels.length - 1)) * chartHeight
        ctx.fillText(value.toString(), padding - 10, y + 4)
      })
      
      // X轴标签
      ctx.textAlign = 'center'
      points.forEach(point => {
        const date = new Date(point.date + 'T00:00:00') // 确保正确解析日期
        const label = `${date.getMonth() + 1}/${date.getDate()}`
        ctx.fillText(label, point.x, height - padding + 20)
      })
    }

    const loadStats = async () => {
      try {
        const response = await emailAPI.getStats()
        Object.assign(stats, response.data)
        
        // 等待DOM更新后绘制图表
        await nextTick()
        drawChart()
      } catch (error) {
        console.error('加载统计失败:', error)
      }
    }

    const loadEmails = async (page = 1) => {
      loadingEmails.value = true
      try {
        const response = await emailAPI.getEmails(page, 20)
        emails.value = response.data.emails || []
        totalEmails.value = response.data.total || 0
        currentPage.value = page
      } catch (error) {
        console.error('加载邮件失败:', error)
        toast.add({
          severity: 'error',
          summary: '加载失败',
          detail: '加载邮件列表失败',
          life: 3000
        })
      } finally {
        loadingEmails.value = false
      }
    }

    const loadConfig = async () => {
      try {
        const response = await configAPI.getConfig()
        config.api_token = response.data.api_token || ''
      } catch (error) {
        console.error('加载配置失败:', error)
      }
    }

    const copyToClipboard = async (text) => {
      try {
        await navigator.clipboard.writeText(text)
        toast.add({
          severity: 'success',
          summary: '复制成功',
          detail: '内容已复制到剪贴板',
          life: 2000
        })
      } catch (error) {
        // 降级方案：使用传统的复制方法
        const textArea = document.createElement('textarea')
        textArea.value = text
        document.body.appendChild(textArea)
        textArea.select()
        try {
          document.execCommand('copy')
          toast.add({
            severity: 'success',
            summary: '复制成功',
            detail: '内容已复制到剪贴板',
            life: 2000
          })
        } catch (fallbackError) {
          toast.add({
            severity: 'error',
            summary: '复制失败',
            detail: '无法复制到剪贴板',
            life: 3000
          })
        }
        document.body.removeChild(textArea)
      }
    }


    const refreshData = async () => {
      refreshing.value = true
      try {
        await Promise.all([
          loadStats(),
          loadEmails(currentPage.value),
          loadConfig()
        ])
        toast.add({
          severity: 'success',
          summary: '✅ 刷新成功',
          detail: 'MailCat 数据已更新',
          life: 2000
        })
      } catch (error) {
        toast.add({
          severity: 'error',
          summary: '刷新失败',
          detail: '数据刷新失败',
          life: 3000
        })
      } finally {
        refreshing.value = false
      }
    }

    const onPageChange = (event) => {
      loadEmails(event.page + 1)
    }

    const viewEmailDetail = (emailOrEvent) => {
      // @row-click 传递 { originalEvent, data, index }，按钮直接传递 data
      const email = emailOrEvent.data || emailOrEvent
      selectedEmailId.value = email.id
      showEmailDialog.value = true
    }

    const handleLogout = () => {
      confirm.require({
        message: '确定要退出登录吗？',
        header: '确认退出',
        icon: 'pi pi-exclamation-triangle',
        accept: async () => {
          try {
            await authAPI.logout()
          } catch (error) {
            console.error('退出登录失败:', error)
          }
          localStorage.removeItem('admin_session')
          router.push('/login')
        }
      })
    }

    onMounted(() => {
      loadStats()
      loadEmails()
      loadConfig()
    })

    return {
      refreshing,
      chartCanvas,
      stats,
      emails,
      totalEmails,
      currentPage,
      loadingEmails,
      config,
      showEmailDialog,
      selectedEmailId,
      apiEndpoints,
      formatTime,
      loadEmails,
      refreshData,
      onPageChange,
      viewEmailDetail,
      handleLogout,
      copyToClipboard
    }
  }
}
</script>

<style scoped>
.dashboard-container {
  min-height: 100vh;
  background: var(--background-secondary);
}

/* 顶部导航栏 */
.main-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-lg) var(--spacing-xl);
  border-bottom: 1px solid var(--border-color);
  background: var(--surface);
  box-shadow: var(--shadow-small);
  position: sticky;
  top: 0;
  z-index: 100;
}

.app-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
  letter-spacing: -0.025em;
}

.header-right {
  display: flex;
  gap: var(--spacing-md);
}

/* 主内容区域 */
.content-area {
  padding: var(--spacing-xl);
  max-width: 1400px;
  margin: 0 auto;
}

/* 统计卡片 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: var(--spacing-lg);
  margin-bottom: var(--spacing-xl);
}

.stat-card {
  background: var(--surface);
  border: 1px solid var(--border-color);
  transition: all 0.15s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-medium);
  border-color: var(--border-color-dark);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.stat-icon {
  width: 56px;
  height: 56px;
  background: var(--background-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-large);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-primary);
  font-size: 1.375rem;
  transition: all 0.15s ease;
}

.stat-card:hover .stat-icon {
  background: var(--primary-color);
  color: var(--text-inverse);
  border-color: var(--primary-color);
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 2rem;
  font-weight: 800;
  color: var(--text-primary);
  line-height: 1;
  margin-bottom: var(--spacing-xs);
  font-family: var(--font-mono);
}

.stat-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

/* 图表卡片 */
.chart-card {
  margin-bottom: var(--spacing-xl);
  background: var(--surface);
  border: 1px solid var(--border-color);
}

.chart-card h3 {
  margin: 0;
  color: var(--text-primary);
  font-size: 1.25rem;
  font-weight: 700;
}

.chart-container {
  width: 100%;
  height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chart-container canvas {
  max-width: 100%;
  height: auto;
}

/* 设置网格 */
.settings-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: var(--spacing-lg);
  margin-bottom: var(--spacing-xl);
}

.settings-card {
  background: var(--surface);
  border: 1px solid var(--border-color);
}

.settings-card h3 {
  margin: 0;
  color: var(--text-primary);
  font-size: 1.25rem;
  font-weight: 700;
}

.settings-form {
  max-width: 100%;
}

.form-group {
  margin-bottom: var(--spacing-lg);
}

.form-group label {
  display: block;
  margin-bottom: var(--spacing-sm);
  color: var(--text-primary);
  font-size: 0.875rem;
  font-weight: 600;
  letter-spacing: 0.025em;
}

.readonly-input {
  background: var(--background-secondary) !important;
  color: var(--text-secondary) !important;
  cursor: not-allowed;
}

.input-with-copy {
  display: flex;
  gap: var(--spacing-sm);
  align-items: center;
}

.long-input {
  flex: 1;
  min-width: 0;
  font-family: var(--font-mono);
  font-size: 0.8125rem;
}

.copy-btn {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-medium);
  transition: all 0.15s ease;
}

.copy-btn:hover {
  background: var(--surface-hover);
  transform: translateY(-1px);
  box-shadow: var(--shadow-small);
}

.copy-btn:active {
  transform: translateY(0);
}

.save-btn {
  margin-top: var(--spacing-md);
}

/* 邮件卡片 */
.emails-card {
  background: var(--surface);
  border: 1px solid var(--border-color);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h3 {
  margin: 0;
  color: var(--text-primary);
  font-size: 1.25rem;
  font-weight: 700;
  letter-spacing: -0.025em;
}

/* 表格样式 */
.emails-table {
  background: transparent;
}

.emails-table :deep(.p-datatable-tbody > tr:hover) {
  background: var(--surface-hover);
  cursor: pointer;
}

.email-subject {
  font-weight: 600;
  color: var(--text-primary);
  cursor: pointer;
}

.email-subject:hover {
  color: var(--accent-color);
}

.email-time {
  color: var(--text-secondary);
  font-size: 0.8125rem;
  font-family: var(--font-mono);
  font-weight: 500;
}

/* 动画效果 */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.stat-card,
.chart-card,
.settings-card,
.emails-card {
  animation: fadeIn 0.3s ease-out;
}

.stat-card:nth-child(1) { animation-delay: 0.1s; }
.stat-card:nth-child(2) { animation-delay: 0.2s; }

/* 响应式设计 */
@media (max-width: 1200px) {
  .settings-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .content-area {
    padding: var(--spacing-md);
  }
  
  .main-header {
    padding: var(--spacing-md) var(--spacing-lg);
    flex-direction: column;
    gap: var(--spacing-md);
  }
  
  .header-right {
    width: 100%;
    justify-content: center;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
    gap: var(--spacing-md);
  }
  
  .settings-grid {
    grid-template-columns: 1fr;
    gap: var(--spacing-md);
  }
  
  .app-title {
    font-size: 1.5rem;
    text-align: center;
  }
}

@media (max-width: 480px) {
  .content-area {
    padding: var(--spacing-sm);
  }
  
  .main-header {
    padding: var(--spacing-sm) var(--spacing-md);
  }
  
  .app-title {
    font-size: 1.375rem;
  }
  
  .stat-value {
    font-size: 1.75rem;
  }
  
  .stat-icon {
    width: 48px;
    height: 48px;
    font-size: 1.25rem;
  }
  
  .chart-container {
    height: 250px;
  }
}

/* 查询示例样式 */
.query-examples {
  background: var(--background-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-medium);
  padding: var(--spacing-lg);
  margin-top: var(--spacing-sm);
}

.example-section {
  margin-bottom: var(--spacing-lg);
}

.example-section:last-child {
  margin-bottom: 0;
}

.example-section h4 {
  color: var(--text-primary);
  font-size: 1rem;
  font-weight: 600;
  margin: 0 0 var(--spacing-md) 0;
}

.method-section {
  margin-bottom: var(--spacing-md);
}

.method-section:last-child {
  margin-bottom: 0;
}

.method-section strong {
  color: var(--text-primary);
  font-size: 0.875rem;
  font-weight: 600;
  display: block;
  margin-bottom: var(--spacing-sm);
}

.example-item {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
  margin-bottom: var(--spacing-sm);
  padding: var(--spacing-sm);
  background: var(--surface);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-small);
}

.example-item:last-child {
  margin-bottom: 0;
}

.example-item code {
  font-family: var(--font-mono);
  font-size: 0.75rem;
  color: var(--primary-color);
  background: var(--background-secondary);
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--radius-small);
  border: 1px solid var(--border-color);
  word-break: break-all;
  line-height: 1.4;
}

.example-desc {
  font-size: 0.8125rem;
  color: var(--text-secondary);
  font-style: italic;
}

.param-list {
  margin: 0;
  padding-left: var(--spacing-lg);
}

.param-list li {
  margin-bottom: var(--spacing-xs);
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.param-list li:last-child {
  margin-bottom: 0;
}

.param-list code {
  font-family: var(--font-mono);
  font-size: 0.8125rem;
  color: var(--primary-color);
  background: var(--background-secondary);
  padding: 2px var(--spacing-xs);
  border-radius: var(--radius-small);
  border: 1px solid var(--border-color);
}

/* 响应式调整 */
@media (max-width: 768px) {
  .query-examples {
    padding: var(--spacing-md);
  }
  
  .example-item code {
    font-size: 0.6875rem;
  }
  
  .example-desc {
    font-size: 0.75rem;
  }
}
</style>