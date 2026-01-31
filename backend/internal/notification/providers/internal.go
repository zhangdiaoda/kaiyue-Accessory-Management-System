package providers

import (
	"warehouse/internal/model"
	"warehouse/internal/notification"

	"gorm.io/gorm"
)

// InternalProvider 站内信通知提供者
type InternalProvider struct {
	db *gorm.DB
}

// NewInternalProvider 创建站内信通知提供者
func NewInternalProvider(db *gorm.DB) *InternalProvider {
	return &InternalProvider{
		db: db,
	}
}

// GetType 获取通知类型
func (p *InternalProvider) GetType() notification.ProviderType {
	return notification.ProviderInternal
}

// Validate 验证配置
func (p *InternalProvider) Validate() error {
	return nil
}

// Send 发送站内信
func (p *InternalProvider) Send(notif *notification.Notification) error {
	message := &model.InternalMessage{
		SenderID:   1, // 系统管理员ID,可以配置
		ReceiverID: notif.ReceiverID,
		Title:      notif.Title,
		Content:    notif.Content,
		IsRead:     false,
	}

	return p.db.Create(message).Error
}

// SendToAll 发送给所有用户
func (p *InternalProvider) SendToAll(notif *notification.Notification) error {
	message := &model.InternalMessage{
		SenderID:   1,
		ReceiverID: 0, // 0表示全体成员
		Title:      notif.Title,
		Content:    notif.Content,
		IsRead:     false,
	}

	return p.db.Create(message).Error
}
