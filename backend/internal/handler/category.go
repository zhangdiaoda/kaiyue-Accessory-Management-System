package handler

import (
	"net/http"
	"warehouse/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	db *gorm.DB
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

// GetCategoryList 获取分类列表
func (h *CategoryHandler) GetCategoryList(c *gin.Context) {
	var categories []model.PartCategory
	h.db.Order("parent_id ASC, sort_order ASC").Find(&categories)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data":    categories,
	})
}

// CreateCategory 创建分类
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var category model.PartCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	if err := h.db.Create(&category).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    category,
	})
}

// UpdateCategory 更新分类
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category model.PartCategory

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	if err := h.db.Model(&model.PartCategory{}).Where("id = ?", id).Updates(&category).Error; err != nil {
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

// DeleteCategory 删除分类
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	// 检查是否有子分类
	var count int64
	h.db.Model(&model.PartCategory{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "该分类下有子分类，无法删除",
		})
		return
	}

	// 检查是否有配件使用该分类
	h.db.Model(&model.Part{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "该分类下有配件，无法删除",
		})
		return
	}

	if err := h.db.Delete(&model.PartCategory{}, id).Error; err != nil {
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
