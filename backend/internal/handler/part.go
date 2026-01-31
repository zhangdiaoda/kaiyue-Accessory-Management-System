package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"warehouse/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type PartHandler struct {
	db *gorm.DB
	ni *NotificationIntegrator
}

func NewPartHandler(db *gorm.DB, ni *NotificationIntegrator) *PartHandler {
	return &PartHandler{db: db, ni: ni}
}

// GetPartList 获取配件列表
func (h *PartHandler) GetPartList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	var parts []model.Part
	var total int64

	query := h.db.Model(&model.Part{})

	// 搜索
	if keyword != "" {
		query = query.Where("name LIKE ? OR part_no LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 统计总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&parts)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data": gin.H{
			"total":   total,
			"records": parts,
			"current": page,
			"size":    pageSize,
		},
	})
}

// CreatePart 创建配件
func (h *PartHandler) CreatePart(c *gin.Context) {
	var part model.Part
	if err := c.ShouldBindJSON(&part); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	if err := h.db.Create(&part).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建失败: " + err.Error(),
		})
		return
	}

	// 异步发送库存预警通知
	if h.ni != nil {
		go h.ni.CheckAndNotifyStockWarning(part.ID)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    part,
	})
}

// UpdatePart 更新配件
func (h *PartHandler) UpdatePart(c *gin.Context) {
	id := c.Param("id")
	var part model.Part

	if err := c.ShouldBindJSON(&part); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	if err := h.db.Model(&model.Part{}).Where("id = ?", id).Updates(&part).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "更新失败",
		})
		return
	}

	// 异步发送库存预警通知
	if h.ni != nil {
		idUint, _ := strconv.ParseUint(id, 10, 32)
		go h.ni.CheckAndNotifyStockWarning(uint(idUint))
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
	})
}

// DeletePart 删除配件
func (h *PartHandler) DeletePart(c *gin.Context) {
	id := c.Param("id")

	if err := h.db.Delete(&model.Part{}, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "删除失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// GetLowStockParts 获取低库存配件
func (h *PartHandler) GetLowStockParts(c *gin.Context) {
	var parts []model.Part

	h.db.Where("stock_quantity < warning_threshold").
		Order("stock_quantity ASC").
		Find(&parts)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data":    parts,
	})
}

// DownloadTemplate 下载导入模板
func (h *PartHandler) DownloadTemplate(c *gin.Context) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// 设置表头
	headers := []string{"配件编号*", "名称*", "分类", "规格", "单位", "库存数量", "预警阈值", "单价", "备注"}
	for i, v := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue("Sheet1", cell, v)
	}

	// 设置示例数据
	examples := []interface{}{"P-001", "示例配件", "通用", "规格A", "个", 100, 10, 99.5, "备注信息"}
	for i, v := range examples {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue("Sheet1", cell, v)
	}

	// 强制下载
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=part_import_template.xlsx")

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// ImportParts 导入配件
func (h *PartHandler) ImportParts(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "上传失败"})
		return
	}
	defer file.Close()

	f, err := excelize.OpenReader(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Excel解析失败"})
		return
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "读取Sheet1失败"})
		return
	}

	success := 0
	failed := 0

	for i, row := range rows {
		if i == 0 {
			continue
		} // 跳过表头
		if len(row) < 2 {
			continue
		} // 必填校验

		partNo := row[0]
		name := row[1]
		if partNo == "" || name == "" {
			failed++
			continue
		}

		// 处理可选字段
		categoryName := ""
		if len(row) > 2 {
			categoryName = row[2]
		}
		spec := ""
		if len(row) > 3 {
			spec = row[3]
		}
		unit := "件"
		if len(row) > 4 && row[4] != "" {
			unit = row[4]
		}

		stock := 0
		if len(row) > 5 {
			stock, _ = strconv.Atoi(row[5])
		}

		warning := 10
		if len(row) > 6 {
			warning, _ = strconv.Atoi(row[6])
		}

		price := 0.0
		if len(row) > 7 {
			price, _ = strconv.ParseFloat(row[7], 64)
		}

		remark := ""
		if len(row) > 8 {
			remark = row[8]
		}

		// 处理分类
		var catID uint = 0
		if categoryName != "" {
			var cat model.PartCategory
			if err := h.db.FirstOrCreate(&cat, model.PartCategory{Name: categoryName}).Error; err == nil {
				catID = cat.ID
			}
		}

		// 存入数据库 (Upsert: 根据PartNo更新或插入)
		var part model.Part
		result := h.db.Where("part_no = ?", partNo).First(&part)

		if result.Error == gorm.ErrRecordNotFound {
			// Create
			h.db.Create(&model.Part{
				PartNo: partNo, Name: name, CategoryID: catID, Specification: spec,
				Unit: unit, StockQuantity: stock, WarningThreshold: warning, Price: price, Remark: remark,
			})
			success++
		} else {
			// Update
			h.db.Model(&part).Updates(model.Part{
				Name: name, CategoryID: catID, Specification: spec,
				Unit: unit, StockQuantity: stock, WarningThreshold: warning, Price: price, Remark: remark,
			})
			success++

			// 异步发送库存预警通知
			if h.ni != nil {
				go h.ni.CheckAndNotifyStockWarning(part.ID)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": fmt.Sprintf("导入完成，成功 %d 条，失败 %d 条", success, failed),
	})
}

// ExportParts 导出配件
func (h *PartHandler) ExportParts(c *gin.Context) {
	var parts []model.Part
	// 联表查询分类名称
	// 注意：Category 关联尚未明确主要在 model/models.go 中是否定义了 Preload 用的外键关系
	// 简单起见，这里假设 Category 存在。若不存在，此行可能报错，需手动处理
	// 查看 models.go 发现 Part 定义并没有 Category 字段做关联（只存了 CategoryID）
	// 因此这里使用 Joins 或手动 ID Map

	h.db.Find(&parts)
	// 获取所有分类
	var cats []model.PartCategory
	h.db.Find(&cats)
	catMap := make(map[uint]string)
	for _, c := range cats {
		catMap[c.ID] = c.Name
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// 表头
	headers := []string{"ID", "配件编号", "名称", "分类", "规格", "单位", "库存", "阈值", "单价", "备注", "创建时间"}
	for i, v := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue("Sheet1", cell, v)
	}

	// 填充数据
	for i, p := range parts {
		rowIdx := i + 2
		catName := catMap[p.CategoryID]
		vals := []interface{}{
			p.ID, p.PartNo, p.Name, catName, p.Specification, p.Unit,
			p.StockQuantity, p.WarningThreshold, p.Price, p.Remark, p.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		for j, v := range vals {
			cell, _ := excelize.CoordinatesToCellName(j+1, rowIdx)
			f.SetCellValue("Sheet1", cell, v)
		}
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=parts_export_%s.xlsx", time.Now().Format("200601021504")))

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
