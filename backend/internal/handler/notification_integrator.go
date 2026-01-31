package handler

import (
	"fmt"
	"strings"
	"time"
	"warehouse/internal/model"
	"warehouse/internal/notification"

	"gorm.io/gorm"
)

// NotificationIntegrator 通知集成助手
type NotificationIntegrator struct {
	db      *gorm.DB
	manager *notification.Manager
	rg      *ReportGenerator
}

// NewNotificationIntegrator 创建通知集成助手
func NewNotificationIntegrator(db *gorm.DB, manager *notification.Manager) *NotificationIntegrator {
	baseURL := "http://localhost:8080" // 可配置
	return &NotificationIntegrator{
		db:      db,
		manager: manager,
		rg:      NewReportGenerator(db, baseURL),
	}
}

// sendToAllRegistered 发送通知到所有已注册且订阅了该场景的渠道
func (ni *NotificationIntegrator) sendToAllRegistered(scene notification.SceneType, notif *notification.Notification) {
	providers := ni.manager.GetProvidersByScene(scene)
	if len(providers) > 0 {
		ni.manager.SendNotificationAsync(notif, providers)
	}
}

// CheckAndNotifyStockWarning 检查并发送库存预警通知
func (ni *NotificationIntegrator) CheckAndNotifyStockWarning(partID uint) {
	if ni.manager == nil {
		return
	}

	var part model.Part
	if err := ni.db.First(&part, partID).Error; err != nil {
		return
	}

	if part.StockQuantity < part.WarningThreshold {
		notif := notification.BuildStockWarningNotification(
			part.Name,
			part.PartNo,
			part.StockQuantity,
			part.WarningThreshold,
		)
		ni.sendToAllRegistered(notification.SceneStockWarning, notif)
	}
}

// NotifyBorrowCreated 发送领用创建通知
func (ni *NotificationIntegrator) NotifyBorrowCreated(employeeName, partName string, quantity int) {
	if ni.manager == nil {
		return
	}

	notif := notification.BuildBorrowCreatedNotification(employeeName, partName, quantity)
	ni.sendToAllRegistered(notification.SceneBorrowCreated, notif)
}

// SendDailyReport 发送每日报表
func (ni *NotificationIntegrator) SendDailyReport() {
	if ni.manager == nil {
		return
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var borrowCount int64
	ni.db.Model(&model.BorrowRecord{}).Where("borrow_time >= ?", today).Count(&borrowCount)

	var returnCount int64
	ni.db.Model(&model.BorrowRecord{}).Where("return_time >= ?", today).Count(&returnCount)

	var lowStockParts []model.Part
	ni.db.Where("stock_quantity < warning_threshold").Limit(10).Find(&lowStockParts)
	lowStockCount := len(lowStockParts)

	// 生成Excel报表
	var excelLink string
	if ni.rg != nil {
		if filename, err := ni.rg.GenerateDailyReportExcel(now); err == nil {
			excelLink = ni.rg.GetDownloadURL(filename)
		}
	}

	// 构建 Markdown 表格内容
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("### 📊 %s 每日业务概览\n\n", now.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("- ✅ 今日领用: **%d** 条\n", borrowCount))
	sb.WriteString(fmt.Sprintf("- ✅ 今日归还: **%d** 条\n", returnCount))
	sb.WriteString(fmt.Sprintf("- ⚠️ 库存预警: **%d** 项\n\n", lowStockCount))

	if lowStockCount > 0 {
		sb.WriteString("#### 🚨 低库存配件 Top 10\n")
		sb.WriteString("| 配件名称 | 当前库存 | 预警阈值 |\n")
		sb.WriteString("| :--- | :--- | :--- |\n")
		for _, p := range lowStockParts {
			sb.WriteString(fmt.Sprintf("| %s | %d | %d |\n", p.Name, p.StockQuantity, p.WarningThreshold))
		}
		sb.WriteString("\n")
	}

	// 添加Excel下载链接
	if excelLink != "" {
		sb.WriteString(fmt.Sprintf("📥 [下载Excel详细报表](%s)\n\n", excelLink))
		sb.WriteString("> 报表有效期：7天\n")
	}

	sb.WriteString(fmt.Sprintf("> 自动推送时间: %s", now.Format("15:04:05")))

	notif := &notification.Notification{
		Scene:   notification.SceneDailyReport,
		Title:   fmt.Sprintf("📊 每日报表 (%s)", now.Format("01-02")),
		Content: sb.String(),
	}
	ni.sendToAllRegistered(notification.SceneDailyReport, notif)
}

// CheckAndNotifyOverdueReturn 检查并通知超期未归还
func (ni *NotificationIntegrator) CheckAndNotifyOverdueReturn() {
	if ni.manager == nil {
		return
	}

	threshold := time.Now().AddDate(0, 0, -7)

	var overdueRecords []model.BorrowRecord
	ni.db.Where("status IN ('BORROWED', 'PARTIAL_RETURNED') AND borrow_time < ?", threshold).Find(&overdueRecords)

	for _, record := range overdueRecords {
		var employee model.Employee
		var part model.Part
		ni.db.First(&employee, record.EmployeeID)
		ni.db.First(&part, record.PartID)

		overdueDays := int(time.Since(record.BorrowTime).Hours() / 24)
		notif := notification.BuildReturnReminderNotification(employee.Name, part.Name, record.BorrowTime, overdueDays)
		ni.sendToAllRegistered(notification.SceneReturnReminder, notif)
	}
}

// NotifyCustomReport 发送自定义报表
func (ni *NotificationIntegrator) NotifyCustomReport(title, content string) {
	if ni.manager == nil {
		return
	}
	notif := notification.BuildDetailedReportNotification(title, content)
	ni.sendToAllRegistered(notification.SceneDailyReport, notif)
}

// NotifyReturnCreated 发送归还创建通知
func (ni *NotificationIntegrator) NotifyReturnCreated(employeeName, partName string, returnQty, damagedQty int) {
	if ni.manager == nil {
		return
	}

	notif := notification.BuildReturnCreatedNotification(employeeName, partName, returnQty, damagedQty)
	ni.sendToAllRegistered(notification.SceneReturnCreated, notif)
}

// NotifyRestock 发送补货通知
func (ni *NotificationIntegrator) NotifyRestock(partName string, quantity int, newStock int) {
	if ni.manager == nil {
		return
	}

	notif := &notification.Notification{
		Scene:   "restock", // 新增补货场景
		Title:   "📥 补货入库通知",
		Content: fmt.Sprintf("配件【%s】已完成补货入库\n本次入库: %d\n当前总库存: %d", partName, quantity, newStock),
	}
	ni.sendToAllRegistered("restock", notif)
}

// SendWeeklyReport 发送周报
func (ni *NotificationIntegrator) SendWeeklyReport() {
	if ni.manager == nil {
		return
	}

	now := time.Now()
	// 计算本周一00:00:00
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // 周日算作7
	}
	weekStart := now.AddDate(0, 0, -(weekday - 1))
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())

	// 统计本周数据
	var borrowCount int64
	ni.db.Model(&model.BorrowRecord{}).Where("borrow_time >= ?", weekStart).Count(&borrowCount)

	var returnCount int64
	ni.db.Model(&model.BorrowRecord{}).Where("return_time >= ?", weekStart).Count(&returnCount)

	var inboundCount int64
	ni.db.Model(&model.InboundRecord{}).Where("inbound_time >= ?", weekStart).Count(&inboundCount)

	// 本周总入库数量
	var totalInbound int
	ni.db.Model(&model.InboundRecord{}).
		Where("inbound_time >= ?", weekStart).
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&totalInbound)

	// 低库存配件
	var lowStockParts []model.Part
	ni.db.Where("stock_quantity < warning_threshold").Limit(5).Find(&lowStockParts)

	// 构建周报内容
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("### 📈 %s 周报\n\n", now.Format("2006年第01周")))
	sb.WriteString(fmt.Sprintf("**统计时间**: %s ~ %s\n\n", weekStart.Format("01-02"), now.Format("01-02")))
	sb.WriteString(fmt.Sprintf("- 📦 本周领用: **%d** 笔\n", borrowCount))
	sb.WriteString(fmt.Sprintf("- ✅ 本周归还: **%d** 笔\n", returnCount))
	sb.WriteString(fmt.Sprintf("- 📥 本周入库: **%d** 笔 (数量: %d)\n", inboundCount, totalInbound))
	sb.WriteString(fmt.Sprintf("- ⚠️ 低库存配件: **%d** 项\n\n", len(lowStockParts)))

	if len(lowStockParts) > 0 {
		sb.WriteString("#### 🚨 需关注低库存配件\n")
		sb.WriteString("| 配件名称 | 当前库存 | 预警阈值 |\n")
		sb.WriteString("| :--- | :--- | :--- |\n")
		for _, p := range lowStockParts {
			sb.WriteString(fmt.Sprintf("| %s | %d | %d |\n", p.Name, p.StockQuantity, p.WarningThreshold))
		}
	}

	notif := &notification.Notification{
		Scene:   notification.SceneWeeklyReport,
		Title:   "📈 周报",
		Content: sb.String(),
	}
	ni.sendToAllRegistered(notification.SceneWeeklyReport, notif)
}

// SendMonthlyReport 发送月报
func (ni *NotificationIntegrator) SendMonthlyReport() {
	if ni.manager == nil {
		return
	}

	now := time.Now()
	// 计算本月1号00:00:00
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// 统计本月数据
	var borrowCount int64
	ni.db.Model(&model.BorrowRecord{}).Where("borrow_time >= ?", monthStart).Count(&borrowCount)

	var returnCount int64
	ni.db.Model(&model.BorrowRecord{}).Where("return_time >= ?", monthStart).Count(&returnCount)

	var inboundCount int64
	ni.db.Model(&model.InboundRecord{}).Where("inbound_time >= ?", monthStart).Count(&inboundCount)

	// 本月总入库数量
	var totalInbound int
	ni.db.Model(&model.InboundRecord{}).
		Where("inbound_time >= ?", monthStart).
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&totalInbound)

	// 本月活跃员工Top 10
	type EmployeeStat struct {
		EmployeeID   uint
		EmployeeName string
		BorrowCount  int64
	}
	var topEmployees []EmployeeStat
	ni.db.Table("borrow_records").
		Select("employee_id, COUNT(*) as borrow_count").
		Where("borrow_time >= ?", monthStart).
		Group("employee_id").
		Order("borrow_count DESC").
		Limit(10).
		Scan(&topEmployees)

	// 填充员工姓名
	for i := range topEmployees {
		var emp model.Employee
		ni.db.First(&emp, topEmployees[i].EmployeeID)
		topEmployees[i].EmployeeName = emp.Name
	}

	// 本月热门配件Top 10
	type PartStat struct {
		PartID      uint
		PartName    string
		BorrowCount int64
	}
	var topParts []PartStat
	ni.db.Table("borrow_records").
		Select("part_id, COUNT(*) as borrow_count").
		Where("borrow_time >= ?", monthStart).
		Group("part_id").
		Order("borrow_count DESC").
		Limit(10).
		Scan(&topParts)

	// 填充配件名称
	for i := range topParts {
		var part model.Part
		ni.db.First(&part, topParts[i].PartID)
		topParts[i].PartName = part.Name
	}

	// 构建月报内容
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("### 📊 %s 月度报告\n\n", now.Format("2006年01月")))
	sb.WriteString(fmt.Sprintf("**统计时间**: %s ~ %s\n\n", monthStart.Format("2006-01-02"), now.Format("2006-01-02")))

	sb.WriteString("#### 📈 业务数据总览\n")
	sb.WriteString(fmt.Sprintf("- 📦 本月领用: **%d** 笔\n", borrowCount))
	sb.WriteString(fmt.Sprintf("- ✅ 本月归还: **%d** 笔\n", returnCount))
	sb.WriteString(fmt.Sprintf("- 📥 本月入库: **%d** 笔 (数量: %d)\n\n", inboundCount, totalInbound))

	if len(topEmployees) > 0 {
		sb.WriteString("#### 👥 活跃员工 Top 10\n")
		sb.WriteString("| 排名 | 员工姓名 | 领用次数 |\n")
		sb.WriteString("| :--- | :--- | :--- |\n")
		for i, e := range topEmployees {
			sb.WriteString(fmt.Sprintf("| %d | %s | %d |\n", i+1, e.EmployeeName, e.BorrowCount))
		}
		sb.WriteString("\n")
	}

	if len(topParts) > 0 {
		sb.WriteString("#### 🔧 热门配件 Top 10\n")
		sb.WriteString("| 排名 | 配件名称 | 被领用次数 |\n")
		sb.WriteString("| :--- | :--- | :--- |\n")
		for i, p := range topParts {
			sb.WriteString(fmt.Sprintf("| %d | %s | %d |\n", i+1, p.PartName, p.BorrowCount))
		}
	}

	sb.WriteString(fmt.Sprintf("\n\u003e 生成时间: %s", now.Format("2006-01-02 15:04:05")))

	notif := &notification.Notification{
		Scene:   notification.SceneMonthlyReport,
		Title:   fmt.Sprintf("📊 %s月报", now.Format("01月")),
		Content: sb.String(),
	}
	ni.sendToAllRegistered(notification.SceneMonthlyReport, notif)
}
