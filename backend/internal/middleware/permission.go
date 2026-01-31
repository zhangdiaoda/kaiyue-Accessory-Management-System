package middleware

import (
	"net/http"
	"warehouse/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PermissionMiddleware 权限检查中间件
// 检查用户是否拥有指定权限
func PermissionMiddleware(db *gorm.DB, requiredPermission string) gin.HandlerFunc {
	permService := service.NewPermissionService(db)

	return func(c *gin.Context) {
		// 获取用户信息
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
			})
			c.Abort()
			return
		}

		// 超级管理员拥有所有权限
		if role == "SUPER_ADMIN" {
			c.Next()
			return
		}

		// 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
			})
			c.Abort()
			return
		}

		// 检查权限
		hasPermission, err := permService.CheckPermission(userID.(uint), role.(string), requiredPermission)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "权限检查失败",
			})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "无权限执行此操作",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
