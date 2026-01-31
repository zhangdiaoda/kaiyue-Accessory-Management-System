package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"warehouse/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type EmployeeHandler struct {
	db *gorm.DB
}

func NewEmployeeHandler(db *gorm.DB) *EmployeeHandler {
	return &EmployeeHandler{db: db}
}

// GetEmployeeList 获取员工列表
func (h *EmployeeHandler) GetEmployeeList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")
	status := c.Query("status") // 1:在职 0:离职

	var employees []model.Employee
	var total int64

	query := h.db.Model(&model.Employee{})

	if keyword != "" {
		query = query.Where("name LIKE ? OR employee_no LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&employees)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data": gin.H{
			"total":   total,
			"records": employees,
			"current": page,
			"size":    pageSize,
		},
	})
}

// CreateEmployee 创建员工
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var employee model.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	if err := h.db.Create(&employee).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    employee,
	})
}

// UpdateEmployee 更新员工
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee model.Employee

	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	if err := h.db.Model(&model.Employee{}).Where("id = ?", id).Updates(&employee).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
	})
}

// DeleteEmployee 删除员工
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	id := c.Param("id")

	if err := h.db.Delete(&model.Employee{}, id).Error; err != nil {
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

// GetAllEmployees 获取所有在职员工（用于下拉选择）
func (h *EmployeeHandler) GetAllEmployees(c *gin.Context) {
	var employees []model.Employee
	h.db.Where("status = 1").Order("name ASC").Find(&employees)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data":    employees,
	})
}

// DownloadTemplate 下载员工导入模板
func (h *EmployeeHandler) DownloadTemplate(c *gin.Context) {
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)

	f.SetCellValue(sheet, "A1", "工号")
	f.SetCellValue(sheet, "B1", "姓名")
	f.SetCellValue(sheet, "C1", "部门")
	f.SetCellValue(sheet, "D1", "职位")
	f.SetCellValue(sheet, "E1", "手机号")

	// 示例数据
	f.SetCellValue(sheet, "A2", "E1001")
	f.SetCellValue(sheet, "B2", "张三")
	f.SetCellValue(sheet, "C2", "生产部")
	f.SetCellValue(sheet, "D2", "操作工")
	f.SetCellValue(sheet, "E2", "13800138000")

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=employee_template.xlsx")
	f.Write(c.Writer)
}

// ImportEmployees 导入员工
func (h *EmployeeHandler) ImportEmployees(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "读取文件失败"})
		return
	}

	f, err := excelize.OpenReader(file)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "解析Excel失败"})
		return
	}

	rows, _ := f.GetRows("Sheet1")
	if len(rows) <= 1 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "模板为空或数据不足"})
		return
	}

	count := 0
	for i, row := range rows {
		if i == 0 || len(row) < 2 {
			continue
		}

		employee := model.Employee{
			EmployeeNo: row[0],
			Name:       row[1],
			Status:     1, // 默认在职
		}
		if len(row) > 2 {
			employee.Department = row[2]
		}
		if len(row) > 3 {
			employee.Position = row[3]
		}
		if len(row) > 4 {
			employee.Phone = row[4]
		}

		// 检查工号是否重复
		var exist model.Employee
		if err := h.db.Where("employee_no = ?", employee.EmployeeNo).First(&exist).Error; err == nil {
			// 如果存在则更新
			h.db.Model(&exist).Updates(employee)
		} else {
			h.db.Create(&employee)
		}
		count++
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": fmt.Sprintf("成功导入 %d 条记录", count)})
}

// ExportEmployees 导出员工
func (h *EmployeeHandler) ExportEmployees(c *gin.Context) {
	var employees []model.Employee
	h.db.Find(&employees)

	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"工号", "姓名", "部门", "职位", "手机号", "状态", "创建时间"}
	for i, head := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, head)
	}

	for i, e := range employees {
		rowIdx := i + 2
		status := "离职"
		if e.Status == 1 {
			status = "在职"
		}
		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowIdx), e.EmployeeNo)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", rowIdx), e.Name)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", rowIdx), e.Department)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", rowIdx), e.Position)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", rowIdx), e.Phone)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", rowIdx), status)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", rowIdx), e.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=employees_export.xlsx")
	f.Write(c.Writer)
}
