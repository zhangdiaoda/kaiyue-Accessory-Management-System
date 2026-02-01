package handler

import (
	"fmt"
	"log"
	"sync"
	"warehouse/internal/model"
	"warehouse/internal/service"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// NotificationScheduler 通知定时任务调度器
type NotificationScheduler struct {
	ni      *NotificationIntegrator
	db      *gorm.DB
	cron    *cron.Cron
	lock    sync.Mutex
	entries map[string]cron.EntryID
}

// NewNotificationScheduler 创建通知调度器
func NewNotificationScheduler(ni *NotificationIntegrator) *NotificationScheduler {
	c := cron.New(cron.WithSeconds())
	return &NotificationScheduler{
		ni:      ni,
		db:      ni.db,
		cron:    c,
		entries: make(map[string]cron.EntryID),
	}
}

// Start 启动调度器并注册默认任务
func (s *NotificationScheduler) Start() {
	s.lock.Lock()
	defer s.lock.Unlock()

	// 初始化默认配置
	s.initDefaultScheduleConfig()

	// 从数据库读取配置并注册任务
	s.registerTasksFromConfig()

	s.cron.Start()
	log.Println("🚀 通知定时任务系统已启动")
}

// initDefaultScheduleConfig 初始化默认调度配置
func (s *NotificationScheduler) initDefaultScheduleConfig() {
	defaults := map[string]string{
		"schedule_daily_report_cron":   "0 30 8 * * *", // 每天 8:30
		"schedule_overdue_check_cron":  "0 0 9 * * *",  // 每天 9:00
		"schedule_weekly_report_cron":  "0 0 9 * * 1",  // 每周一 9:00
		"schedule_monthly_report_cron": "0 0 9 1 * *",  // 每月1号 9:00
		"schedule_log_cleanup_cron":    "0 0 2 * * *",  // 每天凌晨2点清理日志
	}

	for key, value := range defaults {
		var config model.SysConfig
		err := s.db.Where("config_key = ?", key).First(&config).Error
		if err == gorm.ErrRecordNotFound {
			s.db.Create(&model.SysConfig{
				ConfigKey:   key,
				ConfigValue: value,
				Description: s.getDescriptionForKey(key),
			})
		}
	}
}

// getDescriptionForKey 获取配置项描述
func (s *NotificationScheduler) getDescriptionForKey(key string) string {
	descriptions := map[string]string{
		"schedule_daily_report_cron":   "每日报表推送时间（Cron表达式）",
		"schedule_overdue_check_cron":  "超期检查时间（Cron表达式）",
		"schedule_weekly_report_cron":  "周报推送时间（Cron表达式）",
		"schedule_monthly_report_cron": "月报推送时间（Cron表达式）",
	}
	return descriptions[key]
}

// registerTasksFromConfig 从数据库配置注册任务
func (s *NotificationScheduler) registerTasksFromConfig() {
	// 注册每日报表
	if cronExpr := s.getCronConfig("schedule_daily_report_cron"); cronExpr != "" {
		s.registerTask("daily_report", cronExpr, func() {
			log.Println("⏰ 执行定时任务: 发送每日报表")
			s.ni.SendDailyReport()
		})
	}

	// 注册超期检查
	if cronExpr := s.getCronConfig("schedule_overdue_check_cron"); cronExpr != "" {
		s.registerTask("overdue_reminder", cronExpr, func() {
			log.Println("⏰ 执行定时任务: 检查超期未归还")
			s.ni.CheckAndNotifyOverdueReturn()
		})
	}

	// 注册周报（如果有实现）
	if cronExpr := s.getCronConfig("schedule_weekly_report_cron"); cronExpr != "" {
		s.registerTask("weekly_report", cronExpr, func() {
			log.Println("⏰ 执行定时任务: 发送周报")
			s.ni.SendWeeklyReport()
		})
	}

	// 注册月报
	if cronExpr := s.getCronConfig("schedule_monthly_report_cron"); cronExpr != "" {
		s.registerTask("monthly_report", cronExpr, func() {
			log.Println("⏰ 执行定时任务: 发送月报")
			s.ni.SendMonthlyReport()
		})
	}

	// 注册日志清理任务
	if cronExpr := s.getCronConfig("schedule_log_cleanup_cron"); cronExpr != "" {
		s.registerTask("log_cleanup", cronExpr, func() {
			s.performLogCleanup()
		})
	}
}

// getCronConfig 从数据库获取Cron配置
func (s *NotificationScheduler) getCronConfig(key string) string {
	var config model.SysConfig
	if err := s.db.Where("config_key = ?", key).First(&config).Error; err == nil {
		return config.ConfigValue
	}
	return ""
}

// registerTask 注册单个定时任务
func (s *NotificationScheduler) registerTask(name, cronExpr string, job func()) {
	id, err := s.cron.AddFunc(cronExpr, job)
	if err == nil {
		s.entries[name] = id
		log.Printf("✅ 注册定时任务成功: %s -> %s", name, cronExpr)
	} else {
		log.Printf("❌ 注册定时任务失败 [%s]: %v", name, err)
	}
}

// UpdateSchedule 更新定时任务配置
func (s *NotificationScheduler) UpdateSchedule(taskName, cronExpr string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	// 验证cron表达式
	if _, err := cron.ParseStandard(cronExpr); err != nil {
		return fmt.Errorf("无效的Cron表达式: %v", err)
	}

	// 移除旧任务
	if oldID, exists := s.entries[taskName]; exists {
		s.cron.Remove(oldID)
		delete(s.entries, taskName)
	}

	// 根据任务名注册新任务
	var job func()
	switch taskName {
	case "daily_report":
		job = func() {
			log.Println("⏰ 执行定时任务: 发送每日报表")
			s.ni.SendDailyReport()
		}
	case "overdue_reminder":
		job = func() {
			log.Println("⏰ 执行定时任务: 检查超期未归还")
			s.ni.CheckAndNotifyOverdueReturn()
		}
	case "weekly_report":
		job = func() {
			log.Println("⏰ 执行定时任务: 发送周报")
			s.ni.SendWeeklyReport()
		}
	case "monthly_report":
		job = func() {
			log.Println("⏰ 执行定时任务: 发送月报")
			s.ni.SendMonthlyReport()
		}
	case "log_cleanup":
		job = func() {
			s.performLogCleanup()
		}
	default:
		return fmt.Errorf("未知的任务名称: %s", taskName)
	}

	// 注册新任务
	id, err := s.cron.AddFunc(cronExpr, job)
	if err != nil {
		return fmt.Errorf("注册任务失败: %v", err)
	}

	s.entries[taskName] = id
	log.Printf("🔄 更新定时任务成功: %s -> %s", taskName, cronExpr)
	return nil
}

// Stop 停止调度器
func (s *NotificationScheduler) Stop() {
	s.cron.Stop()
	log.Println("🛑 通知定时任务系统已停止")
}

// RunDailyReportNow 立即手动执行每日报表
func (s *NotificationScheduler) RunDailyReportNow() {
	go s.ni.SendDailyReport()
}

// RunOverdueCheckNow 立即手动执行超期检查
func (s *NotificationScheduler) RunOverdueCheckNow() {
	go s.ni.CheckAndNotifyOverdueReturn()
}

// performLogCleanup 执行日志清理任务
func (s *NotificationScheduler) performLogCleanup() {
	log.Println("⏰ 执行定时任务: 清理操作日志")

	// 创建审计服务
	auditService := service.NewAuditService(s.db)

	// 获取清理配置
	config, err := auditService.GetCleanupConfig()
	if err != nil || !config.Enabled {
		log.Println("⏸️  日志清理未启用或配置错误,跳过")
		return
	}

	// 1. 按大小清理(如果超过阈值)
	sizeMB, _ := auditService.GetTableSize()
	log.Printf("📊 当前日志表大小: %.2f MB, 阈值: %d MB", sizeMB, config.SizeThreshold)

	if sizeMB > float64(config.SizeThreshold) {
		deletedCount, err := auditService.CleanupBySize(float64(config.SizeThreshold), config.KeepCount)
		if err != nil {
			log.Printf("❌ 按大小清理日志失败: %v", err)
		} else {
			log.Printf("✅ 按大小清理完成,删除了 %d 条旧日志", deletedCount)
		}
	}

	// 2. 按天数清理
	if config.DaysToKeep > 0 {
		deletedCount, err := auditService.DeleteOldLogs(config.DaysToKeep)
		if err != nil {
			log.Printf("❌ 按天数清理日志失败: %v", err)
		} else if deletedCount > 0 {
			log.Printf("✅ 按天数清理完成,删除了 %d 条过期日志", deletedCount)
		}
	}

	// 更新最后清理时间
	auditService.UpdateLastCleanupTime()
}
