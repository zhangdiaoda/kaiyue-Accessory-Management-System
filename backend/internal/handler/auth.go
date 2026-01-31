package handler

import (
	"net/http"
	"warehouse/internal/model"
	"warehouse/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token      string `json:"token"`
	Username   string `json:"username"`
	RealName   string `json:"real_name"`
	Role       string `json:"role"`
	Department string `json:"department"`
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 查询用户
	var user model.User
	if err := h.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "用户名或密码错误",
		})
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "用户名或密码错误",
		})
		return
	}

	// 检查用户状态
	if user.Status == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "账号已被禁用",
		})
		return
	}

	// 生成Token
	token, err := utils.GenerateToken(user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "生成Token失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data": LoginResponse{
			Token:      token,
			Username:   user.Username,
			RealName:   user.RealName,
			Role:       user.Role,
			Department: user.Department,
		},
	})
}

// GetUserInfo 获取当前用户信息
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	username, _ := c.Get("username")

	var user model.User
	if err := h.db.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data": LoginResponse{
			Username:   user.Username,
			RealName:   user.RealName,
			Role:       user.Role,
			Department: user.Department,
		},
	})
}

// Logout 退出登录
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "退出登录成功",
	})
}

// UpdateProfile 更新个人信息
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	username, _ := c.Get("username")

	var req struct {
		RealName string `json:"real_name"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var user model.User
	if err := h.db.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "用户不存在"})
		return
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Password != "" {
		hashedPassword, _ := utils.HashPassword(req.Password)
		updates["password"] = hashedPassword
	}

	if err := h.db.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "操作成功"})
}

// GetUserList 获取系统用户列表
func (h *AuthHandler) GetUserList(c *gin.Context) {
	var users []model.User
	if err := h.db.Order("id desc").Find(&users).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "获取用户列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data":    users,
	})
}

// CreateUser 创建系统用户
func (h *AuthHandler) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 检查用户名是否存在
	var count int64
	h.db.Model(&model.User{}).Where("username = ?", user.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "用户名已存在"})
		return
	}

	// 加密密码
	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功"})
}

// UpdateUser 更新系统用户
func (h *AuthHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var user model.User
	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "用户不存在"})
		return
	}

	updates := map[string]interface{}{
		"real_name":  req.RealName,
		"role":       req.Role,
		"department": req.Department,
		"phone":      req.Phone,
		"status":     req.Status,
	}

	// 如果传入了新密码，则更新密码
	if req.Password != "" {
		hashedPassword, _ := utils.HashPassword(req.Password)
		updates["password"] = hashedPassword
	}

	if err := h.db.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "操作成功"})
}

// DeleteUser 删除系统用户
func (h *AuthHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	// 禁止删除admin
	var user model.User
	h.db.First(&user, id)
	if user.Username == "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "超级管理员不能删除"})
		return
	}

	if err := h.db.Delete(&model.User{}, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "操作成功"})
}

// ResetPassword 重置密码
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	hashedPassword, _ := utils.HashPassword(req.Password)
	if err := h.db.Model(&model.User{}).Where("id = ?", id).Update("password", hashedPassword).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "重置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "密码重置成功"})
}
