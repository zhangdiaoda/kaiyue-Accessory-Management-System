package notification

import (
	"time"
)

// ProviderType 通知提供者类型
type ProviderType string

const (
	ProviderDingTalk ProviderType = "dingtalk"
	ProviderWechat   ProviderType = "wechat"
	ProviderInternal ProviderType = "internal"
)

// SceneType 通知场景类型
type SceneType string

const (
	SceneStockWarning       SceneType = "stock_warning"       // 库存预警
	SceneBorrowCreated      SceneType = "borrow_created"      // 领用通知
	SceneReturnReminder     SceneType = "return_reminder"     // 归还提醒
	SceneReturnCreated      SceneType = "return_created"      // 归还通知
	SceneRestock            SceneType = "restock"             // 补货通知
	SceneDailyReport        SceneType = "daily_report"        // 每日报表
	SceneWeeklyReport       SceneType = "weekly_report"       // 周报推送
	SceneMonthlyReport      SceneType = "monthly_report"      // 月报推送
	SceneSystemAnnouncement SceneType = "system_announcement" // 系统公告
)

// Notification 通知消息结构
type Notification struct {
	Scene      SceneType              `json:"scene"`       // 通知场景
	Title      string                 `json:"title"`       // 标题
	Content    string                 `json:"content"`     // 内容
	ReceiverID uint                   `json:"receiver_id"` // 接收人ID
	OpenID     string                 `json:"openid"`      // 微信OpenID
	Extra      map[string]interface{} `json:"extra"`       // 额外数据
}

// NotificationProvider 通知提供者接口
type NotificationProvider interface {
	// Send 发送通知
	Send(notification *Notification) error

	// GetType 获取通知类型
	GetType() ProviderType

	// Validate 验证配置
	Validate() error
}

// QueueTask 队列任务
type QueueTask struct {
	Provider     ProviderType  `json:"provider"`
	Notification *Notification `json:"notification"`
	RetryCount   int           `json:"retry_count"`
	ScheduledAt  time.Time     `json:"scheduled_at"`
}
