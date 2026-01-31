package model

import (
	"time"
)

// NotificationConfig 通知配置模型
type NotificationConfig struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	ConfigName       string    `gorm:"size:100;not null" json:"config_name"`
	ProviderType     string    `gorm:"size:20;not null" json:"provider_type"` // dingtalk/wechat/internal
	IsEnabled        bool      `gorm:"default:true" json:"is_enabled"`
	ConfigData       string    `gorm:"type:text" json:"config_data"`       // JSON格式配置
	SubscribedScenes string    `gorm:"type:text" json:"subscribed_scenes"` // 订阅场景列表 (JSON数组)
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (NotificationConfig) TableName() string {
	return "notification_config"
}

// NotificationLog 通知发送记录
type NotificationLog struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	ProviderType   string     `gorm:"size:20;not null" json:"provider_type"`
	SceneType      string     `gorm:"size:50;not null" json:"scene_type"`
	Title          string     `gorm:"size:200" json:"title"`
	Content        string     `gorm:"type:text" json:"content"`
	ReceiverID     *uint      `json:"receiver_id"`
	ReceiverOpenID string     `gorm:"size:100" json:"receiver_openid"`
	Status         string     `gorm:"size:20;default:'pending'" json:"status"` // pending/success/failed
	ErrorMsg       string     `gorm:"type:text" json:"error_msg"`
	RetryCount     int        `gorm:"default:0" json:"retry_count"`
	SentAt         *time.Time `json:"sent_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

func (NotificationLog) TableName() string {
	return "notification_log"
}

// UserNotificationSetting 用户通知设置
type UserNotificationSetting struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	UserID         uint       `gorm:"not null" json:"user_id"`
	SceneType      string     `gorm:"size:50;not null" json:"scene_type"`
	EnableDingTalk bool       `gorm:"default:false" json:"enable_dingtalk"`
	EnableWechat   bool       `gorm:"default:false" json:"enable_wechat"`
	EnableInternal bool       `gorm:"default:true" json:"enable_internal"`
	QuietStartTime *time.Time `json:"quiet_start_time"`
	QuietEndTime   *time.Time `json:"quiet_end_time"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (UserNotificationSetting) TableName() string {
	return "user_notification_setting"
}

// WechatUserBinding 微信用户绑定
type WechatUserBinding struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	UserID            uint       `gorm:"not null;uniqueIndex" json:"user_id"`
	OpenID            string     `gorm:"size:100;not null;uniqueIndex" json:"openid"`
	SubscribeStatus   bool       `gorm:"default:true" json:"subscribe_status"`
	SubscribeScenes   string     `gorm:"type:text" json:"subscribe_scenes"` // JSON数组
	BindTime          time.Time  `gorm:"autoCreateTime" json:"bind_time"`
	LastSubscribeTime *time.Time `json:"last_subscribe_time"`
}

func (WechatUserBinding) TableName() string {
	return "wechat_user_binding"
}
