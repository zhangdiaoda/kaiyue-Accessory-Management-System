package providers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"warehouse/internal/notification"
)

// DingTalkProvider 钉钉通知提供者
type DingTalkProvider struct {
	WebhookURL string   `json:"webhook_url"`
	Secret     string   `json:"secret"`     // 加签密钥
	AtMobiles  []string `json:"at_mobiles"` // @的手机号列表
	IsAtAll    bool     `json:"is_at_all"`  // 是否@所有人
}

// DingTalkConfig 钉钉配置
type DingTalkConfig struct {
	WebhookURL string   `json:"webhook_url"`
	Secret     string   `json:"secret"`
	AtMobiles  []string `json:"at_mobiles"`
	IsAtAll    bool     `json:"is_at_all"`
}

// NewDingTalkProvider 创建钉钉通知提供者
func NewDingTalkProvider(config *DingTalkConfig) *DingTalkProvider {
	return &DingTalkProvider{
		WebhookURL: config.WebhookURL,
		Secret:     config.Secret,
		AtMobiles:  config.AtMobiles,
		IsAtAll:    config.IsAtAll,
	}
}

// GetType 获取通知类型
func (p *DingTalkProvider) GetType() notification.ProviderType {
	return notification.ProviderDingTalk
}

// Validate 验证配置
func (p *DingTalkProvider) Validate() error {
	if p.WebhookURL == "" {
		return fmt.Errorf("webhook_url不能为空")
	}
	return nil
}

// Send 发送通知
func (p *DingTalkProvider) Send(notif *notification.Notification) error {
	// 1. 生成签名
	timestamp := time.Now().UnixMilli()
	sign := p.generateSign(timestamp)

	// 2. 构建请求URL
	url := fmt.Sprintf("%s&timestamp=%d&sign=%s", p.WebhookURL, timestamp, sign)

	// 3. 构建消息体
	message := p.buildMarkdownMessage(notif)

	// 4. 发送HTTP请求
	return p.sendRequest(url, message)
}

// generateSign 生成加签
func (p *DingTalkProvider) generateSign(timestamp int64) string {
	if p.Secret == "" {
		return ""
	}

	stringToSign := fmt.Sprintf("%d\n%s", timestamp, p.Secret)
	h := hmac.New(sha256.New, []byte(p.Secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// buildMarkdownMessage 构建Markdown消息
func (p *DingTalkProvider) buildMarkdownMessage(notif *notification.Notification) map[string]interface{} {
	return map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": notif.Title,
			"text":  notif.Content,
		},
		"at": map[string]interface{}{
			"atMobiles": p.AtMobiles,
			"isAtAll":   p.IsAtAll,
		},
	}
}

// buildTextMessage 构建Text消息
func (p *DingTalkProvider) buildTextMessage(notif *notification.Notification) map[string]interface{} {
	content := fmt.Sprintf("**%s**\n\n%s", notif.Title, notif.Content)
	return map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": content,
		},
		"at": map[string]interface{}{
			"atMobiles": p.AtMobiles,
			"isAtAll":   p.IsAtAll,
		},
	}
}

// sendRequest 发送HTTP请求
func (p *DingTalkProvider) sendRequest(url string, message map[string]interface{}) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("消息序列化失败: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查钉钉返回的错误码
	if errcode, ok := result["errcode"].(float64); ok && errcode != 0 {
		errmsg := result["errmsg"].(string)
		return fmt.Errorf("钉钉返回错误: %s (错误码: %.0f)", errmsg, errcode)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP错误: %d", resp.StatusCode)
	}

	return nil
}

// SendTestMessage 发送测试消息
func (p *DingTalkProvider) SendTestMessage() error {
	testNotif := &notification.Notification{
		Scene:   notification.SceneSystemAnnouncement,
		Title:   "测试消息",
		Content: "## 仓储管理系统\n\n这是一条测试消息，钉钉Webhook配置正确！✅",
	}
	return p.Send(testNotif)
}
