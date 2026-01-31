package handler

import (
	"net/http"
	"strconv"
	"time"
	"warehouse/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InboundHandler struct {
	db *gorm.DB
	ni *NotificationIntegrator
}

func NewInboundHandler(db *gorm.DB, ni *NotificationIntegrator) *InboundHandler {
	return &InboundHandler{db: db, ni: ni}
}

// RestockPart 补货/入库
func (h *InboundHandler) RestockPart(c *gin.Context) {
	var req struct {
		PartID   uint   `json:"part_id" binding:"required"`
		Quantity int    `json:"quantity" binding:"required,gt=0"`
		Remark   string `json:"remark"`
		BatchNo  string `json:"batch_no"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 获取当前操作人 (从中间件)
	operatorName := "Unknown"
	if v, exists := c.Get("username"); exists {
		operatorName = v.(string)
	}

	// 开启事务
	err := h.db.Transaction(func(tx *gorm.DB) error {
		// 1. 检查配件是否存在
		var part model.Part
		if err := tx.First(&part, req.PartID).Error; err != nil {
			return err
		}

		// 2. 更新库存
		part.StockQuantity += req.Quantity
		if err := tx.Model(&part).Update("stock_quantity", part.StockQuantity).Error; err != nil {
			return err
		}

		// 3. 创建入库记录
		record := model.InboundRecord{
			PartID:    req.PartID,
			Quantity:  req.Quantity,
			Operator:  operatorName,
			BatchNo:   req.BatchNo,
			Remark:    req.Remark,
			CreatedAt: time.Now(),
		}
		if err := tx.Create(&record).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "补货失败: " + err.Error()})
		return
	}

	// 异步发送补货通知
	if h.ni != nil {
		var part model.Part
		h.db.First(&part, req.PartID)
		go h.ni.NotifyRestock(part.Name, req.Quantity, part.StockQuantity)
		// 同时检查是否解除了库存预警
		go h.ni.CheckAndNotifyStockWarning(req.PartID)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "补货成功"})
}

// GetInboundList 获取入库记录列表
func (h *InboundHandler) GetInboundList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	partName := c.Query("part_name")

	var records []model.InboundRecord
	var total int64

	query := h.db.Model(&model.InboundRecord{}).Preload("Part")

	if partName != "" {
		// 关联查询
		query = query.Joins("JOIN part ON part.id = inbound_record.part_id").
			Where("part.name LIKE ? OR part.part_no LIKE ?", "%"+partName+"%", "%"+partName+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("inbound_record.created_at DESC").Find(&records)

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
