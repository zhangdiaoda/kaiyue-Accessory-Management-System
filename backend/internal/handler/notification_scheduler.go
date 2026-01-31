package handler

import (
	"fmt"
	"log"
	"sync"
	"warehouse/internal/model"

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
