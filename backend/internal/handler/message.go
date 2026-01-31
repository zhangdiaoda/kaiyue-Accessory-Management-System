package handler

import (
	"net/http"
	"warehouse/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MessageHandler struct {
	db *gorm.DB
}

func NewMessageHandler(db *gorm.DB) *MessageHandler {
	return &MessageHandler{db: db}
}

// GetMessages 获取用户的站内信列表（收件+发件同步展示）
func (h *MessageHandler) GetMessages(c *gin.Context) {
	username, _ := c.Get("username")
	var user model.User
	h.db.Where("username = ?", username).First(&user)

	type DisplayMessage struct {
		model.InternalMessage
		SenderName   string `json:"sender_name"`
		ReceiverName string `json:"receiver_name"`
	}

	var messages []DisplayMessage
	err := h.db.Table("internal_message").
		Select("internal_message.*, sender.real_name as sender_name, CASE WHEN internal_message.receiver_id = 0 THEN '全体成员' ELSE receiver.real_name END as receiver_name").
		Joins("LEFT JOIN sys_user as sender ON internal_message.sender_id = sender.id").
		Joins("LEFT JOIN sys_user as receiver ON internal_message.receiver_id = receiver.id").
		Where("internal_message.receiver_id = ? OR internal_message.receiver_id = 0 OR internal_message.sender_id = ?", user.ID, user.ID).
		Order("internal_message.created_at desc").
		Scan(&messages).Error

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "获取消息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data":    messages,
	})
}

// SendMessage 发送消息
func (h *MessageHandler) SendMessage(c *gin.Context) {
	username, _ := c.Get("username")
	var sender model.User
	h.db.Where("username = ?", username).First(&sender)

	var msg model.InternalMessage
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	msg.SenderID = sender.ID
	if err := h.db.Create(&msg).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "发送失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "消息已发送"})
}

// MarkAsRead 标记为已读
func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Model(&model.InternalMessage{}).Where("id = ?", id).Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "操作成功"})
}

// GetUnreadCount 获取未读消息数
func (h *MessageHandler) GetUnreadCount(c *gin.Context) {
	username, _ := c.Get("username")
	var user model.User
	h.db.Where("username = ?", username).First(&user)

	var count int64
	h.db.Model(&model.InternalMessage{}).Where("(receiver_id = ? OR (receiver_id = 0 AND sender_id != ?)) AND is_read = ?", user.ID, user.ID, false).Count(&count)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data":    count,
	})
}
