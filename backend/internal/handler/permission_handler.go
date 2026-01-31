package handler

import (
	"net/http"
	"strconv"
	"warehouse/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PermissionHandler 权限管理Handler
type PermissionHandler struct {
	db          *gorm.DB
	permService *service.PermissionService
}

// NewPermissionHandler 创建权限Handler
func NewPermissionHandler(db *gorm.DB) *PermissionHandler {
	return &PermissionHandler{
		db:          db,
		permService: service.NewPermissionService(db),
	}
}

// GetAllPermissions 获取所有权限（按分类分组）
func (h *PermissionHandler) GetAllPermissions(c *gin.Context) {
	permissions, err := h.permService.GetPermissionsByCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取权限失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": permissions,
	})
}

// GetUserPermissions 获取用户权限
func (h *PermissionHandler) GetUserPermissions(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的用户ID",
		})
		return
	}

	// 获取用户信息
	user, err := h.permService.GetUserInfo(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "用户不存在",
		})
		return
	}

	// 获取用户权限
	permissions, err := h.permService.GetUserPermissions(uint(userID), user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取用户权限失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": permissions,
	})
}

// SetUserPermissions 设置用户权限
func (h *PermissionHandler) SetUserPermissions(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的用户ID",
		})
		return
	}

	var req struct {
		Permissions []string `json:"permissions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	if err := h.permService.SetUserPermissions(uint(userID), req.Permissions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "设置权限失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "权限设置成功",
	})
}

// GetRolePermissions 获取角色默认权限
func (h *PermissionHandler) GetRolePermissions(c *gin.Context) {
	role := c.Param("role")

	permissions, err := h.permService.GetRolePermissions(role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取角色权限失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": permissions,
	})
}

// SetRolePermissions 设置角色默认权限
func (h *PermissionHandler) SetRolePermissions(c *gin.Context) {
	role := c.Param("role")

	var req struct {
		Permissions []string `json:"permissions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	if err := h.permService.SetRolePermissions(role, req.Permissions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "设置角色权限失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "角色权限设置成功",
	})
}
