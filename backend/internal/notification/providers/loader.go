package providers

import (
	"encoding/json"
	"fmt"
	"warehouse/internal/model"
	"warehouse/internal/notification"

	"gorm.io/gorm"
)

// LoadAllProviders 从数据库加载所有启用的通知渠道并注册到管理器
func LoadAllProviders(db *gorm.DB, m *notification.Manager) {
	// 先清除已有的提供者
	m.ClearProviders()

	// 1. 注册站内信提供者 (默认启用)
	m.RegisterProvider(NewInternalProvider(db))

	// 2. 加载并注册钉钉提供者
	var dingTalkConfig model.NotificationConfig
	if err := db.Where("provider_type = ? AND is_enabled = ?", "dingtalk", true).First(&dingTalkConfig).Error; err == nil {
		var dtConfig DingTalkConfig
		if err := json.Unmarshal([]byte(dingTalkConfig.ConfigData), &dtConfig); err == nil {
			m.RegisterProvider(NewDingTalkProvider(&dtConfig))
			fmt.Println("🔔 通知系统: 钉钉渠道已加载")
		}
	}

	// 3. 加载并注册微信提供者
	var wechatConfig model.NotificationConfig
	if err := db.Where("provider_type = ? AND is_enabled = ?", "wechat", true).First(&wechatConfig).Error; err == nil {
		var wxConfig WechatConfig
		if err := json.Unmarshal([]byte(wechatConfig.ConfigData), &wxConfig); err == nil {
			// 这里创建一个简单的缓存，如果需要持久化缓存可以在外部持有
			wxConfig.TokenCache = notification.NewSimpleCache()
			m.RegisterProvider(NewWechatProvider(&wxConfig))
			fmt.Println("🔔 通知系统: 微信渠道已加载")
		}
	}
}
