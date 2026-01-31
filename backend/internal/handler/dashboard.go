package handler

import (
	"net/http"
	"warehouse/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	db *gorm.DB
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

// GetDashboardStats 获取仪表盘统计数据
func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	var stats struct {
		PartCount     int64 `json:"part_count"`     // 配件种类
		TotalStock    int64 `json:"total_stock"`    // 总库存数量
		WarningCount  int64 `json:"warning_count"`  // 低库存预警数
		MonthlyBorrow int64 `json:"monthly_borrow"` // 本月领用次数
		BorrowedCount int64 `json:"borrowed_count"` // 未归还数量
		EmployeeCount int64 `json:"employee_count"` // 在职员工数
	}

	// 基础统计
	h.db.Model(&model.Part{}).Count(&stats.PartCount)
	h.db.Model(&model.Part{}).Select("COALESCE(SUM(stock_quantity), 0)").Scan(&stats.TotalStock)
	h.db.Model(&model.Part{}).Where("stock_quantity < warning_threshold").Count(&stats.WarningCount)
	h.db.Model(&model.BorrowRecord{}).Where("DATE_FORMAT(borrow_time, '%Y-%m') = DATE_FORMAT(NOW(), '%Y-%m')").Count(&stats.MonthlyBorrow)
	h.db.Model(&model.BorrowRecord{}).Where("status IN ('BORROWED', 'PARTIAL_RETURNED')").Count(&stats.BorrowedCount)
	h.db.Model(&model.Employee{}).Where("status = 1").Count(&stats.EmployeeCount)

	// 新增：报废责任排行 (Top 5)
	type Scrapper struct {
		Name       string `json:"name"`
		Department string `json:"department"`
		Count      int64  `json:"count"`
	}
	var topScrappers []Scrapper
	h.db.Table("scrap_inventory").
		Select("employee.name, employee.department, SUM(scrap_inventory.quantity) as count").
		Joins("LEFT JOIN employee ON scrap_inventory.employee_id = employee.id").
		Group("employee.id").
		Order("count DESC").
		Limit(5).
		Scan(&topScrappers)

	// 新增：月度趋势 (最近12个月 领用 vs 归还)
	type TrendStat struct {
		Month  string `json:"month"`
		Borrow int64  `json:"borrow"`
		Return int64  `json:"return"`
	}
	var trends []TrendStat

	// 使用Raw SQL处理复杂的日期聚合，保证12个月覆盖需要更复杂的逻辑，这里简化为查询存在的记录
	h.db.Raw(`
		SELECT 
			DATE_FORMAT(borrow_time, '%Y-%m') as month,
			SUM(borrow_quantity) as borrow,
			SUM(return_quantity) as 'return'
		FROM borrow_record
		WHERE borrow_time >= DATE_SUB(NOW(), INTERVAL 11 MONTH)
		GROUP BY month
		ORDER BY month ASC
	`).Scan(&trends)

	// 新增：分类库存分布
	type CategoryStat struct {
		Name  string `json:"name"`
		Value int64  `json:"value"`
	}
	var catStats []CategoryStat
	h.db.Table("part").
		Select("part_category.name, SUM(part.stock_quantity) as value").
		Joins("LEFT JOIN part_category ON part.category_id = part_category.id").
		Where("part_category.id IS NOT NULL").
		Group("part_category.name").
		Order("value DESC").
		Limit(8).
		Scan(&catStats)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data": gin.H{
			"base":           stats,
			"top_scrappers":  topScrappers,
			"monthly_trend":  trends,
			"category_stats": catStats,
		},
	})
}
