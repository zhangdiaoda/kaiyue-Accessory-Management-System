package handler

import (
	"net/http"
	"warehouse/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SystemHandler struct {
	db *gorm.DB
}

func NewSystemHandler(db *gorm.DB) *SystemHandler {
	return &SystemHandler{db: db}
}

// GetConfig 获取系统配置
func (h *SystemHandler) GetConfig(c *gin.Context) {
	var configs []model.SysConfig
	if err := h.db.Find(&configs).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "获取配置失败"})
		return
	}

	configMap := make(map[string]string)
	for _, cfg := range configs {
		configMap[cfg.ConfigKey] = cfg.ConfigValue
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data":    configMap,
	})
}

// GetBrandingConfig 获取公开的品牌配置（登录页用）
func (h *SystemHandler) GetBrandingConfig(c *gin.Context) {
	var configs []model.SysConfig
	keys := []string{"system_name", "company_name", "brand_logo", "copyright", "login_subtitle"}
	if err := h.db.Where("config_key IN ?", keys).Find(&configs).Error; err != nil {

		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "获取配置失败"})
		return
	}

	configMap := make(map[string]string)
	for _, cfg := range configs {
		configMap[cfg.ConfigKey] = cfg.ConfigValue
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data":    configMap,
	})
}

// UpdateConfig 更新系统配置
func (h *SystemHandler) UpdateConfig(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	for key, value := range req {
		var config model.SysConfig
		result := h.db.Where("config_key = ?", key).First(&config)
		if result.Error == gorm.ErrRecordNotFound {
			h.db.Create(&model.SysConfig{ConfigKey: key, ConfigValue: value})
		} else {
			h.db.Model(&config).Update("config_value", value)
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "配置更新成功"})
}

// GetAnnouncements 获取公告列表
func (h *SystemHandler) GetAnnouncements(c *gin.Context) {
	var list []model.Announcement
	query := h.db.Order("created_at desc")

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&list).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "获取公告失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data":    list,
	})
}

// CreateAnnouncement 创建公告
func (h *SystemHandler) CreateAnnouncement(c *gin.Context) {
	var announcement model.Announcement
	if err := c.ShouldBindJSON(&announcement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := h.db.Create(&announcement).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "公告发布成功"})
}

// UpdateAnnouncement 更新公告
func (h *SystemHandler) UpdateAnnouncement(c *gin.Context) {
	id := c.Param("id")
	var announcement model.Announcement
	if err := h.db.First(&announcement, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "公告不存在"})
		return
	}

	if err := c.ShouldBindJSON(&announcement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := h.db.Save(&announcement).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "操作成功"})
}

// DeleteAnnouncement 删除公告
func (h *SystemHandler) DeleteAnnouncement(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&model.Announcement{}, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "操作成功"})
}
