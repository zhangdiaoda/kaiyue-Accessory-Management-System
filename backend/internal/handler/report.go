package handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReportHandler struct {
	db *gorm.DB
	ni *NotificationIntegrator
}

func NewReportHandler(db *gorm.DB, ni *NotificationIntegrator) *ReportHandler {
	return &ReportHandler{db: db, ni: ni}
}

// GetPartReport 按配件统计
func (h *ReportHandler) GetPartReport(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	type PartStat struct {
		PartID       uint    `json:"part_id"`
		PartNo       string  `json:"part_no"`
		PartName     string  `json:"part_name"`
		BorrowCount  int64   `json:"borrow_count"`
		TotalBorrow  int64   `json:"total_borrow"`
		TotalReturn  int64   `json:"total_return"`
		TotalDamaged int64   `json:"total_damaged"`
		DamageRate   float64 `json:"damage_rate"`
	}

	var stats []PartStat

	h.db.Table("borrow_record").
		Select(`
			part.id as part_id,
			part.part_no,
			part.name as part_name,
			COUNT(*) as borrow_count,
			SUM(borrow_record.borrow_quantity) as total_borrow,
			SUM(borrow_record.return_quantity) as total_return,
			SUM(borrow_record.damaged_quantity) as total_damaged,
			ROUND(SUM(borrow_record.damaged_quantity) * 100.0 / NULLIF(SUM(borrow_record.borrow_quantity), 0), 2) as damage_rate
		`).
		Joins("LEFT JOIN part ON borrow_record.part_id = part.id").
		Where("borrow_record.borrow_time BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59").
		Group("part.id").
		Order("total_borrow DESC").
		Scan(&stats)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data":    stats,
	})
}

// GetEmployeeReport 按员工统计（支持配件筛选）
func (h *ReportHandler) GetEmployeeReport(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	partID := c.Query("part_id") // 新增：配件ID筛选

	type EmployeeStat struct {
		EmployeeID   uint   `json:"employee_id"`
		EmployeeNo   string `json:"employee_no"`
		EmployeeName string `json:"employee_name"`
		Department   string `json:"department"`
		BorrowCount  int64  `json:"borrow_count"`
		TotalBorrow  int64  `json:"total_borrow"`
		TotalDamaged int64  `json:"total_damaged"`
	}

	var stats []EmployeeStat

	query := h.db.Table("borrow_record").
		Select(`
			employee.id as employee_id,
			employee.employee_no,
			employee.name as employee_name,
			employee.department,
			COUNT(*) as borrow_count,
			SUM(borrow_record.borrow_quantity) as total_borrow,
			SUM(borrow_record.damaged_quantity) as total_damaged
		`).
		Joins("LEFT JOIN employee ON borrow_record.employee_id = employee.id").
		Where("borrow_record.borrow_time BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59")

	// 如果指定了配件ID，添加筛选条件
	if partID != "" {
		query = query.Where("borrow_record.part_id = ?", partID)
	}

	query.Group("employee.id").
		Order("total_borrow DESC").
		Scan(&stats)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data":    stats,
	})
}

// GetDepartmentReport 按部门统计
func (h *ReportHandler) GetDepartmentReport(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	type DepartmentStat struct {
		Department    string `json:"department"`
		EmployeeCount int64  `json:"employee_count"`
		BorrowCount   int64  `json:"borrow_count"`
		TotalBorrow   int64  `json:"total_borrow"`
		TotalDamaged  int64  `json:"total_damaged"`
	}

	var stats []DepartmentStat

	h.db.Table("borrow_record").
		Select(`
			employee.department,
			COUNT(DISTINCT employee.id) as employee_count,
			COUNT(*) as borrow_count,
			SUM(borrow_record.borrow_quantity) as total_borrow,
			SUM(borrow_record.damaged_quantity) as total_damaged
		`).
		Joins("LEFT JOIN employee ON borrow_record.employee_id = employee.id").
		Where("borrow_record.borrow_time BETWEEN ? AND ? AND employee.department IS NOT NULL",
			startDate+" 00:00:00", endDate+" 23:59:59").
		Group("employee.department").
		Order("total_borrow DESC").
		Scan(&stats)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data":    stats,
	})
}

// GetDetailedReport 获取详细报表（每人每月每产品）
func (h *ReportHandler) GetDetailedReport(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -3, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	employeeID := c.Query("employee_id")
	partID := c.Query("part_id")

	type DetailRecord struct {
		Month        string `json:"month"`
		EmployeeID   uint   `json:"employee_id"`
		EmployeeNo   string `json:"employee_no"`
		EmployeeName string `json:"employee_name"`
		Department   string `json:"department"`
		PartID       uint   `json:"part_id"`
		PartNo       string `json:"part_no"`
		PartName     string `json:"part_name"`
		BorrowCount  int64  `json:"borrow_count"`
		TotalBorrow  int64  `json:"total_borrow"`
		TotalReturn  int64  `json:"total_return"`
		TotalDamaged int64  `json:"total_damaged"`
		Unreturned   int64  `json:"unreturned"`
	}

	var records []DetailRecord

	query := h.db.Table("borrow_record").
		Select(`
			DATE_FORMAT(borrow_record.borrow_time, '%Y-%m') as month,
			employee.id as employee_id,
			employee.employee_no,
			employee.name as employee_name,
			employee.department,
			part.id as part_id,
			part.part_no,
			part.name as part_name,
			COUNT(*) as borrow_count,
			SUM(borrow_record.borrow_quantity) as total_borrow,
			SUM(borrow_record.return_quantity) as total_return,
			SUM(borrow_record.damaged_quantity) as total_damaged,
			SUM(borrow_record.borrow_quantity - borrow_record.return_quantity - borrow_record.damaged_quantity) as unreturned
		`).
		Joins("LEFT JOIN employee ON borrow_record.employee_id = employee.id").
		Joins("LEFT JOIN part ON borrow_record.part_id = part.id").
		Where("borrow_record.borrow_time BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59")

	// 可选筛选条件
	if employeeID != "" {
		query = query.Where("employee.id = ?", employeeID)
	}
	if partID != "" {
		query = query.Where("part.id = ?", partID)
	}

	query.Group("month, employee.id, part.id").
		Order("month DESC, employee.name, part.name").
		Scan(&records)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data":    records,
	})
}

// PushReport 推送报表到通知渠道 (支持多维度表格)
func (h *ReportHandler) PushReport(c *gin.Context) {
	if h.ni == nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "通知集成未初始化"})
		return
	}

	dimension := c.DefaultQuery("dimension", "detail")
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	var sb strings.Builder
	var title string

	switch dimension {
	case "part":
		title = "📊 配件领用排行推送"
		sb.WriteString("## 📦 配件领用排行 (Top 15)\n\n")
		sb.WriteString(fmt.Sprintf("周期: `%s` 至 `%s`\n\n", startDate, endDate))
		sb.WriteString("| 配件名称 | 编号 | 领用总数 | 损毁 |\n")
		sb.WriteString("| :--- | :--- | :--- | :--- |\n")

		var stats []struct {
			Name         string
			PartNo       string
			TotalBorrow  int
			TotalDamaged int
		}
		h.db.Table("borrow_record").
			Select("part.name, part.part_no, SUM(borrow_record.borrow_quantity) as total_borrow, SUM(borrow_record.damaged_quantity) as total_damaged").
			Joins("LEFT JOIN part ON borrow_record.part_id = part.id").
			Where("borrow_record.borrow_time BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59").
			Group("part.id").Order("total_borrow DESC").Limit(15).Scan(&stats)

		for _, s := range stats {
			sb.WriteString(fmt.Sprintf("| %s | %s | %d | %d |\n", s.Name, s.PartNo, s.TotalBorrow, s.TotalDamaged))
		}

	case "employee":
		title = "📊 员工领用排行推送"
		sb.WriteString("## 👤 员工领用排行 (Top 15)\n\n")
		sb.WriteString(fmt.Sprintf("周期: `%s` 至 `%s`\n\n", startDate, endDate))
		sb.WriteString("| 员工姓名 | 部门 | 领用总数 | 损毁 |\n")
		sb.WriteString("| :--- | :--- | : :--- | :--- |\n")

		var stats []struct {
			Name         string
			Department   string
			TotalBorrow  int
			TotalDamaged int
		}
		h.db.Table("borrow_record").
			Select("employee.name, employee.department, SUM(borrow_record.borrow_quantity) as total_borrow, SUM(borrow_record.damaged_quantity) as total_damaged").
			Joins("LEFT JOIN employee ON borrow_record.employee_id = employee.id").
			Where("borrow_record.borrow_time BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59").
			Group("employee.id").Order("total_borrow DESC").Limit(15).Scan(&stats)

		for _, s := range stats {
			sb.WriteString(fmt.Sprintf("| %s | %s | %d | %d |\n", s.Name, s.Department, s.TotalBorrow, s.TotalDamaged))
		}

	case "department":
		title = "📊 部门领用排行推送"
		sb.WriteString("## 🏢 部门领用排行\n\n")
		sb.WriteString(fmt.Sprintf("周期: `%s` 至 `%s`\n\n", startDate, endDate))
		sb.WriteString("| 部门 | 员工数 | 领用总数 | 损毁 |\n")
		sb.WriteString("| :--- | :--- | :--- | :--- |\n")

		var stats []struct {
			Department    string
			EmployeeCount int
			TotalBorrow   int
			TotalDamaged  int
		}
		h.db.Table("borrow_record").
			Select("employee.department, COUNT(DISTINCT employee.id) as employee_count, SUM(borrow_record.borrow_quantity) as total_borrow, SUM(borrow_record.damaged_quantity) as total_damaged").
			Joins("LEFT JOIN employee ON borrow_record.employee_id = employee.id").
			Where("borrow_record.borrow_time BETWEEN ? AND ? AND employee.department IS NOT NULL", startDate+" 00:00:00", endDate+" 23:59:59").
			Group("employee.department").Order("total_borrow DESC").Scan(&stats)

		for _, s := range stats {
			sb.WriteString(fmt.Sprintf("| %s | %d | %d | %d |\n", s.Department, s.EmployeeCount, s.TotalBorrow, s.TotalDamaged))
		}

	default: // detail
		title = "📊 业务明细推送"
		sb.WriteString("## 📋 业务明细统计报表\n\n")
		sb.WriteString(fmt.Sprintf("周期: `%s` 至 `%s`\n\n", startDate, endDate))
		sb.WriteString("| 员工 | 领用配件 | 数量 |\n")
		sb.WriteString("| :--- | :--- | :--- |\n")

		var ranks []struct {
			EmployeeName string
			PartName     string
			TotalBorrow  int
		}
		h.db.Table("borrow_record").
			Select("employee.name as employee_name, part.name as part_name, SUM(borrow_record.borrow_quantity) as total_borrow").
			Joins("LEFT JOIN employee ON borrow_record.employee_id = employee.id").
			Joins("LEFT JOIN part ON borrow_record.part_id = part.id").
			Where("borrow_record.borrow_time BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59").
			Group("employee.id, part.id").Order("total_borrow DESC").Limit(15).Scan(&ranks)

		for _, r := range ranks {
			sb.WriteString(fmt.Sprintf("| %s | %s | %d |\n", r.EmployeeName, r.PartName, r.TotalBorrow))
		}
	}

	sb.WriteString(fmt.Sprintf("\n> 推送时间: %s", time.Now().Format("2006-01-02 15:04:05")))
	h.ni.NotifyCustomReport(title, sb.String())

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "报表已生成并推送至群聊",
	})
}

// DownloadReport 下载报表文件
func (h *ReportHandler) DownloadReport(c *gin.Context) {
	filename := c.Param("filename")
	reportDir := "temp/reports"

	// 文件名安全验证（防止路径遍历攻击）
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "非法文件名"})
		return
	}

	// 验证文件名格式（白名单）
	validPrefixes := []string{"daily_report_", "monthly_report_", "weekly_report_"}
	isValid := false
	for _, prefix := range validPrefixes {
		if strings.HasPrefix(filename, prefix) && strings.HasSuffix(filename, ".xlsx") {
			isValid = true
			break
		}
	}

	if !isValid {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "不支持的文件类型"})
		return
	}

	// 构建文件路径
	filepath := fmt.Sprintf("%s/%s", reportDir, filename)

	// 检查文件是否存在
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "文件不存在或已过期"})
		return
	}

	// 设置响应头并发送文件
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.File(filepath)
}
