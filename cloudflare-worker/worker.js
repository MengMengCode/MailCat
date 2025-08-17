// Cloudflare Worker 示例代码
// 用于接收邮件并转发到Go API

export default {
  // 处理HTTP请求
  async fetch(request, env, ctx) {
    const url = new URL(request.url);
    
    // 处理根路径，显示状态页面
    if (url.pathname === '/') {
      return await this.handleStatusPage(env);
    }
    
    // 处理health检查
    if (url.pathname === '/health') {
      return await this.handleHealthCheck(env);
    }
    
    // 其他路径返回404
    return new Response('Not Found', { status: 404 });
  },

  // 处理状态页面
  async handleStatusPage(env) {
    const healthStatus = await this.checkAPIHealth(env);
    
    const html = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>邮件接收器 Worker 状态</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background: white;
            border-radius: 8px;
            padding: 30px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .status {
            display: flex;
            align-items: center;
            margin: 20px 0;
            padding: 15px;
            border-radius: 6px;
            font-weight: 500;
        }
        .status.success {
            background-color: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }
        .status.error {
            background-color: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }
        .status-icon {
            margin-right: 10px;
            font-size: 20px;
        }
        .info {
            background-color: #f8f9fa;
            padding: 15px;
            border-radius: 6px;
            margin: 20px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>📧 邮件接收器 Worker 状态</h1>
        
        <div class="status ${healthStatus.success ? 'success' : 'error'}">
            <span class="status-icon">${healthStatus.success ? '✅' : '❌'}</span>
            <div>
                <strong>API连接状态: ${healthStatus.success ? '成功' : '失败'}</strong>
                <br>
                <small>${healthStatus.message}</small>
            </div>
        </div>
        
        ${!healthStatus.success ? `
        <div class="info">
            <h3>错误详情</h3>
            ${healthStatus.debug_info && healthStatus.debug_info.error ? `<p><strong>错误信息:</strong> <code>${healthStatus.debug_info.error}</code></p>` : ''}
            ${healthStatus.debug_info && healthStatus.debug_info.solution ? `
            <div style="background-color: #fff3cd; border: 1px solid #ffeaa7; padding: 10px; border-radius: 4px; margin-top: 10px;">
                <strong>💡 解决方案:</strong> ${healthStatus.debug_info.solution}
            </div>
            ` : ''}
        </div>
        ` : ''}
        
        <div class="info">
            <h3>注意事项</h3>
            <p>🌐 <strong>API_ENDPOINT必须使用域名，不支持IP地址</strong></p>
        </div>
    </div>
</body>
</html>`;
    
    return new Response(html, {
      headers: { 'Content-Type': 'text/html; charset=utf-8' }
    });
  },

  // 处理health检查API
  async handleHealthCheck(env) {
    const healthStatus = await this.checkAPIHealth(env);
    
    return new Response(JSON.stringify({
      status: healthStatus.success ? 'healthy' : 'unhealthy',
      message: healthStatus.message,
      timestamp: new Date().toISOString(),
      api_endpoint: env.API_ENDPOINT || null,
      token_configured: !!env.API_TOKEN
    }), {
      headers: {
        'Content-Type': 'application/json',
        'Access-Control-Allow-Origin': '*'
      }
    });
  },

  // 检查API健康状态
  async checkAPIHealth(env) {
    if (!env.API_ENDPOINT) {
      return {
        success: false,
        message: 'API_ENDPOINT 环境变量未配置',
        debug_info: null
      };
    }
    
    if (!env.API_TOKEN) {
      return {
        success: false,
        message: 'API_TOKEN 环境变量未配置',
        debug_info: null
      };
    }
    
    const healthUrl = env.API_ENDPOINT + '/health';
    
    // 检查是否是IP地址或localhost
    const isIPAddress = /^https?:\/\/(\d+\.\d+\.\d+\.\d+|localhost|127\.0\.0\.1)/.test(env.API_ENDPOINT);
    
    if (isIPAddress) {
      return {
        success: false,
        message: 'Cloudflare Worker仅支持域名访问',
        debug_info: {
          url: healthUrl,
          method: 'GET',
          error: 'Cloudflare限制：Worker无法访问IP地址或localhost，必须使用域名',
          solution: '请将API_ENDPOINT改为域名格式，如: https://your-domain.com 或 https://api.example.com'
        }
      };
    }
    
    // 检查是否使用了有效的域名格式
    if (!env.API_ENDPOINT.match(/^https?:\/\/[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}/)) {
      return {
        success: false,
        message: 'API端点格式不正确',
        debug_info: {
          url: healthUrl,
          method: 'GET',
          error: 'API_ENDPOINT必须是完整的域名URL',
          solution: '请使用格式: https://your-domain.com 或 http://your-domain.com'
        }
      };
    }
    
    try {
      // 尝试访问API的health端点
      const response = await fetch(healthUrl, {
        method: 'GET',
        headers: {
          'Authorization': 'Bearer ' + env.API_TOKEN,
          'User-Agent': 'Cloudflare-Worker-Health-Check'
        },
        // 设置超时
        signal: AbortSignal.timeout(5000)
      });
      
      const responseText = await response.text();
      
      if (response.ok) {
        return {
          success: true,
          message: `API服务器响应正常 (状态码: ${response.status})`,
          debug_info: {
            url: healthUrl,
            method: 'GET',
            status: response.status,
            response: responseText
          }
        };
      } else {
        return {
          success: false,
          message: `API服务器响应异常 (状态码: ${response.status})`,
          debug_info: {
            url: healthUrl,
            method: 'GET',
            status: response.status,
            response: responseText
          }
        };
      }
    } catch (error) {
      return {
        success: false,
        message: `无法连接到API服务器: ${error.message}`,
        debug_info: {
          url: healthUrl,
          method: 'GET',
          error: error.message
        }
      };
    }
  },

  // 处理邮件
  async email(message, env, ctx) {
    console.log('=== Email Processing Started ===');
    
    try {
      // 检查环境变量
      if (!env.API_ENDPOINT) {
        console.error('API_ENDPOINT not configured');
        message.setReject('API endpoint not configured');
        return;
      }
      
      if (!env.API_TOKEN) {
        console.error('API_TOKEN not configured');
        message.setReject('API token not configured');
        return;
      }
      
      // 详细记录邮件信息
      console.log('Email from:', message.from);
      console.log('Email to:', message.to);
      console.log('Email headers:', JSON.stringify(Object.fromEntries(message.headers)));
      
      // 获取邮件主题
      const subject = message.headers.get('subject') || message.headers.get('Subject') || '';
      console.log('Email subject:', subject);
      
      // 解析邮件内容 - 使用不同的方法尝试获取内容
      let body = '';
      let htmlBody = '';
      let rawContent = '';
      
      try {
        // 尝试获取原始内容
        const reader = message.raw.getReader();
        const chunks = [];
        let done = false;
        
        while (!done) {
          const { value, done: readerDone } = await reader.read();
          done = readerDone;
          if (value) {
            chunks.push(value);
          }
        }
        
        const rawBytes = new Uint8Array(chunks.reduce((acc, chunk) => acc + chunk.length, 0));
        let offset = 0;
        for (const chunk of chunks) {
          rawBytes.set(chunk, offset);
          offset += chunk.length;
        }
        
        rawContent = new TextDecoder().decode(rawBytes);
        console.log('Raw email content length:', rawContent.length);
        console.log('Raw email preview:', rawContent.substring(0, 500));
        
      } catch (e) {
        console.warn('Failed to get raw content:', e.message);
      }
      
      try {
        body = await message.text();
        console.log('Text body length:', body.length);
        console.log('Text body preview:', body.substring(0, 200));
      } catch (e) {
        console.warn('Failed to get text body:', e.message);
      }
      
      try {
        htmlBody = await message.html() || '';
        console.log('HTML body length:', htmlBody.length);
      } catch (e) {
        console.warn('Failed to get HTML body:', e.message);
      }
      
      // 如果没有获取到正文，尝试从原始内容中提取
      if (!body && rawContent) {
        // 简单的邮件内容提取
        const lines = rawContent.split('\n');
        let inBody = false;
        const bodyLines = [];
        
        for (const line of lines) {
          if (inBody) {
            bodyLines.push(line);
          } else if (line.trim() === '') {
            inBody = true; // 空行后开始邮件正文
          }
        }
        
        if (bodyLines.length > 0) {
          body = bodyLines.join('\n').trim();
          console.log('Extracted body from raw content:', body.substring(0, 200));
        }
      }
      
      const emailData = {
        from: message.from,
        to: message.to,
        subject: subject,
        body: body || rawContent || '(无法解析邮件内容)',
        html_body: htmlBody,
        headers: Object.fromEntries(message.headers)
      };
      
      console.log('Final email data:', {
        from: emailData.from,
        to: emailData.to,
        subject: emailData.subject,
        bodyLength: emailData.body.length,
        htmlBodyLength: emailData.html_body.length,
        headerCount: Object.keys(emailData.headers).length
      });

      // 发送到Go API
      const response = await fetch(env.API_ENDPOINT + '/api/v1/emails', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + env.API_TOKEN
        },
        body: JSON.stringify(emailData)
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error('API request failed:', response.status, errorText);
        message.setReject(`API request failed: ${response.status}`);
        return;
      }

      const responseData = await response.text();
      console.log('Email successfully sent to API:', responseData);
      console.log('=== Email Processing Completed ===');
      
    } catch (error) {
      console.error('Error processing email:', error.message, error.stack);
      message.setReject(`Internal error: ${error.message}`);
    }
  }
};

// 环境变量配置说明：
// API_ENDPOINT: Go API服务器地址，例如 https://your-domain.com
// API_TOKEN: API认证令牌，与config.yaml中的auth_token相同