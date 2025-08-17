# MailCat

基于Go + Vue.js的现代化邮件接收和管理系统，可以接收来自Cloudflare Worker的邮件数据并存储到SQLite3数据库中，提供美观的Web管理界面。

## 功能特性

- 🚀 接收来自Cloudflare Worker转发的邮件数据
- 💾 将邮件存储到本地SQLite3数据库
- 🌐 现代化Vue.js前端管理界面
- 📊 实时邮件统计和管理
- 🔍 邮件搜索和详情查看
- 🔐 安全的身份验证机制
- 📱 响应式设计，支持移动端
- ⚡ RESTful API支持
- 📄 分页查询支持
- 🏥 健康检查端点

## 项目结构

```
mailcat/
├── main.go                           # 主程序入口
├── go.mod                           # Go模块文件
├── go.sum                           # Go依赖锁定文件
├── README.md                        # 项目文档
├── .gitignore                       # Git忽略文件
├── config/
│   └── config.yaml                  # 配置文件
├── internal/                        # Go后端代码
│   ├── config/
│   │   └── config.go               # 配置管理
│   ├── models/
│   │   └── email.go                # 数据模型
│   ├── database/
│   │   └── database.go             # 数据库操作
│   ├── handlers/
│   │   ├── email.go                # 邮件API处理器
│   │   └── admin.go                # 管理员API处理器
│   ├── router/
│   │   └── router.go               # 路由设置
│   └── utils/
│       └── email_parser.go         # 邮件解析工具
├── web/                            # 前端资源
│   └── frontend/                   # Vue.js前端应用
│       ├── index.html              # 入口HTML
│       ├── package.json            # 前端依赖
│       ├── package-lock.json       # 前端依赖锁定
│       ├── vite.config.js          # Vite构建配置
│       └── src/                    # Vue源码
│           ├── main.js             # 前端入口
│           ├── App.vue             # 根组件
│           ├── router/
│           │   └── index.js        # 前端路由
│           ├── views/              # 页面组件
│           │   ├── Login.vue       # 登录页面
│           │   └── Dashboard.vue   # 管理面板
│           ├── components/         # 通用组件
│           │   └── EmailDetailDialog.vue  # 邮件详情对话框
│           └── services/
│               └── api.js          # API服务
├── cloudflare-worker/              # Cloudflare Worker代码
│   └── worker.js                   # Worker脚本
└── data/                          # 数据库文件目录（自动创建）
    └── emails.db                  # SQLite数据库
```

## 快速开始

### 1. 安装后端依赖

```bash
go mod tidy
```

### 2. 构建前端

```bash
cd web/frontend
npm install
npm run build
cd ../..
```

### 3. 配置

编辑 `config/config.yaml` 文件：

```yaml
server:
  port: "8080"
  host: "0.0.0.0"

database:
  path: "./data/emails.db"

api:
  auth_token: "your-secret-token-here"  # 请修改为安全的令牌

admin:
  password: "your-admin-password"       # 请修改为安全的密码
```

### 4. 运行服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

### 5. 访问管理界面

- 管理面板: `http://localhost:8080/admin/`
- 登录页面: `http://localhost:8080/admin/login`
- API健康检查: `http://localhost:8080/health`

默认管理员密码请在配置文件中设置。

### 6. 配置Cloudflare Worker

1. 在Cloudflare Dashboard中创建新的Worker
2. 复制 `cloudflare-worker/worker.js` 中的代码
3. 设置环境变量：
   
   在Worker设置页面的"变量和机密"部分添加：
   
   | 类型 | 名称 | 值 |
   |------|------|-----|
   | 纯文本 | `API_ENDPOINT` | `https://your-domain.com` |
   | 机密 | `API_TOKEN` | `your-secret-token-here` |
   
   **重要提示：**
   - 🌐 **`API_ENDPOINT` 必须使用域名，不能使用IP地址**（Cloudflare Worker限制）
   - ✅ 正确格式：`https://api.example.com` 或 `http://your-domain.com`
   - ❌ 错误格式：`http://192.168.1.100:8080` 或 `http://localhost:8080`
   - 🔒 `API_TOKEN` 建议设置为"机密"类型而不是"纯文本"，这样更安全
   - 🔑 `API_TOKEN` 的值必须与您的 `config/config.yaml` 中的 `auth_token` 完全一致
   - 🔐 推荐使用HTTPS以确保数据传输安全

4. 配置邮件路由规则
5. 部署后访问Worker域名查看连接状态

#### Worker状态页面

部署Worker后，可以通过以下方式检查连接状态：

- **状态页面**: 访问Worker域名根路径（如 `https://your-worker.your-subdomain.workers.dev/`）
  - 显示与Go API的连接状态
  - 显示配置信息
  - 提供实时状态检查

- **Health API**: 访问 `/health` 端点获取JSON格式的状态信息
  ```json
  {
    "status": "healthy",
    "message": "API服务器响应正常 (状态码: 200)",
    "timestamp": "2024-01-01T12:00:00.000Z",
    "api_endpoint": "https://your-domain.com",
    "token_configured": true
  }
  ```

## API端点

### 公开端点

#### 健康检查
```
GET /health
```

### 邮件API端点（需要API Token认证）

#### 接收邮件（由Cloudflare Worker调用）
```
POST /api/v1/emails
Authorization: Bearer <your-api-token>
Content-Type: application/json

{
  "from": "sender@example.com",
  "to": "recipient@example.com",
  "subject": "邮件主题",
  "body": "邮件正文",
  "html_body": "<html>HTML邮件内容</html>",
  "headers": {
    "header-name": "header-value"
  }
}
```

#### 获取邮件列表
```
GET /api/v1/emails?page=1&limit=20
Authorization: Bearer <your-api-token>
# 或者使用查询参数
GET /api/v1/emails?token=your-api-token&page=1&limit=20
```

#### 获取单个邮件
```
GET /api/v1/emails/{id}
Authorization: Bearer <your-api-token>
```

### 管理员Web端点

#### 前端应用
```
GET /admin/                    # 管理面板首页
GET /admin/login              # 登录页面
GET /admin/dashboard          # 仪表板页面
```

#### 管理员认证
```
POST /admin/login
Content-Type: application/json

{
  "password": "your-admin-password"
}
```

```
POST /admin/logout
```

### 管理员API端点（需要管理员Session认证）

#### 获取统计信息
```
GET /admin/api/stats
Cookie: admin_session=<session-token>
```

#### 获取邮件列表（管理员视图）
```
GET /admin/api/emails?page=1&limit=20
Cookie: admin_session=<session-token>
```

#### 获取单个邮件详情
```
GET /admin/api/emails/{id}
Cookie: admin_session=<session-token>
```

#### 获取配置信息
```
GET /admin/api/config
Cookie: admin_session=<session-token>
```

#### 保存配置信息
```
POST /admin/api/config
Cookie: admin_session=<session-token>
Content-Type: application/json

{
  "api_token": "new-api-token",
  "admin_password": "new-admin-password"
}
```

### API查询参数说明

#### 邮件列表查询参数：
- `page`: 页码（默认：1，最小值：1）
- `limit`: 每页数量（默认：20，范围：1-100）
- `token`: API认证令牌（仅用于API端点，也可使用Authorization头部）

#### 查询示例：
```bash
# 查询第1页，每页20条（默认）
GET /api/v1/emails?token=your-api-token

# 查询第2页，每页50条
GET /api/v1/emails?token=your-api-token&page=2&limit=50

# 获取最多100条邮件
GET /api/v1/emails?token=your-api-token&limit=100
```

## 响应格式

### 邮件列表响应
```json
{
  "emails": [
    {
      "id": 1,
      "from": "sender@example.com",
      "to": "recipient@example.com",
      "subject": "邮件主题",
      "body": "邮件正文",
      "html_body": "<html>HTML内容</html>",
      "headers": "{\"header\":\"value\"}",
      "received_at": "2024-01-01T12:00:00Z",
      "created_at": "2024-01-01T12:00:00Z"
    }
  ],
  "total": 100,
  "page": 1,
  "limit": 20
}
```

## 身份验证

### API认证（用于邮件API）

API使用Bearer Token进行身份验证。可以通过以下方式提供令牌：

1. **HTTP Header**: `Authorization: Bearer <your-api-token>`
2. **Query Parameter**: `?token=<your-api-token>`

API Token在 `config/config.yaml` 中的 `api.auth_token` 字段配置。

### 管理员认证（用于Web管理界面）

管理员使用Session Cookie进行身份验证：

1. **登录**: POST `/admin/login` 使用密码登录
2. **Session**: 登录成功后会设置 `admin_session` Cookie
3. **认证**: 后续请求会自动携带Cookie进行认证
4. **登出**: POST `/admin/logout` 清除Session

管理员密码在 `config/config.yaml` 中的 `admin.password` 字段配置。

### 安全建议

- 🔐 使用强密码作为API Token和管理员密码
- 🔒 在生产环境中使用HTTPS
- 🔑 定期更换认证凭据
- 🚫 不要在日志中记录敏感信息

## 数据库

使用SQLite3数据库存储邮件数据。数据库文件会自动创建在配置指定的路径。

### 邮件表结构
```sql
CREATE TABLE emails (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    from_address TEXT NOT NULL,
    to_address TEXT NOT NULL,
    subject TEXT,
    body TEXT,
    html_body TEXT,
    headers TEXT,
    received_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## 部署

### 使用systemd（Linux）

1. 创建服务文件 `/etc/systemd/system/mailcat.service`：

```ini
[Unit]
Description=MailCat Service
After=network.target

[Service]
Type=simple
User=your-user
WorkingDirectory=/path/to/mailcat
ExecStart=/path/to/mailcat/mailcat
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

2. 启用并启动服务：

```bash
sudo systemctl enable mailcat
sudo systemctl start mailcat
```

### 使用Docker

创建 `Dockerfile`：

```dockerfile
# 前端构建阶段
FROM node:18-alpine AS frontend-builder
WORKDIR /app/web/frontend
COPY web/frontend/package*.json ./
RUN npm ci
COPY web/frontend/ ./
RUN npm run build

# 后端构建阶段
FROM golang:1.21-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o mailcat main.go

# 最终运行阶段
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=backend-builder /app/mailcat .
COPY --from=backend-builder /app/config ./config
COPY --from=frontend-builder /app/web/frontend/dist ./web/dist
CMD ["./mailcat"]
```

构建并运行：

```bash
docker build -t mailcat .
docker run -p 8080:8080 -v $(pwd)/data:/root/data mailcat
```

### 使用Docker Compose

创建 `docker-compose.yml`：

```yaml
version: '3.8'
services:
  mailcat:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./data:/root/data
      - ./config:/root/config
    environment:
      - GIN_MODE=release
    restart: unless-stopped
```

运行：

```bash
docker-compose up -d
```

## 安全建议

1. 使用强密码作为API令牌
2. 在生产环境中使用HTTPS
3. 定期备份数据库文件
4. 限制API访问的IP地址
5. 监控日志文件

## 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查数据库文件路径权限
   - 确保目录存在
   - 验证SQLite3是否正确安装

2. **API认证失败**
   - 检查API令牌是否正确
   - 确认Authorization Header格式正确
   - 验证config.yaml中的auth_token配置

3. **管理员登录失败**
   - 检查管理员密码是否正确
   - 确认config.yaml中的admin.password配置
   - 清除浏览器Cookie后重试

4. **前端页面无法访问**
   - 确认前端已正确构建：`cd web/frontend && npm run build`
   - 检查web/dist目录是否存在
   - 验证Vite构建是否成功

5. **Cloudflare Worker无法连接**
   - 检查API端点地址（必须使用域名，不能使用IP）
   - 确认防火墙设置
   - 验证SSL证书
   - 检查Worker环境变量配置

6. **邮件内容显示异常**
   - 检查邮件编码格式
   - 验证MIME解析是否正确
   - 查看服务器日志获取详细错误信息

### 开发调试

1. **启用调试模式**
   ```bash
   GIN_MODE=debug go run main.go
   ```

2. **查看详细日志**
   - 检查控制台输出
   - 使用浏览器开发者工具查看网络请求
   - 检查API响应状态码和错误信息

3. **前端开发模式**
   ```bash
   cd web/frontend
   npm run dev
   ```
   然后访问 `http://localhost:5173` 进行前端开发调试。

## 许可证

MIT License