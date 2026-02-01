package service

import (
	"fmt"
	"time"
	"warehouse/internal/model"

	"gorm.io/gorm"
)

// AuditService 审计日志服务
type AuditService struct {
	db *gorm.DB
}

// NewAuditService 创建审计日志服务
func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{db: db}
}

// RecordOperation 记录操作日志
func (s *AuditService) RecordOperation(log *model.OperationLog) error {
	return s.db.Create(log).Error
}

// GetOperationLogs 查询操作日志
type LogFilter struct {
	UserID    *uint
	Username  *string
	Operation *string
	Module    *string
	Status    *string
	StartTime *time.Time
	EndTime   *time.Time
	PageNum   int
	PageSize  int
}

type LogListResult struct {
	List  []model.OperationLog `json:"list"`
	Total int64                `json:"total"`
}

func (s *AuditService) GetOperationLogs(filter LogFilter) (*LogListResult, error) {
	query := s.db.Model(&model.OperationLog{})

	// 应用过滤条件
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Username != nil && *filter.Username != "" {
		query = query.Where("username LIKE ?", "%"+*filter.Username+"%")
	}
	if filter.Operation != nil && *filter.Operation != "" {
		query = query.Where("operation = ?", *filter.Operation)
	}
	if filter.Module != nil && *filter.Module != "" {
		query = query.Where("module = ?", *filter.Module)
	}
	if filter.Status != nil && *filter.Status != "" {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.StartTime != nil {
		query = query.Where("created_at >= ?", *filter.StartTime)
	}
	if filter.EndTime != nil {
		query = query.Where("created_at <= ?", *filter.EndTime)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页
	if filter.PageNum < 1 {
		filter.PageNum = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 20
	}
	offset := (filter.PageNum - 1) * filter.PageSize

	var logs []model.OperationLog
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&logs).Error; err != nil {
		return nil, err
	}

	return &LogListResult{
		List:  logs,
		Total: total,
	}, nil
}

// GetLogDetail 获取日志详情
func (s *AuditService) GetLogDetail(logID uint) (*model.OperationLog, error) {
	var log model.OperationLog
	if err := s.db.First(&log, logID).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

// DeleteOldLogs 删除过期日志（保留最近N天）
func (s *AuditService) DeleteOldLogs(daysToKeep int) (int64, error) {
	cutoffTime := time.Now().AddDate(0, 0, -daysToKeep)
	result := s.db.Where("created_at < ?", cutoffTime).Delete(&model.OperationLog{})
	return result.RowsAffected, result.Error
}

// GetOperationStats 获取操作统计
type OperationStats struct {
	TotalOperations int64            `json:"total_operations"`
	SuccessCount    int64            `json:"success_count"`
	FailedCount     int64            `json:"failed_count"`
	ByModule        map[string]int64 `json:"by_module"`
	ByOperation     map[string]int64 `json:"by_operation"`
}

func (s *AuditService) GetOperationStats(startTime, endTime time.Time) (*OperationStats, error) {
	stats := &OperationStats{
		ByModule:    make(map[string]int64),
		ByOperation: make(map[string]int64),
	}

	query := s.db.Model(&model.OperationLog{}).
		Where("created_at BETWEEN ? AND ?", startTime, endTime)

	// 总操作数
	if err := query.Count(&stats.TotalOperations).Error; err != nil {
		return nil, err
	}

	// 成功/失败数
	if err := s.db.Model(&model.OperationLog{}).
		Where("created_at BETWEEN ? AND ? AND status = ?", startTime, endTime, "SUCCESS").
		Count(&stats.SuccessCount).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&model.OperationLog{}).
		Where("created_at BETWEEN ? AND ? AND status = ?", startTime, endTime, "FAILED").
		Count(&stats.FailedCount).Error; err != nil {
		return nil, err
	}

	// 按模块统计
	type ModuleCount struct {
		Module string
		Count  int64
	}
	var moduleCounts []ModuleCount
	if err := s.db.Model(&model.OperationLog{}).
		Select("module, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Group("module").
		Scan(&moduleCounts).Error; err != nil {
		return nil, err
	}
	for _, mc := range moduleCounts {
		stats.ByModule[mc.Module] = mc.Count
	}

	// 按操作类型统计
	type OperationCount struct {
		Operation string
		Count     int64
	}
	var operationCounts []OperationCount
	if err := s.db.Model(&model.OperationLog{}).
		Select("operation, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Group("operation").
		Scan(&operationCounts).Error; err != nil {
		return nil, err
	}
	for _, oc := range operationCounts {
		stats.ByOperation[oc.Operation] = oc.Count
	}

	return stats, nil
}

// ClearAllLogs 清空所有操作日志
func (s *AuditService) ClearAllLogs() (int64, error) {
	result := s.db.Unscoped().Where("1 = 1").Delete(&model.OperationLog{})
	return result.RowsAffected, result.Error
}

// GetTableSize 获取日志表大小(MB)
func (s *AuditService) GetTableSize() (float64, error) {
	var result struct {
		SizeMB float64 `gorm:"column:size_mb"`
	}

	err := s.db.Raw(`
		SELECT 
			ROUND(((data_length + index_length) / 1024 / 1024), 2) AS size_mb
		FROM information_schema.TABLES
		WHERE table_schema = DATABASE()
		  AND table_name = 'sys_operation_log'
	`).Scan(&result).Error

	if err != nil {
		return 0, err
	}
	return result.SizeMB, nil
}

// CleanupBySize 按大小清理日志(保留最近的记录)
func (s *AuditService) CleanupBySize(maxSizeMB float64, keepCount int) (int64, error) {
	// 先检查当前大小
	currentSize, err := s.GetTableSize()
	if err != nil {
		return 0, err
	}

	// 如果未超过阈值,不处理
	if currentSize <= maxSizeMB {
		return 0, nil
	}

	// 获取总记录数
	var totalCount int64
	if err := s.db.Model(&model.OperationLog{}).Count(&totalCount).Error; err != nil {
		return 0, err
	}

	// 计算需要删除的记录数
	deleteCount := totalCount - int64(keepCount)
	if deleteCount <= 0 {
		return 0, nil
	}

	// 获取需要保留的最小ID(最新的keepCount条记录)
	var minIDToKeep uint
	if err := s.db.Model(&model.OperationLog{}).
		Select("id").
		Order("created_at DESC").
		Limit(keepCount).
		Offset(keepCount-1).
		Pluck("id", &minIDToKeep).Error; err != nil {
		return 0, err
	}

	// 删除旧记录
	result := s.db.Unscoped().Where("id < ?", minIDToKeep).Delete(&model.OperationLog{})
	return result.RowsAffected, result.Error
}

// GetCleanupConfig 获取清理配置
type CleanupConfig struct {
	Enabled       bool   `json:"enabled"`
	Schedule      string `json:"schedule"`       // daily, weekly, monthly
	SizeThreshold int    `json:"size_threshold"` // MB
	DaysToKeep    int    `json:"days_to_keep"`
	KeepCount     int    `json:"keep_count"`      // 按大小清理时保留的记录数
	LastCleanupAt string `json:"last_cleanup_at"` // 最后清理时间
}

func (s *AuditService) GetCleanupConfig() (*CleanupConfig, error) {
	config := &CleanupConfig{
		Enabled:       false,
		Schedule:      "weekly",
		SizeThreshold: 1024,
		DaysToKeep:    30,
		KeepCount:     10000,
	}

	// 从sys_config表读取配置
	var configs []model.SysConfig
	if err := s.db.Where("config_key LIKE 'log_cleanup_%'").Find(&configs).Error; err != nil {
		return config, nil // 返回默认配置
	}

	for _, c := range configs {
		switch c.ConfigKey {
		case "log_cleanup_enabled":
			config.Enabled = c.ConfigValue == "true"
		case "log_cleanup_schedule":
			config.Schedule = c.ConfigValue
		case "log_cleanup_size_threshold":
			fmt.Sscanf(c.ConfigValue, "%d", &config.SizeThreshold)
		case "log_cleanup_days_to_keep":
			fmt.Sscanf(c.ConfigValue, "%d", &config.DaysToKeep)
		case "log_cleanup_keep_count":
			fmt.Sscanf(c.ConfigValue, "%d", &config.KeepCount)
		case "log_cleanup_last_at":
			config.LastCleanupAt = c.ConfigValue
		}
	}

	return config, nil
}

// SaveCleanupConfig 保存清理配置
func (s *AuditService) SaveCleanupConfig(config *CleanupConfig) error {
	configs := map[string]string{
		"log_cleanup_enabled":        fmt.Sprintf("%t", config.Enabled),
		"log_cleanup_schedule":       config.Schedule,
		"log_cleanup_size_threshold": fmt.Sprintf("%d", config.SizeThreshold),
		"log_cleanup_days_to_keep":   fmt.Sprintf("%d", config.DaysToKeep),
		"log_cleanup_keep_count":     fmt.Sprintf("%d", config.KeepCount),
	}

	for key, value := range configs {
		var existingConfig model.SysConfig
		err := s.db.Where("config_key = ?", key).First(&existingConfig).Error

		if err == gorm.ErrRecordNotFound {
			// 创建新配置
			s.db.Create(&model.SysConfig{
				ConfigKey:   key,
				ConfigValue: value,
			})
		} else if err == nil {
			// 更新现有配置
			s.db.Model(&existingConfig).Update("config_value", value)
		}
	}

	return nil
}

// UpdateLastCleanupTime 更新最后清理时间
func (s *AuditService) UpdateLastCleanupTime() error {
	now := time.Now().Format(time.RFC3339)
	var config model.SysConfig
	err := s.db.Where("config_key = ?", "log_cleanup_last_at").First(&config).Error

	if err == gorm.ErrRecordNotFound {
		return s.db.Create(&model.SysConfig{
			ConfigKey:   "log_cleanup_last_at",
			ConfigValue: now,
		}).Error
	}

	return s.db.Model(&config).Update("config_value", now).Error
}
