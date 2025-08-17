# 多阶段构建 Dockerfile for MailCat

# 第一阶段：构建前端
FROM node:18-alpine AS frontend-builder

WORKDIR /app/frontend

# 复制前端依赖文件
COPY web/frontend/package*.json ./

# 安装前端依赖（包括构建工具）
RUN npm ci

# 复制前端源码
COPY web/frontend/ ./

# 构建前端
RUN npm run build

# 第二阶段：构建Go应用
FROM golang:1.21-alpine AS go-builder

# 安装必要的包，包括构建工具和SQLite开发包
RUN apk add --no-cache gcc musl-dev sqlite-dev linux-headers build-base

WORKDIR /app

# 复制Go模块文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制Go源码
COPY . .

# 从前端构建阶段复制构建好的静态文件
COPY --from=frontend-builder /app/dist ./web/dist

# 构建Go应用，使用兼容musl的编译选项
RUN CGO_ENABLED=1 GOOS=linux \
    go build -a -installsuffix cgo \
    -tags "sqlite_omit_load_extension" \
    -ldflags "-s -w" \
    -o mailcat .

# 第三阶段：运行时镜像
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates sqlite tzdata

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /app

# 创建非root用户
RUN addgroup -g 1001 -S mailcat && \
    adduser -S mailcat -u 1001 -G mailcat

# 复制构建好的应用
COPY --from=go-builder /app/mailcat .

# 复制配置文件模板
COPY --from=go-builder /app/config ./config

# 创建数据目录并设置权限
RUN mkdir -p /app/data && \
    chown -R mailcat:mailcat /app

# 切换到非root用户
USER mailcat

# 暴露端口
EXPOSE 8080

# 创建数据卷挂载点
VOLUME ["/app/data"]

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动应用
CMD ["./mailcat"]