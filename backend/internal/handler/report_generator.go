package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"warehouse/internal/model"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// ReportGenerator 报表生成器
type ReportGenerator struct {
	db        *gorm.DB
	reportDir string
	baseURL   string
}

// NewReportGenerator 创建报表生成器
func NewReportGenerator(db *gorm.DB, baseURL string) *ReportGenerator {
	reportDir := "temp/reports"
	// 确保目录存在
	os.MkdirAll(reportDir, 0755)

	return &ReportGenerator{
		db:        db,
		reportDir: reportDir,
		baseURL:   baseURL,
	}
}

// GenerateDailyReportExcel 生成日报Excel
func (rg *ReportGenerator) GenerateDailyReportExcel(date time.Time) (string, error) {
	f := excelize.NewFile()
	defer f.Close()

	// 创建概览Sheet
	sheet := "概览"
	f.SetSheetName("Sheet1", sheet)

	// 设置标题
	f.SetCellValue(sheet, "A1", "凯越机械仓储管理系统")
	f.SetCellValue(sheet, "A2", "每日业务报表")
	f.SetCellValue(sheet, "A3", fmt.Sprintf("日期: %s", date.Format("2006-01-02")))

	// 样式设置
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 16},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	f.SetCellStyle(sheet, "A1", "E1", titleStyle)

	// 统计数据
	today := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	var borrowCount int64
	rg.db.Model(&model.BorrowRecord{}).Where("borrow_time >= ?", today).Count(&borrowCount)

	var returnCount int64
	rg.db.Model(&model.BorrowRecord{}).Where("return_time >= ?", today).Count(&returnCount)

	var lowStockParts []model.Part
	rg.db.Where("stock_quantity < warning_threshold").Find(&lowStockParts)

	// 数据概览
	f.SetCellValue(sheet, "A5", "数据指标")
	f.SetCellValue(sheet, "B5", "数值")
	f.SetCellValue(sheet, "A6", "今日领用")
	f.SetCellValue(sheet, "B6", borrowCount)
	f.SetCellValue(sheet, "A7", "今日归还")
	f.SetCellValue(sheet, "B7", returnCount)
	f.SetCellValue(sheet, "A8", "库存预警")
	f.SetCellValue(sheet, "B8", len(lowStockParts))

	// 低库存配件详情
	if len(lowStockParts) > 0 {
		detailSheet := "低库存配件"
		f.NewSheet(detailSheet)

		f.SetCellValue(detailSheet, "A1", "配件编号")
		f.SetCellValue(detailSheet, "B1", "配件名称")
		f.SetCellValue(detailSheet, "C1", "当前库存")
		f.SetCellValue(detailSheet, "D1", "预警阈值")
		f.SetCellValue(detailSheet, "E1", "缺口")

		headerStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
		})
		f.SetCellStyle(detailSheet, "A1", "E1", headerStyle)

		for i, part := range lowStockParts {
			row := i + 2
			f.SetCellValue(detailSheet, fmt.Sprintf("A%d", row), part.PartNo)
			f.SetCellValue(detailSheet, fmt.Sprintf("B%d", row), part.Name)
			f.SetCellValue(detailSheet, fmt.Sprintf("C%d", row), part.StockQuantity)
			f.SetCellValue(detailSheet, fmt.Sprintf("D%d", row), part.WarningThreshold)
			f.SetCellValue(detailSheet, fmt.Sprintf("E%d", row), part.WarningThreshold-part.StockQuantity)
		}

		// 设置列宽
		f.SetColWidth(detailSheet, "A", "A", 15)
		f.SetColWidth(detailSheet, "B", "B", 25)
		f.SetColWidth(detailSheet, "C", "E", 12)
	}

	// 保存文件
	filename := fmt.Sprintf("daily_report_%s.xlsx", date.Format("20060102"))
	filepath := filepath.Join(rg.reportDir, filename)

	if err := f.SaveAs(filepath); err != nil {
		return "", err
	}

	return filename, nil
}

// GenerateMonthlyReportExcel 生成月报Excel
func (rg *ReportGenerator) GenerateMonthlyReportExcel(date time.Time) (string, error) {
	f := excelize.NewFile()
	defer f.Close()

	// 创建概览Sheet
	sheet := "概览"
	f.SetSheetName("Sheet1", sheet)

	// 设置标题
	f.SetCellValue(sheet, "A1", "凯越机械仓储管理系统")
	f.SetCellValue(sheet, "A2", "月度业务报告")
	f.SetCellValue(sheet, "A3", fmt.Sprintf("月份: %s", date.Format("2006年01月")))

	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 16},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	f.SetCellStyle(sheet, "A1", "E1", titleStyle)

	// 统计本月数据
	monthStart := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())

	var borrowCount int64
	rg.db.Model(&model.BorrowRecord{}).Where("borrow_time >= ?", monthStart).Count(&borrowCount)

	var returnCount int64
	rg.db.Model(&model.BorrowRecord{}).Where("return_time >= ?", monthStart).Count(&returnCount)

	var inboundCount int64
	rg.db.Model(&model.InboundRecord{}).Where("inbound_time >= ?", monthStart).Count(&inboundCount)

	// 数据概览
	f.SetCellValue(sheet, "A5", "数据指标")
	f.SetCellValue(sheet, "B5", "数值")
	f.SetCellValue(sheet, "A6", "本月领用")
	f.SetCellValue(sheet, "B6", borrowCount)
	f.SetCellValue(sheet, "A7", "本月归还")
	f.SetCellValue(sheet, "B7", returnCount)
	f.SetCellValue(sheet, "A8", "本月入库")
	f.SetCellValue(sheet, "B8", inboundCount)

	// 活跃员工TOP10
	type EmployeeStat struct {
		EmployeeID   uint
		EmployeeName string
		BorrowCount  int64
	}
	var topEmployees []EmployeeStat
	rg.db.Table("borrow_records").
		Select("employee_id, COUNT(*) as borrow_count").
		Where("borrow_time >= ?", monthStart).
		Group("employee_id").
		Order("borrow_count DESC").
		Limit(10).
		Scan(&topEmployees)

	// 填充员工姓名
	for i := range topEmployees {
		var emp model.Employee
		rg.db.First(&emp, topEmployees[i].EmployeeID)
		topEmployees[i].EmployeeName = emp.Name
	}

	if len(topEmployees) > 0 {
		empSheet := "活跃员工TOP10"
		f.NewSheet(empSheet)

		f.SetCellValue(empSheet, "A1", "排名")
		f.SetCellValue(empSheet, "B1", "员工姓名")
		f.SetCellValue(empSheet, "C1", "领用次数")

		headerStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
		})
		f.SetCellStyle(empSheet, "A1", "C1", headerStyle)

		for i, emp := range topEmployees {
			row := i + 2
			f.SetCellValue(empSheet, fmt.Sprintf("A%d", row), i+1)
			f.SetCellValue(empSheet, fmt.Sprintf("B%d", row), emp.EmployeeName)
			f.SetCellValue(empSheet, fmt.Sprintf("C%d", row), emp.BorrowCount)
		}

		f.SetColWidth(empSheet, "A", "A", 8)
		f.SetColWidth(empSheet, "B", "B", 20)
		f.SetColWidth(empSheet, "C", "C", 12)
	}

	// 热门配件TOP10
	type PartStat struct {
		PartID      uint
		PartName    string
		BorrowCount int64
	}
	var topParts []PartStat
	rg.db.Table("borrow_records").
		Select("part_id, COUNT(*) as borrow_count").
		Where("borrow_time >= ?", monthStart).
		Group("part_id").
		Order("borrow_count DESC").
		Limit(10).
		Scan(&topParts)

	// 填充配件名称
	for i := range topParts {
		var part model.Part
		rg.db.First(&part, topParts[i].PartID)
		topParts[i].PartName = part.Name
	}

	if len(topParts) > 0 {
		partSheet := "热门配件TOP10"
		f.NewSheet(partSheet)

		f.SetCellValue(partSheet, "A1", "排名")
		f.SetCellValue(partSheet, "B1", "配件名称")
		f.SetCellValue(partSheet, "C1", "被领用次数")

		headerStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
		})
		f.SetCellStyle(partSheet, "A1", "C1", headerStyle)

		for i, part := range topParts {
			row := i + 2
			f.SetCellValue(partSheet, fmt.Sprintf("A%d", row), i+1)
			f.SetCellValue(partSheet, fmt.Sprintf("B%d", row), part.PartName)
			f.SetCellValue(partSheet, fmt.Sprintf("C%d", row), part.BorrowCount)
		}

		f.SetColWidth(partSheet, "A", "A", 8)
		f.SetColWidth(partSheet, "B", "B", 25)
		f.SetColWidth(partSheet, "C", "C", 12)
	}

	// 保存文件
	filename := fmt.Sprintf("monthly_report_%s.xlsx", date.Format("200601"))
	filepath := filepath.Join(rg.reportDir, filename)

	if err := f.SaveAs(filepath); err != nil {
		return "", err
	}

	return filename, nil
}

// GetDownloadURL 获取下载URL
func (rg *ReportGenerator) GetDownloadURL(filename string) string {
	return fmt.Sprintf("%s/api/reports/download/%s", rg.baseURL, filename)
}

// CleanOldReports 清理7天前的报表文件
func (rg *ReportGenerator) CleanOldReports() error {
	files, err := os.ReadDir(rg.reportDir)
	if err != nil {
		return err
	}

	cutoff := time.Now().AddDate(0, 0, -7)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			filepath := filepath.Join(rg.reportDir, file.Name())
			os.Remove(filepath)
		}
	}

	return nil
}
