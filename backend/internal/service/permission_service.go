package service

import (
	"errors"
	"warehouse/internal/model"

	"gorm.io/gorm"
)

// PermissionService 权限服务
type PermissionService struct {
	db *gorm.DB
}

// NewPermissionService 创建权限服务
func NewPermissionService(db *gorm.DB) *PermissionService {
	return &PermissionService{db: db}
}

// GetAllPermissions 获取所有权限
func (s *PermissionService) GetAllPermissions() ([]model.Permission, error) {
	var permissions []model.Permission
	err := s.db.Order("category, id").Find(&permissions).Error
	return permissions, err
}

// GetPermissionsByCategory 按分类获取权限
func (s *PermissionService) GetPermissionsByCategory() (map[string][]model.Permission, error) {
	permissions, err := s.GetAllPermissions()
	if err != nil {
		return nil, err
	}

	result := make(map[string][]model.Permission)
	for _, perm := range permissions {
		result[perm.Category] = append(result[perm.Category], perm)
	}
	return result, nil
}

// GetUserPermissions 获取用户权限
// 返回用户的个人权限 + 角色默认权限
func (s *PermissionService) GetUserPermissions(userID uint, role string) ([]string, error) {
	// 超级管理员拥有所有权限
	if role == "SUPER_ADMIN" {
		var permissions []model.Permission
		if err := s.db.Find(&permissions).Error; err != nil {
			return nil, err
		}
		codes := make([]string, len(permissions))
		for i, p := range permissions {
			codes[i] = p.Code
		}
		return codes, nil
	}

	permSet := make(map[string]bool)

	// 1. 获取角色默认权限
	var rolePerms []model.RolePermission
	if err := s.db.Where("role = ?", role).Find(&rolePerms).Error; err != nil {
		return nil, err
	}
	for _, rp := range rolePerms {
		permSet[rp.PermissionCode] = true
	}

	// 2. 获取用户个人权限（扩展权限）
	var userPerms []model.UserPermission
	if err := s.db.Where("user_id = ?", userID).Find(&userPerms).Error; err != nil {
		return nil, err
	}
	for _, up := range userPerms {
		permSet[up.PermissionCode] = true
	}

	// 转换为切片
	permissions := make([]string, 0, len(permSet))
	for code := range permSet {
		permissions = append(permissions, code)
	}

	return permissions, nil
}

// GetRolePermissions 获取角色默认权限
func (s *PermissionService) GetRolePermissions(role string) ([]string, error) {
	var rolePerms []model.RolePermission
	if err := s.db.Where("role = ?", role).Find(&rolePerms).Error; err != nil {
		return nil, err
	}

	codes := make([]string, len(rolePerms))
	for i, rp := range rolePerms {
		codes[i] = rp.PermissionCode
	}
	return codes, nil
}

// SetUserPermissions 设置用户权限（个人扩展权限）
// 会覆盖现有的个人权限，不影响角色默认权限
func (s *PermissionService) SetUserPermissions(userID uint, permissionCodes []string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除现有权限
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserPermission{}).Error; err != nil {
			return err
		}

		// 添加新权限
		for _, code := range permissionCodes {
			userPerm := model.UserPermission{
				UserID:         userID,
				PermissionCode: code,
			}
			if err := tx.Create(&userPerm).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// SetRolePermissions 设置角色默认权限
func (s *PermissionService) SetRolePermissions(role string, permissionCodes []string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除现有权限
		if err := tx.Where("role = ?", role).Delete(&model.RolePermission{}).Error; err != nil {
			return err
		}

		// 添加新权限
		for _, code := range permissionCodes {
			rolePerm := model.RolePermission{
				Role:           role,
				PermissionCode: code,
			}
			if err := tx.Create(&rolePerm).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// CheckPermission 检查用户是否拥有指定权限
func (s *PermissionService) CheckPermission(userID uint, role string, permissionCode string) (bool, error) {
	// 超级管理员拥有所有权限
	if role == "SUPER_ADMIN" {
		return true, nil
	}

	// 检查角色默认权限
	var count int64
	if err := s.db.Model(&model.RolePermission{}).
		Where("role = ? AND permission_code = ?", role, permissionCode).
		Count(&count).Error; err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}

	// 检查用户个人权限
	if err := s.db.Model(&model.UserPermission{}).
		Where("user_id = ? AND permission_code = ?", userID, permissionCode).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetUserInfo 获取用户信息（用于权限检查）
func (s *PermissionService) GetUserInfo(userID uint) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}
