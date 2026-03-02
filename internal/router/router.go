package router

import (
	"mailcat/internal/database"
	"mailcat/internal/handlers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter(db *database.DB, authToken string, adminPassword string) *gin.Engine {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)
	
	r := gin.Default()

	// 请求体大小限制（10MB，防止大邮件 DoS）
	r.MaxMultipartMemory = 10 << 20
	
	// 设置模板分隔符，避免与Vue.js语法冲突
	r.Delims("{[{", "}]}")
	
	// 静态文件服务 - 仅服务Vue构建后的资源
	r.Static("/assets", "./web/dist/assets")

	// 安全响应头中间件
	r.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	})
	
	// 同源部署，不启用 CORS（浏览器默认阻止跨域请求）

	// 请求体大小限制中间件
	r.Use(func(c *gin.Context) {
		if c.Request.ContentLength > 10*1024*1024 { // 10MB
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "Request body too large",
			})
			return
		}
		c.Next()
	})
	
	// 创建邮件处理器
	emailHandler := handlers.NewEmailHandler(db, authToken)
	
	// 创建管理员处理器
	adminHandler := handlers.NewAdminHandler(db, authToken, adminPassword)
	
	// 公开端点
	r.GET("/health", emailHandler.HealthCheck)
	
	// 根路径重定向到管理员登录页面
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/admin/login")
	})
	
	// API路由组
	api := r.Group("/api/v1")
	{
		// 邮件接收端点（需要认证）
		api.POST("/emails", emailHandler.AuthMiddleware(), emailHandler.ReceiveEmail)
		
		// 邮件读取端点（需要认证）
		api.GET("/emails", emailHandler.AuthMiddleware(), emailHandler.GetEmails)
		api.GET("/emails/:id", emailHandler.AuthMiddleware(), emailHandler.GetEmailByID)
	}
	
	// 管理员路由组
	admin := r.Group("/admin")
	{
		// 静态资源必须在其他路由之前定义
		admin.Static("/assets", "./web/dist/assets")
		
		// Vue 前端应用 - 服务构建后的静态文件
		admin.StaticFile("/", "./web/dist/index.html")
		admin.StaticFile("/login", "./web/dist/index.html")
		admin.StaticFile("/dashboard", "./web/dist/index.html")
		
		// API 路由
		admin.POST("/login", adminHandler.Login)
		admin.POST("/logout", adminHandler.Logout)
		
		// 管理员API路由
		adminAPI := admin.Group("/api")
		adminAPI.Use(adminHandler.AdminAuthMiddleware())
		{
			adminAPI.GET("/stats", adminHandler.GetStats)
			adminAPI.GET("/emails", adminHandler.GetAdminEmails)
			adminAPI.GET("/emails/:id", emailHandler.GetEmailByID)
			adminAPI.GET("/config", adminHandler.GetConfig)
			adminAPI.POST("/config", adminHandler.SaveConfig)
		}
		
	}
	
	return r
}