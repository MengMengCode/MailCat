<div align="center">

# 📧 MailCat

**现代化邮件接收与管理系统**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![Vue.js](https://img.shields.io/badge/Vue.js-3.x-4FC08D?style=flat-square&logo=vue.js)](https://vuejs.org/)
[![License](https://img.shields.io/badge/License-GPL%20v3-blue.svg?style=flat-square)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-ghcr.io-2496ED?style=flat-square&logo=docker)](https://github.com/MengMengCode/MailCat/pkgs/container/mailcat)
[![CI](https://img.shields.io/github/actions/workflow/status/MengMengCode/MailCat/docker-publish.yml?style=flat-square&label=Build)](https://github.com/MengMengCode/MailCat/actions)

</div>

---

## 📋 目录

- [✨ 功能特性](#-功能特性)
- [🔒 安全特性](#-安全特性)
- [🚀 快速开始](#-快速开始)
  - [Docker Compose 部署（推荐）](#1-docker-compose-部署推荐)
  - [Docker Run 部署](#2-docker-run-部署)
  - [源码构建运行](#3-源码构建运行)
- [☁️ Cloudflare Worker 配置](#️-cloudflare-worker-配置)
- [📡 API 使用说明](#-api-使用说明)
- [⚙️ 配置说明](#️-配置说明)
- [🔄 升级指南](#-升级指南)
- [🤝 贡献指南](#-贡献指南)
- [📄 许可证](#-许可证)

---

## ✨ 功能特性

MailCat 是一个基于 **Go + Vue.js** 的现代化邮件接收与管理系统，具有以下特性：

🔹 **轻量高效** - 基于 Go 语言开发，性能优异，资源占用低  
🔹 **现代化界面** - Vue.js 3 + PrimeVue 构建的响应式 Web 界面  
🔹 **云端集成** - 完美集成 Cloudflare Worker，实现邮件转发  
🔹 **数据持久化** - 使用 SQLite3 数据库，轻量且可靠  
🔹 **RESTful API** - 提供完整的 API 接口，支持第三方集成  
🔹 **容器化部署** - 支持 Docker 一键部署，镜像托管于 GitHub Container Registry  
🔹 **安全认证** - 双端哈希密码传输、随机 Session、速率限制  
🔹 **分页查询** - 支持大量邮件的分页浏览和管理  

---

## 🔒 安全特性

| 特性 | 说明 |
|------|------|
| **双端密码哈希** | 前端 SHA-256 哈希后传输，服务端 HMAC 安全比较，密码不明文传输 |
| **随机 Session Token** | 每次登录生成加密安全的随机 Token，不可预测 |
| **登录速率限制** | 5 次失败后锁定 15 分钟，防暴力破解 |
| **XSS 防护** | 邮件 HTML 使用 `sandbox` iframe 渲染，隔离恶意脚本 |
| **安全响应头** | 自动添加 `X-Content-Type-Options`、`X-Frame-Options`、`X-XSS-Protection` 等 |
| **请求体限制** | 10MB 请求体大小限制，防止 DoS 攻击 |
| **CORS 收紧** | 默认禁止跨域请求 |
| **API Token 脱敏** | 管理面板仅显示脱敏后的 Token |

---

## 🚀 快速开始

### 1. Docker Compose 部署（推荐）

创建 `docker-compose.yml` 文件：

```yaml
version: '3.8'

services:
  mailcat:
    image: ghcr.io/mengmengcode/mailcat:latest
    container_name: mailcat
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      # API 认证令牌 - 用于 Cloudflare Worker 调用 API（请修改为安全的随机字符串）
      - MAILCAT_API_AUTH_TOKEN=your_secure_api_token_here
      # 管理员密码 - 用于 Web 管理界面登录（请修改为强密码）
      - MAILCAT_ADMIN_PASSWORD=your_secure_admin_password_here
      # 时区设置
      - TZ=Asia/Shanghai
    volumes:
      # 数据持久化 - SQLite 数据库文件
      - mailcat_data:/app/data
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  mailcat_data:
    driver: local
```

启动服务：

```bash
docker compose up -d
```

✅ 服务启动后访问：**http://your-server-ip:8080**

---

### 2. Docker Run 部署

```bash
docker run -d \
  --name mailcat \
  --restart unless-stopped \
  -p 8080:8080 \
  -e MAILCAT_API_AUTH_TOKEN=your_secure_api_token_here \
  -e MAILCAT_ADMIN_PASSWORD=your_secure_admin_password_here \
  -e TZ=Asia/Shanghai \
  -v mailcat_data:/app/data \
  ghcr.io/mengmengcode/mailcat:latest
```

---

### 3. 源码构建运行

```bash
# 克隆项目
git clone https://github.com/MengMengCode/MailCat.git
cd MailCat

# 构建前端资源
cd web/frontend
npm install && npm run build
cd ../..

# 设置环境变量
export MAILCAT_API_AUTH_TOKEN=your_secure_api_token_here
export MAILCAT_ADMIN_PASSWORD=your_secure_admin_password_here

# 安装 Go 依赖并启动
go mod tidy
go run main.go
```

✅ 服务启动后访问：**http://localhost:8080**

## ☁️ Cloudflare Worker 配置

### 步骤 1：创建 Worker

1. 登录 [Cloudflare Dashboard](https://dash.cloudflare.com/)
2. 进入 **Workers & Pages** 页面
3. 点击 **创建应用程序** → **创建 Worker**

### 步骤 2：部署代码

复制 [`cloudflare-worker/worker.js`](cloudflare-worker/worker.js) 中的代码到 Worker 编辑器

### 步骤 3：配置环境变量

在 Worker 的 **设置** → **变量和机密** 中添加：

| 变量类型 | 变量名称 | 变量值 | 说明 |
|---------|---------|--------|------|
| **环境变量** | `API_ENDPOINT` | `https://your.domain.com` | MailCat 服务地址 |
| **机密** | `API_TOKEN` | `your_secure_api_token_here` | API 认证令牌 |

> ⚠️ **重要提醒**
> - `API_ENDPOINT` 必须使用完整域名（不支持 IP 或 localhost）
> - `API_TOKEN` 必须与 MailCat 服务配置保持一致
> - 强烈建议使用 HTTPS 确保数据传输安全

### 步骤 4：配置邮件路由

1. 在 Cloudflare Dashboard 中进入你的域名管理
2. 转到 **电子邮件** → **电子邮件路由**
3. 添加路由规则，将邮件转发到 Worker

### 步骤 5：测试连接

访问 Worker 域名，查看连接状态和健康检查结果。

---

## 📡 API 使用说明

### 基础信息

- **基础 URL**：`https://your.domain.com/api/v1`
- **认证方式**：`Authorization: Bearer <token>` 请求头
- **数据格式**：JSON

### 邮件查询接口

#### 端点地址
```
GET /api/v1/emails
```

#### 认证方式

```bash
curl -H "Authorization: Bearer your_auth_token" \
     https://your.domain.com/api/v1/emails
```

#### 查询参数

| 参数名 | 类型 | 默认值 | 范围 | 说明 |
|--------|------|--------|------|------|
| `page` | integer | `1` | ≥ 1 | 页码 |
| `limit` | integer | `20` | 1-100 | 每页数量 |

#### 使用示例

**默认查询（第1页，20条）**
```bash
curl -H "Authorization: Bearer your_auth_token" \
     "https://your.domain.com/api/v1/emails"
```

**分页查询（第2页，50条）**
```bash
curl -H "Authorization: Bearer your_auth_token" \
     "https://your.domain.com/api/v1/emails?page=2&limit=50"
```

**获取单封邮件详情**
```bash
curl -H "Authorization: Bearer your_auth_token" \
     "https://your.domain.com/api/v1/emails/1"
```

#### 响应示例

```json
{
  "emails": [
    {
      "id": 1,
      "from": "sender@example.com",
      "to": "recipient@yourdomain.com",
      "subject": "欢迎使用 MailCat",
      "body": "这是邮件的纯文本内容",
      "html_body": "<p>这是 <strong>HTML</strong> 格式的邮件内容</p>",
      "headers": "{\"Content-Type\":\"text/html; charset=utf-8\",\"Date\":\"Mon, 01 Jan 2025 12:00:00 +0000\"}",
      "received_at": "2025-01-01T12:00:00Z",
      "created_at": "2025-01-01T12:00:00Z"
    }
  ],
  "pagination": {
    "total": 150,
    "page": 1,
    "limit": 20,
    "total_pages": 8
  }
}
```

#### 错误响应

```json
{
  "error": "Unauthorized",
  "message": "Invalid or missing authentication token",
  "code": 401
}
```

---



## ⚙️ 配置说明

### 环境变量

| 变量名 | 必填 | 默认值 | 说明 |
|--------|:----:|--------|------|
| `MAILCAT_API_AUTH_TOKEN` | ✅ | - | API 认证令牌（Cloudflare Worker 使用） |
| `MAILCAT_ADMIN_PASSWORD` | ✅ | - | 管理员登录密码 |
| `MAILCAT_SERVER_PORT` | ❌ | `8080` | 服务监听端口 |
| `MAILCAT_SERVER_HOST` | ❌ | `0.0.0.0` | 服务监听地址 |
| `MAILCAT_DATABASE_PATH` | ❌ | `./data/emails.db` | SQLite 数据库文件路径 |
| `TZ` | ❌ | `UTC` | 时区设置，建议 `Asia/Shanghai` |

### 配置文件

项目也支持通过 [`config/config.yaml`](config/config.yaml) 进行配置（环境变量优先级更高）：

```yaml
server:
  port: "8080"
  host: "0.0.0.0"

database:
  path: "./data/emails.db"

api:
  auth_token: ""  # 建议通过环境变量 MAILCAT_API_AUTH_TOKEN 设置

admin:
  password: ""  # 建议通过环境变量 MAILCAT_ADMIN_PASSWORD 设置
```

> ⚠️ **安全提醒**：请勿将真实的 Token 和密码提交到版本控制中，推荐使用环境变量或 `.env` 文件。

---

## 🔄 升级指南

### 从旧版本升级

MailCat 保证**数据库格式向后兼容**，升级不会影响已有数据。

**Docker 用户：**

```bash
# 拉取最新镜像
docker pull ghcr.io/mengmengcode/mailcat:latest

# 重新创建容器（数据卷自动保留）
docker compose down
docker compose up -d
```

**源码用户：**

```bash
git pull origin main
cd web/frontend && npm install && npm run build && cd ../..
go mod tidy
go build -o mailcat .
```

### 注意事项

- ✅ 数据库 SQLite 文件完全兼容，无需迁移
- ✅ 现有 Cloudflare Worker 配置无需修改
- ⚠️ API 认证已移除 URL 参数传 Token 的方式（安全原因），请改用 `Authorization: Bearer <token>` 请求头

---

## 🤝 贡献指南

我们欢迎所有形式的贡献！

### 开发环境搭建

1. **克隆项目**
   ```bash
   git clone https://github.com/your-repo/mailcat.git
   cd mailcat
   ```

2. **后端开发**
   ```bash
   go mod tidy
   go run main.go
   ```

3. **前端开发**
   ```bash
   cd web/frontend
   npm install
   npm run dev
   ```

---

## 📄 许可证

本项目基于 [GNU General Public License v3.0](LICENSE) 开源协议发布。

---

<div align="center">

**⭐ 如果这个项目对您有帮助，请给我们一个 Star！**

Made with ❤️ by MailCat Team

</div>