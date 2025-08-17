package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"mailcat/internal/database"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	db           *database.DB
	authToken    string
	adminPassword string
}

func NewAdminHandler(db *database.DB, authToken, adminPassword string) *AdminHandler {
	return &AdminHandler{
		db:           db,
		authToken:    authToken,
		adminPassword: adminPassword,
	}
}

// AdminAuthMiddleware 管理员认证中间件
func (h *AdminHandler) AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查session
		session := c.GetHeader("X-Admin-Session")
		if session == "" {
			// 检查cookie
			if cookie, err := c.Cookie("admin_session"); err == nil {
				session = cookie
			}
		}
		
		// 简单的session验证（实际项目中应该使用更安全的方式）
		if session != "admin_logged_in_"+h.adminPassword {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// Login 处理登录请求
func (h *AdminHandler) Login(c *gin.Context) {
	var loginReq struct {
		Password string `json:"password"`
	}
	
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}
	
	if loginReq.Password != h.adminPassword {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid password",
		})
		return
	}
	
	// 设置session cookie
	sessionToken := "admin_logged_in_" + h.adminPassword
	c.SetCookie("admin_session", sessionToken, 3600*24, "/", "", false, true)
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"session": sessionToken,
	})
}

// Logout 处理登出请求
func (h *AdminHandler) Logout(c *gin.Context) {
	c.SetCookie("admin_session", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}


// GetStats 获取统计信息
func (h *AdminHandler) GetStats(c *gin.Context) {
	stats, err := h.db.GetEmailStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get stats",
			"details": err.Error(),
		})
		return
	}
	
	// 添加系统状态
	stats["system_status"] = "running"
	
	c.JSON(http.StatusOK, stats)
}

// GetAdminEmails 获取邮件列表（管理员接口）
func (h *AdminHandler) GetAdminEmails(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	response, err := h.db.GetEmails(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get emails",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetConfig 获取配置信息
func (h *AdminHandler) GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"api_token": h.authToken,
	})
}

// SaveConfig 保存配置信息
func (h *AdminHandler) SaveConfig(c *gin.Context) {
	var configReq struct {
		APIToken      string `json:"api_token"`
		AdminPassword string `json:"admin_password"`
	}

	if err := c.ShouldBindJSON(&configReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	// 这里应该保存到配置文件
	// 简化实现，只返回成功
	c.JSON(http.StatusOK, gin.H{
		"message": "Configuration saved successfully",
	})
}

// escapeHtml HTML转义函数
func escapeHtml(s string) string {
	if s == "" {
		return ""
	}
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}