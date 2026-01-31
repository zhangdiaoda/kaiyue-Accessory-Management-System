package model

import (
	"time"
)

// InboundRecord 入库/补货记录
type InboundRecord struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	PartID     uint      `gorm:"not null;index" json:"part_id"`
	Part       Part      `gorm:"foreignKey:PartID" json:"part"`
	Quantity   int       `gorm:"not null" json:"quantity"`         // 补货数量
	OperatorID uint      `json:"operator_id"`                      // 操作人ID (可选)
	Operator   string    `json:"operator"`                         // 操作人姓名 (冗余存储方便查询)
	BatchNo    string    `gorm:"type:varchar(50)" json:"batch_no"` // 批次号/采购单号
	Remark     string    `gorm:"type:text" json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
}

// TableName 指定表名
func (InboundRecord) TableName() string {
	return "inbound_record"
}
