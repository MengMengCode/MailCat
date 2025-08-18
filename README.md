<div align="center">

# 📧 MailCat

**现代化邮件接收与管理系统**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![Vue.js](https://img.shields.io/badge/Vue.js-3.x-4FC08D?style=flat-square&logo=vue.js)](https://vuejs.org/)
[![License](https://img.shields.io/badge/License-GPL%20v3-blue.svg?style=flat-square)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker)](https://hub.docker.com/)

</div>

---

## 📋 目录

- [✨ 功能特性](#-功能特性)
- [🚀 快速开始](#-快速开始)
  - [源码构建运行](#1-源码构建运行)
  - [Docker Compose 部署](#2-docker-compose-部署)
  - [Docker Run 部署](#3-docker-run-部署)
- [☁️ Cloudflare Worker 配置](#️-cloudflare-worker-配置)
- [📡 API 使用说明](#-api-使用说明)
- [🖼️ 界面预览](#️-界面预览)
- [⚙️ 配置说明](#️-配置说明)
- [🤝 贡献指南](#-贡献指南)
- [📄 许可证](#-许可证)

---

## ✨ 功能特性

MailCat 是一个基于 **Go + Vue.js** 的现代化邮件接收与管理系统，具有以下特性：

🔹 **轻量高效** - 基于 Go 语言开发，性能优异，资源占用低  
🔹 **现代化界面** - Vue.js 3 + Element Plus 构建的响应式 Web 界面  
🔹 **云端集成** - 完美集成 Cloudflare Worker，实现邮件转发  
🔹 **数据持久化** - 使用 SQLite3 数据库，轻量且可靠  
🔹 **RESTful API** - 提供完整的 API 接口，支持第三方集成  
🔹 **容器化部署** - 支持 Docker 一键部署，开箱即用  
🔹 **安全认证** - 支持 Token 认证和管理员密码保护  
🔹 **分页查询** - 支持大量邮件的分页浏览和管理  

---

## 🚀 快速开始

### 1. 源码构建运行

```bash
# 克隆项目
git clone https://github.com/MengMengCode/MailCat.git
cd mailcat

# 安装 Go 依赖
go mod tidy

# 构建前端资源
cd web/frontend
npm install && npm run build
cd ../..

# 启动服务
go run main.go
```

✅ 服务启动后访问：**http://server-ip:8080**

---

### 2. Docker Compose 部署

创建 `docker-compose.yml` 文件：

```yaml
version: '3.8'

services:
  mailcat:
    image: mengmengcode/mailcat:latest
    container_name: mailcat
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - MAILCAT_API_AUTH_TOKEN=your_secure_api_token_here
      - MAILCAT_ADMIN_PASSWORD=your_secure_admin_password_here
    volumes:
      - mailcat_data:/app/data
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  mailcat_data:
    driver: local
```

启动服务：

```bash
docker-compose up -d
```

---

### 3. Docker Run 部署

```bash
docker run -d \
  --name mailcat \
  --restart unless-stopped \
  -p 8080:8080 \
  -e MAILCAT_API_AUTH_TOKEN=your_secure_api_token_here \
  -e MAILCAT_ADMIN_PASSWORD=your_secure_admin_password_here \
  -v mailcat_data:/app/data \
  mengmengcode/mailcat:latest
```

---

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
- **认证方式**：Bearer Token 或 URL 参数
- **数据格式**：JSON

### 邮件查询接口

#### 端点地址
```
GET /api/v1/emails
```

#### 认证方式

**方式一：URL 参数**
```
https://your.domain.com/api/v1/emails?token=your_auth_token
```

**方式二：请求头**
```bash
curl -H "Authorization: Bearer your_auth_token" \
     https://your.domain.com/api/v1/emails
```

#### 查询参数

| 参数名 | 类型 | 默认值 | 范围 | 说明 |
|--------|------|--------|------|------|
| `page` | integer | `1` | ≥ 1 | 页码 |
| `limit` | integer | `20` | 1-100 | 每页数量 |
| `token` | string | - | - | 认证令牌（可选，如使用请求头认证） |

#### 使用示例

**默认查询（第1页，20条）**
```bash
curl "https://your.domain.com/api/v1/emails?token=your_auth_token"
```

**分页查询（第2页，50条）**
```bash
curl "https://your.domain.com/api/v1/emails?token=your_auth_token&page=2&limit=50"
```

**获取所有邮件（分页遍历）**
```bash
# 第一次请求获取总数
curl "https://your.domain.com/api/v1/emails?token=your_auth_token&limit=100"

# 根据返回的 total 字段计算总页数，然后逐页查询
curl "https://your.domain.com/api/v1/emails?token=your_auth_token&page=2&limit=100"
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

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `MAILCAT_API_AUTH_TOKEN` | - | API 认证令牌（必填） |
| `MAILCAT_ADMIN_PASSWORD` | - | 管理员密码（必填） |
| `MAILCAT_PORT` | `8080` | 服务监听端口 |
| `MAILCAT_DB_PATH` | `./data/mailcat.db` | SQLite 数据库文件路径 |

### 配置文件

项目支持通过 [`config/config.yaml`](config/config.yaml) 进行配置：

```yaml
server:
  port: 8080
  host: "0.0.0.0"

database:
  path: "./data/mailcat.db"

auth:
  api_token: "your_secure_api_token_here"
  admin_password: "your_secure_admin_password_here"
```

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