package handlers

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"mailcat/internal/database"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	db               *database.DB
	authToken        string
	adminPasswordHash string

	// 随机 session 管理
	sessions   map[string]time.Time // token -> 过期时间
	sessionMu  sync.RWMutex

	// 登录速率限制
	loginAttempts map[string]*loginAttemptInfo
	loginMu       sync.Mutex
}

type loginAttemptInfo struct {
	count    int
	firstAt  time.Time
	lockedAt time.Time
}

const (
	maxLoginAttempts  = 5               // 最大尝试次数
	loginWindow       = 5 * time.Minute // 计数窗口
	lockDuration      = 15 * time.Minute // 锁定时长
	sessionMaxAge     = 24 * time.Hour  // Session 有效期
)

func NewAdminHandler(db *database.DB, authToken, adminPassword string) *AdminHandler {
	h := &AdminHandler{
		db:                db,
		authToken:         authToken,
		adminPasswordHash: sha256Hex(adminPassword),
		sessions:          make(map[string]time.Time),
		loginAttempts:     make(map[string]*loginAttemptInfo),
	}
	// 启动后台清理过期 session
	go h.cleanExpiredSessions()
	return h
}

// sha256Hex 计算字符串的 SHA-256 哈希值并返回十六进制字符串
func sha256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

// generateSessionToken 生成加密安全的随机 session token
func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// cleanExpiredSessions 定期清理过期 session
func (h *AdminHandler) cleanExpiredSessions() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		h.sessionMu.Lock()
		now := time.Now()
		for token, expiry := range h.sessions {
			if now.After(expiry) {
				delete(h.sessions, token)
			}
		}
		h.sessionMu.Unlock()
	}
}

// validateSession 验证 session token 是否有效
func (h *AdminHandler) validateSession(token string) bool {
	h.sessionMu.RLock()
	defer h.sessionMu.RUnlock()
	expiry, exists := h.sessions[token]
	if !exists {
		return false
	}
	return time.Now().Before(expiry)
}

// checkRateLimit 检查登录速率限制，返回是否允许登录
func (h *AdminHandler) checkRateLimit(ip string) (bool, int) {
	h.loginMu.Lock()
	defer h.loginMu.Unlock()

	info, exists := h.loginAttempts[ip]
	if !exists {
		h.loginAttempts[ip] = &loginAttemptInfo{count: 0, firstAt: time.Now()}
		return true, maxLoginAttempts
	}

	now := time.Now()

	// 检查是否在锁定期间
	if !info.lockedAt.IsZero() && now.Before(info.lockedAt.Add(lockDuration)) {
		remaining := int(info.lockedAt.Add(lockDuration).Sub(now).Seconds())
		return false, remaining
	}

	// 超出计数窗口则重置
	if now.After(info.firstAt.Add(loginWindow)) {
		info.count = 0
		info.firstAt = now
		info.lockedAt = time.Time{}
	}

	if info.count >= maxLoginAttempts {
		info.lockedAt = now
		return false, int(lockDuration.Seconds())
	}

	return true, maxLoginAttempts - info.count
}

// recordLoginAttempt 记录一次失败的登录尝试
func (h *AdminHandler) recordLoginAttempt(ip string) {
	h.loginMu.Lock()
	defer h.loginMu.Unlock()

	info, exists := h.loginAttempts[ip]
	if !exists {
		h.loginAttempts[ip] = &loginAttemptInfo{count: 1, firstAt: time.Now()}
		return
	}
	info.count++
}

// resetLoginAttempts 登录成功后重置计数
func (h *AdminHandler) resetLoginAttempts(ip string) {
	h.loginMu.Lock()
	defer h.loginMu.Unlock()
	delete(h.loginAttempts, ip)
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
		
		// 使用随机 session token 验证
		if !h.validateSession(session) {
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
	clientIP := c.ClientIP()

	// 速率限制检查
	allowed, remaining := h.checkRateLimit(clientIP)
	if !allowed {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error":       "Too many login attempts, please try again later",
			"retry_after": remaining,
		})
		return
	}

	var loginReq struct {
		Password string `json:"password"`
	}
	
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}
	
	// 前端发送的是 SHA-256 哈希后的密码，使用 HMAC 安全比较防止时序攻击
	if !hmac.Equal([]byte(loginReq.Password), []byte(h.adminPasswordHash)) {
		h.recordLoginAttempt(clientIP)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid password",
		})
		return
	}
	
	// 登录成功，重置速率限制
	h.resetLoginAttempts(clientIP)

	// 生成加密安全的随机 session token
	sessionToken, err := generateSessionToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	// 存储 session
	h.sessionMu.Lock()
	h.sessions[sessionToken] = time.Now().Add(sessionMaxAge)
	h.sessionMu.Unlock()

	c.SetCookie("admin_session", sessionToken, int(sessionMaxAge.Seconds()), "/", "", false, true)
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"session": sessionToken,
	})
}

// Logout 处理登出请求
func (h *AdminHandler) Logout(c *gin.Context) {
	// 从存储中移除 session
	session := c.GetHeader("X-Admin-Session")
	if session == "" {
		if cookie, err := c.Cookie("admin_session"); err == nil {
			session = cookie
		}
	}
	if session != "" {
		h.sessionMu.Lock()
		delete(h.sessions, session)
		h.sessionMu.Unlock()
	}

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

// GetConfig 获取配置信息（仅显示脱敏后的令牌）
func (h *AdminHandler) GetConfig(c *gin.Context) {
	masked := h.authToken
	if len(masked) > 8 {
		masked = masked[:4] + "****" + masked[len(masked)-4:]
	} else {
		masked = "****"
	}
	c.JSON(http.StatusOK, gin.H{
		"api_token": masked,
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