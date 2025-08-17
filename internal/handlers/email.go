package handlers

import (
	"encoding/base64"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"mailcat/internal/database"
	"mailcat/internal/models"
	"mailcat/internal/utils"
	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	db        *database.DB
	authToken string
}

func NewEmailHandler(db *database.DB, authToken string) *EmailHandler {
	return &EmailHandler{
		db:        db,
		authToken: authToken,
	}
}

// AuthMiddleware 验证API令牌
func (h *EmailHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("token")
		}

		if token != "Bearer "+h.authToken && token != h.authToken {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ReceiveEmail 接收来自Cloudflare Worker的邮件
func (h *EmailHandler) ReceiveEmail(c *gin.Context) {
	var emailReq models.EmailRequest
	if err := c.ShouldBindJSON(&emailReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	email, err := h.db.SaveEmail(&emailReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save email",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Email received successfully",
		"email":   email,
	})
}

// GetEmails 获取邮件列表
func (h *EmailHandler) GetEmails(c *gin.Context) {
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
			"details": err.Error(),
		})
		return
	}

	// 清理返回字段：只返回核心字段（发件人，收件人，收件时间，主题，内容）
	optimizedEmails := make([]gin.H, len(response.Emails))
	for i, email := range response.Emails {
		// 解析和清理邮件内容
		body := email.Body
		htmlBody := email.HTMLBody
		
		// 首先检查是否是纯Base64编码的内容
		if body != "" && isBase64Content(body) {
			if decoded, err := base64.StdEncoding.DecodeString(body); err == nil {
				body = string(decoded)
			}
		}
		
		// 优先尝试解析MIME格式内容（包含base64解码）
		if body != "" && isMIMEContent(body) {
			if parsedContent, parseErr := utils.ParseMIMEContent(body); parseErr == nil {
				if parsedContent.TextBody != "" {
					body = parsedContent.TextBody
				}
				if parsedContent.HTMLBody != "" {
					htmlBody = parsedContent.HTMLBody
				}
			}
		} else if body != "" && (strings.Contains(body, "Content-Type:") || strings.Contains(body, "boundary=")) {
			// 如果body包含multipart数据，尝试标准解析
			if parsedContent, parseErr := utils.ParseEmailFromRaw(body); parseErr == nil {
				if parsedContent.TextBody != "" {
					body = parsedContent.TextBody
				}
				if parsedContent.HTMLBody != "" {
					htmlBody = parsedContent.HTMLBody
				}
			}
		}
		
		// 如果还是空的，尝试从raw_email解析
		if (body == "" || htmlBody == "") && email.RawEmail != "" {
			rawEmail := email.RawEmail
			
			// 检查raw_email是否是Base64编码
			if isBase64Content(rawEmail) {
				if decoded, err := base64.StdEncoding.DecodeString(rawEmail); err == nil {
					rawEmail = string(decoded)
				}
			}
			
			if isMIMEContent(rawEmail) {
				if parsedContent, parseErr := utils.ParseMIMEContent(rawEmail); parseErr == nil {
					if body == "" && parsedContent.TextBody != "" {
						body = parsedContent.TextBody
					}
					if htmlBody == "" && parsedContent.HTMLBody != "" {
						htmlBody = parsedContent.HTMLBody
					}
				}
			} else if parsedContent, parseErr := utils.ParseEmailFromRaw(rawEmail); parseErr == nil {
				if body == "" && parsedContent.TextBody != "" {
					body = parsedContent.TextBody
				}
				if htmlBody == "" && parsedContent.HTMLBody != "" {
					htmlBody = parsedContent.HTMLBody
				}
			}
		}

		// 清理发件人和收件人字段，移除多余的格式
		from := cleanEmailAddress(email.From)
		to := cleanEmailAddress(email.To)

		optimizedEmails[i] = gin.H{
			"id":          email.ID,
			"from":        from,              // 发件人（已清理）
			"to":          to,                // 收件人（已清理）
			"received_at": email.ReceivedAt,  // 收件时间
			"subject":     email.Subject,     // 主题
			"content":     body,              // 纯文本内容（已解析和清理）
		}
	}

	optimizedResponse := gin.H{
		"emails": optimizedEmails,
		"total":  response.Total,
		"page":   response.Page,
		"limit":  response.Limit,
	}

	c.JSON(http.StatusOK, optimizedResponse)
}

// GetEmailByID 根据ID获取单个邮件
func (h *EmailHandler) GetEmailByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email ID",
		})
		return
	}

	email, err := h.db.GetEmailByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Email not found",
		})
		return
	}

	// 解析和清理邮件内容
	body := email.Body
	htmlBody := email.HTMLBody
	
	// 首先检查是否是纯Base64编码的内容
	if body != "" && isBase64Content(body) {
		if decoded, err := base64.StdEncoding.DecodeString(body); err == nil {
			body = string(decoded)
		}
	}
	
	// 优先尝试解析MIME格式内容（包含base64解码）
	if body != "" && isMIMEContent(body) {
		if parsedContent, parseErr := utils.ParseMIMEContent(body); parseErr == nil {
			if parsedContent.TextBody != "" {
				body = parsedContent.TextBody
			}
			if parsedContent.HTMLBody != "" {
				htmlBody = parsedContent.HTMLBody
			}
		}
	} else if body != "" && (strings.Contains(body, "Content-Type:") || strings.Contains(body, "boundary=")) {
		// 如果body包含multipart数据，尝试标准解析
		if parsedContent, parseErr := utils.ParseEmailFromRaw(body); parseErr == nil {
			if parsedContent.TextBody != "" {
				body = parsedContent.TextBody
			}
			if parsedContent.HTMLBody != "" {
				htmlBody = parsedContent.HTMLBody
			}
		}
	}
	
	// 如果还是空的，尝试从raw_email解析
	if (body == "" || htmlBody == "") && email.RawEmail != "" {
		rawEmail := email.RawEmail
		
		// 检查raw_email是否是Base64编码
		if isBase64Content(rawEmail) {
			if decoded, err := base64.StdEncoding.DecodeString(rawEmail); err == nil {
				rawEmail = string(decoded)
			}
		}
		
		if isMIMEContent(rawEmail) {
			if parsedContent, parseErr := utils.ParseMIMEContent(rawEmail); parseErr == nil {
				if body == "" && parsedContent.TextBody != "" {
					body = parsedContent.TextBody
				}
				if htmlBody == "" && parsedContent.HTMLBody != "" {
					htmlBody = parsedContent.HTMLBody
				}
			}
		} else if parsedContent, parseErr := utils.ParseEmailFromRaw(rawEmail); parseErr == nil {
			if body == "" && parsedContent.TextBody != "" {
				body = parsedContent.TextBody
			}
			if htmlBody == "" && parsedContent.HTMLBody != "" {
				htmlBody = parsedContent.HTMLBody
			}
		}
	}

	// 如果没有HTML内容但有纯文本内容，将纯文本转换为HTML
	if htmlBody == "" && body != "" {
		htmlBody = textToHTML(body)
	}

	// 清理发件人和收件人字段
	from := cleanEmailAddress(email.From)
	to := cleanEmailAddress(email.To)

	// 返回完整的邮件详情，包括头部信息
	response := gin.H{
		"id":          email.ID,
		"from":        from,              // 发件人（已清理）
		"to":          to,                // 收件人（已清理）
		"received_at": email.ReceivedAt,  // 收件时间
		"created_at":  email.CreatedAt,   // 创建时间
		"subject":     email.Subject,     // 主题
		"body":        body,              // 纯文本内容（已解析和清理）
		"html_body":   htmlBody,          // HTML内容（仅用于详情查看）
		"headers":     email.Headers,     // 邮件头部信息
	}

	c.JSON(http.StatusOK, response)
}

// isMIMEContent 检查内容是否为MIME格式
func isMIMEContent(content string) bool {
	// 检查是否包含MIME边界标识符
	lines := strings.Split(content, "\r\n")
	for _, line := range lines {
		// 查找以 -- 开头的边界线，且长度合理
		if strings.HasPrefix(line, "--") && len(line) > 10 {
			// 进一步验证是否包含Content-Type
			if strings.Contains(content, "Content-Type:") {
				return true
			}
		}
	}
	return false
}

// cleanEmailAddress 清理邮件地址，移除多余的格式
func cleanEmailAddress(address string) string {
	if address == "" {
		return ""
	}
	
	// 移除引号和尖括号，提取纯邮件地址
	// 例如: "lammy2021" <lammy2021@126.com> -> lammy2021@126.com
	re := regexp.MustCompile(`<([^>]+)>`)
	matches := re.FindStringSubmatch(address)
	if len(matches) > 1 {
		return matches[1]
	}
	
	// 移除引号
	address = strings.Trim(address, `"`)
	address = strings.TrimSpace(address)
	
	return address
}

// HealthCheck 健康检查端点
func (h *EmailHandler) HealthCheck(c *gin.Context) {
	// 检查是否提供了认证信息
	token := c.GetHeader("Authorization")
	if token == "" {
		token = c.Query("token")
	}
	
	// 如果提供了认证信息，验证它
	authenticated := false
	if token != "" {
		if token == "Bearer "+h.authToken || token == h.authToken {
			authenticated = true
		}
	}
	
	response := gin.H{
		"status": "ok",
		"message": "Email receiver service is running",
		"authenticated": authenticated,
	}
	
	// 如果提供了无效的认证信息，返回401而不是403
	if token != "" && !authenticated {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"message": "Invalid authentication token",
			"authenticated": false,
			"debug_info": gin.H{
				"received_token": token,
				"expected_format": "Bearer " + h.authToken,
			},
		})
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// isBase64Content 检查内容是否是Base64编码
func isBase64Content(content string) bool {
	// 移除换行符和空格
	cleanContent := strings.ReplaceAll(content, "\r\n", "")
	cleanContent = strings.ReplaceAll(cleanContent, "\n", "")
	cleanContent = strings.ReplaceAll(cleanContent, " ", "")
	
	// Base64内容应该只包含Base64字符集
	if len(cleanContent) == 0 {
		return false
	}
	
	// 检查长度是否是4的倍数（Base64特征）
	if len(cleanContent)%4 != 0 {
		return false
	}
	
	// 检查是否只包含Base64字符
	for _, char := range cleanContent {
		if !((char >= 'A' && char <= 'Z') ||
			 (char >= 'a' && char <= 'z') ||
			 (char >= '0' && char <= '9') ||
			 char == '+' || char == '/' || char == '=') {
			return false
		}
	}
	
	// 尝试解码以验证是否是有效的Base64
	_, err := base64.StdEncoding.DecodeString(cleanContent)
	return err == nil
}

// textToHTML 将纯文本转换为HTML格式
func textToHTML(text string) string {
	if text == "" {
		return ""
	}
	
	// 转义HTML特殊字符
	text = strings.ReplaceAll(text, "&", "&")
	text = strings.ReplaceAll(text, "<", "<")
	text = strings.ReplaceAll(text, ">", ">")
	text = strings.ReplaceAll(text, "\"", "&quot;")
	text = strings.ReplaceAll(text, "'", "&#39;")
	
	// 将换行符转换为<br>标签
	text = strings.ReplaceAll(text, "\r\n", "<br>")
	text = strings.ReplaceAll(text, "\n", "<br>")
	text = strings.ReplaceAll(text, "\r", "<br>")
	
	// 将多个空格转换为&nbsp;
	text = regexp.MustCompile(`  +`).ReplaceAllStringFunc(text, func(spaces string) string {
		result := ""
		for i := 0; i < len(spaces); i++ {
			if i == 0 {
				result += " "
			} else {
				result += "&nbsp;"
			}
		}
		return result
	})
	
	// 检测并转换URL为链接
	urlRegex := regexp.MustCompile(`https?://[^\s<>"]+`)
	text = urlRegex.ReplaceAllStringFunc(text, func(url string) string {
		return `<a href="` + url + `" target="_blank" rel="noopener noreferrer">` + url + `</a>`
	})
	
	// 检测并转换邮箱地址为链接
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	text = emailRegex.ReplaceAllStringFunc(text, func(email string) string {
		return `<a href="mailto:` + email + `">` + email + `</a>`
	})
	
	// 包装在HTML结构中
	return `<div style="font-family: monospace; white-space: pre-wrap; word-wrap: break-word;">` + text + `</div>`
}