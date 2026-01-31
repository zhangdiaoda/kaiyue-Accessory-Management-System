package middleware

import (
	"net/http"
	"strings"
	"warehouse/internal/model"
	"warehouse/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := ""
		auth := c.GetHeader("Authorization")

		if auth != "" {
			// 从Header获取
			tokenString = strings.TrimPrefix(auth, "Bearer ")
		} else {
			// 尝试从Query获取
			tokenString = c.Query("token")
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权，请先登录",
			})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Token无效或已过期",
			})
			c.Abort()
			return
		}

		// 查询用户信息
		var user model.User
		if err := db.Where("username = ?", claims.Username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "用户不存在",
			})
			c.Abort()
			return
		}

		// 将完整用户信息存入上下文
		c.Set("user_id", user.ID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("real_name", user.RealName)
		c.Next()
	}
}

// CORS中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
