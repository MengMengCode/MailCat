package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"strings"
)

// EmailContent 解析后的邮件内容
type EmailContent struct {
	TextBody string `json:"text_body"`
	HTMLBody string `json:"html_body"`
	Subject  string `json:"subject"`
	From     string `json:"from"`
	To       string `json:"to"`
}

// ParseEmailContent 解析邮件内容
func ParseEmailContent(rawEmail string, headers map[string]string) (*EmailContent, error) {
	content := &EmailContent{
		Subject: headers["subject"],
		From:    headers["from"],
		To:      headers["to"],
	}

	// 检查Content-Type
	contentType := headers["content-type"]
	if contentType == "" {
		// 如果没有Content-Type，尝试直接使用原始内容
		content.TextBody = rawEmail
		return content, nil
	}

	// 解析Content-Type
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		// 如果解析失败，尝试直接使用原始内容
		content.TextBody = rawEmail
		return content, nil
	}

	switch {
	case strings.HasPrefix(mediaType, "multipart/"):
		return parseMultipartEmail(rawEmail, mediaType, params)
	case mediaType == "text/plain":
		decoded, err := decodeContent(rawEmail, headers)
		if err != nil {
			content.TextBody = rawEmail
		} else {
			content.TextBody = decoded
		}
		return content, nil
	case mediaType == "text/html":
		decoded, err := decodeContent(rawEmail, headers)
		if err != nil {
			content.HTMLBody = rawEmail
		} else {
			content.HTMLBody = decoded
		}
		return content, nil
	default:
		content.TextBody = rawEmail
		return content, nil
	}
}

// parseMultipartEmail 解析multipart邮件
func parseMultipartEmail(rawEmail, mediaType string, params map[string]string) (*EmailContent, error) {
	content := &EmailContent{}
	
	boundary := params["boundary"]
	if boundary == "" {
		return content, fmt.Errorf("multipart email missing boundary")
	}

	// 创建multipart reader
	reader := multipart.NewReader(strings.NewReader(rawEmail), boundary)
	
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		// 读取part内容
		partContent, err := io.ReadAll(part)
		if err != nil {
			part.Close()
			continue
		}
		part.Close()

		// 获取part的Content-Type
		partContentType := part.Header.Get("Content-Type")
		partMediaType, _, _ := mime.ParseMediaType(partContentType)

		// 解码内容
		decodedContent := string(partContent)
		
		// 检查Content-Transfer-Encoding
		encoding := part.Header.Get("Content-Transfer-Encoding")
		if encoding != "" {
			decoded, err := decodeByEncoding(string(partContent), encoding)
			if err == nil {
				decodedContent = decoded
			}
		}

		// 根据Content-Type分配内容
		switch partMediaType {
		case "text/plain":
			if content.TextBody == "" {
				content.TextBody = decodedContent
			}
		case "text/html":
			if content.HTMLBody == "" {
				content.HTMLBody = decodedContent
			}
		}
	}

	return content, nil
}

// decodeContent 根据编码解码内容
func decodeContent(content string, headers map[string]string) (string, error) {
	encoding := headers["content-transfer-encoding"]
	if encoding == "" {
		return content, nil
	}

	return decodeByEncoding(content, encoding)
}

// decodeByEncoding 根据指定编码解码内容
func decodeByEncoding(content, encoding string) (string, error) {
	encoding = strings.ToLower(strings.TrimSpace(encoding))
	
	switch encoding {
	case "base64":
		decoded, err := base64.StdEncoding.DecodeString(content)
		if err != nil {
			return content, err
		}
		return string(decoded), nil
		
	case "quoted-printable":
		reader := quotedprintable.NewReader(strings.NewReader(content))
		decoded, err := io.ReadAll(reader)
		if err != nil {
			return content, err
		}
		return string(decoded), nil
		
	case "7bit", "8bit", "binary":
		return content, nil
		
	default:
		return content, nil
	}
}

// ParseEmailFromRaw 从原始邮件数据解析
func ParseEmailFromRaw(rawEmail string) (*EmailContent, error) {
	// 尝试解析为标准邮件格式
	msg, err := mail.ReadMessage(strings.NewReader(rawEmail))
	if err != nil {
		return nil, fmt.Errorf("failed to parse email: %w", err)
	}

	content := &EmailContent{
		Subject: msg.Header.Get("Subject"),
		From:    msg.Header.Get("From"),
		To:      msg.Header.Get("To"),
	}

	// 读取邮件体
	body, err := io.ReadAll(msg.Body)
	if err != nil {
		return content, err
	}

	// 获取Content-Type
	contentType := msg.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		content.TextBody = string(body)
		return content, nil
	}

	switch {
	case strings.HasPrefix(mediaType, "multipart/"):
		return parseMultipartFromBody(string(body), mediaType, params)
	case mediaType == "text/plain":
		encoding := msg.Header.Get("Content-Transfer-Encoding")
		decoded, err := decodeByEncoding(string(body), encoding)
		if err != nil {
			content.TextBody = string(body)
		} else {
			content.TextBody = decoded
		}
		return content, nil
	case mediaType == "text/html":
		encoding := msg.Header.Get("Content-Transfer-Encoding")
		decoded, err := decodeByEncoding(string(body), encoding)
		if err != nil {
			content.HTMLBody = string(body)
		} else {
			content.HTMLBody = decoded
		}
		return content, nil
	default:
		content.TextBody = string(body)
		return content, nil
	}
}

// parseMultipartFromBody 从邮件体解析multipart内容
func parseMultipartFromBody(body, mediaType string, params map[string]string) (*EmailContent, error) {
	content := &EmailContent{}
	
	boundary := params["boundary"]
	if boundary == "" {
		return content, fmt.Errorf("multipart email missing boundary")
	}

	reader := multipart.NewReader(strings.NewReader(body), boundary)
	
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		partContent, err := io.ReadAll(part)
		if err != nil {
			part.Close()
			continue
		}
		part.Close()

		partContentType := part.Header.Get("Content-Type")
		partMediaType, _, _ := mime.ParseMediaType(partContentType)

		decodedContent := string(partContent)
		encoding := part.Header.Get("Content-Transfer-Encoding")
		if encoding != "" {
			decoded, err := decodeByEncoding(string(partContent), encoding)
			if err == nil {
				decodedContent = decoded
			}
		}

		switch partMediaType {
		case "text/plain":
			if content.TextBody == "" {
				content.TextBody = decodedContent
			}
		case "text/html":
			if content.HTMLBody == "" {
				content.HTMLBody = decodedContent
			}
		}
	}

	return content, nil
}

// TryParseEmailContent 尝试多种方式解析邮件内容
func TryParseEmailContent(rawBody, htmlBody string, headersJSON string) (string, string, error) {
	// 如果已经有内容，直接返回
	if rawBody != "" || htmlBody != "" {
		return rawBody, htmlBody, nil
	}

	// 解析headers
	var headers map[string]string
	if headersJSON != "" {
		if err := json.Unmarshal([]byte(headersJSON), &headers); err != nil {
			return rawBody, htmlBody, fmt.Errorf("failed to parse headers: %w", err)
		}
	}

	// 检查是否是multipart邮件
	contentType := headers["content-type"]
	if contentType == "" {
		return rawBody, htmlBody, nil
	}

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return rawBody, htmlBody, nil
	}

	// 如果不是multipart，无法进一步解析
	if !strings.HasPrefix(mediaType, "multipart/") {
		return rawBody, htmlBody, nil
	}

	// 尝试从headers中提取原始邮件内容
	// 这里需要原始的邮件数据，但我们只有headers
	// 在实际应用中，应该在接收邮件时就解析内容
	
	return rawBody, htmlBody, fmt.Errorf("multipart email requires raw email data for parsing")
}

// ParseMIMEContent 解析MIME格式的邮件内容，包含base64解码
func ParseMIMEContent(content string) (*EmailContent, error) {
	result := &EmailContent{}
	
	// 检测MIME边界
	boundary := detectMIMEBoundary(content)
	if boundary == "" {
		// 不是multipart格式，直接返回原内容
		result.TextBody = content
		return result, nil
	}
	
	// 按行分割内容
	lines := strings.Split(content, "\r\n")
	
	var currentPart strings.Builder
	var currentEncoding string
	var currentContentType string
	var inHeaders bool = false
	var inContent bool = false
	
	for _, line := range lines {
		// 检查是否是边界线
		if strings.HasPrefix(line, "--"+boundary) {
			// 处理之前的部分
			if inContent && currentPart.Len() > 0 {
				partContent := strings.TrimSpace(currentPart.String())
				decodedContent, err := decodePartContent(partContent, currentEncoding)
				if err == nil {
					assignContentByType(result, decodedContent, currentContentType)
				}
			}
			
			// 检查是否是结束边界
			if strings.HasSuffix(line, "--") {
				break
			}
			
			// 重置状态，开始新的部分
			currentPart.Reset()
			currentEncoding = ""
			currentContentType = ""
			inHeaders = true
			inContent = false
			continue
		}
		
		if inHeaders {
			if line == "" {
				// 空行表示头部结束，内容开始
				inHeaders = false
				inContent = true
				continue
			}
			
			// 解析头部信息
			if strings.HasPrefix(line, "Content-Type:") {
				currentContentType = strings.TrimSpace(strings.TrimPrefix(line, "Content-Type:"))
			} else if strings.HasPrefix(line, "Content-Transfer-Encoding:") {
				currentEncoding = strings.TrimSpace(strings.TrimPrefix(line, "Content-Transfer-Encoding:"))
			}
		} else if inContent {
			// 添加内容行
			if currentPart.Len() > 0 {
				currentPart.WriteString("\r\n")
			}
			currentPart.WriteString(line)
		}
	}
	
	// 处理最后一个部分
	if inContent && currentPart.Len() > 0 {
		partContent := strings.TrimSpace(currentPart.String())
		decodedContent, err := decodePartContent(partContent, currentEncoding)
		if err == nil {
			assignContentByType(result, decodedContent, currentContentType)
		}
	}
	
	return result, nil
}

// detectMIMEBoundary 检测MIME边界标识符
func detectMIMEBoundary(content string) string {
	lines := strings.Split(content, "\r\n")
	
	for _, line := range lines {
		// 查找以 -- 开头的边界线
		if strings.HasPrefix(line, "--") && len(line) > 2 {
			// 提取边界标识符（去掉开头的 --）
			boundary := strings.TrimPrefix(line, "--")
			// 去掉可能的结束标记 --
			boundary = strings.TrimSuffix(boundary, "--")
			
			// 验证这是一个有效的边界（应该包含字母数字字符）
			if len(boundary) > 0 && isValidBoundary(boundary) {
				return boundary
			}
		}
	}
	
	return ""
}

// isValidBoundary 验证边界标识符是否有效
func isValidBoundary(boundary string) bool {
	// 边界应该只包含字母、数字和一些特殊字符
	for _, char := range boundary {
		if !((char >= 'a' && char <= 'z') ||
			 (char >= 'A' && char <= 'Z') ||
			 (char >= '0' && char <= '9') ||
			 char == '_' || char == '-' || char == '=' || char == '.') {
			return false
		}
	}
	return len(boundary) > 5 // 边界应该有一定长度
}

// decodePartContent 根据编码类型解码内容
func decodePartContent(content, encoding string) (string, error) {
	encoding = strings.ToLower(strings.TrimSpace(encoding))
	
	switch encoding {
	case "base64":
		// 清理base64字符串（移除换行符和空格）
		cleanContent := strings.ReplaceAll(content, "\r\n", "")
		cleanContent = strings.ReplaceAll(cleanContent, "\n", "")
		cleanContent = strings.ReplaceAll(cleanContent, " ", "")
		
		decoded, err := base64.StdEncoding.DecodeString(cleanContent)
		if err != nil {
			return content, err
		}
		return string(decoded), nil
	case "quoted-printable":
		reader := quotedprintable.NewReader(strings.NewReader(content))
		decoded, err := io.ReadAll(reader)
		if err != nil {
			return content, err
		}
		return string(decoded), nil
	default:
		return content, nil
	}
}

// assignContentByType 根据内容类型分配解码后的内容
func assignContentByType(result *EmailContent, content, contentType string) {
	contentType = strings.ToLower(strings.TrimSpace(contentType))
	
	if strings.Contains(contentType, "text/plain") {
		if result.TextBody == "" {
			result.TextBody = content
		}
	} else if strings.Contains(contentType, "text/html") {
		if result.HTMLBody == "" {
			result.HTMLBody = content
		}
	}
}