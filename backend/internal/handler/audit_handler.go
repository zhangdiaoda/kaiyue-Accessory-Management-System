package handler

import (
	"net/http"
	"strconv"
	"time"
	"warehouse/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuditHandler 审计日志Handler
type AuditHandler struct {
	db           *gorm.DB
	auditService *service.AuditService
}

// NewAuditHandler 创建审计日志Handler
func NewAuditHandler(db *gorm.DB) *AuditHandler {
	return &AuditHandler{
		db:           db,
		auditService: service.NewAuditService(db),
	}
}

// GetOperationLogs 查询操作日志
func (h *AuditHandler) GetOperationLogs(c *gin.Context) {
	var filter service.LogFilter

	// 解析查询参数
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		userID, _ := strconv.ParseUint(userIDStr, 10, 32)
		uid := uint(userID)
		filter.UserID = &uid
	}

	if username := c.Query("username"); username != "" {
		filter.Username = &username
	}

	if operation := c.Query("operation"); operation != "" {
		filter.Operation = &operation
	}

	if module := c.Query("module"); module != "" {
		filter.Module = &module
	}

	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}

	if startTimeStr := c.Query("start_time"); startTimeStr != "" {
		if startTime, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			filter.StartTime = &startTime
		}
	}

	if endTimeStr := c.Query("end_time"); endTimeStr != "" {
		if endTime, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			filter.EndTime = &endTime
		}
	}

	filter.PageNum, _ = strconv.Atoi(c.DefaultQuery("page_num", "1"))
	filter.PageSize, _ = strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.auditService.GetOperationLogs(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询日志失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// GetLogDetail 获取日志详情
func (h *AuditHandler) GetLogDetail(c *gin.Context) {
	logIDStr := c.Param("id")
	logID, err := strconv.ParseUint(logIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的日志ID",
		})
		return
	}

	log, err := h.auditService.GetLogDetail(uint(logID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "日志不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": log,
	})
}

// GetOperationStats 获取操作统计
func (h *AuditHandler) GetOperationStats(c *gin.Context) {
	// 默认查询最近7天的数据
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -7)

	if startTimeStr := c.Query("start_time"); startTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			startTime = t
		}
	}

	if endTimeStr := c.Query("end_time"); endTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			endTime = t
		}
	}

	stats, err := h.auditService.GetOperationStats(startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取统计数据失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": stats,
	})
}

// ClearLogs 手动清空所有日志
func (h *AuditHandler) ClearLogs(c *gin.Context) {
	deletedCount, err := h.auditService.ClearAllLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "清空日志失败",
		})
		return
	}

	// 更新最后清理时间
	h.auditService.UpdateLastCleanupTime()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "清空成功",
		"data": gin.H{
			"deleted_count": deletedCount,
		},
	})
}

// GetLogSize 获取日志表大小
func (h *AuditHandler) GetLogSize(c *gin.Context) {
	sizeMB, err := h.auditService.GetTableSize()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取日志大小失败",
		})
		return
	}

	// 获取日志总数
	var totalCount int64
	h.db.Model(&service.LogFilter{}).Count(&totalCount)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"size_mb":     sizeMB,
			"total_count": totalCount,
		},
	})
}

// GetCleanupConfig 获取清理配置
func (h *AuditHandler) GetCleanupConfig(c *gin.Context) {
	config, err := h.auditService.GetCleanupConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取配置失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": config,
	})
}

// UpdateCleanupConfig 更新清理配置
func (h *AuditHandler) UpdateCleanupConfig(c *gin.Context) {
	var config service.CleanupConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	if err := h.auditService.SaveCleanupConfig(&config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "保存配置失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "配置已保存",
	})
}
