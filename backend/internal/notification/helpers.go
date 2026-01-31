package notification

import (
	"fmt"
	"time"
)

// BuildStockWarningNotification 构建库存预警通知
func BuildStockWarningNotification(partName, partNo string, currentStock, minStock int) *Notification {
	return &Notification{
		Scene:   SceneStockWarning,
		Title:   fmt.Sprintf("⚠️ 库存预警:%s", partName),
		Content: fmt.Sprintf("配件【%s】(编号:%s)库存不足!\n当前库存: %d\n最低阈值: %d\n请及时补货!", partName, partNo, currentStock, minStock),
		Extra: map[string]interface{}{
			"part_name":     partName,
			"part_no":       partNo,
			"current_stock": currentStock,
			"min_stock":     minStock,
		},
	}
}

// BuildBorrowNotification 构建领用通知
func BuildBorrowNotification(employeeName, partName string, quantity int, borrowDate time.Time) *Notification {
	return &Notification{
		Scene:   SceneBorrowCreated,
		Title:   "📦 领用通知",
		Content: fmt.Sprintf("员工【%s】领用了配件【%s】\n数量: %d\n时间: %s", employeeName, partName, quantity, borrowDate.Format("2006-01-02 15:04:05")),
		Extra: map[string]interface{}{
			"employee_name": employeeName,
			"part_name":     partName,
			"quantity":      quantity,
			"borrow_date":   borrowDate,
		},
	}
}

// BuildReturnReminderNotification 构建归还提醒通知
func BuildReturnReminderNotification(employeeName, partName string, borrowDate time.Time, overdueDays int) *Notification {
	return &Notification{
		Scene:   SceneReturnReminder,
		Title:   "⏰ 归还提醒",
		Content: fmt.Sprintf("员工【%s】于 %s 领用的配件【%s】已超期 %d 天未归还,请及时催促归还!", employeeName, borrowDate.Format("2006-01-02"), partName, overdueDays),
		Extra: map[string]interface{}{
			"employee_name": employeeName,
			"part_name":     partName,
			"borrow_date":   borrowDate,
			"overdue_days":  overdueDays,
		},
	}
}

// BuildDailyReportNotification 构建每日报表通知
func BuildDailyReportNotification(date time.Time, borrowCount, returnCount, lowStockCount int) *Notification {
	return &Notification{
		Scene:   SceneDailyReport,
		Title:   fmt.Sprintf("📊 %s 每日报表", date.Format("2006-01-02")),
		Content: fmt.Sprintf("今日数据汇总:\n✅ 领用记录: %d 条\n✅ 归还记录: %d 条\n⚠️ 库存预警: %d 项", borrowCount, returnCount, lowStockCount),
		Extra: map[string]interface{}{
			"date":            date,
			"borrow_count":    borrowCount,
			"return_count":    returnCount,
			"low_stock_count": lowStockCount,
		},
	}
}

// BuildWeeklyReportNotification 构建周报通知
func BuildWeeklyReportNotification(startDate, endDate time.Time, borrowCount, returnCount, lowStockCount int) *Notification {
	return &Notification{
		Scene:   SceneWeeklyReport,
		Title:   fmt.Sprintf("📈 周报(%s - %s)", startDate.Format("01-02"), endDate.Format("01-02")),
		Content: fmt.Sprintf("本周数据汇总:\n✅ 领用记录: %d 条\n✅ 归还记录: %d 条\n⚠️ 库存预警: %d 项", borrowCount, returnCount, lowStockCount),
		Extra: map[string]interface{}{
			"start_date":      startDate,
			"end_date":        endDate,
			"borrow_count":    borrowCount,
			"return_count":    returnCount,
			"low_stock_count": lowStockCount,
		},
	}
}

// BuildSystemAnnouncementNotification 构建系统公告通知
func BuildSystemAnnouncementNotification(title, content string) *Notification {
	return &Notification{
		Scene:   SceneSystemAnnouncement,
		Title:   fmt.Sprintf("📢 系统公告:%s", title),
		Content: content,
		Extra:   map[string]interface{}{},
	}
}

// BuildBorrowCreatedNotification 构建领用通知(简化版,自动使用当前时间)
func BuildBorrowCreatedNotification(employeeName, partName string, quantity int) *Notification {
	return BuildBorrowNotification(employeeName, partName, quantity, time.Now())
}

// BuildDetailedReportNotification 构建明细报表通知
func BuildDetailedReportNotification(title, content string) *Notification {
	return &Notification{
		Scene:   SceneDailyReport,
		Title:   title,
		Content: content,
	}
}

// BuildReturnNotification 构建归还通知
func BuildReturnNotification(employeeName, partName string, returnQty, damagedQty int, returnDate time.Time) *Notification {
	content := fmt.Sprintf("员工【%s】归还了配件【%s】\n正常归还: %d\n损毁报废: %d\n时间: %s",
		employeeName, partName, returnQty, damagedQty, returnDate.Format("2006-01-02 15:04:05"))

	return &Notification{
		Scene:   SceneReturnCreated,
		Title:   "✅ 归还通知",
		Content: content,
		Extra: map[string]interface{}{
			"employee_name": employeeName,
			"part_name":     partName,
			"return_qty":    returnQty,
			"damaged_qty":   damagedQty,
			"return_date":   returnDate,
		},
	}
}

// BuildReturnCreatedNotification 构建归还通知(简化版)
func BuildReturnCreatedNotification(employeeName, partName string, returnQty, damagedQty int) *Notification {
	return BuildReturnNotification(employeeName, partName, returnQty, damagedQty, time.Now())
}
