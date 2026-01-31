package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"warehouse/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BorrowHandler struct {
	db *gorm.DB
	ni *NotificationIntegrator
}

func NewBorrowHandler(db *gorm.DB, ni *NotificationIntegrator) *BorrowHandler {
	return &BorrowHandler{db: db, ni: ni}
}

// CreateBorrowRecord 创建领用记录
func (h *BorrowHandler) CreateBorrowRecord(c *gin.Context) {
	var req struct {
		EmployeeID     uint   `json:"employee_id" binding:"required"`
		PartID         uint   `json:"part_id" binding:"required"`
		BorrowQuantity int    `json:"borrow_quantity" binding:"required,gt=0"`
		Remark         string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 获取当前用户ID（管理员）
	username, _ := c.Get("username")
	var admin model.User
	h.db.Where("username = ?", username).First(&admin)

	// 检查配件库存
	var part model.Part
	if err := h.db.First(&part, req.PartID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "配件不存在",
		})
		return
	}

	if part.StockQuantity < req.BorrowQuantity {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": fmt.Sprintf("库存不足，当前库存：%d", part.StockQuantity),
		})
		return
	}

	// 开启事务
	tx := h.db.Begin()

	// 创建领用记录
	record := model.BorrowRecord{
		RecordNo:       fmt.Sprintf("BR%s%06d", time.Now().Format("20060102"), time.Now().Unix()%1000000),
		EmployeeID:     req.EmployeeID,
		PartID:         req.PartID,
		BorrowQuantity: req.BorrowQuantity,
		Status:         "BORROWED",
		BorrowTime:     time.Now(),
		BorrowAdminID:  admin.ID,
		Remark:         req.Remark,
	}

	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建记录失败",
		})
		return
	}

	// 扣减库存
	if err := tx.Model(&model.Part{}).Where("id = ?", req.PartID).
		Update("stock_quantity", gorm.Expr("stock_quantity - ?", req.BorrowQuantity)).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "扣减库存失败",
		})
		return
	}

	tx.Commit()

	// 异步发送通知
	if h.ni != nil {
		// 1. 发送领用通知
		var employee model.Employee
		h.db.First(&employee, req.EmployeeID)
		go h.ni.NotifyBorrowCreated(employee.Name, part.Name, req.BorrowQuantity)

		// 2. 检查并发送库存预警通知(因为领用后库存可能低于阈值)
		go h.ni.CheckAndNotifyStockWarning(req.PartID)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "领用登记成功",
		"data":    record,
	})
}

// ReturnBorrowRecord 旧件处置登记（回收或报废）
func (h *BorrowHandler) ReturnBorrowRecord(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		ReturnQuantity  int `json:"return_quantity" binding:"min=0"`
		DamagedQuantity int `json:"damaged_quantity" binding:"min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 获取当前管理员
	username, _ := c.Get("username")
	var admin model.User
	h.db.Where("username = ?", username).First(&admin)

	// 获取领用记录
	var record model.BorrowRecord
	if err := h.db.First(&record, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "记录不存在",
		})
		return
	}

	// 计算已归还总数
	totalReturned := record.ReturnQuantity + record.DamagedQuantity + req.ReturnQuantity + req.DamagedQuantity
	if totalReturned > record.BorrowQuantity {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "处置数量超过持有数量",
		})
		return
	}

	// 开启事务
	tx := h.db.Begin()

	// 更新记录
	now := time.Now()
	updates := map[string]interface{}{
		"return_quantity":  record.ReturnQuantity + req.ReturnQuantity,
		"damaged_quantity": record.DamagedQuantity + req.DamagedQuantity,
		"return_time":      &now,
		"return_admin_id":  admin.ID,
	}

	// 判断状态
	if totalReturned == record.BorrowQuantity {
		updates["status"] = "RETURNED"
	} else {
		updates["status"] = "PARTIAL_RETURNED"
	}

	if err := tx.Model(&record).Updates(updates).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "更新记录失败",
		})
		return
	}

	// 1. 如果有回收，计入旧件库
	if req.ReturnQuantity > 0 {
		var oldInv model.OldPartInventory
		// 累加原有记录或创建新记录
		if err := tx.Where("part_id = ? AND employee_id = ?", record.PartID, record.EmployeeID).First(&oldInv).Error; err == nil {
			if err := tx.Model(&oldInv).Update("quantity", gorm.Expr("quantity + ?", req.ReturnQuantity)).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, gin.H{"code": 500, "message": "更新旧件库存失败"})
				return
			}
		} else {
			if err := tx.Create(&model.OldPartInventory{
				PartID:     record.PartID,
				EmployeeID: record.EmployeeID,
				Quantity:   req.ReturnQuantity,
			}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, gin.H{"code": 500, "message": "创建旧件库存失败"})
				return
			}
		}
	}

	// 2. 如果有报废，计入废品库
	if req.DamagedQuantity > 0 {
		var scrapInv model.ScrapInventory
		if err := tx.Where("part_id = ? AND employee_id = ?", record.PartID, record.EmployeeID).First(&scrapInv).Error; err == nil {
			if err := tx.Model(&scrapInv).Update("quantity", gorm.Expr("quantity + ?", req.DamagedQuantity)).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, gin.H{"code": 500, "message": "更新废品库存失败"})
				return
			}
		} else {
			if err := tx.Create(&model.ScrapInventory{
				PartID:     record.PartID,
				EmployeeID: record.EmployeeID,
				Quantity:   req.DamagedQuantity,
			}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, gin.H{"code": 500, "message": "创建废品库存失败"})
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "提交失败"})
		return
	}

	// 异步发送通知
	if h.ni != nil {
		var employee model.Employee
		var part model.Part
		h.db.First(&employee, record.EmployeeID)
		h.db.First(&part, record.PartID)
		go h.ni.NotifyReturnCreated(employee.Name, part.Name, req.ReturnQuantity, req.DamagedQuantity)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "处置成功",
	})
}

// CheckUnreturned 检查员工+配件的待回收旧件记录
func (h *BorrowHandler) CheckUnreturned(c *gin.Context) {
	employeeID := c.Query("employee_id")
	partID := c.Query("part_id")

	if employeeID == "" || partID == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	var record model.BorrowRecord
	result := h.db.Where("employee_id = ? AND part_id = ? AND status IN ('BORROWED', 'PARTIAL_RETURNED')",
		employeeID, partID).
		Order("borrow_time DESC").
		First(&record)

	if result.Error != nil {
		// 没有未归还记录
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "无待回收记录",
			"data": gin.H{
				"has_unreturned": false,
			},
		})
		return
	}

	// 计算未归还数量
	unreturned := record.BorrowQuantity - record.ReturnQuantity - record.DamagedQuantity

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "有待回收旧件",
		"data": gin.H{
			"has_unreturned":    true,
			"record_id":         record.ID,
			"borrow_quantity":   record.BorrowQuantity,
			"returned_quantity": record.ReturnQuantity,
			"damaged_quantity":  record.DamagedQuantity,
			"unreturned":        unreturned,
			"borrow_time":       record.BorrowTime,
		},
	})
}

// GetBorrowRecordList 获取领用记录列表
func (h *BorrowHandler) GetBorrowRecordList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status := c.Query("status")

	var records []struct {
		model.BorrowRecord
		EmployeeName string `json:"employee_name"`
		PartName     string `json:"part_name"`
		PartNo       string `json:"part_no"`
	}
	var total int64

	query := h.db.Table("borrow_record").
		Select("borrow_record.*, employee.name as employee_name, part.name as part_name, part.part_no").
		Joins("LEFT JOIN employee ON borrow_record.employee_id = employee.id").
		Joins("LEFT JOIN part ON borrow_record.part_id = part.id")

	if status != "" {
		query = query.Where("borrow_record.status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("borrow_record.created_at DESC").Scan(&records)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data": gin.H{
			"total":   total,
			"records": records,
			"current": page,
			"size":    pageSize,
		},
	})
}

// GetOldPartInventory 获取旧件库列表 (带溯源)
func (h *BorrowHandler) GetOldPartInventory(c *gin.Context) {
	type DisplayItem struct {
		model.OldPartInventory
		PartName     string `json:"part_name"`
		PartNo       string `json:"part_no"`
		EmployeeName string `json:"employee_name"`
		Department   string `json:"department"`
	}

	var list []DisplayItem
	err := h.db.Table("old_part_inventory").
		Select("old_part_inventory.*, part.name as part_name, part.part_no as part_no, employee.name as employee_name, employee.department").
		Joins("LEFT JOIN part ON old_part_inventory.part_id = part.id").
		Joins("LEFT JOIN employee ON old_part_inventory.employee_id = employee.id").
		Order("old_part_inventory.updated_at desc").
		Scan(&list).Error

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "获取旧件库失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": list})
}

// GetScrapInventory 获取废品库列表 (带溯源)
func (h *BorrowHandler) GetScrapInventory(c *gin.Context) {
	type DisplayItem struct {
		model.ScrapInventory
		PartName     string `json:"part_name"`
		PartNo       string `json:"part_no"`
		EmployeeName string `json:"employee_name"`
		Department   string `json:"department"`
	}

	var list []DisplayItem
	err := h.db.Table("scrap_inventory").
		Select("scrap_inventory.*, part.name as part_name, part.part_no as part_no, employee.name as employee_name, employee.department").
		Joins("LEFT JOIN part ON scrap_inventory.part_id = part.id").
		Joins("LEFT JOIN employee ON scrap_inventory.employee_id = employee.id").
		Order("scrap_inventory.updated_at desc").
		Scan(&list).Error

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "获取废品库失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": list})
}
