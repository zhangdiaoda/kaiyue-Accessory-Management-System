package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"warehouse/internal/notification"
)

// WechatProvider 微信公众号通知提供者
type WechatProvider struct {
	AppID      string `json:"app_id"`
	AppSecret  string `json:"app_secret"`
	TemplateID string `json:"template_id"`
	tokenCache TokenCache
}

// WechatConfig 微信配置
type WechatConfig struct {
	AppID      string     `json:"app_id"`
	AppSecret  string     `json:"app_secret"`
	TemplateID string     `json:"template_id"`
	TokenCache TokenCache `json:"-"`
}

// TokenCache Access Token缓存接口
type TokenCache interface {
	Get(key string) (string, error)
	Set(key string, value string, expiration time.Duration) error
}

// NewWechatProvider 创建微信通知提供者
func NewWechatProvider(config *WechatConfig) *WechatProvider {
	return &WechatProvider{
		AppID:      config.AppID,
		AppSecret:  config.AppSecret,
		TemplateID: config.TemplateID,
		tokenCache: config.TokenCache,
	}
}

// GetType 获取通知类型
func (p *WechatProvider) GetType() notification.ProviderType {
	return notification.ProviderWechat
}

// Validate 验证配置
func (p *WechatProvider) Validate() error {
	if p.AppID == "" {
		return fmt.Errorf("app_id不能为空")
	}
	if p.AppSecret == "" {
		return fmt.Errorf("app_secret不能为空")
	}
	if p.TemplateID == "" {
		return fmt.Errorf("template_id不能为空")
	}
	return nil
}

// Send 发送订阅通知
func (p *WechatProvider) Send(notif *notification.Notification) error {
	// 1. 获取Access Token
	token, err := p.getAccessToken()
	if err != nil {
		return fmt.Errorf("获取access_token失败: %w", err)
	}

	// 2. 构建API URL
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/message/subscribe/bizsend?access_token=%s",
		token,
	)

	// 3. 构建消息数据
	data := p.buildMessageData(notif)

	// 4. 发送请求
	return p.sendRequest(url, data)
}

// getAccessToken 获取Access Token(带缓存)
func (p *WechatProvider) getAccessToken() (string, error) {
	cacheKey := "wechat:access_token"

	// 1. 尝试从缓存获取
	if p.tokenCache != nil {
		token, err := p.tokenCache.Get(cacheKey)
		if err == nil && token != "" {
			return token, nil
		}
	}

	// 2. 请求微信API获取
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		p.AppID, p.AppSecret,
	)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	// 检查错误
	if errcode, ok := result["errcode"]; ok && errcode != nil {
		errmsg := result["errmsg"].(string)
		return "", fmt.Errorf("微信返回错误: %s (错误码: %v)", errmsg, errcode)
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("获取access_token失败")
	}

	// 3. 缓存Token(7000秒,提前200秒过期)
	if p.tokenCache != nil {
		p.tokenCache.Set(cacheKey, token, 7000*time.Second)
	}

	return token, nil
}

// buildMessageData 构建消息数据
func (p *WechatProvider) buildMessageData(notif *notification.Notification) map[string]interface{} {
	// 截取内容长度(thing类型最多20个字符)
	content := notif.Content
	if len([]rune(content)) > 20 {
		content = string([]rune(content)[:17]) + "..."
	}

	title := notif.Title
	if len([]rune(title)) > 20 {
		title = string([]rune(title)[:17]) + "..."
	}

	return map[string]interface{}{
		"touser":      notif.OpenID,
		"template_id": p.TemplateID,
		"data": map[string]interface{}{
			"thing1": map[string]string{"value": title},
			"time2":  map[string]string{"value": time.Now().Format("2006-01-02 15:04")},
			"thing3": map[string]string{"value": content},
		},
	}
}

// sendRequest 发送HTTP请求
func (p *WechatProvider) sendRequest(url string, data map[string]interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("消息序列化失败: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查微信返回的错误码
	if errcode, ok := result["errcode"].(float64); ok && errcode != 0 {
		errmsg := result["errmsg"].(string)
		return fmt.Errorf("微信返回错误: %s (错误码: %.0f)", errmsg, errcode)
	}

	return nil
}

// SendTestMessage 发送测试消息(需要有效的OpenID)
func (p *WechatProvider) SendTestMessage(openID string) error {
	testNotif := &notification.Notification{
		Scene:   notification.SceneSystemAnnouncement,
		Title:   "测试消息",
		Content: "配置正确",
		OpenID:  openID,
	}
	return p.Send(testNotif)
}
