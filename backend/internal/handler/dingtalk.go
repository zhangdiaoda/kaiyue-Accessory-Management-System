package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DingTalkHandler struct {
}

func NewDingTalkHandler() *DingTalkHandler {
	return &DingTalkHandler{}
}

// DingTalkMessage 钉钉消息结构
type DingTalkMessage struct {
	MsgType  string          `json:"msgtype"`
	Markdown MarkdownContent `json:"markdown"`
}

type MarkdownContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// SendReport 发送报表到钉钉
func (h *DingTalkHandler) SendReport(c *gin.Context) {
	var req struct {
		WebhookURL string `json:"webhook_url" binding:"required"`
		Title      string `json:"title" binding:"required"`
		Content    string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 构建钉钉消息
	message := DingTalkMessage{
		MsgType: "markdown",
		Markdown: MarkdownContent{
			Title: req.Title,
			Text:  req.Content,
		},
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "消息构建失败",
		})
		return
	}

	// 发送到钉钉
	resp, err := http.Post(req.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "发送失败: " + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// 读取响应
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
		// 检查钉钉返回的错误码
		if errcode, ok := result["errcode"].(float64); ok && errcode != 0 {
			errmsg := result["errmsg"].(string)
			c.JSON(http.StatusOK, gin.H{
				"code":    500,
				"message": fmt.Sprintf("钉钉返回错误: %s (错误码: %.0f)", errmsg, errcode),
			})
			return
		}
	}

	if resp.StatusCode != 200 {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": fmt.Sprintf("钉钉返回HTTP错误: %d", resp.StatusCode),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "发送成功",
	})
}

// TestWebhook 测试Webhook连接
func (h *DingTalkHandler) TestWebhook(c *gin.Context) {
	var req struct {
		WebhookURL string `json:"webhook_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 发送测试消息
	message := DingTalkMessage{
		MsgType: "markdown",
		Markdown: MarkdownContent{
			Title: "测试消息",
			Text:  "## 仓储管理系统\n\n这是一条测试消息，Webhook配置正确！",
		},
	}

	jsonData, _ := json.Marshal(message)
	resp, err := http.Post(req.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "连接失败",
		})
		return
	}
	defer resp.Body.Close()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "测试成功",
	})
}
