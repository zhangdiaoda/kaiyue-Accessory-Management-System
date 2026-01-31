package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"
	"warehouse/internal/model"
	"warehouse/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuditMiddleware 审计日志中间件
// 仅记录仓库管理员的操作
func AuditMiddleware(db *gorm.DB) gin.HandlerFunc {
	auditService := service.NewAuditService(db)

	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()

		// 备份请求体（用于日志记录）
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 执行请求
		c.Next()

		// 获取用户信息
		username, exists := c.Get("username")
		if !exists {
			return // 未登录的请求不记录
		}

		role, _ := c.Get("role")
		if role != "WAREHOUSE_ADMIN" {
			return // 只记录仓库管理员的操作
		}

		userID, _ := c.Get("user_id")
		realName, _ := c.Get("real_name")

		// 计算执行时间
		duration := time.Since(startTime).Milliseconds()

		// 确定操作状态
		status := "SUCCESS"
		if c.Writer.Status() >= 400 {
			status = "FAILED"
		}

		// 创建日志记录
		log := &model.OperationLog{
			UserID:        userID.(uint),
			Username:      username.(string),
			RealName:      realName.(string),
			Operation:     getOperationType(c.Request.Method, c.FullPath()),
			Module:        getModule(c.FullPath()),
			Description:   getDescription(c),
			RequestMethod: c.Request.Method,
			RequestURL:    c.Request.URL.String(),
			RequestParams: sanitizeParams(requestBody),
			IPAddress:     c.ClientIP(),
			UserAgent:     c.Request.UserAgent(),
			Status:        status,
			ExecutionTime: int(duration),
		}

		// 异步记录日志（不阻塞主流程）
		go func() {
			if err := auditService.RecordOperation(log); err != nil {
				// 记录日志失败，打印错误（不影响业务）
				println("审计日志记录失败:", err.Error())
			}
		}()
	}
}

// getOperationType 根据请求方法和路径获取操作类型
func getOperationType(method, path string) string {
	switch method {
	case "POST":
		if contains(path, "import") {
			return "导入"
		}
		return "创建"
	case "PUT", "PATCH":
		return "更新"
	case "DELETE":
		return "删除"
	case "GET":
		if contains(path, "export") || contains(path, "download") {
			return "导出"
		}
		return "查询"
	default:
		return method
	}
}

// getModule 根据路径获取模块名称
func getModule(path string) string {
	if contains(path, "/parts") {
		return "配件管理"
	} else if contains(path, "/borrow") {
		return "领用管理"
	} else if contains(path, "/employee") {
		return "员工管理"
	} else if contains(path, "/report") {
		return "报表管理"
	} else if contains(path, "/system") {
		return "系统管理"
	}
	return "其他"
}

// getDescription 生成操作描述
func getDescription(c *gin.Context) string {
	operation := getOperationType(c.Request.Method, c.FullPath())
	module := getModule(c.FullPath())
	return operation + " - " + module
}

// sanitizeParams 清理敏感参数（如密码）
func sanitizeParams(body []byte) string {
	if len(body) == 0 {
		return ""
	}

	// 尝试解析JSON
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return string(body) // 非JSON直接返回
	}

	// 移除敏感字段
	sensitiveFields := []string{"password", "token", "secret"}
	for _, field := range sensitiveFields {
		if _, exists := data[field]; exists {
			data[field] = "***"
		}
	}

	result, _ := json.Marshal(data)
	return string(result)
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[0:len(substr)] == substr || s[len(s)-len(substr):] == substr || bytes.Contains([]byte(s), []byte(substr))))
}
