package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"warehouse/internal/model"
	"warehouse/internal/notification"
	"warehouse/internal/notification/providers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NotificationHandler struct {
	db        *gorm.DB
	manager   *notification.Manager
	scheduler *NotificationScheduler
}

func NewNotificationHandler(db *gorm.DB, manager *notification.Manager, scheduler *NotificationScheduler) *NotificationHandler {
	return &NotificationHandler{
		db:        db,
		manager:   manager,
		scheduler: scheduler,
	}
}

// RunDailyReportNow 立即运行每日报表
func (h *NotificationHandler) RunDailyReportNow(c *gin.Context) {
	if h.scheduler == nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "调度器未初始化"})
		return
	}
	h.scheduler.RunDailyReportNow()
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "已成功触发每日报表"})
}

// RunOverdueCheckNow 立即运行超期检查
func (h *NotificationHandler) RunOverdueCheckNow(c *gin.Context) {
	if h.scheduler == nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "调度器未初始化"})
		return
	}
	h.scheduler.RunOverdueCheckNow()
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "已成功触发超期归还检查"})
}

// GetConfigs 获取所有通知配置
func (h *NotificationHandler) GetConfigs(c *gin.Context) {
	var configs []model.NotificationConfig
	if err := h.db.Find(&configs).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data":    configs,
	})
}

// GetConfig 获取指定类型的通知配置
func (h *NotificationHandler) GetConfig(c *gin.Context) {
	providerType := c.Param("type")

	var config model.NotificationConfig
	err := h.db.Where("provider_type = ?", providerType).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "配置不存在",
			"data":    nil,
		})
		return
	} else if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data":    config,
	})
}

// UpdateConfig 更新通知配置
func (h *NotificationHandler) UpdateConfig(c *gin.Context) {
	var req struct {
		ProviderType string `json:"provider_type" binding:"required"`
		ConfigName   string `json:"config_name" binding:"required"`
		ConfigData   string `json:"config_data" binding:"required"`
		IsEnabled    bool   `json:"is_enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var config model.NotificationConfig
	err := h.db.Where("provider_type = ?", req.ProviderType).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		// 创建新配置
		config = model.NotificationConfig{
			ProviderType: req.ProviderType,
			ConfigName:   req.ConfigName,
			ConfigData:   req.ConfigData,
			IsEnabled:    req.IsEnabled,
		}
		if err := h.db.Create(&config).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 500, "message": "创建失败"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "查询失败"})
		return
	} else {
		// 更新配置
		config.ConfigName = req.ConfigName
		config.ConfigData = req.ConfigData
		config.IsEnabled = req.IsEnabled
		if err := h.db.Save(&config).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 500, "message": "更新失败"})
			return
		}
	}

	// 保存成功后，立即重新加载所有提供者，使配置生效
	providers.LoadAllProviders(h.db, h.manager)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "保存成功",
		"data":    config,
	})
}

// TestDingTalk 测试钉钉通知
func (h *NotificationHandler) TestDingTalk(c *gin.Context) {
	var req struct {
		WebhookURL string   `json:"webhook_url" binding:"required"`
		Secret     string   `json:"secret"`
		AtMobiles  []string `json:"at_mobiles"`
		IsAtAll    bool     `json:"is_at_all"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 创建临时提供者发送测试消息
	provider := providers.NewDingTalkProvider(&providers.DingTalkConfig{
		WebhookURL: req.WebhookURL,
		Secret:     req.Secret,
		AtMobiles:  req.AtMobiles,
		IsAtAll:    req.IsAtAll,
	})

	if err := provider.SendTestMessage(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "发送失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "测试消息发送成功",
	})
}

// TestWechat 测试微信通知
func (h *NotificationHandler) TestWechat(c *gin.Context) {
	var req struct {
		AppID      string `json:"app_id" binding:"required"`
		AppSecret  string `json:"app_secret" binding:"required"`
		TemplateID string `json:"template_id" binding:"required"`
		OpenID     string `json:"openid" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 创建临时提供者发送测试消息
	cache := notification.NewSimpleCache()
	provider := providers.NewWechatProvider(&providers.WechatConfig{
		AppID:      req.AppID,
		AppSecret:  req.AppSecret,
		TemplateID: req.TemplateID,
		TokenCache: cache,
	})

	if err := provider.SendTestMessage(req.OpenID); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "发送失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "测试消息发送成功",
	})
}

// GetLogs 获取通知日志
func (h *NotificationHandler) GetLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	providerType := c.Query("provider_type")
	status := c.Query("status")

	query := h.db.Model(&model.NotificationLog{})

	if providerType != "" {
		query = query.Where("provider_type = ?", providerType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var logs []model.NotificationLog
	offset := (page - 1) * pageSize
	if err := query.Order("created_at desc").Limit(pageSize).Offset(offset).Find(&logs).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data": gin.H{
			"list":  logs,
			"total": total,
			"page":  page,
		},
	})
}

// GetStats 获取统计信息
func (h *NotificationHandler) GetStats(c *gin.Context) {
	stats := h.manager.GetStats()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data":    stats,
	})
}

// SendManualNotification 手动发送通知
func (h *NotificationHandler) SendManualNotification(c *gin.Context) {
	var req struct {
		Title     string   `json:"title" binding:"required"`
		Content   string   `json:"content" binding:"required"`
		Providers []string `json:"providers" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	notif := notification.BuildSystemAnnouncementNotification(req.Title, req.Content)

	var providerTypes []notification.ProviderType
	for _, p := range req.Providers {
		providerTypes = append(providerTypes, notification.ProviderType(p))
	}

	h.manager.SendNotificationAsync(notif, providerTypes)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "通知已加入发送队列",
	})
}

// GetUserSettings 获取用户通知设置
func (h *NotificationHandler) GetUserSettings(c *gin.Context) {
	username, _ := c.Get("username")
	var user model.User
	h.db.Where("username = ?", username).First(&user)

	var settings []model.UserNotificationSetting
	h.db.Where("user_id = ?", user.ID).Find(&settings)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data":    settings,
	})
}

// UpdateUserSetting 更新用户通知设置
func (h *NotificationHandler) UpdateUserSetting(c *gin.Context) {
	username, _ := c.Get("username")
	var user model.User
	h.db.Where("username = ?", username).First(&user)

	var req model.UserNotificationSetting
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	req.UserID = user.ID

	var setting model.UserNotificationSetting
	err := h.db.Where("user_id = ? AND scene_type = ?", user.ID, req.SceneType).First(&setting).Error
	if err == gorm.ErrRecordNotFound {
		if err := h.db.Create(&req).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 500, "message": "创建失败"})
			return
		}
	} else {
		setting.EnableDingTalk = req.EnableDingTalk
		setting.EnableWechat = req.EnableWechat
		setting.EnableInternal = req.EnableInternal
		if err := h.db.Save(&setting).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 500, "message": "更新失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "保存成功",
	})
}

// BindWechatUser 绑定微信用户
func (h *NotificationHandler) BindWechatUser(c *gin.Context) {
	username, _ := c.Get("username")
	var user model.User
	h.db.Where("username = ?", username).First(&user)

	var req struct {
		OpenID string `json:"openid" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	binding := model.WechatUserBinding{
		UserID:          user.ID,
		OpenID:          req.OpenID,
		SubscribeStatus: true,
		SubscribeScenes: "[]",
	}

	if err := h.db.Create(&binding).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "绑定失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "绑定成功",
	})
}

// GetWechatBinding 获取微信绑定信息
func (h *NotificationHandler) GetWechatBinding(c *gin.Context) {
	username, _ := c.Get("username")
	var user model.User
	h.db.Where("username = ?", username).First(&user)

	var binding model.WechatUserBinding
	err := h.db.Where("user_id = ?", user.ID).First(&binding).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "未绑定",
			"data":    nil,
		})
		return
	} else if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data":    binding,
	})
}

// UpdateSubscribeScenes 更新订阅场景
func (h *NotificationHandler) UpdateSubscribeScenes(c *gin.Context) {
	username, _ := c.Get("username")
	var user model.User
	h.db.Where("username = ?", username).First(&user)

	var req struct {
		Scenes []string `json:"scenes" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	scenesJSON, _ := json.Marshal(req.Scenes)

	var binding model.WechatUserBinding
	if err := h.db.Where("user_id = ?", user.ID).First(&binding).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "未绑定微信"})
		return
	}

	binding.SubscribeScenes = string(scenesJSON)
	if err := h.db.Save(&binding).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
	})
}

// GetScheduleConfigs 获取所有调度配置
func (h *NotificationHandler) GetScheduleConfigs(c *gin.Context) {
	var configs []model.SysConfig
	if err := h.db.Where("config_key LIKE ?", "schedule_%").Find(&configs).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data":    configs,
	})
}

// UpdateScheduleConfig 更新调度配置
func (h *NotificationHandler) UpdateScheduleConfig(c *gin.Context) {
	var req struct {
		ConfigKey   string `json:"config_key" binding:"required"`
		ConfigValue string `json:"config_value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 更新数据库配置
	var config model.SysConfig
	if err := h.db.Where("config_key = ?", req.ConfigKey).First(&config).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "配置不存在"})
		return
	}

	config.ConfigValue = req.ConfigValue
	if err := h.db.Save(&config).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "保存失败"})
		return
	}

	// 更新调度器中的任务
	taskName := h.getTaskNameFromConfigKey(req.ConfigKey)
	if taskName != "" && h.scheduler != nil {
		if err := h.scheduler.UpdateSchedule(taskName, req.ConfigValue); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 500, "message": "更新调度失败: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "配置已更新并生效",
	})
}

// getTaskNameFromConfigKey 从配置键获取任务名称
func (h *NotificationHandler) getTaskNameFromConfigKey(configKey string) string {
	mapping := map[string]string{
		"schedule_daily_report_cron":   "daily_report",
		"schedule_overdue_check_cron":  "overdue_reminder",
		"schedule_weekly_report_cron":  "weekly_report",
		"schedule_monthly_report_cron": "monthly_report",
	}
	return mapping[configKey]
}
