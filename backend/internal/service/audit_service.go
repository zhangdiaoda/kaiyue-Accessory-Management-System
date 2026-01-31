package service

import (
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
