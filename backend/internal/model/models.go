package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Username   string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password   string    `gorm:"size:100;not null" json:"-"` // -表示不序列化到JSON
	RealName   string    `gorm:"size:50;not null" json:"real_name"`
	Role       string    `gorm:"size:20;not null" json:"role"` // SUPER_ADMIN, WAREHOUSE_ADMIN
	Department string    `gorm:"size:50" json:"department"`
	Phone      string    `gorm:"size:20" json:"phone"`
	Status     int       `gorm:"default:1" json:"status"` // 0:禁用 1:启用
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "sys_user"
}

// Employee 员工模型
type Employee struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	EmployeeNo string    `gorm:"uniqueIndex;size:50;not null" json:"employee_no"`
	Name       string    `gorm:"size:50;not null" json:"name"`
	Department string    `gorm:"size:50" json:"department"`
	Position   string    `gorm:"size:50" json:"position"`
	Phone      string    `gorm:"size:20" json:"phone"`
	Status     int       `gorm:"default:1" json:"status"` // 0:离职 1:在职
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// PartCategory 配件分类模型
type PartCategory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	ParentID  uint      `gorm:"default:0" json:"parent_id"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Part 配件模型
type Part struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	PartNo           string    `gorm:"uniqueIndex;size:50;not null" json:"part_no"`
	Name             string    `gorm:"size:100;not null" json:"name"`
	CategoryID       uint      `gorm:"not null" json:"category_id"`
	Specification    string    `gorm:"size:200" json:"specification"`
	Unit             string    `gorm:"size:20;default:'件'" json:"unit"`
	StockQuantity    int       `gorm:"default:0" json:"stock_quantity"`
	WarningThreshold int       `gorm:"default:10" json:"warning_threshold"`
	Price            float64   `gorm:"type:decimal(10,2)" json:"price"`
	Remark           string    `gorm:"type:text" json:"remark"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// BorrowRecord 领用记录模型
type BorrowRecord struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	RecordNo        string     `gorm:"uniqueIndex;size:50;not null" json:"record_no"`
	EmployeeID      uint       `gorm:"not null" json:"employee_id"`
	PartID          uint       `gorm:"not null" json:"part_id"`
	BorrowQuantity  int        `gorm:"not null" json:"borrow_quantity"`
	ReturnQuantity  int        `gorm:"default:0" json:"return_quantity"`
	DamagedQuantity int        `gorm:"default:0" json:"damaged_quantity"`
	Status          string     `gorm:"size:20;not null" json:"status"` // BORROWED, RETURNED, PARTIAL_RETURNED
	BorrowTime      time.Time  `gorm:"not null" json:"borrow_time"`
	BorrowAdminID   uint       `gorm:"not null" json:"borrow_admin_id"`
	ReturnTime      *time.Time `json:"return_time,omitempty"`
	ReturnAdminID   *uint      `json:"return_admin_id,omitempty"`
	Remark          string     `gorm:"type:text" json:"remark"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// SysConfig 系统配置模型
type SysConfig struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ConfigKey   string    `gorm:"uniqueIndex;size:100;not null" json:"config_key"`
	ConfigValue string    `gorm:"type:text" json:"config_value"`
	Description string    `gorm:"size:200" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (SysConfig) TableName() string {
	return "sys_config"
}

// Announcement 系统公告模型
type Announcement struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:200;not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Type      string    `gorm:"size:20;not null" json:"type"` // POPUP, SCROLL
	Status    int       `gorm:"default:1" json:"status"`      // 0: 隐藏, 1: 显示
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Announcement) TableName() string {
	return "announcement"
}

// InternalMessage 站内信模型
type InternalMessage struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	SenderID   uint      `json:"sender_id"`
	ReceiverID uint      `json:"receiver_id"`
	Title      string    `gorm:"size:200;not null" json:"title"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	IsRead     bool      `gorm:"default:false" json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

func (InternalMessage) TableName() string {
	return "internal_message"
}

// OldPartInventory 旧件库模型 (存放回收的旧件)
type OldPartInventory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PartID     uint      `gorm:"not null" json:"part_id"`
	EmployeeID uint      `gorm:"not null" json:"employee_id"`
	Quantity   int       `gorm:"not null" json:"quantity"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (OldPartInventory) TableName() string {
	return "old_part_inventory"
}

// ScrapInventory 废品库模型 (存放报废的配件)
type ScrapInventory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PartID     uint      `gorm:"not null" json:"part_id"`
	EmployeeID uint      `gorm:"not null" json:"employee_id"`
	Quantity   int       `gorm:"not null" json:"quantity"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (ScrapInventory) TableName() string {
	return "scrap_inventory"
}

// Permission 权限模型
type Permission struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Category    string    `gorm:"size:50;not null;index" json:"category"`
	Description string    `gorm:"size:200" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Permission) TableName() string {
	return "sys_permission"
}

// UserPermission 用户权限关联模型
type UserPermission struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"not null;index" json:"user_id"`
	PermissionCode string    `gorm:"size:50;not null" json:"permission_code"`
	CreatedAt      time.Time `json:"created_at"`
}

func (UserPermission) TableName() string {
	return "sys_user_permission"
}

// RolePermission 角色默认权限模型
type RolePermission struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Role           string    `gorm:"size:50;not null;index" json:"role"`
	PermissionCode string    `gorm:"size:50;not null" json:"permission_code"`
	CreatedAt      time.Time `json:"created_at"`
}

func (RolePermission) TableName() string {
	return "sys_role_permission"
}

// OperationLog 操作日志模型
type OperationLog struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"not null;index" json:"user_id"`
	Username       string    `gorm:"size:50;not null" json:"username"`
	RealName       string    `gorm:"size:50;not null" json:"real_name"`
	Operation      string    `gorm:"size:50;not null;index" json:"operation"`
	Module         string    `gorm:"size:50;not null;index" json:"module"`
	Description    string    `gorm:"type:text" json:"description"`
	RequestMethod  string    `gorm:"size:10" json:"request_method"`
	RequestURL     string    `gorm:"size:500" json:"request_url"`
	RequestParams  string    `gorm:"type:text" json:"request_params"`
	ResponseResult string    `gorm:"type:text" json:"response_result"`
	IPAddress      string    `gorm:"size:50" json:"ip_address"`
	UserAgent      string    `gorm:"size:500" json:"user_agent"`
	Status         string    `gorm:"size:20;index" json:"status"` // SUCCESS, FAILED
	ErrorMessage   string    `gorm:"type:text" json:"error_message"`
	ExecutionTime  int       `json:"execution_time"` // 执行时长(ms)
	CreatedAt      time.Time `gorm:"index" json:"created_at"`
}

func (OperationLog) TableName() string {
	return "sys_operation_log"
}
